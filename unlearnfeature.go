package main

import (
	"errors"
	"fmt"
	"regexp"

	"gopkg.in/redis.v5"

	"github.com/bwmarrin/discordgo"
)

// UnlearnFeature allows crbot to unlearn existing calls
type UnlearnFeature struct {
	featureRegistry *FeatureRegistry
	redisClient     *redis.Client
}

// NewUnlearnFeature returns a new UnlearnFeature.
func NewUnlearnFeature(featureRegistry *FeatureRegistry, redisClient *redis.Client) *UnlearnFeature {
	return &UnlearnFeature{
		featureRegistry: featureRegistry,
		redisClient:     redisClient,
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

	callRegexp := regexp.MustCompile("(?s)^[[:alnum:]].*$")

	// Show help when not enough data is present, or malicious data is present.
	if len(splitContent) < 2 || !callRegexp.MatchString(splitContent[1])) {
		return &Command{
			Type: Type_Help,
			Help: &HelpData{
				Type: Type_Unlearn,
			},
		}, nil
	}

	// Only unlearn commands that aren't built-in and exist
	if !f.redisClient.HExists(Redis_Hash, splitContent[1]).Val() || f.featureRegistry.GetTypeFromName(splitContent[1]) != Type_None {
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
	MsgUnlearnFail    = "I can't unlearn ?%s"
	MsgUnlearnSuccess = "Forgot about %s"
)

// Execute replies over the given channel indicating successful unlearning, or
// failure to unlearn.
func (f *UnlearnFeature) Execute(s *discordgo.Session, channel string, command *Command) {
	if command.Unlearn == nil {
		fatal("Incorrectly generated unlearn command", errors.New("wat"))
	}
	if !command.Unlearn.CallOpen {
		s.ChannelMessageSend(channel, fmt.Sprintf(MsgUnlearnFail, command.Unlearn.Call))
		return
	}

	// Remove the command.
	if !f.redisClient.HExists(Redis_Hash, command.Unlearn.Call).Val() {
		fatal("Tried to unlearn command that doesn't exist: "+command.Unlearn.Call, errors.New("wat"))
	}
	f.redisClient.HDel(Redis_Hash, command.Unlearn.Call)

	// Send ack.
	s.ChannelMessageSend(channel, fmt.Sprintf(MsgUnlearnSuccess, command.Unlearn.Call))
}
