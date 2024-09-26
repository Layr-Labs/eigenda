package semver

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/Layr-Labs/eigenda/api/grpc/node"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func ScanOperators(operators map[core.OperatorID]*core.IndexedOperatorInfo, numWorkers int, nodeInfoTimeout time.Duration, logger logging.Logger) map[string]int {
	var wg sync.WaitGroup
	var mu sync.Mutex
	semvers := make(map[string]int)
	operatorChan := make(chan core.OperatorID, len(operators))
	worker := func() {
		for operatorId := range operatorChan {
			operatorSocket := core.OperatorSocket(operators[operatorId].Socket)
			dispersalSocket := operatorSocket.GetDispersalSocket()
			semver := GetSemverInfo(context.Background(), dispersalSocket, operatorId, logger, nodeInfoTimeout)

			mu.Lock()
			semvers[semver]++
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
	for operatorId := range operators {
		operatorChan <- operatorId
	}
	close(operatorChan)

	// Wait for all workers to finish
	wg.Wait()
	return semvers
}

// query operator host info endpoint if available
func GetSemverInfo(ctx context.Context, socket string, operatorId core.OperatorID, logger logging.Logger, timeout time.Duration) string {
	conn, err := grpc.Dial(socket, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return "unreachable"
	}
	defer conn.Close()
	ctxWithTimeout, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	client := node.NewDispersalClient(conn)
	reply, err := client.NodeInfo(ctxWithTimeout, &node.NodeInfoRequest{})
	if err != nil {
		var semver string
		if strings.Contains(err.Error(), "unknown method NodeInfo") {
			semver = "<0.8.0"
		} else if strings.Contains(err.Error(), "unknown service") {
			semver = "filtered"
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
