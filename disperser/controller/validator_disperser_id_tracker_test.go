package controller

import (
	"sync"
	"testing"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/stretchr/testify/assert"
)

func TestNewValidatorDisperserIDTracker(t *testing.T) {
	tracker := NewValidatorDisperserIDTracker([]uint32{1, 0})
	assert.NotNil(t, tracker)
	assert.Equal(t, []uint32{1, 0}, tracker.disperserIDs)
}

func TestValidatorDisperserIDTracker_GetDisperserID(t *testing.T) {
	tracker := NewValidatorDisperserIDTracker([]uint32{1, 0})
	validatorID := core.OperatorID{1, 2, 3}

	t.Run("unknown validator returns primary", func(t *testing.T) {
		disperserID := tracker.GetDisperserID(validatorID)
		assert.Equal(t, uint32(1), disperserID)
	})

	t.Run("known primary validator returns primary", func(t *testing.T) {
		tracker.RecordSuccess(validatorID, 1)
		disperserID := tracker.GetDisperserID(validatorID)
		assert.Equal(t, uint32(1), disperserID)
	})

	t.Run("known fallback validator returns fallback", func(t *testing.T) {
		validatorID2 := core.OperatorID{4, 5, 6}
		tracker.RecordSuccess(validatorID2, 0)
		disperserID := tracker.GetDisperserID(validatorID2)
		assert.Equal(t, uint32(0), disperserID)
	})
}

func TestValidatorDisperserIDTracker_RecordSuccess(t *testing.T) {
	tracker := NewValidatorDisperserIDTracker([]uint32{1, 0})
	validatorID := core.OperatorID{1, 2, 3}

	t.Run("record primary success", func(t *testing.T) {
		tracker.RecordSuccess(validatorID, 1)
		assert.Equal(t, uint32(1), tracker.GetDisperserID(validatorID))

		stats := tracker.GetStats()
		assert.Equal(t, 1, stats.CountByID[1])
		assert.Equal(t, 0, stats.CountByID[0])
	})

	t.Run("record fallback success", func(t *testing.T) {
		validatorID2 := core.OperatorID{4, 5, 6}
		tracker.RecordSuccess(validatorID2, 0)
		assert.Equal(t, uint32(0), tracker.GetDisperserID(validatorID2))

		stats := tracker.GetStats()
		assert.Equal(t, 1, stats.CountByID[1])
		assert.Equal(t, 1, stats.CountByID[0])
	})

	t.Run("record unknown disperser ID", func(t *testing.T) {
		validatorID3 := core.OperatorID{7, 8, 9}
		// Record success with unknown ID (e.g., 99)
		tracker.RecordSuccess(validatorID3, 99)
		// Should still return primary (unknown validators default to primary)
		assert.Equal(t, uint32(1), tracker.GetDisperserID(validatorID3))
	})
}

func TestValidatorDisperserIDTracker_RecordFailure(t *testing.T) {
	tracker := NewValidatorDisperserIDTracker([]uint32{1, 0})
	validatorID := core.OperatorID{1, 2, 3}

	t.Run("record failure does not change tracker", func(t *testing.T) {
		// Failure doesn't update the tracker - just returns the first ID for unknown validators
		tracker.RecordFailure(validatorID, 1)
		assert.Equal(t, uint32(1), tracker.GetDisperserID(validatorID))

		stats := tracker.GetStats()
		assert.Equal(t, 0, stats.CountByID[1])
		assert.Equal(t, 0, stats.CountByID[0])
	})

	t.Run("successful ID persists after failure record", func(t *testing.T) {
		validatorID2 := core.OperatorID{4, 5, 6}
		tracker.RecordSuccess(validatorID2, 0)
		tracker.RecordFailure(validatorID2, 0)
		// Should still be 0
		assert.Equal(t, uint32(0), tracker.GetDisperserID(validatorID2))
	})
}

func TestValidatorDisperserIDTracker_GetStats(t *testing.T) {
	tracker := NewValidatorDisperserIDTracker([]uint32{1, 0})

	t.Run("initial stats", func(t *testing.T) {
		stats := tracker.GetStats()
		assert.Equal(t, 0, stats.CountByID[1])
		assert.Equal(t, 0, stats.CountByID[0])
		assert.Equal(t, 0, stats.UnknownCount)
	})

	t.Run("mixed validator distribution", func(t *testing.T) {
		// Add 2 validators using ID 1
		tracker.RecordSuccess(core.OperatorID{1}, 1)
		tracker.RecordSuccess(core.OperatorID{2}, 1)

		// Add 3 validators using ID 0
		tracker.RecordSuccess(core.OperatorID{3}, 0)
		tracker.RecordSuccess(core.OperatorID{4}, 0)
		tracker.RecordSuccess(core.OperatorID{5}, 0)

		stats := tracker.GetStats()
		assert.Equal(t, 2, stats.CountByID[1])
		assert.Equal(t, 3, stats.CountByID[0])
		assert.Equal(t, 0, stats.UnknownCount)
	})
}

func TestValidatorDisperserIDTracker_ThreadSafety(t *testing.T) {
	tracker := NewValidatorDisperserIDTracker([]uint32{1, 0})

	// Run concurrent operations
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			validatorID := core.OperatorID{byte(idx)}

			// Concurrent reads and writes
			tracker.GetDisperserID(validatorID)
			tracker.RecordSuccess(validatorID, 1)
			tracker.RecordFailure(validatorID, 1)
			tracker.GetStats()
		}(i)
	}

	wg.Wait()

	// Verify we can still get stats without panic
	stats := tracker.GetStats()
	assert.GreaterOrEqual(t, stats.CountByID[0]+stats.CountByID[1], 0)
}
