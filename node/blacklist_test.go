package node

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// TestBlacklistCreate verifies that a new Blacklist struct can be created
// with empty entries and zero LastUpdated time.
func TestBlacklistCreate(t *testing.T) {
	blacklist := Blacklist{
		Entries:     []BlacklistEntry{},
		LastUpdated: 0,
	}

	require.Equal(t, 0, len(blacklist.Entries))
	require.Equal(t, uint64(0), blacklist.LastUpdated)
}

// TestBlacklistAddEntry tests the AddEntry method which adds new blacklist entries
// with specified disperser ID, context ID, and reason.
// It verifies:
// - Entry data is stored correctly
// - Timestamp is set appropriately
// - LastUpdated field is updated when new entries are added
// - Multiple entries can be added successfully
func TestBlacklistAddEntry(t *testing.T) {
	blacklist := Blacklist{
		Entries:     []BlacklistEntry{},
		LastUpdated: 0,
	}

	// Add an entry and verify it was added correctly
	blacklist.AddEntry(1, "test-context", "violation")

	require.Equal(t, 1, len(blacklist.Entries))
	require.Equal(t, uint32(1), blacklist.Entries[0].DisperserID)
	require.Equal(t, "test-context", blacklist.Entries[0].Metadata.ContextId)
	require.Equal(t, "violation", blacklist.Entries[0].Metadata.Reason)
	require.Equal(t, blacklist.LastUpdated, blacklist.Entries[0].Timestamp)

	// Verify timestamp is recent
	now := uint64(time.Now().Unix())
	require.LessOrEqual(t, blacklist.LastUpdated, now)
	require.GreaterOrEqual(t, blacklist.LastUpdated, now-5) // Within 5 seconds

	// Add a second entry
	time.Sleep(1 * time.Second) // Ensure timestamp changes
	prevTimestamp := blacklist.LastUpdated
	blacklist.AddEntry(2, "another-context", "another-violation")

	require.Equal(t, 2, len(blacklist.Entries))
	require.Greater(t, blacklist.LastUpdated, prevTimestamp)
	require.Equal(t, uint32(2), blacklist.Entries[1].DisperserID)
	require.Equal(t, "another-context", blacklist.Entries[1].Metadata.ContextId)
	require.Equal(t, "another-violation", blacklist.Entries[1].Metadata.Reason)
}

// TestBlacklistSerialization validates the serialization and deserialization
// functionality of Blacklist using ToBytes and FromBytes methods.
// It ensures:
// - A populated blacklist can be converted to bytes
// - Those bytes can be converted back to an identical blacklist
// - All fields (including nested metadata) are preserved correctly
func TestBlacklistSerialization(t *testing.T) {
	original := Blacklist{
		Entries: []BlacklistEntry{
			{
				DisperserID: 1,
				Metadata: BlacklistMetadata{
					ContextId: "test-context",
					Reason:    "violation",
				},
				Timestamp: 12345,
			},
			{
				DisperserID: 2,
				Metadata: BlacklistMetadata{
					ContextId: "another-context",
					Reason:    "another-violation",
				},
				Timestamp: 67890,
			},
		},
		LastUpdated: 67890,
	}

	// Serialize to bytes
	bytes, err := original.ToBytes()
	require.NoError(t, err)
	require.NotEmpty(t, bytes)

	// Deserialize from bytes
	recreated := Blacklist{}
	err = recreated.FromBytes(bytes)
	require.NoError(t, err)

	// Verify the recreated blacklist matches the original
	require.Equal(t, len(original.Entries), len(recreated.Entries))
	require.Equal(t, original.LastUpdated, recreated.LastUpdated)

	for i, entry := range original.Entries {
		require.Equal(t, entry.DisperserID, recreated.Entries[i].DisperserID)
		require.Equal(t, entry.Metadata.ContextId, recreated.Entries[i].Metadata.ContextId)
		require.Equal(t, entry.Metadata.Reason, recreated.Entries[i].Metadata.Reason)
		require.Equal(t, entry.Timestamp, recreated.Entries[i].Timestamp)
	}
}

// TestBlacklistEmptySerialization tests serialization and deserialization
// specifically for an empty blacklist with no entries.
// It ensures that:
// - An empty blacklist can be serialized without errors
// - Deserializing it results in a proper empty blacklist
func TestBlacklistEmptySerialization(t *testing.T) {
	empty := Blacklist{
		Entries:     []BlacklistEntry{},
		LastUpdated: 0,
	}

	// Serialize empty blacklist
	bytes, err := empty.ToBytes()
	require.NoError(t, err)
	require.NotEmpty(t, bytes)

	// Deserialize from bytes
	recreated := Blacklist{}
	err = recreated.FromBytes(bytes)
	require.NoError(t, err)

	// Verify the recreated blacklist is empty
	require.Equal(t, 0, len(recreated.Entries))
	require.Equal(t, uint64(0), recreated.LastUpdated)
}

// TestBlacklistInvalidDeserialization tests error handling when attempting
// to deserialize invalid JSON data into a Blacklist struct.
// It verifies that the FromBytes method properly returns an error
// when given malformed input.
func TestBlacklistInvalidDeserialization(t *testing.T) {
	tests := []struct {
		name        string
		invalidJson []byte
	}{
		{
			name:        "completely invalid JSON",
			invalidJson: []byte(`{invalid json`),
		},
		{
			name:        "wrong type for disperser_id",
			invalidJson: []byte(`{"entries": [{"disperser_id": "not-a-number", "metadata": {"context_id": "test", "reason": "test"}, "timestamp": 123}], "last_updated": 456}`),
		},
		{
			name:        "wrong type for timestamp",
			invalidJson: []byte(`{"entries": [{"disperser_id": 123, "metadata": {"context_id": "test", "reason": "test"}, "timestamp": "not-a-number"}], "last_updated": 456}`),
		},
		{
			name:        "wrong type for last_updated",
			invalidJson: []byte(`{"entries": [], "last_updated": "not-a-number"}`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			blacklist := Blacklist{}
			err := blacklist.FromBytes(tt.invalidJson)
			require.Error(t, err)
		})
	}
}
