package vote

import (
	"fmt"
	"math"
	"time"

	"github.com/jakevoytko/crbot/api"
	"github.com/jakevoytko/crbot/model"
)

// Returns the full status line of an in-progress vote.
func StatusLine(clock model.UTCClock, vote *model.Vote) string {
	// Add the vote totals.
	statusStr := statusString(vote)
	votesFor := len(vote.VotesFor)
	votesAgainst := len(vote.VotesAgainst)
	votesForStr := MsgOneVoteFor
	if votesFor != 1 {
		votesForStr = fmt.Sprintf(MsgVotesFor, votesFor)
	}
	votesAgainstStr := MsgOneVoteAgainst
	if votesAgainst != 1 {
		votesAgainstStr = fmt.Sprintf(MsgVotesAgainst, votesAgainst)
	}

	timeString := TimeString(clock, vote.TimestampEnd)

	return statusStr + ". " + votesForStr + ", " + votesAgainstStr + ". " + timeString
}

// Returns the full status line of a concluded vote.
func CompletedStatusLine(vote *model.Vote) string {
	statusStr := MsgStatusInconclusive
	switch vote.VoteOutcome {
	case model.VoteOutcomeNotDone:
		statusStr = MsgStatusInconclusive // Don't know how this would happen.
	case model.VoteOutcomePassed:
		statusStr = MsgStatusVotePassed
	case model.VoteOutcomeFailed:
		statusStr = MsgStatusVoteFailed
	case model.VoteOutcomeNotEnough:
		statusStr = MsgStatusInconclusive
	}

	votesFor := len(vote.VotesFor)
	votesAgainst := len(vote.VotesAgainst)

	votesForStr := MsgOneVoteFor
	if votesFor != 1 {
		votesForStr = fmt.Sprintf(MsgVotesFor, votesFor)
	}
	votesAgainstStr := MsgOneVoteAgainst
	if votesAgainst != 1 {
		votesAgainstStr = fmt.Sprintf(MsgVotesAgainst, votesAgainst)
	}
	return statusStr + " " + votesForStr + ", " + votesAgainstStr
}

func statusString(vote *model.Vote) string {
	if vote.HasEnoughVotes() {
		switch vote.CalculateActiveStatus() {
		case model.VoteOutcomePassed:
			return MsgStatusVotePassing

		default:
			return MsgStatusVoteFailing
		}
	}
	return MsgStatusVotesNeeded
}

const (
	MsgNoTimeRemaining       = "No time remaining in vote"
	MsgMinutesRemaining      = "%v minutes remaining"
	MsgSecondsRemaining      = "%v seconds remaining"
	MsgMillisecondsRemaining = "%v milliseconds remaining"
)

func TimeString(clock model.UTCClock, timestampEnd time.Time) string {
	currentTime := clock.Now()
	remaining := timestampEnd.Sub(currentTime)
	timeString := MsgNoTimeRemaining

	// This shows the rounded-up time. Some examples:
	// time.Duration(30) * time.Minute -> 30 minutes
	// time.Duration(30) * time.Minute - time.Nanosecond -> 30 minutes
	// time.Duration(30) * time.Minute - time.Minute + time.Nanosecond -> 30 minutes
	// time.Duration(30) * time.Minute - time.Minute -> 29 minutes
	if remaining > time.Minute {
		minutes := int(math.Ceil(float64(remaining) / float64(time.Minute)))
		timeString = fmt.Sprintf(MsgMinutesRemaining, minutes)
	} else if remaining > time.Second {
		seconds := int(math.Ceil(float64(remaining) / float64(time.Second)))
		timeString = fmt.Sprintf(MsgSecondsRemaining, seconds)
	} else if remaining > time.Millisecond {
		milliseconds := int(math.Ceil(float64(remaining) / float64(time.Millisecond)))
		timeString = fmt.Sprintf(MsgMillisecondsRemaining, milliseconds)
	}

	// If it's less than a millisecond or already over, just go with the "no time
	// remaining" message.
	return timeString
}

// Iterates through active vote pointers to see if any require timers to be
// re-fired.
func handleVotesOnInitialLoad(s api.DiscordSession, modelHelper *ModelHelper, clock model.UTCClock, timer model.UTCTimer, commandChannel chan<- *model.Command) error {
	votes, err := modelHelper.MostRecentVotes()
	if err != nil {
		return err
	}

	now := clock.Now()

	for _, vote := range votes {
		// Start a timer so that this vote can conclude. This will conclude expired
		// votes immediately (negative durations cause timers to fire).
		if vote.VoteOutcome == model.VoteOutcomeNotDone {
			timer.ExecuteAfter(vote.TimestampEnd.Sub(now), func() {
				commandChannel <- &model.Command{
					Type:      model.Type_VoteConclude,
					ChannelID: vote.ChannelID,
				}
			})
		}
	}

	return nil
}
