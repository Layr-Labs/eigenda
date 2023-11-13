package indexer

import (
	"context"
	"errors"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/indexer"
	"github.com/Layr-Labs/eigenda/indexer/eth"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

type IndexedChainState struct {
	core.ChainState

	Indexer *indexer.Indexer
}

var _ core.IndexedChainState = (*IndexedChainState)(nil)

// TODO: Pass in dependencies instead of creating them here

func NewIndexedChainState(
	config *indexer.Config,
	eigenDAServiceManagerAddr gethcommon.Address,
	chainState core.ChainState,
	headerStore indexer.HeaderStore,
	client common.EthClient,
	rpcClient common.RPCEthClient,
	logger common.Logger,
) (*IndexedChainState, error) {

	pubKeyFilterer, err := NewOperatorPubKeysFilterer(eigenDAServiceManagerAddr, client)
	if err != nil {
		return nil, err
	}

	socketsFilterer, err := NewOperatorSocketsFilterer(eigenDAServiceManagerAddr, client)
	if err != nil {
		return nil, err
	}

	handlers := []indexer.AccumulatorHandler{
		{
			Acc:      NewOperatorPubKeysAccumulator(logger),
			Filterer: pubKeyFilterer,
			Status:   indexer.Good,
		},
		{
			Acc:      NewOperatorSocketsAccumulator(logger),
			Filterer: socketsFilterer,
			Status:   indexer.Good,
		},
	}

	headerSrvc := eth.NewHeaderService(logger, rpcClient)
	upgrader := &Upgrader{}
	indexer := indexer.NewIndexer(
		config,
		handlers,
		headerSrvc,
		headerStore,
		upgrader,
		logger,
	)

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
			return nil, errors.New("aggregate key for quorum not found")
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
	header, err := ics.Indexer.HeaderStore.GetLatestHeader(false)
	if err != nil {
		return 0, err
	}
	return uint(header.Number), nil
}

func (ics *IndexedChainState) getObjects(blockNumber uint) (*OperatorPubKeys, OperatorSockets, error) {

	queryHeader := &indexer.Header{
		Number: uint64(blockNumber),
	}

	obj, _, err := ics.Indexer.HeaderStore.GetObject(queryHeader, ics.Indexer.Handlers[0].Acc)
	if err != nil {
		return nil, nil, err
	}

	pubkeys, ok := obj.(*OperatorPubKeys)
	if !ok {
		return nil, nil, ErrWrongObjectFromIndexer
	}

	obj, _, err = ics.Indexer.HeaderStore.GetObject(queryHeader, ics.Indexer.Handlers[1].Acc)
	if err != nil {
		return nil, nil, err
	}

	sockets, ok := obj.(OperatorSockets)
	if !ok {
		return nil, nil, ErrWrongObjectFromIndexer
	}

	return pubkeys, sockets, nil

}
