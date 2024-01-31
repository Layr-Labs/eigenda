package dataapi

import (
	"context"
	"net"
	"sort"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/gammazero/workerpool"
)

type OperatorOnlineStatus struct {
	OperatorInfo        *Operator
	IndexedOperatorInfo *core.IndexedOperatorInfo
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
	s.logger.Info("Time taken to get deregistered operators for days: %v", time.Since(startTime))
	sort.Slice(DeregisteredOperatorMetadata, func(i, j int) bool {
		return DeregisteredOperatorMetadata[i].BlockNumber < DeregisteredOperatorMetadata[j].BlockNumber
	})

	return DeregisteredOperatorMetadata, nil
}

func processOperatorOnlineCheck(deregisteredOperatorState *IndexedDeregisteredOperatorState, operatorOnlineStatusresultsChan chan<- *DeregisteredOperatorMetadata, logger common.Logger) {
	operators := deregisteredOperatorState.Operators
	wp := workerpool.New(poolSize)

	for _, operatorInfo := range operators {
		operatorStatus := OperatorOnlineStatus{
			OperatorInfo:        operatorInfo.Metadata,
			IndexedOperatorInfo: operatorInfo.IndexedOperatorInfo,
		}

		// Submit each operator status check to the worker pool
		wp.Submit(func() {
			checkIsOnlineAndProcessOperator(operatorStatus, operatorOnlineStatusresultsChan, logger)
		})
	}

	wp.StopWait() // Wait for all submitted tasks to complete and stop the pool
}

func checkIsOnlineAndProcessOperator(operatorStatus OperatorOnlineStatus, operatorOnlineStatusresultsChan chan<- *DeregisteredOperatorMetadata, logger common.Logger) {
	socket := core.OperatorSocket(operatorStatus.IndexedOperatorInfo.Socket).GetRetrievalSocket()
	isOnline := checkIsOperatorOnline(socket)

	// Log the online status
	if isOnline {
		logger.Debug("Operator %v is online at %s", operatorStatus.IndexedOperatorInfo, socket)
	} else {
		logger.Debug("Operator %v is offline at %s", operatorStatus.IndexedOperatorInfo, socket)
	}

	// Create the metadata regardless of online status
	metadata := &DeregisteredOperatorMetadata{
		OperatorId:  string(operatorStatus.OperatorInfo.OperatorId[:]),
		BlockNumber: uint(operatorStatus.OperatorInfo.BlockNumber),
		Socket:      socket,
		IsOnline:    isOnline,
	}

	// Send the metadata to the results channel
	operatorOnlineStatusresultsChan <- metadata
}

// method to check if operator is online
// Note: This method is least intrusive wat to check if operator is online
// AlternateSolution: Should we add an endpt to check if operator is online?
func checkIsOperatorOnline(socket string) bool {
	timeout := time.Second * 10
	conn, err := net.DialTimeout("tcp", socket, timeout)
	if err != nil {
		return false
	}
	defer conn.Close() // Close the connection after checking
	return true
}
