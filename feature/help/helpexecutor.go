package help

import (
	"errors"

	"github.com/jakevoytko/crbot/api"
	"github.com/jakevoytko/crbot/feature"
	"github.com/jakevoytko/crbot/log"
	"github.com/jakevoytko/crbot/model"
)

const (
	MsgDefaultHelp = "Type `?help` for this message, `?list` to list all commands, or `?help <command>` to get help for a particular command."
)

// HelpExecutor prints help text for help commands.
type HelpExecutor struct {
	featureRegistry *feature.Registry
}

// NewHelpExecutor works as advertised.
func NewHelpExecutor(featureRegistry *feature.Registry) *HelpExecutor {
	return &HelpExecutor{
		featureRegistry: featureRegistry,
	}
}

// GetType returns the type.
func (e *HelpExecutor) GetType() int {
	return model.Type_Help
}

// Execute replies over the given channel with a help message.
func (e *HelpExecutor) Execute(s api.DiscordSession, channel model.Snowflake, command *model.Command) {
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
