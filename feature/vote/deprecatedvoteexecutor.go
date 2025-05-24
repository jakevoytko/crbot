package vote

import (
	"github.com/jakevoytko/crbot/api"
	"github.com/jakevoytko/crbot/log"
	"github.com/jakevoytko/crbot/model"
)

// DeprecatedVoteExecutor executes a vote begin command
type DeprecatedVoteExecutor struct {
	commandType int
}

// NewDeprecatedVoteExecutor works as advertised
func NewDeprecatedVoteExecutor(commandType int) *DeprecatedVoteExecutor {
	return &DeprecatedVoteExecutor{
		commandType: commandType,
	}
}

// GetType returns the type of this feature.
func (e *DeprecatedVoteExecutor) GetType() int {
	return e.commandType
}

// PublicOnly returns whether the executor should be intercepted in a private channel.
func (e *DeprecatedVoteExecutor) PublicOnly() bool {
	return true
}

// Execute starts a new vote if one is not already active. It also starts a
// timer to use to conclude the vote.
func (e *DeprecatedVoteExecutor) Execute(s api.DiscordSession, channelID model.Snowflake, command *model.Command) {
	_, err := s.ChannelMessageSend(channelID.Format(), "Deprecated. Please use Discord polls.")
	if err != nil {
		log.Fatal("Unable to send deprecated message to user", err)
	}
	return
}
