package mock

import (
	"context"

	"github.com/Layr-Labs/eigenda/chainio"
	"github.com/Layr-Labs/eigenda/chainio/thegraph"
	"github.com/stretchr/testify/mock"
)

type MockIndexedChainState struct {
	mock.Mock
}

var _ thegraph.IndexedChainState = (*MockIndexedChainState)(nil)

func (m *MockIndexedChainState) GetIndexedOperatorState(ctx context.Context, blockNumber uint, quorums []chainio.QuorumID) (*chainio.IndexedOperatorState, error) {
	args := m.Called()
	var value *chainio.IndexedOperatorState
	if args.Get(0) != nil {
		value = args.Get(0).(*chainio.IndexedOperatorState)
	}
	return value, args.Error(1)
}

func (m *MockIndexedChainState) GetIndexedOperatorInfoByOperatorId(ctx context.Context, operatorId chainio.OperatorID, blockNumber uint32) (*chainio.IndexedOperatorInfo, error) {
	args := m.Called()
	var value *chainio.IndexedOperatorInfo
	if args.Get(0) != nil {
		value = args.Get(0).(*chainio.IndexedOperatorInfo)
	}
	return value, args.Error(1)
}
