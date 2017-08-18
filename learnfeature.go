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

// Parsers gets the learn feature parsers.
func (f *LearnFeature) Parsers() []Parser {
	return []Parser{
		NewLearnParser(f.featureRegistry, f.commandMap),
		NewUnlearnParser(f.featureRegistry, f.commandMap),
	}
}

// CommandInterceptors returns nothing.
func (f *LearnFeature) CommandInterceptors() []CommandInterceptor {
	return []CommandInterceptor{}
}

// FallbackParser returns the custom parser, to recognize custom ? commands. It
// should be the only fallback parser in the project.
func (f *LearnFeature) FallbackParser() Parser {
	return NewCustomParser(f.commandMap)
}

func (f *LearnFeature) Executors() []Executor {
	return []Executor{
		NewLearnExecutor(f.commandMap),
		NewUnlearnExecutor(f.commandMap),
		NewCustomExecutor(f.commandMap),
	}
}

///////////////////////////////////////////////////////////////////////////////
// Messages
///////////////////////////////////////////////////////////////////////////////

const (
	MsgCustomNeedsArgs     = "This command takes args. Please type `?command <more text>` instead of `?command`"
	MsgHelpLearn           = "Type `?learn <call> <the response the bot should read>`. When you type `?call`, the bot will reply with the response.\n\nThe first character of the call must be alphanumeric, and the first character of the response must not begin with /, ?, or !\n\nUse $1 in the response to substitute all arguments"
	MsgHelpUnlearn         = "Type `?unlearn <call>` to forget a user-defined command."
	MsgLearnFail           = "I already know ?%s"
	MsgLearnSuccess        = "Learned about %s"
	MsgUnlearnFail         = "I can't unlearn `?%s`"
	MsgUnlearnMustBePublic = "I can't unlearn in a private message."
	MsgUnlearnSuccess      = "Forgot about %s"
)

///////////////////////////////////////////////////////////////////////////////
// Parsers
///////////////////////////////////////////////////////////////////////////////

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

// HelpText explains how to use ?learn.
func (p *LearnParser) HelpText(command string) (string, error) {
	return MsgHelpLearn, nil
}

// Parse parses the given learn command.
func (f *LearnParser) Parse(splitContent []string) (*Command, error) {
	if splitContent[0] != f.GetName() {
		fatal("parseLearn called with non-learn command", errors.New("wat"))
	}
	splitContent = CollapseWhitespace(splitContent, 1)
	splitContent = CollapseWhitespace(splitContent, 2)

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

// UnlearnParser parses ?unlearn commands.
type UnlearnParser struct {
	featureRegistry *FeatureRegistry
	commandMap      StringMap
}

// NewUnlearnParser works as advertised.
func NewUnlearnParser(featureRegistry *FeatureRegistry, commandMap StringMap) *UnlearnParser {
	return &UnlearnParser{
		featureRegistry: featureRegistry,
		commandMap:      commandMap,
	}
}

// GetName returns the named type of this feature.
func (p *UnlearnParser) GetName() string {
	return Name_Unlearn
}

// HelpText returns the help text for ?unlearn.
func (p *UnlearnParser) HelpText(command string) (string, error) {
	return MsgHelpUnlearn, nil
}

// Parse parses the given unlearn command.
func (p *UnlearnParser) Parse(splitContent []string) (*Command, error) {
	if splitContent[0] != p.GetName() {
		fatal("parseUnlearn called with non-unlearn command", errors.New("wat"))
	}

	splitContent = CollapseWhitespace(splitContent, 1)

	callRegexp := regexp.MustCompile("^[[:alnum:]].*$")

	// Show help when not enough data is present, or malicious data is present.
	if len(splitContent) < 2 || !callRegexp.MatchString(splitContent[1]) {
		return &Command{
			Type: Type_Help,
			Help: &HelpData{
				Command: Name_Unlearn,
			},
		}, nil
	}

	// Only unlearn commands that aren't built-in and exist
	has, err := p.commandMap.Has(splitContent[1])
	if err != nil {
		return nil, err
	}
	if !has || p.featureRegistry.IsInvokable(splitContent[1]) {
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

// CustomParser parses all fallthrough commands.
type CustomParser struct {
	commandMap StringMap
}

// NewCustomParser works as advertised.
func NewCustomParser(commandMap StringMap) *CustomParser {
	return &CustomParser{
		commandMap: commandMap,
	}
}

// GetName returns nothing, since it doesn't have a user-invokable name.
func (p *CustomParser) GetName() string {
	return ""
}

// HelpText returns help text for the given custom command.
func (p *CustomParser) HelpText(command string) (string, error) {
	ok, err := p.commandMap.Has(command)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", nil
	}

	val, err := p.commandMap.Get(command)
	if err != nil {
		return "", err
	}
	response := "?" + command
	if strings.Contains(val, "$1") {
		response = response + " <args>"
	}
	return response, nil
}

// Parse parses the given custom command.
func (f *CustomParser) Parse(splitContent []string) (*Command, error) {
	// TODO(jake): Drop this and external hash check, handle missing commands solely in execute.
	has, err := f.commandMap.Has(splitContent[0][1:])
	if err != nil {
		return nil, err
	}
	if !has {
		fatal("parseCustom called with non-custom command", errors.New("wat"))
	}
	return &Command{
		Type: Type_Custom,
		Custom: &CustomData{
			Call: splitContent[0][1:],
			Args: strings.Join(splitContent[1:], " "),
		},
	}, nil
}

///////////////////////////////////////////////////////////////////////////////
// Executors
///////////////////////////////////////////////////////////////////////////////

// LearnExecutor learns a user-generated command.
type LearnExecutor struct {
	commandMap StringMap
}

// NewLearnExecutor works as advertised.
func NewLearnExecutor(commandMap StringMap) *LearnExecutor {
	return &LearnExecutor{commandMap: commandMap}
}

// GetType returns the type of this feature.
func (f *LearnExecutor) GetType() int {
	return Type_Learn
}

// Execute replies over the given channel with a help message.
func (f *LearnExecutor) Execute(s DiscordSession, channel string, command *Command) {
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

type UnlearnExecutor struct {
	commandMap StringMap
}

func NewUnlearnExecutor(commandMap StringMap) *UnlearnExecutor {
	return &UnlearnExecutor{commandMap: commandMap}
}

// GetType returns the type of this feature.
func (f *UnlearnExecutor) GetType() int {
	return Type_Unlearn
}

// Execute replies over the given channel indicating successful unlearning, or
// failure to unlearn.
func (e *UnlearnExecutor) Execute(s DiscordSession, channel string, command *Command) {
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
	if has, err := e.commandMap.Has(command.Unlearn.Call); !has || err != nil {
		if has {
			fatal("Tried to unlearn command that doesn't exist: "+command.Unlearn.Call, errors.New("wat"))
		}
		fatal("Error in UnlearnFeature#execute, testing a command", err)
	}
	if err := e.commandMap.Delete(command.Unlearn.Call); err != nil {
		fatal("Unsuccessful unlearning a key; Dying since it might work with a restart", err)
	}

	// Send ack.
	s.ChannelMessageSend(channel, fmt.Sprintf(MsgUnlearnSuccess, command.Unlearn.Call))
}

type CustomExecutor struct {
	commandMap StringMap
}

func NewCustomExecutor(commandMap StringMap) *CustomExecutor {
	return &CustomExecutor{commandMap: commandMap}
}

// GetType returns the type of this feature.
func (e *CustomExecutor) GetType() int {
	return Type_Custom
}

// Execute returns the response if possible.
func (e *CustomExecutor) Execute(s DiscordSession, channel string, command *Command) {
	if command.Custom == nil {
		fatal("Incorrectly generated learn command", errors.New("wat"))
	}

	has, err := e.commandMap.Has(command.Custom.Call)
	if err != nil {
		fatal("Error testing custom feature", err)
	}
	if !has {
		fatal("Accidentally found a mismatched call/response pair", errors.New("Call response mismatch"))
	}

	response, err := e.commandMap.Get(command.Custom.Call)
	if err != nil {
		fatal("Error reading custom response", err)
	}

	if strings.Contains(response, "$1") {
		if command.Custom.Args == "" {
			response = MsgCustomNeedsArgs
		} else {
			response = strings.Replace(response, "$1", command.Custom.Args, 4)
		}
	}
	s.ChannelMessageSend(channel, response)
}
