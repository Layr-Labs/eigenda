package mock

import (
	"context"
	"testing"

	"github.com/Layr-Labs/eigenda/node"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMockBlacklistStore(t *testing.T) {
	ctx := context.Background()

	// Test that we can use the interface with the mock implementation
	mockStore := NewMockBlacklistStore()

	// Verify the mock implements the interface
	var _ node.BlacklistStore = mockStore

	// Test IsBlacklisted with mock
	mockStore.On("IsBlacklisted", ctx, uint32(123)).Return(true)
	result := mockStore.IsBlacklisted(ctx, 123)
	assert.True(t, result)
	mockStore.AssertExpectations(t)

	// Test AddEntry with mock
	mockStore.On("AddEntry", ctx, uint32(456), "test_context", "test_reason").Return(nil)
	err := mockStore.AddEntry(ctx, 456, "test_context", "test_reason")
	require.NoError(t, err)
	mockStore.AssertExpectations(t)

	// Test HasDisperserID with mock
	mockStore.On("HasDisperserID", ctx, uint32(789)).Return(false)
	hasDisperser := mockStore.HasDisperserID(ctx, 789)
	assert.False(t, hasDisperser)
	mockStore.AssertExpectations(t)

	// Test GetByDisperserID with mock
	blacklist := &node.Blacklist{
		Entries: []node.BlacklistEntry{
			{
				DisperserID: 100,
				Metadata: node.BlacklistMetadata{
					ContextId: "test",
					Reason:    "testing",
				},
				Timestamp: 1234567890,
			},
		},
		LastUpdated: 1234567890,
	}
	mockStore.On("GetByDisperserID", ctx, uint32(100)).Return(blacklist, nil)

	result_blacklist, err := mockStore.GetByDisperserID(ctx, 100)
	require.NoError(t, err)
	assert.Equal(t, blacklist, result_blacklist)
	mockStore.AssertExpectations(t)
}

// TestBlacklistStoreMockIntegration demonstrates how to use the mock
// in a larger test scenario where you might inject it as a dependency
func TestBlacklistStoreMockIntegration(t *testing.T) {
	ctx := context.Background()
	mockStore := NewMockBlacklistStore()

	// Example function that takes a BlacklistStore interface
	checkIfBlacklisted := func(store node.BlacklistStore, disperserID uint32) bool {
		return store.IsBlacklisted(ctx, disperserID)
	}

	// Setup mock expectations
	mockStore.On("IsBlacklisted", ctx, uint32(999)).Return(false)

	// Test the function with our mock
	result := checkIfBlacklisted(mockStore, 999)
	assert.False(t, result)

	mockStore.AssertExpectations(t)
}
