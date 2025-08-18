package quorumscan

import (
	"math/big"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigensdk-go/logging"
)

type QuorumMetrics struct {
	Operators        []string           `json:"operators"`
	OperatorStake    map[string]float64 `json:"operator_stake"`
	OperatorStakePct map[string]float64 `json:"operator_stake_pct"`
	OperatorSocket   map[string]string  `json:"operator_socket"`
	BlockNumber      uint               `json:"block_number"`
}

func QuorumScan(operators map[core.OperatorID]*core.IndexedOperatorInfo, operatorState *core.OperatorState, logger logging.Logger) map[uint8]*QuorumMetrics {
	metrics := make(map[uint8]*QuorumMetrics)
	for operatorId := range operators {

		// Calculate stake percentage for each quorum
		for quorum, totalOperatorInfo := range operatorState.Totals {
			if _, exists := metrics[quorum]; !exists {
				metrics[quorum] = &QuorumMetrics{
					Operators:        []string{},
					OperatorStakePct: make(map[string]float64),
					OperatorStake:    make(map[string]float64),
					OperatorSocket:   make(map[string]string),
					BlockNumber:      operatorState.BlockNumber,
				}
			}
			stakePercentage := float64(0)
			if stake, ok := operatorState.Operators[quorum][operatorId]; ok {
				totalStake := new(big.Float).SetInt(totalOperatorInfo.Stake)
				operatorStake := new(big.Float).SetInt(stake.Stake)
				stakePercentage, _ = new(big.Float).Mul(big.NewFloat(100), new(big.Float).Quo(operatorStake, totalStake)).Float64()
				stakeValue, _ := operatorStake.Float64()
				metrics[quorum].Operators = append(metrics[quorum].Operators, operatorId.Hex())
				metrics[quorum].OperatorStake[operatorId.Hex()] = stakeValue
				metrics[quorum].OperatorStakePct[operatorId.Hex()] = stakePercentage
				metrics[quorum].OperatorSocket[operatorId.Hex()] = operators[operatorId].Socket
			}
		}
	}

	return metrics
}
