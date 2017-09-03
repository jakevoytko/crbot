package main

import (
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/jakevoytko/crbot/api"
	"github.com/jakevoytko/crbot/app"
	"github.com/jakevoytko/crbot/feature"
	"github.com/jakevoytko/crbot/log"
	"github.com/jakevoytko/crbot/model"
)

// ModerationFeature is a Feature that executes moderation events.
type ModerationFeature struct {
	featureRegistry *feature.Registry
	config          *app.Config
}

// NewModerationFeature returns a new ModerationFeature.
func NewModerationFeature(featureRegistry *feature.Registry, config *app.Config) *ModerationFeature {
	return &ModerationFeature{
		featureRegistry: featureRegistry,
		config:          config,
	}
}

// Parsers returns the parsers.
func (f *ModerationFeature) Parsers() []feature.Parser {
	return []feature.Parser{}
}

// CommandInterceptors returns command interceptors.
func (f *ModerationFeature) CommandInterceptors() []feature.CommandInterceptor {
	return []feature.CommandInterceptor{
		NewRickListCommandInterceptor(f.config),
	}
}

// FallbackParser returns nil.
func (f *ModerationFeature) FallbackParser() feature.Parser {
	return nil
}

// RickListCommandInterceptor asserts that the
type RickListCommandInterceptor struct {
	rickList []int64
}

// NewRickListCommandInterceptor returns a new ricklist command interceptor.
func NewRickListCommandInterceptor(config *app.Config) *RickListCommandInterceptor {
	return &RickListCommandInterceptor{
		rickList: config.RickList,
	}
}

// Intercept checks whether the command is forbidden by the ricklist.
func (i *RickListCommandInterceptor) Intercept(command *model.Command, s api.DiscordSession, m *discordgo.MessageCreate) (*model.Command, error) {
	// Check moderation.
	// RickList
	// - RickListed users can only use ?learn in private channels, without it responding with
	//   a rickroll.
	if channel, err := s.Channel(m.ChannelID); err == nil && channel.IsPrivate && command.Type != model.Type_Learn {
		for _, ricked := range i.rickList {
			if strconv.FormatInt(ricked, 10) == m.Author.ID {
				return &model.Command{
					Type: model.Type_RickList,
				}, nil
			}
		}
	}

	return command, nil
}

// Executors gets the executors.
func (f *ModerationFeature) Executors() []feature.Executor {
	return []feature.Executor{NewRickListExecutor(f.featureRegistry)}
}

const (
	MsgRickList = "https://www.youtube.com/watch?v=dQw4w9WgXcQ"
)

// RickListExecutor prints a rick roll.
type RickListExecutor struct {
	featureRegistry *feature.Registry
}

// NewRickListExecutor works as advertised.
func NewRickListExecutor(featureRegistry *feature.Registry) *RickListExecutor {
	return &RickListExecutor{
		featureRegistry: featureRegistry,
	}
}

// GetType returns the type.
func (e *RickListExecutor) GetType() int {
	return model.Type_RickList
}

// Execute replies over the given channel with a rick roll.
func (e *RickListExecutor) Execute(s api.DiscordSession, channel string, command *model.Command) {
	if _, err := s.ChannelMessageSend(channel, MsgRickList); err != nil {
		log.Info("Failed to send ricklist message", err)
	}
}
