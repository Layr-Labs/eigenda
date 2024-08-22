package subgraph

import (
	"context"
	"fmt"
	"sync"
	"time"

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
		QueryBatchesByBlockTimestampRange(ctx context.Context, start, end uint64) ([]*Batches, error)
		QueryOperators(ctx context.Context, first int) ([]*Operator, error)
		QueryBatchNonSigningOperatorIdsInInterval(ctx context.Context, intervalSeconds int64) ([]*BatchNonSigningOperatorIds, error)
		QueryBatchNonSigningInfo(ctx context.Context, startTime, endTime int64) ([]*BatchNonSigningInfo, error)
		QueryDeregisteredOperatorsGreaterThanBlockTimestamp(ctx context.Context, blockTimestamp uint64) ([]*Operator, error)
		QueryRegisteredOperatorsGreaterThanBlockTimestamp(ctx context.Context, blockTimestamp uint64) ([]*Operator, error)
		QueryOperatorInfoByOperatorIdAtBlockNumber(ctx context.Context, operatorId string, blockNumber uint32) (*IndexedOperatorInfo, error)
		QueryOperatorAddedToQuorum(ctx context.Context, startBlock, endBlock uint32) ([]*OperatorQuorum, error)
		QueryOperatorRemovedFromQuorum(ctx context.Context, startBlock, endBlock uint32) ([]*OperatorQuorum, error)
	}

	api struct {
		uiMonitoringGql  *graphql.Client
		operatorStateGql *graphql.Client
	}
)

var _ Api = (*api)(nil)

func NewApi(uiMonitoringSocketAddr string, operatorStateSocketAddr string) *api {
	once.Do(func() {
		uiMonitoringGql := graphql.NewClient(uiMonitoringSocketAddr, nil)
		operatorStateGql := graphql.NewClient(operatorStateSocketAddr, nil)
		instance = &api{
			uiMonitoringGql:  uiMonitoringGql,
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
	err := a.uiMonitoringGql.Query(ctx, result, variables)
	if err != nil {
		return nil, err
	}

	return result.Batches, nil
}

func (a *api) QueryBatchesByBlockTimestampRange(ctx context.Context, start, end uint64) ([]*Batches, error) {
	variables := map[string]any{
		"first":              graphql.Int(maxEntriesPerQuery),
		"blockTimestamp_gte": graphql.Int(start),
		"blockTimestamp_lte": graphql.Int(end),
	}
	skip := 0
	query := new(queryBatchesByBlockTimestampRange)
	result := make([]*Batches, 0)
	for {
		variables["first"] = graphql.Int(maxEntriesPerQuery)
		variables["skip"] = graphql.Int(skip)

		err := a.uiMonitoringGql.Query(ctx, &query, variables)
		if err != nil {
			return nil, err
		}

		if len(query.Batches) == 0 {
			break
		}
		result = append(result, query.Batches...)
		skip += maxEntriesPerQuery
	}

	return result, nil
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

func (a *api) QueryOperatorDeregistrations(ctx context.Context, first int) ([]*Operator, error) {
	variables := map[string]any{
		"first": graphql.Int(first),
	}
	result := new(queryOperatorDeregistereds)
	err := a.operatorStateGql.Query(ctx, result, variables)
	if err != nil {
		return nil, err
	}

	return result.OperatorRegistereds, nil
}

func (a *api) QueryBatchNonSigningInfo(ctx context.Context, startTime, endTime int64) ([]*BatchNonSigningInfo, error) {

	variables := map[string]any{
		"blockTimestamp_gt": graphql.Int(startTime),
		"blockTimestamp_lt": graphql.Int(endTime),
	}
	skip := 0

	result := new(queryBatchNonSigningInfo)
	batchNonSigningInfo := make([]*BatchNonSigningInfo, 0)
	for {
		variables["first"] = graphql.Int(maxEntriesPerQuery)
		variables["skip"] = graphql.Int(skip)

		err := a.uiMonitoringGql.Query(ctx, &result, variables)
		if err != nil {
			return nil, err
		}

		if len(result.BatchNonSigningInfo) == 0 {
			break
		}
		batchNonSigningInfo = append(batchNonSigningInfo, result.BatchNonSigningInfo...)

		skip += maxEntriesPerQuery
	}

	return batchNonSigningInfo, nil
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

		err := a.uiMonitoringGql.Query(ctx, &result, variables)
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

func (a *api) QueryRegisteredOperatorsGreaterThanBlockTimestamp(ctx context.Context, blockTimestamp uint64) ([]*Operator, error) {
	variables := map[string]any{
		"blockTimestamp_gt": graphql.Int(blockTimestamp),
	}
	query := new(queryOperatorRegisteredsGTBlockTimestamp)
	err := a.operatorStateGql.Query(ctx, &query, variables)
	if err != nil {
		return nil, err
	}
	return query.OperatorRegistereds, nil
}

func (a *api) QueryDeregisteredOperatorsGreaterThanBlockTimestamp(ctx context.Context, blockTimestamp uint64) ([]*Operator, error) {
	variables := map[string]any{
		"blockTimestamp_gt": graphql.Int(blockTimestamp),
	}
	query := new(queryOperatorDeregisteredsGTBlockTimestamp)
	err := a.operatorStateGql.Query(ctx, &query, variables)
	if err != nil {
		return nil, err
	}
	return query.OperatorDeregistereds, nil
}

func (a *api) QueryOperatorInfoByOperatorIdAtBlockNumber(ctx context.Context, operatorId string, blockNumber uint32) (*IndexedOperatorInfo, error) {
	var (
		query     queryOperatorById
		variables = map[string]any{
			"id": graphql.String(fmt.Sprintf("0x%s", operatorId)),
		}
	)
	err := a.operatorStateGql.Query(context.Background(), &query, variables)
	if err != nil {
		return nil, err
	}

	return &query.Operator, nil
}

// QueryOperatorAddedToQuorum finds operators' quorum opt-in history in range [startBlock, endBlock].
func (a *api) QueryOperatorAddedToQuorum(ctx context.Context, startBlock, endBlock uint32) ([]*OperatorQuorum, error) {
	if startBlock > endBlock {
		return nil, fmt.Errorf("endBlock must be no less than startBlock, startBlock: %d, endBlock: %d", startBlock, endBlock)
	}
	variables := map[string]any{
		"blockNumber_gt": graphql.Int(startBlock - 1),
		"blockNumber_lt": graphql.Int(endBlock + 1),
	}
	skip := 0
	result := new(queryOperatorAddedToQuorum)
	addedToQuorums := make([]*OperatorQuorum, 0)
	for {
		variables["first"] = graphql.Int(maxEntriesPerQuery)
		variables["skip"] = graphql.Int(skip)
		err := a.operatorStateGql.Query(ctx, &result, variables)
		if err != nil {
			return nil, err
		}
		if len(result.OperatorAddedToQuorum) == 0 {
			break
		}
		addedToQuorums = append(addedToQuorums, result.OperatorAddedToQuorum...)
		skip += maxEntriesPerQuery
	}
	return addedToQuorums, nil
}

// QueryOperatorRemovedFromQuorum finds operators' quorum opt-out history in range [startBlock, endBlock].
func (a *api) QueryOperatorRemovedFromQuorum(ctx context.Context, startBlock, endBlock uint32) ([]*OperatorQuorum, error) {
	if startBlock > endBlock {
		return nil, fmt.Errorf("endBlock must be no less than startBlock, startBlock: %d, endBlock: %d", startBlock, endBlock)
	}
	variables := map[string]any{
		"blockNumber_gt": graphql.Int(startBlock - 1),
		"blockNumber_lt": graphql.Int(endBlock + 1),
	}
	skip := 0
	result := new(queryOperatorRemovedFromQuorum)
	removedFromQuorums := make([]*OperatorQuorum, 0)
	for {
		variables["first"] = graphql.Int(maxEntriesPerQuery)
		variables["skip"] = graphql.Int(skip)
		err := a.operatorStateGql.Query(ctx, &result, variables)
		if err != nil {
			return nil, err
		}
		if len(result.OperatorRemovedFromQuorum) == 0 {
			break
		}
		removedFromQuorums = append(removedFromQuorums, result.OperatorRemovedFromQuorum...)
		skip += maxEntriesPerQuery
	}
	return removedFromQuorums, nil
}
