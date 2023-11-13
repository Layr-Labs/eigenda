package mock

import (
	"context"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/thegraph"
	"github.com/stretchr/testify/mock"
)

type MockIndexedChainState struct {
	mock.Mock
}

var _ thegraph.IndexedChainState = (*MockIndexedChainState)(nil)

func (m *MockIndexedChainState) GetIndexedOperatorState(ctx context.Context, blockNumber uint, quorums []core.QuorumID) (*core.IndexedOperatorState, error) {
	args := m.Called()
	var value *core.IndexedOperatorState
	if args.Get(0) != nil {
		value = args.Get(0).(*core.IndexedOperatorState)
	}
	return value, args.Error(1)
}

func (m *MockIndexedChainState) GetIndexedOperatorInfoByOperatorId(ctx context.Context, operatorId core.OperatorID, blockNumber uint32) (*core.IndexedOperatorInfo, error) {
	args := m.Called()
	var value *core.IndexedOperatorInfo
	if args.Get(0) != nil {
		value = args.Get(0).(*core.IndexedOperatorInfo)
	}
	return value, args.Error(1)
}
