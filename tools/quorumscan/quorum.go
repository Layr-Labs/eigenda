package quorumscan

import (
	"fmt"
	"math/big"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigensdk-go/logging"
)

type QuorumMetrics struct {
	Operators        []string           `json:"operators"`
	OperatorStake    map[string]float64 `json:"operator_stake"`
	OperatorStakePct map[string]float64 `json:"operator_stake_pct"`
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
				}
			}
			stakePercentage := float64(0)
			effectiveStakePercentage := float64(0)
			if stake, ok := operatorState.Operators[quorum][operatorId]; ok {
				totalStake := new(big.Float).SetInt(totalOperatorInfo.Stake)
				totalEffectiveStake := new(big.Float).SetInt(totalOperatorInfo.EffectiveStake)
				operatorStake := new(big.Float).SetInt(stake.Stake)
				operatorEffectiveStake := new(big.Float).SetInt(stake.EffectiveStake)
				stakePercentage, _ = new(big.Float).Mul(big.NewFloat(100), new(big.Float).Quo(operatorStake, totalStake)).Float64()
				effectiveStakePercentage, _ = new(big.Float).Mul(big.NewFloat(100), new(big.Float).Quo(operatorEffectiveStake, totalEffectiveStake)).Float64()
				stakeValue, _ := operatorStake.Float64()
				effectiveStakeValue, _ := operatorEffectiveStake.Float64()
				metrics[quorum].Operators = append(metrics[quorum].Operators, operatorId.Hex())
				metrics[quorum].OperatorStake[operatorId.Hex()] = stakeValue
				metrics[quorum].OperatorStakePct[operatorId.Hex()] = stakePercentage
				fmt.Println("effectiveStakePercentage", effectiveStakePercentage, "effectiveStakeValue", effectiveStakeValue)
				// metrics[quorum].OperatorEffectiveStake[operatorId.Hex()] = effectiveStakeValue
				// metrics[quorum].OperatorEffectiveStakePct[operatorId.Hex()] = effectiveStakePercentage
			}
		}
	}

	return metrics
}
