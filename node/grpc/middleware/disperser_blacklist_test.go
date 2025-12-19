package middleware

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestDisperserBlacklist_TTL(t *testing.T) {
	t.Parallel()

	b := NewDisperserBlacklist(nil, 10*time.Minute, 2*time.Minute, 3)

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

	b0 := NewDisperserBlacklist(nil, 0, 2*time.Minute, 3)
	b0.Blacklist(id, now, "reason")
	require.False(t, b0.IsBlacklisted(id, now))

	bNeg := NewDisperserBlacklist(nil, -1*time.Second, 2*time.Minute, 3)
	bNeg.Blacklist(id, now, "reason")
	require.False(t, bNeg.IsBlacklisted(id, now))

	// Nil blacklist should behave as disabled.
	var bNil *DisperserBlacklist
	require.False(t, bNil.IsBlacklisted(id, now))
	bNil.Blacklist(id, now, "reason") // should not panic
}

func TestDisperserBlacklist_StrikeThreshold(t *testing.T) {
	t.Parallel()

	b := NewDisperserBlacklist(nil, 10*time.Minute, 2*time.Minute, 3)
	now := time.Unix(1000, 0)
	id := uint32(7)

	b.RecordInvalid(id, now, "bad1")
	require.False(t, b.IsBlacklisted(id, now))

	b.RecordInvalid(id, now.Add(30*time.Second), "bad2")
	require.False(t, b.IsBlacklisted(id, now.Add(30*time.Second)))

	// Third invalid within the window triggers blacklisting.
	b.RecordInvalid(id, now.Add(60*time.Second), "bad3")
	require.True(t, b.IsBlacklisted(id, now.Add(60*time.Second)))

	// After TTL expires, it should be forgiven (strikes cleared).
	require.False(t, b.IsBlacklisted(id, now.Add(11*time.Minute)))

	// One invalid after forgiveness should not immediately re-ban.
	b.RecordInvalid(id, now.Add(11*time.Minute), "bad4")
	require.False(t, b.IsBlacklisted(id, now.Add(11*time.Minute)))
}

func TestDisperserBlacklist_StrikeWindow(t *testing.T) {
	t.Parallel()

	b := NewDisperserBlacklist(nil, 10*time.Minute, 2*time.Minute, 3)
	now := time.Unix(1000, 0)
	id := uint32(8)

	// Two invalids, but then we wait past the strike window before the third.
	b.RecordInvalid(id, now, "bad1")
	b.RecordInvalid(id, now.Add(30*time.Second), "bad2")
	b.RecordInvalid(id, now.Add(3*time.Minute), "bad3")

	// The first two are outside the 2m window at the time of the third, so no ban.
	require.False(t, b.IsBlacklisted(id, now.Add(3*time.Minute)))
}
