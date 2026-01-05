package middleware

import (
	"sync"
	"time"

	"github.com/Layr-Labs/eigensdk-go/logging"
	"golang.org/x/time/rate"
)

// DisperserRateLimiter applies a token-bucket rate limit per disperser ID.
// The limiter is local (per process) and best-effort.
type DisperserRateLimiter struct {
	logger logging.Logger

	limit rate.Limit
	burst int

	mu    sync.Mutex
	state map[uint32]*rate.Limiter
}

// NewDisperserRateLimiter creates a per-disperser rate limiter. If limitPerSecond <= 0 or
// burst <= 0, rate limiting is disabled.
func NewDisperserRateLimiter(logger logging.Logger, limitPerSecond float64, burst int) *DisperserRateLimiter {
	return &DisperserRateLimiter{
		logger: logger,
		limit:  rate.Limit(limitPerSecond),
		burst:  burst,
		state:  make(map[uint32]*rate.Limiter),
	}
}

// Allow returns true if a request for the disperser is permitted at time now.
// Each call consumes one token; tokens are replenished over time up to burst.
func (l *DisperserRateLimiter) Allow(disperserID uint32, now time.Time) bool {
	if l == nil || l.limit <= 0 || l.burst <= 0 {
		return true
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	limiter := l.getOrCreateLimiterLocked(disperserID)
	return limiter.AllowN(now, 1)
}

func (l *DisperserRateLimiter) getOrCreateLimiterLocked(disperserID uint32) *rate.Limiter {
	limiter, ok := l.state[disperserID]
	if !ok || limiter == nil {
		limiter = rate.NewLimiter(l.limit, l.burst)
		l.state[disperserID] = limiter
	}
	return limiter
}
