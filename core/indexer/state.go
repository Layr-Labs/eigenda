package indexer

import (
	"context"
	"encoding/binary"
	"errors"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/indexer"
	lru "github.com/hashicorp/golang-lru/v2"
)

type IndexedChainState struct {
	core.ChainState

	Indexer indexer.Indexer

	operatorStateCache *lru.Cache[string, *core.IndexedOperatorState]
}

var _ core.IndexedChainState = (*IndexedChainState)(nil)

func NewIndexedChainState(
	chainState core.ChainState,
	indexer indexer.Indexer,
	cacheSize int,
) (*IndexedChainState, error) {
	operatorStateCache := (*lru.Cache[string, *core.IndexedOperatorState])(nil)
	var err error
	if cacheSize > 0 {
		operatorStateCache, err = lru.New[string, *core.IndexedOperatorState](cacheSize)
		if err != nil {
			return nil, err
		}
	}
	return &IndexedChainState{
		ChainState: chainState,
		Indexer:    indexer,

		operatorStateCache: operatorStateCache,
	}, nil
}

func (ics *IndexedChainState) Start(ctx context.Context) error {
	return ics.Indexer.Index(ctx)
}

func (ics *IndexedChainState) GetIndexedOperatorState(ctx context.Context, blockNumber uint, quorums []core.QuorumID) (*core.IndexedOperatorState, error) {
	cacheKey := computeCacheKey(blockNumber, quorums)
	if ics.operatorStateCache != nil {
		if val, ok := ics.operatorStateCache.Get(cacheKey); ok {
			return val, nil
		}
	}

	pubkeys, sockets, err := ics.getObjects(blockNumber)
	if err != nil {
		return nil, err
	}

	operatorState, err := ics.ChainState.GetOperatorState(ctx, blockNumber, quorums)
	if err != nil {
		return nil, err
	}

	ops := make(map[core.OperatorID]*core.IndexedOperatorInfo, len(pubkeys.Operators))
	for id, op := range pubkeys.Operators {

		socket, ok := sockets[id]
		if !ok {
			return nil, errors.New("socket for operator not found")
		}

		ops[id] = &core.IndexedOperatorInfo{
			PubkeyG1: &core.G1Point{G1Affine: op.PubKeyG1},
			PubkeyG2: &core.G2Point{G2Affine: op.PubKeyG2},
			Socket:   socket,
		}
	}

	aggKeys := make(map[core.QuorumID]*core.G1Point, len(pubkeys.Operators))
	for _, quorum := range quorums {
		key, ok := pubkeys.QuorumTotals[quorum]
		if !ok {
			continue
		}
		aggKeys[quorum] = &core.G1Point{G1Affine: key}
	}

	state := &core.IndexedOperatorState{
		OperatorState:    operatorState,
		IndexedOperators: ops,
		AggKeys:          aggKeys,
	}

	if ics.operatorStateCache != nil {
		ics.operatorStateCache.Add(cacheKey, state)
	}

	return state, nil
}

func (ics *IndexedChainState) GetIndexedOperators(ctx context.Context, blockNumber uint) (map[core.OperatorID]*core.IndexedOperatorInfo, error) {

	pubkeys, sockets, err := ics.getObjects(blockNumber)
	if err != nil {
		return nil, err
	}

	ops := make(map[core.OperatorID]*core.IndexedOperatorInfo, len(pubkeys.Operators))
	for id, op := range pubkeys.Operators {

		socket, ok := sockets[id]
		if !ok {
			return nil, errors.New("socket for operator not found")
		}

		ops[id] = &core.IndexedOperatorInfo{
			PubkeyG1: &core.G1Point{G1Affine: op.PubKeyG1},
			PubkeyG2: &core.G2Point{G2Affine: op.PubKeyG2},
			Socket:   socket,
		}
	}

	return ops, nil
}

func (ics *IndexedChainState) GetCurrentBlockNumber() (uint, error) {
	header, err := ics.Indexer.GetLatestHeader(false)
	if err != nil {
		return 0, err
	}
	return uint(header.Number), nil
}

func (ics *IndexedChainState) getObjects(blockNumber uint) (*OperatorPubKeys, OperatorSockets, error) {

	queryHeader := &indexer.Header{
		Number: uint64(blockNumber),
	}

	obj, err := ics.Indexer.GetObject(queryHeader, 0)
	if err != nil {
		return nil, nil, err
	}

	pubkeys, ok := obj.(*OperatorPubKeys)
	if !ok {
		return nil, nil, ErrWrongObjectFromIndexer
	}

	obj, err = ics.Indexer.GetObject(queryHeader, 1)
	if err != nil {
		return nil, nil, err
	}

	sockets, ok := obj.(OperatorSockets)
	if !ok {
		return nil, nil, ErrWrongObjectFromIndexer
	}

	return pubkeys, sockets, nil

}

func computeCacheKey(blockNumber uint, quorumIDs []uint8) string {
	bytes := make([]byte, 8+len(quorumIDs))
	binary.LittleEndian.PutUint64(bytes, uint64(blockNumber))
	copy(bytes[8:], quorumIDs)
	return string(bytes)
}
