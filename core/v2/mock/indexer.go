package mock

import (
	"context"

	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/core/v2/thegraph"
	"github.com/stretchr/testify/mock"
)

type MockIndexedChainState struct {
	mock.Mock
}

var _ thegraph.IndexedChainState = (*MockIndexedChainState)(nil)

func (m *MockIndexedChainState) GetIndexedOperatorState(ctx context.Context, blockNumber uint, quorums []corev2.QuorumID) (*corev2.IndexedOperatorState, error) {
	args := m.Called()
	var value *corev2.IndexedOperatorState
	if args.Get(0) != nil {
		value = args.Get(0).(*corev2.IndexedOperatorState)
	}
	return value, args.Error(1)
}

func (m *MockIndexedChainState) GetIndexedOperatorInfoByOperatorId(ctx context.Context, operatorId corev2.OperatorID, blockNumber uint32) (*corev2.IndexedOperatorInfo, error) {
	args := m.Called()
	var value *corev2.IndexedOperatorInfo
	if args.Get(0) != nil {
		value = args.Get(0).(*corev2.IndexedOperatorInfo)
	}
	return value, args.Error(1)
}
