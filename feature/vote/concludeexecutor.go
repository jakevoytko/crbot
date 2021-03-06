package vote

import (
	"fmt"
	"strings"

	"github.com/jakevoytko/crbot/api"
	"github.com/jakevoytko/crbot/log"
	"github.com/jakevoytko/crbot/model"
)

// ConcludeExecutor concludes the active vote and prints the results
type ConcludeExecutor struct {
	modelHelper *ModelHelper
}

// NewConcludeExecutor works as advertised
func NewConcludeExecutor(modelHelper *ModelHelper) *ConcludeExecutor {
	return &ConcludeExecutor{
		modelHelper: modelHelper,
	}
}

// GetType returns the type of this feature.
func (e *ConcludeExecutor) GetType() int {
	return model.CommandTypeVoteConclude
}

// PublicOnly returns whether the executor should be intercepted in a private
// channel. Since the vote is pinned to a channel, it should have been filtered
// then. If somehow the channel went private at that point, allow it to
// conclude.
func (e *ConcludeExecutor) PublicOnly() bool {
	return false
}

const (
	// MsgVoteConcluded is the header for a concluded vote
	MsgVoteConcluded = "@here -- Vote started by %s has concluded"
)

// Execute executes the command
func (e *ConcludeExecutor) Execute(s api.DiscordSession, channelID model.Snowflake, command *model.Command) {
	vote, err := e.modelHelper.MostRecentVote(channelID)
	if err != nil {
		log.Info("Error grabbing most recent vote", err)
		return
	}
	if vote == nil {
		log.Info("Tried to conclude nonexistant vote", err)
		return
	}

	user, err := s.User(vote.UserID.Format())
	if err != nil {
		log.Info("Error fetching the owner when rendering the status message", err)
		return
	}

	voteOutcome := model.VoteOutcomeNotEnough
	if vote.HasEnoughVotes() {
		voteOutcome = vote.CalculateActiveStatus()
	}
	vote.VoteOutcome = voteOutcome

	err = e.modelHelper.SetVoteOutcome(channelID, voteOutcome)
	if err != nil {
		// Log as info so that this doesn't crash-loop on startup.
		log.Info("Error setting vote outcome", err)
	}

	messages := []string{
		fmt.Sprintf(MsgVoteConcluded, user.Mention()),
		CompletedStatusLine(vote),
	}
	message := strings.Join(messages, "\n")
	if _, err := s.ChannelMessageSend(channelID.Format(), message); err != nil {
		log.Info("Error sending conclude message", err)
	}
}
