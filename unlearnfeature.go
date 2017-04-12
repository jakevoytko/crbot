package main

import (
	"errors"
	"fmt"
	"regexp"
)

// UnlearnFeature allows crbot to unlearn existing calls
type UnlearnFeature struct {
	featureRegistry *FeatureRegistry
	commandMap      StringMap
}

// NewUnlearnFeature returns a new UnlearnFeature.
func NewUnlearnFeature(featureRegistry *FeatureRegistry, commandMap StringMap) *UnlearnFeature {
	return &UnlearnFeature{
		featureRegistry: featureRegistry,
		commandMap:      commandMap,
	}
}

// GetName returns the named type of this feature.
func (f *UnlearnFeature) GetName() string {
	return Name_Unlearn
}

// GetType returns the type of this feature.
func (f *UnlearnFeature) GetType() int {
	return Type_Unlearn
}

// Invokable returns whether the user can execute this command by name.
func (f *UnlearnFeature) Invokable() bool {
	return true
}

// Parse parses the given unlearn command.
func (f *UnlearnFeature) Parse(splitContent []string) (*Command, error) {
	if splitContent[0] != f.GetName() {
		fatal("parseUnlearn called with non-unlearn command", errors.New("wat"))
	}

	callRegexp := regexp.MustCompile("^[[:alnum:]].*$")

	// Show help when not enough data is present, or malicious data is present.
	if len(splitContent) < 2 || !callRegexp.MatchString(splitContent[1]) {
		return &Command{
			Type: Type_Help,
			Help: &HelpData{
				Type: Type_Unlearn,
			},
		}, nil
	}

	// Only unlearn commands that aren't built-in and exist
	has, err := f.commandMap.Has(splitContent[1])
	if err != nil {
		return nil, err
	}
	if !has || f.featureRegistry.GetTypeFromName(splitContent[1]) != Type_None {
		return &Command{
			Type: Type_Unlearn,
			Unlearn: &UnlearnData{
				CallOpen: false,
				Call:     splitContent[1],
			},
		}, nil
	}

	// Everything is good.
	return &Command{
		Type: Type_Unlearn,
		Unlearn: &UnlearnData{
			CallOpen: true,
			Call:     splitContent[1],
		},
	}, nil
}

const (
	MsgUnlearnFail         = "I can't unlearn `?%s`"
	MsgUnlearnMustBePublic = "I can't unlearn in a private message."
	MsgUnlearnSuccess      = "Forgot about %s"
)

// Execute replies over the given channel indicating successful unlearning, or
// failure to unlearn.
func (f *UnlearnFeature) Execute(s DiscordSession, channel string, command *Command) {
	if command.Unlearn == nil {
		fatal("Incorrectly generated unlearn command", errors.New("wat"))
	}

	// Get the current channel and check if we're being asked to unlearn in a
	// private message.
	discordChannel, err := s.Channel(channel)
	if err != nil {
		fatal("This message didn't come from a valid channel", errors.New("wat"))
	}
	if discordChannel.IsPrivate {
		s.ChannelMessageSend(channel, MsgUnlearnMustBePublic)
		return
	}

	if !command.Unlearn.CallOpen {
		s.ChannelMessageSend(channel, fmt.Sprintf(MsgUnlearnFail, command.Unlearn.Call))
		return
	}

	// Remove the command.
	if has, err := f.commandMap.Has(command.Unlearn.Call); !has || err != nil {
		if has {
			fatal("Tried to unlearn command that doesn't exist: "+command.Unlearn.Call, errors.New("wat"))
		}
		fatal("Error in UnlearnFeature#execute, testing a command", err)
	}
	if err := f.commandMap.Delete(command.Unlearn.Call); err != nil {
		fatal("Unsuccessful unlearning a key; Dying since it might work with a restart", err)
	}

	// Send ack.
	s.ChannelMessageSend(channel, fmt.Sprintf(MsgUnlearnSuccess, command.Unlearn.Call))
}
