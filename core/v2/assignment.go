package v2

import (
	"fmt"
	"math/big"
	"sort"

	"github.com/Layr-Labs/eigenda/core"
)

const stakeCapIterations = 1

// GetAssignments calculates chunk assignments for operators in a quorum based on their stake
// numIterations specifies the number of iterative capping rounds to apply (default: 2)
func GetAssignments(state *core.OperatorState, blobParams *core.BlobVersionParameters, quorum uint8) (map[core.OperatorID]Assignment, error) {
	if state == nil {
		return nil, fmt.Errorf("state cannot be nil")
	}

	if blobParams == nil {
		return nil, fmt.Errorf("blob params cannot be nil")
	}

	ops, ok := state.Operators[quorum]
	if !ok {
		return nil, fmt.Errorf("no operators found for quorum %d", quorum)
	}

	numOps := len(ops)
	if uint32(numOps) > blobParams.MaxNumOperators {
		return nil, fmt.Errorf("too many operators (%d) to get assignments: max number of operators is %d", numOps, blobParams.MaxNumOperators)
	}

	// Early return for empty operator set
	if numOps == 0 {
		return make(map[core.OperatorID]Assignment), nil
	}

	// Calculate effective stakes using the recursive helper function with specified iterations
	// Cap percentage is 1/codingRate - this ensures that with higher redundancy (higher coding rate),
	// we allow a smaller percentage of total stake for any single operator
	capPercentage := 1.0 / float64(blobParams.CodingRate)
	totalEffectiveStake, err := calculateOperatorEffectiveStake(ops, state.Totals[quorum].Stake, capPercentage, stakeCapIterations)
	if err != nil {
		return nil, err
	}

	type operatorAssignment struct {
		id     core.OperatorID
		index  uint32
		chunks uint32
		stake  *big.Int
	}

	numOperatorsBig := big.NewInt(int64(numOps))
	numChunksBig := big.NewInt(int64(blobParams.NumChunks))

	// Calculate number of chunks - numOperators once and reuse
	diffChunksOps := new(big.Int).Sub(numChunksBig, numOperatorsBig)
	chunkAssignments := make([]operatorAssignment, 0, numOps)
	// Calculate initial chunk assignments based on effective stake
	totalCalculatedChunks := uint32(0)
	for ID, r := range ops {
		// Calculate chunks for this operator: (effectiveStake * (numChunks - numOperators)) / totalEffectiveStake (rounded up)
		num := new(big.Int).Mul(r.EffectiveStake, diffChunksOps)
		chunks := uint32(core.RoundUpDivideBig(num, totalEffectiveStake).Uint64())

		chunkAssignments = append(chunkAssignments, operatorAssignment{
			id:     ID,
			index:  uint32(r.Index),
			chunks: chunks,
			stake:  r.EffectiveStake, // Use effective stake for sorting
		})

		totalCalculatedChunks += chunks
	}

	// Sort by effective stake (decreasing) with index as tie-breaker
	sort.Slice(chunkAssignments, func(i, j int) bool {
		stakeCmp := chunkAssignments[i].stake.Cmp(chunkAssignments[j].stake)
		if stakeCmp == 0 {
			return chunkAssignments[i].index < chunkAssignments[j].index
		}
		return stakeCmp > 0 // Sort in descending order
	})

	// Distribute any remaining chunks
	delta := int(blobParams.NumChunks) - int(totalCalculatedChunks)
	if delta < 0 {
		return nil, fmt.Errorf("total chunks %d exceeds maximum %d", totalCalculatedChunks, blobParams.NumChunks)
	}

	assignments := make(map[core.OperatorID]Assignment, numOps)
	index := uint32(0)

	// Assign chunks to operators
	for i, a := range chunkAssignments {
		// Add remaining chunks to operators with highest effective stake first
		if i < delta {
			a.chunks++
		}

		// Always add operators to the assignments map, even with zero chunks
		assignments[a.id] = Assignment{
			StartIndex: index,
			NumChunks:  a.chunks,
		}

		index += a.chunks
	}

	return assignments, nil
}

// GetAssignment returns the assignment for a specific operator
func GetAssignment(state *core.OperatorState, blobParams *core.BlobVersionParameters, quorum core.QuorumID, id core.OperatorID) (Assignment, error) {
	if blobParams == nil {
		return Assignment{}, fmt.Errorf("blob params cannot be nil")
	}

	assignments, err := GetAssignments(state, blobParams, quorum)
	if err != nil {
		return Assignment{}, err
	}

	assignment, ok := assignments[id]
	if !ok {
		return Assignment{}, ErrNotFound
	}

	return assignment, nil
}

// GetChunkLength calculates the chunk length based on blob length and parameters
func GetChunkLength(blobLength uint32, blobParams *core.BlobVersionParameters) (uint32, error) {
	if blobLength == 0 {
		return 0, fmt.Errorf("blob length must be greater than 0")
	}

	if blobParams == nil {
		return 0, fmt.Errorf("blob params cannot be nil")
	}

	// Check that the blob length is a power of 2 using bit manipulation
	if blobLength&(blobLength-1) != 0 {
		return 0, fmt.Errorf("blob length %d is not a power of 2", blobLength)
	}

	chunkLength := blobLength * blobParams.CodingRate / blobParams.NumChunks
	if chunkLength == 0 {
		chunkLength = 1
	}

	return chunkLength, nil
}

// applyStakeCap applies a cap to an operator's stake
// Returns min(operatorStake, stakeCap)
func applyStakeCap(operatorStake, stakeCap *big.Int) (*big.Int, error) {
	if operatorStake == nil || operatorStake.Sign() < 0 {
		return nil, fmt.Errorf("operator stake must be non-negative")
	}

	if stakeCap == nil || stakeCap.Sign() <= 0 {
		return nil, fmt.Errorf("stake cap must be positive")
	}

	// Effective stake = min(nominal stake, stake cap)
	effectiveStake := new(big.Int)
	if operatorStake.Cmp(stakeCap) <= 0 {
		effectiveStake.Set(operatorStake)
	} else {
		effectiveStake.Set(stakeCap)
	}

	return effectiveStake, nil
}

// calculateOperatorEffectiveStake calculates the effective stake for all operators by applying
// a cap recursively for the specified number of iterations.
// Returns the effective stakes for all operators and the new total effective stake.
func calculateOperatorEffectiveStake(
	operators map[core.OperatorID]*core.OperatorInfo,
	totalStake *big.Int,
	capPercentage float64,
	numIterations int,
) (*big.Int, error) {
	if totalStake == nil || totalStake.Sign() <= 0 {
		return nil, fmt.Errorf("total stake must be positive")
	}

	if capPercentage <= 0 || capPercentage > 1 {
		return nil, fmt.Errorf("cap percentage must be between 0 and 1")
	}

	if numIterations <= 0 {
		// Base case: calculate final total effective stake
		totalEffectiveStake := new(big.Int)
		for _, opInfo := range operators {
			if opInfo.EffectiveStake != nil {
				totalEffectiveStake = new(big.Int).Add(totalEffectiveStake, opInfo.EffectiveStake)
			} else {
				// If EffectiveStake is not set, use original stake
				totalEffectiveStake = new(big.Int).Add(totalEffectiveStake, opInfo.Stake)
			}
		}
		return totalEffectiveStake, nil
	}

	// Apply cap to all operators for this iteration
	iterationTotalStake := new(big.Int).Set(totalStake)
	newTotalEffectiveStake := new(big.Int)

	// Calculate stake cap: capPercentage * totalStake
	stakeCap := new(big.Int).Mul(iterationTotalStake, big.NewInt(int64(capPercentage*100)))
	stakeCap = stakeCap.Div(stakeCap, big.NewInt(100))

	for _, opInfo := range operators {
		// Determine which stake to use as input for this iteration
		var stakeToUse *big.Int
		if numIterations == 1 || opInfo.EffectiveStake == nil {
			stakeToUse = opInfo.Stake
		} else {
			stakeToUse = opInfo.EffectiveStake
		}

		// Apply stake cap to this operator
		effectiveStake, err := applyStakeCap(stakeToUse, stakeCap)
		if err != nil {
			return nil, err
		}

		// Update operator's effective stake
		opInfo.EffectiveStake = effectiveStake

		// Add to the new total
		newTotalEffectiveStake = new(big.Int).Add(newTotalEffectiveStake, effectiveStake)
	}

	// Recursive call for the next iteration, using the new total
	return calculateOperatorEffectiveStake(operators, newTotalEffectiveStake, capPercentage, numIterations-1)
}
