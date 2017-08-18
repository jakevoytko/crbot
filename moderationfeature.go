package main

import (
	"strconv"

	"github.com/bwmarrin/discordgo"
)

// ModerationFeature is a Feature that executes moderation events.
type ModerationFeature struct {
	featureRegistry *FeatureRegistry
	config          *Config
}

// NewModerationFeature returns a new ModerationFeature.
func NewModerationFeature(featureRegistry *FeatureRegistry, config *Config) *ModerationFeature {
	return &ModerationFeature{
		featureRegistry: featureRegistry,
		config:          config,
	}
}

// Parsers returns the parsers.
func (f *ModerationFeature) Parsers() []Parser {
	return []Parser{}
}

// CommandInterceptors returns command interceptors.
func (f *ModerationFeature) CommandInterceptors() []CommandInterceptor {
	return []CommandInterceptor{
		NewRickListCommandInterceptor(f.config),
	}
}

// FallbackParser returns nil.
func (f *ModerationFeature) FallbackParser() Parser {
	return nil
}

// RickListCommandInterceptor asserts that the
type RickListCommandInterceptor struct {
	rickList []int64
}

// NewRickListCommandInterceptor returns a new ricklist command interceptor.
func NewRickListCommandInterceptor(config *Config) *RickListCommandInterceptor {
	return &RickListCommandInterceptor{
		rickList: config.RickList,
	}
}

// Intercept checks whether the command is forbidden by the ricklist.
func (i *RickListCommandInterceptor) Intercept(command *Command, s DiscordSession, m *discordgo.MessageCreate) (*Command, error) {
	// Check moderation.
	// RickList
	// - RickListed users can only use ?learn in private channels, without it responding with
	//   a rickroll.
	if channel, err := s.Channel(m.ChannelID); err == nil && channel.IsPrivate && command.Type != Type_Learn {
		for _, ricked := range i.rickList {
			if strconv.FormatInt(ricked, 10) == m.Author.ID {
				return &Command{
					Type: Type_RickList,
				}, nil
			}
		}
	}

	return command, nil
}

// Executors gets the executors.
func (f *ModerationFeature) Executors() []Executor {
	return []Executor{NewRickListExecutor(f.featureRegistry)}
}

const (
	MsgRickList = "https://www.youtube.com/watch?v=dQw4w9WgXcQ"
)

// RickListExecutor prints a rick roll.
type RickListExecutor struct {
	featureRegistry *FeatureRegistry
}

// NewRickListExecutor works as advertised.
func NewRickListExecutor(featureRegistry *FeatureRegistry) *RickListExecutor {
	return &RickListExecutor{
		featureRegistry: featureRegistry,
	}
}

// GetType returns the type.
func (e *RickListExecutor) GetType() int {
	return Type_RickList
}

// Execute replies over the given channel with a rick roll.
func (e *RickListExecutor) Execute(s DiscordSession, channel string, command *Command) {
	if _, err := s.ChannelMessageSend(channel, MsgRickList); err != nil {
		info("Failed to send ricklist message", err)
	}
}
