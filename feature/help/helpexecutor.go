package help

import (
	"errors"

	"github.com/jakevoytko/crbot/api"
	"github.com/jakevoytko/crbot/feature"
	"github.com/jakevoytko/crbot/log"
	"github.com/jakevoytko/crbot/model"
)

const (
	// MsgDefaultHelp is what is shown when ?help is entered with no arguments.
	MsgDefaultHelp = "Type `?help` for this message, `?list` to list all commands, or `?help <command>` to get help for a particular command."
)

// Executor prints help text for help commands.
type Executor struct {
	featureRegistry *feature.Registry
}

// NewExecutor works as advertised.
func NewExecutor(featureRegistry *feature.Registry) *Executor {
	return &Executor{
		featureRegistry: featureRegistry,
	}
}

// GetType returns the type.
func (e *Executor) GetType() int {
	return model.CommandTypeHelp
}

// PublicOnly returns whether the executor should be intercepted in a private channel.
func (e *Executor) PublicOnly() bool {
	return false
}

// Execute replies over the given channel with a help message.
func (e *Executor) Execute(s api.DiscordSession, channel model.Snowflake, command *model.Command) {
	if command.Help == nil {
		log.Fatal("Incorrectly generated help command", errors.New("wat"))
	}

	// Try to parse custom commands before fallback commands. This matches actual
	// command invocation order, and will avoid returning help text when there
	// happens to be a custom command that was overwritten by a new builtin.
	parser := e.featureRegistry.GetParserByName(command.Help.Command)
	if parser != nil {
		// Use the builtin parsers to generate help text.
		helpText, err := parser.HelpText(command.Help.Command)
		if err != nil {
			log.Info("Failed to generate help text for command "+command.Help.Command, err)
			return
		}
		if _, err := s.ChannelMessageSend(channel.Format(), helpText); err != nil {
			log.Info("Failed to send default help message", err)
		}
		return
	}

	// Use the fallback parser to generate help text.
	fallbackHelpText, err := e.featureRegistry.FallbackParser.HelpText(command.Help.Command)
	if err != nil {
		log.Info("Failed to generate backup help text for command "+command.Help.Command, err)
		return
	}
	if fallbackHelpText == "" {
		fallbackHelpText = MsgDefaultHelp
	}

	if _, err := s.ChannelMessageSend(channel.Format(), fallbackHelpText); err != nil {
		log.Info("Failed to send fallback or default help message", err)
	}
}
