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

// OperatorList wraps a set of operators with their IDs and addresses.
type OperatorList struct {
	operatorIds []core.OperatorID

	// The addressToId and idToAddress provide 1:1 mapping of operator ID and address
	// for operators provided in the "operatorIds" above.
	addressToId map[string]core.OperatorID
	idToAddress map[core.OperatorID]string
}

func NewOperatorList() *OperatorList {
	return &OperatorList{
		operatorIds: make([]core.OperatorID, 0),
		addressToId: make(map[string]core.OperatorID),
		idToAddress: make(map[core.OperatorID]string),
	}
}

func (o *OperatorList) GetOperatorIds() []core.OperatorID {
	return o.operatorIds
}

func (o *OperatorList) Add(id core.OperatorID, address string) {
	if _, exists := o.idToAddress[id]; exists {
		return
	}
	if _, exists := o.addressToId[address]; exists {
		return
	}

	o.addressToId[address] = id
	o.idToAddress[id] = address
	o.operatorIds = append(o.operatorIds, id)
}

func (o *OperatorList) GetAddress(id string) (string, bool) {
	opID, err := core.OperatorIDFromHex(id)
	if err != nil {
		return "", false
	}
	address, exists := o.idToAddress[opID]
	return address, exists
}

func (o *OperatorList) GetID(address string) (core.OperatorID, bool) {
	id, exists := o.addressToId[address]
	return id, exists
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

func (oh *OperatorHandler) ProbeV2OperatorPorts(ctx context.Context, operatorId string) (*OperatorPortCheckResponse, error) {
	operatorInfo, err := oh.subgraphClient.QueryOperatorInfoByOperatorId(ctx, operatorId)
	if err != nil {
		oh.logger.Warn("failed to fetch operator info", "operatorId", operatorId, "error", err)
		return &OperatorPortCheckResponse{}, err
	}

	operatorSocket := core.OperatorSocket(operatorInfo.Socket)

	retrievalOnline, retrievalStatus := false, "v2 retrieval port closed or unreachable"
	retrievalSocket := operatorSocket.GetV2RetrievalSocket()
	if retrievalSocket == "" {
		retrievalStatus = "v2 retrieval port is not registered"
	} else {
		retrievalPortOpen := checkIsOperatorPortOpen(retrievalSocket, 3, oh.logger)
		if retrievalPortOpen {
			retrievalOnline, retrievalStatus = checkServiceOnline(ctx, "validator.Retrieval", retrievalSocket, 3*time.Second)
		}
	}

	dispersalOnline, dispersalStatus := false, "v2 dispersal port closed or unreachable"
	dispersalSocket := operatorSocket.GetV2DispersalSocket()
	if dispersalSocket == "" {
		dispersalStatus = "v2 dispersal port is not registered"
	} else {
		dispersalPortOpen := checkIsOperatorPortOpen(dispersalSocket, 3, oh.logger)
		if dispersalPortOpen {
			dispersalOnline, dispersalStatus = checkServiceOnline(ctx, "validator.Dispersal", dispersalSocket, 3*time.Second)
		}
	}

	// Create the metadata regardless of online status
	portCheckResponse := &OperatorPortCheckResponse{
		OperatorId:      operatorId,
		DispersalSocket: dispersalSocket,
		DispersalStatus: dispersalStatus,
		DispersalOnline: dispersalOnline,
		RetrievalSocket: retrievalSocket,
		RetrievalOnline: retrievalOnline,
		RetrievalStatus: retrievalStatus,
	}

	// Log the online status
	oh.logger.Info("v2 operator port check response", "response", portCheckResponse)

	// Send the metadata to the results channel
	return portCheckResponse, nil
}

func (oh *OperatorHandler) ProbeV1OperatorPorts(ctx context.Context, operatorId string) (*OperatorPortCheckResponse, error) {
	operatorInfo, err := oh.subgraphClient.QueryOperatorInfoByOperatorId(ctx, operatorId)
	if err != nil {
		oh.logger.Warn("failed to fetch operator info", "operatorId", operatorId, "error", err)
		return &OperatorPortCheckResponse{}, err
	}

	operatorSocket := core.OperatorSocket(operatorInfo.Socket)
	retrievalSocket := operatorSocket.GetV1RetrievalSocket()
	retrievalPortOpen := checkIsOperatorPortOpen(retrievalSocket, 3, oh.logger)
	retrievalOnline, retrievalStatus := false, "v1 retrieval port closed or unreachable"
	if retrievalPortOpen {
		retrievalOnline, retrievalStatus = checkServiceOnline(ctx, "node.Retrieval", retrievalSocket, 3*time.Second)
	}

	dispersalSocket := operatorSocket.GetV1DispersalSocket()
	dispersalPortOpen := checkIsOperatorPortOpen(dispersalSocket, 3, oh.logger)
	dispersalOnline, dispersalStatus := false, "v1 dispersal port closed or unreachable"
	if dispersalPortOpen {
		dispersalOnline, dispersalStatus = checkServiceOnline(ctx, "node.Dispersal", dispersalSocket, 3*time.Second)
	}

	// Create the metadata regardless of online status
	portCheckResponse := &OperatorPortCheckResponse{
		OperatorId:      operatorId,
		DispersalSocket: dispersalSocket,
		DispersalStatus: dispersalStatus,
		DispersalOnline: dispersalOnline,
		RetrievalSocket: retrievalSocket,
		RetrievalOnline: retrievalOnline,
		RetrievalStatus: retrievalStatus,
	}

	// Log the online status
	oh.logger.Info("v1 operator port check response", "response", portCheckResponse)

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

func (oh *OperatorHandler) GetOperatorsStakeAtBlock(ctx context.Context, operatorId string, currentBlock uint32) (*OperatorsStakeResponse, error) {
	state, err := oh.chainState.GetOperatorState(ctx, uint(currentBlock), []core.QuorumID{0, 1, 2})
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
	return &OperatorsStakeResponse{
		StakeRankedOperators: stakeRanked,
	}, nil
}

func (oh *OperatorHandler) GetOperatorsStake(ctx context.Context, operatorId string) (*OperatorsStakeResponse, error) {
	currentBlock, err := oh.indexedChainState.GetCurrentBlockNumber()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch current block number: %w", err)
	}
	return oh.GetOperatorsStakeAtBlock(ctx, operatorId, uint32(currentBlock))
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

// CreateOperatorQuorumIntervals creates OperatorQuorumIntervals that are within the
// the block interval [startBlock, endBlock] for operators specified in OperatorList.
//
// Note: the returned result OperatorQuorumIntervals[op][q] means a sequence of increasing
// and non-overlapping block intervals during which the operator "op" is registered in
// quorum "q".
func (oh *OperatorHandler) CreateOperatorQuorumIntervals(
	ctx context.Context,
	operatorList *OperatorList,
	operatorQuorumEvents *OperatorQuorumEvents,
	startBlock, endBlock uint32,
) (OperatorQuorumIntervals, []uint8, error) {
	// Get operators' quorums at startBlock.
	quorumSeen := make(map[uint8]struct{}, 0)

	bitmaps, err := oh.chainReader.GetQuorumBitmapForOperatorsAtBlockNumber(ctx, operatorList.operatorIds, startBlock)
	if err != nil {
		return nil, nil, err
	}
	operatorInitialQuorum := make(map[string][]uint8)
	for i := range bitmaps {
		opQuorumIDs := eth.BitmapToQuorumIds(bitmaps[i])
		operatorInitialQuorum[operatorList.operatorIds[i].Hex()] = opQuorumIDs
		for _, q := range opQuorumIDs {
			quorumSeen[q] = struct{}{}
		}
	}

	// Get all quorums.
	allQuorums := make([]uint8, 0)
	for q := range quorumSeen {
		allQuorums = append(allQuorums, q)
	}

	// Get quorum change events from [startBlock+1, endBlock] for operators in operator set.
	addedToQuorum, removedFromQuorum, err := oh.getOperatorQuorumEvents(operatorQuorumEvents, operatorList)
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

func (oh *OperatorHandler) getOperatorQuorumEvents(
	operatorQuorumEvents *OperatorQuorumEvents,
	operatorList *OperatorList,
) (map[string][]*OperatorQuorum, map[string][]*OperatorQuorum, error) {
	addedToQuorum := make(map[string][]*OperatorQuorum)
	removedFromQuorum := make(map[string][]*OperatorQuorum)
	// Make quorum events organize by operatorID (instead of address) and drop those who
	// are not in the operator set.
	for op, events := range operatorQuorumEvents.AddedToQuorum {
		if id, ok := operatorList.GetID(op); ok {
			addedToQuorum[id.Hex()] = events
		}
	}
	for op, events := range operatorQuorumEvents.RemovedFromQuorum {
		if id, ok := operatorList.GetID(op); ok {
			removedFromQuorum[id.Hex()] = events
		}
	}
	return addedToQuorum, removedFromQuorum, nil
}
