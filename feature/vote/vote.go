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
	Message        string
	TimestampStart time.Time
	TimestampEnd   time.Time
	VotesFor       []int64
	VotesAgainst   []int64
	VoteOutcome    int
}

// NewVote works as advertised.
func NewVote(voteID int, userID int64, message string, timestampStart, timestampEnd time.Time, votesFor, votesAgainst []int64, voteOutcome int) *Vote {
	return &Vote{
		VoteID:         voteID,
		UserID:         userID,
		Message:        message,
		TimestampStart: timestampStart,
		TimestampEnd:   timestampEnd,
		VotesFor:       votesFor,
		VotesAgainst:   votesAgainst,
		VoteOutcome:    voteOutcome,
	}
}

// IsActive returns whether there are enough votes to claim confidence.
func (v *Vote) HasEnoughVotes() bool {
	return len(v.VotesFor) >= 5 || len(v.VotesAgainst) >= 5
}

// CalculateActiveStatus compares the vote totals and returns what the outcome
// would be. This ignores the recorded outcome, and the number of votes.
func (v *Vote) CalculateActiveStatus() int {
	if len(v.VotesFor) > len(v.VotesAgainst) {
		return VoteOutcomePassed
	}
	return VoteOutcomeFailed
}
