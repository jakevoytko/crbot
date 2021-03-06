package vote

import (
	"fmt"
	"strings"

	"github.com/jakevoytko/crbot/api"
	"github.com/jakevoytko/crbot/log"
	"github.com/jakevoytko/crbot/model"
)

// BallotExecutor executes a vote
type BallotExecutor struct {
	modelHelper *ModelHelper
}

// NewBallotExecutor works as advertised
func NewBallotExecutor(modelHelper *ModelHelper) *BallotExecutor {
	return &BallotExecutor{
		modelHelper: modelHelper,
	}
}

// GetType returns the type of this feature.
func (e *BallotExecutor) GetType() int {
	return model.CommandTypeVoteBallot
}

// PublicOnly returns whether the executor should be intercepted in a private channel.
func (e *BallotExecutor) PublicOnly() bool {
	return true
}

const (
	// MsgAlreadyVoted returns that the user already voted
	MsgAlreadyVoted = "%v already voted"
	// MsgVotedAgainst returns that the user voted against the vote
	MsgVotedAgainst = "%v voted no"
	// MsgVotedInFavor returns that the user voted for the vote
	MsgVotedInFavor = "%v voted yes"
)

// Execute runs the command
func (e *BallotExecutor) Execute(s api.DiscordSession, channelID model.Snowflake, command *model.Command) {
	userID, err := model.ParseSnowflake(command.Author.ID)
	if err != nil {
		log.Fatal("Error parsing discord user ID", err)
	}

	vote, err := e.modelHelper.CastBallot(channelID, userID, command.Ballot.InFavor)
	switch err {
	case ErrorNoVoteActive:
		if _, err := s.ChannelMessageSend(channelID.Format(), MsgNoActiveVote); err != nil {
			log.Fatal("Unable to send no-active-vote message to user", err)
		}
		return

	case ErrorAlreadyVoted:
		if _, err := s.ChannelMessageSend(channelID.Format(), fmt.Sprintf(MsgAlreadyVoted, command.Author.Mention())); err != nil {
			log.Fatal("Unable to send already voted message to user", err)
		}
		return
	}

	voteMessage := fmt.Sprintf(MsgVotedAgainst, command.Author.Mention())
	if command.Ballot.InFavor {
		voteMessage = fmt.Sprintf(MsgVotedInFavor, command.Author.Mention())
	}

	messages := []string{voteMessage, StatusLine(e.modelHelper.UTCClock, vote)}
	message := strings.Join(messages, "\n")
	if _, err := s.ChannelMessageSend(channelID.Format(), message); err != nil {
		log.Info("Failed to send ballot status message", err)
	}
}
