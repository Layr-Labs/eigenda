package dataapi

import (
	"context"
	"errors"
	"net"
	"sort"
	"time"

	"github.com/Layr-Labs/eigenda/core"
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
	operatorOnlineStatusresultsChan chan *DeregisteredOperatorMetadata
)

func (s *server) getDeregisteredOperatorForDays(ctx context.Context, days int32) ([]*DeregisteredOperatorMetadata, error) {
	// Track time taken to get deregistered operators
	startTime := time.Now()

	indexedDeregisteredOperatorState, err := s.subgraphClient.QueryIndexedDeregisteredOperatorsForTimeWindow(ctx, days)
	if err != nil {
		return nil, err
	}

	// Convert the map to a slice.
	operators := indexedDeregisteredOperatorState.Operators

	operatorOnlineStatusresultsChan = make(chan *DeregisteredOperatorMetadata, len(operators))
	processOperatorOnlineCheck(indexedDeregisteredOperatorState, operatorOnlineStatusresultsChan, s.logger)

	// Collect results of work done
	DeregisteredOperatorMetadata := make([]*DeregisteredOperatorMetadata, 0, len(operators))
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

func processOperatorOnlineCheck(deregisteredOperatorState *IndexedDeregisteredOperatorState, operatorOnlineStatusresultsChan chan<- *DeregisteredOperatorMetadata, logger logging.Logger) {
	operators := deregisteredOperatorState.Operators
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

func checkIsOnlineAndProcessOperator(operatorStatus OperatorOnlineStatus, operatorOnlineStatusresultsChan chan<- *DeregisteredOperatorMetadata, logger logging.Logger) {
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
	metadata := &DeregisteredOperatorMetadata{
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
	operatorInfo, err := s.subgraphClient.QueryOperatorInfoByOperatorId(context.Background(), operatorId)
	if err != nil {
		s.logger.Warn("failed to fetch operator info", "error", err)
		return &OperatorPortCheckResponse{}, errors.New("not found")
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
	s.logger.Info("operator port check response", portCheckResponse)

	// Send the metadata to the results channel
	return portCheckResponse, nil
}

// method to check if operator is online
// Note: This method is least intrusive way to check if operator is online
// AlternateSolution: Should we add an endpt to check if operator is online?
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
