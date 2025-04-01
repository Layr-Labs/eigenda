package v2

import (
	"fmt"
	"math/big"
	"sort"

	"github.com/Layr-Labs/eigenda/core"
)

// GetAssignments calculates chunk assignments for operators in a quorum based on their stake
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

	type operatorAssignment struct {
		id     core.OperatorID
		index  uint32
		chunks uint32
		stake  *big.Int
	}

	numOperatorsBig := big.NewInt(int64(numOps))
	numChunksBig := big.NewInt(int64(blobParams.NumChunks))
	totalStake := state.Totals[quorum].Stake

	// Calculate number of chunks - numOperators once and reuse
	diffChunksOps := new(big.Int).Sub(numChunksBig, numOperatorsBig)
	chunkAssignments := make([]operatorAssignment, 0, numOps)
	// Calculate initial chunk assignments based on stake
	totalCalculatedChunks := uint32(0)
	for ID, r := range ops {
		// Calculate chunks for this operator: (stake * (numChunks - numOperators)) / totalStake (rounded up)
		num := new(big.Int).Mul(r.Stake, diffChunksOps)
		chunks := uint32(core.RoundUpDivideBig(num, totalStake).Uint64())

		chunkAssignments = append(chunkAssignments, operatorAssignment{
			id:     ID,
			index:  uint32(r.Index),
			chunks: chunks,
			stake:  r.Stake,
		})

		totalCalculatedChunks += chunks
	}

	// Sort by stake (decreasing) with index as tie-breaker
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
		// Add remaining chunks to operators with highest stake first
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
