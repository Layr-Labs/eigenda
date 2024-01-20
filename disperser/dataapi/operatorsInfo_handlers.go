package dataapi

import (
	"context"
	"net"
	"sort"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/core"
)

type OperatorOnlineStatus struct {
	OperatorInfo        *Operator
	IndexedOperatorInfo *core.IndexedOperatorInfo
}

var (
	// TODO: this should be configurable
	numWorkers                      = 10
	operatorOnlineStatusChan        chan OperatorOnlineStatus
	operatorOnlineStatusresultsChan chan *DeregisteredOperatorMetadata
)

func (s *server) getDeregisterdOperatorForDays(ctx context.Context, days int32) ([]*DeregisteredOperatorMetadata, error) {
	// Track Time taken to get deregistered operators
	startTime := time.Now()

	indexedDeregisteredOperatorState, err := s.subgraphClient.QueryIndexedDeregisteredOperatorsForTimeWindow(ctx, days)
	if err != nil {
		return nil, err
	}

	// Convert the map to a slice.
	operators := indexedDeregisteredOperatorState.Operators

	operatorOnlineStatusChan = make(chan OperatorOnlineStatus, len(operators))
	operatorOnlineStatusresultsChan = make(chan *DeregisteredOperatorMetadata, len(operators))
	processOperatorsInParallel(indexedDeregisteredOperatorState, operatorOnlineStatusChan, operatorOnlineStatusresultsChan, s.logger)

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

// method to check if operator is online
func checkIsOperatorOnline(ipAddress string) bool {
	timeout := time.Second * 10
	conn, err := net.DialTimeout("tcp", ipAddress, timeout)
	if err != nil {
		return false
	}
	defer conn.Close() // Close the connection after checking
	return true
}

// Helper Function to Process Operators in Parallel
func processOperatorsInParallel(deRegisteredOperatorState *IndexedDeregisteredOperatorState, operatorOnlineStatusChan chan OperatorOnlineStatus, operatorOnlineStatusresultsChan chan<- *DeregisteredOperatorMetadata, logger common.Logger) {

	operators := deRegisteredOperatorState.Operators
	// Start worker goroutines
	for i := 0; i < numWorkers; i++ {
		go func() {
			for item := range operatorOnlineStatusChan {
				ipAddress := core.OperatorSocket(item.IndexedOperatorInfo.Socket).GetRetrievalSocket()
				isOnline := checkIsOperatorOnline(ipAddress)

				// Log the online status
				if isOnline {
					logger.Debug("Operator %v is online at %s", item.IndexedOperatorInfo, ipAddress)
				} else {
					logger.Debug("Operator %v is offline at %s", item.IndexedOperatorInfo, ipAddress)
				}

				// Create the metadata regardless of online status
				metadata := &DeregisteredOperatorMetadata{
					OperatorId:  string(item.OperatorInfo.OperatorId[:]),
					BlockNumber: uint(item.OperatorInfo.BlockNumber),
					IpAddress:   ipAddress,
					IsOnline:    isOnline,
				}

				// Send the metadata to the results channel
				operatorOnlineStatusresultsChan <- metadata
			}
		}()
	}

	// Send work to the workers
	for _, operatorInfo := range operators {
		operatorOnlineStatus := OperatorOnlineStatus{
			OperatorInfo:        operatorInfo.Metadata,
			IndexedOperatorInfo: operatorInfo.IndexedOperatorInfo,
		}
		operatorOnlineStatusChan <- operatorOnlineStatus
	}
	close(operatorOnlineStatusChan) // Close the channel after sending all tasks

}
