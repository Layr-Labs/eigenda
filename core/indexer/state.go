package indexer

import (
	"context"
	"errors"
	"fmt"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/indexer"
)

type IndexedChainState struct {
	core.ChainState

	Indexer indexer.Indexer
}

var _ core.IndexedChainState = (*IndexedChainState)(nil)

func NewIndexedChainState(
	chainState core.ChainState,
	indexer indexer.Indexer,
) (*IndexedChainState, error) {

	return &IndexedChainState{
		ChainState: chainState,
		Indexer:    indexer,
	}, nil
}

func (ics *IndexedChainState) Start(ctx context.Context) error {
	return ics.Indexer.Index(ctx)
}

func (ics *IndexedChainState) GetIndexedOperatorState(ctx context.Context, blockNumber uint, quorums []core.QuorumID) (*core.IndexedOperatorState, error) {

	pubkeys, sockets, err := ics.getObjects(blockNumber)
	if err != nil {
		return nil, fmt.Errorf("unable to complete ics.getObjects(%d): %s)", blockNumber, err)
	}

	operatorState, err := ics.ChainState.GetOperatorState(ctx, blockNumber, quorums)
	if err != nil {
		return nil, fmt.Errorf("unable to complete ics.ChainState.GetOperatorState(%d, %v): %s", blockNumber, quorums, err)
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

	return state, nil
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
		return nil, nil, fmt.Errorf("(1) unable to call Indexer.GetObject({Number = %d}, 0): %s", blockNumber, err)
	}

	pubkeys, ok := obj.(*OperatorPubKeys)
	if !ok {
		return nil, nil, ErrWrongObjectFromIndexer
	}

	obj, err = ics.Indexer.GetObject(queryHeader, 1)
	if err != nil {
		return nil, nil, fmt.Errorf("(2) unable to call Indexer.GetObject(%v, 1): %s", queryHeader, err)
	}

	sockets, ok := obj.(OperatorSockets)
	if !ok {
		return nil, nil, ErrWrongObjectFromIndexer
	}

	return pubkeys, sockets, nil

}
