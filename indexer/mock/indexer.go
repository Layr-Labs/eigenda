package mock

import (
	"context"

	"github.com/Layr-Labs/eigenda/indexer"
	"github.com/stretchr/testify/mock"
)

type MockIndexer struct {
	mock.Mock
}

var _ indexer.Indexer = (*MockIndexer)(nil)

func (m *MockIndexer) Index(ctx context.Context) error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockIndexer) HandleAccumulator(acc indexer.Accumulator, f indexer.Filterer, headers indexer.Headers) error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockIndexer) GetLatestHeader(finalized bool) (*indexer.Header, error) {
	args := m.Called(finalized)
	return args.Get(0).(*indexer.Header), args.Error(1)
}

func (m *MockIndexer) GetObject(header *indexer.Header, handlerIndex int) (indexer.AccumulatorObject, error) {
	args := m.Called(header, handlerIndex)
	return args.Get(0).(indexer.AccumulatorObject), args.Error(1)
}
