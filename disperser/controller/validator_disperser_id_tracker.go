package controller

import (
	"sync"

	"github.com/Layr-Labs/eigenda/core"
)

// ValidatorDisperserIDTracker tracks which disperser ID each validator accepts.
// This allows the controller to avoid retrying with wrong IDs once a validator's
// accepted ID is known.
type ValidatorDisperserIDTracker struct {
	mu sync.RWMutex
	// validatorAcceptedID maps validator ID to the disperser ID they accept
	// If a validator is not in the map, we haven't determined their accepted ID yet
	validatorAcceptedID map[core.OperatorID]uint32
	// disperserIDs is the ordered list of disperser IDs to try (in priority order)
	disperserIDs []uint32
}

// ValidatorDisperserIDStats contains statistics about validator disperser ID distribution.
type ValidatorDisperserIDStats struct {
	// CountByID maps disperser ID to the number of validators using it
	CountByID map[uint32]int
	// UnknownCount is the number of validators with unknown/undetermined ID
	UnknownCount int
}

// NewValidatorDisperserIDTracker creates a new tracker.
// disperserIDs should be provided in priority order (first ID is tried first).
func NewValidatorDisperserIDTracker(disperserIDs []uint32) *ValidatorDisperserIDTracker {
	return &ValidatorDisperserIDTracker{
		validatorAcceptedID: make(map[core.OperatorID]uint32),
		disperserIDs:        disperserIDs,
	}
}

// GetDisperserID returns the disperser ID to use for the given validator.
// If the validator's accepted ID is known, it returns that ID.
// If unknown, it returns the first ID in the priority list.
func (t *ValidatorDisperserIDTracker) GetDisperserID(validatorID core.OperatorID) uint32 {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if acceptedID, ok := t.validatorAcceptedID[validatorID]; ok {
		return acceptedID
	}

	// Unknown validator - return first/highest priority ID
	if len(t.disperserIDs) > 0 {
		return t.disperserIDs[0]
	}

	return 0 // Shouldn't happen if tracker is properly initialized
}

// RecordSuccess records that a validator successfully accepted a request with the given disperser ID.
func (t *ValidatorDisperserIDTracker) RecordSuccess(validatorID core.OperatorID, disperserID uint32) {
	t.mu.Lock()
	defer t.mu.Unlock()

	// Verify this is a known disperser ID
	isKnown := false
	for _, id := range t.disperserIDs {
		if id == disperserID {
			isKnown = true
			break
		}
	}

	if isKnown {
		t.validatorAcceptedID[validatorID] = disperserID
	}
}

// RecordFailure records that a validator rejected a request with the given disperser ID.
// This doesn't change the tracked ID - the controller will try the next ID in the list.
func (t *ValidatorDisperserIDTracker) RecordFailure(validatorID core.OperatorID, disperserID uint32) {
	// Failure doesn't update the tracker - we'll try the next ID in the priority list
	// The tracker only records successful IDs
}

// GetStats returns statistics about validator disperser ID distribution.
func (t *ValidatorDisperserIDTracker) GetStats() ValidatorDisperserIDStats {
	t.mu.RLock()
	defer t.mu.RUnlock()

	stats := ValidatorDisperserIDStats{
		CountByID: make(map[uint32]int),
	}

	for _, acceptedID := range t.validatorAcceptedID {
		stats.CountByID[acceptedID]++
	}

	return stats
}

// GetNextDisperserID returns the next disperser ID to try after the given ID failed.
// Returns the next ID in priority order, or 0 if no more IDs to try.
func (t *ValidatorDisperserIDTracker) GetNextDisperserID(currentID uint32) uint32 {
	t.mu.RLock()
	defer t.mu.RUnlock()

	for i, id := range t.disperserIDs {
		if id == currentID && i+1 < len(t.disperserIDs) {
			return t.disperserIDs[i+1]
		}
	}

	return 0 // No more IDs to try
}
