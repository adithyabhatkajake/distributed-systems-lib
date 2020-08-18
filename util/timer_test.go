package util_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/adithyabhatkajake/libe2c/util"
)

// Check if the correct callback is called after the time
func TestTimerWaiting(t *testing.T) {
	success := false
	timer := util.NewTimer(func() {
		success = true
	})
	timer.SetTime(time.Second * 4)
	timer.Start()
	<-time.After(time.Second * 4)
	// By now success should be true
	require.Equal(t, true, success)
}

// Check if the timer is indeed cancelled, and the goroutine returns
func TestTimerCancel(t *testing.T) {
	success := false
	timer := util.NewTimer(func() {
		success = true
	})
	timer.SetTime(time.Second * 4)
	timer.Start()
	<-time.After(time.Second * 2)
	timer.Cancel()
	<-time.After(time.Second * 3)
	// By now, success should still be false
	require.Equal(t, false, success)
}
