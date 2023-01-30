package vote

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jakevoytko/crbot/model"
	stringmap "github.com/jakevoytko/go-stringmap"
)

// ModelHelper adds helper functions for using votes with a
// stringmap. Currently, it handles a majority of business logic. However, if
// functionality is ever added to allow a moderator to edit a vote after it has
// happened, the business logic will need to be pulled of here and leave this to
// be strictly a data structure that only validates its own sanity. The only
// logic that seems like it should be handled but isn't is the logic for when a
// vote outcome is recorded, because that is handled in an evented manner
// (i.e. when a timer ends).
type ModelHelper struct {
	StringMap stringmap.StringMap
	UTCClock  model.UTCClock
}

// NewModelHelper works as advertised.
func NewModelHelper(stringMap stringmap.StringMap, utcClock model.UTCClock) *ModelHelper {
	return &ModelHelper{
		StringMap: stringMap,
		UTCClock:  utcClock,
	}
}

const (
	// KeyMostRecentVoteID is the key/value store key for the vote
	KeyMostRecentVoteID = "most-recent-vote-id-channel-%v"
	// RedisMostRecentVoteID is the vote id of the most recent vote
	RedisMostRecentVoteID = "most-recent-vote-id-channel-*"
	// KeyVoteTemplate is the map from vote to channel
	KeyVoteTemplate = "vote-%v-channel-%v"
	// VoteDuration is the duration of a vote
	VoteDuration = time.Duration(30) * time.Minute
)

// ErrorOnlyOneVote indicates that a second vote cannot be started
var ErrorOnlyOneVote = errors.New("tried to start vote when one is already active")

// ErrorNoVoteActive indicates that the user can't vote if no vote is active
var ErrorNoVoteActive = errors.New("cannot vote when there is no active vote")

// ErrorAlreadyVoted indicates that the user can't vote twice
var ErrorAlreadyVoted = errors.New("user already voted")

// ErrorVoteHasOutcome indicates that the application already set the vote outcome, and to give up
var ErrorVoteHasOutcome = errors.New("cannot change vote outcome")

// IsVoteActive returns whether there is a most-recent, active vote.
func (h *ModelHelper) IsVoteActive(channelID model.Snowflake) (bool, error) {
	vote, err := h.MostRecentVote(channelID)
	if err != nil {
		return false, err
	}
	if vote == nil {
		return false, nil
	}

	currentTime := h.UTCClock.Now()

	// Bail if the current time is not within the vote's range.
	return vote.VoteOutcome == model.VoteOutcomeNotDone &&
		currentTime.Sub(vote.TimestampStart) >= 0 && vote.TimestampEnd.Sub(currentTime) > 0, nil
}

// MostRecentVotes returns every most recent vote in the database.
func (h *ModelHelper) MostRecentVotes() ([]*model.Vote, error) {
	results, err := h.StringMap.ScanKeys(RedisMostRecentVoteID)
	if err != nil {
		return nil, err
	}

	votes := []*model.Vote{}
	for _, key := range results {
		mostRecentVoteKey, err := h.StringMap.Get(key)
		if err != nil {
			return nil, err
		}

		ok, err := h.StringMap.Has(mostRecentVoteKey)
		if err != nil {
			return nil, err
		}
		if !ok {
			// Missing vote for some reason. Ignore.
			continue
		}

		mostRecentVote, err := h.StringMap.Get(mostRecentVoteKey)
		if err != nil {
			return nil, err
		}

		var vote model.Vote
		err = json.Unmarshal([]byte(mostRecentVote), &vote)
		if err != nil {
			return nil, err
		}
		votes = append(votes, &vote)
	}

	return votes, nil
}

// MostRecentVote returns the active vote, or nil if none present. Returns an error
// on i/o problems.
func (h *ModelHelper) MostRecentVote(channelID model.Snowflake) (*model.Vote, error) {
	reifiedKey := fmt.Sprintf(KeyMostRecentVoteID, channelID)
	ok, err := h.StringMap.Has(reifiedKey)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}

	mostRecentVoteID, err := h.StringMap.Get(reifiedKey)
	if err != nil {
		return nil, err
	}
	mostRecentVote, err := h.StringMap.Get(mostRecentVoteID)
	if err != nil {
		return nil, err
	}

	var deserializedVote model.Vote
	err = json.Unmarshal([]byte(mostRecentVote), &deserializedVote)
	if err != nil {
		return nil, err
	}
	return &deserializedVote, nil
}

// MostRecentVoteID returns the most recent ID. Returns `0, nil` if no vote has
// ever been executed.
func (h *ModelHelper) MostRecentVoteID(channelID model.Snowflake) (int, error) {
	vote, err := h.MostRecentVote(channelID)
	if err != nil {
		return 0, err
	}
	if vote == nil {
		return 0, nil
	}

	return vote.VoteID, nil
}

// StartNewVote starts and returns a new vote. Returns ErrorOnlyOneVote if
// another vote was active when trying to start this one.
func (h *ModelHelper) StartNewVote(channelID, userID model.Snowflake, message string) (*model.Vote, error) {
	// Don't overwrite an existing vote.
	if ok, err := h.IsVoteActive(channelID); ok || err != nil {
		if err != nil {
			return nil, err
		}
		return nil, ErrorOnlyOneVote
	}

	// Generate a new Vote with no ballots, starting now in UTC.
	mostRecentVote, err := h.MostRecentVote(channelID)
	if err != nil {
		return nil, err
	}
	nextVoteID := 1
	if mostRecentVote != nil {
		nextVoteID = mostRecentVote.VoteID + 1
	}
	voteStart := h.UTCClock.Now()
	voteEnd := voteStart.Add(VoteDuration)
	vote := model.NewVote(
		nextVoteID, channelID, userID, message, voteStart, voteEnd, []model.Snowflake{}, []model.Snowflake{}, model.VoteOutcomeNotDone)

	err = h.writeVote(vote)
	if err != nil {
		return nil, err
	}

	return vote, nil
}

// CastBallot casts a ballot against the current poll for the given user. On
// success, it returns the vote with the ballot incorporated. Returns
// ErrorNoVoteActive if there is no active poll, or the inner error if a
// component errored. Returns ErrorAlreadyVoted if the user is present in either
// list.
func (h *ModelHelper) CastBallot(channelID model.Snowflake, userID model.Snowflake, inFavor bool) (*model.Vote, error) {
	ok, err := h.IsVoteActive(channelID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrorNoVoteActive
	}

	vote, err := h.MostRecentVote(channelID)
	if err != nil {
		return nil, err
	}
	if vote == nil {
		return nil, ErrorNoVoteActive
	}

	// Ensure the user hasn't already voted.
	for _, id := range vote.VotesFor {
		if id == userID {
			return nil, ErrorAlreadyVoted
		}
	}
	for _, id := range vote.VotesAgainst {
		if id == userID {
			return nil, ErrorAlreadyVoted
		}
	}

	if inFavor {
		vote.VotesFor = append(vote.VotesFor, userID)
	} else {
		vote.VotesAgainst = append(vote.VotesAgainst, userID)
	}

	err = h.writeVote(vote)
	if err != nil {
		return nil, err
	}

	return vote, nil
}

// SetVoteOutcome terminates an active vote with the given outcome
func (h *ModelHelper) SetVoteOutcome(channelID model.Snowflake, voteOutcome int) error {
	vote, err := h.MostRecentVote(channelID)
	if err != nil {
		return err
	}
	if vote == nil {
		return ErrorNoVoteActive
	}

	// Ensure the user hasn't already set an outcome.
	if vote.VoteOutcome != model.VoteOutcomeNotDone {
		return ErrorVoteHasOutcome
	}
	vote.VoteOutcome = voteOutcome

	err = h.writeVote(vote)
	if err != nil {
		return err
	}

	return nil
}

func (h *ModelHelper) writeVote(vote *model.Vote) error {
	// Serialize and write.
	serializedVote, err := json.Marshal(vote)
	if err != nil {
		return err
	}

	// Write the vote.
	voteKey := fmt.Sprintf(KeyVoteTemplate, vote.VoteID, vote.ChannelID)
	err = h.StringMap.Set(voteKey, string(serializedVote))
	if err != nil {
		return err
	}

	// Write the metadata afterwards, so it's guaranteed to always point to a
	// valid vote record.
	err = h.StringMap.Set(fmt.Sprintf(KeyMostRecentVoteID, vote.ChannelID), voteKey)
	if err != nil {
		return err
	}

	return nil
}
