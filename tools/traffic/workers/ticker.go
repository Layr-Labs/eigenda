package workers

import "time"

// InterceptableTicker is a wrapper around the time.Ticker struct.
// It allows for deterministic time passage to be simulated in tests.
type InterceptableTicker interface {
	// getTimeChannel returns the channel that the ticker sends ticks on. Equivalent to time.Ticker.C.
	getTimeChannel() <-chan time.Time
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

func (s *standardTicker) getTimeChannel() <-chan time.Time {
	return s.ticker.C
}

// MockTicker is a deterministic implementation of the InterceptableTicker interface.
type MockTicker struct {
	channel chan time.Time
	now     time.Time
}

// NewMockTicker creates a new InterceptableTicker that can be deterministically controlled in tests.
func NewMockTicker(now time.Time) InterceptableTicker {
	return &MockTicker{
		channel: make(chan time.Time),
		now:     now,
	}
}

func (m *MockTicker) getTimeChannel() <-chan time.Time {
	return m.channel
}

// Tick advances the ticker by the specified duration.
func (m *MockTicker) Tick(elapsedTime time.Duration) {
	m.now = m.now.Add(elapsedTime)
	m.channel <- m.now
}
