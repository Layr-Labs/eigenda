package chainstate

import (
	"context"
	"fmt"

	"github.com/Layr-Labs/eigenda/chainstate/store"
	"github.com/Layr-Labs/eigenda/chainstate/types"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigensdk-go/logging"
)

// IndexedChainState adapts the chainstate indexer to the core.IndexedChainState
// interface, so consumers such as the disperser controller can obtain indexed
// operator state from the local in-memory store instead of the operator-state
// subgraph.
//
// The embedded core.ChainState supplies the on-chain half of the interface
// (stake-weighted operator state read directly from contracts); this type only
// implements the "indexed" half (aggregate public keys and per-operator BLS
// keys / sockets) by reading from the indexer's store. This mirrors how
// core/thegraph composes an on-chain ChainState with subgraph-sourced indexed
// data, so the two implementations are drop-in interchangeable.
type IndexedChainState struct {
	core.ChainState

	indexer *Indexer
	store   store.Store
	logger  logging.Logger
}

var _ core.IndexedChainState = (*IndexedChainState)(nil)

// NewIndexedChainState wraps an indexer and an on-chain ChainState as a
// core.IndexedChainState. The chainState is used verbatim for the ChainState
// portion of the interface; the indexer's store backs the indexed queries.
func NewIndexedChainState(
	indexer *Indexer,
	chainState core.ChainState,
	logger logging.Logger,
) *IndexedChainState {
	return &IndexedChainState{
		ChainState: chainState,
		indexer:    indexer,
		store:      indexer.GetStore(),
		logger:     logger.With("component", "ChainStateIndexedChainState"),
	}
}

// Start begins the indexer's background indexing. It returns once indexing has
// started; callers that need the store to be caught up to a particular block
// before querying must poll GetLastIndexedBlock themselves.
func (ics *IndexedChainState) Start(ctx context.Context) error {
	if err := ics.indexer.Start(ctx); err != nil {
		return fmt.Errorf("failed to start indexer: %w", err)
	}
	return nil
}

// GetIndexedOperatorState returns the IndexedOperatorState for the given block
// and quorums: the on-chain operator state, plus each quorum's aggregate public
// key and each operator's indexed info (BLS keys and socket) from the store.
//
// The assembly mirrors core/thegraph's implementation, including its handling of
// operators that are present on-chain but missing from the index (an error) and
// indexed operators that are not part of any requested quorum at the reference
// block (filtered out).
func (ics *IndexedChainState) GetIndexedOperatorState(
	ctx context.Context,
	blockNumber uint,
	quorums []core.QuorumID,
) (*core.IndexedOperatorState, error) {
	operatorState, err := ics.ChainState.GetOperatorState(ctx, blockNumber, quorums)
	if err != nil {
		return nil, fmt.Errorf("failed to get operator state: %w", err)
	}

	aggKeys := make(map[core.QuorumID]*core.G1Point)
	for _, quorum := range quorums {
		apk, err := ics.getQuorumAPK(ctx, quorum, uint64(blockNumber))
		if err != nil {
			return nil, fmt.Errorf("failed to get aggregate public key for quorum %d: %w", quorum, err)
		}
		// A requested quorum with no APK snapshot is an error, not a silent
		// omission: returning a partial AggKeys map would let a consumer verify
		// signatures against an incomplete key set. This matches core/thegraph,
		// which surfaces a missing APK as an error.
		if apk == nil {
			return nil, fmt.Errorf("no aggregate public key found for quorum %d at block %d", quorum, blockNumber)
		}
		aggKeys[quorum] = apk
	}

	indexedOperators, err := ics.GetIndexedOperators(ctx, blockNumber)
	if err != nil {
		return nil, err
	}

	// Every operator present in the on-chain state must be present in the index.
	operatorSeen := make(map[core.OperatorID]struct{})
	for _, quorumOperators := range operatorState.Operators {
		for operatorID := range quorumOperators {
			if indexedOperators[operatorID] == nil {
				return nil, fmt.Errorf("operator %s not found in indexed state", operatorID.Hex())
			}
			operatorSeen[operatorID] = struct{}{}
		}
	}

	// Drop indexed operators that are not part of any requested quorum at the
	// reference block (e.g. operators that registered after blockNumber).
	for operatorID := range indexedOperators {
		if _, ok := operatorSeen[operatorID]; !ok {
			delete(indexedOperators, operatorID)
		}
	}

	return &core.IndexedOperatorState{
		OperatorState:    operatorState,
		IndexedOperators: indexedOperators,
		AggKeys:          aggKeys,
	}, nil
}

// GetIndexedOperators returns the indexed info (BLS keys and socket) for all
// operators known to the store that are registered as of blockNumber, keyed by
// operator ID.
//
// Like core/thegraph, this may over-fetch (an operator that registered after
// blockNumber but is still registered is included); callers that need results
// scoped to a reference block filter against the on-chain operator state, as
// GetIndexedOperatorState does.
func (ics *IndexedChainState) GetIndexedOperators(
	ctx context.Context,
	blockNumber uint,
) (map[core.OperatorID]*core.IndexedOperatorInfo, error) {
	// A zero/empty filter returns all operators; we filter by registration below.
	operators, err := ics.store.ListOperators(ctx, types.OperatorFilter{}, 0, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to list operators: %w", err)
	}

	result := make(map[core.OperatorID]*core.IndexedOperatorInfo, len(operators))
	for _, op := range operators {
		if !registeredAt(op, uint64(blockNumber)) {
			continue
		}
		info, err := toIndexedOperatorInfo(op)
		if err != nil {
			return nil, fmt.Errorf("operator %s: %w", op.ID.Hex(), err)
		}
		result[op.ID] = info
	}
	return result, nil
}

// GetIndexedOperatorInfoByOperatorId returns the indexed info for a single
// operator. blockNumber is accepted for interface compatibility; the indexed
// fields (BLS keys, socket) are taken from the operator's current stored state,
// matching core/thegraph which reads these immutable-ish fields from the latest
// indexed block.
func (ics *IndexedChainState) GetIndexedOperatorInfoByOperatorId(
	ctx context.Context,
	operatorID core.OperatorID,
	blockNumber uint32,
) (*core.IndexedOperatorInfo, error) {
	op, err := ics.store.GetOperator(ctx, operatorID)
	if err != nil {
		return nil, fmt.Errorf("failed to get operator %s: %w", operatorID.Hex(), err)
	}
	return toIndexedOperatorInfo(op)
}

// getQuorumAPK returns the aggregate public key for a quorum as of blockNumber,
// i.e. the latest snapshot at or before blockNumber, or nil if none exists.
// This replicates the subgraph's blockNumber_lte lookup semantics.
func (ics *IndexedChainState) getQuorumAPK(
	ctx context.Context,
	quorumID core.QuorumID,
	blockNumber uint64,
) (*core.G1Point, error) {
	apks, err := ics.store.ListQuorumAPKs(ctx, types.QuorumAPKFilter{
		QuorumID: quorumID,
		MaxBlock: blockNumber,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list quorum APKs: %w", err)
	}
	if len(apks) == 0 {
		return nil, nil
	}
	// ListQuorumAPKs returns snapshots sorted ascending by block number, so the
	// last entry is the latest snapshot at or before blockNumber.
	return apks[len(apks)-1].APK, nil
}

// registeredAt reports whether the operator was registered as of blockNumber.
func registeredAt(op *types.Operator, blockNumber uint64) bool {
	if op.RegisteredAtBlockNumber > blockNumber {
		return false
	}
	if op.DeregisteredAtBlockNumber != nil && *op.DeregisteredAtBlockNumber <= blockNumber {
		return false
	}
	return true
}

// toIndexedOperatorInfo projects a stored operator onto the core indexed-info
// shape. It errors if the operator is missing the BLS keys or socket that the
// indexed info requires.
func toIndexedOperatorInfo(op *types.Operator) (*core.IndexedOperatorInfo, error) {
	if op.BLSPubKeyG1 == nil || op.BLSPubKeyG2 == nil {
		return nil, fmt.Errorf("operator missing BLS public keys")
	}
	if op.Socket == "" {
		return nil, fmt.Errorf("operator missing socket")
	}
	return &core.IndexedOperatorInfo{
		PubkeyG1: op.BLSPubKeyG1,
		PubkeyG2: op.BLSPubKeyG2,
		Socket:   op.Socket,
	}, nil
}
