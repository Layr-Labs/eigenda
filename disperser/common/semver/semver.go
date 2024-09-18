package semver

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/Layr-Labs/eigenda/api/grpc/node"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/disperser/dataapi"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// GetOperatorInfo retrieves information about an operator registered ip/ports
func GetOperatorInfo(subgraphClient dataapi.SubgraphClient, operatorId string, logger logging.Logger) (*core.IndexedOperatorInfo, error) {
	operatorInfo, err := subgraphClient.QueryOperatorInfoByOperatorId(context.Background(), operatorId)
	if err != nil {
		logger.Warn("failed to fetch operator info", "operatorId", operatorId, "error", err)
		return nil, fmt.Errorf("operator info not found for operatorId %s", operatorId)
	}
	return operatorInfo, nil
}

// scanOperators scans for available operators.
func ScanOperatorsHostInfo(ctx context.Context, subgraphClient dataapi.SubgraphClient, logger logging.Logger) (map[string]int, error) {
	registrations, err := subgraphClient.QueryOperatorsWithLimit(ctx, 10000)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch indexed registered operator state - %s", err)
	}
	deregistrations, err := subgraphClient.QueryOperatorDeregistrations(ctx, 10000)
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
			operatorInfo, err := GetOperatorInfo(subgraphClient, operatorId, logger)
			if err != nil {
				mu.Lock()
				semvers["not-found"]++
				mu.Unlock()
				continue
			}
			operatorSocket := core.OperatorSocket(operatorInfo.Socket)
			dispersalSocket := operatorSocket.GetDispersalSocket()
			semverInfo := GetSemverInfo(context.Background(), dispersalSocket, operatorId, logger)

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

	return semvers, nil
}

// query operator host info endpoint if available
func GetSemverInfo(ctx context.Context, socket string, operatorId string, logger logging.Logger) string {
	conn, err := grpc.Dial(socket, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return "unreachable"
	}
	defer conn.Close()
	client := node.NewDispersalClient(conn)
	ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Second*time.Duration(3))
	defer cancel()
	reply, err := client.NodeInfo(ctxWithTimeout, &node.NodeInfoRequest{})
	if err != nil {
		var semver string
		if strings.Contains(err.Error(), "Unimplemented") {
			semver = "<0.8.0"
		} else if strings.Contains(err.Error(), "DeadlineExceeded") {
			semver = "timeout"
		} else if strings.Contains(err.Error(), "Unavailable") {
			semver = "refused"
		} else {
			semver = "error"
		}

		logger.Warn("NodeInfo", "operatorId", operatorId, "semver", semver, "error", err)
		return semver
	}

	// local node source compiles without semver
	if reply.Semver == "" {
		reply.Semver = "src-compile"
	}

	logger.Info("NodeInfo", "operatorId", operatorId, "socker", socket, "semver", reply.Semver, "os", reply.Os, "arch", reply.Arch, "numCpu", reply.NumCpu, "memBytes", reply.MemBytes)
	return reply.Semver
}
