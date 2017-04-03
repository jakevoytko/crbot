package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"sort"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

// ListFeature is a Feature that lists commands that are available.
type ListFeature struct {
	featureRegistry *FeatureRegistry
	commandMap      StringMap
}

// NewListFeature returns a new ListFeature.
func NewListFeature(featureRegistry *FeatureRegistry, commandMap StringMap) *ListFeature {
	return &ListFeature{
		featureRegistry: featureRegistry,
		commandMap:      commandMap,
	}
}

// GetName returns the named type of this feature.
func (f *ListFeature) GetName() string {
	return Name_List
}

// GetType returns the type of this feature.
func (f *ListFeature) GetType() int {
	return Type_List
}

// Invokable indicates whether the user can invoke this feature by name.
func (f *ListFeature) Invokable() bool {
	return true
}

// Parse parses the given list command.
func (f *ListFeature) Parse(splitContent []string) (*Command, error) {
	if splitContent[0] != f.GetName() {
		fatal("parseList called with non-list command", errors.New("wat"))
	}
	return &Command{
		Type: Type_List,
	}, nil
}

const (
	MsgGistAddress       = "The list of commands is here: "
	MsgGistPostFail      = "Unable to connect to Gist service. Give it a few minutes and try again"
	MsgGistResponseFail  = "Failure reading response from Gist service"
	MsgGistSerializeFail = "Unable to serialize Gist"
	MsgGistStatusCode    = "Failed to upload Gist :("
	MsgGistUrlFail       = "Failed getting url from Gist service"
	MsgListBuiltins      = "List of builtins:"
	MsgListCustom        = "List of learned commands:"
)

// Execute uploads the command list to github and pings the gist link in chat.
func (f *ListFeature) Execute(s *discordgo.Session, channel string, command *Command) {
	builtins := f.featureRegistry.GetInvokableFeatureNames()
	all, err := f.commandMap.GetAll()
	if err != nil {
		fatal("Error reading all commands", err)
	}
	custom := make([]string, 0, len(all))
	for name := range all {
		custom = append(custom, name)
	}

	sort.Strings(builtins)
	sort.Strings(custom)

	var buffer bytes.Buffer
	buffer.WriteString(MsgListBuiltins)
	buffer.WriteString("\n")
	for _, name := range builtins {
		buffer.WriteString(" - ")
		buffer.WriteString(name)
		buffer.WriteString("\n")
	}

	buffer.WriteString("\n")

	buffer.WriteString(MsgListCustom)
	buffer.WriteString("\n")
	for _, name := range custom {
		buffer.WriteString(" - ?")
		buffer.WriteString(name)
		buffer.WriteString("\n")
	}

	url, err := uploadCommandList(buffer.String())
	if err != nil {
		s.ChannelMessageSend(channel, err.Error())
		return
	}
	s.ChannelMessageSend(channel, MsgGistAddress+": "+url)
}

///////////////////////////////////////////////////////////////////////////////
// Gist handling
///////////////////////////////////////////////////////////////////////////////
type Gist struct {
	Description string           `json:"description"`
	Public      bool             `json:"public"`
	Files       map[string]*File `json:"files"`
}

// A file represents the contents of a Gist.
type File struct {
	Content string `json:"content"`
}

// simpleGist returns a Gist object with just the given contents.
func simpleGist(contents string) *Gist {
	return &Gist{
		Public:      false,
		Description: "CRBot command list",
		Files: map[string]*File{
			"commands": &File{
				Content: contents,
			},
		},
	}
}

func uploadCommandList(contents string) (string, error) {
	gist := simpleGist(contents)
	serializedGist, err := json.Marshal(gist)
	if err != nil {
		info("Error marshalling gist", err)
		return "", errors.New(MsgGistSerializeFail)
	}
	response, err := http.Post(
		"https://api.github.com/gists", "application/json", bytes.NewBuffer(serializedGist))
	if err != nil {
		info("Error POSTing gist", err)
		return "", errors.New(MsgGistPostFail)
	} else if response.StatusCode != 201 {
		info("Bad status code", errors.New("Code: "+strconv.Itoa(response.StatusCode)))
		return "", errors.New(MsgGistStatusCode)
	}

	responseMap := map[string]interface{}{}
	if err := json.NewDecoder(response.Body).Decode(&responseMap); err != nil {
		info("Error reading gist response", err)
		return "", errors.New(MsgGistResponseFail)
	}

	if finalUrl, ok := responseMap["html_url"]; ok {
		return finalUrl.(string), nil
	}
	return "", errors.New(MsgGistUrlFail)
}
