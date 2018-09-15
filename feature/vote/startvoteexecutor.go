package vote

import (
	"fmt"

	"github.com/jakevoytko/crbot/api"
	"github.com/jakevoytko/crbot/log"
	"github.com/jakevoytko/crbot/model"
)

// StartVoteExecutor executes a vote begin command
type StartVoteExecutor struct {
	modelHelper    *ModelHelper
	commandChannel chan<- *model.Command
	utcTimer       model.UTCTimer
}

// NewStartVoteExecutor works as advertised
func NewStartVoteExecutor(modelHelper *ModelHelper, commandChannel chan<- *model.Command, utcTimer model.UTCTimer) *StartVoteExecutor {
	return &StartVoteExecutor{
		modelHelper:    modelHelper,
		commandChannel: commandChannel,
		utcTimer:       utcTimer,
	}
}

// GetType returns the type of this feature.
func (e *StartVoteExecutor) GetType() int {
	return model.CommandTypeVote
}

// PublicOnly returns whether the executor should be intercepted in a private channel.
func (e *StartVoteExecutor) PublicOnly() bool {
	return true
}

const (
	// MsgActiveVote prints that a vote is already active
	MsgActiveVote = "Cannot start a vote while another is in progress. Type `?votestatus` for more info"
	// MsgBroadcastNewVote prints that a new vote is happening
	MsgBroadcastNewVote = "@everyone -- %s started a new vote: %s\n\nType `?yes` or `?no` to vote. 30 minutes remaining."
)

// Execute starts a new vote if one is not already active. It also starts a
// timer to use to conclude the vote.
func (e *StartVoteExecutor) Execute(s api.DiscordSession, channelID model.Snowflake, command *model.Command) {
	ok, err := e.modelHelper.IsVoteActive(channelID)
	if err != nil {
		log.Fatal("Error occurred while calling for active vote", err)
	}
	if ok {
		_, err := s.ChannelMessageSend(channelID.Format(), MsgActiveVote)
		if err != nil {
			log.Fatal("Unable to send vote-already-active message to user", err)
		}
		return
	}

	userID, err := model.ParseSnowflake(command.Author.ID)
	if err != nil {
		log.Info("Error parsing command user ID", err)
		return
	}
	vote, err := e.modelHelper.StartNewVote(channelID, userID, command.Vote.Message)
	if err != nil {
		log.Fatal("error starting new vote", err)
	}

	broadcastMessage := fmt.Sprintf(MsgBroadcastNewVote, command.Author.Mention(), command.Vote.Message)
	_, err = s.ChannelMessageSend(channelID.Format(), broadcastMessage)
	if err != nil {
		log.Fatal("Unable to broadcast new message across the channel", err)
	}

	// After the vote has expired, send a conclude command so the status can be
	// written to storage and printed to the users.
	e.utcTimer.ExecuteAfter(vote.TimestampEnd.Sub(vote.TimestampStart), func() {
		e.commandChannel <- &model.Command{
			Type:      model.CommandTypeVoteConclude,
			ChannelID: channelID,
		}
	})
}
