package semver

import (
	"context"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/Layr-Labs/eigenda/api/grpc/node"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type SemverMetrics struct {
	Semver                string            `json:"semver"`
	Operators             uint8             `json:"count"`
	OperatorIds           []string          `json:"operators"`
	QuorumStakePercentage map[uint8]float64 `json:"stake_percentage"`
}

func ScanOperators(operators map[core.OperatorID]*core.IndexedOperatorInfo, operatorState *core.OperatorState, useRetrievalSocket bool, numWorkers int, nodeInfoTimeout time.Duration, logger logging.Logger) map[string]*SemverMetrics {
	var wg sync.WaitGroup
	var mu sync.Mutex
	semvers := make(map[string]*SemverMetrics)
	operatorChan := make(chan core.OperatorID, len(operators))
	worker := func() {
		for operatorId := range operatorChan {
			operatorSocket := core.OperatorSocket(operators[operatorId].Socket)
			var socket string
			if useRetrievalSocket {
				socket = operatorSocket.GetRetrievalSocket()
			} else {
				socket = operatorSocket.GetDispersalSocket()
			}
			semver := GetSemverInfo(context.Background(), socket, useRetrievalSocket, operatorId, logger, nodeInfoTimeout)

			mu.Lock()
			if _, exists := semvers[semver]; !exists {
				semvers[semver] = &SemverMetrics{
					Semver:                semver,
					Operators:             1,
					OperatorIds:           []string{operatorId.Hex()},
					QuorumStakePercentage: make(map[uint8]float64),
				}
			} else {
				semvers[semver].Operators += 1
				semvers[semver].OperatorIds = append(semvers[semver].OperatorIds, operatorId.Hex())
			}

			// Calculate stake percentage for each quorum
			for quorum, totalOperatorInfo := range operatorState.Totals {
				stakePercentage := float64(0)
				if stake, ok := operatorState.Operators[quorum][operatorId]; ok {
					totalStake := new(big.Float).SetInt(totalOperatorInfo.Stake)
					operatorStake := new(big.Float).SetInt(stake.Stake)
					stakePercentage, _ = new(big.Float).Mul(big.NewFloat(100), new(big.Float).Quo(operatorStake, totalStake)).Float64()
				}

				if _, exists := semvers[semver].QuorumStakePercentage[quorum]; !exists {
					semvers[semver].QuorumStakePercentage[quorum] = stakePercentage
				} else {
					semvers[semver].QuorumStakePercentage[quorum] += stakePercentage
				}
			}
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
func GetSemverInfo(ctx context.Context, socket string, userRetrievalClient bool, operatorId core.OperatorID, logger logging.Logger, timeout time.Duration) string {
	conn, err := grpc.NewClient(socket, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return "unreachable"
	}
	defer conn.Close()
	ctxWithTimeout, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	var reply *node.NodeInfoReply
	if userRetrievalClient {
		client := node.NewRetrievalClient(conn)
		reply, err = client.NodeInfo(ctxWithTimeout, &node.NodeInfoRequest{})
	} else {
		client := node.NewDispersalClient(conn)
		reply, err = client.NodeInfo(ctxWithTimeout, &node.NodeInfoRequest{})
	}
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

		logger.Warn("NodeInfo", "operatorId", operatorId.Hex(), "semver", semver, "error", err)
		return semver
	}

	// local node source compiles without semver
	if reply.Semver == "" {
		reply.Semver = "0.8.4"
	}

	logger.Info("NodeInfo", "operatorId", operatorId.Hex(), "socket", socket, "userRetrievalClient", userRetrievalClient, "semver", reply.Semver, "os", reply.Os, "arch", reply.Arch, "numCpu", reply.NumCpu, "memBytes", reply.MemBytes)
	return reply.Semver
}
