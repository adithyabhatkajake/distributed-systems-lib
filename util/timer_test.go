package util_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/adithyabhatkajake/libchatter/util"
)

// Check if the correct callback is called after the time
func TestTimerWaiting(t *testing.T) {
	success := false
	timer := util.NewTimer(func() {
		success = true
	})
	timer.SetTime(time.Second * 4)
	timer.Start()
	<-time.After(time.Second * 5)
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
	<-time.After(time.Second * 10)
	// By now, success should still be false
	require.Equal(t, false, success)
}

func TestTimerReset(t *testing.T) {
	count := 1
	timer := util.NewTimer(func() {
		count++
	})
	timer.SetTime(time.Second * 2)
	timer.Start()
	<-time.After(time.Second)
	timer.Reset()
	<-time.After(time.Second * 4)
	// By now count should be 2, not 3 or 1
	require.Equal(t, 2, count)
}
