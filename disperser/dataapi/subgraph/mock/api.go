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

func (m *MockSubgraphApi) QueryOperators(ctx context.Context, first int) ([]*subgraph.OperatorRegistered, error) {
	args := m.Called()

	var value []*subgraph.OperatorRegistered
	if args.Get(0) != nil {
		value = args.Get(0).([]*subgraph.OperatorRegistered)

		if len(value) > first {
			value = value[:first]
		}
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
