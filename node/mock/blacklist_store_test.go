package mock

import (
	"context"
	"crypto/sha256"
	"fmt"
	"sync"
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
	retrieved, err := store.Get(testKey)
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

	nonExistentKey := []byte("does-not-exist")

	// Should return error
	_, err = store.Get(nonExistentKey)
	require.Error(t, err)
}

// TestBlacklistStoreIsBlacklistedWithMockTime tests various blacklisting scenarios using mock time
func TestBlacklistStoreIsBlacklistedWithMockTime(t *testing.T) {
	logger, err := common.NewLogger(common.DefaultTextLoggerConfig())
	require.NoError(t, err)

	ctx := context.Background()

	t.Run("SingleEntryExactlyOneHour", func(t *testing.T) {
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

	// Test Case 5: 14 days after last update, should not be blacklisted
	t.Run("FourteenDaysAfterLastUpdate1Entry", func(t *testing.T) {
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

		disperserId := uint32(205)

		// Add four entries (should behave same as 3+)
		err = store.AddEntry(ctx, disperserId, "context1", "violation1")
		require.NoError(t, err)

		// check if disperser is blacklisted
		require.True(t, store.IsBlacklisted(ctx, disperserId))

		// check if disperser is not blacklisted after 1 hour but entry exists
		mockTime.SinceFunc = func(t time.Time) time.Duration {
			return time.Hour
		}
		require.False(t, store.IsBlacklisted(ctx, disperserId))
		require.True(t, store.HasDisperserID(ctx, disperserId))

		// 14 days after last update, should not be blacklisted
		mockTime.SinceFunc = func(t time.Time) time.Duration {
			return 14 * 24 * time.Hour
		}
		require.False(t, store.IsBlacklisted(ctx, disperserId))
		require.False(t, store.HasDisperserID(ctx, disperserId))
	})

	t.Run("FourteenDaysAfterLastUpdate2Entries", func(t *testing.T) {
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

		disperserId := uint32(206)

		// Add two entries (should behave same as 3+)
		err = store.AddEntry(ctx, disperserId, "context1", "violation1")
		require.NoError(t, err)
		err = store.AddEntry(ctx, disperserId, "context2", "violation2")
		require.NoError(t, err)

		// check if disperser is blacklisted
		require.True(t, store.IsBlacklisted(ctx, disperserId))

		// check if disperser is not blacklisted after 1 hour but entry exists
		mockTime.SinceFunc = func(t time.Time) time.Duration {
			return 24 * time.Hour
		}
		require.False(t, store.IsBlacklisted(ctx, disperserId))
		require.True(t, store.HasDisperserID(ctx, disperserId))

		// 14 days after last update, should not be blacklisted
		mockTime.SinceFunc = func(t time.Time) time.Duration {
			return 14 * 24 * time.Hour
		}
		require.False(t, store.IsBlacklisted(ctx, disperserId))
		require.False(t, store.HasDisperserID(ctx, disperserId))
	})

	t.Run("FourteenDaysAfterLastUpdate3Entries", func(t *testing.T) {
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

		disperserId := uint32(207)

		// Add three entries (should behave same as 3+)
		err = store.AddEntry(ctx, disperserId, "context1", "violation1")
		require.NoError(t, err)
		err = store.AddEntry(ctx, disperserId, "context2", "violation2")
		require.NoError(t, err)
		err = store.AddEntry(ctx, disperserId, "context3", "violation3")
		require.NoError(t, err)

		// check if disperser is blacklisted
		require.True(t, store.IsBlacklisted(ctx, disperserId))
		// check if disperser is blacklisted for 1 week
		mockTime.SinceFunc = func(t time.Time) time.Duration {
			return 6 * 24 * time.Hour
		}
		require.True(t, store.IsBlacklisted(ctx, disperserId))

		// check if disperser is not blacklisted after 7 days
		mockTime.SinceFunc = func(t time.Time) time.Duration {
			return 7 * 24 * time.Hour
		}
		require.False(t, store.IsBlacklisted(ctx, disperserId))

		// check if disperser entry is deleted after 14 days
		mockTime.SinceFunc = func(t time.Time) time.Duration {
			return 14 * 24 * time.Hour
		}
		require.False(t, store.IsBlacklisted(ctx, disperserId))
		require.False(t, store.HasDisperserID(ctx, disperserId))

	})
}

// TestBlacklistStoreReadWriteExclusion tests basic locking behavior
func TestBlacklistStoreReadWriteExclusion(t *testing.T) {
	logger, err := common.NewLogger(common.DefaultTextLoggerConfig())
	require.NoError(t, err)

	testDir := t.TempDir()
	store, err := node.NewLevelDBBlacklistStore(testDir, logger, false, false, node.DefaultTime)
	require.NoError(t, err)

	ctx := context.Background()

	// Test that concurrent operations don't cause data corruption
	t.Run("ConcurrentReadWriteIntegrity", func(t *testing.T) {
		const disperserId = uint32(400)
		const numOperations = 50

		var wg sync.WaitGroup

		// Start multiple writers
		for i := 0; i < 5; i++ {
			wg.Add(1)
			go func(writerID int) {
				defer wg.Done()
				for j := 0; j < numOperations; j++ {
					_ = store.AddEntry(ctx, disperserId, fmt.Sprintf("writer%d-op%d", writerID, j), "violation")
				}
			}(i)
		}

		// Start multiple readers
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func(readerID int) {
				defer wg.Done()
				for j := 0; j < numOperations; j++ {
					_ = store.HasDisperserID(ctx, disperserId)
					_, _ = store.GetByDisperserID(ctx, disperserId)
					_ = store.IsBlacklisted(ctx, disperserId)
				}
			}(i)
		}

		// Wait for all operations to complete
		wg.Wait()

		// Verify final state is consistent
		blacklist, err := store.GetByDisperserID(ctx, disperserId)
		if err == nil {
			// If blacklist exists, it should have some entries
			require.Greater(t, len(blacklist.Entries), 0, "Blacklist should have entries")
			require.LessOrEqual(t, len(blacklist.Entries), 5*numOperations, "Should not have more entries than written")
		}
		// If err != nil, that's also valid (blacklist might have been cleaned up)
	})
}

// TestBlacklistStoreWriteExclusion tests that write operations are properly serialized
func TestBlacklistStoreWriteExclusion(t *testing.T) {
	logger, err := common.NewLogger(common.DefaultTextLoggerConfig())
	require.NoError(t, err)

	testDir := t.TempDir()
	store, err := node.NewLevelDBBlacklistStore(testDir, logger, false, false, node.DefaultTime)
	require.NoError(t, err)

	ctx := context.Background()
	const disperserId = uint32(500)
	const numWriters = 5

	startCh := make(chan struct{})
	completionOrder := make(chan int, numWriters)

	// Start multiple writers
	for i := 0; i < numWriters; i++ {
		go func(writerID int) {
			<-startCh // Wait for start signal

			// Each writer adds multiple entries
			for j := 0; j < 3; j++ {
				err := store.AddEntry(ctx, disperserId, fmt.Sprintf("writer%d-entry%d", writerID, j), "violation")
				require.NoError(t, err)
			}

			completionOrder <- writerID
		}(i)
	}

	// Start all writers at the "same" time
	close(startCh)

	// Collect completion order
	var completedWriters []int
	for i := 0; i < numWriters; i++ {
		writerID := <-completionOrder
		completedWriters = append(completedWriters, writerID)
	}

	// Verify all writers completed
	require.Len(t, completedWriters, numWriters)

	// Verify final state - should have numWriters * 3 entries
	blacklist, err := store.GetByDisperserID(ctx, disperserId)
	require.NoError(t, err)
	require.Len(t, blacklist.Entries, numWriters*3)
}

// TestBlacklistStoreIsBlacklistedConcurrency tests concurrent IsBlacklisted calls
func TestBlacklistStoreIsBlacklistedConcurrency(t *testing.T) {
	logger, err := common.NewLogger(common.DefaultTextLoggerConfig())
	require.NoError(t, err)

	testDir := t.TempDir()

	// Create mock time to control blacklist behavior
	mockTime := &MockTime{}
	baseTime := time.Date(2023, 6, 15, 10, 0, 0, 0, time.UTC)
	mockTime.NowFunc = func() time.Time { return baseTime }
	mockTime.UnixFunc = func(sec int64, nsec int64) time.Time { return time.Unix(sec, nsec) }
	mockTime.SinceFunc = func(t time.Time) time.Duration { return 30 * time.Minute } // Not blacklisted

	store, err := node.NewLevelDBBlacklistStore(testDir, logger, false, false, mockTime)
	require.NoError(t, err)

	ctx := context.Background()
	disperserId := uint32(600)

	// Add an entry
	err = store.AddEntry(ctx, disperserId, "context1", "violation1")
	require.NoError(t, err)

	const numCheckers = 20
	const checksPerChecker = 50

	startCh := make(chan struct{})
	resultCh := make(chan bool, numCheckers*checksPerChecker)

	// Start multiple IsBlacklisted checkers
	for i := 0; i < numCheckers; i++ {
		go func() {
			<-startCh
			for j := 0; j < checksPerChecker; j++ {
				isBlacklisted := store.IsBlacklisted(ctx, disperserId)
				resultCh <- isBlacklisted
			}
		}()
	}

	// Start all checkers
	start := time.Now()
	close(startCh)

	// Collect results
	trueCount := 0
	for i := 0; i < numCheckers*checksPerChecker; i++ {
		if <-resultCh {
			trueCount++
		}
	}

	duration := time.Since(start)
	t.Logf("Concurrent IsBlacklisted calls completed in %v", duration)

	// All should return true (blacklisted) since we set 30 minute duration
	require.Equal(t, numCheckers*checksPerChecker, trueCount)

	// Should complete quickly due to concurrent reads
	maxExpectedDuration := time.Millisecond * 200
	require.Less(t, duration, maxExpectedDuration, "Concurrent IsBlacklisted calls took too long")
}

// TestBlacklistStoreRaceConditions tests for race conditions using mixed operations
func TestBlacklistStoreRaceConditions(t *testing.T) {
	logger, err := common.NewLogger(common.DefaultTextLoggerConfig())
	require.NoError(t, err)

	testDir := t.TempDir()
	store, err := node.NewLevelDBBlacklistStore(testDir, logger, false, false, node.DefaultTime)
	require.NoError(t, err)

	ctx := context.Background()
	const numGoroutines = 10
	const operationsPerGoroutine = 50

	startCh := make(chan struct{})
	doneCh := make(chan struct{}, numGoroutines)

	// Mix of different operations
	for i := 0; i < numGoroutines; i++ {
		go func(goroutineID int) {
			defer func() { doneCh <- struct{}{} }()
			<-startCh

			for j := 0; j < operationsPerGoroutine; j++ {
				disperserId := uint32(700 + (goroutineID*10 + j%10)) // Spread across different IDs

				switch j % 6 {
				case 0: // Add entry
					_ = store.AddEntry(ctx, disperserId, fmt.Sprintf("g%d-j%d", goroutineID, j), "violation")
				case 1: // Check if exists
					_ = store.HasDisperserID(ctx, disperserId)
				case 2: // Get by disperser ID
					_, _ = store.GetByDisperserID(ctx, disperserId)
				case 3: // Check if blacklisted
					_ = store.IsBlacklisted(ctx, disperserId)
				case 4: // Delete
					_ = store.DeleteByDisperserID(ctx, disperserId)
				case 5: // Add another entry
					_ = store.AddEntry(ctx, disperserId, fmt.Sprintf("g%d-j%d-second", goroutineID, j), "another violation")
				}
			}
		}(i)
	}

	// Start all goroutines
	close(startCh)

	// Wait for all to complete
	for i := 0; i < numGoroutines; i++ {
		select {
		case <-doneCh:
			// Continue
		case <-time.After(time.Second * 10):
			t.Fatal("Race condition test timed out - possible deadlock")
		}
	}

	// If we get here without panics or deadlocks, the locking is working
	t.Log("Race condition test completed successfully")
}
