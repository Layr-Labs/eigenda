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

type (
	SubgraphClient interface {
		QueryBatchesWithLimit(ctx context.Context, limit, skip int) ([]*Batch, error)
		QueryOperatorsWithLimit(ctx context.Context, limit int) ([]*Operator, error)
		QueryBatchNonSigningOperatorIdsInInterval(ctx context.Context, intervalSeconds int64) (map[string]int, error)
		QueryBatchNonSigningInfoInInterval(ctx context.Context, intervalSeconds int64) ([]*BatchNonSigningInfo, error)
		QueryOperatorQuorumEvent(ctx context.Context, startBlock, endBlock uint32) (*OperatorQuorumEvents, error)
		QueryIndexedDeregisteredOperatorsForTimeWindow(ctx context.Context, days int32) (*IndexedDeregisteredOperatorState, error)
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
	DeregisteredOperatorInfo struct {
		IndexedOperatorInfo *core.IndexedOperatorInfo
		// BlockNumber is the block number at which the operator was deregistered.
		BlockNumber          uint
		Metadata             *Operator
		OperatorProcessError string
	}
	IndexedDeregisteredOperatorState struct {
		Operators map[core.OperatorID]*DeregisteredOperatorInfo
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
	subgraphClient struct {
		api    subgraph.Api
		logger logging.Logger
	}
)

var _ SubgraphClient = (*subgraphClient)(nil)

func NewSubgraphClient(api subgraph.Api, logger logging.Logger) *subgraphClient {
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

func (sc *subgraphClient) QueryBatchNonSigningInfoInInterval(ctx context.Context, intervalSeconds int64) ([]*BatchNonSigningInfo, error) {
	batchNonSigningInfoGql, err := sc.api.QueryBatchNonSigningInfo(ctx, intervalSeconds)
	if err != nil {
		return nil, err
	}
	batchNonSigningInfo := make([]*BatchNonSigningInfo, len(batchNonSigningInfoGql))
	for i, infoGql := range batchNonSigningInfoGql {
		info, err := convertNonSigningInfo(infoGql)
		if err != nil {
			return nil, err
		}
		batchNonSigningInfo[i] = info
	}
	return batchNonSigningInfo, nil
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
			err = errQ
		}
		operatorAddedQuorum = added
	})

	pool.Submit(func() {
		removed, errQ := sc.api.QueryOperatorRemovedFromQuorum(ctx, startBlock, endBlock)

		if errQ != nil {
			err = errQ
		}
		operatorRemovedQuorum = removed
	})
	pool.StopWait()

	if err != nil {
		return nil, err
	}

	addedQuorum, err := parseOperatorQuorum(operatorAddedQuorum)
	if err != nil {
		return nil, err
	}
	removedQuorum, err := parseOperatorQuorum(operatorRemovedQuorum)
	if err != nil {
		return nil, err
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
		var operatorId [32]byte

		if err != nil && operator == nil {
			sc.logger.Warn("failed to convert", "err", err, "operator", deregisteredOperator)
			continue
		}

		// Copy the operator id to a 32 byte array.
		copy(operatorId[:], operator.OperatorId)

		operatorInfo, err := sc.api.QueryOperatorInfoByOperatorIdAtBlockNumber(ctx, operatorId, uint32(operator.BlockNumber))
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

		operators[operatorId] = &DeregisteredOperatorInfo{
			IndexedOperatorInfo:  indexedOperatorInfo,
			BlockNumber:          uint(operator.BlockNumber),
			Metadata:             operator,
			OperatorProcessError: "",
		}
	}

	return &IndexedDeregisteredOperatorState{
		Operators: operators,
	}, nil
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

// This helper function adds an operator with an error message to the operators map.
func addOperatorWithErrorDetail(operators map[core.OperatorID]*DeregisteredOperatorInfo, operator *Operator, operatorId [32]byte, errorMessage string) {
	operators[operatorId] = &DeregisteredOperatorInfo{
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
			return nil, err
		}
		blockTimestamp, err := strconv.ParseUint(string(opq.BlockTimestamp), 10, 64)
		if err != nil {
			return nil, err
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
				return nil, err
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
			return nil, err
		}
		quorums[i] = uint8(quorum)
	}
	blockNum, err := strconv.ParseUint(string(infoGql.BatchHeader.ReferenceBlockNumber), 10, 64)
	if err != nil {
		return nil, err
	}
	confirmBlockNum, err := strconv.ParseUint(string(infoGql.BlockNumber), 10, 64)
	if err != nil {
		return nil, err
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
