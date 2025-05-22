package v2

import (
	"fmt"
	"sort"

	"github.com/Layr-Labs/eigenda/core"
)

func getOrderedOperators(state *core.OperatorState, quorum core.QuorumID) ([]core.OperatorID, map[core.OperatorID]*core.OperatorInfo, error) {

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

func GetAssignmentsForQuorum(state *core.OperatorState, blobParams *core.BlobVersionParameters, quorum core.QuorumID) (map[core.OperatorID]*Assignment, error) {

	orderedOps, operators, err := getOrderedOperators(state, quorum)
	if err != nil {
		return nil, fmt.Errorf("failed to get ordered operators for quorum %d: %w", quorum, err)
	}

	numOperators := len(orderedOps)

	// if numOperators > blobParams.MaxNumOperators {
	// 	return nil, fmt.Errorf("too many operators for quorum %d", quorum)
	// }

	effectiveNumChunks := blobParams.NumChunks - uint32(numOperators)

	total, ok := state.Totals[quorum]
	if !ok {
		return nil, fmt.Errorf("no total found for quorum %d", quorum)
	}

	assignments := make(map[core.OperatorID]*Assignment, len(operators))

	offset := uint32(0)

	totalChunks := 0
	for _, id := range orderedOps {

		if _, ok := operators[id]; !ok {
			continue
		}

		chunksForOperator := core.RoundUpDivide(uint64(effectiveNumChunks)*operators[id].Stake.Uint64(), total.Stake.Uint64())

		totalChunks += int(chunksForOperator)

		assignments[id] = &Assignment{
			Indices: make([]uint32, chunksForOperator),
		}

		for j := range assignments[id].Indices {
			assignments[id].Indices[j] = offset
			offset++
		}

	}

	return assignments, nil
}

func AddAssignmentsForQuorum(assignments map[core.OperatorID]*Assignment, state *core.OperatorState, blobParams *core.BlobVersionParameters, quorum core.QuorumID) (map[core.OperatorID]*Assignment, error) {

	dummyAssignments, err := GetAssignmentsForQuorum(state, blobParams, quorum)
	if err != nil {
		return nil, fmt.Errorf("failed to get assignments for quorum %d: %w", quorum, err)
	}

	orderedOps, _, err := getOrderedOperators(state, quorum)
	if err != nil {
		return nil, fmt.Errorf("failed to get ordered operators for quorum %d: %w", quorum, err)
	}

	usedIndices := make(map[uint32]struct{})

	newAssignments := make(map[core.OperatorID]*Assignment)

	for _, id := range orderedOps {
		newAssignmentIndicesCount := len(dummyAssignments[id].Indices)

		if newAssignmentIndicesCount > len(assignments[id].Indices) {
			newAssignmentIndicesCount = len(assignments[id].Indices)
		}

		newAssignments[id] = &Assignment{
			Indices: assignments[id].Indices[:newAssignmentIndicesCount],
		}

		if newAssignmentIndicesCount < len(assignments[id].Indices) {
			excessIndices := assignments[id].Indices[newAssignmentIndicesCount:]
			for _, index := range excessIndices {
				usedIndices[index] = struct{}{}
			}
		}
	}

	availableIndices := make([]uint32, 0)
	for i := uint32(0); i < blobParams.NumChunks; i++ {
		if _, ok := usedIndices[i]; !ok {
			availableIndices = append(availableIndices, i)
		}
	}

	for _, id := range orderedOps {

		newAssignmentIndicesCount := len(dummyAssignments[id].Indices)
		if newAssignmentIndicesCount > len(newAssignments[id].Indices) {

			// Add available indices to new assignments
			newAssignments[id].Indices = append(newAssignments[id].Indices, availableIndices[:newAssignmentIndicesCount]...)

			// Remove used indices from available indices
			availableIndices = availableIndices[newAssignmentIndicesCount:]
		}
	}

	return newAssignments, nil
}

func MergeAssignmentsAndCap(assignments []map[core.OperatorID]*Assignment, blobParams *core.BlobVersionParameters) map[core.OperatorID]Assignment {

	_mergedAssignments := make(map[core.OperatorID]*Assignment)
	indexMaps := make(map[core.OperatorID]map[uint32]struct{})

	maxChunks := blobParams.NumChunks / blobParams.CodingRate

	for _, assignment := range assignments {
		for id, a := range assignment {

			if _, ok := _mergedAssignments[id]; !ok {
				// Take all indices if less than maxChunks, otherwise take only maxChunks
				indicesLen := uint32(len(a.Indices))
				if indicesLen > maxChunks {
					indicesLen = maxChunks
				}

				_mergedAssignments[id] = &Assignment{
					Indices: a.Indices[:indicesLen],
				}
				indexMaps[id] = make(map[uint32]struct{})
				for _, index := range a.Indices[:indicesLen] {
					indexMaps[id][index] = struct{}{}
				}
				continue
			}

			for _, index := range a.Indices {

				if uint32(len(_mergedAssignments[id].Indices)) >= maxChunks {
					break
				}

				if _, ok := indexMaps[id][index]; ok {
					continue
				}
				_mergedAssignments[id].Indices = append(_mergedAssignments[id].Indices, index)
				indexMaps[id][index] = struct{}{}
			}
		}
	}

	mergedAssignments := make(map[core.OperatorID]Assignment)
	for id, a := range _mergedAssignments {
		mergedAssignments[id] = Assignment{
			Indices: a.Indices,
		}
	}

	return mergedAssignments
}

// GetAssignments calculates chunk assignments for operators in a quorum based on their stake
func GetAssignments(state *core.OperatorState, blobParams *core.BlobVersionParameters, quorums []core.QuorumID, blobKey []byte) (map[core.OperatorID]Assignment, error) {
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
			assignments, err := GetAssignmentsForQuorum(state, blobParams, q)
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

// GetAssignment returns the assignment for a specific operator
func GetAssignment(state *core.OperatorState, blobParams *core.BlobVersionParameters, quorums []core.QuorumID, blobKey []byte, id core.OperatorID) (Assignment, error) {

	if blobParams == nil {
		return Assignment{}, fmt.Errorf("blob params cannot be nil")
	}

	assignments, err := GetAssignments(state, blobParams, quorums, blobKey)
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
