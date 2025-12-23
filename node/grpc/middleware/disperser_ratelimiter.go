package middleware

import (
	"sync"
	"time"

	"github.com/Layr-Labs/eigensdk-go/logging"
)

// DisperserRateLimiter applies a token-bucket rate limit per disperser ID.
// The limiter is local (per process) and best-effort.
type DisperserRateLimiter struct {
	logger logging.Logger

	limitPerSecond float64
	burst          int

	mu    sync.Mutex
	state map[uint32]*disperserLimiterState
}

type disperserLimiterState struct {
	tokens     float64
	lastRefill time.Time
}

// NewDisperserRateLimiter creates a per-disperser rate limiter. If limitPerSecond <= 0 or
// burst <= 0, rate limiting is disabled.
func NewDisperserRateLimiter(logger logging.Logger, limitPerSecond float64, burst int) *DisperserRateLimiter {
	return &DisperserRateLimiter{
		logger:         logger,
		limitPerSecond: limitPerSecond,
		burst:          burst,
		state:          make(map[uint32]*disperserLimiterState),
	}
}

// Allow returns true if a request for the disperser is permitted at time now.
// Each call consumes one token; tokens are replenished over time up to burst.
func (l *DisperserRateLimiter) Allow(disperserID uint32, now time.Time) bool {
	if l == nil || l.limitPerSecond <= 0 || l.burst <= 0 {
		return true
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	st := l.getOrCreateStateLocked(disperserID, now)

	// Refill based on elapsed time.
	elapsed := now.Sub(st.lastRefill).Seconds()
	if elapsed > 0 {
		st.tokens = minFloat(float64(l.burst), st.tokens+elapsed*l.limitPerSecond)
		st.lastRefill = now
	}

	if st.tokens < 1 {
		return false
	}

	st.tokens -= 1
	return true
}

func (l *DisperserRateLimiter) getOrCreateStateLocked(disperserID uint32, now time.Time) *disperserLimiterState {
	st, ok := l.state[disperserID]
	if !ok || st == nil {
		st = &disperserLimiterState{
			tokens:     float64(l.burst),
			lastRefill: now,
		}
		l.state[disperserID] = st
	}
	return st
}

func minFloat(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
