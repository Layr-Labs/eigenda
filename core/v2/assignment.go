package v2

import (
	"fmt"
	"math/big"
	"sort"

	"github.com/Layr-Labs/eigenda/core"
)

func getOrderedOperators(
	state *core.OperatorState,
	quorum core.QuorumID,
) ([]core.OperatorID, map[core.OperatorID]*core.OperatorInfo, error) {

	if state == nil {
		return nil, nil, fmt.Errorf("state cannot be nil")
	}

	operators, ok := state.Operators[quorum]
	if !ok || len(operators) == 0 {
		return nil, nil, fmt.Errorf("no operators found for quorum %d", quorum)
	}

	orderedOps := make([]core.OperatorID, 0, len(operators))
	for id := range operators {
		orderedOps = append(orderedOps, id)
	}

	sort.Slice(orderedOps, func(i, j int) bool {
		return orderedOps[i].Hex() < orderedOps[j].Hex()
	})

	return orderedOps, operators, nil
}

// GetAssignmentsForQuorum calculates chunk assignments for the validators in a single quorum, independently
// of any other quorums. Not all of the chunks in the encoded blob will be assigned; only enough to satisfy the
// reconstruction threshold for the blob.
func GetAssignmentsForQuorum(
	state *core.OperatorState,
	blobParams *core.BlobVersionParameters,
	quorum core.QuorumID,
) (map[core.OperatorID]*Assignment, []core.OperatorID, error) {

	orderedOps, operators, err := getOrderedOperators(state, quorum)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get ordered operators for quorum %d: %w", quorum, err)
	}

	if len(orderedOps) > int(blobParams.MaxNumOperators) {
		return nil, nil, fmt.Errorf("too many operators for quorum %d", quorum)
	}

	effectiveNumChunks := blobParams.NumChunks - blobParams.MaxNumOperators

	total, ok := state.Totals[quorum]
	if !ok {
		return nil, nil, fmt.Errorf("no total found for quorum %d", quorum)
	}

	assignments := make(map[core.OperatorID]*Assignment, len(operators))

	offset := uint32(0)

	totalChunks := 0
	for _, id := range orderedOps {

		operator, ok := operators[id]
		if !ok {
			return nil, nil, fmt.Errorf("operator %s not found for quorum %d", id, quorum)
		}

		chunksForOperator := uint32(core.RoundUpDivideBig(new(big.Int).Mul(big.NewInt(int64(effectiveNumChunks)), operator.Stake), total.Stake).Uint64())

		totalChunks += int(chunksForOperator)

		assignments[id] = &Assignment{
			Indices: make([]uint32, chunksForOperator),
		}

		for j := range assignments[id].Indices {
			assignments[id].Indices[j] = offset
			offset++
		}

	}

	return assignments, orderedOps, nil
}

// AddAssignmentsForQuorum uses an existing quorum assignment as a baseline and creates a new assignment for a separate
// quorum which maximizes the overlap of the assignments for each validator. This is done through two steps:
// 1. For each validator, as many chunks as possible are taken from the existing assignments for the first quorum,
// 2. Any unused chunks are then distributed among the validators who still need additional chunks to meet their alloted number.
// This has the property that the total number of chunks assigned to an operator across the two quorums will be equal to that
// of the quorum in which it has the largest allocation. (AddAssignmentsForQuorum can be used iteratively with more than two quorums
// in order to maximize overlap, but will not preserve this property.)
func AddAssignmentsForQuorum(
	assignments map[core.OperatorID]*Assignment,
	state *core.OperatorState,
	blobParams *core.BlobVersionParameters,
	quorum core.QuorumID,
) (map[core.OperatorID]*Assignment, error) {

	dummyAssignments, orderedOps, err := GetAssignmentsForQuorum(state, blobParams, quorum)
	if err != nil {
		return nil, fmt.Errorf("failed to get assignments for quorum %d: %w", quorum, err)
	}

	usedIndices := make(map[uint32]struct{})

	newAssignments := make(map[core.OperatorID]*Assignment)

	for _, id := range orderedOps {
		newAssignmentIndicesCount := len(dummyAssignments[id].Indices)

		if _, ok := assignments[id]; !ok {
			newAssignments[id] = &Assignment{
				Indices: make([]uint32, 0, newAssignmentIndicesCount),
			}
			continue
		}

		if newAssignmentIndicesCount > len(assignments[id].Indices) {
			newAssignmentIndicesCount = len(assignments[id].Indices)
		}

		newAssignments[id] = &Assignment{
			Indices: assignments[id].Indices[:newAssignmentIndicesCount],
		}

		for _, index := range newAssignments[id].Indices {
			usedIndices[index] = struct{}{}
		}
	}

	availableIndices := make([]uint32, 0, blobParams.NumChunks)
	for i := uint32(0); i < blobParams.NumChunks; i++ {
		if _, ok := usedIndices[i]; !ok {
			availableIndices = append(availableIndices, i)
		}
	}

	for _, id := range orderedOps {

		newAssignmentIndicesCount := len(dummyAssignments[id].Indices)
		if newAssignmentIndicesCount > len(newAssignments[id].Indices) {

			indicesToAdd := newAssignmentIndicesCount - len(newAssignments[id].Indices)

			// Add available indices to new assignments
			newAssignments[id].Indices = append(newAssignments[id].Indices, availableIndices[:indicesToAdd]...)

			// Remove used indices from available indices
			availableIndices = availableIndices[indicesToAdd:]
		}
	}

	return newAssignments, nil
}

// MergeAssignmentsAndCap merges a list of assignments into a single assignment which contains the union of the
// indices from each of the input assignments. The number of indices for each operator is capped at the maximum
// number of chunks needed to construct a blob. This is because once a validator has enough unique chunks to reconstruct
// a blob, the relationship of these chunk indices to those held by other validators is irrelevant.
func MergeAssignmentsAndCap(
	assignments []map[core.OperatorID]*Assignment,
	blobParams *core.BlobVersionParameters,
) map[core.OperatorID]Assignment {

	mergedAssignments := make(map[core.OperatorID]*Assignment)
	indexMaps := make(map[core.OperatorID]map[uint32]struct{})

	maxChunks := blobParams.NumChunks / blobParams.CodingRate

	for _, assignment := range assignments {
		for id, a := range assignment {

			if _, ok := mergedAssignments[id]; !ok {
				// Take all indices if less than maxChunks, otherwise take only maxChunks
				indicesLen := uint32(len(a.Indices))
				if indicesLen > maxChunks {
					indicesLen = maxChunks
				}

				mergedAssignments[id] = &Assignment{
					Indices: a.Indices[:indicesLen],
				}
				indexMaps[id] = make(map[uint32]struct{})
				for _, index := range a.Indices[:indicesLen] {
					indexMaps[id][index] = struct{}{}
				}
				continue
			}

			for _, index := range a.Indices {

				if uint32(len(mergedAssignments[id].Indices)) >= maxChunks {
					break
				}

				if _, ok := indexMaps[id][index]; ok {
					continue
				}
				mergedAssignments[id].Indices = append(mergedAssignments[id].Indices, index)
				indexMaps[id][index] = struct{}{}
			}
		}
	}

	mergedAssignmentsFinal := make(map[core.OperatorID]Assignment)
	for id, a := range mergedAssignments {
		mergedAssignmentsFinal[id] = Assignment{
			Indices: a.Indices,
		}
	}

	return mergedAssignmentsFinal
}

// GetAssignmentsForBlob calculates chunk assignments for the validators in a set of quorums based on their stake.
// The quorums passed into GetAssignmentsForBlob should be the full set of quorums contained in the blob header.
// Moreover, the OperatorState must include the operator state maps for each of the quorums specified.
// GetAssignmentsForBlob will attempt to construct maximally overlapping assignments for each quorum, and then merge them together.
// The number of chunks assigned to each operator is capped at the maximum number of chunks needed to construct a blob.
func GetAssignmentsForBlob(
	state *core.OperatorState,
	blobParams *core.BlobVersionParameters,
	quorums []core.QuorumID,
) (map[core.OperatorID]Assignment, error) {
	if state == nil {
		return nil, fmt.Errorf("state cannot be nil")
	}

	if blobParams == nil {
		return nil, fmt.Errorf("blob params cannot be nil")
	}

	// Sort quorums
	sort.Slice(quorums, func(i, j int) bool {
		return quorums[i] < quorums[j]
	})

	assignmentsList := make([]map[core.OperatorID]*Assignment, len(quorums))
	for i, q := range quorums {
		if i == 0 {
			assignments, _, err := GetAssignmentsForQuorum(state, blobParams, q)
			if err != nil {
				return nil, fmt.Errorf("failed to get assignments for quorum %d: %w", q, err)
			}
			assignmentsList[i] = assignments
			continue
		}

		assignments, err := AddAssignmentsForQuorum(assignmentsList[0], state, blobParams, q)
		if err != nil {
			return nil, fmt.Errorf("failed to add assignments for quorum %d: %w", q, err)
		}
		assignmentsList[i] = assignments
	}

	mergedAssignments := MergeAssignmentsAndCap(assignmentsList, blobParams)

	return mergedAssignments, nil
}

// GetAssignmentForBlob returns the assignment for a specific operator for a specific blob. The quorums passed into
// GetAssignmentsForBlob should be the full set of quorums contained in the blob header. Moreover, the OperatorState
// must include the operator state maps for each of the quorums specified. GetAssignmentForBlob calls
// GetAssignmentsForBlob under the hood.
func GetAssignmentForBlob(
	state *core.OperatorState,
	blobParams *core.BlobVersionParameters,
	quorums []core.QuorumID,
	id core.OperatorID,
) (Assignment, error) {

	if blobParams == nil {
		return Assignment{}, fmt.Errorf("blob params cannot be nil")
	}

	assignments, err := GetAssignmentsForBlob(state, blobParams, quorums)
	if err != nil {
		return Assignment{}, err
	}

	assignment, ok := assignments[id]
	if !ok {
		return Assignment{}, fmt.Errorf("assignment not found for operator %s", id)
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
