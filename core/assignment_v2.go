package core

import (
	"fmt"
	"math/big"
	"sort"
)

type BlobParameters struct {
	CodingRate              uint
	ReconstructionThreshold float64
	NumChunks               uint
}

var (
	ParametersMap = map[byte]BlobParameters{
		0: {CodingRate: 8, ReconstructionThreshold: 0.22, NumChunks: 8192},
	}
)

// Implementation

// AssignmentCoordinator is responsible for taking the current OperatorState and the security requirements represented by a
// given QuorumResults and determining or validating system parameters that will satisfy these security requirements given the
// OperatorStates. There are two classes of parameters that must be determined or validated: 1) the chunk indices that will be
// assigned to each DA node, and 2) the length of each chunk.
type AssignmentCoordinatorV2 interface {

	// GetAssignments calculates the full set of node assignments.
	GetAssignments(state *OperatorState, blobVersion byte, quorum QuorumID) (map[OperatorID]Assignment, error)

	// GetOperatorAssignment calculates the assignment for a specific DA node
	GetAssignment(state *OperatorState, blobVersion byte, quorum QuorumID, id OperatorID) (Assignment, error)

	// ValidateChunkLength validates that the chunk length for the given quorum satisfies all protocol constraints
	GetChunkLength(blobVersion byte, blobLength uint) (uint, error)
}

type StdAssignmentCoordinatorV2 struct {
}

var _ AssignmentCoordinatorV2 = (*StdAssignmentCoordinatorV2)(nil)

func (c *StdAssignmentCoordinatorV2) GetAssignments(state *OperatorState, blobVersion byte, quorum QuorumID) (map[OperatorID]Assignment, error) {

	params := ParametersMap[blobVersion]

	n := big.NewInt(int64(len(state.Operators[quorum])))
	m := big.NewInt(int64(params.NumChunks))

	type assignment struct {
		id     OperatorID
		index  uint
		chunks uint
		stake  *big.Int
	}

	chunkAssignments := make([]assignment, 0, len(state.Operators[quorum]))
	for ID, r := range state.Operators[quorum] {

		num := new(big.Int).Mul(r.Stake, new(big.Int).Sub(m, n))
		denom := state.Totals[quorum].Stake

		chunks := roundUpDivideBig(num, denom)

		// delta := new(big.Int).Sub(new(big.Int).Mul(r.Stake, m), new(big.Int).Mul(denom, chunks))

		chunkAssignments = append(chunkAssignments, assignment{id: ID, index: r.Index, chunks: uint(chunks.Uint64()), stake: r.Stake})
	}

	// Sort chunk decreasing by stake or operator ID in case of a tie
	sort.Slice(chunkAssignments, func(i, j int) bool {
		if chunkAssignments[i].stake.Cmp(chunkAssignments[j].stake) == 0 {
			return chunkAssignments[i].index < chunkAssignments[j].index
		}
		return chunkAssignments[i].stake.Cmp(chunkAssignments[j].stake) == 1
	})

	mp := 0
	for _, a := range chunkAssignments {
		mp += int(a.chunks)
	}

	delta := int(params.NumChunks) - mp
	if delta < 0 {
		return nil, fmt.Errorf("total chunks %d exceeds maximum %d", mp, params.NumChunks)
	}

	assignments := make(map[OperatorID]Assignment, len(chunkAssignments))
	index := uint(0)
	for i, a := range chunkAssignments {
		if i < delta {
			a.chunks++
		}

		assignment := Assignment{
			StartIndex: index,
			NumChunks:  a.chunks,
		}

		assignments[a.id] = assignment
		index += a.chunks
	}

	return assignments, nil

}

func (c *StdAssignmentCoordinatorV2) GetAssignment(state *OperatorState, blobVersion byte, quorum QuorumID, id OperatorID) (Assignment, error) {

	assignments, err := c.GetAssignments(state, blobVersion, quorum)
	if err != nil {
		return Assignment{}, err
	}

	assignment, ok := assignments[id]
	if !ok {
		return Assignment{}, ErrNotFound
	}

	return assignment, nil
}

func (c *StdAssignmentCoordinatorV2) GetChunkLength(blobVersion byte, blobLength uint) (uint, error) {

	return blobLength * ParametersMap[blobVersion].CodingRate / ParametersMap[blobVersion].NumChunks, nil

}