package list

import (
	"bytes"
	"sort"
	"strings"

	"github.com/jakevoytko/crbot/api"
	"github.com/jakevoytko/crbot/feature"
	"github.com/jakevoytko/crbot/log"
	"github.com/jakevoytko/crbot/model"
)

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

// PublicOnly returns whether the executor should be intercepted in a private channel.
func (e *ListExecutor) PublicOnly() bool {
	return false
}

// Execute uploads the command list to github and pings the gist link in chat.
func (e *ListExecutor) Execute(s api.DiscordSession, channel model.Snowflake, command *model.Command) {
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
		s.ChannelMessageSend(channel.Format(), err.Error())
		return
	}
	s.ChannelMessageSend(channel.Format(), MsgGistAddress+": "+url)
}
