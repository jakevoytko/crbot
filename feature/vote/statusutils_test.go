package vote

import (
	"testing"
	"time"

	testutil "github.com/jakevoytko/crbot/testutil"
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
