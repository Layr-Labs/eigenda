package v2_test

import (
	"context"
	"math/rand"
	"testing"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/mock"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/stretchr/testify/assert"
)

func TestOperatorAssignmentsV2(t *testing.T) {

	state := dat.GetTotalOperatorState(context.Background(), 0)
	operatorState := state.OperatorState

	assignments, err := corev2.GetAssignments(operatorState, blobParams, 0)
	assert.NoError(t, err)
	expectedAssignments := map[core.OperatorID]corev2.Assignment{
		mock.MakeOperatorId(0): {
			StartIndex: 7802,
			NumChunks:  390,
		},
		mock.MakeOperatorId(1): {
			StartIndex: 7022,
			NumChunks:  780,
		},
		mock.MakeOperatorId(2): {
			StartIndex: 5852,
			NumChunks:  1170,
		},
		mock.MakeOperatorId(3): {
			StartIndex: 4291,
			NumChunks:  1561,
		},
		mock.MakeOperatorId(4): {
			StartIndex: 2340,
			NumChunks:  1951,
		},
		mock.MakeOperatorId(5): {
			StartIndex: 0,
			NumChunks:  2340,
		},
	}

	for operatorID, assignment := range assignments {

		assert.Equal(t, assignment, expectedAssignments[operatorID])

		assignment, err := corev2.GetAssignment(operatorState, blobParams, 0, operatorID)
		assert.NoError(t, err)

		assert.Equal(t, assignment, expectedAssignments[operatorID])

	}

}

func TestAssignmentWithTooManyOperators(t *testing.T) {

	numOperators := blobParams.MaxNumOperators + 1

	stakes := map[core.QuorumID]map[core.OperatorID]int{
		0: {},
	}
	for i := 0; i < int(numOperators); i++ {
		stakes[0][mock.MakeOperatorId(i)] = rand.Intn(100) + 1
	}

	dat, err := mock.NewChainDataMock(stakes)
	if err != nil {
		t.Fatal(err)
	}

	state := dat.GetTotalOperatorState(context.Background(), 0)

	assert.Equal(t, len(state.Operators[0]), int(numOperators))

	_, err = corev2.GetAssignments(state.OperatorState, blobParams, 0)
	assert.Error(t, err)

}

func FuzzOperatorAssignmentsV2(f *testing.F) {

	// Add distributions to fuzz

	for i := 1; i < 100; i++ {
		f.Add(i)
	}

	for i := 0; i < 100; i++ {
		f.Add(rand.Intn(2048) + 100)
	}

	for i := 0; i < 5; i++ {
		f.Add(int(blobParams.MaxNumOperators))
	}

	f.Fuzz(func(t *testing.T, numOperators int) {

		// Generate a random slice of integers of length n

		stakes := map[core.QuorumID]map[core.OperatorID]int{
			0: {},
		}
		for i := 0; i < numOperators; i++ {
			stakes[0][mock.MakeOperatorId(i)] = rand.Intn(100) + 1
		}

		dat, err := mock.NewChainDataMock(stakes)
		if err != nil {
			t.Fatal(err)
		}

		state := dat.GetTotalOperatorState(context.Background(), 0)

		assignments, err := corev2.GetAssignments(state.OperatorState, blobParams, 0)
		assert.NoError(t, err)

		// Check that the total number of chunks is correct
		totalChunks := uint32(0)
		for _, assignment := range assignments {
			totalChunks += assignment.NumChunks
		}
		assert.Equal(t, totalChunks, blobParams.NumChunks)

		// Check that each operator's assignment satisfies the security requirement
		for operatorID, assignment := range assignments {

			totalStake := uint32(state.Totals[0].Stake.Uint64())
			myStake := uint32(state.Operators[0][operatorID].Stake.Uint64())

			reconstructionThreshold := 0.22
			LHS := assignment.NumChunks * totalStake * blobParams.CodingRate * uint32(reconstructionThreshold*100)
			RHS := 100 * myStake * blobParams.NumChunks

			assert.GreaterOrEqual(t, LHS, RHS)
		}
	})
}

func TestChunkLength(t *testing.T) {
	pairs := []struct {
		blobLength  uint32
		chunkLength uint32
	}{
		{512, 1},
		{1024, 1},
		{2048, 2},
		{4096, 4},
		{8192, 8},
	}

	for _, pair := range pairs {
		chunkLength, err := corev2.GetChunkLength(pair.blobLength, blobParams)
		assert.NoError(t, err)
		assert.Equal(t, pair.chunkLength, chunkLength)
	}
}

func TestInvalidChunkLength(t *testing.T) {
	invalidLengths := []uint32{
		0,
		3,
		5,
		6,
		7,
		9,
		10,
		11,
		12,
		13,
		14,
		15,
		31,
		63,
		127,
		255,
		511,
		1023,
	}

	for _, length := range invalidLengths {
		_, err := corev2.GetChunkLength(length, blobParams)
		assert.Error(t, err)
	}
}

func TestNilStateAssignments(t *testing.T) {
	_, err := corev2.GetAssignments(nil, blobParams, 0)
	assert.Error(t, err)
}

func TestNonExistentQuorum(t *testing.T) {
	state := dat.GetTotalOperatorState(context.Background(), 0)
	nonExistentQuorum := uint8(99) // Assuming this quorum doesn't exist
	_, err := corev2.GetAssignments(state.OperatorState, blobParams, nonExistentQuorum)
	assert.Error(t, err)
}

func TestNonExistentOperator(t *testing.T) {
	state := dat.GetTotalOperatorState(context.Background(), 0)
	nonExistentOperatorID := mock.MakeOperatorId(999) // Assuming this operator doesn't exist
	_, err := corev2.GetAssignment(state.OperatorState, blobParams, 0, nonExistentOperatorID)
	assert.Error(t, err)
	assert.Equal(t, corev2.ErrNotFound, err)
}

func TestSingleOperator(t *testing.T) {
	stakes := map[core.QuorumID]map[core.OperatorID]int{
		0: {
			mock.MakeOperatorId(0): 100,
		},
	}

	dat, err := mock.NewChainDataMock(stakes)
	assert.NoError(t, err)

	state := dat.GetTotalOperatorState(context.Background(), 0)
	assignments, err := corev2.GetAssignments(state.OperatorState, blobParams, 0)
	assert.NoError(t, err)

	assert.Len(t, assignments, 1)
	assignment, exists := assignments[mock.MakeOperatorId(0)]
	assert.True(t, exists)
	assert.Equal(t, uint32(0), assignment.StartIndex)
	assert.Equal(t, blobParams.NumChunks, assignment.NumChunks)
}

func TestEqualStakeOperators(t *testing.T) {
	// Test with operators having equal stake to ensure correct index-based sorting
	numOperators := 5
	stakes := map[core.QuorumID]map[core.OperatorID]int{
		0: {},
	}

	// All operators have the same stake (100)
	for i := 0; i < numOperators; i++ {
		stakes[0][mock.MakeOperatorId(i)] = 100
	}

	dat, err := mock.NewChainDataMock(stakes)
	assert.NoError(t, err)

	state := dat.GetTotalOperatorState(context.Background(), 0)
	assignments, err := corev2.GetAssignments(state.OperatorState, blobParams, 0)
	assert.NoError(t, err)

	// Since all operators have equal stake, they should be sorted by index
	// Check that operators with lower indices have higher number of chunks
	// or same number of chunks with earlier start indices
	var prevNumChunks uint32
	var prevIndex uint32

	// Sort operators by their index for verification
	var operatorIDs []int
	for i := 0; i < numOperators; i++ {
		operatorIDs = append(operatorIDs, i)
	}

	for _, i := range operatorIDs {
		operatorID := mock.MakeOperatorId(i)
		assignment := assignments[operatorID]

		if i > 0 {
			// Either the current operator should have fewer chunks,
			// or if equal chunks, a higher start index
			if prevNumChunks == assignment.NumChunks {
				assert.True(t, prevIndex < assignment.StartIndex,
					"Operators with same stake should be sorted by index")
			} else {
				assert.True(t, prevNumChunks >= assignment.NumChunks,
					"Operators with lower indices should have more chunks")
			}
		}

		prevNumChunks = assignment.NumChunks
		prevIndex = assignment.StartIndex
	}

	// Verify total chunks is correct
	totalChunks := uint32(0)
	for _, assignment := range assignments {
		totalChunks += assignment.NumChunks
	}
	assert.Equal(t, blobParams.NumChunks, totalChunks)
}

func TestExactlyMaxOperators(t *testing.T) {
	numOperators := blobParams.MaxNumOperators

	stakes := map[core.QuorumID]map[core.OperatorID]int{
		0: {},
	}
	for i := 0; i < int(numOperators); i++ {
		stakes[0][mock.MakeOperatorId(i)] = rand.Intn(100) + 1
	}

	dat, err := mock.NewChainDataMock(stakes)
	assert.NoError(t, err)

	state := dat.GetTotalOperatorState(context.Background(), 0)
	assert.Equal(t, len(state.Operators[0]), int(numOperators))

	// Test should not error with exactly the max number of operators
	assignments, err := corev2.GetAssignments(state.OperatorState, blobParams, 0)
	assert.NoError(t, err)
	assert.Len(t, assignments, int(numOperators))

	// Check total chunks
	totalChunks := uint32(0)
	for _, assignment := range assignments {
		totalChunks += assignment.NumChunks
	}
	assert.Equal(t, blobParams.NumChunks, totalChunks)
}

func TestZeroStakeOperator(t *testing.T) {
	// Test with one operator having zero stake
	stakes := map[core.QuorumID]map[core.OperatorID]int{
		0: {
			mock.MakeOperatorId(0): 100,
			mock.MakeOperatorId(1): 0, // Zero stake
			mock.MakeOperatorId(2): 200,
		},
	}

	dat, err := mock.NewChainDataMock(stakes)
	assert.NoError(t, err)

	state := dat.GetTotalOperatorState(context.Background(), 0)
	assignments, err := corev2.GetAssignments(state.OperatorState, blobParams, 0)
	assert.NoError(t, err)

	// Operator with zero stake should get zero chunks
	zeroStakeOp := mock.MakeOperatorId(1)
	assignment, exists := assignments[zeroStakeOp]
	assert.True(t, exists)
	assert.Equal(t, uint32(0), assignment.NumChunks)

	// Verify total chunks is correct
	totalChunks := uint32(0)
	for _, assignment := range assignments {
		totalChunks += assignment.NumChunks
	}
	assert.Equal(t, blobParams.NumChunks, totalChunks)
}

func TestImbalancedStakeDistribution(t *testing.T) {
	// Test with highly imbalanced stake distribution
	stakes := map[core.QuorumID]map[core.OperatorID]int{
		0: {
			mock.MakeOperatorId(0): 1000000, // Very high stake
			mock.MakeOperatorId(1): 1,       // Very low stake
			mock.MakeOperatorId(2): 1,       // Very low stake
			mock.MakeOperatorId(3): 1,       // Very low stake
		},
	}

	dat, err := mock.NewChainDataMock(stakes)
	assert.NoError(t, err)

	state := dat.GetTotalOperatorState(context.Background(), 0)
	assignments, err := corev2.GetAssignments(state.OperatorState, blobParams, 0)
	assert.NoError(t, err)

	// High stake operator should get almost all chunks
	highStakeOp := mock.MakeOperatorId(0)
	assignment, exists := assignments[highStakeOp]
	assert.True(t, exists)
	assert.True(t, assignment.NumChunks > blobParams.NumChunks*9/10,
		"High stake operator should get most chunks")

	// Verify total chunks is correct
	totalChunks := uint32(0)
	for _, assignment := range assignments {
		totalChunks += assignment.NumChunks
	}
	assert.Equal(t, blobParams.NumChunks, totalChunks)
}

func TestNilBlobParamsForChunkLength(t *testing.T) {
	_, err := corev2.GetChunkLength(1024, nil)
	assert.Error(t, err)
}

func TestNilBlobParamsForGetAssignment(t *testing.T) {
	state := dat.GetTotalOperatorState(context.Background(), 0)
	_, err := corev2.GetAssignment(state.OperatorState, nil, 0, mock.MakeOperatorId(0))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "blob params cannot be nil")
}

func TestConsecutiveAssignments(t *testing.T) {
	// Test that assignments have consecutive, non-overlapping chunk ranges
	numOperators := 10
	stakes := map[core.QuorumID]map[core.OperatorID]int{
		0: {},
	}

	// Generate random stakes
	for i := 0; i < numOperators; i++ {
		stakes[0][mock.MakeOperatorId(i)] = rand.Intn(1000) + 1
	}

	dat, err := mock.NewChainDataMock(stakes)
	assert.NoError(t, err)

	state := dat.GetTotalOperatorState(context.Background(), 0)
	assignments, err := corev2.GetAssignments(state.OperatorState, blobParams, 0)
	assert.NoError(t, err)

	// Create a map to track which chunks are assigned
	chunkAssignments := make(map[uint32]core.OperatorID)

	// Track which indexes are the start or end of assignments
	chunkRanges := make(map[uint32]bool)

	// Fill the assignments map and track the ranges
	for opID, assignment := range assignments {
		startIdx := assignment.StartIndex
		endIdx := startIdx + assignment.NumChunks

		// Mark start and end points
		chunkRanges[startIdx] = true
		chunkRanges[endIdx] = true

		// Check each chunk in this assignment
		for i := startIdx; i < endIdx; i++ {
			_, alreadyAssigned := chunkAssignments[i]
			assert.False(t, alreadyAssigned, "Chunk %d is assigned more than once", i)
			chunkAssignments[i] = opID
		}
	}

	// Make sure all chunks from 0 to NumChunks-1 are assigned
	for i := uint32(0); i < blobParams.NumChunks; i++ {
		_, assigned := chunkAssignments[i]
		assert.True(t, assigned, "Chunk %d is not assigned to any operator", i)
	}

	// Check specific properties of assignment boundaries

	// 1. Should start at 0
	assert.True(t, chunkRanges[0], "Assignments should start at index 0")

	// 2. Should end at NumChunks
	assert.True(t, chunkRanges[blobParams.NumChunks], "Assignments should end at index NumChunks")

	// 3. For every operator, check that their assignment is a continuous range
	for opID, assignment := range assignments {
		startIdx := assignment.StartIndex
		endIdx := startIdx + assignment.NumChunks

		// Skip operators with 0 chunks
		if assignment.NumChunks == 0 {
			continue
		}

		// Check that all chunks in range belong to the same operator
		for i := startIdx; i < endIdx; i++ {
			assignedOp, exists := chunkAssignments[i]
			assert.True(t, exists, "Chunk %d should be assigned", i)
			assert.Equal(t, opID, assignedOp, "Chunk %d should be assigned to operator %s", i, opID)
		}

		// Check boundaries - chunk before start should be a different operator (if it exists)
		if startIdx > 0 {
			prevChunkOp, exists := chunkAssignments[startIdx-1]
			assert.True(t, exists, "Chunk %d should be assigned", startIdx-1)
			assert.NotEqual(t, opID, prevChunkOp,
				"Chunk before range start should belong to a different operator")
		}

		// Check that chunk at end belongs to a different operator (if within bounds)
		if endIdx < blobParams.NumChunks {
			nextChunkOp, exists := chunkAssignments[endIdx]
			assert.True(t, exists, "Chunk %d should be assigned", endIdx)
			assert.NotEqual(t, opID, nextChunkOp,
				"Chunk at range end should belong to a different operator")
		}
	}

	// Check that the full range of chunks is covered
	totalAssignedChunks := uint32(0)
	for _, assignment := range assignments {
		totalAssignedChunks += assignment.NumChunks
	}
	assert.Equal(t, blobParams.NumChunks, totalAssignedChunks,
		"Total assigned chunks should equal the blob parameter NumChunks")
}
