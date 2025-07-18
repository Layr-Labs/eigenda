package dataapi

import (
	"context"
	"encoding/hex"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/disperser/dataapi/subgraph"
	"github.com/Layr-Labs/eigensdk-go/logging"

	"github.com/gammazero/workerpool"
)

const (
	maxWorkerPoolSize = 10
)

// Define the type for the enum.
type OperatorState int

const (
	Deregistered OperatorState = iota // iota starts at 0
	Registered
)

type (
	SubgraphClient interface {
		QueryBatchesWithLimit(ctx context.Context, limit, skip int) ([]*Batch, error)
		QueryOperatorsWithLimit(ctx context.Context, limit int) ([]*Operator, error)
		QueryBatchNonSigningOperatorIdsInInterval(ctx context.Context, intervalSeconds int64) (map[string]int, error)
		QueryBatchNonSigningInfoInInterval(ctx context.Context, startTime, endTime int64) ([]*BatchNonSigningInfo, error)
		QueryOperatorQuorumEvent(ctx context.Context, startBlock, endBlock uint32) (*OperatorQuorumEvents, error)
		QueryIndexedOperatorsWithStateForTimeWindow(ctx context.Context, days int32, state OperatorState) (*IndexedQueriedOperatorInfo, error)
		QueryOperatorInfoByOperatorId(ctx context.Context, operatorId string) (*core.IndexedOperatorInfo, error)
		QueryOperatorEjectionsForTimeWindow(ctx context.Context, days int32, operatorId string, first uint, skip uint) ([]*QueriedOperatorEjections, error)
		QueryReservations(ctx context.Context, currentTimestamp uint64, limit, skip int) ([]*Reservation, error)
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
		Id              string
		OperatorId      string
		Operator        string
		BlockTimestamp  uint64
		BlockNumber     uint64
		TransactionHash string
	}
	OperatorQuorum struct {
		Operator       string
		QuorumNumbers  []byte
		BlockNumber    uint32
		BlockTimestamp uint64
	}
	OperatorQuorumEvents struct {
		// AddedToQuorum is mapping from operator address to a list of sorted events
		// (ascending by BlockNumber) where the operator was added to quorums.
		AddedToQuorum map[string][]*OperatorQuorum
		// RemovedFromQuorum is mapping from operator address to a list of sorted events
		// (ascending by BlockNumber) where the operator was removed from quorums.
		RemovedFromQuorum map[string][]*OperatorQuorum
	}
	QueriedOperatorInfo struct {
		IndexedOperatorInfo *core.IndexedOperatorInfo
		// BlockNumber is the block number at which the operator was deregistered.
		BlockNumber          uint
		Metadata             *Operator
		OperatorProcessError string
	}
	IndexedQueriedOperatorInfo struct {
		Operators map[core.OperatorID]*QueriedOperatorInfo
	}

	NonSigner struct {
		OperatorId string
		Count      int
	}
	BatchNonSigningInfo struct {
		BlockNumber          uint32
		QuorumNumbers        []uint8
		ReferenceBlockNumber uint32
		// The operatorIds of nonsigners for the batch.
		NonSigners []string
	}
	Reservation struct {
		Account      string
		EndTimestamp int64
	}
	subgraphClient struct {
		api    subgraph.Api
		logger logging.Logger
	}
)

var _ SubgraphClient = (*subgraphClient)(nil)

func NewSubgraphClient(api subgraph.Api, logger logging.Logger) *subgraphClient {
	return &subgraphClient{api: api, logger: logger.With("component", "SubgraphClient")}
}

func (sc *subgraphClient) QueryBatchesWithLimit(ctx context.Context, limit, skip int) ([]*Batch, error) {
	subgraphBatches, err := sc.api.QueryBatches(ctx, true, "blockTimestamp", limit, skip)
	if err != nil {
		return nil, fmt.Errorf("failed to query batches: %w", err)
	}
	batches, err := convertBatches(subgraphBatches)
	if err != nil {
		return nil, fmt.Errorf("failed to convert batches: %w", err)
	}
	return batches, nil
}

func (sc *subgraphClient) QueryOperatorsWithLimit(ctx context.Context, limit int) ([]*Operator, error) {
	operatorsGql, err := sc.api.QueryOperators(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query operators: %w", err)
	}
	operators := make([]*Operator, len(operatorsGql))
	for i, operatorGql := range operatorsGql {
		operator, err := convertOperator(operatorGql)
		if err != nil {
			return nil, fmt.Errorf("failed to convert operator at index %d: %w", i, err)
		}
		operators[i] = operator
	}
	return operators, nil
}

func (sc *subgraphClient) QueryOperatorInfoByOperatorId(ctx context.Context, operatorId string) (*core.IndexedOperatorInfo, error) {
	operatorInfo, err := sc.api.QueryOperatorInfoByOperatorIdAtBlockNumber(ctx, operatorId, 0)
	if err != nil {
		sc.logger.Error(fmt.Sprintf("failed to query operator info for operator %s", operatorId))
		return nil, fmt.Errorf("failed to query operator info for operator %s: %w", operatorId, err)
	}

	indexedOperatorInfo, err := ConvertOperatorInfoGqlToIndexedOperatorInfo(operatorInfo)
	if err != nil {
		errorMessage := fmt.Sprintf("failed to convert operator info gql to indexed operator info for operator %s", operatorId)
		sc.logger.Error(errorMessage)
		return nil, fmt.Errorf("failed to convert operator info for operator %s: %w", operatorId, err)
	}
	return indexedOperatorInfo, nil
}

func (sc *subgraphClient) QueryBatchNonSigningInfoInInterval(ctx context.Context, startTime, endTime int64) ([]*BatchNonSigningInfo, error) {
	batchNonSigningInfoGql, err := sc.api.QueryBatchNonSigningInfo(ctx, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to query batch non-signing info for interval %d-%d: %w", startTime, endTime, err)
	}
	batchNonSigningInfo := make([]*BatchNonSigningInfo, len(batchNonSigningInfoGql))
	for i, infoGql := range batchNonSigningInfoGql {
		info, err := convertNonSigningInfo(infoGql)
		if err != nil {
			return nil, fmt.Errorf("failed to convert non-signing info at index %d: %w", i, err)
		}
		batchNonSigningInfo[i] = info
	}
	return batchNonSigningInfo, nil
}

func (sc *subgraphClient) QueryBatchNonSigningOperatorIdsInInterval(ctx context.Context, intervalSeconds int64) (map[string]int, error) {
	batchNonSigningOperatorIdsGql, err := sc.api.QueryBatchNonSigningOperatorIdsInInterval(ctx, intervalSeconds)
	if err != nil {
		return nil, fmt.Errorf("failed to query batch non-signing operator IDs for interval %d seconds: %w", intervalSeconds, err)
	}
	batchNonSigningOperatorIds := make(map[string]int, len(batchNonSigningOperatorIdsGql))
	for _, batchNonSigningOperatorIdsGql := range batchNonSigningOperatorIdsGql {
		for _, nonSigner := range batchNonSigningOperatorIdsGql.NonSigning.NonSigners {
			batchNonSigningOperatorIds[string(nonSigner.OperatorId)]++
		}
	}
	return batchNonSigningOperatorIds, nil
}

func (sc *subgraphClient) QueryOperatorQuorumEvent(ctx context.Context, startBlock, endBlock uint32) (*OperatorQuorumEvents, error) {
	var (
		operatorAddedQuorum   []*subgraph.OperatorQuorum
		operatorRemovedQuorum []*subgraph.OperatorQuorum
		err                   error
		pool                  = workerpool.New(maxWorkerPoolSize)
	)

	pool.Submit(func() {
		added, errQ := sc.api.QueryOperatorAddedToQuorum(ctx, startBlock, endBlock)
		if errQ != nil {
			err = fmt.Errorf("failed to query operators added to quorum for blocks %d-%d: %w", startBlock, endBlock, errQ)
		}
		operatorAddedQuorum = added
	})

	pool.Submit(func() {
		removed, errQ := sc.api.QueryOperatorRemovedFromQuorum(ctx, startBlock, endBlock)

		if errQ != nil {
			err = fmt.Errorf("failed to query operators removed from quorum for blocks %d-%d: %w", startBlock, endBlock, errQ)
		}
		operatorRemovedQuorum = removed
	})
	pool.StopWait()

	if err != nil {
		return nil, err
	}

	addedQuorum, err := parseOperatorQuorum(operatorAddedQuorum)
	if err != nil {
		return nil, fmt.Errorf("failed to parse added operator quorum events: %w", err)
	}
	removedQuorum, err := parseOperatorQuorum(operatorRemovedQuorum)
	if err != nil {
		return nil, fmt.Errorf("failed to parse removed operator quorum events: %w", err)
	}

	addedQuorumMap := make(map[string][]*OperatorQuorum)
	for _, opq := range addedQuorum {
		if _, ok := addedQuorumMap[opq.Operator]; !ok {
			addedQuorumMap[opq.Operator] = make([]*OperatorQuorum, 0)
		}
		addedQuorumMap[opq.Operator] = append(addedQuorumMap[opq.Operator], opq)
	}

	removedQuorumMap := make(map[string][]*OperatorQuorum)
	for _, opq := range removedQuorum {
		if _, ok := removedQuorumMap[opq.Operator]; !ok {
			removedQuorumMap[opq.Operator] = make([]*OperatorQuorum, 0)
		}
		removedQuorumMap[opq.Operator] = append(removedQuorumMap[opq.Operator], opq)
	}

	return &OperatorQuorumEvents{
		AddedToQuorum:     addedQuorumMap,
		RemovedFromQuorum: removedQuorumMap,
	}, nil
}

func (sc *subgraphClient) QueryIndexedOperatorsWithStateForTimeWindow(ctx context.Context, days int32, state OperatorState) (*IndexedQueriedOperatorInfo, error) {
	// Query all operators in the last N days.
	lastNDayInSeconds := uint64(time.Now().Add(-time.Duration(days) * 24 * time.Hour).Unix())
	var operators map[core.OperatorID]*QueriedOperatorInfo
	switch state {
	case Deregistered:
		// Get OperatorsInfo for DeRegistered Operators
		deregisteredOperators, err := sc.api.QueryDeregisteredOperatorsGreaterThanBlockTimestamp(ctx, lastNDayInSeconds)
		if err != nil {
			return nil, fmt.Errorf("failed to query deregistered operators for %d days: %w", days, err)
		}

		operators = make(map[core.OperatorID]*QueriedOperatorInfo, len(deregisteredOperators))
		getOperatorInfoForQueriedOperators(sc, ctx, operators, deregisteredOperators)
	case Registered:
		// Get OperatorsInfo for Registered Operators
		registeredOperators, err := sc.api.QueryRegisteredOperatorsGreaterThanBlockTimestamp(ctx, lastNDayInSeconds)
		if err != nil {
			return nil, fmt.Errorf("failed to query registered operators for %d days: %w", days, err)
		}

		operators = make(map[core.OperatorID]*QueriedOperatorInfo, len(registeredOperators))
		getOperatorInfoForQueriedOperators(sc, ctx, operators, registeredOperators)

	default:
		return nil, fmt.Errorf("invalid operator state: %d", state)
	}

	return &IndexedQueriedOperatorInfo{
		Operators: operators,
	}, nil
}

func (sc *subgraphClient) QueryOperatorEjectionsForTimeWindow(ctx context.Context, days int32, operatorId string, first uint, skip uint) ([]*QueriedOperatorEjections, error) {
	// Query all operators in the last N days.
	lastNDayInSeconds := uint64(time.Now().Add(-time.Duration(days) * 24 * time.Hour).Unix())

	var err error
	var ejections []*subgraph.OperatorEjection

	if operatorId == "" {
		ejections, err = sc.api.QueryOperatorEjectionsGteBlockTimestamp(ctx, lastNDayInSeconds, first, skip)
		if err != nil {
			return nil, fmt.Errorf("failed to query operator ejections for %d days: %w", days, err)
		}
	} else {
		ejections, err = sc.api.QueryOperatorEjectionsGteBlockTimestampByOperatorId(ctx, lastNDayInSeconds, operatorId, first, skip)
		if err != nil {
			return nil, fmt.Errorf("failed to query operator ejections for operator %s for %d days: %w", operatorId, days, err)
		}
	}

	queriedEjections := make([]*QueriedOperatorEjections, len(ejections))
	for i, ejection := range ejections {
		blockNumber, err := strconv.ParseUint(string(ejection.BlockNumber), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse block number for ejection at index %d: %w", i, err)
		}

		timestamp, err := strconv.ParseInt(string(ejection.BlockTimestamp), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse block timestamp for ejection at index %d: %w", i, err)
		}

		t := time.Unix(timestamp, 0)
		blockTimestamp := t.Format(time.RFC3339)
		queriedEjections[i] = &QueriedOperatorEjections{
			OperatorId:      string(ejection.OperatorId),
			Quorum:          uint8(ejection.QuorumNumber),
			BlockNumber:     blockNumber,
			BlockTimestamp:  blockTimestamp,
			TransactionHash: string(ejection.TransactionHash),
		}
	}

	return queriedEjections, nil
}

func (sc *subgraphClient) QueryIndexedDeregisteredOperatorsForTimeWindow(ctx context.Context, days int32) (*IndexedQueriedOperatorInfo, error) {
	// Query all deregistered operators in the last N days.
	lastNDayInSeconds := uint64(time.Now().Add(-time.Duration(days) * 24 * time.Hour).Unix())
	deregisteredOperators, err := sc.api.QueryDeregisteredOperatorsGreaterThanBlockTimestamp(ctx, lastNDayInSeconds)
	if err != nil {
		return nil, fmt.Errorf("failed to query deregistered operators for %d days: %w", days, err)
	}

	operators := make(map[core.OperatorID]*QueriedOperatorInfo, len(deregisteredOperators))
	// Get OpeatroInfo for DeRegistered Operators
	getOperatorInfoForQueriedOperators(sc, ctx, operators, deregisteredOperators)

	return &IndexedQueriedOperatorInfo{
		Operators: operators,
	}, nil
}

func (sc *subgraphClient) QueryIndexedRegisteredOperatorsForTimeWindow(ctx context.Context, days int32) (*IndexedQueriedOperatorInfo, error) {
	// Query all registered operators in the last N days.
	lastNDayInSeconds := uint64(time.Now().Add(-time.Duration(days) * 24 * time.Hour).Unix())
	registeredOperators, err := sc.api.QueryRegisteredOperatorsGreaterThanBlockTimestamp(ctx, lastNDayInSeconds)
	if err != nil {
		return nil, fmt.Errorf("failed to query registered operators for %d days: %w", days, err)
	}

	operators := make(map[core.OperatorID]*QueriedOperatorInfo, len(registeredOperators))

	// Get OpeatroInfo for Registered Operators
	getOperatorInfoForQueriedOperators(sc, ctx, operators, registeredOperators)

	return &IndexedQueriedOperatorInfo{
		Operators: operators,
	}, nil

}

func (sc *subgraphClient) QueryReservations(ctx context.Context, currentTimestamp uint64, limit, skip int) ([]*Reservation, error) {
	reservationsGql, err := sc.api.QueryReservations(ctx, currentTimestamp, limit, skip)
	if err != nil {
		return nil, fmt.Errorf("failed to query reservations: %w", err)
	}

	reservations := make([]*Reservation, len(reservationsGql))
	for i, resGql := range reservationsGql {
		endTimestamp, err := strconv.ParseInt(string(resGql.EndTimestamp), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse end timestamp for reservation at index %d: %w", i, err)
		}
		reservations[i] = &Reservation{
			Account:      string(resGql.Account),
			EndTimestamp: endTimestamp,
		}
	}
	return reservations, nil
}

func getOperatorInfoForQueriedOperators(sc *subgraphClient, ctx context.Context, operators map[core.OperatorID]*QueriedOperatorInfo, queriedOperators []*subgraph.Operator) {

	for i := range queriedOperators {
		queriedOperator := queriedOperators[i]
		operator, err := convertOperator(queriedOperator)
		var operatorId [32]byte

		if err != nil && operator == nil {
			sc.logger.Warn("failed to convert", "err", err, "operator", queriedOperator)
			continue
		}

		// Copy the operator id to a 32 byte array.
		copy(operatorId[:], operator.OperatorId)

		operatorInfo, err := sc.api.QueryOperatorInfoByOperatorIdAtBlockNumber(ctx, operator.OperatorId, uint32(operator.BlockNumber))
		if err != nil {
			operatorIdString := "0x" + hex.EncodeToString(operatorId[:])
			errorMessage := fmt.Sprintf("query operator info by operator id at block number failed: %d for operator %s", uint32(operator.BlockNumber), operatorIdString)
			addOperatorWithErrorDetail(operators, operator, operatorId, errorMessage)
			sc.logger.Warn(errorMessage)
			continue
		}
		indexedOperatorInfo, err := ConvertOperatorInfoGqlToIndexedOperatorInfo(operatorInfo)
		if err != nil {
			operatorIdString := "0x" + hex.EncodeToString(operatorId[:])
			errorMessage := fmt.Sprintf("failed to convert operator info gql to indexed operator info at blocknumber: %d for operator %s", uint32(operator.BlockNumber), operatorIdString)
			addOperatorWithErrorDetail(operators, operator, operatorId, errorMessage)
			sc.logger.Warn(errorMessage)
			continue
		}

		operators[operatorId] = &QueriedOperatorInfo{
			IndexedOperatorInfo:  indexedOperatorInfo,
			BlockNumber:          uint(operator.BlockNumber),
			Metadata:             operator,
			OperatorProcessError: "",
		}
	}
}

func convertBatches(subgraphBatches []*subgraph.Batches) ([]*Batch, error) {
	batches := make([]*Batch, len(subgraphBatches))
	for i, batch := range subgraphBatches {
		batchId, err := strconv.ParseUint(string(batch.BatchId), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse batch ID at index %d: %w", i, err)
		}
		timestamp, err := strconv.ParseUint(string(batch.BlockTimestamp), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse block timestamp at index %d: %w", i, err)
		}
		blockNum, err := strconv.ParseUint(string(batch.BlockNumber), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse block number at index %d: %w", i, err)
		}
		gasFees, err := convertGasFees(batch.GasFees)
		if err != nil {
			return nil, fmt.Errorf("failed to convert gas fees at index %d: %w", i, err)
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
		return nil, fmt.Errorf("failed to parse gas used: %w", err)
	}
	gasPrice, err := strconv.ParseUint(string(gasFees.GasPrice), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse gas price: %w", err)
	}
	txFee, err := strconv.ParseUint(string(gasFees.TxFee), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse transaction fee: %w", err)
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
		return nil, fmt.Errorf("failed to parse operator block timestamp: %w", err)
	}
	blockNum, err := strconv.ParseUint(string(operator.BlockNumber), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse operator block number: %w", err)
	}

	return &Operator{
		Id:              string(operator.Id),
		OperatorId:      string(operator.OperatorId),
		Operator:        string(operator.Operator),
		BlockTimestamp:  timestamp,
		BlockNumber:     blockNum,
		TransactionHash: string(operator.TransactionHash),
	}, nil
}

// This helper function adds an operator with an error message to the operators map.
func addOperatorWithErrorDetail(operators map[core.OperatorID]*QueriedOperatorInfo, operator *Operator, operatorId [32]byte, errorMessage string) {
	operators[operatorId] = &QueriedOperatorInfo{
		IndexedOperatorInfo:  nil,
		BlockNumber:          uint(operator.BlockNumber),
		Metadata:             operator,
		OperatorProcessError: errorMessage,
	}
}

func parseOperatorQuorum(operatorQuorum []*subgraph.OperatorQuorum) ([]*OperatorQuorum, error) {
	parsed := make([]*OperatorQuorum, len(operatorQuorum))
	for i, opq := range operatorQuorum {
		blockNum, err := strconv.ParseUint(string(opq.BlockNumber), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse operator quorum block number at index %d: %w", i, err)
		}
		blockTimestamp, err := strconv.ParseUint(string(opq.BlockTimestamp), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse operator quorum block timestamp at index %d: %w", i, err)
		}
		if len(opq.QuorumNumbers) < 2 || len(opq.QuorumNumbers)%2 != 0 {
			return nil, fmt.Errorf("the QuorumNumbers is expected to start with 0x and have an even length, QuorumNumbers: %s", string(opq.QuorumNumbers))
		}
		// The quorum numbers string starts with "0x", so we should skip it.
		quorumStr := string(opq.QuorumNumbers)[2:]
		quorumNumbers := make([]byte, 0)
		for i := 0; i < len(quorumStr); i += 2 {
			pair := quorumStr[i : i+2]
			quorum, err := strconv.Atoi(pair)
			if err != nil {
				return nil, fmt.Errorf("failed to parse quorum number pair '%s' at index %d: %w", pair, i, err)
			}
			quorumNumbers = append(quorumNumbers, uint8(quorum))
		}
		parsed[i] = &OperatorQuorum{
			Operator:       string(opq.Operator),
			QuorumNumbers:  quorumNumbers,
			BlockNumber:    uint32(blockNum),
			BlockTimestamp: blockTimestamp,
		}
	}
	// Sort the quorum events by ascending order of block number.
	sort.SliceStable(parsed, func(i, j int) bool {
		if parsed[i].BlockNumber == parsed[j].BlockNumber {
			return parsed[i].Operator < parsed[j].Operator
		}
		return parsed[i].BlockNumber < parsed[j].BlockNumber
	})
	return parsed, nil
}

func convertNonSigningInfo(infoGql *subgraph.BatchNonSigningInfo) (*BatchNonSigningInfo, error) {
	quorums := make([]uint8, len(infoGql.BatchHeader.QuorumNumbers))
	for i, q := range infoGql.BatchHeader.QuorumNumbers {
		quorum, err := strconv.ParseUint(string(q), 10, 8)
		if err != nil {
			return nil, fmt.Errorf("failed to parse quorum number at index %d: %w", i, err)
		}
		quorums[i] = uint8(quorum)
	}
	blockNum, err := strconv.ParseUint(string(infoGql.BatchHeader.ReferenceBlockNumber), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse reference block number: %w", err)
	}
	confirmBlockNum, err := strconv.ParseUint(string(infoGql.BlockNumber), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse confirmation block number: %w", err)
	}
	nonSigners := make([]string, len(infoGql.NonSigning.NonSigners))
	for i, nonSigner := range infoGql.NonSigning.NonSigners {
		nonSigners[i] = string(nonSigner.OperatorId)
	}

	return &BatchNonSigningInfo{
		BlockNumber:          uint32(confirmBlockNum),
		QuorumNumbers:        quorums,
		ReferenceBlockNumber: uint32(blockNum),
		NonSigners:           nonSigners,
	}, nil
}
