package testutil

import "time"

// A DDO that stores a function, and the duration until it should fire.
type timedCallback struct {
	Duration time.Duration
	CB       func()
}

// FakeUTCTimer is a faked UTC timer for testing.
type FakeUTCTimer struct {
	timedCBs []timedCallback
}

// NewFakeUTCTimer works as advertised.
func NewFakeUTCTimer() *FakeUTCTimer {
	return &FakeUTCTimer{
		timedCBs: []timedCallback{},
	}
}

// ExecuteAfter stores the callback and the duration. If the duration is 0 or
// negative, it will fire immediately.
func (c *FakeUTCTimer) ExecuteAfter(d time.Duration, cb func()) {
	if d <= time.Duration(0) {
		cb()
		return
	}
	c.timedCBs = append(c.timedCBs, timedCallback{d, cb})
}

// ElapseTime subtracts duration from each stored callback. For each duration
// that becomes negative, it is removed from the list and the associated
// callback is called.
func (c *FakeUTCTimer) ElapseTime(duration time.Duration) {
	result := []timedCallback{}
	for _, t := range c.timedCBs {
		t.Duration = t.Duration - duration
		if t.Duration <= time.Duration(0) {
			t.CB()
			continue
		}
		result = append(result, t)
	}
	c.timedCBs = result
}
