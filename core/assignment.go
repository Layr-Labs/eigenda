package core

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"
)

const (
	percentMultiplier = 100

	// minChunkLength is the minimum chunk length supported. Generally speaking, it doesn't make sense for a chunk to be
	// smaller than the proof overhead, which is equal to one G1 point.
	MinChunkLength = 1

	// maxRequiredNumChunks is the maximum number of chunks that can be required for a single quorum. Encoding costs scale
	// as N*log(N), with N being the number of chunks. The value of 8192 was chosen to ensure that the encoding costs for
	// a single quorum are reasonable, while still allowing for a single operator to have O(0.01%) of the total data.
	MaxRequiredNumChunks = 8192
)

var (
	ErrChunkLengthTooSmall = errors.New("chunk length too small")
	ErrChunkLengthTooLarge = errors.New("chunk length too large")
	ErrNotFound            = errors.New("not found")
)

// Assignment

type OperatorID [32]byte

func (id OperatorID) Hex() string {
	return hex.EncodeToString(id[:])
}

// The "s" is an operatorId in hex string format, which may or may not have the "0x" prefix.
func OperatorIDFromHex(s string) (OperatorID, error) {
	opID := [32]byte{}
	s = strings.TrimPrefix(s, "0x")
	if len(s) != 64 {
		return OperatorID(opID), errors.New("operatorID hex string must be 64 bytes, or 66 bytes if starting with 0x")
	}
	opIDslice, err := hex.DecodeString(s)
	if err != nil {
		return OperatorID(opID), err
	}
	copy(opID[:], opIDslice)
	return OperatorID(opID), nil
}

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
// given QuorumResults and determining or validating system parameters that will satisfy these security requirements given the
// OperatorStates. There are two classes of parameters that must be determined or validated: 1) the chunk indices that will be
// assigned to each DA node, and 2) the length of each chunk.
type AssignmentCoordinator interface {

	// GetAssignments calculates the full set of node assignments.
	GetAssignments(state *OperatorState, blobLength uint, info *BlobQuorumInfo) (map[OperatorID]Assignment, AssignmentInfo, error)

	// GetOperatorAssignment calculates the assignment for a specific DA node
	GetOperatorAssignment(state *OperatorState, header *BlobHeader, quorum QuorumID, id OperatorID) (Assignment, AssignmentInfo, error)

	// ValidateChunkLength validates that the chunk length for the given quorum satisfies all protocol constraints
	ValidateChunkLength(state *OperatorState, blobLength uint, info *BlobQuorumInfo) (bool, error)

	// CalculateChunkLength will find the max chunk length (as a power of 2) which satisfies the protocol constraints. If
	// targetNumChunks is non-zero, then CalculateChunkLength will return the smaller of 1) the smallest chunk length which
	// results in a number of chunks less than or equal to targetNumChunks and 2) the largest chunk length which satisfies
	// the protocol constraints.
	CalculateChunkLength(state *OperatorState, blobLength, targetNumChunks uint, param *SecurityParam) (uint, error)
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

		gammaChunkLength := big.NewInt(int64(info.ChunkLength) * int64((info.ConfirmationThreshold - info.AdversaryThreshold)))
		if gammaChunkLength.Cmp(big.NewInt(0)) <= 0 {
			return nil, AssignmentInfo{}, fmt.Errorf("gammaChunkLength must be greater than 0")
		}
		if totalStakes.Cmp(big.NewInt(0)) == 0 {
			return nil, AssignmentInfo{}, fmt.Errorf("total stake in quorum %d must be greater than 0", quorum)
		}
		denom := new(big.Int).Mul(gammaChunkLength, totalStakes)
		if denom.Cmp(big.NewInt(0)) == 0 {
			return nil, AssignmentInfo{}, fmt.Errorf("gammaChunkLength %d and total stake %d in quorum %d must be greater than 0", gammaChunkLength, totalStakes, quorum)
		}
		m := RoundUpDivideBig(num, denom)

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

	quorumInfo := header.GetQuorumInfo(quorum)
	if quorumInfo == nil {
		return Assignment{}, AssignmentInfo{}, fmt.Errorf("invalid request: quorum ID %d not found in blob header", quorum)
	}

	assignments, info, err := c.GetAssignments(state, header.Length, quorumInfo)
	if err != nil {
		return Assignment{}, AssignmentInfo{}, err
	}

	assignment, ok := assignments[id]
	if !ok {
		return Assignment{}, AssignmentInfo{}, ErrNotFound
	}

	return assignment, info, nil
}

func (c *StdAssignmentCoordinator) ValidateChunkLength(state *OperatorState, blobLength uint, info *BlobQuorumInfo) (bool, error) {

	// Check that the chunk length meets the minimum requirement
	if info.ChunkLength < MinChunkLength {
		return false, fmt.Errorf("%w: chunk length: %d, min chunk length: %d", ErrChunkLengthTooSmall, info.ChunkLength, MinChunkLength)
	}

	// Get minimum stake amont
	minStake := state.Totals[info.QuorumID].Stake
	for _, r := range state.Operators[info.QuorumID] {
		if r.Stake.Cmp(minStake) < 0 {
			minStake = r.Stake
		}
	}

	totalStake := state.Totals[info.QuorumID].Stake
	if info.ChunkLength != MinChunkLength {
		if totalStake.Cmp(big.NewInt(0)) == 0 {
			return false, fmt.Errorf("total stake in quorum %d must be greater than 0", info.QuorumID)
		}
		num := new(big.Int).Mul(big.NewInt(2*int64(blobLength*percentMultiplier)), minStake)
		denom := new(big.Int).Mul(big.NewInt(int64(info.ConfirmationThreshold-info.AdversaryThreshold)), totalStake)
		maxChunkLength := uint(RoundUpDivideBig(num, denom).Uint64())

		maxChunkLength2 := RoundUpDivide(2*blobLength*percentMultiplier, MaxRequiredNumChunks*uint(info.ConfirmationThreshold-info.AdversaryThreshold))

		if maxChunkLength < maxChunkLength2 {
			maxChunkLength = maxChunkLength2
		}

		maxChunkLength = uint(NextPowerOf2(maxChunkLength))

		if info.ChunkLength > maxChunkLength {
			return false, fmt.Errorf("%w: chunk length: %d, max chunk length: %d", ErrChunkLengthTooLarge, info.ChunkLength, maxChunkLength)
		}

	}

	return true, nil

}

// CalculateChunkLength will find the max chunk length (as a power of 2) which satisfies the protocol constraints. It does this by
// doubling the chunk length (multiplicative binary search) until it is too large or we are beneath the targetNumChunks.
// This will always give the largest acceptable chunk length. The loop will always stop because the chunk length will eventually be
// too large for the constraint in ValidateChunkLength
func (c *StdAssignmentCoordinator) CalculateChunkLength(state *OperatorState, blobLength, targetNumChunks uint, param *SecurityParam) (uint, error) {

	chunkLength := uint(MinChunkLength) * 2

	for {

		quorumInfo := &BlobQuorumInfo{
			SecurityParam: *param,
			ChunkLength:   chunkLength,
		}

		ok, err := c.ValidateChunkLength(state, blobLength, quorumInfo)
		if err != nil || !ok {
			return chunkLength / 2, nil
		}

		if targetNumChunks != 0 {

			_, info, err := c.GetAssignments(state, blobLength, quorumInfo)
			if err != nil {
				return 0, err
			}

			if info.TotalChunks <= targetNumChunks {
				return chunkLength, nil
			}
		}

		chunkLength *= 2

	}

}
