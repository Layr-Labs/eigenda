package v2

import (
	"fmt"
	"math/big"
	"sort"

	"github.com/Layr-Labs/eigenda/core"
)

func GetAssignments(state *core.OperatorState, blobVersion byte, quorum uint8) (map[core.OperatorID]Assignment, error) {

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

	numOperators := big.NewInt(int64(len(ops)))
	numChunks := big.NewInt(int64(params.NumChunks))

	type assignment struct {
		id     core.OperatorID
		index  uint32
		chunks uint32
		stake  *big.Int
	}

	chunkAssignments := make([]assignment, 0, len(ops))
	for ID, r := range state.Operators[quorum] {

		num := new(big.Int).Mul(r.Stake, new(big.Int).Sub(numChunks, numOperators))
		denom := state.Totals[quorum].Stake

		chunks := core.RoundUpDivideBig(num, denom)

		chunkAssignments = append(chunkAssignments, assignment{id: ID, index: uint32(r.Index), chunks: uint32(chunks.Uint64()), stake: r.Stake})
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

	assignments := make(map[core.OperatorID]Assignment, len(chunkAssignments))
	index := uint32(0)
	for i, a := range chunkAssignments {
		if i < delta {
			a.chunks++
		}

		assignments[a.id] = Assignment{
			StartIndex: index,
			NumChunks:  a.chunks,
		}
		index += a.chunks
	}

	return assignments, nil

}

func GetAssignment(state *core.OperatorState, blobVersion byte, quorum core.QuorumID, id core.OperatorID) (Assignment, error) {

	assignments, err := GetAssignments(state, blobVersion, quorum)
	if err != nil {
		return Assignment{}, err
	}

	assignment, ok := assignments[id]
	if !ok {
		return Assignment{}, ErrNotFound
	}

	return assignment, nil
}

func GetChunkLength(blobVersion byte, blobLength uint32) (uint32, error) {

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
