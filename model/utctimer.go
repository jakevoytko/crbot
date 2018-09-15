package model

import "time"

// UTCTimer is an interface that allows a closure to be executed after a certain
// amount of time. This is non-blocking.
type UTCTimer interface {
	ExecuteAfter(time.Duration, func())
}

// SystemUTCTimer is a real implementation of UTCTimer.
type SystemUTCTimer struct{}

// NewSystemUTCTimer works as advertised.
func NewSystemUTCTimer() *SystemUTCTimer {
	return &SystemUTCTimer{}
}

// ExecuteAfter executes the given closure after the given duration has
// elapsed. This is non-blocking.
func (c *SystemUTCTimer) ExecuteAfter(howLong time.Duration, cb func()) {
	go func() {
		timer := time.NewTimer(howLong)
		<-timer.C
		cb()
	}()
}
