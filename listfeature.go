package main

import (
	"bytes"
	"errors"
	"sort"
)

// ListFeature is a Feature that lists commands that are available.
type ListFeature struct {
	featureRegistry *FeatureRegistry
	commandMap      StringMap
	gist            Gist
}

// NewListFeature returns a new ListFeature.
func NewListFeature(
	featureRegistry *FeatureRegistry,
	commandMap StringMap,
	gist Gist) *ListFeature {

	return &ListFeature{
		featureRegistry: featureRegistry,
		commandMap:      commandMap,
		gist:            gist,
	}
}

// Parsers returns the parsers.
func (f *ListFeature) Parsers() []Parser {
	return []Parser{NewListParser()}
}

// FallbackParser returns nil.
func (f *ListFeature) FallbackParser() Parser {
	return nil
}

func (f *ListFeature) Executors() []Executor {
	return []Executor{NewListExecutor(f.featureRegistry, f.commandMap, f.gist)}
}

// ListParser parses ?list commands.
type ListParser struct{}

// NewListParser works as advertised.
func NewListParser() *ListParser {
	return &ListParser{}
}

// GetName returns the named type of this feature.
func (p *ListParser) GetName() string {
	return Name_List
}

const (
	MsgHelpList = "Type `?list` to get the URL of a Gist with all builtin and learned commands"
)

// HelpText explains how to use ?list.
func (p *ListParser) HelpText() string {
	return MsgHelpList
}

// Parse parses the given list command.
func (p *ListParser) Parse(splitContent []string) (*Command, error) {
	if splitContent[0] != p.GetName() {
		fatal("parseList called with non-list command", errors.New("wat"))
	}
	return &Command{
		Type: Type_List,
	}, nil
}

const (
	MsgGistAddress       = "The list of commands is here"
	MsgGistPostFail      = "Unable to connect to Gist service. Give it a few minutes and try again"
	MsgGistResponseFail  = "Failure reading response from Gist service"
	MsgGistSerializeFail = "Unable to serialize Gist"
	MsgGistStatusCode    = "Failed to upload Gist :("
	MsgGistUrlFail       = "Failed getting url from Gist service"
	MsgListBuiltins      = "List of builtins:"
	MsgListCustom        = "List of learned commands:"
)

type ListExecutor struct {
	featureRegistry *FeatureRegistry
	commandMap      StringMap
	gist            Gist
}

func NewListExecutor(featureRegistry *FeatureRegistry, commandMap StringMap, gist Gist) *ListExecutor {
	return &ListExecutor{
		featureRegistry: featureRegistry,
		commandMap:      commandMap,
		gist:            gist,
	}
}

// GetType returns the type of this feature.
func (e *ListExecutor) GetType() int {
	return Type_List
}

// Execute uploads the command list to github and pings the gist link in chat.
func (e *ListExecutor) Execute(s DiscordSession, channel string, command *Command) {
	builtins := e.featureRegistry.GetInvokableFeatureNames()
	all, err := e.commandMap.GetAll()
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

	url, err := e.gist.Upload(buffer.String())
	if err != nil {
		s.ChannelMessageSend(channel, err.Error())
		return
	}
	s.ChannelMessageSend(channel, MsgGistAddress+": "+url)
}
