package vote

import "time"

const (
	// These are serialized and stored, so they cannot change.
	VoteOutcomeNotDone   = 1
	VoteOutcomePassed    = 2
	VoteOutcomeFailed    = 3
	VoteOutcomeNotEnough = 4
)

// Vote is the JSON-serialized and -deserialized implementation of a single vote.
type Vote struct {
	VoteID         int
	UserID         int64
	TimestampStart time.Time
	TimestampEnd   time.Time
	VotesFor       []int64
	VotesAgainst   []int64
	VoteOutcome    int
}

// NewVote works as advertised.
func NewVote(voteID int, userID int64, timestampStart, timestampEnd time.Time, votesFor, votesAgainst []int64, voteOutcome int) *Vote {
	return &Vote{
		VoteID:         voteID,
		UserID:         userID,
		TimestampStart: timestampStart,
		TimestampEnd:   timestampEnd,
		VotesFor:       votesFor,
		VotesAgainst:   votesAgainst,
		VoteOutcome:    voteOutcome,
	}
}
