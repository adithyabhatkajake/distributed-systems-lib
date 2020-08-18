package util

/* References */
/* https://gobyexample.com/timers */
/* https://gobyexample.com/tickers */

import (
	"time" // https://golang.org/pkg/time/
)

// Timer is an extension of golang's timer, that allows cancelling
type Timer struct {
	callable func()
	dur      time.Duration
	cancel   chan struct{}
}

// NewTimer returns a new timer with the given callback function
func NewTimer(call func()) *Timer {
	t := &Timer{callable: call}
	t.cancel = make(chan struct{})
	return t
}

// SetTime sets the waiting time for the timer
func (t *Timer) SetTime(dur time.Duration) {
	t.dur = dur
}

// Start begins the countdown for the timer
func (t *Timer) Start() {
	// Start a thread that waits for the timer to finish
	go func() {
		// Wait for timer to finish or if cancelled
		select {
		case <-time.After(t.dur):
			t.callable()
		case <-t.cancel:
			return
		}
	}()
}

// Cancel cancels the timer
func (t *Timer) Cancel() {
	var a struct{}
	t.cancel <- a
}
