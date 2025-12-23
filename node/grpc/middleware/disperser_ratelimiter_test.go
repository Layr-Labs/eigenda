package middleware

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestDisperserRateLimiter_AllowsUntilBurst(t *testing.T) {
	t.Parallel()

	limiter := NewDisperserRateLimiter(nil, 1, 3) // 1 rps, burst 3
	id := uint32(123)
	now := time.Unix(1000, 0)

	require.True(t, limiter.Allow(id, now))
	require.True(t, limiter.Allow(id, now))
	require.True(t, limiter.Allow(id, now))
	// Burst exhausted
	require.False(t, limiter.Allow(id, now))
}

func TestDisperserRateLimiter_RefillsOverTime(t *testing.T) {
	t.Parallel()

	limiter := NewDisperserRateLimiter(nil, 1, 2) // 1 rps, burst 2
	id := uint32(7)
	start := time.Unix(1000, 0)

	require.True(t, limiter.Allow(id, start))
	require.True(t, limiter.Allow(id, start))
	require.False(t, limiter.Allow(id, start))

	// After 1s, one token should refill.
	require.True(t, limiter.Allow(id, start.Add(1*time.Second)))
	// But not yet two.
	require.False(t, limiter.Allow(id, start.Add(1*time.Second)))

	// After enough time, burst should be full again.
	require.True(t, limiter.Allow(id, start.Add(3*time.Second)))
	require.True(t, limiter.Allow(id, start.Add(3*time.Second)))
	require.False(t, limiter.Allow(id, start.Add(3*time.Second)))
}

func TestDisperserRateLimiter_DisabledWhenZeroOrNil(t *testing.T) {
	t.Parallel()

	id := uint32(42)
	now := time.Unix(1000, 0)

	limiterZero := NewDisperserRateLimiter(nil, 0, 1)
	require.True(t, limiterZero.Allow(id, now))

	limiterBurstZero := NewDisperserRateLimiter(nil, 1, 0)
	require.True(t, limiterBurstZero.Allow(id, now))

	var limiterNil *DisperserRateLimiter
	require.True(t, limiterNil.Allow(id, now))
}
