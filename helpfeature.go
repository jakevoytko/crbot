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

// Parsers returns the parsers.
func (f *HelpFeature) Parsers() []Parser {
	return []Parser{NewHelpParser(f.featureRegistry)}
}

// FallbackParser returns nil.
func (f *HelpFeature) FallbackParser() Parser {
	return nil
}

// Executors gets the executors.
func (f *HelpFeature) Executors() []Executor {
	return []Executor{NewHelpExecutor(f.featureRegistry)}
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
	MsgHelpHelp = "Get help for help."
)

// GetHelpText returns the help text.
func (p *HelpParser) HelpText(command string) (string, error) {
	return MsgHelpHelp, nil
}

// Parse parses the given help command.
func (p *HelpParser) Parse(splitContent []string) (*Command, error) {
	if splitContent[0] != p.GetName() {
		fatal("parseHelp called with non-help command", errors.New("wat"))
	}
	splitContent = CollapseWhitespace(splitContent, 1)

	userCommand := ""
	if len(splitContent) > 1 {
		userCommand = splitContent[1]
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

// HelpExecutor prints help text for help commands.
type HelpExecutor struct {
	featureRegistry *FeatureRegistry
}

// NewHelpExecutor works as advertised.
func NewHelpExecutor(featureRegistry *FeatureRegistry) *HelpExecutor {
	return &HelpExecutor{
		featureRegistry: featureRegistry,
	}
}

// GetType returns the type.
func (e *HelpExecutor) GetType() int {
	return Type_Help
}

// Execute replies over the given channel with a help message.
func (e *HelpExecutor) Execute(s DiscordSession, channel string, command *Command) {
	if command.Help == nil {
		fatal("Incorrectly generated help command", errors.New("wat"))
	}

	// Try to parse custom commands before fallback commands. This matches actual
	// command invocation order, and will avoid returning help text when there
	// happens to be a custom command that was overwritten by a new builtin.
	parser := e.featureRegistry.GetParserByName(command.Help.Command)
	if parser != nil {
		// Use the builtin parsers to generate help text.
		helpText, err := parser.HelpText(command.Help.Command)
		if err != nil {
			info("Failed to generate help text for command "+command.Help.Command, err)
			return
		}
		if _, err := s.ChannelMessageSend(channel, helpText); err != nil {
			info("Failed to send default help message", err)
		}
		return
	}

	// Use the fallback parser to generate help text.
	fallbackHelpText, err := e.featureRegistry.FallbackParser.HelpText(command.Help.Command)
	if err != nil {
		info("Failed to generate backup help text for command "+command.Help.Command, err)
		return
	}
	if fallbackHelpText == "" {
		fallbackHelpText = MsgDefaultHelp
	}

	if _, err := s.ChannelMessageSend(channel, fallbackHelpText); err != nil {
		info("Failed to send fallback or default help message", err)
	}
}
