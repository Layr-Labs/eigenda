package mock

import (
	"context"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/stretchr/testify/mock"
)

var _ core.IndexedChainState = (*IndexedChainState)(nil)

// IndexedChainState is a mock implementation of core.IndexedChainState.
type IndexedChainState struct {
	Mock mock.Mock
}

func (m *IndexedChainState) GetCurrentBlockNumber() (uint, error) {
	args := m.Mock.Called()
	return args.Get(0).(uint), args.Error(1)
}

func (m *IndexedChainState) GetOperatorState(
	ctx context.Context,
	blockNumber uint,
	quorums []core.QuorumID) (*core.OperatorState, error) {

	args := m.Mock.Called(blockNumber, quorums)
	return args.Get(0).(*core.OperatorState), args.Error(1)
}

func (m *IndexedChainState) GetOperatorStateByOperator(
	ctx context.Context,
	blockNumber uint,
	operator core.OperatorID) (*core.OperatorState, error) {

	args := m.Mock.Called(blockNumber, operator)
	return args.Get(0).(*core.OperatorState), args.Error(1)
}

func (m *IndexedChainState) GetOperatorSocket(
	ctx context.Context,
	blockNumber uint,
	operator core.OperatorID) (string, error) {

	args := m.Mock.Called(blockNumber, operator)
	return args.Get(0).(string), args.Error(1)
}

func (m *IndexedChainState) GetIndexedOperatorState(
	ctx context.Context,
	blockNumber uint,
	quorums []core.QuorumID) (*core.IndexedOperatorState, error) {

	args := m.Mock.Called(blockNumber, quorums)
	return args.Get(0).(*core.IndexedOperatorState), args.Error(1)
}

func (m *IndexedChainState) GetIndexedOperators(
	ctx context.Context,
	blockNumber uint) (map[core.OperatorID]*core.IndexedOperatorInfo, error) {

	args := m.Mock.Called(blockNumber)
	return args.Get(0).(map[core.OperatorID]*core.IndexedOperatorInfo), args.Error(1)
}

func (m *IndexedChainState) Start(context context.Context) error {
	args := m.Mock.Called()
	return args.Error(0)
}
