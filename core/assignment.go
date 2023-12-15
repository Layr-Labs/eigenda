package core

import (
	"errors"
	"math/big"
)

// Assignment

type OperatorID = [32]byte

type OperatorIndex = uint

type ChunkNumber = uint

// AssignmentInfo contains the global information associated with a group of assignments, such as the total number of chunks
type AssignmentInfo struct {
	TotalChunks ChunkNumber
}

// Assignment contains information about the set of chunks that a specific node will receive
type Assignment struct {
	StartIndex ChunkNumber
	NumChunks  ChunkNumber
}

// GetIndices generates the list of ChunkIndices associated with a given assignment
func (c *Assignment) GetIndices() []ChunkNumber {
	indices := make([]ChunkNumber, c.NumChunks)
	for ind := range indices {
		indices[ind] = c.StartIndex + ChunkNumber(ind)
	}
	return indices
}

// AssignmentCoordinator is responsible for taking the current OperatorState and the security requirements represented by a
// given  QuorumResults and determining or validating system parameters that will satisfy these security requirements given the
// OperatorStates. There are two classes of parameters that must be determined or validate: 1) the chunk indices that will be
// assigned to each DA node, and 2) the size of each chunk.
type AssignmentCoordinatorOld interface {

	// GetAssignments calculates the full set of node assignments. The assignment of indices to nodes depends only on the OperatorState
	// for a given quorum and the quantizationFactor. In particular, it does not depend on the security parameters.
	GetAssignments(state *OperatorState, quorumID QuorumID, quantizationFactor uint) (map[OperatorID]Assignment, AssignmentInfo, error)

	// GetOperatorAssignment calculates the assignment for a specific DA node
	GetOperatorAssignment(state *OperatorState, quorum QuorumID, quantizationFactor uint, id OperatorID) (Assignment, AssignmentInfo, error)

	// GetMinimumChunkLength calculates the minimum chunkSize that is sufficient for a given blob for each quorum
	GetMinimumChunkLength(numOperators, blobLength, quantizationFactor uint, quorumThreshold, adversaryThreshold uint8) (uint, error)

	// GetChunkLengthFromHeader calculates the chunk length from the blob header
	GetChunkLengthFromHeader(state *OperatorState, header *BlobQuorumInfo) (uint, error)
}

// Implementation

const PercentMultiplier = 100

var (
	ErrNotFound = errors.New("not found")
)

// type StdAssignmentCoordinatorOld struct {
// }

// var _ AssignmentCoordinatorOld = (*StdAssignmentCoordinatorOld)(nil)

// func (c *StdAssignmentCoordinatorOld) GetAssignments(state *OperatorState, quorum QuorumID, quantizationFactor uint) (map[OperatorID]Assignment, AssignmentInfo, error) {

// 	numOperators := len(state.Operators[quorum])
// 	numOperatorsBig := new(big.Int).SetUint64(uint64(numOperators))

// 	quantizationFactorBig := new(big.Int).SetUint64(uint64(quantizationFactor))

// 	chunksByOperator := make([]uint, numOperators)

// 	// Get NumPar
// 	numChunks := uint(0)
// 	totalStakes := state.Totals[quorum].Stake
// 	for _, r := range state.Operators[quorum] {

// 		m := new(big.Int).Mul(numOperatorsBig, r.Stake)
// 		m = m.Mul(m, quantizationFactorBig)
// 		m = roundUpDivideBig(m, totalStakes)

// 		numChunks += uint(m.Uint64())
// 		chunksByOperator[r.Index] = uint(m.Uint64())
// 	}

// 	currentIndex := uint(0)
// 	assignments := make([]Assignment, numOperators)

// 	headerHash := [32]byte{}

// 	for orderedInd := range chunksByOperator {

// 		// Find the operator that should be at index currentIndex
// 		operatorInd := getOperatorAtIndex(headerHash, orderedInd, numOperators)
// 		m := chunksByOperator[operatorInd]

// 		assignments[operatorInd] = Assignment{
// 			StartIndex: currentIndex,
// 			NumChunks:  m,
// 		}
// 		currentIndex += m
// 	}

// 	assignmentMap := make(map[OperatorID]Assignment)

// 	for id, opInfo := range state.Operators[quorum] {
// 		assignment := assignments[opInfo.Index]
// 		assignmentMap[id] = assignment
// 	}

// 	return assignmentMap, AssignmentInfo{
// 		TotalChunks: numChunks,
// 	}, nil

// }

// // getOperatorAtIndex returns the operator at a given index within the reordered sequence.
// // We reorder the sequence by letting the reordered_index = operator_index + headerHash.
// // Thus, get get the operator at a given reordered_index, we simply reverse:
// // operator_index = reordered_index - headerHash
// func getOperatorAtIndex(headerHash [32]byte, index, numOperators int) int {
// 	indexBig := new(big.Int).SetUint64(uint64(index))
// 	offset := new(big.Int).SetBytes(headerHash[:])

// 	operatorIndex := new(big.Int).Sub(indexBig, offset)

// 	operatorIndex.Mod(operatorIndex, new(big.Int).SetUint64(uint64(numOperators)))

// 	return int(operatorIndex.Uint64())
// }

// func (c *StdAssignmentCoordinatorOld) GetOperatorAssignment(state *OperatorState, quorum QuorumID, quantizationFactor uint, id OperatorID) (Assignment, AssignmentInfo, error) {

// 	assignments, info, err := c.GetAssignments(state, quorum, quantizationFactor)
// 	if err != nil {
// 		return Assignment{}, AssignmentInfo{}, err
// 	}

// 	assignment, ok := assignments[id]
// 	if !ok {
// 		return Assignment{}, AssignmentInfo{}, ErrNotFound
// 	}

// 	return assignment, info, nil
// }

// func (c *StdAssignmentCoordinatorOld) GetMinimumChunkLength(numOperators, blobLength, quantizationFactor uint, quorumThreshold, adversaryThreshold uint8) (uint, error) {

// 	if adversaryThreshold >= quorumThreshold {
// 		return 0, errors.New("invalid header: quorum threshold does not exceed adversary threshold")
// 	}

// 	numSys := roundUpDivide(uint(quorumThreshold-adversaryThreshold)*numOperators*quantizationFactor, PercentMultiplier)
// 	chunkLength := roundUpDivide(blobLength, numSys)
// 	return chunkLength, nil

// }

// func (c *StdAssignmentCoordinatorOld) GetChunkLengthFromHeader(state *OperatorState, header *BlobQuorumInfo) (uint, error) {

// 	// Validate the chunk length
// 	numOperators := uint(len(state.Operators[header.QuorumID]))
// 	chunkLength := header.EncodedBlobLength / (header.QuantizationFactor * numOperators)

// 	if chunkLength*header.QuantizationFactor*numOperators != header.EncodedBlobLength {
// 		return 0, errors.New("invalid header")
// 	}

// 	return chunkLength, nil
// }

func roundUpDivideBig(a, b *big.Int) *big.Int {

	one := new(big.Int).SetUint64(1)
	res := new(big.Int)
	a.Add(a, b)
	a.Sub(a, one)
	res.Div(a, b)
	return res

}

func roundUpDivide(a, b uint) uint {
	return (a + b - 1) / b

}
