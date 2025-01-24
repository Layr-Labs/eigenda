package dataapi

import (
	"context"
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/eth"
	"github.com/Layr-Labs/eigenda/disperser/common/semver"
	"github.com/Layr-Labs/eigenda/operators"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection/grpc_reflection_v1"
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
	retrievalPortOpen := checkIsOperatorPortOpen(retrievalSocket, 3, oh.logger)
	retrievalOnline, retrievalStatus := false, "port closed or unreachable"
	if retrievalPortOpen {
		retrievalOnline, retrievalStatus = checkServiceOnline(ctx, "node.Retrieval", retrievalSocket, 3*time.Second)
	}

	v1DispersalSocket := operatorSocket.GetV1DispersalSocket()
	v1DispersalPortOpen := checkIsOperatorPortOpen(v1DispersalSocket, 3, oh.logger)
	v1DispersalOnline, v1DispersalStatus := false, "port closed or unreachable"
	if v1DispersalPortOpen {
		v1DispersalOnline, v1DispersalStatus = checkServiceOnline(ctx, "node.Dispersal", v1DispersalSocket, 3*time.Second)
	}

	v2DispersalOnline, v2DispersalStatus := false, ""
	v2DispersalSocket := operatorSocket.GetV2DispersalSocket()
	if v2DispersalSocket == "" {
		v2DispersalStatus = "v2 dispersal port is not registered"
	} else {
		v2DispersalPortOpen := checkIsOperatorPortOpen(v2DispersalSocket, 3, oh.logger)
		if !v2DispersalPortOpen {
			v2DispersalStatus = "port closed or unreachable"
		} else {
			v2DispersalOnline, v2DispersalStatus = checkServiceOnline(ctx, "node.v2.Dispersal", v2DispersalSocket, 3*time.Second)
		}
	}

	// Create the metadata regardless of online status
	portCheckResponse := &OperatorPortCheckResponse{
		OperatorId:        operatorId,
		DispersalSocket:   v1DispersalSocket,
		DispersalStatus:   v1DispersalStatus,
		DispersalOnline:   v1DispersalOnline,
		V2DispersalSocket: v2DispersalSocket,
		V2DispersalOnline: v2DispersalOnline,
		V2DispersalStatus: v2DispersalStatus,
		RetrievalSocket:   retrievalSocket,
		RetrievalOnline:   retrievalOnline,
		RetrievalStatus:   retrievalStatus,
	}

	// Log the online status
	oh.logger.Info("operator port check response", "response", portCheckResponse)

	// Send the metadata to the results channel
	return portCheckResponse, nil
}

// query operator host info endpoint if available
func checkServiceOnline(ctx context.Context, serviceName string, socket string, timeout time.Duration) (bool, string) {
	conn, err := grpc.NewClient(socket, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return false, err.Error()
	}
	defer conn.Close()
	ctxWithTimeout, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Create a reflection client
	reflectionClient := grpc_reflection_v1.NewServerReflectionClient(conn)

	// Send ListServices request
	stream, err := reflectionClient.ServerReflectionInfo(ctxWithTimeout)
	if err != nil {
		return false, err.Error()
	}

	// Send the ListServices request
	listReq := &grpc_reflection_v1.ServerReflectionRequest{
		MessageRequest: &grpc_reflection_v1.ServerReflectionRequest_ListServices{},
	}
	if err := stream.Send(listReq); err != nil {
		return false, err.Error()
	}

	// Get the response
	r, err := stream.Recv()
	if err != nil {
		return false, err.Error()
	}

	// Check if the service exists
	if list := r.GetListServicesResponse(); list != nil {
		for _, service := range list.GetService() {
			if service.GetName() == serviceName {
				return true, fmt.Sprintf("%s is available", serviceName)
			}
		}
	}
	return false, fmt.Sprintf("grpc available but %s service not found at %s", serviceName, socket)
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

func (oh *OperatorHandler) CreateOperatorQuorumIntervals(ctx context.Context, nonsigners []core.OperatorID, nonsignerAddressToId map[string]core.OperatorID, startBlock, endBlock uint32) (OperatorQuorumIntervals, []uint8, error) {
	// Get operators' initial quorums (at startBlock).
	quorumSeen := make(map[uint8]struct{}, 0)

	bitmaps, err := oh.chainReader.GetQuorumBitmapForOperatorsAtBlockNumber(ctx, nonsigners, startBlock)
	if err != nil {
		return nil, nil, err
	}
	operatorInitialQuorum := make(map[string][]uint8)
	for i := range bitmaps {
		opQuorumIDs := eth.BitmapToQuorumIds(bitmaps[i])
		operatorInitialQuorum[nonsigners[i].Hex()] = opQuorumIDs
		for _, q := range opQuorumIDs {
			quorumSeen[q] = struct{}{}
		}
	}

	// Get all quorums.
	allQuorums := make([]uint8, 0)
	for q := range quorumSeen {
		allQuorums = append(allQuorums, q)
	}

	// Get operators' quorum change events from [startBlock+1, endBlock].
	addedToQuorum, removedFromQuorum, err := oh.getOperatorQuorumEvents(ctx, startBlock, endBlock, nonsignerAddressToId)
	if err != nil {
		return nil, nil, err
	}

	// Create operators' quorum intervals.
	operatorQuorumIntervals, err := CreateOperatorQuorumIntervals(startBlock, endBlock, operatorInitialQuorum, addedToQuorum, removedFromQuorum)
	if err != nil {
		return nil, nil, err
	}

	return operatorQuorumIntervals, allQuorums, nil
}

func (oh *OperatorHandler) getOperatorQuorumEvents(ctx context.Context, startBlock, endBlock uint32, nonsignerAddressToId map[string]core.OperatorID) (map[string][]*OperatorQuorum, map[string][]*OperatorQuorum, error) {
	addedToQuorum := make(map[string][]*OperatorQuorum)
	removedFromQuorum := make(map[string][]*OperatorQuorum)
	if startBlock == endBlock {
		return addedToQuorum, removedFromQuorum, nil
	}
	operatorQuorumEvents, err := oh.subgraphClient.QueryOperatorQuorumEvent(ctx, startBlock+1, endBlock)
	if err != nil {
		return nil, nil, err
	}
	// Make quorum events organize by operatorID (instead of address) and drop those who
	// are not nonsigners.
	for op, events := range operatorQuorumEvents.AddedToQuorum {
		if id, ok := nonsignerAddressToId[op]; ok {
			addedToQuorum[id.Hex()] = events
		}
	}
	for op, events := range operatorQuorumEvents.RemovedFromQuorum {
		if id, ok := nonsignerAddressToId[op]; ok {
			removedFromQuorum[id.Hex()] = events
		}
	}
	return addedToQuorum, removedFromQuorum, nil
}
