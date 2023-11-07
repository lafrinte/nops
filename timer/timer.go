package timer

import (
	"sync"
	"time"
)

func initTimer(t *time.Timer, timeout time.Duration) *time.Timer {
	if t == nil {
		return time.NewTimer(timeout)
	}

	if t.Reset(timeout) {
		panic("active timer trapped into initTimer()")
	}

	return t
}

func stopTimer(t *time.Timer) {
	if !t.Stop() {
		// Collect possibly added time from the channel
		// if timer has been stopped and nobody collected its value.
		select {
		case <-t.C:
		default:
		}
	}
}

// AcquireTimer returns a time.Timer from the pool and updates it to
// send the current time on its channel after at least timeout.
func AcquireTimer(timeout time.Duration) *time.Timer {
	v := timerPool.Get()
	if v == nil {
		return time.NewTimer(timeout)
	}
	t := v.(*time.Timer)
	initTimer(t, timeout)
	return t
}

// ReleaseTimer returns the time.Timer acquired via AcquireTimer to the pool
// and prevents the Timer from firing.
func ReleaseTimer(t *time.Timer) {
	stopTimer(t)
	timerPool.Put(t)
}

var timerPool sync.Pool
