package dataapi

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/Layr-Labs/eigenda/core"
)

func (s *server) getDeregisterdOperatorForDays(ctx context.Context, days int32) ([]*DeregisteredOperatorInfo, error) {

	indexedDeregisteredOperatorState, err := s.subgraphClient.QueryIndexedDeregisteredOperatorsForTimeWindow(ctx, days)
	if err != nil {
		return nil, err
	}

	// Convert the map to a slice.

	operators := indexedDeregisteredOperatorState.Operators
	deRegisteredOperators := make([]*DeregisteredOperatorInfo, 0, len(operators))

	for _, operatorInfo := range operators {
		indexedOperatorInfo := operatorInfo.IndexedOperatorInfo
		// Try pinging deregistered if offline add to the list
		if !checkIsOperatorOnline(indexedOperatorInfo.Socket) {
			fmt.Printf("DataAPI Operator %v is offline\n", indexedOperatorInfo)
			deRegisteredOperators = append(deRegisteredOperators, operatorInfo)
		}
	}

	return deRegisteredOperators, nil
}

// method to check if operator is online
func checkIsOperatorOnline(socket string) bool {

	ipAddress := core.OperatorSocket(socket).GetRetrievalSocket()
	timeout := time.Second * 10
	conn, err := net.DialTimeout("tcp", ipAddress, timeout)
	if err != nil {
		// The server is not responding or not reachable
		return false
	}
	defer conn.Close() // Close the connection after checking
	return true
}
