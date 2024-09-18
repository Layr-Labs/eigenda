package dataapi

import (
	"context"
	"fmt"
	"net"
	"sort"
	"sync"
	"time"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/disperser/common/semver"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/gammazero/workerpool"
)

type OperatorOnlineStatus struct {
	OperatorInfo         *Operator
	IndexedOperatorInfo  *core.IndexedOperatorInfo
	OperatorProcessError string
}

var (
	// TODO: Poolsize should be configurable
	// Observe performance and tune accordingly
	poolSize                        = 50
	operatorOnlineStatusresultsChan chan *QueriedStateOperatorMetadata
)

// Function to get registered operators for given number of days
// Queries subgraph for deregistered operators
// Process operator online status
// Returns list of Operators with their online status, socket address and block number they deregistered
func (s *server) getDeregisteredOperatorForDays(ctx context.Context, days int32) ([]*QueriedStateOperatorMetadata, error) {
	// Track time taken to get deregistered operators
	startTime := time.Now()

	indexedDeregisteredOperatorState, err := s.subgraphClient.QueryIndexedOperatorsWithStateForTimeWindow(ctx, days, Deregistered)
	if err != nil {
		return nil, err
	}

	// Convert the map to a slice.
	operators := indexedDeregisteredOperatorState.Operators

	operatorOnlineStatusresultsChan = make(chan *QueriedStateOperatorMetadata, len(operators))
	processOperatorOnlineCheck(indexedDeregisteredOperatorState, operatorOnlineStatusresultsChan, s.logger)

	// Collect results of work done
	DeregisteredOperatorMetadata := make([]*QueriedStateOperatorMetadata, 0, len(operators))
	for range operators {
		metadata := <-operatorOnlineStatusresultsChan
		DeregisteredOperatorMetadata = append(DeregisteredOperatorMetadata, metadata)
	}

	// Log the time taken
	s.logger.Info("Time taken to get deregistered operators for days", "duration", time.Since(startTime))
	sort.Slice(DeregisteredOperatorMetadata, func(i, j int) bool {
		return DeregisteredOperatorMetadata[i].BlockNumber < DeregisteredOperatorMetadata[j].BlockNumber
	})

	return DeregisteredOperatorMetadata, nil
}

// Function to get registered operators for given number of days
// Queries subgraph for registered operators
// Process operator online status
// Returns list of Operators with their online status, socket address and block number they registered
func (s *server) getRegisteredOperatorForDays(ctx context.Context, days int32) ([]*QueriedStateOperatorMetadata, error) {
	// Track time taken to get registered operators
	startTime := time.Now()

	indexedRegisteredOperatorState, err := s.subgraphClient.QueryIndexedOperatorsWithStateForTimeWindow(ctx, days, Registered)
	if err != nil {
		return nil, err
	}

	// Convert the map to a slice.
	operators := indexedRegisteredOperatorState.Operators

	operatorOnlineStatusresultsChan = make(chan *QueriedStateOperatorMetadata, len(operators))
	processOperatorOnlineCheck(indexedRegisteredOperatorState, operatorOnlineStatusresultsChan, s.logger)

	// Collect results of work done
	RegisteredOperatorMetadata := make([]*QueriedStateOperatorMetadata, 0, len(operators))
	for range operators {
		metadata := <-operatorOnlineStatusresultsChan
		RegisteredOperatorMetadata = append(RegisteredOperatorMetadata, metadata)
	}

	// Log the time taken
	s.logger.Info("Time taken to get registered operators for days", "duration", time.Since(startTime))
	sort.Slice(RegisteredOperatorMetadata, func(i, j int) bool {
		return RegisteredOperatorMetadata[i].BlockNumber < RegisteredOperatorMetadata[j].BlockNumber
	})

	return RegisteredOperatorMetadata, nil
}

func processOperatorOnlineCheck(queriedOperatorsInfo *IndexedQueriedOperatorInfo, operatorOnlineStatusresultsChan chan<- *QueriedStateOperatorMetadata, logger logging.Logger) {
	operators := queriedOperatorsInfo.Operators
	wp := workerpool.New(poolSize)

	for _, operatorInfo := range operators {
		operatorStatus := OperatorOnlineStatus{
			OperatorInfo:         operatorInfo.Metadata,
			IndexedOperatorInfo:  operatorInfo.IndexedOperatorInfo,
			OperatorProcessError: operatorInfo.OperatorProcessError,
		}

		// Submit each operator status check to the worker pool
		wp.Submit(func() {
			checkIsOnlineAndProcessOperator(operatorStatus, operatorOnlineStatusresultsChan, logger)
		})
	}

	wp.StopWait() // Wait for all submitted tasks to complete and stop the pool
}

func checkIsOnlineAndProcessOperator(operatorStatus OperatorOnlineStatus, operatorOnlineStatusresultsChan chan<- *QueriedStateOperatorMetadata, logger logging.Logger) {
	var isOnline bool
	var socket string
	if operatorStatus.IndexedOperatorInfo != nil {
		socket = core.OperatorSocket(operatorStatus.IndexedOperatorInfo.Socket).GetRetrievalSocket()
		isOnline = checkIsOperatorOnline(socket, 10, logger)
	}

	// Log the online status
	if isOnline {
		logger.Debug("Operator is online", "operatorInfo", operatorStatus.IndexedOperatorInfo, "socket", socket)
	} else {
		logger.Debug("Operator is offline", "operatorInfo", operatorStatus.IndexedOperatorInfo, "socket", socket)
	}

	// Create the metadata regardless of online status
	metadata := &QueriedStateOperatorMetadata{
		OperatorId:           string(operatorStatus.OperatorInfo.OperatorId[:]),
		BlockNumber:          uint(operatorStatus.OperatorInfo.BlockNumber),
		Socket:               socket,
		IsOnline:             isOnline,
		OperatorProcessError: operatorStatus.OperatorProcessError,
	}

	// Send the metadata to the results channel
	operatorOnlineStatusresultsChan <- metadata
}

// Check that the socketString is not private/unspecified
func ValidOperatorIP(address string, logger logging.Logger) bool {
	host, _, err := net.SplitHostPort(address)
	if err != nil {
		logger.Error("Failed to split host port", "address", address, "error", err)
		return false
	}
	ips, err := net.LookupIP(host)
	if err != nil {
		logger.Error("Error resolving operator host IP", "host", host, "error", err)
		return false
	}
	ipAddr := ips[0]
	if ipAddr == nil {
		logger.Error("IP address is nil", "host", host, "ips", ips)
		return false
	}
	isValid := !ipAddr.IsPrivate() && !ipAddr.IsUnspecified()
	logger.Debug("Operator IP validation", "address", address, "host", host, "ips", ips, "ipAddr", ipAddr, "isValid", isValid)

	return isValid
}

func (s *server) probeOperatorPorts(ctx context.Context, operatorId string) (*OperatorPortCheckResponse, error) {
	operatorInfo, err := s.getOperatorInfo(ctx, operatorId)
	if err != nil {
		s.logger.Warn("failed to fetch operator info", "operatorId", operatorId, "error", err)
		return &OperatorPortCheckResponse{}, err
	}

	operatorSocket := core.OperatorSocket(operatorInfo.Socket)
	retrievalSocket := operatorSocket.GetRetrievalSocket()
	retrievalOnline := checkIsOperatorOnline(retrievalSocket, 3, s.logger)

	dispersalSocket := operatorSocket.GetDispersalSocket()
	dispersalOnline := checkIsOperatorOnline(dispersalSocket, 3, s.logger)

	// Create the metadata regardless of online status
	portCheckResponse := &OperatorPortCheckResponse{
		OperatorId:      operatorId,
		DispersalSocket: dispersalSocket,
		RetrievalSocket: retrievalSocket,
		DispersalOnline: dispersalOnline,
		RetrievalOnline: retrievalOnline,
	}

	// Log the online status
	s.logger.Info("operator port check response", "response", portCheckResponse)

	// Send the metadata to the results channel
	return portCheckResponse, nil
}

func (s *server) getOperatorInfo(ctx context.Context, operatorId string) (*core.IndexedOperatorInfo, error) {
	operatorInfo, err := s.subgraphClient.QueryOperatorInfoByOperatorId(ctx, operatorId)
	if err != nil {
		s.logger.Warn("failed to fetch operator info", "operatorId", operatorId, "error", err)
		return nil, fmt.Errorf("operator info not found for operatorId %s", operatorId)
	}
	return operatorInfo, nil
}

func (s *server) scanOperatorsHostInfo(ctx context.Context, logger logging.Logger) (*SemverReportResponse, error) {
	registrations, err := s.subgraphClient.QueryOperatorsWithLimit(context.Background(), 10000)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch indexed registered operator state - %s", err)
	}
	deregistrations, err := s.subgraphClient.QueryOperatorDeregistrations(context.Background(), 10000)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch indexed deregistered operator state - %s", err)
	}

	operators := make(map[string]int)

	// Add registrations
	for _, registration := range registrations {
		logger.Info("Operator", "operatorId", string(registration.OperatorId), "info", registration)
		operators[string(registration.OperatorId)]++
	}
	// Deduct deregistrations
	for _, deregistration := range deregistrations {
		operators[string(deregistration.OperatorId)]--
	}

	activeOperators := make([]string, 0)
	for operatorId, count := range operators {
		if count > 0 {
			activeOperators = append(activeOperators, operatorId)
		}
	}
	logger.Info("Active operators found", "count", len(activeOperators))

	var wg sync.WaitGroup
	var mu sync.Mutex
	numWorkers := 10
	operatorChan := make(chan string, len(activeOperators))
	semvers := make(map[string]int)
	worker := func() {
		for operatorId := range operatorChan {
			operatorInfo, err := s.getOperatorInfo(ctx, operatorId)
			if err != nil {
				mu.Lock()
				semvers["not-found"]++
				mu.Unlock()
				continue
			}
			operatorSocket := core.OperatorSocket(operatorInfo.Socket)
			dispersalSocket := operatorSocket.GetDispersalSocket()
			semverInfo := semver.GetSemverInfo(context.Background(), dispersalSocket, operatorId, false, logger)

			mu.Lock()
			semvers[semverInfo]++
			mu.Unlock()
		}
		wg.Done()
	}

	// Launch worker goroutines
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker()
	}

	// Send operator IDs to the channel
	for _, operatorId := range activeOperators {
		operatorChan <- operatorId
	}
	close(operatorChan)

	// Wait for all workers to finish
	wg.Wait()

	// Create HostInfoReportResponse instance
	semverReport := &SemverReportResponse{
		Semver: semvers,
	}

	// Publish semver report metrics
	s.metrics.UpdateSemverCounts(semvers)

	return semverReport, nil
}

// method to check if operator is online via socket dial
func checkIsOperatorOnline(socket string, timeoutSecs int, logger logging.Logger) bool {
	if !ValidOperatorIP(socket, logger) {
		logger.Error("port check blocked invalid operator IP", "socket", socket)
		return false
	}
	timeout := time.Second * time.Duration(timeoutSecs)
	conn, err := net.DialTimeout("tcp", socket, timeout)
	if err != nil {
		logger.Warn("port check timeout", "socket", socket, "timeout", timeoutSecs, "error", err)
		return false
	}
	defer conn.Close() // Close the connection after checking
	return true
}
