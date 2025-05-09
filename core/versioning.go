package core

import (
	"context"
	"math/big"
	"time"

	pbvalidator "github.com/Layr-Labs/eigenda/api/grpc/validator"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type OperatorInfoVerbose struct {
	OperatorID OperatorID
	Socket     OperatorSocket
	Stake      *big.Int
	NodeInfo   *pbvalidator.GetNodeInfoReply
}

type OperatorStateVerbose map[QuorumID]map[OperatorIndex]OperatorInfoVerbose

// GetOperatorVerboseState returns the verbose state of all operators within the supplied quorums including their node info.
// The returned state is for the block number supplied.
func GetOperatorVerboseState(ctx context.Context, stakesWithSocket OperatorStakesWithSocket, quorums []QuorumID, blockNumber uint32) (OperatorStateVerbose, error) {
	quorumBytes := make([]byte, len(quorums))
	for ind, quorum := range quorums {
		quorumBytes[ind] = byte(uint8(quorum))
	}

	state := make(OperatorStateVerbose, len(quorums))
	totalOperators := 0
	successfulNodeInfoFetches := 0
	failedNodeInfoFetches := 0

	for _, quorumID := range quorums {
		state[quorumID] = make(map[OperatorIndex]OperatorInfoVerbose, len(stakesWithSocket[quorumID]))

		for j, op := range stakesWithSocket[quorumID] {
			totalOperators++
			operatorIndex := OperatorIndex(j)

			nodeVersion, err := GetNodeInfoFromSocket(ctx, op.Socket)
			if err != nil {
				failedNodeInfoFetches++

				// Instead of failing completely, continue with nil NodeInfo for this operator
				state[quorumID][operatorIndex] = OperatorInfoVerbose{
					OperatorID: op.OperatorID,
					Socket:     op.Socket,
					Stake:      op.Stake,
					NodeInfo:   nil,
				}
				continue
			}
			successfulNodeInfoFetches++

			state[quorumID][operatorIndex] = OperatorInfoVerbose{
				OperatorID: op.OperatorID,
				Socket:     op.Socket,
				Stake:      op.Stake,
				NodeInfo:   nodeVersion,
			}
		}
	}

	return state, nil
}

// GetNodeInfoFromSocket pings the operator's endpoint and returns NodeInfoReply
func GetNodeInfoFromSocket(ctx context.Context, socket OperatorSocket) (*pbvalidator.GetNodeInfoReply, error) {
	endpoint := socket.GetV2DispersalSocket()

	conn, err := grpc.NewClient(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := pbvalidator.NewDispersalClient(conn)
	ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := client.GetNodeInfo(ctxTimeout, &pbvalidator.GetNodeInfoRequest{})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// CalculateQuorumRolloutReadiness computes the stake percentage and readiness for each quorum.
func CalculateQuorumRolloutReadiness(
	ops OperatorStateVerbose,
	requiredVersion string,
	threshold float64,
) (map[QuorumID]float64, map[QuorumID]bool) {
	totalStakeByQuorum := make(map[QuorumID]*big.Int)
	upgradedStakeByQuorum := make(map[QuorumID]*big.Int)

	for quorumID, opMap := range ops {
		totalStake := big.NewInt(0)
		upgradedStake := big.NewInt(0)
		operatorCount := 0
		upgradedOperatorCount := 0

		for _, opState := range opMap {
			operatorCount++
			if opState.Stake == nil {
				continue
			}

			totalStake.Add(totalStake, opState.Stake)

			if opState.NodeInfo != nil {
				if opState.NodeInfo.Semver == requiredVersion {
					upgradedOperatorCount++
					upgradedStake.Add(upgradedStake, opState.Stake)
				}
			}
		}

		totalStakeByQuorum[quorumID] = totalStake
		upgradedStakeByQuorum[quorumID] = upgradedStake
	}

	pctByQuorum := make(map[QuorumID]float64)
	readyByQuorum := make(map[QuorumID]bool)
	for quorum, total := range totalStakeByQuorum {
		upgraded := upgradedStakeByQuorum[quorum]
		pct := 0.0
		if total.Cmp(big.NewInt(0)) > 0 && upgraded != nil {
			pct, _ = new(big.Rat).SetFrac(upgraded, total).Float64()
		}
		isReady := pct >= threshold

		pctByQuorum[quorum] = pct
		readyByQuorum[quorum] = isReady
	}

	return pctByQuorum, readyByQuorum
}
