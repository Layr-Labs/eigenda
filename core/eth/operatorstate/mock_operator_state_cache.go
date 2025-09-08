package operatorstate

import (
	"context"
	"fmt"
	"sync"

	"github.com/Layr-Labs/eigenda/core"
)

var _ OperatorStateCache = (*MockOperatorStateCache)(nil)

// A mock implementation of the OperatorStateCache interface for testing purposes. States returned must be manually
// set using SetOperatorState.
type MockOperatorStateCache struct {
	// A "cache" of operator states, indexed by reference block number.
	cache sync.Map
}

// Create a new mock operator state cache. This cache does not have any initial data, and must be populated using
// SetOperatorState before it can be used.
func NewMockOperatorStateCache() *MockOperatorStateCache {
	return &MockOperatorStateCache{
		cache: sync.Map{},
	}
}

func (m *MockOperatorStateCache) GetOperatorState(
	_ context.Context,
	referenceBlockNumber uint64,
	quorums []core.QuorumID,
) (*core.OperatorState, error) {

	unfilteredState, ok := m.cache.Load(referenceBlockNumber)
	if !ok {
		return nil, fmt.Errorf("referenceBlockNumber %d not found in mock cache", referenceBlockNumber)
	}

	filteredState, err := filterByQuorum(unfilteredState.(*core.OperatorState), quorums)
	if err != nil {
		return nil, fmt.Errorf("failed to filter operator state by quorum: %w", err)
	}

	return filteredState, nil
}

// Set the operator state for a specific reference block number.
func (m *MockOperatorStateCache) SetOperatorState(
	_ context.Context,
	referenceBlockNumber uint64,
	operatorState *core.OperatorState,
) {
	m.cache.Store(referenceBlockNumber, operatorState)
}
