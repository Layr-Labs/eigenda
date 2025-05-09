package core

import (
	"context"
	"fmt"
	"math/big"
	"time"

	pbvalidator "github.com/Layr-Labs/eigenda/api/grpc/validator"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// OperatorInfoVerbose contains information about an operator including its node information
type OperatorInfoVerbose struct {
	OperatorID OperatorID
	Socket     OperatorSocket
	Stake      *big.Int
	NodeInfo   *pbvalidator.GetNodeInfoReply
}

type OperatorStateVerbose map[QuorumID]map[OperatorIndex]OperatorInfoVerbose

// VersionCheckConfig contains configuration for version checking
type VersionCheckConfig struct {
	// RequiredVersion is the minimum required version in semver format (e.g. ">=0.9.0-rc.1")
	RequiredVersion string
	// StakeThreshold is the percentage of stake that needs to be at the required version
	// for a quorum to be considered ready (e.g. 0.8 for 80%)
	StakeThreshold float64
	// CheckInterval is how often to check version info
	CheckInterval time.Duration
	// Enabled determines whether to enforce version checks
	Enabled bool
}

// VersionChecker provides version checking functionality
type VersionChecker struct {
	config      VersionCheckConfig
	logger      logging.Logger
	chainReader Reader

	operatorStateByQuorum   map[QuorumID]map[OperatorIndex]OperatorInfoVerbose
	rolloutReadyByQuorum    map[QuorumID]bool
	rolloutStakePctByQuorum map[QuorumID]float64
}

// NewVersionChecker creates a new version checker
func NewVersionChecker(config VersionCheckConfig, logger logging.Logger, chainReader Reader) *VersionChecker {
	return &VersionChecker{
		config:                  config,
		logger:                  logger,
		chainReader:             chainReader,
		operatorStateByQuorum:   make(map[QuorumID]map[OperatorIndex]OperatorInfoVerbose),
		rolloutReadyByQuorum:    make(map[QuorumID]bool),
		rolloutStakePctByQuorum: make(map[QuorumID]float64),
	}
}

// QuorumVersionProbe performs a check of operator node versions, updating internal state
func (vc *VersionChecker) QuorumVersionProbe(ctx context.Context, quorumIds []QuorumID) error {
	vc.logger.Debug("Starting operator version check")

	currentBlock, err := vc.chainReader.GetCurrentBlockNumber(ctx)
	if err != nil {
		vc.logger.Error("failed to get current block number", "err", err)
		return fmt.Errorf("failed to get current block number: %w", err)
	}
	stakesWithSocket, err := vc.chainReader.GetOperatorStakesWithSocketForQuorums(ctx, quorumIds, currentBlock)
	if err != nil {
		vc.logger.Error("failed to get operator stakes with socket", "err", err)
		return fmt.Errorf("failed to get operator stakes with socket: %w", err)
	}
	operatorState, err := GetOperatorVerboseState(ctx, stakesWithSocket, quorumIds)
	if err != nil {
		return fmt.Errorf("failed to get operator info for quorums: %w", err)
	}

	pctByQuorum, rolloutReady := CalculateQuorumRolloutReadiness(
		operatorState,
		vc.config.RequiredVersion,
		vc.config.StakeThreshold,
	)

	// Log detailed results for debugging
	for quorum, pct := range pctByQuorum {
		ready := rolloutReady[quorum]
		vc.logger.Debug("Operator version rollout check result",
			"quorum", quorum,
			"required_version", vc.config.RequiredVersion,
			"stake_pct", pct,
			"threshold", vc.config.StakeThreshold,
			"rollout_ready", ready,
			"upgradedThreshold", pct >= vc.config.StakeThreshold)
	}

	// If no quorums were found, log a warning
	if len(pctByQuorum) == 0 {
		vc.logger.Warn("No quorums found in rollout readiness calculation")
	}

	// Update internal state
	for quorum, pct := range pctByQuorum {
		vc.rolloutStakePctByQuorum[quorum] = pct
	}
	for quorum, ready := range rolloutReady {
		vc.rolloutReadyByQuorum[quorum] = ready
	}
	for quorum, state := range operatorState {
		vc.operatorStateByQuorum[quorum] = state
	}

	return nil
}

// IsQuorumRolloutReady checks if specified quorums are ready for the new version
func (vc *VersionChecker) IsQuorumRolloutReady(quorumIDs []QuorumID) (bool, []string) {
	if !vc.config.Enabled {
		vc.logger.Info("Version checking is disabled, assuming all quorums are ready")
		return true, nil
	}

	notReady := make([]string, 0)
	for _, quorum := range quorumIDs {
		ready, exists := vc.rolloutReadyByQuorum[quorum]
		if !exists {
			vc.logger.Warn("Quorum not found in rollout readiness map", "quorumID", quorum)
			ready = false // Default to not ready if missing
		}

		pct, pctExists := vc.rolloutStakePctByQuorum[quorum]
		if !pctExists {
			vc.logger.Warn("Quorum not found in stake percentage map", "quorumID", quorum)
			pct = 0 // Default to 0% if missing
		}

		vc.logger.Info("Checking quorum readiness", "quorumID", quorum, "ready", ready, "stakePct", pct)
		if !ready {
			threshold := vc.config.StakeThreshold * 100
			notReady = append(notReady, fmt.Sprintf("quorum %d: %.2f%% of %.2f%% upgraded", quorum, pct*100, threshold))
		}
	}

	if len(notReady) > 0 {
		return false, notReady
	}

	vc.logger.Info("All specified quorums are rollout ready")
	return true, nil
}

// GetRolloutStatus returns the current rollout status for all quorums
func (vc *VersionChecker) GetRolloutStatus() (map[QuorumID]float64, map[QuorumID]bool) {
	return vc.rolloutStakePctByQuorum, vc.rolloutReadyByQuorum
}

// RefreshTrackedQuorums refreshes the version info for all currently tracked quorums in the cache.
func (vc *VersionChecker) RefreshTrackedQuorums(ctx context.Context) error {
	quorumIds := make([]QuorumID, 0, len(vc.rolloutReadyByQuorum))
	for quorum := range vc.rolloutReadyByQuorum {
		quorumIds = append(quorumIds, quorum)
	}
	return vc.QuorumVersionProbe(ctx, quorumIds)
}

// GetOperatorVerboseState returns the verbose state of all operators within the supplied quorums including their node info.
// The returned state is for the block number supplied.
func GetOperatorVerboseState(ctx context.Context, stakesWithSocket OperatorStakesWithSocket, quorums []QuorumID) (OperatorStateVerbose, error) {
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
				// // Check if version meets required constraint
				// currentVersion, err := version.NewVersion(opState.NodeInfo.Semver)
				// if err != nil {
				// 	continue
				// }

				// constraint, err := version.NewConstraint(requiredVersion)
				// if err != nil {
				// 	continue
				// }

				// if constraint.Check(currentVersion) {
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
