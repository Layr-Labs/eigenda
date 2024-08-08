package test

import (
	"github.com/Layr-Labs/eigenda/tools/traffic/workers"
	"time"
)

var _ workers.InterceptableTicker = (*mockTicker)(nil)

// mockTicker is a deterministic implementation of the InterceptableTicker interface.
type mockTicker struct {
	channel chan time.Time
	now     time.Time
}

// newMockTicker creates a new InterceptableTicker that can be deterministically controlled in tests.
func newMockTicker(now time.Time) *mockTicker {
	return &mockTicker{
		channel: make(chan time.Time),
		now:     now,
	}
}

func (m *mockTicker) GetTimeChannel() <-chan time.Time {
	return m.channel
}

// Tick advances the ticker by the specified duration.
func (m *mockTicker) Tick(elapsedTime time.Duration) {
	m.now = m.now.Add(elapsedTime)
	m.channel <- m.now
}
