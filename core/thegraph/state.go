package thegraph

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/shurcooL/graphql"
)

const (
	defaultInterval      = time.Second
	maxInterval          = 5 * time.Minute
	maxEntriesPerQuery   = 1000
	startRetriesInterval = time.Second * 5
	startMaxRetries      = 6
)

type (
	IndexedChainState interface {
		GetIndexedOperatorState(ctx context.Context, blockNumber uint, quorums []core.QuorumID) (*core.IndexedOperatorState, error)
		GetIndexedOperatorInfoByOperatorId(ctx context.Context, operatorId core.OperatorID, blockNumber uint32) (*core.IndexedOperatorInfo, error)
	}

	AggregatePubkeyKeyGql struct {
		Apk_X graphql.String `graphql:"apk_X"`
		Apk_Y graphql.String `graphql:"apk_Y"`
	}

	SocketUpdates struct {
		Socket graphql.String
	}

	IndexedOperatorInfoGql struct {
		Id         graphql.String
		PubkeyG1_X graphql.String   `graphql:"pubkeyG1_X"`
		PubkeyG1_Y graphql.String   `graphql:"pubkeyG1_Y"`
		PubkeyG2_X []graphql.String `graphql:"pubkeyG2_X"`
		PubkeyG2_Y []graphql.String `graphql:"pubkeyG2_Y"`
		// Socket is the socket address of the operator, in the form "host:port"
		SocketUpdates []SocketUpdates `graphql:"socketUpdates(first: 1, orderBy: blockNumber, orderDirection: desc)"`
	}

	// The indexed operator state consists of both mutable properties (socket) and  immutable properties
	// (everyhing else: pubkeyG1, pubkeyG2, id). For the socket, we always want the latest value, irrespective
	// of the reference block number. For the immutable properties, we can also use the value from the latest block
	// since value cannot change. Thus, we always pull the state from the latest block indexed by the subgraph.
	//
	// Note that the deregistrationBlockNumber will only be set if the operator has deregistered from all quorums. By using
	// the latest block, we allow the false-positive case where an operator was deregistered from all quorums at the reference
	// block, but then re-registered afterward. Note that this can over-fetch operators but never under-fetch.  We filter out
	// any extra operators in GetIndexedOperatorState.
	QueryOperatorsGql struct {
		Operators []IndexedOperatorInfoGql `graphql:"operators(first: $first, skip: $skip, orderBy: id, orderDirection: desc, where: {deregistrationBlockNumber_gt: $blockNumber})"`
	}

	QueryOperatorByIdGql struct {
		Operator IndexedOperatorInfoGql `graphql:"operator(id: $id)"`
	}

	QueryQuorumAPKGql struct {
		QuorumAPK []AggregatePubkeyKeyGql `graphql:"quorumApks(first: $first,orderDirection:$orderDirection,orderBy:$orderBy,where: {quorumNumber: $quorumNumber,blockNumber_lte: $blockNumber})"`
	}

	queryFirstOperatorGql struct {
		Operators []IndexedOperatorInfoGql `graphql:"operators(first: $first)"`
	}

	GraphQLQuerier interface {
		Query(ctx context.Context, q any, variables map[string]any) error
	}

	indexedChainState struct {
		core.ChainState
		querier GraphQLQuerier

		logger logging.Logger
	}
)

var _ IndexedChainState = (*indexedChainState)(nil)

func NewIndexedChainState(cs core.ChainState, querier GraphQLQuerier, logger logging.Logger) *indexedChainState {
	return &indexedChainState{
		ChainState: cs,
		querier:    querier,
		logger:     logger.With("component", "IndexedChainState"),
	}
}

func (ics *indexedChainState) Start(ctx context.Context) error {
	retries := float64(startMaxRetries)
	for {
		err := ics.querier.Query(ctx, &queryFirstOperatorGql{}, map[string]any{
			"first": graphql.Int(1),
		})
		if err == nil {
			return nil
		}
		ics.logger.Error("Error connecting to subgraph", "err", err)
		if retries <= 0 {
			return errors.New("subgraph timeout")
		}
		retrySec := math.Pow(2, retries)
		time.Sleep(time.Duration(retrySec) * startRetriesInterval)
		retries--
	}
}

func (ics *indexedChainState) GetIndexedOperatorState(ctx context.Context, blockNumber uint, quorums []core.QuorumID) (*core.IndexedOperatorState, error) {
	operatorState, err := ics.ChainState.GetOperatorState(ctx, blockNumber, quorums)
	if err != nil {
		return nil, err
	}

	aggregatePublicKeys := ics.getQuorumAPKs(ctx, quorums, uint32(blockNumber))
	aggKeys := make(map[uint8]*core.G1Point)
	for _, apk := range aggregatePublicKeys {
		if apk.Err != nil {
			ics.logger.Warn("Error getting aggregate public key", "err", apk.Err)
			continue
		}
		if apk.Err == nil && apk.AggregatePubk != nil {
			aggKeys[apk.QuorumNumber] = apk.AggregatePubk
		}
	}
	if len(aggKeys) == 0 {
		return nil, fmt.Errorf("no aggregate public keys found for any of the specified quorums at block number %d", blockNumber)
	}

	indexedOperators, err := ics.getRegisteredIndexedOperatorInfo(ctx, uint32(blockNumber))
	if err != nil {
		return nil, err
	}

	// Detect missing operators
	operatorSeen := make(map[core.OperatorID]struct{})
	for _, quorumOperators := range operatorState.Operators {
		for operatorID := range quorumOperators {
			if indexedOperators[operatorID] == nil {
				return nil, fmt.Errorf("operator %s not found in indexed state", operatorID.Hex())
			}
			operatorSeen[operatorID] = struct{}{}
		}
	}

	// Filter out the operators who are not part of any quorum. This can happen if the operator registers or re-registers
	// after the reference block number.
	for operatorID := range indexedOperators {
		if _, ok := operatorSeen[operatorID]; !ok {
			delete(indexedOperators, operatorID)
		}
	}

	state := &core.IndexedOperatorState{
		OperatorState:    operatorState,
		IndexedOperators: indexedOperators,
		AggKeys:          aggKeys,
	}
	return state, nil
}

// GetIndexedOperatorInfoByOperatorId returns the IndexedOperatorInfo for the operator with the given operatorId at the given block number
func (ics *indexedChainState) GetIndexedOperatorInfoByOperatorId(ctx context.Context, operatorId core.OperatorID, blockNumber uint32) (*core.IndexedOperatorInfo, error) {
	var (
		query     QueryOperatorByIdGql
		variables = map[string]any{
			"id": graphql.String(fmt.Sprintf("0x%s", operatorId.Hex())),
		}
	)
	err := ics.querier.Query(context.Background(), &query, variables)
	if err != nil {
		ics.logger.Error("Error requesting info for operator", "err", err, "operatorId", operatorId.Hex(), "blockNumber", blockNumber)
		return nil, err
	}

	return convertIndexedOperatorInfoGqlToIndexedOperatorInfo(&query.Operator)
}

type quorumAPK struct {
	QuorumNumber  uint8
	AggregatePubk *core.G1Point
	Err           error
}

// GetQuorumAPKs returns the Aggregate Public Keys for the given quorums at the given block number
func (ics *indexedChainState) getQuorumAPKs(ctx context.Context, quorumIDs []core.QuorumID, blockNumber uint32) map[uint8]*quorumAPK {
	quorumAPKs := make(map[uint8]*quorumAPK)
	for i := range quorumIDs {
		id := quorumIDs[i]
		apk, err := ics.getQuorumAPK(ctx, id, blockNumber)
		if err != nil {
			quorumAPKs[id] = &quorumAPK{
				QuorumNumber:  uint8(id),
				AggregatePubk: nil,
				Err:           err,
			}
			continue
		}
		if apk == nil {
			quorumAPKs[id] = &quorumAPK{
				QuorumNumber:  uint8(id),
				AggregatePubk: nil,
				Err:           fmt.Errorf("quorum APK not found for quorum %d", id),
			}
			continue
		}
		quorumAPKs[id] = &quorumAPK{
			QuorumNumber:  uint8(id),
			AggregatePubk: apk,
			Err:           nil,
		}
	}
	return quorumAPKs
}

// GetQuorumAPK returns the Aggregate Public Key for the given quorum at the given block number
func (ics *indexedChainState) getQuorumAPK(ctx context.Context, quorumID core.QuorumID, blockNumber uint32) (*core.G1Point, error) {
	var (
		query     QueryQuorumAPKGql
		variables = map[string]any{
			"first":          graphql.Int(1),
			"orderDirection": graphql.String("desc"),
			"orderBy":        graphql.String("blockNumber"),
			"blockNumber":    graphql.Int(blockNumber),
			"quorumNumber":   graphql.Int(quorumID),
		}
	)
	err := ics.querier.Query(ctx, &query, variables)
	if err != nil {
		ics.logger.Error("Error requesting for apk", "err", err)
		return nil, err
	}

	if len(query.QuorumAPK) == 0 {
		ics.logger.Errorf("no quorum APK found for quorum %d, block number %d", quorumID, blockNumber)
		return nil, errors.New("no quorum APK found")
	}

	quorumAPKPoint := new(bn254.G1Affine)
	_, err = quorumAPKPoint.X.SetString(string(query.QuorumAPK[0].Apk_X))
	if err != nil {
		return nil, err
	}
	_, err = quorumAPKPoint.Y.SetString(string(query.QuorumAPK[0].Apk_Y))
	if err != nil {
		return nil, err
	}
	return &core.G1Point{G1Affine: quorumAPKPoint}, nil
}

// GetRegisteredIndexedOperatorInfo returns the IndexedOperatorInfo for all registered operators at the given block number keyed by operatorId
func (ics *indexedChainState) getRegisteredIndexedOperatorInfo(ctx context.Context, blockNumber uint32) (map[core.OperatorID]*core.IndexedOperatorInfo, error) {
	operatorsGql, err := ics.getAllOperatorsRegisteredAtBlockNumberWithPagination(ctx, blockNumber)
	if err != nil {
		return nil, err
	}

	operators := make(map[core.OperatorID]*core.IndexedOperatorInfo, len(operatorsGql))
	for i := range operatorsGql {
		operator := operatorsGql[i]
		operatorIndexedInfo, err := convertIndexedOperatorInfoGqlToIndexedOperatorInfo(&operator)
		if err != nil {
			return nil, err
		}

		// convert graphql.String to [32]byte
		// example: "0x0000000000000000000000000000000000000000000000000000000000000001" -> [32]byte{0x01}
		operatorId, err := core.OperatorIDFromHex(string(operator.Id))
		if err != nil {
			return nil, err
		}
		operators[operatorId] = operatorIndexedInfo
	}
	return operators, nil
}

func (ics *indexedChainState) getAllOperatorsRegisteredAtBlockNumberWithPagination(ctx context.Context, blockNumber uint32) ([]IndexedOperatorInfoGql, error) {
	operators := make([]IndexedOperatorInfoGql, 0)
	for {
		var (
			query     QueryOperatorsGql
			variables = map[string]any{
				"first":       graphql.Int(maxEntriesPerQuery),
				"skip":        graphql.Int(len(operators)), // skip is the number of operators already retrieved
				"blockNumber": graphql.Int(blockNumber),
			}
		)
		err := ics.querier.Query(ctx, &query, variables)
		if err != nil {
			ics.logger.Error("Error requesting for operators", "err", err)
			return nil, err
		}

		if len(query.Operators) == 0 {
			break
		}
		operators = append(operators, query.Operators...)
	}
	return operators, nil
}

func convertIndexedOperatorInfoGqlToIndexedOperatorInfo(operator *IndexedOperatorInfoGql) (*core.IndexedOperatorInfo, error) {

	if len(operator.SocketUpdates) == 0 {
		return nil, errors.New("no socket found for operator")
	}

	pubkeyG1 := new(bn254.G1Affine)
	_, err := pubkeyG1.X.SetString(string(operator.PubkeyG1_X))
	if err != nil {
		return nil, err
	}
	_, err = pubkeyG1.Y.SetString(string(operator.PubkeyG1_Y))
	if err != nil {
		return nil, err
	}

	pubkeyG2 := new(bn254.G2Affine)
	_, err = pubkeyG2.X.A1.SetString(string(operator.PubkeyG2_X[0]))
	if err != nil {
		return nil, err
	}
	_, err = pubkeyG2.X.A0.SetString(string(operator.PubkeyG2_X[1]))
	if err != nil {
		return nil, err
	}
	_, err = pubkeyG2.Y.A1.SetString(string(operator.PubkeyG2_Y[0]))
	if err != nil {
		return nil, err
	}
	_, err = pubkeyG2.Y.A0.SetString(string(operator.PubkeyG2_Y[1]))
	if err != nil {
		return nil, err
	}

	return &core.IndexedOperatorInfo{
		PubkeyG1: &core.G1Point{G1Affine: pubkeyG1},
		PubkeyG2: &core.G2Point{G2Affine: pubkeyG2},
		Socket:   string(operator.SocketUpdates[0].Socket),
	}, nil
}
