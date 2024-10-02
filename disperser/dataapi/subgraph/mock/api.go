package mock

import (
	"cmp"
	"context"
	"slices"

	"github.com/Layr-Labs/eigenda/disperser/dataapi/subgraph"
	"github.com/stretchr/testify/mock"
)

type MockSubgraphApi struct {
	mock.Mock
}

var _ subgraph.Api = (*MockSubgraphApi)(nil)

func (m *MockSubgraphApi) QueryBatches(ctx context.Context, descending bool, orderByField string, first, skip int) ([]*subgraph.Batches, error) {
	args := m.Called()

	var value []*subgraph.Batches
	if args.Get(0) != nil {
		value = args.Get(0).([]*subgraph.Batches)

		if orderByField == "blockTimestamp" {
			slices.SortStableFunc(value, func(a, b *subgraph.Batches) int {
				return cmp.Compare(a.BlockTimestamp, b.BlockTimestamp)
			})
		}
		if descending {
			slices.Reverse(value)
		}
		if skip > 0 && len(value) > skip {
			value = value[skip:]
		}
		if first > 0 && len(value) > first {
			value = value[:first]
		}
	}

	return value, args.Error(1)
}

func (m *MockSubgraphApi) QueryBatchesByBlockTimestampRange(ctx context.Context, start, end uint64) ([]*subgraph.Batches, error) {
	args := m.Called()
	var value []*subgraph.Batches
	if args.Get(0) != nil {
		value = args.Get(0).([]*subgraph.Batches)
	}

	return value, args.Error(1)
}

func (m *MockSubgraphApi) QueryOperators(ctx context.Context, first int) ([]*subgraph.Operator, error) {
	args := m.Called()

	var value []*subgraph.Operator
	if args.Get(0) != nil {
		value = args.Get(0).([]*subgraph.Operator)

		if len(value) > first {
			value = value[:first]
		}
	}

	return value, args.Error(1)
}

func (m *MockSubgraphApi) QueryOperatorsDeregistered(ctx context.Context, first int) ([]*subgraph.Operator, error) {
	args := m.Called()

	var value []*subgraph.Operator
	if args.Get(0) != nil {
		value = args.Get(0).([]*subgraph.Operator)

		if len(value) > first {
			value = value[:first]
		}
	}

	return value, args.Error(1)
}

func (m *MockSubgraphApi) QueryBatchNonSigningInfo(ctx context.Context, startTime, endTime int64) ([]*subgraph.BatchNonSigningInfo, error) {
	args := m.Called(startTime, endTime)

	var value []*subgraph.BatchNonSigningInfo
	if args.Get(0) != nil {
		value = args.Get(0).([]*subgraph.BatchNonSigningInfo)
	}

	return value, args.Error(1)
}

func (m *MockSubgraphApi) QueryBatchNonSigningOperatorIdsInInterval(ctx context.Context, first int64) ([]*subgraph.BatchNonSigningOperatorIds, error) {
	args := m.Called()

	var value []*subgraph.BatchNonSigningOperatorIds
	if args.Get(0) != nil {
		value = args.Get(0).([]*subgraph.BatchNonSigningOperatorIds)

		if len(value) > int(first) {
			value = value[:first]
		}
	}

	return value, args.Error(1)
}

func (m *MockSubgraphApi) QueryRegisteredOperatorsGreaterThanBlockTimestamp(ctx context.Context, blockTimestamp uint64) ([]*subgraph.Operator, error) {
	args := m.Called()

	var value []*subgraph.Operator
	if args.Get(0) != nil {
		value = args.Get(0).([]*subgraph.Operator)
	}

	return value, args.Error(1)
}

func (m *MockSubgraphApi) QueryDeregisteredOperatorsGreaterThanBlockTimestamp(ctx context.Context, blockTimestamp uint64) ([]*subgraph.Operator, error) {
	args := m.Called()

	var value []*subgraph.Operator
	if args.Get(0) != nil {
		value = args.Get(0).([]*subgraph.Operator)
	}

	return value, args.Error(1)
}

func (m *MockSubgraphApi) QueryOperatorInfoByOperatorIdAtBlockNumber(ctx context.Context, operatorId string, blockNumber uint32) (*subgraph.IndexedOperatorInfo, error) {
	args := m.Called()

	var value *subgraph.IndexedOperatorInfo
	if args.Get(0) != nil {
		value = args.Get(0).(*subgraph.IndexedOperatorInfo)
	}

	return value, args.Error(1)
}

func (m *MockSubgraphApi) QueryOperatorAddedToQuorum(ctx context.Context, startBlock, endBlock uint32) ([]*subgraph.OperatorQuorum, error) {
	args := m.Called()

	var value []*subgraph.OperatorQuorum
	if args.Get(0) != nil {
		value = args.Get(0).([]*subgraph.OperatorQuorum)
	}

	return value, args.Error(1)
}

func (m *MockSubgraphApi) QueryOperatorRemovedFromQuorum(ctx context.Context, startBlock, endBlock uint32) ([]*subgraph.OperatorQuorum, error) {
	args := m.Called()

	var value []*subgraph.OperatorQuorum
	if args.Get(0) != nil {
		value = args.Get(0).([]*subgraph.OperatorQuorum)
	}

	return value, args.Error(1)
}
