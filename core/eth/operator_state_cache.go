package eth

import (
	"context"
	"fmt"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
	lru "github.com/hashicorp/golang-lru/v2"
)

// the size of the index lock used by the OperatorStateCache
const indexLockSize = 64

// Responsible for fetching and caching operator state for a given reference block number and quorums.
//
// This utility is thread safe, and can be used in performance sensitive, multithreaded environments.
type OperatorStateCache interface {
	// GetOperatorState retrieves the operator state for a given reference block number and quorums
	GetOperatorState(
		ctx context.Context,
		referenceBlockNumber uint64,
		quorums []core.QuorumID,
	) (*core.OperatorState, error)
}

var _ OperatorStateCache = (*operatorStateCache)(nil)

// A standard implementation of the OperatorStateCache interface.
type operatorStateCache struct {
	// indexes chain data, required to get operator public keys
	chainState core.ChainState

	// used to get a list of quorums registered at a given reference block number
	quorumScanner QuorumScanner

	// A cache for operator state, indexed by reference block number.
	// This cache implementation is thread safe.
	cache *lru.Cache[uint64, *core.OperatorState]

	// Used to prevent simultaneous lookup for a particular reference block number. Not used to protect data
	// structures against concurrent access.
	indexLock *common.IndexLock
}

// Create a new caching wrapper around ChainState for fetching operator state.
func NewOperatorStateCache(
	contractBackend bind.ContractBackend,
	chainState core.ChainState,
	registryCoordinatorAddress gethcommon.Address,
	cacheSize int,
) (OperatorStateCache, error) {

	cache, err := lru.New[uint64, *core.OperatorState](cacheSize)
	if err != nil {
		return nil, fmt.Errorf("NewOperatorStateCache: %w", err)
	}

	qs, err := NewQuorumScanner(contractBackend, registryCoordinatorAddress)
	if err != nil {
		return nil, fmt.Errorf("NewQuorumScanner: %w", err)
	}

	return &operatorStateCache{
		chainState:    chainState,
		quorumScanner: qs,
		cache:         cache,
		indexLock:     common.NewIndexLock(indexLockSize),
	}, nil
}

// GetOperatorState retrieves the operator state for a given reference block number and quorums.
func (c *operatorStateCache) GetOperatorState(
	ctx context.Context,
	referenceBlockNumber uint64,
	quorums []core.QuorumID,
) (*core.OperatorState, error) {

	// Acquire a lock that prevents simultaneous lookups for the same reference block number.
	c.indexLock.Lock(referenceBlockNumber)
	defer c.indexLock.Unlock(referenceBlockNumber)

	// Check if the operator state is already cached
	if state, found := c.cache.Get(referenceBlockNumber); found {
		return state, nil
	}

	// Fetch the operator state for all quorums.
	allQuorums, err := c.quorumScanner.GetQuorums(ctx, referenceBlockNumber)
	if err != nil {
		return nil, fmt.Errorf("getAllQuorums: %w", err)
	}
	state, err := c.chainState.GetOperatorState(ctx, uint(referenceBlockNumber), allQuorums)
	if err != nil {
		return nil, fmt.Errorf("GetOperatorState: %w", err)
	}

	// Cache the fetched operator state.
	c.cache.Add(referenceBlockNumber, state)

	// Only return data on the specified quorums.
	filteredState, err := filterByQuorum(state, quorums)
	if err != nil {
		return nil, fmt.Errorf("filterByQuorum: %w", err)
	}

	return filteredState, nil
}

// The code expects an operator state with an exact set of quorums, so filter out any extras. Easier to do this
// than to rewrite existing code that expects a specific set of quorums.
func filterByQuorum(
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
			return nil, fmt.Errorf("quorum %d not found in operator state", quorumID)
		}
		totals, ok := state.Totals[quorumID]
		if !ok {
			return nil, fmt.Errorf("totals for quorum %d not found in operator state", quorumID)
		}
		filteredState.Operators[quorumID] = operators
		filteredState.Totals[quorumID] = totals
	}

	return filteredState, nil
}
