package indexer_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/indexer"
	coremock "github.com/Layr-Labs/eigenda/core/mock"
	indexermock "github.com/Layr-Labs/eigenda/indexer/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type testComponents struct {
	ChainState        *coremock.ChainDataMock
	Indexer           *indexermock.MockIndexer
	IndexedChainState *indexer.IndexedChainState
}

func TestIndexedOperatorStateCache(t *testing.T) {
	c := createTestComponents(t)
	pubKeys := &indexer.OperatorPubKeys{}
	c.Indexer.On("GetObject", mock.Anything, 0).Return(pubKeys, nil)
	sockets := indexer.OperatorSockets{
		core.OperatorID{0, 1}: "socket1",
	}
	c.Indexer.On("GetObject", mock.Anything, 1).Return(sockets, nil)

	operatorState := &core.OperatorState{
		Operators: map[core.QuorumID]map[core.OperatorID]*core.OperatorInfo{
			0: {
				core.OperatorID{0}: {
					Stake: big.NewInt(100),
					Index: 0,
				},
			},
		},
	}
	c.ChainState.On("GetOperatorState", mock.Anything, uint(100), []core.QuorumID{0}).Return(operatorState, nil)
	c.ChainState.On("GetOperatorState", mock.Anything, uint(100), []core.QuorumID{1}).Return(operatorState, nil)
	c.ChainState.On("GetOperatorState", mock.Anything, uint(101), []core.QuorumID{0, 1}).Return(operatorState, nil)

	ctx := context.Background()
	// Get the operator state for block 100 and quorum 0
	_, err := c.IndexedChainState.GetIndexedOperatorState(ctx, uint(100), []core.QuorumID{0})
	assert.NoError(t, err)
	c.ChainState.AssertNumberOfCalls(t, "GetOperatorState", 1)

	// Get the operator state for block 100 and quorum 0 again
	_, err = c.IndexedChainState.GetIndexedOperatorState(ctx, uint(100), []core.QuorumID{0})
	assert.NoError(t, err)
	c.ChainState.AssertNumberOfCalls(t, "GetOperatorState", 1)

	// Get the operator state for block 100 and quorum 1
	_, err = c.IndexedChainState.GetIndexedOperatorState(ctx, uint(100), []core.QuorumID{1})
	assert.NoError(t, err)
	c.ChainState.AssertNumberOfCalls(t, "GetOperatorState", 2)

	// Get the operator state for block 101 and quorum 0 & 1
	_, err = c.IndexedChainState.GetIndexedOperatorState(ctx, uint(101), []core.QuorumID{0, 1})
	assert.NoError(t, err)
	c.ChainState.AssertNumberOfCalls(t, "GetOperatorState", 3)

	// Get the operator state for block 101 and quorum 0 & 1 again
	_, err = c.IndexedChainState.GetIndexedOperatorState(ctx, uint(101), []core.QuorumID{0, 1})
	assert.NoError(t, err)
	c.ChainState.AssertNumberOfCalls(t, "GetOperatorState", 3)
}

func createTestComponents(t *testing.T) *testComponents {
	chainState := &coremock.ChainDataMock{}
	idx := &indexermock.MockIndexer{}
	ics, err := indexer.NewIndexedChainState(chainState, idx, 1)
	assert.NoError(t, err)
	return &testComponents{
		ChainState:        chainState,
		Indexer:           idx,
		IndexedChainState: ics,
	}
}
