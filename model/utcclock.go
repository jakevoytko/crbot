package model

import "time"

// UTCClock is an interface that returns a UTC time.
type UTCClock interface {
	Now() time.Time
}

// SystemUTCClock is a UTCClock that returns UTC based on the current system's
// time.
type SystemUTCClock struct{}

// NewSystemUTCClock works as advertised.
func NewSystemUTCClock() *SystemUTCClock {
	return &SystemUTCClock{}
}

// Now returns the UTC time for the current system.
func (c *SystemUTCClock) Now() time.Time {
	return time.Now().UTC()
}
