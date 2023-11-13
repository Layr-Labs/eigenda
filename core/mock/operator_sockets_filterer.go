package mock

import (
	"context"

	"github.com/Layr-Labs/eigenda/core"
	coreindexer "github.com/Layr-Labs/eigenda/core/indexer"
	"github.com/Layr-Labs/eigenda/indexer"
	"github.com/stretchr/testify/mock"
)

type MockOperatorSocketsFilterer struct {
	mock.Mock
}

var _ coreindexer.OperatorSocketsFilterer = (*MockOperatorSocketsFilterer)(nil)

func (t *MockOperatorSocketsFilterer) FilterHeaders(headers indexer.Headers) ([]indexer.HeaderAndEvents, error) {
	args := t.Called()
	result := args.Get(0)
	return result.([]indexer.HeaderAndEvents), args.Error(1)
}

func (t *MockOperatorSocketsFilterer) GetSyncPoint(latestHeader *indexer.Header) (uint64, error) {
	args := t.Called()
	result := args.Get(0)
	return result.(uint64), args.Error(1)
}

func (t *MockOperatorSocketsFilterer) SetSyncPoint(latestHeader *indexer.Header) error {
	args := t.Called()
	return args.Error(0)
}

func (t *MockOperatorSocketsFilterer) FilterFastMode(headers indexer.Headers) (*indexer.Header, indexer.Headers, error) {
	args := t.Called()
	result1 := args.Get(0)
	result2 := args.Get(1)
	return result1.(*indexer.Header), result2.(indexer.Headers), args.Error(2)
}

func (t *MockOperatorSocketsFilterer) WatchOperatorSocketUpdate(ctx context.Context, operatorId core.OperatorID) (chan string, error) {
	args := t.Called()
	result := args.Get(0)
	return result.(chan string), args.Error(1)
}
