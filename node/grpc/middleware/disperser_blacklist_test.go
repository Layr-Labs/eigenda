package middleware

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestDisperserBlacklist_TTL(t *testing.T) {
	t.Parallel()

	b := NewDisperserBlacklist(nil, 10*time.Minute)

	now := time.Unix(1000, 0)
	id := uint32(123)

	require.False(t, b.IsBlacklisted(id, now))

	b.Blacklist(id, now, "reason")
	require.True(t, b.IsBlacklisted(id, now))

	// Still blacklisted within the TTL.
	require.True(t, b.IsBlacklisted(id, now.Add(9*time.Minute)))

	// Expired after the TTL and should be pruned.
	require.False(t, b.IsBlacklisted(id, now.Add(11*time.Minute)))
	require.False(t, b.IsBlacklisted(id, now.Add(12*time.Minute)))
}

func TestDisperserBlacklist_DisabledWhenTTLZeroOrNegative(t *testing.T) {
	t.Parallel()

	now := time.Unix(1000, 0)
	id := uint32(42)

	b0 := NewDisperserBlacklist(nil, 0)
	b0.Blacklist(id, now, "reason")
	require.False(t, b0.IsBlacklisted(id, now))

	bNeg := NewDisperserBlacklist(nil, -1*time.Second)
	bNeg.Blacklist(id, now, "reason")
	require.False(t, bNeg.IsBlacklisted(id, now))

	// Nil blacklist should behave as disabled.
	var bNil *DisperserBlacklist
	require.False(t, bNil.IsBlacklisted(id, now))
	bNil.Blacklist(id, now, "reason") // should not panic
}
