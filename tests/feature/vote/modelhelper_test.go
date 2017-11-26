package vote

import (
	"reflect"
	"testing"
	"time"

	"github.com/jakevoytko/crbot/feature/vote"
	"github.com/jakevoytko/crbot/model"
	"github.com/jakevoytko/crbot/testutil"
)

const (
	Channel1 = model.Snowflake(1)
	Channel2 = model.Snowflake(2)
	UserID1  = model.Snowflake(8675309)
	UserID2  = model.Snowflake(9000000)
)

func TestIsVoteActive_NoMemory(t *testing.T) {
	modelHelper, _ := initializeTests()

	assertIsVoteActive(t, modelHelper, Channel1, false)
	assertIsVoteActive(t, modelHelper, Channel2, false)
}

func TestIsVoteActive_ActiveVote(t *testing.T) {
	modelHelper, _ := initializeTests()

	// Active on one channel
	assertStartNewVote(t, modelHelper, Channel1, UserID1)
	assertIsVoteActive(t, modelHelper, Channel1, true)

	// But not the other
	assertIsVoteActive(t, modelHelper, Channel2, false)
}

func TestIsVoteActive_2Votes(t *testing.T) {
	modelHelper, _ := initializeTests()

	// Active on one channel
	assertStartNewVote(t, modelHelper, Channel1, UserID1)
	assertIsVoteActive(t, modelHelper, Channel1, true)

	assertIsVoteActive(t, modelHelper, Channel2, false)

	// And then on the other
	assertStartNewVote(t, modelHelper, Channel2, UserID1)
	assertIsVoteActive(t, modelHelper, Channel2, true)
}

func TestIsVoteActive_TimeExpires(t *testing.T) {
	modelHelper, clock := initializeTests()

	assertStartNewVote(t, modelHelper, Channel1, UserID1)
	assertIsVoteActive(t, modelHelper, Channel1, true)

	// Advance the clock 1ns before the vote will expire, and then verify that
	// advancing it 1ns causes the vote to expire.
	clock.Advance(vote.VoteDuration - time.Duration(1)*time.Nanosecond)
	assertIsVoteActive(t, modelHelper, Channel1, true)
	clock.Advance(time.Duration(1) * time.Nanosecond)
	assertIsVoteActive(t, modelHelper, Channel1, false)
}

func TestIsVoteActive_TimeExpiresForBothVotes(t *testing.T) {
	modelHelper, clock := initializeTests()

	assertStartNewVote(t, modelHelper, Channel1, UserID1)
	assertStartNewVote(t, modelHelper, Channel2, UserID1)
	assertIsVoteActive(t, modelHelper, Channel1, true)
	assertIsVoteActive(t, modelHelper, Channel2, true)

	// Advance the clock 1ns before the vote will expire, and then verify that
	// advancing it 1ns causes the vote to expire.
	clock.Advance(vote.VoteDuration - time.Duration(1)*time.Nanosecond)
	assertIsVoteActive(t, modelHelper, Channel1, true)
	assertIsVoteActive(t, modelHelper, Channel2, true)
	clock.Advance(time.Duration(1) * time.Nanosecond)
	assertIsVoteActive(t, modelHelper, Channel1, false)
	assertIsVoteActive(t, modelHelper, Channel2, false)
}

func TestIsVoteActive_OffsetStartTimes(t *testing.T) {
	modelHelper, clock := initializeTests()

	// Start a vote
	assertStartNewVote(t, modelHelper, Channel1, UserID1)
	assertIsVoteActive(t, modelHelper, Channel1, true)

	// Advance the clock 15 minutes. Start a second vote and assert that they are
	// both active.
	advanceTime := time.Duration(15) * time.Minute
	clock.Advance(advanceTime)
	assertStartNewVote(t, modelHelper, Channel2, UserID1)
	assertIsVoteActive(t, modelHelper, Channel1, true)
	assertIsVoteActive(t, modelHelper, Channel2, true)

	// Advance the clock 15m... assert that the first vote has expired.
	clock.Advance(advanceTime)
	assertIsVoteActive(t, modelHelper, Channel1, false)
	assertIsVoteActive(t, modelHelper, Channel2, true)

	// And then assert that the vote can expire.
	clock.Advance(advanceTime)
	assertIsVoteActive(t, modelHelper, Channel1, false)
	assertIsVoteActive(t, modelHelper, Channel2, false)
}

func TestMostRecentVote_NoMemory(t *testing.T) {
	modelHelper, _ := initializeTests()

	assertMostRecentVote(t, modelHelper, Channel1, nil)
	assertMostRecentVote(t, modelHelper, Channel2, nil)
}

func TestMostRecentVote_ActiveVote(t *testing.T) {
	modelHelper, _ := initializeTests()

	newVote := assertStartNewVote(t, modelHelper, Channel1, UserID1)
	assertMostRecentVote(t, modelHelper, Channel1, newVote)
	assertMostRecentVote(t, modelHelper, Channel2, nil)
}

func TestMostRecentVote_2ActiveVotes(t *testing.T) {
	modelHelper, _ := initializeTests()

	newVote1 := assertStartNewVote(t, modelHelper, Channel1, UserID1)
	newVote2 := assertStartNewVote(t, modelHelper, Channel2, UserID1)
	assertMostRecentVote(t, modelHelper, Channel1, newVote1)
	assertMostRecentVote(t, modelHelper, Channel2, newVote2)
}

func TestMostRecentVote_ReturnsExpiredVote(t *testing.T) {
	modelHelper, clock := initializeTests()

	newVote := assertStartNewVote(t, modelHelper, Channel1, UserID1)
	assertIsVoteActive(t, modelHelper, Channel1, true)

	// Check that the correct vote is still returned, even though it has expired.
	clock.Advance(vote.VoteDuration)
	assertIsVoteActive(t, modelHelper, Channel1, false)
	assertMostRecentVote(t, modelHelper, Channel1, newVote)
}

func TestMostRecentVoteID_NoMemory(t *testing.T) {
	modelHelper, _ := initializeTests()

	assertMostRecentVoteID(t, modelHelper, Channel1, 0)
	assertMostRecentVoteID(t, modelHelper, Channel2, 0)
}

func TestMostRecentVoteID_ActiveVote(t *testing.T) {
	modelHelper, _ := initializeTests()

	newVote := assertStartNewVote(t, modelHelper, Channel1, UserID1)
	assertMostRecentVoteID(t, modelHelper, Channel1, newVote.VoteID)
	assertMostRecentVoteID(t, modelHelper, Channel2, 0)
}

func TestMostRecentVoteID_2ActiveVotes(t *testing.T) {
	modelHelper, _ := initializeTests()

	newVote1 := assertStartNewVote(t, modelHelper, Channel1, UserID1)
	newVote2 := assertStartNewVote(t, modelHelper, Channel2, UserID1)
	assertMostRecentVoteID(t, modelHelper, Channel1, newVote1.VoteID)
	assertMostRecentVoteID(t, modelHelper, Channel2, newVote2.VoteID)
}

func TestMostRecentVoteID_ReturnsExpiredVote(t *testing.T) {
	modelHelper, clock := initializeTests()

	newVote := assertStartNewVote(t, modelHelper, Channel1, UserID1)
	clock.Advance(vote.VoteDuration)
	assertMostRecentVoteID(t, modelHelper, Channel1, newVote.VoteID)
}

func TestStartNewVote_AddingFromNoMemory(t *testing.T) {
	modelHelper, _ := initializeTests()

	vote := assertStartNewVote(t, modelHelper, Channel1, UserID1)

	if vote == nil {
		t.Fatalf("Should have returned a new vote")
	}
	if vote.VoteID != 1 {
		t.Errorf("Wrong vote ID")
	}
	if vote.UserID != 8675309 {
		t.Errorf("Wrong user ID")
	}
	if len(vote.VotesFor) != 0 || len(vote.VotesAgainst) != 0 {
		t.Errorf("New vote should have no recorded votes")
	}
	if vote.VoteOutcome != model.VoteOutcomeNotDone {
		t.Errorf("Vote outcome should be NotDone")
	}

	assertMostRecentVote(t, modelHelper, Channel1, vote)

	assertMostRecentVoteID(t, modelHelper, Channel1, vote.VoteID)
}

func TestStartNewVote_DoesNotOverrideExistingVote(t *testing.T) {
	modelHelper, _ := initializeTests()

	assertStartNewVote(t, modelHelper, Channel1, UserID1)
	assertStartNewVoteFails(t, modelHelper, Channel1, 11)
}

func TestCastBallot_NoVoteActive(t *testing.T) {
	modelHelper, _ := initializeTests()

	assertEarlyBallotFails(t, modelHelper, Channel1, UserID1, true /* inFavor */)
	assertEarlyBallotFails(t, modelHelper, Channel1, UserID1, false /* inFavor */)
}

func TestCastBallot_VoteFor(t *testing.T) {
	modelHelper, _ := initializeTests()

	assertStartNewVote(t, modelHelper, Channel1, UserID1)
	assertCastBallot(t, modelHelper, Channel1, UserID1, true)
}

func TestCastBallot_VoteAgainst(t *testing.T) {
	modelHelper, _ := initializeTests()

	assertStartNewVote(t, modelHelper, Channel1, UserID1)
	assertCastBallot(t, modelHelper, Channel1, 1 /* userID */, false)
}

func TestCastBallot_CanNotVoteTwiceAfterInFavor(t *testing.T) {
	modelHelper, _ := initializeTests()

	assertStartNewVote(t, modelHelper, Channel1, UserID1)
	assertCastBallot(t, modelHelper, Channel1, UserID1, true)
	assertCannotVoteAgain(t, modelHelper, Channel1, UserID1, true)
	assertCannotVoteAgain(t, modelHelper, Channel1, UserID1, false)
}

func TestCastBallot_CanNotVoteTwiceAfterOpposed(t *testing.T) {
	modelHelper, _ := initializeTests()

	assertStartNewVote(t, modelHelper, Channel1, UserID1)
	assertCastBallot(t, modelHelper, Channel1, UserID1, false)
	assertCannotVoteAgain(t, modelHelper, Channel1, UserID1, true)
	assertCannotVoteAgain(t, modelHelper, Channel1, UserID1, false)
}

func TestCastBallot_LotsOfVotes(t *testing.T) {
	modelHelper, _ := initializeTests()

	assertStartNewVote(t, modelHelper, Channel1, 1 /* userID */)
	for userID := 0; userID < 10; userID++ {
		inFavor := userID%2 == 0
		assertCastBallot(t, modelHelper, Channel1, model.Snowflake(userID), inFavor)
	}
}

func TestCastBallot_LotsOfVotes_2Channels(t *testing.T) {
	modelHelper, _ := initializeTests()

	assertStartNewVote(t, modelHelper, Channel1, 1 /* userID */)
	assertStartNewVote(t, modelHelper, Channel2, 1 /* userID */)
	for userID := 0; userID < 10; userID++ {
		inFavor := userID%2 == 0
		// Votes are opposite per-channel for testing.
		assertCastBallot(t, modelHelper, Channel1, model.Snowflake(userID), inFavor)
		assertCastBallot(t, modelHelper, Channel2, model.Snowflake(userID), !inFavor)
	}
}

func TestCastBallot_InFavorFailsOnExpiredVote(t *testing.T) {
	modelHelper, clock := initializeTests()

	assertStartNewVote(t, modelHelper, Channel1, UserID1)
	clock.Advance(vote.VoteDuration)
	assertCannotVoteWhenExpired(t, modelHelper, Channel1, UserID1, true)
}

func TestCastBallot_OpposedFailsOnExpiredVote(t *testing.T) {
	modelHelper, clock := initializeTests()

	assertStartNewVote(t, modelHelper, Channel1, UserID1)
	clock.Advance(vote.VoteDuration)
	assertCannotVoteWhenExpired(t, modelHelper, Channel1, UserID1, false)
}

func TestSetVoteOutcome_NotStarted(t *testing.T) {
	modelHelper, _ := initializeTests()

	assertEarlySetVoteOutcomeFails(t, modelHelper, Channel1, model.VoteOutcomePassed)
	assertEarlySetVoteOutcomeFails(t, modelHelper, Channel1, model.VoteOutcomeFailed)
	assertEarlySetVoteOutcomeFails(t, modelHelper, Channel1, model.VoteOutcomeNotEnough)
	assertEarlySetVoteOutcomeFails(t, modelHelper, Channel1, model.VoteOutcomeNotDone)
}

func TestSetVoteOutcome_Passed(t *testing.T) {
	modelHelper, _ := initializeTests()

	assertStartNewVote(t, modelHelper, Channel1, UserID1)
	assertSetVoteOutcome(t, modelHelper, Channel1, model.VoteOutcomePassed)
	assertSetVoteOutcomeFails(t, modelHelper, Channel1, model.VoteOutcomeFailed)
	assertSetVoteOutcomeFails(t, modelHelper, Channel1, model.VoteOutcomeNotEnough)
	assertSetVoteOutcomeFails(t, modelHelper, Channel1, model.VoteOutcomeNotDone)
}

func TestSetVoteOutcome_Passed_2Channels(t *testing.T) {
	modelHelper, _ := initializeTests()

	assertStartNewVote(t, modelHelper, Channel1, UserID1)
	assertStartNewVote(t, modelHelper, Channel2, UserID1)
	assertSetVoteOutcome(t, modelHelper, Channel1, model.VoteOutcomePassed)
	assertSetVoteOutcome(t, modelHelper, Channel2, model.VoteOutcomePassed)
}

func TestSetVoteOutcome_Failed(t *testing.T) {
	modelHelper, _ := initializeTests()

	assertStartNewVote(t, modelHelper, Channel1, UserID1)
	assertSetVoteOutcome(t, modelHelper, Channel1, model.VoteOutcomeFailed)
	assertSetVoteOutcomeFails(t, modelHelper, Channel1, model.VoteOutcomePassed)
	assertSetVoteOutcomeFails(t, modelHelper, Channel1, model.VoteOutcomeNotEnough)
	assertSetVoteOutcomeFails(t, modelHelper, Channel1, model.VoteOutcomeNotDone)
}

func TestSetVoteOutcome_NotEnough(t *testing.T) {
	modelHelper, _ := initializeTests()

	assertStartNewVote(t, modelHelper, Channel1, UserID1)
	assertSetVoteOutcome(t, modelHelper, Channel1, model.VoteOutcomeNotEnough)
	assertSetVoteOutcomeFails(t, modelHelper, Channel1, model.VoteOutcomeFailed)
	assertSetVoteOutcomeFails(t, modelHelper, Channel1, model.VoteOutcomePassed)
	assertSetVoteOutcomeFails(t, modelHelper, Channel1, model.VoteOutcomeNotDone)
}

func TestSetVoteOutcome_NotDoneCanBeOverwritten(t *testing.T) {
	modelHelper, _ := initializeTests()

	assertStartNewVote(t, modelHelper, Channel1, UserID1)
	assertSetVoteOutcome(t, modelHelper, Channel1, model.VoteOutcomeNotDone)
	assertSetVoteOutcome(t, modelHelper, Channel1, model.VoteOutcomePassed)
	assertSetVoteOutcomeFails(t, modelHelper, Channel1, model.VoteOutcomeNotDone)
}

func TestSetVoteOutcome_SucceedsWhenExpired(t *testing.T) {
	modelHelper, clock := initializeTests()

	assertStartNewVote(t, modelHelper, Channel1, UserID1)
	clock.Advance(vote.VoteDuration)
	assertSetVoteOutcome(t, modelHelper, Channel1, model.VoteOutcomePassed)
}

func initializeTests() (*vote.ModelHelper, *testutil.FakeUTCClock) {
	stringMap := testutil.NewInMemoryStringMap()
	clock := testutil.NewFakeUTCClock()
	return vote.NewModelHelper(stringMap, clock), clock
}

func assertIsVoteActive(t *testing.T, modelHelper *vote.ModelHelper, channel model.Snowflake, active bool) {
	t.Helper()

	ok, err := modelHelper.IsVoteActive(channel)
	if err != nil {
		t.Errorf("Should not have errored")
	}
	if ok != active {
		t.Errorf("Input mismatch")
	}
}

func assertMostRecentVote(t *testing.T, modelHelper *vote.ModelHelper, channelID model.Snowflake, vote *model.Vote) {
	t.Helper()

	mostRecentVote, err := modelHelper.MostRecentVote(channelID)
	if err != nil {
		t.Errorf("Should not have errored pulling most-recent vote: %v", err)
	}
	if !reflect.DeepEqual(mostRecentVote, vote) {
		t.Errorf("Stored vote is not equal to most recent vote")
	}
}

func assertMostRecentVoteID(t *testing.T, modelHelper *vote.ModelHelper, channelID model.Snowflake, voteID int) {
	t.Helper()

	mostRecentVoteID, err := modelHelper.MostRecentVoteID(channelID)
	if err != nil {
		t.Errorf("Should not have errored pulling most-recent vote ID: %v", err)
	}
	if mostRecentVoteID != voteID {
		t.Errorf("Stored vote ID not equal to most recent vote ID")
	}
}

func assertStartNewVoteFails(t *testing.T, modelHelper *vote.ModelHelper, channelID, userID model.Snowflake) {
	t.Helper()

	_, err := modelHelper.StartNewVote(channelID, userID, "hug Jake")
	if err != vote.ErrorOnlyOneVote {
		t.Errorf("Should have failed to add vote with ID %v", userID)
	}
}

func assertStartNewVote(t *testing.T, modelHelper *vote.ModelHelper, channelID, userID model.Snowflake) *model.Vote {
	t.Helper()

	vote, err := modelHelper.StartNewVote(channelID, userID, "hug Jake")
	if err != nil {
		t.Errorf("Should have started a vote with ID %v", userID)
	}
	return vote
}

func assertEarlyBallotFails(t *testing.T, modelHelper *vote.ModelHelper, channelID, userID model.Snowflake, inFavor bool) {
	t.Helper()

	_, err := modelHelper.CastBallot(channelID, userID, inFavor)
	if err != vote.ErrorNoVoteActive {
		t.Errorf("Should have failed due to no active vote")
	}
}

func assertCastBallot(t *testing.T, modelHelper *vote.ModelHelper, channelID, userID model.Snowflake, inFavor bool) {
	t.Helper()

	vote, err := modelHelper.CastBallot(channelID, userID, inFavor)
	if err != nil {
		t.Errorf("Unexpected ballot failure")
	}

	// Hunt for the userID in the cast votes.
	found := false
	voteSide := vote.VotesFor
	if !inFavor {
		voteSide = vote.VotesAgainst
	}
	for _, val := range voteSide {
		if val == userID {
			found = true
		}
	}

	if !found {
		t.Errorf("Vote did not successfully happen")
	}

	// Assert that the vote was serialized.
	assertMostRecentVote(t, modelHelper, channelID, vote)
}

func assertCannotVoteAgain(t *testing.T, modelHelper *vote.ModelHelper, channelID, userID model.Snowflake, inFavor bool) {
	t.Helper()

	_, err := modelHelper.CastBallot(channelID, userID, inFavor)
	if err != vote.ErrorAlreadyVoted {
		t.Errorf("Expected ballot to fail: %v", err)
	}
}

func assertCannotVoteWhenExpired(t *testing.T, modelHelper *vote.ModelHelper, channelID, userID model.Snowflake, inFavor bool) {
	t.Helper()

	_, err := modelHelper.CastBallot(channelID, userID, inFavor)
	if err != vote.ErrorNoVoteActive {
		t.Errorf("Expected ballot to fail: %v", err)
	}
}

func assertEarlySetVoteOutcomeFails(t *testing.T, modelHelper *vote.ModelHelper, channelID model.Snowflake, voteOutcome int) {
	t.Helper()

	err := modelHelper.SetVoteOutcome(channelID, voteOutcome)
	if err != vote.ErrorNoVoteActive {
		t.Errorf("Expected vote to already have an outcome")
	}
}

func assertSetVoteOutcome(t *testing.T, modelHelper *vote.ModelHelper, channelID model.Snowflake, voteOutcome int) {
	t.Helper()

	err := modelHelper.SetVoteOutcome(channelID, voteOutcome)
	if err != nil {
		t.Errorf("Unexpected ballot failure")
	}

	vote, err := modelHelper.MostRecentVote(channelID)
	if err != nil {
		t.Fatalf("Unexpected failure %v", err)
	}

	if vote.VoteOutcome != voteOutcome {
		t.Errorf("Unexpected stored vote outcome %v", vote.VoteOutcome)
	}
}

func assertSetVoteOutcomeFails(t *testing.T, modelHelper *vote.ModelHelper, channelID model.Snowflake, voteOutcome int) {
	t.Helper()

	err := modelHelper.SetVoteOutcome(channelID, voteOutcome)
	if err != vote.ErrorVoteHasOutcome {
		t.Errorf("Expected failure due to existing outcome: %v", err)
	}
}
