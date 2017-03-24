package main

import (
	"errors"
	"strings"

	"gopkg.in/redis.v5"

	"github.com/bwmarrin/discordgo"
)

// CustomFeature is the fallback Feature that issues user-defined call and
// response commands. Accordingly, CustomFeature probably has a mouth on it.
type CustomFeature struct {
	redisClient *redis.Client
}

// NewCustomFeature returns a new CustomFeature.
func NewCustomFeature(redisClient *redis.Client) *CustomFeature {
	return &CustomFeature{
		redisClient: redisClient,
	}
}

// GetName returns the named type of this feature.
func (f *CustomFeature) GetName() string {
	return ""
}

// GetType returns the type of this feature.
func (f *CustomFeature) GetType() int {
	return Type_Custom
}

// Invokable returns whether users can invoke this from the command line.
func (f *CustomFeature) Invokable() bool {
	return false
}

// Parse parses the given custom command.
func (f *CustomFeature) Parse(splitContent []string) (*Command, error) {
	// TODO(jake): Drop this and external hash check, handle missing commands solely in execute.
	if !f.redisClient.HExists(Redis_Hash, splitContent[0][1:]).Val() {
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
func (f *CustomFeature) Execute(s *discordgo.Session, channel string, command *Command) {
	if command.Custom == nil {
		fatal("Incorrectly generated learn command", errors.New("wat"))
	}

	if !f.redisClient.HExists(Redis_Hash, command.Custom.Call).Val() {
		fatal("Accidentally found a mismatched call/response pair", errors.New("Call response mismatch"))
	}

	response := f.redisClient.HGet(Redis_Hash, command.Custom.Call).Val()

	if strings.Contains(response, "$1") {
		if command.Custom.Args == "" {
			response = MsgCustomNeedsArgs
		} else {
			response = strings.Replace(response, "$1", command.Custom.Args, 1)
		}
	}
	s.ChannelMessageSend(channel, response)
}
