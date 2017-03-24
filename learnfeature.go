package main

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"gopkg.in/redis.v5"

	"github.com/bwmarrin/discordgo"
)

// LearnFeature allows crbot to learn new calls and responses
type LearnFeature struct {
	featureRegistry *FeatureRegistry
	redisClient     *redis.Client
}

// NewLearnFeature returns a new LearnFeature.
func NewLearnFeature(featureRegistry *FeatureRegistry, redisClient *redis.Client) *LearnFeature {
	return &LearnFeature{
		featureRegistry: featureRegistry,
		redisClient:     redisClient,
	}
}

// GetName returns the named type of this feature.
func (f *LearnFeature) GetName() string {
	return Name_Learn
}

// GetType returns the type of this feature.
func (f *LearnFeature) GetType() int {
	return Type_Learn
}

// Invokable returns whether the user can execute this command by name.
func (f *LearnFeature) Invokable() bool {
	return true
}

// Parse parses the given learn command.
func (f *LearnFeature) Parse(splitContent []string) (*Command, error) {
	if splitContent[0] != f.GetName() {
		fatal("parseLearn called with non-learn command", errors.New("wat"))
	}

	callRegexp := regexp.MustCompile("(?s)^[[:alnum:]].*$")
	responseRegexp := regexp.MustCompile("(?s)^[^/?!].*$")

	// Show help when not enough data is present, or malicious data is present.
	if len(splitContent) < 3 || !callRegexp.MatchString(splitContent[1]) || !responseRegexp.MatchString(splitContent[2]) {
		return &Command{
			Type: Type_Help,
			Help: &HelpData{
				Type: Type_Learn,
			},
		}, nil
	}

	// Don't overwrite old or builtin commands.
	if f.redisClient.HExists(Redis_Hash, splitContent[1]).Val() || f.featureRegistry.GetTypeFromName(splitContent[1]) != Type_None {
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
func (f *LearnFeature) Execute(s *discordgo.Session, channel string, command *Command) {
	if command.Learn == nil {
		fatal("Incorrectly generated learn command", errors.New("wat"))
	}
	if !command.Learn.CallOpen {
		s.ChannelMessageSend(channel, fmt.Sprintf(MsgLearnFail, command.Learn.Call))
		return
	}

	// Teach the command.
	if f.redisClient.HExists(Redis_Hash, command.Learn.Call).Val() {
		fatal("Collision when adding a call for "+command.Learn.Call, errors.New("wat"))
	}
	f.redisClient.HSet(Redis_Hash, command.Learn.Call, command.Learn.Response)

	// Send ack.
	s.ChannelMessageSend(channel, fmt.Sprintf(MsgLearnSuccess, command.Learn.Call))
}
