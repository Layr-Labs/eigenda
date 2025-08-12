package reservation

import (
	"fmt"
	"time"
)

// TimeMovedBackwardError is returned when a timestamp is observed that is before a previously observed timestamp.
//
// This should not normally happen, but with clock drift and NTP adjustments, system clocks can occasionally jump
// backward. This error allows the system to handle such cases gracefully rather than fatally erroring.
type TimeMovedBackwardError struct {
	// The current time that was provided
	CurrentTime time.Time
	// The previously observed time that is after CurrentTime
	PreviousTime time.Time
}

// Implements the error interface
func (e *TimeMovedBackwardError) Error() string {
	return fmt.Sprintf("time moved backward: current time %s is before previous time %s (delta: %v)",
		e.CurrentTime.Format(time.RFC3339Nano),
		e.PreviousTime.Format(time.RFC3339Nano),
		e.PreviousTime.Sub(e.CurrentTime))
}
