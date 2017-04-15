package main

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// LearnFeature allows crbot to learn new calls and responses
type LearnFeature struct {
	featureRegistry *FeatureRegistry
	commandMap      StringMap
}

// NewLearnFeature returns a new LearnFeature.
func NewLearnFeature(featureRegistry *FeatureRegistry, commandMap StringMap) *LearnFeature {
	return &LearnFeature{
		featureRegistry: featureRegistry,
		commandMap:      commandMap,
	}
}

// GetType returns the type of this feature.
func (f *LearnFeature) GetType() int {
	return Type_Learn
}

// Parsers gets the learn feature parsers.
func (f *LearnFeature) Parsers() []Parser {
	return []Parser{
		NewLearnParser(f.featureRegistry, f.commandMap),
	}
}

// FallbackParser returns nil.
func (f *LearnFeature) FallbackParser() Parser {
	return nil
}

// LearnParser parses ?learn commands.
type LearnParser struct {
	featureRegistry *FeatureRegistry
	commandMap      StringMap
}

// NewLearnParser works as advertised.
func NewLearnParser(featureRegistry *FeatureRegistry, commandMap StringMap) *LearnParser {
	return &LearnParser{
		featureRegistry: featureRegistry,
		commandMap:      commandMap,
	}
}

// GetName returns the named type of this feature.
func (p *LearnParser) GetName() string {
	return Name_Learn
}

const (
	MsgHelpLearn = "Type `?learn <call> <the response the bot should read>`. When you type `?call`, the bot will reply with the response.\n\nThe first character of the call must be alphanumeric, and the first character of the response must not begin with /, ?, or !\n\nUse $1 in the response to substitute all arguments"
)

// HelpText explains how to use ?learn.
func (p *LearnParser) HelpText() string {
	return MsgHelpLearn
}

// Parse parses the given learn command.
func (f *LearnParser) Parse(splitContent []string) (*Command, error) {
	if splitContent[0] != f.GetName() {
		fatal("parseLearn called with non-learn command", errors.New("wat"))
	}

	callRegexp := regexp.MustCompile("^[[:alnum:]].*$")
	responseRegexp := regexp.MustCompile("(?s)^[^/?!].*$")

	// Show help when not enough data is present, or malicious data is present.
	if len(splitContent) < 3 || !callRegexp.MatchString(splitContent[1]) || !responseRegexp.MatchString(splitContent[2]) {
		return &Command{
			Type: Type_Help,
			Help: &HelpData{
				Command: Name_Learn,
			},
		}, nil
	}

	// Don't overwrite old or builtin commands.
	has, err := f.commandMap.Has(splitContent[1])
	if err != nil {
		return nil, err
	}
	if has || f.featureRegistry.IsInvokable(splitContent[1]) {
		return &Command{
			Type: Type_Learn,
			Learn: &LearnData{
				CallOpen: false,
				Call:     splitContent[1],
			},
		}, nil
	}

	// Everything is good.
	response := strings.Join(splitContent[2:], " ")
	return &Command{
		Type: Type_Learn,
		Learn: &LearnData{
			CallOpen: true,
			Call:     splitContent[1],
			Response: response,
		},
	}, nil
}

const (
	MsgLearnFail    = "I already know ?%s"
	MsgLearnSuccess = "Learned about %s"
)

// Execute replies over the given channel with a help message.
func (f *LearnFeature) Execute(s DiscordSession, channel string, command *Command) {
	if command.Learn == nil {
		fatal("Incorrectly generated learn command", errors.New("wat"))
	}
	if !command.Learn.CallOpen {
		s.ChannelMessageSend(channel, fmt.Sprintf(MsgLearnFail, command.Learn.Call))
		return
	}

	// Teach the command.
	if has, err := f.commandMap.Has(command.Learn.Call); err != nil || has {
		if has {
			fatal("Collision when adding a call for "+command.Learn.Call, errors.New("wat"))
		}
		fatal("Error in LearnFeature#Execute, testing a command", err)
	}
	if err := f.commandMap.Set(command.Learn.Call, command.Learn.Response); err != nil {
		fatal("Error storing a learn command. Dying since it might work with restart", err)
	}

	// Send ack.
	s.ChannelMessageSend(channel, fmt.Sprintf(MsgLearnSuccess, command.Learn.Call))
}
