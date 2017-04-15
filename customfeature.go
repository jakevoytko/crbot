package main

import (
	"errors"
	"strings"
)

// CustomFeature is the fallback Feature that issues user-defined call and
// response commands. Accordingly, CustomFeature probably has a mouth on it.
type CustomFeature struct {
	commandMap StringMap
}

// NewCustomFeature returns a new CustomFeature.
func NewCustomFeature(commandMap StringMap) *CustomFeature {
	return &CustomFeature{
		commandMap: commandMap,
	}
}

// GetType returns the type of this feature.
func (f *CustomFeature) GetType() int {
	return Type_Custom
}

// Parsers returns nothing, since the custom parsers is a fallthrough parser.
func (f *CustomFeature) Parsers() []Parser {
	return []Parser{}
}

// FallbackParser returns the custom parser.
func (f *CustomFeature) FallbackParser() Parser {
	return NewCustomParser(f.commandMap)
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

// HelpText panics, since it should never be invoked.
func (p *CustomParser) HelpText() string {
	panic("CustomParser.HelpText cannot be called")
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

const (
	MsgCustomNeedsArgs = "This command takes args. Please type `?command <more text>` instead of `?command`"
)

// Execute returns the response if possible.
func (f *CustomFeature) Execute(s DiscordSession, channel string, command *Command) {
	if command.Custom == nil {
		fatal("Incorrectly generated learn command", errors.New("wat"))
	}

	has, err := f.commandMap.Has(command.Custom.Call)
	if err != nil {
		fatal("Error testing custom feature", err)
	}
	if !has {
		fatal("Accidentally found a mismatched call/response pair", errors.New("Call response mismatch"))
	}

	response, err := f.commandMap.Get(command.Custom.Call)
	if err != nil {
		fatal("Error reading custom response", err)
	}

	if strings.Contains(response, "$1") {
		if command.Custom.Args == "" {
			response = MsgCustomNeedsArgs
		} else {
			response = strings.Replace(response, "$1", command.Custom.Args, 1)
		}
	}
	s.ChannelMessageSend(channel, response)
}
