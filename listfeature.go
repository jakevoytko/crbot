package main

import (
	"bytes"
	"errors"
	"sort"
	"strings"

	"github.com/jakevoytko/crbot/api"
	"github.com/jakevoytko/crbot/feature"
	"github.com/jakevoytko/crbot/log"
	"github.com/jakevoytko/crbot/model"
)

// ListFeature is a Feature that lists commands that are available.
type ListFeature struct {
	featureRegistry *feature.Registry
	commandMap      model.StringMap
	gist            api.Gist
}

// NewListFeature returns a new ListFeature.
func NewListFeature(
	featureRegistry *feature.Registry,
	commandMap model.StringMap,
	gist api.Gist) *ListFeature {

	return &ListFeature{
		featureRegistry: featureRegistry,
		commandMap:      commandMap,
		gist:            gist,
	}
}

// Parsers returns the parsers.
func (f *ListFeature) Parsers() []feature.Parser {
	return []feature.Parser{NewListParser()}
}

// CommandInterceptors returns nothing.
func (f *ListFeature) CommandInterceptors() []feature.CommandInterceptor {
	return []feature.CommandInterceptor{}
}

// FallbackParser returns nil.
func (f *ListFeature) FallbackParser() feature.Parser {
	return nil
}

func (f *ListFeature) Executors() []feature.Executor {
	return []feature.Executor{NewListExecutor(f.featureRegistry, f.commandMap, f.gist)}
}

// ListParser parses ?list commands.
type ListParser struct{}

// NewListParser works as advertised.
func NewListParser() *ListParser {
	return &ListParser{}
}

// GetName returns the named type of this feature.
func (p *ListParser) GetName() string {
	return model.Name_List
}

const (
	MsgHelpList = "Type `?list` to get the URL of a Gist with all builtin and learned commands"
)

// HelpText explains how to use ?list.
func (p *ListParser) HelpText(command string) (string, error) {
	return MsgHelpList, nil
}

// Parse parses the given list command.
func (p *ListParser) Parse(splitContent []string) (*model.Command, error) {
	if splitContent[0] != p.GetName() {
		log.Fatal("parseList called with non-list command", errors.New("wat"))
	}
	return &model.Command{
		Type: model.Type_List,
	}, nil
}

const (
	MsgGistAddress  = "The list of commands is here"
	MsgListBuiltins = "List of builtins:"
	MsgListCustom   = "List of learned commands:"
)

type ListExecutor struct {
	featureRegistry *feature.Registry
	commandMap      model.StringMap
	gist            api.Gist
}

func NewListExecutor(featureRegistry *feature.Registry, commandMap model.StringMap, gist api.Gist) *ListExecutor {
	return &ListExecutor{
		featureRegistry: featureRegistry,
		commandMap:      commandMap,
		gist:            gist,
	}
}

// GetType returns the type of this feature.
func (e *ListExecutor) GetType() int {
	return model.Type_List
}

// Execute uploads the command list to github and pings the gist link in chat.
func (e *ListExecutor) Execute(s api.DiscordSession, channel string, command *model.Command) {
	builtins := e.featureRegistry.GetInvokableFeatureNames()
	all, err := e.commandMap.GetAll()
	if err != nil {
		log.Fatal("Error reading all commands", err)
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
		helpText, err := e.featureRegistry.GetParserByName(name).HelpText(name)

		// Log and continue if something goes wrong, to give the gist a chance of publishing.
		if err == nil {
			buffer.WriteString(": ")
			buffer.WriteString(helpText)
		} else {
			log.Info("Error getting builtin help text", err)
		}

		buffer.WriteString("\n")
	}

	buffer.WriteString("\n")

	buffer.WriteString(MsgListCustom)
	buffer.WriteString("\n")
	for _, name := range custom {
		buffer.WriteString(" - ?")
		buffer.WriteString(name)
		if strings.Contains(all[name], "$1") {
			buffer.WriteString(" <args>")
		}
		buffer.WriteString("\n")
	}

	url, err := e.gist.Upload(buffer.String())
	if err != nil {
		s.ChannelMessageSend(channel, err.Error())
		return
	}
	s.ChannelMessageSend(channel, MsgGistAddress+": "+url)
}
