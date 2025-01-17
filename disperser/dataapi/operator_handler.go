package dataapi

import (
	"context"
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/disperser/common/semver"
	"github.com/Layr-Labs/eigenda/operators"
	"github.com/Layr-Labs/eigensdk-go/logging"
)

// OperatorHandler handles operations to collect and process operators info.
type OperatorHandler struct {
	// For visibility
	logger  logging.Logger
	metrics *Metrics

	// For accessing operator info
	chainReader       core.Reader
	chainState        core.ChainState
	indexedChainState core.IndexedChainState
	subgraphClient    SubgraphClient
}

func NewOperatorHandler(logger logging.Logger, metrics *Metrics, chainReader core.Reader, chainState core.ChainState, indexedChainState core.IndexedChainState, subgraphClient SubgraphClient) *OperatorHandler {
	return &OperatorHandler{
		logger:            logger,
		metrics:           metrics,
		chainReader:       chainReader,
		chainState:        chainState,
		indexedChainState: indexedChainState,
		subgraphClient:    subgraphClient,
	}
}

func (oh *OperatorHandler) ProbeOperatorHosts(ctx context.Context, operatorId string) (*OperatorPortCheckResponse, error) {
	operatorInfo, err := oh.subgraphClient.QueryOperatorInfoByOperatorId(ctx, operatorId)
	if err != nil {
		oh.logger.Warn("failed to fetch operator info", "operatorId", operatorId, "error", err)
		return &OperatorPortCheckResponse{}, err
	}

	operatorSocket := core.OperatorSocket(operatorInfo.Socket)
	retrievalSocket := operatorSocket.GetRetrievalSocket()
	retrievalOnline := checkIsOperatorOnline(retrievalSocket, 3, oh.logger)

	dispersalSocket := operatorSocket.GetDispersalSocket()
	dispersalOnline := checkIsOperatorOnline(dispersalSocket, 3, oh.logger)

	// Create the metadata regardless of online status
	portCheckResponse := &OperatorPortCheckResponse{
		OperatorId:      operatorId,
		DispersalSocket: dispersalSocket,
		RetrievalSocket: retrievalSocket,
		DispersalOnline: dispersalOnline,
		RetrievalOnline: retrievalOnline,
	}

	// Log the online status
	oh.logger.Info("operator port check response", "response", portCheckResponse)

	// Send the metadata to the results channel
	return portCheckResponse, nil
}

func (oh *OperatorHandler) GetOperatorsStake(ctx context.Context, operatorId string) (*OperatorsStakeResponse, error) {
	currentBlock, err := oh.indexedChainState.GetCurrentBlockNumber()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch current block number: %w", err)
	}
	state, err := oh.chainState.GetOperatorState(ctx, currentBlock, []core.QuorumID{0, 1, 2})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch indexed operator state: %w", err)
	}

	tqs, quorumsStake := operators.GetRankedOperators(state)
	oh.metrics.UpdateOperatorsStake(tqs, quorumsStake)

	stakeRanked := make(map[string][]*OperatorStake)
	for q, operators := range quorumsStake {
		quorum := fmt.Sprintf("%d", q)
		stakeRanked[quorum] = make([]*OperatorStake, 0)
		for i, op := range operators {
			if len(operatorId) == 0 || operatorId == op.OperatorId.Hex() {
				stakeRanked[quorum] = append(stakeRanked[quorum], &OperatorStake{
					QuorumId:        quorum,
					OperatorId:      op.OperatorId.Hex(),
					StakePercentage: op.StakeShare / 100.0,
					Rank:            i + 1,
				})
			}
		}
	}
	stakeRanked["total"] = make([]*OperatorStake, 0)
	for i, op := range tqs {
		if len(operatorId) == 0 || operatorId == op.OperatorId.Hex() {
			stakeRanked["total"] = append(stakeRanked["total"], &OperatorStake{
				QuorumId:        "total",
				OperatorId:      op.OperatorId.Hex(),
				StakePercentage: op.StakeShare / 100.0,
				Rank:            i + 1,
			})
		}
	}
	return &OperatorsStakeResponse{
		StakeRankedOperators: stakeRanked,
	}, nil
}

func (s *OperatorHandler) ScanOperatorsHostInfo(ctx context.Context) (*SemverReportResponse, error) {
	currentBlock, err := s.indexedChainState.GetCurrentBlockNumber()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch current block number: %w", err)
	}
	operators, err := s.indexedChainState.GetIndexedOperators(context.Background(), currentBlock)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch indexed operator info: %w", err)
	}

	// check operator socket registration against the indexed state
	for operatorID, operatorInfo := range operators {
		socket, err := s.chainState.GetOperatorSocket(context.Background(), currentBlock, operatorID)
		if err != nil {
			s.logger.Warn("failed to get operator socket", "operatorId", operatorID.Hex(), "error", err)
			continue
		}
		if socket != operatorInfo.Socket {
			s.logger.Warn("operator socket mismatch", "operatorId", operatorID.Hex(), "socket", socket, "operatorInfo", operatorInfo.Socket)
		}
	}

	s.logger.Info("Queried indexed operators", "operators", len(operators), "block", currentBlock)
	operatorState, err := s.chainState.GetOperatorState(context.Background(), currentBlock, []core.QuorumID{0, 1, 2})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch operator state: %w", err)
	}

	nodeInfoWorkers := 20
	nodeInfoTimeout := time.Duration(1 * time.Second)
	useRetrievalClient := false
	semvers := semver.ScanOperators(operators, operatorState, useRetrievalClient, nodeInfoWorkers, nodeInfoTimeout, s.logger)

	// Create HostInfoReportResponse instance
	semverReport := &SemverReportResponse{
		Semver: semvers,
	}

	// Publish semver report metrics
	s.metrics.UpdateSemverCounts(semvers)

	s.logger.Info("Semver scan completed", "semverReport", semverReport)
	return semverReport, nil

}
