package workers

import "time"

// InterceptableTicker is a wrapper around the time.Ticker struct.
// It allows for deterministic time passage to be simulated in tests.
type InterceptableTicker interface {
	// getTimeChannel returns the channel that the ticker sends ticks on. Equivalent to time.Ticker.C.
	GetTimeChannel() <-chan time.Time
}

// standardTicker behaves exactly like a time.Ticker, for use in production code.
type standardTicker struct {
	ticker *time.Ticker
}

// NewTicker creates a new InterceptableTicker that behaves like a time.Ticker.
func NewTicker(d time.Duration) InterceptableTicker {
	return &standardTicker{
		ticker: time.NewTicker(d),
	}
}

func (s *standardTicker) GetTimeChannel() <-chan time.Time {
	return s.ticker.C
}
