package v2

import (
	"fmt"
	"math"
	"math/big"
	"math/rand"

	"encoding/binary"

	"github.com/Layr-Labs/eigenda/core"
)

type SampleConfig struct {
	NumUnits       uint32
	SamplesPerUnit uint32
}

var DefaultSampleConfig = SampleConfig{
	NumUnits:       393,
	SamplesPerUnit: 20,
}

// getMaxStakePercentage calculates the maximum stake percentage across specified quorums for a specific operator. The function assumes that it will be passed an operator state that contains the quorums for which the operator is a member. Therefore, it does not return an error if a quorum is not found in the state.
func getMaxStakePercentage(state *core.OperatorState, blobParams *core.BlobVersionParameters, quorums []core.QuorumID, id core.OperatorID) (*big.Float, error) {
	maxStakePercentage := new(big.Float)
	found := false

	for _, q := range quorums {
		total, ok := state.Totals[q]
		if !ok {
			continue
		}

		ops, ok := state.Operators[q]
		if !ok || len(ops) == 0 {
			return nil, fmt.Errorf("no operators found for quorum %d", q)
		}

		numOps := len(ops)
		if uint32(numOps) > blobParams.MaxNumOperators {
			return nil, fmt.Errorf("too many operators (%d) to get assignments: max number of operators is %d", numOps, blobParams.MaxNumOperators)
		}

		// Get the stake for this operator in this quorum
		stake, ok := ops[id]
		if !ok {
			continue // Skip if operator is not in this quorum
		}

		found = true
		stakePercentage := new(big.Float).Quo(
			new(big.Float).SetInt(stake.Stake),
			new(big.Float).SetInt(total.Stake),
		)
		if maxStakePercentage.Cmp(stakePercentage) < 0 {
			maxStakePercentage = stakePercentage
		}
	}

	if !found {
		return nil, ErrNotFound
	}

	return maxStakePercentage, nil
}

// getRelevantOperators gets all of the operators that are relevant to the blob
func getRelevantOperators(state *core.OperatorState, quorums []core.QuorumID) (map[core.OperatorID]struct{}, error) {
	if state == nil {
		return nil, fmt.Errorf("state cannot be nil")
	}

	operators := make(map[core.OperatorID]struct{}, 0)
	for _, q := range quorums {
		ops, ok := state.Operators[q]
		if !ok || len(ops) == 0 {
			return nil, fmt.Errorf("no operators found for quorum %d", q)
		}

		for id := range ops {
			operators[id] = struct{}{}
		}
	}
	return operators, nil
}

// GetAssignments calculates chunk assignments for operators in a quorum based on their stake
func GetAssignments(state *core.OperatorState, blobParams *core.BlobVersionParameters, quorums []core.QuorumID, blobKey []byte) (map[core.OperatorID]Assignment, error) {
	if state == nil {
		return nil, fmt.Errorf("state cannot be nil")
	}

	if blobParams == nil {
		return nil, fmt.Errorf("blob params cannot be nil")
	}

	ops, err := getRelevantOperators(state, quorums)
	if err != nil {
		return nil, fmt.Errorf("failed to get relevant operators: %w", err)
	}

	numOps := len(ops)

	// Early return for empty operator set
	if numOps == 0 {
		return make(map[core.OperatorID]Assignment), nil
	}

	assignments := make(map[core.OperatorID]Assignment, numOps)
	for id := range ops {
		assignment, err := GetAssignment(state, blobParams, quorums, blobKey, id)
		if err != nil {
			return nil, fmt.Errorf("failed to get assignment for operator %v: %w", id, err)
		}
		assignments[id] = assignment
	}

	return assignments, nil
}

// GetAssignment returns the assignment for a specific operator
func GetAssignment(state *core.OperatorState, blobParams *core.BlobVersionParameters, quorums []core.QuorumID, blobKey []byte, id core.OperatorID) (Assignment, error) {
	if state == nil {
		return Assignment{}, fmt.Errorf("state cannot be nil")
	}

	if blobParams == nil {
		return Assignment{}, fmt.Errorf("blob params cannot be nil")
	}

	maxStakePercentage, err := getMaxStakePercentage(state, blobParams, quorums, id)
	if err != nil {
		return Assignment{}, fmt.Errorf("failed to get max stake percentage: %w", err)
	}

	if maxStakePercentage.Cmp(big.NewFloat(0)) == 0 {
		return Assignment{}, fmt.Errorf("max stake percentage is zero")
	}

	// Calculate number of samples based on max stake percentage
	maxStakeFloat, _ := maxStakePercentage.Float64()
	numSamples := uint32(math.Ceil(maxStakeFloat * float64(DefaultSampleConfig.NumUnits*DefaultSampleConfig.SamplesPerUnit)))

	// Create a deterministic random number generator using the blob key as seed
	// We also mix in the operator ID to ensure different operators get different assignments
	seed := make([]byte, 40) // 32 bytes for blob key + 8 bytes for operator ID
	copy(seed[:32], blobKey[:])
	copy(seed[32:], id[:])
	rng := rand.New(rand.NewSource(int64(binary.BigEndian.Uint64(seed[:8]))))

	// Generate random indices without replacement for this operator
	indices := make([]uint32, numSamples)
	available := make([]uint32, blobParams.NumChunks)
	for i := uint32(0); i < blobParams.NumChunks; i++ {
		available[i] = i
	}

	// Select numSamples indices randomly without replacement
	for i := uint32(0); i < numSamples && i < blobParams.NumChunks; i++ {
		// Pick a random index from remaining available indices
		j := rng.Intn(len(available))
		indices[i] = available[j]
		// Remove the selected index by swapping with the last element and reducing slice length
		available[j] = available[len(available)-1]
		available = available[:len(available)-1]
	}

	return Assignment{
		Indices: indices,
	}, nil
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
