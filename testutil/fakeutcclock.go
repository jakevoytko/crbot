package testutil

import "time"

// FakeUTClock is a mockable UTC clock for testing.
type FakeUTCClock struct {
	currentTime time.Time
}

// NewFakeUTCClock works as advertised.
func NewFakeUTCClock() *FakeUTCClock {
	return &FakeUTCClock{
		currentTime: time.Date(
			2017,
			time.January,
			1, /* day */
			1, /* hour */
			1, /* minute */
			0, /* second */
			0, /* nsec */
			time.UTC),
	}
}

// Now returns the mocked UTC time.
func (c *FakeUTCClock) Now() time.Time {
	return c.currentTime
}

// Advance advances the internal clock by the given duration.
func (c *FakeUTCClock) Advance(d time.Duration) {
	c.currentTime = c.currentTime.Add(d)
}
