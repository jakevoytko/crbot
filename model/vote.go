package model

import "time"

// Vote outcomes used for storage.
const (
	// These are serialized and stored, so they cannot change.
	VoteOutcomeNotDone   = 1
	VoteOutcomePassed    = 2
	VoteOutcomeFailed    = 3
	VoteOutcomeNotEnough = 4
)

// Vote is the JSON-serialized and -deserialized implementation of a single vote.
// TODO(jake): When there is a testrunner with vote-specific functionality, move
// this back into the vote package
type Vote struct {
	VoteID         int
	ChannelID      Snowflake
	UserID         Snowflake
	Message        string
	TimestampStart time.Time
	TimestampEnd   time.Time
	VotesFor       []Snowflake
	VotesAgainst   []Snowflake
	VoteOutcome    int
}

// NewVote works as advertised.
func NewVote(voteID int, channelID, userID Snowflake, message string, timestampStart, timestampEnd time.Time, votesFor, votesAgainst []Snowflake, voteOutcome int) *Vote {
	return &Vote{
		VoteID:         voteID,
		ChannelID:      channelID,
		UserID:         userID,
		Message:        message,
		TimestampStart: timestampStart,
		TimestampEnd:   timestampEnd,
		VotesFor:       votesFor,
		VotesAgainst:   votesAgainst,
		VoteOutcome:    voteOutcome,
	}
}

// HasEnoughVotes returns whether there are enough votes to claim confidence.
func (v *Vote) HasEnoughVotes() bool {
	return len(v.VotesFor)+len(v.VotesAgainst) >= 5
}

// CalculateActiveStatus compares the vote totals and returns what the outcome
// would be. This ignores the recorded outcome, and the number of votes.
func (v *Vote) CalculateActiveStatus() int {
	if len(v.VotesFor) > len(v.VotesAgainst) {
		return VoteOutcomePassed
	}
	return VoteOutcomeFailed
}
