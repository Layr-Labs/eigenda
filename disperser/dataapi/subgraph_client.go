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
	maxWorkerPoolSize = 10
)

type (
	SubgraphClient interface {
		QueryBatchesWithLimit(ctx context.Context, limit, skip int) ([]*Batch, error)
		QueryOperatorsWithLimit(ctx context.Context, limit int) ([]*Operator, error)
		QueryBatchNonSigningOperatorIdsInInterval(ctx context.Context, intervalSeconds int64) (map[string]int, error)
		QueryIndexedDeregisteredOperatorsForTimeWindow(ctx context.Context, days int32) (*IndexedDeregisteredOperatorState, error)
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
		IndexedOperatorInfo *core.IndexedOperatorInfo
		// BlockNumber is the block number at which the operator was deregistered.
		BlockNumber uint
		Metadata    *Operator
	}
	IndexedDeregisteredOperatorState struct {
		Operators map[core.OperatorID]*DeregisteredOperatorInfo
	}
	// OperatorInterval describes a time interval where the operator is live in
	// EigenDA.
	OperatorInterval struct {
		OperatorId string

		// The operator is live from start to end.
		start uint64
		// If the operator is still live now in EigenDA netowrk, end is set to 0.
		end uint64
	}
	// OperatorEvents describes all the registration and deregistration events associated
	// with an operator.
	OperatorEvents struct {
		// Timestamps of operator's registration, in ascending order.
		RegistrationEvents []uint64
		// Timestamps of operator's deregistration, in ascending order.
		DeregistrationEvents []uint64
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
	type timeInterval struct {
		start, end uint64
	}
	// Caching the number of batches in a time interval so we don't need to query
	// subgraph repeatedly. In usual case, most operators will have no opt-in/opt-out
	// events in recent time window, so all of them will just query the same time
	// interval.
	numBatchesCache := make(map[timeInterval]int)
	for _, ie := range intervalEvents {
		interval := ie
		intervalEventsPool.Submit(func() {
			end := interval.end
			if end == 0 {
				end = currentTs
			}
			timeRange := timeInterval{start: interval.start, end: end}
			mu.Lock()
			_, ok := numBatchesCache[timeRange]
			mu.Unlock()
			if !ok {
				batches, err := sc.api.QueryBatchesByBlockTimestampRange(ctx, interval.start, end)
				if err != nil {
					sc.logger.Error("failed to query batches by block timestamp range", "start", interval.start, "end", end, "err", err)
					return
				}
				mu.Lock()
				numBatchesCache[timeRange] = len(batches)
				mu.Unlock()
			}
			mu.Lock()
			numBatchesByOperator[interval.OperatorId] += numBatchesCache[timeRange]
			mu.Unlock()
		})
	}
	intervalEventsPool.StopWait()
	return numBatchesByOperator, nil
}

func (sc *subgraphClient) QueryIndexedDeregisteredOperatorsForTimeWindow(ctx context.Context, days int32) (*IndexedDeregisteredOperatorState, error) {
	// Query all deregistered operators in the last N days.
	lastNDayInSeconds := uint64(time.Now().Add(-time.Duration(days) * 24 * time.Hour).Unix())
	deregisteredOperators, err := sc.api.QueryDeregisteredOperatorsGreaterThanBlockTimestamp(ctx, lastNDayInSeconds)
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
			Metadata:            operator,
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
) ([]OperatorInterval, error) {
	sort.SliceStable(registeredOperators, func(i, j int) bool {
		return registeredOperators[i].BlockTimestamp < registeredOperators[j].BlockTimestamp
	})

	sort.SliceStable(deregisteredOperators, func(i, j int) bool {
		return deregisteredOperators[i].BlockTimestamp < deregisteredOperators[j].BlockTimestamp
	})

	operators := make(map[string]OperatorEvents)
	for operatorId := range nonSigners {
		operators[operatorId] = OperatorEvents{
			RegistrationEvents:   []uint64{},
			DeregistrationEvents: []uint64{},
		}
	}
	for i := range registeredOperators {
		reg := registeredOperators[i]
		operatorId := string(reg.OperatorId)

		// If the operator is not a nonsigner, skip it.
		if _, ok := operators[operatorId]; !ok {
			continue
		}
		timestamp, err := strconv.ParseUint(string(reg.BlockTimestamp), 10, 64)
		if err != nil {
			return nil, err
		}
		operator := operators[operatorId]
		operator.RegistrationEvents = append(operator.RegistrationEvents, timestamp)
		operators[operatorId] = operator
	}

	for i := range deregisteredOperators {
		dereg := deregisteredOperators[i]
		operatorId := string(dereg.OperatorId)

		// If the operator is not a nonsigner, skip it.
		if _, ok := operators[operatorId]; !ok {
			continue
		}
		timestamp, err := strconv.ParseUint(string(dereg.BlockTimestamp), 10, 64)
		if err != nil || timestamp == 0 {
			return nil, err
		}
		operator := operators[operatorId]
		operator.DeregistrationEvents = append(operator.DeregistrationEvents, timestamp)
		operators[operatorId] = operator
	}

	events, err := getOperatorInterval(ctx, operators, blockTimestamp, nonSigners)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func getOperatorInterval(
	ctx context.Context,
	operators map[string]OperatorEvents,
	blockTimestamp uint64,
	nonSigners map[string]int,
) ([]OperatorInterval, error) {
	currentTs := uint64(time.Now().Unix())
	intervals := make([]OperatorInterval, 0)

	// For the time window [blockTimestamp, now], compute the sub intervals during
	// which the operator is live in EigenDA network for validating batches.
	for operatorId := range nonSigners {
		reg := operators[operatorId].RegistrationEvents
		dereg := operators[operatorId].DeregistrationEvents

		// In EigenDA, the registration and deregistration events on timeline for an
		// operator will be like reg-dereg-reg-dereg...
		//
		// The reason is that registering an operator that's already registered will fail
		// and deregistering an operator that's already deregistered will also fail. So
		// the registeration and deregistration will alternate.
		if len(reg)-len(dereg) > 1 || len(dereg)-len(reg) > 1 {
			return nil, fmt.Errorf("The number of registration and deregistration events cannot differ by more than one, num registration events: %d, num deregistration events: %d, operatorId: %s", len(reg), len(dereg), operatorId)
		}

		// Note: if an operator registered at block A and then deregistered
		// at block B, the range of liveness will be [A, B), i.e. the operator
		// will not be responsible for signing at block B.

		if len(reg) == 0 && len(dereg) == 0 {
			// The operator has no registration/deregistration events: it's live
			// for the entire time window.
			intervals = append(intervals, OperatorInterval{
				OperatorId: operatorId,
				start:      blockTimestamp,
				end:        currentTs,
			})
		} else if len(reg) == 0 {
			// The operator has only deregistration event: it's live from the beginning
			// of the time window until the deregistration.
			intervals = append(intervals, OperatorInterval{
				OperatorId: operatorId,
				start:      blockTimestamp,
				end:        dereg[0] - 1,
			})
		} else if len(dereg) == 0 {
			// The operator has only registration event: it's live from registration to
			// the end of the time window.
			intervals = append(intervals, OperatorInterval{
				OperatorId: operatorId,
				start:      reg[0],
				end:        currentTs,
			})
		} else {
			// The operator has both registration and deregistration events in the time
			// window.
			if reg[0] < dereg[0] {
				// The first event in the time window is registration. This means at
				// the beginning (i.e. blockTimestamp) it's not live.
				for i := 0; i < len(reg); i++ {
					if i < len(dereg) {
						intervals = append(intervals, OperatorInterval{
							OperatorId: operatorId,
							start:      reg[i],
							end:        dereg[i] - 1,
						})
					} else {
						intervals = append(intervals, OperatorInterval{
							OperatorId: operatorId,
							start:      reg[i],
							end:        currentTs,
						})
					}
				}
			} else {
				// The first event in the time window is deregistration. This means at
				// the beginning (i.e. blockTimestamp) it's live already.
				intervals = append(intervals, OperatorInterval{
					OperatorId: operatorId,
					start:      blockTimestamp,
					end:        dereg[0] - 1,
				})
				for i := 0; i < len(reg); i++ {
					if i+1 < len(dereg) {
						intervals = append(intervals, OperatorInterval{
							OperatorId: operatorId,
							start:      reg[i],
							end:        dereg[i+1] - 1,
						})
					} else {
						intervals = append(intervals, OperatorInterval{
							OperatorId: operatorId,
							start:      reg[i],
							end:        currentTs,
						})
					}
				}
			}
		}
	}

	// Validate the registration and deregistration events are in timeline order.
	for i := 0; i < len(intervals); i++ {
		if intervals[i].start > intervals[i].end {
			return nil, fmt.Errorf("Start timestamp should not be greater than end or current timestamp for operatorId %s, start timestamp: %d, end or current timestamp: %d", intervals[i].OperatorId, intervals[i].start, intervals[i].end)
		}
		if i > 0 && intervals[i-1].OperatorId == intervals[i].OperatorId && intervals[i-1].end > intervals[i].start {
			return nil, fmt.Errorf("the operator live intervals should never overlap, but found two overlapping intervals [%d, %d] and [%d, %d] for operatorId %s", intervals[i-1].start, intervals[i-1].end, intervals[i].start, intervals[i].end, intervals[i].OperatorId)
		}
	}

	return intervals, nil
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
