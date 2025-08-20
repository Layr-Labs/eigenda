package node

import (
	"context"
	"fmt"

	"github.com/Layr-Labs/eigenda/core"
	lru "github.com/hashicorp/golang-lru/v2"
)

// A light wrapper around ChainState that caches the operator state. If the operator state is requested multiple
// times for the same reference block number, this utility will only fetch the operator state once and cache it.
type OperatorStateCache struct {
	chainState core.ChainState
	cache      *lru.Cache[uint, *core.OperatorState]
}

// Create a new caching wrapper around ChainState for fetching operator state.
func NewOperatorStateCache(
	chainState core.ChainState,
	cacheSize int,
) (*OperatorStateCache, error) {

	cache, err := lru.New[uint, *core.OperatorState](cacheSize)
	if err != nil {
		return nil, fmt.Errorf("NewOperatorStateCache: %w", err)
	}

	return &OperatorStateCache{
		chainState: chainState,
		cache:      cache,
	}, nil
}

// GetOperatorState retrieves the operator state for a given reference block number and quorums.
func (c *OperatorStateCache) GetOperatorState(
	ctx context.Context,
	referenceBlockNumber uint,
	quorums []core.QuorumID,
) (*core.OperatorState, error) {

	// TODO locking, perhaps an index lock?

	// Check if the operator state is already cached
	if state, found := c.cache.Get(referenceBlockNumber); found {
		return state, nil
	}

	// Fetch the operator state for all quorums.
	allQuorums, err := c.getAllQuorums(ctx, referenceBlockNumber)
	if err != nil {
		return nil, fmt.Errorf("getAllQuorums: %w", err)
	}

	state, err := c.chainState.GetOperatorState(ctx, referenceBlockNumber, allQuorums)
	if err != nil {
		return nil, fmt.Errorf("GetOperatorState: %w", err)
	}

	// Cache the fetched operator state
	c.cache.Add(referenceBlockNumber, state)

	// Only return data on the specified quorums.
	filteredState, err := c.filterByQuorum(state, quorums)
	if err != nil {
		return nil, fmt.Errorf("filterByQuorum: %w", err)
	}

	return filteredState, nil
}

func (c *OperatorStateCache) getAllQuorums(ctx context.Context, referenceBlockNumber uint) ([]core.QuorumID, error) {
	// TODO
	return nil, nil
}

// The code expects an operator state with an exact set of quorums, so filter out any extras. Easier to do this
// than to rewrite existing code that expects a specific set of quorums.
func (c *OperatorStateCache) filterByQuorum(
	state *core.OperatorState,
	quorums []core.QuorumID,
) (*core.OperatorState, error) {

	filteredState := &core.OperatorState{
		Operators:   make(map[core.QuorumID]map[core.OperatorID]*core.OperatorInfo, len(quorums)),
		Totals:      make(map[core.QuorumID]*core.OperatorInfo, len(quorums)),
		BlockNumber: state.BlockNumber,
	}

	for _, quorumID := range quorums {
		operators, ok := state.Operators[quorumID]
		if !ok {
			return nil, fmt.Errorf("quorum %s not found in operator state", quorumID)
		}
		totals, ok := state.Totals[quorumID]
		if !ok {
			return nil, fmt.Errorf("totals for quorum %s not found in operator state", quorumID)
		}
		filteredState.Operators[quorumID] = operators
		filteredState.Totals[quorumID] = totals
	}

	return filteredState, nil
}
