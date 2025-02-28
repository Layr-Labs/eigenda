package operators

import (
	"math/big"
	"sort"

	"github.com/Layr-Labs/eigenda/core"
)

type OperatorStakeShare struct {
	OperatorId core.OperatorID
	StakeShare float64
}

// The GetRankedOperators returns ranked operators list, by total-quorum-stake and by individual
// quorums.
func GetRankedOperators(state *core.OperatorState) ([]*OperatorStakeShare, map[uint8][]*OperatorStakeShare) {
	tqsRankedOperators := make([]*OperatorStakeShare, 0)
	quorumRankedOperators := make(map[uint8][]*OperatorStakeShare)
	tqs := make(map[core.OperatorID]*OperatorStakeShare)
	for q, operators := range state.Operators {
		operatorStakeShares := make([]*OperatorStakeShare, 0)
		totalStake := new(big.Float).SetInt(state.Totals[q].Stake)
		for opId, opInfo := range operators {
			opStake := new(big.Float).SetInt(opInfo.Stake)
			share, _ := new(big.Float).Quo(
				new(big.Float).Mul(opStake, big.NewFloat(10000)),
				totalStake).Float64()
			operatorStakeShares = append(operatorStakeShares, &OperatorStakeShare{OperatorId: opId, StakeShare: share})
		}
		// Descending order by stake share in the quorum.
		sort.Slice(operatorStakeShares, func(i, j int) bool {
			if operatorStakeShares[i].StakeShare == operatorStakeShares[j].StakeShare {
				return operatorStakeShares[i].OperatorId.Hex() < operatorStakeShares[j].OperatorId.Hex()
			}
			return operatorStakeShares[i].StakeShare > operatorStakeShares[j].StakeShare
		})

		for _, op := range operatorStakeShares {
			quorumRankedOperators[q] = append(quorumRankedOperators[q], op)
			if _, ok := tqs[op.OperatorId]; !ok {
				tqs[op.OperatorId] = &OperatorStakeShare{OperatorId: op.OperatorId, StakeShare: op.StakeShare}
			} else {
				tqs[op.OperatorId].StakeShare += op.StakeShare
			}
		}
	}
	for _, op := range tqs {
		tqsRankedOperators = append(tqsRankedOperators, op)
	}
	// Descending order by total stake share across the quorums.
	sort.Slice(tqsRankedOperators, func(i, j int) bool {
		if tqsRankedOperators[i].StakeShare == tqsRankedOperators[j].StakeShare {
			return tqsRankedOperators[i].OperatorId.Hex() < tqsRankedOperators[j].OperatorId.Hex()
		}
		return tqsRankedOperators[i].StakeShare > tqsRankedOperators[j].StakeShare
	})
	return tqsRankedOperators, quorumRankedOperators
}
