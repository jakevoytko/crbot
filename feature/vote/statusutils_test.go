package vote

import (
	"testing"
	"time"

	"github.com/jakevoytko/crbot/model"
	"github.com/jakevoytko/crbot/testutil"
)

func TestTimeString(t *testing.T) {
	utcClock := testutil.NewFakeUTCClock()

	testCases := []struct {
		InFuture        time.Duration
		ExpectedMessage string
	}{
		{time.Duration(30) * time.Minute, "30 minutes remaining"},
		{time.Duration(30)*time.Minute - time.Nanosecond, "30 minutes remaining"},
		{time.Duration(29)*time.Minute + time.Nanosecond, "30 minutes remaining"},
		{time.Duration(29) * time.Minute, "29 minutes remaining"},
		{time.Duration(1)*time.Minute + time.Nanosecond, "2 minutes remaining"},
		{time.Duration(1) * time.Minute, "60 seconds remaining"},
		{time.Duration(1)*time.Second + time.Nanosecond, "2 seconds remaining"},
		{time.Duration(1) * time.Second, "1000 milliseconds remaining"},
		{time.Duration(2) * time.Millisecond, "2 milliseconds remaining"},
		{time.Duration(1) * time.Millisecond, "No time remaining in vote"},
		{time.Duration(0), "No time remaining in vote"},
		{time.Duration(-1), "No time remaining in vote"},
	}

	for _, testCase := range testCases {
		testTime := utcClock.Now().Add(testCase.InFuture)
		expected := testCase.ExpectedMessage
		actual := TimeString(utcClock, testTime)
		if expected != actual {
			t.Errorf("Got %v expected %v for time %v", actual, expected, testTime)
		}
	}
}

// TestRunner isn't set up to correctly test initial load, so the tests are
// being done manually here.
func TestHandleVotesOnInitialLoad_EmptyMap(t *testing.T) {
	session := testutil.NewInMemoryDiscordSession()
	stringMap := testutil.NewInMemoryStringMap()
	timer := testutil.NewFakeUTCTimer()
	clock := testutil.NewFakeUTCClock()
	modelHelper := NewModelHelper(stringMap, clock)
	commandChannel := make(chan *model.Command, 10)

	handleVotesOnInitialLoad(session, modelHelper, clock, timer, commandChannel)

	timer.ElapseTime(VoteDuration)
	clock.Advance(VoteDuration)

	select {
	case _, _ = <-commandChannel:
		t.Errorf("Channel should have been empty")
	default:
	}
}

func TestHandleVotesOnInitialLoad_HalfDoneVote(t *testing.T) {
	session := testutil.NewInMemoryDiscordSession()
	stringMap := testutil.NewInMemoryStringMap()
	timer := testutil.NewFakeUTCTimer()
	clock := testutil.NewFakeUTCClock()
	modelHelper := NewModelHelper(stringMap, clock)
	commandChannel := make(chan *model.Command, 10)

	modelHelper.StartNewVote(model.Snowflake(1) /* channelID */, model.Snowflake(2) /* userID */, "oh noes")

	timer.ElapseTime(VoteDuration / 2)
	clock.Advance(VoteDuration / 2)

	handleVotesOnInitialLoad(session, modelHelper, clock, timer, commandChannel)

	// Assert the channel is now empty.
	select {
	case _, _ = <-commandChannel:
		t.Errorf("Channel should have been empty")
	default:
	}

	// Cause the new timer to fire.
	timer.ElapseTime(VoteDuration)
	clock.Advance(VoteDuration)

	// Assert the command currently in the channel.
	select {
	case command, ok := <-commandChannel:
		if !ok {
			t.Errorf("Should have gotten a command from the command channel")
		}
		if command.Type != model.Type_VoteConclude {
			t.Errorf("Expected a conclude vote")
		}
	default:
		t.Errorf("Channel should have not been empty")
	}

	// Assert the channel is now empty.
	select {
	case _, _ = <-commandChannel:
		t.Errorf("Channel should have been empty")
	default:
	}
}

func TestHandleVotesOnInitialLoad_VoteExpired(t *testing.T) {
	session := testutil.NewInMemoryDiscordSession()
	stringMap := testutil.NewInMemoryStringMap()
	timer := testutil.NewFakeUTCTimer()
	clock := testutil.NewFakeUTCClock()
	modelHelper := NewModelHelper(stringMap, clock)
	commandChannel := make(chan *model.Command, 10)

	modelHelper.StartNewVote(model.Snowflake(1) /* channelID */, model.Snowflake(2) /* userID */, "oh noes")

	timer.ElapseTime(VoteDuration)
	clock.Advance(VoteDuration)

	handleVotesOnInitialLoad(session, modelHelper, clock, timer, commandChannel)

	// Assert the command currently in the channel.
	select {
	case command, ok := <-commandChannel:
		if !ok {
			t.Errorf("Should have gotten a command from the command channel")
		}
		if command.Type != model.Type_VoteConclude {
			t.Errorf("Expected a conclude vote")
		}
	default:
		t.Errorf("Channel should have not been empty")
	}

	// Assert the channel is now empty.
	select {
	case _, _ = <-commandChannel:
		t.Errorf("Channel should have been empty")
	default:
	}
}

func TestHandleVotesOnInitialLoad_PreviousExpriedVotesDoNotCauseNewConcludes(t *testing.T) {
	session := testutil.NewInMemoryDiscordSession()
	stringMap := testutil.NewInMemoryStringMap()
	timer := testutil.NewFakeUTCTimer()
	clock := testutil.NewFakeUTCClock()
	modelHelper := NewModelHelper(stringMap, clock)
	commandChannel := make(chan *model.Command, 10)

	modelHelper.StartNewVote(model.Snowflake(1) /* channelID */, model.Snowflake(2) /* userID */, "oh noes")

	timer.ElapseTime(VoteDuration)
	clock.Advance(VoteDuration)

	modelHelper.SetVoteOutcome(model.Snowflake(1) /* channelID */, model.VoteOutcomeNotEnough)

	handleVotesOnInitialLoad(session, modelHelper, clock, timer, commandChannel)

	timer.ElapseTime(VoteDuration)
	clock.Advance(VoteDuration)

	// Assert the channel is now empty.
	select {
	case _, _ = <-commandChannel:
		t.Errorf("Channel should have been empty")
	default:
	}
}
