package moderation

import (
	"github.com/jakevoytko/crbot/api"
	"github.com/jakevoytko/crbot/log"
	"github.com/jakevoytko/crbot/model"
)

const (
	MsgRickList = "https://www.youtube.com/watch?v=dQw4w9WgXcQ"
)

// RickListExecutor prints a rick roll.
type RickListExecutor struct {
}

// NewRickListExecutor works as advertised.
func NewRickListExecutor() *RickListExecutor {
	return &RickListExecutor{}
}

// GetType returns the type.
func (e *RickListExecutor) GetType() int {
	return model.Type_RickList
}

// PublicOnly returns whether the executor should be intercepted in a private channel.
func (e *RickListExecutor) PublicOnly() bool {
	return false
}

// Execute replies over the given channel with a rick roll.
func (e *RickListExecutor) Execute(s api.DiscordSession, channel model.Snowflake, command *model.Command) {
	if _, err := s.ChannelMessageSend(channel.Format(), MsgRickList); err != nil {
		log.Info("Failed to send ricklist message", err)
	}
}
