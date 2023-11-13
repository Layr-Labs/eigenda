package subgraph

import (
	"context"
	"sync"
	"time"

	"github.com/shurcooL/graphql"
)

var (
	once                   sync.Once
	instance               *api
	MAX_ENTITIES_PER_QUERY = 1000
)

type (
	Api interface {
		QueryBatches(ctx context.Context, descending bool, orderByField string, first, skip int) ([]*Batches, error)
		QueryOperators(ctx context.Context, first int) ([]*OperatorRegistered, error)
		QueryBatchNonSigningOperatorIdsInInterval(ctx context.Context, intervalSeconds int64) ([]*BatchNonSigningOperatorIds, error)
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

func (a *api) QueryOperators(ctx context.Context, first int) ([]*OperatorRegistered, error) {
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
		variables["first"] = graphql.Int(MAX_ENTITIES_PER_QUERY)
		variables["skip"] = graphql.Int(skip)

		err := a.uiMonitoringGgl.Query(ctx, &result, variables)
		if err != nil {
			return nil, err
		}

		if len(result.BatchNonSigningOperatorIds) == 0 {
			break
		}
		batchNonSigningOperatorIds = append(batchNonSigningOperatorIds, result.BatchNonSigningOperatorIds...)

		skip += MAX_ENTITIES_PER_QUERY
	}

	result.BatchNonSigningOperatorIds = batchNonSigningOperatorIds
	return result.BatchNonSigningOperatorIds, nil
}
