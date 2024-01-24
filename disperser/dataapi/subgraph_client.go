package dataapi

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/disperser/dataapi/subgraph"
	"github.com/gammazero/workerpool"
)

const (
	_14Days           = 14 * 24 * time.Hour
	maxWorkerPoolSize = 10
)

type (
	SubgraphClient interface {
		QueryBatchesWithLimit(ctx context.Context, limit, skip int) ([]*Batch, error)
		QueryOperatorsWithLimit(ctx context.Context, limit int) ([]*Operator, error)
		QueryBatchNonSigningOperatorIdsInInterval(ctx context.Context, intervalSeconds int64) (map[string]int, error)
		QueryIndexedDeregisteredOperatorsInTheLast14Days(ctx context.Context) (*IndexedDeregisteredOperatorState, error)
		QueryNumBatchesByOperatorsInThePastBlockTimestamp(ctx context.Context, blockTimestamp uint64, nonsigers map[string]int) (map[string]int, error)
	}
	Batch struct {
		Id              []byte
		BatchId         uint64
		BatchHeaderHash []byte
		BlockTimestamp  uint64
		BlockNumber     uint64
		TxHash          []byte
		GasFees         *GasFees
	}
	GasFees struct {
		Id       []byte
		GasUsed  uint64
		GasPrice uint64
		TxFee    uint64
	}
	Operator struct {
		Id              []byte
		OperatorId      []byte
		Operator        []byte
		BlockTimestamp  uint64
		BlockNumber     uint64
		TransactionHash []byte
	}
	DeregisteredOperatorInfo struct {
		*core.IndexedOperatorInfo
		// BlockNumber is the block number at which the operator was deregistered.
		BlockNumber uint
	}
	IndexedDeregisteredOperatorState struct {
		Operators map[core.OperatorID]*DeregisteredOperatorInfo
	}
	OperatorEvents struct {
		OperatorId string
		Events     []uint64
	}
	NonSigner struct {
		OperatorId string
		Count      int
	}
	subgraphClient struct {
		api    subgraph.Api
		logger common.Logger
	}
)

var _ SubgraphClient = (*subgraphClient)(nil)

func NewSubgraphClient(api subgraph.Api, logger common.Logger) *subgraphClient {
	return &subgraphClient{api: api, logger: logger}
}

func (sc *subgraphClient) QueryBatchesWithLimit(ctx context.Context, limit, skip int) ([]*Batch, error) {
	subgraphBatches, err := sc.api.QueryBatches(ctx, true, "blockTimestamp", limit, skip)
	if err != nil {
		return nil, err
	}
	batches, err := convertBatches(subgraphBatches)
	if err != nil {
		return nil, err
	}
	return batches, nil
}

func (sc *subgraphClient) QueryOperatorsWithLimit(ctx context.Context, limit int) ([]*Operator, error) {
	operatorsGql, err := sc.api.QueryOperators(ctx, limit)
	if err != nil {
		return nil, err
	}
	operators := make([]*Operator, len(operatorsGql))
	for i, operatorGql := range operatorsGql {
		operator, err := convertOperator(operatorGql)
		if err != nil {
			return nil, err
		}
		operators[i] = operator
	}
	return operators, nil
}

func (sc *subgraphClient) QueryBatchNonSigningOperatorIdsInInterval(ctx context.Context, intervalSeconds int64) (map[string]int, error) {
	batchNonSigningOperatorIdsGql, err := sc.api.QueryBatchNonSigningOperatorIdsInInterval(ctx, intervalSeconds)
	if err != nil {
		return nil, err
	}
	batchNonSigningOperatorIds := make(map[string]int, len(batchNonSigningOperatorIdsGql))
	for _, batchNonSigningOperatorIdsGql := range batchNonSigningOperatorIdsGql {
		for _, nonSigner := range batchNonSigningOperatorIdsGql.NonSigning.NonSigners {
			batchNonSigningOperatorIds[string(nonSigner.OperatorId)]++
		}
	}
	return batchNonSigningOperatorIds, nil
}

func (sc *subgraphClient) QueryNumBatchesByOperatorsInThePastBlockTimestamp(ctx context.Context, blockTimestamp uint64, nonSigners map[string]int) (map[string]int, error) {
	var (
		registeredOperators   []*subgraph.Operator
		deregisteredOperators []*subgraph.Operator
		err                   error
		pool                  = workerpool.New(maxWorkerPoolSize)
	)

	pool.Submit(func() {
		operators, errQ := sc.api.QueryRegisteredOperatorsGreaterThanBlockTimestamp(ctx, blockTimestamp)
		if errQ != nil {
			err = errQ
		}
		registeredOperators = operators
	})

	pool.Submit(func() {
		operators, errQ := sc.api.QueryDeregisteredOperatorsGreaterThanBlockTimestamp(ctx, blockTimestamp)
		if errQ != nil {
			err = errQ
		}
		deregisteredOperators = operators
	})
	pool.StopWait()

	if err != nil {
		return nil, err
	}

	intervalEvents, err := sc.getOperatorsWithRegisteredDeregisteredIntervalEvents(ctx, registeredOperators, deregisteredOperators, blockTimestamp, nonSigners)
	if err != nil {
		return nil, err
	}

	var (
		mu                   sync.Mutex
		numBatchesByOperator = make(map[string]int, 0)
		intervalEventsPool   = workerpool.New(maxWorkerPoolSize)
		currentTs            = uint64(time.Now().Unix())
	)
	for _, ie := range intervalEvents {
		interval := ie
		intervalEventsPool.Submit(func() {
			end := interval.Events[1]
			if end == 0 {
				end = currentTs
			}
			batches, err := sc.api.QueryBatchesByBlockTimestampRange(ctx, interval.Events[0], end)
			if err != nil {
				sc.logger.Error("failed to query batches by block timestamp range", "start", interval.Events[0], "end", end, "err", err)
				return
			}
			if len(batches) > 0 {
				mu.Lock()
				numBatchesByOperator[interval.OperatorId] += len(batches)
				mu.Unlock()
			}
		})
	}
	intervalEventsPool.StopWait()
	return numBatchesByOperator, nil
}

func (sc *subgraphClient) QueryIndexedDeregisteredOperatorsInTheLast14Days(ctx context.Context) (*IndexedDeregisteredOperatorState, error) {
	last14Days := uint64(time.Now().Add(-_14Days).Unix())
	deregisteredOperators, err := sc.api.QueryDeregisteredOperatorsGreaterThanBlockTimestamp(ctx, last14Days)
	if err != nil {
		return nil, err
	}

	operators := make(map[core.OperatorID]*DeregisteredOperatorInfo, len(deregisteredOperators))
	for i := range deregisteredOperators {
		deregisteredOperator := deregisteredOperators[i]
		operator, err := convertOperator(deregisteredOperator)
		if err != nil {
			return nil, err
		}

		var operatorId [32]byte
		copy(operatorId[:], operator.OperatorId)

		operatorInfo, err := sc.api.QueryOperatorInfoByOperatorIdAtBlockNumber(ctx, operatorId, uint32(operator.BlockNumber))
		if err != nil {
			return nil, err
		}
		indexedOperatorInfo, err := ConvertOperatorInfoGqlToIndexedOperatorInfo(operatorInfo)
		if err != nil {
			return nil, err
		}

		operators[operatorId] = &DeregisteredOperatorInfo{
			IndexedOperatorInfo: indexedOperatorInfo,
			BlockNumber:         uint(operator.BlockNumber),
		}
	}

	return &IndexedDeregisteredOperatorState{
		Operators: operators,
	}, nil
}

func (sc *subgraphClient) getOperatorsWithRegisteredDeregisteredIntervalEvents(
	ctx context.Context,
	registeredOperators []*subgraph.Operator,
	deregisteredOperators []*subgraph.Operator,
	blockTimestamp uint64,
	nonSigners map[string]int,
) ([]OperatorEvents, error) {
	sort.SliceStable(registeredOperators, func(i, j int) bool {
		return registeredOperators[i].BlockTimestamp < registeredOperators[j].BlockTimestamp
	})

	sort.SliceStable(deregisteredOperators, func(i, j int) bool {
		return deregisteredOperators[i].BlockTimestamp < deregisteredOperators[j].BlockTimestamp
	})

	// First position is for registration events and second position is for deregistration events
	operators := make(map[string][][]uint64, 0)
	for operatorId := range nonSigners {
		operators[operatorId] = make([][]uint64, 2)
		operators[operatorId][0] = make([]uint64, 0) // registration events
		operators[operatorId][1] = make([]uint64, 0) // deregistration events
	}
	for i := range registeredOperators {
		operator := registeredOperators[i]
		operatorId := string(operator.OperatorId)

		if _, ok := operators[operatorId]; !ok {
			operators[operatorId] = make([][]uint64, 2)
		}
		timestamp, err := strconv.ParseUint(string(operator.BlockTimestamp), 10, 64)
		if err != nil {
			return nil, err
		}
		operators[operatorId][0] = append(operators[operatorId][0], timestamp)
	}

	for i := range deregisteredOperators {
		operator := deregisteredOperators[i]
		operatorId := string(operator.OperatorId)

		timestamp, err := strconv.ParseUint(string(operator.BlockTimestamp), 10, 64)
		if err != nil || timestamp == 0 {
			return nil, err
		}
		if _, ok := operators[operatorId]; ok {
			operators[operatorId][1] = append(operators[operatorId][1], timestamp)
		}
	}

	currentTs := uint64(time.Now().Unix())
	events := make([]OperatorEvents, 0)

	// For the time window [blockTimestamp, now], compute the sub intervals during
	// which the operator is live in EigenDA network for validating batches.
	for operatorId := range nonSigners {
		reg := operators[operatorId][0]
		dereg := operators[operatorId][1]

		// In EigenDA, the registration and deregistration events on timeline for an
		// operator will be like reg-dereg-reg-dereg...
		//
		// The reason is that registering an operator that's already registered will fail
		// and deregistering an operator that's already deregistered will also fail. So
		// the registeration and deregistration will alternate.
		if len(reg)-len(dereg) > 1 || len(dereg)-len(reg) > 1 {
			return nil, fmt.Errorf("The number of registration and deregistration events cannot differ by more than one, num registration events: %d, num deregistration events: %d, operatorId: %s", len(reg), len(dereg), operatorId)
		}

		if len(reg) == 0 && len(dereg) == 0 {
			// The operator has no registration/deregistration events: it's live
			// for the entire time window.
			events = append(events, OperatorEvents{
				OperatorId: operatorId,
				Events:     []uint64{blockTimestamp, currentTs},
			})
		} else if len(reg) == 0 {
			// The operator has only deregistration event: it's live from the beginning
			// of the time window until the deregistration.
			events = append(events, OperatorEvents{
				OperatorId: operatorId,
				Events:     []uint64{blockTimestamp, dereg[0]},
			})
		} else if len(dereg) == 0 {
			// The operator has only registration event: it's live from registration to
			// the end of the time window.
			events = append(events, OperatorEvents{
				OperatorId: operatorId,
				Events:     []uint64{reg[0], currentTs},
			})
		} else {
			// The operator has both registration and deregistration events in the time
			// window.
			if reg[0] < dereg[0] {
				// The first event in the time window is registration. This means at
				// the beginning (i.e. blockTimestamp) it's not live.
				for i := 0; i < len(reg); i++ {
					if i < len(dereg) {
						events = append(events, OperatorEvents{
							OperatorId: operatorId,
							Events:     []uint64{reg[i], dereg[i]},
						})
					} else {
						events = append(events, OperatorEvents{
							OperatorId: operatorId,
							Events:     []uint64{reg[i], currentTs},
						})
					}
				}
			} else {
				// The first event in the time window is deregistration. This means at
				// the beginning (i.e. blockTimestamp) it's live already.
				events = append(events, OperatorEvents{
					OperatorId: operatorId,
					Events:     []uint64{blockTimestamp, dereg[0]},
				})
				for i := 0; i < len(reg); i++ {
					if i+1 < len(dereg) {
						events = append(events, OperatorEvents{
							OperatorId: operatorId,
							Events:     []uint64{reg[i], dereg[i+1]},
						})
					} else {
						events = append(events, OperatorEvents{
							OperatorId: operatorId,
							Events:     []uint64{reg[i], currentTs},
						})
					}
				}
			}
		}
	}

	// Validate the registration and deregistration events are in timeline order.
	for i := 0; i < len(events); i++ {
		if events[i].Events[0] > events[i].Events[1] {
			return nil, fmt.Errorf("Registration timestamp should not be greater than deregistration or current timestamp for operatorId %s, registration timestamp: %d, deregistration or current timestamp: %d", events[i].OperatorId, events[i].Events[0], events[i].Events[1])
		}
		if i > 0 && events[i-1].OperatorId == events[i].OperatorId && events[i-1].Events[1] > events[i].Events[0] {
			return nil, fmt.Errorf("Registration should not happen when the operator is already registered, but found two consecutive registrations at timestamp %d and %d for operatorId %s", events[i-1].Events[0], events[i].Events[0], events[i].OperatorId)
		}
	}

	return events, nil
}

func convertBatches(subgraphBatches []*subgraph.Batches) ([]*Batch, error) {
	batches := make([]*Batch, len(subgraphBatches))
	for i, batch := range subgraphBatches {
		batchId, err := strconv.ParseUint(string(batch.BatchId), 10, 64)
		if err != nil {
			return nil, err
		}
		timestamp, err := strconv.ParseUint(string(batch.BlockTimestamp), 10, 64)
		if err != nil {
			return nil, err
		}
		blockNum, err := strconv.ParseUint(string(batch.BlockNumber), 10, 64)
		if err != nil {
			return nil, err
		}
		gasFees, err := convertGasFees(batch.GasFees)
		if err != nil {
			return nil, err
		}

		batches[i] = &Batch{
			Id:              []byte(batch.Id),
			BatchId:         batchId,
			BatchHeaderHash: []byte(batch.BatchHeaderHash),
			BlockTimestamp:  timestamp,
			BlockNumber:     blockNum,
			TxHash:          []byte(batch.TxHash),
			GasFees:         gasFees,
		}
	}
	return batches, nil
}

func convertGasFees(gasFees subgraph.GasFees) (*GasFees, error) {
	gasUsed, err := strconv.ParseUint(string(gasFees.GasUsed), 10, 64)
	if err != nil {
		return nil, err
	}
	gasPrice, err := strconv.ParseUint(string(gasFees.GasPrice), 10, 64)
	if err != nil {
		return nil, err
	}
	txFee, err := strconv.ParseUint(string(gasFees.TxFee), 10, 64)
	if err != nil {
		return nil, err
	}
	return &GasFees{
		Id:       []byte(gasFees.Id),
		GasUsed:  gasUsed,
		GasPrice: gasPrice,
		TxFee:    txFee,
	}, nil
}

func convertOperator(operator *subgraph.Operator) (*Operator, error) {
	timestamp, err := strconv.ParseUint(string(operator.BlockTimestamp), 10, 64)
	if err != nil {
		return nil, err
	}
	blockNum, err := strconv.ParseUint(string(operator.BlockNumber), 10, 64)
	if err != nil {
		return nil, err
	}
	return &Operator{
		Id:              []byte(operator.Id),
		OperatorId:      []byte(operator.OperatorId),
		Operator:        []byte(operator.Operator),
		BlockTimestamp:  timestamp,
		BlockNumber:     blockNum,
		TransactionHash: []byte(operator.TransactionHash),
	}, nil
}
