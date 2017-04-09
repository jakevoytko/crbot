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

// GetName returns the named type.
func (f *HelpFeature) GetName() string {
	return Name_Help
}

// GetType returns the type.
func (f *HelpFeature) GetType() int {
	return Type_Help
}

// Invokable returns whether the user can execute this by the command line.
func (f *HelpFeature) Invokable() bool {
	return true
}

// Parse parses the given help command.
func (f *HelpFeature) Parse(splitContent []string) (*Command, error) {
	if splitContent[0] != f.GetName() {
		fatal("parseHelp called with non-help command", errors.New("wat"))
	}
	userType := Type_Unrecognized
	if len(splitContent) > 1 {
		userType = f.featureRegistry.GetTypeFromName(splitContent[1])
	}

	return &Command{
		Type: Type_Help,
		Help: &HelpData{
			Type: userType,
		},
	}, nil
}

const (
	MsgDefaultHelp = "Type `?help` for this message, `?list` to list all commands, or `?help <command>` to get help for a particular command."
	MsgHelpHelp    = "You're probably right. I probably didn't think of this case."
	MsgHelpLearn   = "Type `?learn <call> <the response the bot should read>`. When you type `?call`, the bot will reply with the response.\n\nThe first character of the call must be alphanumeric, and the first character of the response must not begin with /, ?, or !\n\nUse $1 in the response to substitute all arguments"
	MsgHelpList    = "Type `?list` to get the URL of a Gist with all builtin and learned commands"
	MsgHelpUnlearn = "Type `?unlearn <call>` to forget a user-defined command."
)

// Execute replies over the given channel with a help message.
func (f *HelpFeature) Execute(s DiscordSession, channel string, command *Command) {
	if command.Help == nil {
		fatal("Incorrectly generated help command", errors.New("wat"))
	}

	// TODO(jake): Add a help message registry that's shared between the feature
	// registry and the help feature to remove this feature leak.
	switch command.Help.Type {
	default:
		if _, err := s.ChannelMessageSend(channel, MsgDefaultHelp); err != nil {
			info("Failed to send default help message", err)
		}
	case Type_Help:
		if _, err := s.ChannelMessageSend(channel, MsgHelpHelp); err != nil {
			info("Failed to send help help message", err)
		}
	case Type_Learn:
		if _, err := s.ChannelMessageSend(channel, MsgHelpLearn); err != nil {
			info("Failed to send learn help message", err)
		}
	case Type_List:
		if _, err := s.ChannelMessageSend(channel, MsgHelpList); err != nil {
			info("Failed to send list help message", err)
		}
	case Type_Unlearn:
		if _, err := s.ChannelMessageSend(channel, MsgHelpUnlearn); err != nil {
			info("Failed to send unlearn help message", err)
		}
	}
}
