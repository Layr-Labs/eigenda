package subgraph

import (
	"context"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/shurcooL/graphql"
)

var (
	once               sync.Once
	instance           *api
	maxEntriesPerQuery = 1000
)

type (
	Api interface {
		QueryBatches(ctx context.Context, descending bool, orderByField string, first, skip int) ([]*Batches, error)
		QueryOperators(ctx context.Context, first int) ([]*Operator, error)
		QueryBatchNonSigningOperatorIdsInInterval(ctx context.Context, intervalSeconds int64) ([]*BatchNonSigningOperatorIds, error)
		QueryDeregisteredOperatorsGreaterThanBlockTimestampWithPagination(ctx context.Context, blockTimestamp uint64) ([]*Operator, error)
		QueryOperatorInfoByOperatorIdAtBlockNumber(ctx context.Context, operatorId core.OperatorID, blockNumber uint32) (*IndexedOperatorInfo, error)
	}

	api struct {
		uiMonitoringGgl  *graphql.Client
		operatorStateGql *graphql.Client
	}
)

var _ Api = (*api)(nil)

func NewApi(uiMonitoringSocketAddr string, operatorStateSocketAddr string) *api {
	once.Do(func() {
		uiMonitoringGgl := graphql.NewClient(uiMonitoringSocketAddr, nil)
		operatorStateGql := graphql.NewClient(operatorStateSocketAddr, nil)
		instance = &api{
			uiMonitoringGgl:  uiMonitoringGgl,
			operatorStateGql: operatorStateGql,
		}
	})
	return instance
}

func (a *api) QueryBatches(ctx context.Context, descending bool, orderByField string, first, skip int) ([]*Batches, error) {
	order := "asc"
	if descending {
		order = "desc"
	}
	variables := map[string]any{
		"orderDirection": graphql.String(order),
		"orderBy":        graphql.String(orderByField),
		"first":          graphql.Int(first),
		"skip":           graphql.Int(skip),
	}
	result := new(queryBatches)
	err := a.uiMonitoringGgl.Query(ctx, result, variables)
	if err != nil {
		return nil, err
	}

	return result.Batches, nil
}

func (a *api) QueryOperators(ctx context.Context, first int) ([]*Operator, error) {
	variables := map[string]any{
		"first": graphql.Int(first),
	}
	result := new(queryOperatorRegistereds)
	err := a.operatorStateGql.Query(ctx, result, variables)
	if err != nil {
		return nil, err
	}

	return result.OperatorRegistereds, nil
}

func (a *api) QueryBatchNonSigningOperatorIdsInInterval(ctx context.Context, intervalSeconds int64) ([]*BatchNonSigningOperatorIds, error) {
	nonSigningAfter := time.Now().Add(-time.Duration(intervalSeconds) * time.Second).Unix()
	variables := map[string]any{
		"blockTimestamp_gt": graphql.Int(nonSigningAfter),
	}
	skip := 0

	result := new(queryBatchNonSigningOperatorIdsInInterval)
	batchNonSigningOperatorIds := make([]*BatchNonSigningOperatorIds, 0)
	for {
		variables["first"] = graphql.Int(maxEntriesPerQuery)
		variables["skip"] = graphql.Int(skip)

		err := a.uiMonitoringGgl.Query(ctx, &result, variables)
		if err != nil {
			return nil, err
		}

		if len(result.BatchNonSigningOperatorIds) == 0 {
			break
		}
		batchNonSigningOperatorIds = append(batchNonSigningOperatorIds, result.BatchNonSigningOperatorIds...)

		skip += maxEntriesPerQuery
	}

	result.BatchNonSigningOperatorIds = batchNonSigningOperatorIds
	return result.BatchNonSigningOperatorIds, nil
}

func (a *api) QueryDeregisteredOperatorsGreaterThanBlockTimestampWithPagination(ctx context.Context, blockTimestamp uint64) ([]*Operator, error) {
	variables := map[string]any{
		"blockTimestamp_gt": graphql.Int(blockTimestamp),
	}
	skip := 0
	result := new(queryOperatorDeregistereds)
	operators := make([]*Operator, 0)
	for {
		variables["first"] = graphql.Int(maxEntriesPerQuery)
		variables["skip"] = graphql.Int(skip)

		err := a.operatorStateGql.Query(ctx, &result, variables)
		if err != nil {
			return nil, err
		}

		if len(result.OperatorDeregistereds) == 0 {
			break
		}
		operators = append(operators, result.OperatorDeregistereds...)
		skip += maxEntriesPerQuery
	}
	return operators, nil
}

func (a *api) QueryOperatorInfoByOperatorIdAtBlockNumber(ctx context.Context, operatorId core.OperatorID, blockNumber uint32) (*IndexedOperatorInfo, error) {
	var (
		query     queryOperatorById
		variables = map[string]any{
			"id": graphql.String(fmt.Sprintf("0x%s", hex.EncodeToString(operatorId[:]))),
		}
	)
	err := a.operatorStateGql.Query(context.Background(), &query, variables)
	if err != nil {
		return nil, err
	}

	return &query.Operator, nil
}
