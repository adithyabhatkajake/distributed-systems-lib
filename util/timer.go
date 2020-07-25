package util

/* References */
/* https://gobyexample.com/timers */
/* https://gobyexample.com/tickers */

import (
	"time" // https://golang.org/pkg/time/
)

// This is the duration delta used in our consensus protocol
var delta time.Duration

// Set the delta used in our protocol
func InitDelta(t time.Duration) {
	delta = t
}

// Return a timer for Delta
func DeltaTimer() *time.Timer {
	return time.NewTimer(delta)
}

// Return a timer for 2Delta
func DoubleDeltaTimer() *time.Timer {
	return time.NewTimer(2*delta)
}

// Cancel the timer set for whatever time it was set
func CancelTimer(t *time.Timer) {
	t.Stop()
}
