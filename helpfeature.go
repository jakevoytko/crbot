package main

import "errors"

// HelpFeature is a Feature that prints a help prompt for the user.
type HelpFeature struct {
	featureRegistry *FeatureRegistry
}

// NewHelpFeature returns a new HelpFeature.
func NewHelpFeature(featureRegistry *FeatureRegistry) *HelpFeature {
	return &HelpFeature{
		featureRegistry: featureRegistry,
	}
}

// GetType returns the type.
func (f *HelpFeature) GetType() int {
	return Type_Help
}

// Parsers returns the parsers.
func (f *HelpFeature) Parsers() []Parser {
	return []Parser{
		NewHelpParser(f.featureRegistry),
	}
}

// FallbackParser returns nil.
func (f *HelpFeature) FallbackParser() Parser {
	return nil
}

// HelpParser parses ?help commands.
type HelpParser struct {
	featureRegistry *FeatureRegistry
}

// NewHelpParser works as advertised.
func NewHelpParser(featureRegistry *FeatureRegistry) *HelpParser {
	return &HelpParser{
		featureRegistry: featureRegistry,
	}
}

// GetName returns the named type.
func (p *HelpParser) GetName() string {
	return Name_Help
}

const (
	MsgHelpHelp = "Seems like you figured it out."
)

// GetHelpText returns the help text.
func (p *HelpParser) HelpText() string {
	return MsgHelpHelp
}

// Parse parses the given help command.
func (p *HelpParser) Parse(splitContent []string) (*Command, error) {
	if splitContent[0] != p.GetName() {
		fatal("parseHelp called with non-help command", errors.New("wat"))
	}
	userCommand := ""
	if len(splitContent) > 1 {
		if p.featureRegistry.IsInvokable(splitContent[1]) {
			userCommand = splitContent[1]
		}
	}

	return &Command{
		Type: Type_Help,
		Help: &HelpData{
			Command: userCommand,
		},
	}, nil
}

const (
	MsgDefaultHelp = "Type `?help` for this message, `?list` to list all commands, or `?help <command>` to get help for a particular command."
)

// Execute replies over the given channel with a help message.
func (f *HelpFeature) Execute(s DiscordSession, channel string, command *Command) {
	if command.Help == nil {
		fatal("Incorrectly generated help command", errors.New("wat"))
	}

	parser := f.featureRegistry.GetParserByName(command.Help.Command)
	if parser != nil {
		if _, err := s.ChannelMessageSend(channel, parser.HelpText()); err != nil {
			info("Failed to send default help message", err)
		}
	} else {
		if _, err := s.ChannelMessageSend(channel, MsgDefaultHelp); err != nil {
			info("Failed to send default help message", err)
		}
	}
}
