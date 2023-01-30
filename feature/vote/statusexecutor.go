package vote

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jakevoytko/crbot/api"
	"github.com/jakevoytko/crbot/log"
	"github.com/jakevoytko/crbot/model"
)

// StatusExecutor allows the user to check on the progress of the current vote
type StatusExecutor struct {
	modelHelper *ModelHelper
}

// NewStatusExecutor works as advertised
func NewStatusExecutor(modelHelper *ModelHelper) *StatusExecutor {
	return &StatusExecutor{
		modelHelper: modelHelper,
	}
}

// GetType returns the type of this feature.
func (e *StatusExecutor) GetType() int {
	return model.CommandTypeVoteStatus
}

// PublicOnly returns whether the executor should be intercepted in a private channel.
func (e *StatusExecutor) PublicOnly() bool {
	return true
}

const (
	// MsgNoActiveVote prints that there was no active vote
	MsgNoActiveVote = "No active vote"
	// MsgOneVoteAgainst is the unpluralized message for vote against
	MsgOneVoteAgainst = "1 vote against"
	// MsgOneVoteFor is the unpluralized message for vote for
	MsgOneVoteFor = "1 vote for"
	// MsgSpacer is used to separate parts of the output string
	MsgSpacer = "-----"
	// MsgStatusInconclusive prints that a vote was inconclusive
	MsgStatusInconclusive = "Not enough votes were cast."
	// MsgStatusVoteFailed prints that a vote has failed
	MsgStatusVoteFailed = "Vote Failed."
	// MsgStatusVoteFailing prints that a vote is currently failing
	MsgStatusVoteFailing = "Vote is failing"
	// MsgStatusVotePassed indicates that the vote has passed
	MsgStatusVotePassed = "Vote Passed!"
	// MsgStatusVotePassing prints whether the vote is passing
	MsgStatusVotePassing = "Vote is passing"
	// MsgStatusVotesNeeded prints how many votes are needed
	MsgStatusVotesNeeded = "5 votes must be cast before vote can pass"
	// MsgVoteOwner prints who started the vote
	MsgVoteOwner = "Vote started by %s: "
	// MsgVotesAgainst prints the votes against
	MsgVotesAgainst = "%d votes against"
	// MsgVotesFor prints the votes for
	MsgVotesFor = "%d votes for"
)

// Execute prints the status of the current vote.
func (e *StatusExecutor) Execute(s api.DiscordSession, channel model.Snowflake, command *model.Command) {
	ok, err := e.modelHelper.IsVoteActive(channel)
	if err != nil {
		log.Fatal("Error reading vote status", err)
	}
	if !ok {
		if _, err := s.ChannelMessageSend(channel.Format(), MsgNoActiveVote); err != nil {
			log.Fatal("Unable to send no-active-vote message to user", err)
		}
		return
	}

	vote, err := e.modelHelper.MostRecentVote(channel)
	if err != nil {
		log.Fatal("Error pulling most recent vote", err)
	}
	if vote == nil {
		log.Fatal("Nil vote found after vote already active", errors.New("vote should not be null"))
		return
	}

	// The below creates a string like this:
	//
	// Vote started by @SomeoneElse: Votekick @Jake?
	// -----
	// 12 minutes remaining
	// 5 votes must be cast before vote can pass. 3 votes for, 1 vote against. 30 minutes remaining.

	messages := []string{}

	// Add the owner string.
	owner, err := s.User(vote.UserID.Format())
	if err != nil {
		log.Fatal("Error fetching the owner when rendering a vote response", err)
	}
	// Status line and message.
	messages = append(messages, fmt.Sprintf(MsgVoteOwner, owner.Username)+vote.Message)

	// Spacer
	messages = append(messages, MsgSpacer)

	messages = append(messages, StatusLine(e.modelHelper.UTCClock, vote))

	finalMessage := strings.Join(messages, "\n")
	if _, err := s.ChannelMessageSend(channel.Format(), finalMessage); err != nil {
		log.Info("Failed to send vote message", err)
	}
}
