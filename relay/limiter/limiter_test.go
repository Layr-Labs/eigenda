package limiter

import (
	"github.com/stretchr/testify/require"
	"golang.org/x/time/rate"
	"testing"
	"time"
)

// The rate.Limiter library has less documentation than ideal. Although I can figure out what it's doing by reading
// the code, I think it's risky writing things that depend on what may change in the future. In these tests, I verify
// some basic properties of the rate.Limiter library, so that if these properties ever change in the future, the tests
// will fail and we'll know to update the code.

func TestPositiveTokens(t *testing.T) {
	configuredRate := rate.Limit(10.0)
	// "burst" is equivalent to the bucket size, aka the number of tokens that can be stored
	configuredBurst := 10

	// time starts at current time, but advances manually afterward
	now := time.Now()

	rateLimiter := rate.NewLimiter(configuredRate, configuredBurst)

	// number of tokens should equal the burst limit
	require.Equal(t, configuredBurst, int(rateLimiter.TokensAt(now)))

	// moving forward in time should not change the number of tokens
	now = now.Add(time.Second)
	require.Equal(t, configuredBurst, int(rateLimiter.TokensAt(now)))

	// remove each token without advancing time
	for i := 0; i < configuredBurst; i++ {
		require.True(t, rateLimiter.AllowN(now, 1))
		require.Equal(t, configuredBurst-i-1, int(rateLimiter.TokensAt(now)))
	}
	require.Equal(t, 0, int(rateLimiter.TokensAt(now)))

	// removing an additional token should fail
	require.False(t, rateLimiter.AllowN(now, 1))
	require.Equal(t, 0, int(rateLimiter.TokensAt(now)))

	// tokens should return at a rate of once per 100ms
	for i := 0; i < configuredBurst; i++ {
		now = now.Add(100 * time.Millisecond)
		require.Equal(t, i+1, int(rateLimiter.TokensAt(now)))
	}
	require.Equal(t, configuredBurst, int(rateLimiter.TokensAt(now)))

	// remove 7 tokens all at once
	require.True(t, rateLimiter.AllowN(now, 7))
	require.Equal(t, 3, int(rateLimiter.TokensAt(now)))

	// move forward 500ms, returning 5 tokens
	now = now.Add(500 * time.Millisecond)
	require.Equal(t, 8, int(rateLimiter.TokensAt(now)))

	// try to take more than the burst limit
	require.False(t, rateLimiter.AllowN(now, 100))
}

func TestNegativeTokens(t *testing.T) {
	configuredRate := rate.Limit(10.0)
	// "burst" is equivalent to the bucket size, aka the number of tokens that can be stored
	configuredBurst := 10

	// time starts at current time, but advances manually afterward
	now := time.Now()

	rateLimiter := rate.NewLimiter(configuredRate, configuredBurst)

	// number of tokens should equal the burst limit
	require.Equal(t, configuredBurst, int(rateLimiter.TokensAt(now)))

	// remove all tokens then add them back
	require.True(t, rateLimiter.AllowN(now, configuredBurst))
	require.Equal(t, 0, int(rateLimiter.TokensAt(now)))
	for i := 0; i < configuredBurst; i++ {
		require.True(t, rateLimiter.AllowN(now, -1))
		require.Equal(t, i+1, int(rateLimiter.TokensAt(now)))
	}

	// nothing funky should happen when time advances
	now = now.Add(100 * time.Second)
	require.Equal(t, configuredBurst, int(rateLimiter.TokensAt(now)))
}
