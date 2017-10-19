package vote

import (
	"reflect"
	"testing"
	"time"

	"github.com/jakevoytko/crbot/util"
)

func TestIsVoteActive_NoMemory(t *testing.T) {
	modelHelper, _ := initializeTests()

	assertIsVoteActive(t, modelHelper, false)
}

func TestIsVoteActive_ActiveVote(t *testing.T) {
	modelHelper, _ := initializeTests()

	assertStartNewVote(t, modelHelper, 8675309 /* userID */)
	assertIsVoteActive(t, modelHelper, true)
}

func TestIsVoteActive_TimeExpires(t *testing.T) {
	modelHelper, clock := initializeTests()

	assertStartNewVote(t, modelHelper, 8675309 /* userID */)
	assertIsVoteActive(t, modelHelper, true)

	// Advance the clock 1ns before the vote will expire, and then verify that
	// advancing it 1ns causes the vote to expire.
	clock.Advance(VoteDuration - time.Duration(1)*time.Nanosecond)
	assertIsVoteActive(t, modelHelper, true)
	clock.Advance(time.Duration(1) * time.Nanosecond)
	assertIsVoteActive(t, modelHelper, false)
}

func TestMostRecentVote_NoMemory(t *testing.T) {
	modelHelper, _ := initializeTests()

	assertMostRecentVote(t, modelHelper, nil)
}

func TestMostRecentVote_ActiveVote(t *testing.T) {
	modelHelper, _ := initializeTests()

	newVote := assertStartNewVote(t, modelHelper, 8675309 /* userID */)
	assertMostRecentVote(t, modelHelper, newVote)
}

func TestMostRecentVote_ReturnsExpiredVote(t *testing.T) {
	modelHelper, clock := initializeTests()

	newVote := assertStartNewVote(t, modelHelper, 8675309 /* userID */)
	assertIsVoteActive(t, modelHelper, true)

	// Check that the correct vote is still returned, even though it has expired.
	clock.Advance(VoteDuration)
	assertIsVoteActive(t, modelHelper, false)
	assertMostRecentVote(t, modelHelper, newVote)
}

func TestMostRecentVoteID_NoMemory(t *testing.T) {
	modelHelper, _ := initializeTests()

	assertMostRecentVoteID(t, modelHelper, 0)
}

func TestMostRecentVoteID_ActiveVote(t *testing.T) {
	modelHelper, _ := initializeTests()

	newVote := assertStartNewVote(t, modelHelper, 8675309 /* userID */)
	assertMostRecentVoteID(t, modelHelper, newVote.VoteID)
}

func TestMostRecentVoteID_ReturnsExpiredVote(t *testing.T) {
	modelHelper, clock := initializeTests()

	newVote := assertStartNewVote(t, modelHelper, 8675309 /* userID */)
	clock.Advance(VoteDuration)
	assertMostRecentVoteID(t, modelHelper, newVote.VoteID)
}

func TestStartNewVote_AddingFromNoMemory(t *testing.T) {
	modelHelper, _ := initializeTests()

	vote := assertStartNewVote(t, modelHelper, 8675309)

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
	if vote.VoteOutcome != VoteOutcomeNotDone {
		t.Errorf("Vote outcome should be NotDone")
	}

	assertMostRecentVote(t, modelHelper, vote)

	assertMostRecentVoteID(t, modelHelper, vote.VoteID)
}

func TestStartNewVote_DoesNotOverrideExistingVote(t *testing.T) {
	modelHelper, _ := initializeTests()

	assertStartNewVote(t, modelHelper, 8675309 /* userID */)
	assertStartNewVoteFails(t, modelHelper, 11)
}

func TestCastBallot_NoVoteActive(t *testing.T) {
	modelHelper, _ := initializeTests()

	assertEarlyBallotFails(t, modelHelper, 1 /* userID */, true /* inFavor */)
	assertEarlyBallotFails(t, modelHelper, 1 /* userID */, false /* inFavor */)
}

func TestCastBallot_VoteFor(t *testing.T) {
	modelHelper, _ := initializeTests()

	assertStartNewVote(t, modelHelper, 1 /* userID */)
	assertCastBallot(t, modelHelper, 1 /* userID */, true)
}

func TestCastBallot_VoteAgainst(t *testing.T) {
	modelHelper, _ := initializeTests()

	assertStartNewVote(t, modelHelper, 1 /* userID */)
	assertCastBallot(t, modelHelper, 1 /* userID */, false)
}

func TestCastBallot_CanNotVoteTwiceAfterInFavor(t *testing.T) {
	modelHelper, _ := initializeTests()

	assertStartNewVote(t, modelHelper, 1 /* userID */)
	assertCastBallot(t, modelHelper, 1 /* userID */, true)
	assertCannotVoteAgain(t, modelHelper, 1 /* userID */, true)
	assertCannotVoteAgain(t, modelHelper, 1 /* userID */, false)
}

func TestCastBallot_CanNotVoteTwiceAfterOpposed(t *testing.T) {
	modelHelper, _ := initializeTests()

	assertStartNewVote(t, modelHelper, 1 /* userID */)
	assertCastBallot(t, modelHelper, 1 /* userID */, false)
	assertCannotVoteAgain(t, modelHelper, 1 /* userID */, true)
	assertCannotVoteAgain(t, modelHelper, 1 /* userID */, false)
}

func TestCastBallot_LotsOfVotes(t *testing.T) {
	modelHelper, _ := initializeTests()

	assertStartNewVote(t, modelHelper, 1 /* userID */)
	for userID := 0; userID < 10; userID++ {
		inFavor := userID%2 == 0
		assertCastBallot(t, modelHelper, int64(userID), inFavor)
	}
}

func TestCastBallot_InFavorFailsOnExpiredVote(t *testing.T) {
	modelHelper, clock := initializeTests()

	assertStartNewVote(t, modelHelper, 8675309 /* userID */)
	clock.Advance(VoteDuration)
	assertCannotVoteWhenExpired(t, modelHelper, 1 /* userID */, true)
}

func TestCastBallot_OpposedFailsOnExpiredVote(t *testing.T) {
	modelHelper, clock := initializeTests()

	assertStartNewVote(t, modelHelper, 8675309 /* userID */)
	clock.Advance(VoteDuration)
	assertCannotVoteWhenExpired(t, modelHelper, 1 /* userID */, false)
}

func TestSetVoteOutcome_NotStarted(t *testing.T) {
	modelHelper, _ := initializeTests()

	assertEarlySetVoteOutcomeFails(t, modelHelper, VoteOutcomePassed)
	assertEarlySetVoteOutcomeFails(t, modelHelper, VoteOutcomeFailed)
	assertEarlySetVoteOutcomeFails(t, modelHelper, VoteOutcomeNotEnough)
	assertEarlySetVoteOutcomeFails(t, modelHelper, VoteOutcomeNotDone)
}

func TestSetVoteOutcome_Passed(t *testing.T) {
	modelHelper, _ := initializeTests()

	assertStartNewVote(t, modelHelper, 1 /* userID */)
	assertSetVoteOutcome(t, modelHelper, VoteOutcomePassed)
	assertSetVoteOutcomeFails(t, modelHelper, VoteOutcomeFailed)
	assertSetVoteOutcomeFails(t, modelHelper, VoteOutcomeNotEnough)
	assertSetVoteOutcomeFails(t, modelHelper, VoteOutcomeNotDone)
}

func TestSetVoteOutcome_Failed(t *testing.T) {
	modelHelper, _ := initializeTests()

	assertStartNewVote(t, modelHelper, 1 /* userID */)
	assertSetVoteOutcome(t, modelHelper, VoteOutcomeFailed)
	assertSetVoteOutcomeFails(t, modelHelper, VoteOutcomePassed)
	assertSetVoteOutcomeFails(t, modelHelper, VoteOutcomeNotEnough)
	assertSetVoteOutcomeFails(t, modelHelper, VoteOutcomeNotDone)
}

func TestSetVoteOutcome_NotEnough(t *testing.T) {
	modelHelper, _ := initializeTests()

	assertStartNewVote(t, modelHelper, 1 /* userID */)
	assertSetVoteOutcome(t, modelHelper, VoteOutcomeNotEnough)
	assertSetVoteOutcomeFails(t, modelHelper, VoteOutcomeFailed)
	assertSetVoteOutcomeFails(t, modelHelper, VoteOutcomePassed)
	assertSetVoteOutcomeFails(t, modelHelper, VoteOutcomeNotDone)
}

func TestSetVoteOutcome_NotDoneCanBeOverwritten(t *testing.T) {
	modelHelper, _ := initializeTests()

	assertStartNewVote(t, modelHelper, 1 /* userID */)
	assertSetVoteOutcome(t, modelHelper, VoteOutcomeNotDone)
	assertSetVoteOutcome(t, modelHelper, VoteOutcomePassed)
	assertSetVoteOutcomeFails(t, modelHelper, VoteOutcomeNotDone)
}

func TestSetVoteOutcome_SucceedsWhenExpired(t *testing.T) {
	modelHelper, clock := initializeTests()

	assertStartNewVote(t, modelHelper, 1 /* userID */)
	clock.Advance(VoteDuration)
	assertSetVoteOutcome(t, modelHelper, VoteOutcomePassed)
}

func initializeTests() (*ModelHelper, *util.FakeUTCClock) {
	stringMap := util.NewInMemoryStringMap()
	clock := util.NewFakeUTCClock()
	return NewModelHelper(stringMap, clock), clock
}

func assertIsVoteActive(t *testing.T, modelHelper *ModelHelper, active bool) {
	ok, err := modelHelper.IsVoteActive()
	if err != nil {
		t.Errorf("Should not have errored")
	}
	if ok != active {
		t.Errorf("Input mismatch")
	}
}

func assertMostRecentVote(t *testing.T, modelHelper *ModelHelper, vote *Vote) {
	mostRecentVote, err := modelHelper.MostRecentVote()
	if err != nil {
		t.Errorf("Should not have errored pulling most-recent vote: %v", err)
	}
	if !reflect.DeepEqual(mostRecentVote, vote) {
		t.Errorf("Stored vote is not equal to most recent vote")
	}
}

func assertMostRecentVoteID(t *testing.T, modelHelper *ModelHelper, voteID int) {
	mostRecentVoteID, err := modelHelper.MostRecentVoteID()
	if err != nil {
		t.Errorf("Should not have errored pulling most-recent vote ID: %v", err)
	}
	if mostRecentVoteID != voteID {
		t.Errorf("Stored vote ID not equal to most recent vote ID")
	}
}

func assertStartNewVoteFails(t *testing.T, modelHelper *ModelHelper, userID int64) {
	_, err := modelHelper.StartNewVote(userID, "hug Jake")
	if err != ErrorOnlyOneVote {
		t.Errorf("Should have failed to add vote with ID %v", userID)
	}
}

func assertStartNewVote(t *testing.T, modelHelper *ModelHelper, userID int64) *Vote {
	vote, err := modelHelper.StartNewVote(userID, "hug Jake")
	if err != nil {
		t.Errorf("Should have started a vote with ID %v", userID)
	}
	return vote
}

func assertEarlyBallotFails(t *testing.T, modelHelper *ModelHelper, userID int64, inFavor bool) {
	_, err := modelHelper.CastBallot(userID, inFavor)
	if err != ErrorNoVoteActive {
		t.Errorf("Should have failed due to no active vote")
	}
}

func assertCastBallot(t *testing.T, modelHelper *ModelHelper, userID int64, inFavor bool) {
	vote, err := modelHelper.CastBallot(userID, inFavor)
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
	assertMostRecentVote(t, modelHelper, vote)
}

func assertCannotVoteAgain(t *testing.T, modelHelper *ModelHelper, userID int64, inFavor bool) {
	_, err := modelHelper.CastBallot(userID, inFavor)
	if err != ErrorAlreadyVoted {
		t.Errorf("Expected ballot to fail: %v", err)
	}
}

func assertCannotVoteWhenExpired(t *testing.T, modelHelper *ModelHelper, userID int64, inFavor bool) {
	_, err := modelHelper.CastBallot(userID, inFavor)
	if err != ErrorNoVoteActive {
		t.Errorf("Expected ballot to fail: %v", err)
	}
}

func assertEarlySetVoteOutcomeFails(t *testing.T, modelHelper *ModelHelper, voteOutcome int) {
	err := modelHelper.SetVoteOutcome(voteOutcome)
	if err != ErrorNoVoteActive {
		t.Errorf("Expected vote to already have an outcome")
	}
}

func assertSetVoteOutcome(t *testing.T, modelHelper *ModelHelper, voteOutcome int) {
	err := modelHelper.SetVoteOutcome(voteOutcome)
	if err != nil {
		t.Errorf("Unexpected ballot failure")
	}

	vote, err := modelHelper.MostRecentVote()
	if err != nil {
		t.Fatalf("Unexpected failure %v", err)
	}

	if vote.VoteOutcome != voteOutcome {
		t.Errorf("Unexpected stored vote outcome %v", vote.VoteOutcome)
	}
}

func assertSetVoteOutcomeFails(t *testing.T, modelHelper *ModelHelper, voteOutcome int) {
	err := modelHelper.SetVoteOutcome(voteOutcome)
	if err != ErrorVoteHasOutcome {
		t.Errorf("Expected failure due to existing outcome: %v", err)
	}
}
