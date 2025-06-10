package mock

import (
	"context"
	"crypto/sha256"
	"fmt"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/node"
	"github.com/stretchr/testify/require"
)

// TestNewLevelDBBlacklistStore verifies that a new BlacklistStore can be created
// with a levelDB backend at the specified path.
func TestNewLevelDBBlacklistStore(t *testing.T) {
	logger, err := common.NewLogger(common.DefaultTextLoggerConfig())
	require.NoError(t, err)

	testDir := t.TempDir()

	store, err := node.NewLevelDBBlacklistStore(testDir, logger, false, false, node.DefaultTime)
	require.NoError(t, err)
	require.NotNil(t, store)
}

// TestBlacklistStoreAddEntry tests the AddEntry method which adds new blacklist entries
// for dispersers. It verifies:
// - New entries can be added for dispersers not previously blacklisted
// - Existing blacklists can be updated with additional entries
// - Data is properly stored and retrievable
func TestBlacklistStoreAddEntry(t *testing.T) {
	logger, err := common.NewLogger(common.DefaultTextLoggerConfig())
	require.NoError(t, err)

	testDir := t.TempDir()
	store, err := node.NewLevelDBBlacklistStore(testDir, logger, false, false, node.DefaultTime)
	require.NoError(t, err)

	ctx := context.Background()
	disperserId := uint32(123)

	// Initially, disperser should not exist
	require.False(t, store.HasDisperserID(ctx, disperserId))

	// Add first entry
	err = store.AddEntry(ctx, disperserId, "context1", "violation1")
	require.NoError(t, err)

	// Disperser should now exist
	require.True(t, store.HasDisperserID(ctx, disperserId))

	// Get and verify the blacklist
	blacklist, err := store.GetByDisperserID(ctx, disperserId)
	require.NoError(t, err)
	require.Equal(t, 1, len(blacklist.Entries))
	require.Equal(t, disperserId, blacklist.Entries[0].DisperserID)
	require.Equal(t, "context1", blacklist.Entries[0].Metadata.ContextId)
	require.Equal(t, "violation1", blacklist.Entries[0].Metadata.Reason)

	// Add second entry to same disperser
	err = store.AddEntry(ctx, disperserId, "context2", "violation2")
	require.NoError(t, err)

	// Get and verify updated blacklist
	blacklist, err = store.GetByDisperserID(ctx, disperserId)
	require.NoError(t, err)
	require.Equal(t, 2, len(blacklist.Entries))
	require.Equal(t, "context2", blacklist.Entries[1].Metadata.ContextId)
	require.Equal(t, "violation2", blacklist.Entries[1].Metadata.Reason)
}

func TestBlacklistStoreDeleteByDisperserID(t *testing.T) {
	logger, err := common.NewLogger(common.DefaultTextLoggerConfig())
	require.NoError(t, err)

	testDir := t.TempDir()
	store, err := node.NewLevelDBBlacklistStore(testDir, logger, false, false, node.DefaultTime)
	require.NoError(t, err)

	ctx := context.Background()
	disperserId := uint32(123)

	// Initially, disperser should not exist
	require.False(t, store.HasDisperserID(ctx, disperserId))

	// Add first entry
	err = store.AddEntry(ctx, disperserId, "context1", "violation1")
	require.NoError(t, err)

	// Disperser should now exist
	require.True(t, store.HasDisperserID(ctx, disperserId))

	// Get and verify the blacklist
	blacklist, err := store.GetByDisperserID(ctx, disperserId)
	require.NoError(t, err)
	require.Equal(t, 1, len(blacklist.Entries))
	require.Equal(t, disperserId, blacklist.Entries[0].DisperserID)
	require.Equal(t, "context1", blacklist.Entries[0].Metadata.ContextId)
	require.Equal(t, "violation1", blacklist.Entries[0].Metadata.Reason)

	// Delete the disperser
	err = store.DeleteByDisperserID(ctx, disperserId)
	require.NoError(t, err)

	// Disperser should no longer exist
	require.False(t, store.HasDisperserID(ctx, disperserId))

}

// TestBlacklistStoreHasDisperserID tests the HasDisperserID method which checks
// if a disperser ID exists in the blacklist store by hashing the ID.
func TestBlacklistStoreHasDisperserID(t *testing.T) {
	logger, err := common.NewLogger(common.DefaultTextLoggerConfig())
	require.NoError(t, err)

	testDir := t.TempDir()
	store, err := node.NewLevelDBBlacklistStore(testDir, logger, false, false, node.DefaultTime)
	require.NoError(t, err)

	ctx := context.Background()
	disperserId := uint32(456)

	// Should not exist initially
	require.False(t, store.HasDisperserID(ctx, disperserId))

	// Add entry
	err = store.AddEntry(ctx, disperserId, "test-context", "test-reason")
	require.NoError(t, err)

	// Should exist now
	require.True(t, store.HasDisperserID(ctx, disperserId))

	// Different disperser should not exist
	require.False(t, store.HasDisperserID(ctx, uint32(999)))
}

// TestBlacklistStoreGetByDisperserID tests retrieval of blacklist data by disperser ID.
// It verifies that the hashing and retrieval process works correctly.
func TestBlacklistStoreGetByDisperserID(t *testing.T) {
	logger, err := common.NewLogger(common.DefaultTextLoggerConfig())
	require.NoError(t, err)

	testDir := t.TempDir()
	store, err := node.NewLevelDBBlacklistStore(testDir, logger, false, false, node.DefaultTime)
	require.NoError(t, err)

	ctx := context.Background()
	disperserId := uint32(789)

	// Should return error when not found
	_, err = store.GetByDisperserID(ctx, disperserId)
	require.Error(t, err)

	// Add entry
	err = store.AddEntry(ctx, disperserId, "test-context", "test-reason")
	require.NoError(t, err)

	// Should retrieve successfully
	blacklist, err := store.GetByDisperserID(ctx, disperserId)
	require.NoError(t, err)
	require.NotNil(t, blacklist)
	require.Equal(t, 1, len(blacklist.Entries))
}

// TestBlacklistStoreIsBlacklisted tests the blacklisting logic which uses different
// time periods based on the number of violations:
// - 1 entry: blacklisted for 1 hour
// - 2 entries: blacklisted for 1 day
// - 3+ entries: blacklisted for 1 week
func TestBlacklistStoreIsBlacklisted(t *testing.T) {
	logger, err := common.NewLogger(common.DefaultTextLoggerConfig())
	require.NoError(t, err)

	testDir := t.TempDir()
	store, err := node.NewLevelDBBlacklistStore(testDir, logger, false, false, node.DefaultTime)
	require.NoError(t, err)

	ctx := context.Background()
	disperserId := uint32(111)

	// Should not be blacklisted initially
	require.False(t, store.IsBlacklisted(ctx, disperserId))

	// Add first entry (should be blacklisted for 1 hour)
	err = store.AddEntry(ctx, disperserId, "context1", "violation1")
	require.NoError(t, err)
	require.True(t, store.IsBlacklisted(ctx, disperserId))

	// Add second entry (should be blacklisted for 1 day)
	err = store.AddEntry(ctx, disperserId, "context2", "violation2")
	require.NoError(t, err)
	require.True(t, store.IsBlacklisted(ctx, disperserId))

	// Add third entry (should be blacklisted for 1 week)
	err = store.AddEntry(ctx, disperserId, "context3", "violation3")
	require.NoError(t, err)
	require.True(t, store.IsBlacklisted(ctx, disperserId))
}

// TestBlacklistStoreIsBlacklistedExpiration tests the blacklisting logic expiration
func TestBlacklistStoreIsBlacklistedExpiration(t *testing.T) {
	logger, err := common.NewLogger(common.DefaultTextLoggerConfig())
	require.NoError(t, err)

	testDir := t.TempDir()

	// Create mock time
	mockTime := &MockTime{}
	baseTime := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)

	// Set up mock time functions
	mockTime.NowFunc = func() time.Time {
		return baseTime
	}
	mockTime.UnixFunc = func(sec int64, nsec int64) time.Time {
		return time.Unix(sec, nsec)
	}
	mockTime.SinceFunc = func(t time.Time) time.Duration {
		return baseTime.Sub(t)
	}

	store, err := node.NewLevelDBBlacklistStore(testDir, logger, false, false, mockTime)
	require.NoError(t, err)

	ctx := context.Background()
	disperserId := uint32(111)

	// Test Case 1: Single entry - should be blacklisted for 1 hour
	err = store.AddEntry(ctx, disperserId, "context1", "violation1")
	require.NoError(t, err)
	require.True(t, store.IsBlacklisted(ctx, disperserId))

	// Advance time by 30 minutes - should still be blacklisted
	mockTime.SinceFunc = func(t time.Time) time.Duration {
		return 30 * time.Minute
	}
	require.True(t, store.IsBlacklisted(ctx, disperserId))

	// Advance time by 2 hours - should no longer be blacklisted
	mockTime.SinceFunc = func(t time.Time) time.Duration {
		return 2 * time.Hour
	}
	require.False(t, store.IsBlacklisted(ctx, disperserId))

	// Test Case 2: Two entries - should be blacklisted for 1 day
	// Reset time and add second entry
	mockTime.SinceFunc = func(t time.Time) time.Duration {
		return baseTime.Sub(t)
	}
	err = store.AddEntry(ctx, disperserId, "context2", "violation2")
	require.NoError(t, err)
	require.True(t, store.IsBlacklisted(ctx, disperserId))

	// Advance time by 12 hours - should still be blacklisted
	mockTime.SinceFunc = func(t time.Time) time.Duration {
		return 12 * time.Hour
	}
	require.True(t, store.IsBlacklisted(ctx, disperserId))

	// Advance time by 25 hours - should no longer be blacklisted
	mockTime.SinceFunc = func(t time.Time) time.Duration {
		return 25 * time.Hour
	}
	require.False(t, store.IsBlacklisted(ctx, disperserId))

	// Test Case 3: Three entries - should be blacklisted for 1 week
	// Reset time and add third entry
	mockTime.SinceFunc = func(t time.Time) time.Duration {
		return baseTime.Sub(t)
	}
	err = store.AddEntry(ctx, disperserId, "context3", "violation3")
	require.NoError(t, err)
	require.True(t, store.IsBlacklisted(ctx, disperserId))

	// Advance time by 3 days - should still be blacklisted
	mockTime.SinceFunc = func(t time.Time) time.Duration {
		return 72 * time.Hour
	}
	require.True(t, store.IsBlacklisted(ctx, disperserId))

	// Advance time by 8 days - should no longer be blacklisted
	mockTime.SinceFunc = func(t time.Time) time.Duration {
		return 8 * 24 * time.Hour
	}
	require.False(t, store.IsBlacklisted(ctx, disperserId))
}

// TestBlacklistStoreIsBlacklistedExpired tests that blacklist entries expire
// according to their violation count timeouts.
func TestBlacklistStoreIsBlacklistedExpired(t *testing.T) {
	logger, err := common.NewLogger(common.DefaultTextLoggerConfig())
	require.NoError(t, err)

	testDir := t.TempDir()
	store, err := node.NewLevelDBBlacklistStore(testDir, logger, false, false, node.DefaultTime)
	require.NoError(t, err)

	ctx := context.Background()
	disperserId := uint32(222)

	// Create a blacklist with one entry that's 2 hours old (should be expired)
	blacklist := &node.Blacklist{
		Entries: []node.BlacklistEntry{
			{
				DisperserID: disperserId,
				Metadata: node.BlacklistMetadata{
					ContextId: "old-context",
					Reason:    "old-violation",
				},
				Timestamp: uint64(time.Now().Add(-2 * time.Hour).Unix()),
			},
		},
		LastUpdated: uint64(time.Now().Add(-2 * time.Hour).Unix()),
	}

	// Store directly
	data, err := blacklist.ToBytes()
	require.NoError(t, err)

	disperserIdHash := sha256.Sum256(fmt.Appendf(nil, "%d", disperserId))
	err = store.Put(ctx, disperserIdHash[:], data)
	require.NoError(t, err)

	// Should not be blacklisted (expired after 1 hour)
	require.False(t, store.IsBlacklisted(ctx, disperserId))
}

// TestBlacklistStoreGet tests the Get method which retrieves blacklist data by raw key
// and properly deserializes it.
func TestBlacklistStoreGet(t *testing.T) {
	logger, err := common.NewLogger(common.DefaultTextLoggerConfig())
	require.NoError(t, err)

	testDir := t.TempDir()
	store, err := node.NewLevelDBBlacklistStore(testDir, logger, false, false, node.DefaultTime)
	require.NoError(t, err)

	ctx := context.Background()
	testKey := []byte("test-blacklist-key")

	// Create test blacklist
	blacklist := &node.Blacklist{
		Entries: []node.BlacklistEntry{
			{
				DisperserID: 123,
				Metadata: node.BlacklistMetadata{
					ContextId: "test-context",
					Reason:    "test-reason",
				},
				Timestamp: 12345,
			},
		},
		LastUpdated: 12345,
	}

	// Serialize and store
	data, err := blacklist.ToBytes()
	require.NoError(t, err)
	err = store.Put(ctx, testKey, data)
	require.NoError(t, err)

	// Retrieve and verify
	retrieved, err := store.Get(ctx, testKey)
	require.NoError(t, err)
	require.Equal(t, blacklist.LastUpdated, retrieved.LastUpdated)
	require.Equal(t, len(blacklist.Entries), len(retrieved.Entries))
	require.Equal(t, blacklist.Entries[0].DisperserID, retrieved.Entries[0].DisperserID)
}

// TestBlacklistStoreGetNonExistent tests error handling when trying to get
// a blacklist that doesn't exist.
func TestBlacklistStoreGetNonExistent(t *testing.T) {
	logger, err := common.NewLogger(common.DefaultTextLoggerConfig())
	require.NoError(t, err)

	testDir := t.TempDir()
	store, err := node.NewLevelDBBlacklistStore(testDir, logger, false, false, node.DefaultTime)
	require.NoError(t, err)

	ctx := context.Background()
	nonExistentKey := []byte("does-not-exist")

	// Should return error
	_, err = store.Get(ctx, nonExistentKey)
	require.Error(t, err)
}

// TestBlacklistStoreIsBlacklistedWithMockTime tests various blacklisting scenarios using mock time
func TestBlacklistStoreIsBlacklistedWithMockTime(t *testing.T) {
	logger, err := common.NewLogger(common.DefaultTextLoggerConfig())
	require.NoError(t, err)

	testDir := t.TempDir()

	// Create mock time
	mockTime := &MockTime{}
	baseTime := time.Date(2023, 6, 15, 10, 0, 0, 0, time.UTC)

	mockTime.NowFunc = func() time.Time {
		return baseTime
	}
	mockTime.UnixFunc = func(sec int64, nsec int64) time.Time {
		return time.Unix(sec, nsec)
	}

	store, err := node.NewLevelDBBlacklistStore(testDir, logger, false, false, mockTime)
	require.NoError(t, err)

	ctx := context.Background()

	t.Run("SingleEntryExactlyOneHour", func(t *testing.T) {
		disperserId := uint32(201)

		// Add entry
		err = store.AddEntry(ctx, disperserId, "context1", "violation1")
		require.NoError(t, err)

		// Exactly 1 hour later - should NOT be blacklisted (>= 1 hour)
		mockTime.SinceFunc = func(t time.Time) time.Duration {
			return time.Hour
		}
		require.False(t, store.IsBlacklisted(ctx, disperserId))

		// Just under 1 hour - should be blacklisted
		mockTime.SinceFunc = func(t time.Time) time.Duration {
			return time.Hour - time.Second
		}
		require.True(t, store.IsBlacklisted(ctx, disperserId))
	})

	t.Run("TwoEntriesExactlyOneDay", func(t *testing.T) {
		disperserId := uint32(202)

		// Add two entries
		err = store.AddEntry(ctx, disperserId, "context1", "violation1")
		require.NoError(t, err)
		err = store.AddEntry(ctx, disperserId, "context2", "violation2")
		require.NoError(t, err)

		// Exactly 24 hours later - should NOT be blacklisted
		mockTime.SinceFunc = func(t time.Time) time.Duration {
			return 24 * time.Hour
		}
		require.False(t, store.IsBlacklisted(ctx, disperserId))

		// Just under 24 hours - should be blacklisted
		mockTime.SinceFunc = func(t time.Time) time.Duration {
			return 24*time.Hour - time.Second
		}
		require.True(t, store.IsBlacklisted(ctx, disperserId))
	})

	t.Run("ThreeEntriesExactlyOneWeek", func(t *testing.T) {
		disperserId := uint32(203)

		// Add three entries
		err = store.AddEntry(ctx, disperserId, "context1", "violation1")
		require.NoError(t, err)
		err = store.AddEntry(ctx, disperserId, "context2", "violation2")
		require.NoError(t, err)
		err = store.AddEntry(ctx, disperserId, "context3", "violation3")
		require.NoError(t, err)

		// Exactly 1 week later - should NOT be blacklisted
		mockTime.SinceFunc = func(t time.Time) time.Duration {
			return 7 * 24 * time.Hour
		}
		require.False(t, store.IsBlacklisted(ctx, disperserId))

		// Just under 1 week - should be blacklisted
		mockTime.SinceFunc = func(t time.Time) time.Duration {
			return 7*24*time.Hour - time.Second
		}
		require.True(t, store.IsBlacklisted(ctx, disperserId))
	})

	t.Run("FourEntriesSameAsThree", func(t *testing.T) {
		disperserId := uint32(204)

		// Add four entries (should behave same as 3+)
		err = store.AddEntry(ctx, disperserId, "context1", "violation1")
		require.NoError(t, err)
		err = store.AddEntry(ctx, disperserId, "context2", "violation2")
		require.NoError(t, err)
		err = store.AddEntry(ctx, disperserId, "context3", "violation3")
		require.NoError(t, err)
		err = store.AddEntry(ctx, disperserId, "context4", "violation4")
		require.NoError(t, err)

		// Should still be blacklisted for 1 week
		mockTime.SinceFunc = func(t time.Time) time.Duration {
			return 3 * 24 * time.Hour
		}
		require.True(t, store.IsBlacklisted(ctx, disperserId))

		// After 1 week, should not be blacklisted
		mockTime.SinceFunc = func(t time.Time) time.Duration {
			return 7 * 24 * time.Hour
		}
		require.False(t, store.IsBlacklisted(ctx, disperserId))
	})
}

// Note: Interface testing with mock is done in a separate test package to avoid import cycles
