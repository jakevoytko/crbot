package main

import (
	"errors"

	"github.com/jakevoytko/crbot/api"
	"github.com/jakevoytko/crbot/feature"
	"github.com/jakevoytko/crbot/log"
	"github.com/jakevoytko/crbot/model"
)

// HelpFeature is a Feature that prints a help prompt for the user.
type HelpFeature struct {
	featureRegistry *feature.Registry
}

// NewHelpFeature returns a new HelpFeature.
func NewHelpFeature(featureRegistry *feature.Registry) *HelpFeature {
	return &HelpFeature{
		featureRegistry: featureRegistry,
	}
}

// Parsers returns the parsers.
func (f *HelpFeature) Parsers() []feature.Parser {
	return []feature.Parser{NewHelpParser(f.featureRegistry)}
}

// CommandInterceptors returns nothing.
func (f *HelpFeature) CommandInterceptors() []feature.CommandInterceptor {
	return []feature.CommandInterceptor{}
}

// FallbackParser returns nil.
func (f *HelpFeature) FallbackParser() feature.Parser {
	return nil
}

// Executors gets the executors.
func (f *HelpFeature) Executors() []feature.Executor {
	return []feature.Executor{NewHelpExecutor(f.featureRegistry)}
}

// HelpParser parses ?help commands.
type HelpParser struct {
	featureRegistry *feature.Registry
}

// NewHelpParser works as advertised.
func NewHelpParser(featureRegistry *feature.Registry) *HelpParser {
	return &HelpParser{
		featureRegistry: featureRegistry,
	}
}

// GetName returns the named type.
func (p *HelpParser) GetName() string {
	return model.Name_Help
}

const (
	MsgHelpHelp = "Get help for help."
)

// GetHelpText returns the help text.
func (p *HelpParser) HelpText(command string) (string, error) {
	return MsgHelpHelp, nil
}

// Parse parses the given help command.
func (p *HelpParser) Parse(splitContent []string) (*model.Command, error) {
	if splitContent[0] != p.GetName() {
		log.Fatal("parseHelp called with non-help command", errors.New("wat"))
	}
	splitContent = CollapseWhitespace(splitContent, 1)

	userCommand := ""
	if len(splitContent) > 1 {
		userCommand = splitContent[1]
	}

	return &model.Command{
		Type: model.Type_Help,
		Help: &model.HelpData{
			Command: userCommand,
		},
	}, nil
}

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
func (e *HelpExecutor) Execute(s api.DiscordSession, channel string, command *model.Command) {
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
		if _, err := s.ChannelMessageSend(channel, helpText); err != nil {
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

	if _, err := s.ChannelMessageSend(channel, fallbackHelpText); err != nil {
		log.Info("Failed to send fallback or default help message", err)
	}
}
