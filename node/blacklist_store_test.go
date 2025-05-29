package node

import (
	"context"
	"crypto/sha256"
	"fmt"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/stretchr/testify/require"
)

// TestNewLevelDBBlacklistStore verifies that a new BlacklistStore can be created
// with a levelDB backend at the specified path.
func TestNewLevelDBBlacklistStore(t *testing.T) {
	logger, err := common.NewLogger(common.DefaultTextLoggerConfig())
	require.NoError(t, err)

	testDir := t.TempDir()

	store, err := NewLevelDBBlacklistStore(testDir, logger, false, false)
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
	store, err := NewLevelDBBlacklistStore(testDir, logger, false, false)
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

// TestBlacklistStoreHasDisperserID tests the HasDisperserID method which checks
// if a disperser ID exists in the blacklist store by hashing the ID.
func TestBlacklistStoreHasDisperserID(t *testing.T) {
	logger, err := common.NewLogger(common.DefaultTextLoggerConfig())
	require.NoError(t, err)

	testDir := t.TempDir()
	store, err := NewLevelDBBlacklistStore(testDir, logger, false, false)
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

// TestBlacklistStoreHasKey tests the HasKey method which checks if a raw key
// exists in the store.
func TestBlacklistStoreHasKey(t *testing.T) {
	logger, err := common.NewLogger(common.DefaultTextLoggerConfig())
	require.NoError(t, err)

	testDir := t.TempDir()
	store, err := NewLevelDBBlacklistStore(testDir, logger, false, false)
	require.NoError(t, err)

	ctx := context.Background()
	testKey := []byte("test-key")
	testValue := []byte("test-value")

	// Should not exist initially
	require.False(t, store.HasKey(ctx, testKey))

	// Put data directly
	err = store.Put(ctx, testKey, testValue)
	require.NoError(t, err)

	// Should exist now
	require.True(t, store.HasKey(ctx, testKey))
}

// TestBlacklistStoreGetByDisperserID tests retrieval of blacklist data by disperser ID.
// It verifies that the hashing and retrieval process works correctly.
func TestBlacklistStoreGetByDisperserID(t *testing.T) {
	logger, err := common.NewLogger(common.DefaultTextLoggerConfig())
	require.NoError(t, err)

	testDir := t.TempDir()
	store, err := NewLevelDBBlacklistStore(testDir, logger, false, false)
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
	store, err := NewLevelDBBlacklistStore(testDir, logger, false, false)
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

// TestBlacklistStoreIsBlacklistedExpired tests that blacklist entries expire
// according to their violation count timeouts.
func TestBlacklistStoreIsBlacklistedExpired(t *testing.T) {
	logger, err := common.NewLogger(common.DefaultTextLoggerConfig())
	require.NoError(t, err)

	testDir := t.TempDir()
	store, err := NewLevelDBBlacklistStore(testDir, logger, false, false)
	require.NoError(t, err)

	ctx := context.Background()
	disperserId := uint32(222)

	// Create a blacklist with one entry that's 2 hours old (should be expired)
	blacklist := &Blacklist{
		Entries: []BlacklistEntry{
			{
				DisperserID: disperserId,
				Metadata: BlacklistMetadata{
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
	store, err := NewLevelDBBlacklistStore(testDir, logger, false, false)
	require.NoError(t, err)

	ctx := context.Background()
	testKey := []byte("test-blacklist-key")

	// Create test blacklist
	blacklist := &Blacklist{
		Entries: []BlacklistEntry{
			{
				DisperserID: 123,
				Metadata: BlacklistMetadata{
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
	store, err := NewLevelDBBlacklistStore(testDir, logger, false, false)
	require.NoError(t, err)

	ctx := context.Background()
	nonExistentKey := []byte("does-not-exist")

	// Should return error
	_, err = store.Get(ctx, nonExistentKey)
	require.Error(t, err)
}

// TestBlacklistStorePut tests the Put method which stores raw data in the store.
func TestBlacklistStorePut(t *testing.T) {
	logger, err := common.NewLogger(common.DefaultTextLoggerConfig())
	require.NoError(t, err)

	testDir := t.TempDir()
	store, err := NewLevelDBBlacklistStore(testDir, logger, false, false)
	require.NoError(t, err)

	ctx := context.Background()
	testKey := []byte("put-test-key")
	testValue := []byte("put-test-value")

	// Put data
	err = store.Put(ctx, testKey, testValue)
	require.NoError(t, err)

	// Verify it exists
	require.True(t, store.HasKey(ctx, testKey))
}

// Note: Interface testing with mock is done in a separate test package to avoid import cycles
