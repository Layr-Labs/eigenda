package node

import (
	"time"
)

// Time is an interface for mockable time operations
type Time interface {
	Now() time.Time
	Unix(sec int64, nsec int64) time.Time
	Since(t time.Time) time.Duration
}

// RealTime implements Time interface using actual time functions
type RealTime struct{}

// Now returns the current time
func (rt *RealTime) Now() time.Time {
	return time.Now()
}

// Unix returns the local Time corresponding to the given Unix time
func (rt *RealTime) Unix(sec int64, nsec int64) time.Time {
	return time.Unix(sec, nsec)
}

// Since returns the time elapsed since t
func (rt *RealTime) Since(t time.Time) time.Duration {
	return time.Since(t)
}

// DefaultTime is the default time implementation
var DefaultTime Time = &RealTime{}
