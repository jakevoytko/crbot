package vote

import (
	"testing"
	"time"

	"github.com/jakevoytko/crbot/feature/vote"
	"github.com/jakevoytko/crbot/model"
	"github.com/jakevoytko/crbot/testutil"
	stringmap "github.com/jakevoytko/go-stringmap"
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
		actual := vote.TimeString(utcClock, testTime)
		if expected != actual {
			t.Errorf("Got %v expected %v for time %v", actual, expected, testTime)
		}
	}
}

// TestRunner isn't set up to correctly test initial load, so the tests are
// being done manually here.
func TestHandleVotesOnInitialLoad_EmptyMap(t *testing.T) {
	session := testutil.NewInMemoryDiscordSession()
	stringMap := stringmap.NewInMemoryStringMap()
	timer := testutil.NewFakeUTCTimer()
	clock := testutil.NewFakeUTCClock()
	modelHelper := vote.NewModelHelper(stringMap, clock)
	commandChannel := make(chan *model.Command, 10)

	vote.HandleVotesOnInitialLoad(session, modelHelper, clock, timer, commandChannel)

	timer.ElapseTime(vote.VoteDuration)
	clock.Advance(vote.VoteDuration)

	select {
	case <-commandChannel:
		t.Errorf("Channel should have been empty")
	default:
	}
}

func TestHandleVotesOnInitialLoad_HalfDoneVote(t *testing.T) {
	session := testutil.NewInMemoryDiscordSession()
	stringMap := stringmap.NewInMemoryStringMap()
	timer := testutil.NewFakeUTCTimer()
	clock := testutil.NewFakeUTCClock()
	modelHelper := vote.NewModelHelper(stringMap, clock)
	commandChannel := make(chan *model.Command, 10)

	modelHelper.StartNewVote(model.Snowflake(1) /* channelID */, model.Snowflake(2) /* userID */, "oh noes")

	timer.ElapseTime(vote.VoteDuration / 2)
	clock.Advance(vote.VoteDuration / 2)

	vote.HandleVotesOnInitialLoad(session, modelHelper, clock, timer, commandChannel)

	// Assert the channel is now empty.
	select {
	case <-commandChannel:
		t.Errorf("Channel should have been empty")
	default:
	}

	// Cause the new timer to fire.
	timer.ElapseTime(vote.VoteDuration)
	clock.Advance(vote.VoteDuration)

	// Assert the command currently in the channel.
	select {
	case command, ok := <-commandChannel:
		if !ok {
			t.Errorf("Should have gotten a command from the command channel")
		}
		if command.Type != model.CommandTypeVoteConclude {
			t.Errorf("Expected a conclude vote")
		}
	default:
		t.Errorf("Channel should have not been empty")
	}

	// Assert the channel is now empty.
	select {
	case <-commandChannel:
		t.Errorf("Channel should have been empty")
	default:
	}
}

func TestHandleVotesOnInitialLoad_VoteExpired(t *testing.T) {
	session := testutil.NewInMemoryDiscordSession()
	stringMap := stringmap.NewInMemoryStringMap()
	timer := testutil.NewFakeUTCTimer()
	clock := testutil.NewFakeUTCClock()
	modelHelper := vote.NewModelHelper(stringMap, clock)
	commandChannel := make(chan *model.Command, 10)

	modelHelper.StartNewVote(model.Snowflake(1) /* channelID */, model.Snowflake(2) /* userID */, "oh noes")

	timer.ElapseTime(vote.VoteDuration)
	clock.Advance(vote.VoteDuration)

	vote.HandleVotesOnInitialLoad(session, modelHelper, clock, timer, commandChannel)

	// Assert the command currently in the channel.
	select {
	case command, ok := <-commandChannel:
		if !ok {
			t.Errorf("Should have gotten a command from the command channel")
		}
		if command.Type != model.CommandTypeVoteConclude {
			t.Errorf("Expected a conclude vote")
		}
	default:
		t.Errorf("Channel should have not been empty")
	}

	// Assert the channel is now empty.
	select {
	case <-commandChannel:
		t.Errorf("Channel should have been empty")
	default:
	}
}

func TestHandleVotesOnInitialLoad_PreviousExpriedVotesDoNotCauseNewConcludes(t *testing.T) {
	session := testutil.NewInMemoryDiscordSession()
	stringMap := stringmap.NewInMemoryStringMap()
	timer := testutil.NewFakeUTCTimer()
	clock := testutil.NewFakeUTCClock()
	modelHelper := vote.NewModelHelper(stringMap, clock)
	commandChannel := make(chan *model.Command, 10)

	modelHelper.StartNewVote(model.Snowflake(1) /* channelID */, model.Snowflake(2) /* userID */, "oh noes")

	timer.ElapseTime(vote.VoteDuration)
	clock.Advance(vote.VoteDuration)

	modelHelper.SetVoteOutcome(model.Snowflake(1) /* channelID */, model.VoteOutcomeNotEnough)

	vote.HandleVotesOnInitialLoad(session, modelHelper, clock, timer, commandChannel)

	timer.ElapseTime(vote.VoteDuration)
	clock.Advance(vote.VoteDuration)

	// Assert the channel is now empty.
	select {
	case <-commandChannel:
		t.Errorf("Channel should have been empty")
	default:
	}
}
