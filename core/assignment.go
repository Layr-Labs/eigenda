package core

import (
	"errors"
	"fmt"
	"math"
	"math/big"
)

const (
	percentMultiplier    = 100
	minChunkLength       = 1
	maxRequiredNumChunks = 8192
)

var (
	ErrInvalidChunkLength  = errors.New("invalid chunk length")
	ErrChunkLengthTooSmall = errors.New("chunk length too small")
	ErrChunkLengthTooLarge = errors.New("chunk length too large")
	ErrNotFound            = errors.New("not found")
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

// Implementation

// AssignmentCoordinator is responsible for taking the current OperatorState and the security requirements represented by a
// given  QuorumResults and determining or validating system parameters that will satisfy these security requirements given the
// OperatorStates. There are two classes of parameters that must be determined or validate: 1) the chunk indices that will be
// assigned to each DA node, and 2) the size of each chunk.
type AssignmentCoordinator interface {

	// GetAssignments calculates the full set of node assignments. The assignment of indices to nodes depends only on the OperatorState
	// for a given quorum and the quantizationFactor. In particular, it does not depend on the security parameters.
	GetAssignments(state *OperatorState, blobLength uint, info *BlobQuorumInfo) (map[OperatorID]Assignment, AssignmentInfo, error)

	// GetOperatorAssignment calculates the assignment for a specific DA node
	GetOperatorAssignment(state *OperatorState, header *BlobHeader, quorum QuorumID, id OperatorID) (Assignment, AssignmentInfo, error)

	// GetMinimumChunkLength calculates the minimum chunkSize that is sufficient for a given blob for each quorum
	ValidateChunkLength(state *OperatorState, header *BlobHeader, quorum QuorumID) (bool, error)

	CalculateChunkLength(state *OperatorState, blobLength uint, param *SecurityParam) (uint, error)
}

type StdAssignmentCoordinator struct {
}

var _ AssignmentCoordinator = (*StdAssignmentCoordinator)(nil)

func (c *StdAssignmentCoordinator) GetAssignments(state *OperatorState, blobLength uint, info *BlobQuorumInfo) (map[OperatorID]Assignment, AssignmentInfo, error) {

	quorum := info.QuorumID

	numOperators := len(state.Operators[quorum])
	chunksByOperator := make([]uint, numOperators)

	// Get NumPar
	numChunks := uint(0)
	totalStakes := state.Totals[quorum].Stake
	for _, r := range state.Operators[quorum] {

		// m_i = ceil( B*S_i / C \gamma \sum_{j=1}^N S_j )
		num := new(big.Int).Mul(big.NewInt(int64(blobLength*percentMultiplier)), r.Stake)

		gammaChunkLength := big.NewInt(int64(info.ChunkLength) * int64((info.QuorumThreshold - info.AdversaryThreshold)))
		denom := new(big.Int).Mul(gammaChunkLength, totalStakes)
		m := roundUpDivideBig(num, denom)

		numChunks += uint(m.Uint64())
		chunksByOperator[r.Index] = uint(m.Uint64())
	}

	currentIndex := uint(0)
	assignments := make([]Assignment, numOperators)
	for operatorInd := range chunksByOperator {

		// Find the operator that should be at index currentIndex
		m := chunksByOperator[operatorInd]
		assignments[operatorInd] = Assignment{
			StartIndex: currentIndex,
			NumChunks:  m,
		}
		currentIndex += m
	}

	assignmentMap := make(map[OperatorID]Assignment)

	for id, opInfo := range state.Operators[quorum] {
		assignment := assignments[opInfo.Index]
		assignmentMap[id] = assignment
	}

	return assignmentMap, AssignmentInfo{
		TotalChunks: numChunks,
	}, nil

}

func (c *StdAssignmentCoordinator) GetOperatorAssignment(state *OperatorState, header *BlobHeader, quorum QuorumID, id OperatorID) (Assignment, AssignmentInfo, error) {

	assignments, info, err := c.GetAssignments(state, header.Length, header.QuorumInfos[quorum])
	if err != nil {
		return Assignment{}, AssignmentInfo{}, err
	}

	assignment, ok := assignments[id]
	if !ok {
		return Assignment{}, AssignmentInfo{}, ErrNotFound
	}

	return assignment, info, nil
}

func (c *StdAssignmentCoordinator) ValidateChunkLength(state *OperatorState, header *BlobHeader, quorum QuorumID) (bool, error) {

	info := header.QuorumInfos[quorum]

	// Check that the chunk length meets the minimum requirement
	if info.ChunkLength < minChunkLength {
		return false, ErrChunkLengthTooSmall
	}

	// Get minimum stake amont
	minStake := state.Totals[quorum].Stake
	for _, r := range state.Operators[quorum] {
		if r.Stake.Cmp(minStake) < 0 {
			minStake = r.Stake
		}
	}

	totalStake := state.Totals[quorum].Stake
	if info.ChunkLength != minChunkLength {

		num := new(big.Int).Mul(big.NewInt(2*int64(header.Length*percentMultiplier)), minStake)
		denom := new(big.Int).Mul(big.NewInt(int64(info.QuorumThreshold-info.AdversaryThreshold)), totalStake)
		maxChunkLength := uint(roundUpDivideBig(num, denom).Uint64())

		maxChunkLength2 := roundUpDivide(2*header.Length*percentMultiplier, maxRequiredNumChunks*uint(info.QuorumThreshold-info.AdversaryThreshold))

		if maxChunkLength < maxChunkLength2 {
			maxChunkLength = maxChunkLength2
		}

		maxChunkLength = uint(nextPowerOf2(uint64(maxChunkLength)))

		if info.ChunkLength > maxChunkLength {
			fmt.Println("maxChunkLength", maxChunkLength, "info.ChunkLength", info.ChunkLength)
			return false, ErrChunkLengthTooLarge
		}

	}

	return true, nil

}

func (c *StdAssignmentCoordinator) CalculateChunkLength(state *OperatorState, blobLength uint, param *SecurityParam) (uint, error) {

	quorum := param.QuorumID
	numOperators := len(state.Operators[quorum])

	// Get minimum stake amont
	minStake := state.Totals[quorum].Stake
	for _, r := range state.Operators[quorum] {
		if r.Stake.Cmp(minStake) < 0 {
			minStake = r.Stake
		}
	}

	totalStake := state.Totals[quorum].Stake

	numChunks := int(roundUpDivideBig(totalStake, minStake).Uint64())

	if numChunks > maxRequiredNumChunks-numOperators {
		numChunks = maxRequiredNumChunks - numOperators
	}

	chunkLength := roundUpDivide(blobLength*percentMultiplier, uint(numChunks)*uint(param.QuorumThreshold-param.AdversaryThreshold))

	if chunkLength < minChunkLength {
		chunkLength = minChunkLength
	}

	chunkLength = uint(nextPowerOf2(uint64(chunkLength)))

	return chunkLength, nil

}

func roundUpDivideBig(a, b *big.Int) *big.Int {

	one := new(big.Int).SetUint64(1)

	res := new(big.Int).Div(new(big.Int).Sub(new(big.Int).Add(a, b), one), b)
	return res

}

func roundUpDivide(a, b uint) uint {
	return (a + b - 1) / b

}

func nextPowerOf2(d uint64) uint64 {
	nextPower := math.Ceil(math.Log2(float64(d)))
	return uint64(math.Pow(2.0, nextPower))
}
