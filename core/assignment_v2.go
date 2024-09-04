package core

import (
	"fmt"
	"math"
	"math/big"
	"sort"
)

type BlobParameters struct {
	CodingRate              uint
	ReconstructionThreshold float64
	NumChunks               uint
}

func (p BlobParameters) MaxNumOperators() uint {

	return uint(math.Floor(float64(p.NumChunks) * (1 - 1/(p.ReconstructionThreshold*float64(p.CodingRate)))))

}

var (
	ParametersMap = map[byte]BlobParameters{
		0: {CodingRate: 8, ReconstructionThreshold: 0.22, NumChunks: 8192},
	}
)

// AssignmentCoordinator is responsible for assigning chunks to operators in a way that satisfies the security
// requirements of the protocol, as well as the constraints imposed by the specific blob version.
type AssignmentCoordinatorV2 interface {

	// GetAssignments calculates the full set of node assignments
	GetAssignments(state *OperatorState, blobVersion byte, quorum QuorumID) (map[OperatorID]Assignment, error)

	// GetAssignment calculates the assignment for a specific operator
	GetAssignment(state *OperatorState, blobVersion byte, quorum QuorumID, id OperatorID) (Assignment, error)

	// GetChunkLength determines the length of a chunk given the blob version and blob length
	GetChunkLength(blobVersion byte, blobLength uint) (uint, error)
}

type StdAssignmentCoordinatorV2 struct {
}

var _ AssignmentCoordinatorV2 = (*StdAssignmentCoordinatorV2)(nil)

func (c *StdAssignmentCoordinatorV2) GetAssignments(state *OperatorState, blobVersion byte, quorum QuorumID) (map[OperatorID]Assignment, error) {

	params, ok := ParametersMap[blobVersion]
	if !ok {
		return nil, fmt.Errorf("blob version %d not found", blobVersion)
	}

	ops, ok := state.Operators[quorum]
	if !ok {
		return nil, fmt.Errorf("no operators found for quorum %d", quorum)
	}

	if len(ops) > int(params.MaxNumOperators()) {
		return nil, fmt.Errorf("too many operators for blob version %d", blobVersion)
	}

	n := big.NewInt(int64(len(ops)))
	m := big.NewInt(int64(params.NumChunks))

	type assignment struct {
		id     OperatorID
		index  uint
		chunks uint
		stake  *big.Int
	}

	chunkAssignments := make([]assignment, 0, len(ops))
	for ID, r := range state.Operators[quorum] {

		num := new(big.Int).Mul(r.Stake, new(big.Int).Sub(m, n))
		denom := state.Totals[quorum].Stake

		chunks := roundUpDivideBig(num, denom)

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

	if blobLength == 0 {
		return 0, fmt.Errorf("blob length must be greater than 0")
	}

	// Check that the blob length is a power of 2
	if blobLength&(blobLength-1) != 0 {
		return 0, fmt.Errorf("blob length %d is not a power of 2", blobLength)
	}

	if _, ok := ParametersMap[blobVersion]; !ok {
		return 0, fmt.Errorf("blob version %d not found", blobVersion)
	}

	chunkLength := blobLength * ParametersMap[blobVersion].CodingRate / ParametersMap[blobVersion].NumChunks
	if chunkLength == 0 {
		chunkLength = 1
	}

	return chunkLength, nil

}
