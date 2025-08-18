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

func (m *MockIndexedChainState) GetOperatorState(
	ctx context.Context,
	blockNumber uint,
	quorums []core.QuorumID) (*core.OperatorState, error) {

	args := m.Mock.Called(blockNumber, quorums)
	return args.Get(0).(*core.OperatorState), args.Error(1)
}

func (m *MockIndexedChainState) GetOperatorStateWithSocket(
	ctx context.Context,
	blockNumber uint,
	quorums []core.QuorumID) (*core.OperatorState, error) {

	args := m.Mock.Called(blockNumber, quorums)
	return args.Get(0).(*core.OperatorState), args.Error(1)
}

func (m *MockIndexedChainState) GetOperatorStateByOperator(
	ctx context.Context,
	blockNumber uint,
	operator core.OperatorID) (*core.OperatorState, error) {

	args := m.Mock.Called(blockNumber, operator)
	return args.Get(0).(*core.OperatorState), args.Error(1)
}

func (m *MockIndexedChainState) Start(context context.Context) error {
	args := m.Mock.Called()
	return args.Error(0)
}

func (m *MockIndexedChainState) GetCurrentBlockNumber(ctx context.Context) (uint, error) {
	args := m.Mock.Called()
	return args.Get(0).(uint), args.Error(1)
}

func (m *MockIndexedChainState) GetIndexedOperators(
	ctx context.Context,
	blockNumber uint) (map[core.OperatorID]*core.IndexedOperatorInfo, error) {

	args := m.Mock.Called(blockNumber)
	return args.Get(0).(map[core.OperatorID]*core.IndexedOperatorInfo), args.Error(1)
}

func (m *MockIndexedChainState) GetOperatorSocket(
	ctx context.Context,
	blockNumber uint,
	operator core.OperatorID) (string, error) {

	args := m.Mock.Called(blockNumber, operator)
	return args.Get(0).(string), args.Error(1)
}
