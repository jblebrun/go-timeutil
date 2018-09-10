package timeutil

import "time"

// Safely reset a timer, as per go docs
func ResetTimerSafely(timer *time.Timer, newInterval time.Duration) {
	// Proper Reset go nonsense (see docs:
	StopTimerSafely(timer)
	timer.Reset(newInterval)
}

func StopTimerSafely(timer *time.Timer) {
	// Proper Stop go nonsense (see docs:
	if !timer.Stop() {
		//Drain a value if one had been there.
		select {
		case <-timer.C:
		default:
		}
	}
}
