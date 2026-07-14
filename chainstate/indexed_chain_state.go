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
// The three pieces are combined by core.AssembleIndexedOperatorState, which is
// shared with core/thegraph so both implementations apply identical
// missing/extra-operator rules.
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

	state, err := core.AssembleIndexedOperatorState(operatorState, indexedOperators, aggKeys)
	if err != nil {
		return nil, fmt.Errorf("failed to assemble indexed operator state: %w", err)
	}
	return state, nil
}

// GetIndexedOperators returns the indexed info (BLS keys and socket) for every
// operator not deregistered as of blockNumber, keyed by operator ID.
//
// Deliberately mirroring the operator-state subgraph, this filters ONLY on the
// deregistration block (the subgraph's deregistrationBlockNumber_gt condition)
// and not on the registration block. That means it over-fetches: an operator
// that registered after blockNumber, or re-registered after an earlier stint
// that covered blockNumber, is still included. Over-fetching is safe because
// GetIndexedOperatorState filters the result against the on-chain operator
// state; filtering on the stored registration block would instead UNDER-fetch
// after a re-registration (the stored block is the latest registration, hiding
// earlier registration windows) and make assembly fail with a missing operator.
func (ics *IndexedChainState) GetIndexedOperators(
	ctx context.Context,
	blockNumber uint,
) (map[core.OperatorID]*core.IndexedOperatorInfo, error) {
	operators, err := ics.store.ListOperators(ctx, types.OperatorFilter{}, 0, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to list operators: %w", err)
	}

	result := make(map[core.OperatorID]*core.IndexedOperatorInfo, len(operators))
	for _, op := range operators {
		if deregisteredAsOf(op, uint64(blockNumber)) {
			continue
		}
		info, err := toIndexedOperatorInfo(op)
		if err != nil {
			// An incomplete record (missing keys or socket) can exist when the
			// operator's one-time pubkey registration predates the configured
			// start block, or transiently while the indexer is mid-way through
			// the events of a registration transaction. Skip it rather than
			// failing the whole call: if the operator matters for the requested
			// quorums, AssembleIndexedOperatorState reports it as missing; if it
			// doesn't, one bad record must not poison every query.
			ics.logger.Warn("Skipping operator with incomplete indexed record",
				"operator_id", op.ID.Hex(), "error", err)
			continue
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
	apk, err := ics.store.GetLatestQuorumAPK(ctx, quorumID, blockNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest quorum APK: %w", err)
	}
	if apk == nil {
		return nil, nil
	}
	return apk.APK, nil
}

// deregisteredAsOf reports whether the operator had deregistered at or before
// blockNumber (and has not re-registered since; a re-registration clears the
// deregistration marker). This is the inverse of the subgraph's
// deregistrationBlockNumber_gt filter.
func deregisteredAsOf(op *types.Operator, blockNumber uint64) bool {
	return op.DeregisteredAtBlockNumber != nil && *op.DeregisteredAtBlockNumber <= blockNumber
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
