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

const (
	maxNumOperators = 3537
)

func TestOperatorAssignmentsV2(t *testing.T) {

	state := dat.GetTotalOperatorState(context.Background(), 0)
	operatorState := state.OperatorState

	blobVersion := corev2.BlobVersion(0)

	assignments, err := corev2.GetAssignments(operatorState, blobVersion, 0)
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

		assignment, err := corev2.GetAssignment(operatorState, blobVersion, 0, operatorID)
		assert.NoError(t, err)

		assert.Equal(t, assignment, expectedAssignments[operatorID])

	}

}

func TestMaxNumOperators(t *testing.T) {

	assert.Equal(t, corev2.ParametersMap[0].MaxNumOperators(), uint32(maxNumOperators))

}

func TestAssignmentWithTooManyOperators(t *testing.T) {

	numOperators := maxNumOperators + 1

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

	assert.Equal(t, len(state.Operators[0]), numOperators)

	blobVersion := corev2.BlobVersion(0)

	_, err = corev2.GetAssignments(state.OperatorState, blobVersion, 0)
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
		f.Add(maxNumOperators)
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

		blobVersion := corev2.BlobVersion(0)

		assignments, err := corev2.GetAssignments(state.OperatorState, blobVersion, 0)
		assert.NoError(t, err)

		// Check that the total number of chunks is correct
		totalChunks := uint32(0)
		for _, assignment := range assignments {
			totalChunks += assignment.NumChunks
		}
		assert.Equal(t, totalChunks, corev2.ParametersMap[blobVersion].NumChunks)

		// Check that each operator's assignment satisfies the security requirement
		for operatorID, assignment := range assignments {

			totalStake := uint32(state.Totals[0].Stake.Uint64())
			myStake := uint32(state.Operators[0][operatorID].Stake.Uint64())

			LHS := assignment.NumChunks * totalStake * corev2.ParametersMap[blobVersion].CodingRate * uint32(corev2.ParametersMap[blobVersion].ReconstructionThreshold*100)
			RHS := 100 * myStake * corev2.ParametersMap[blobVersion].NumChunks

			assert.GreaterOrEqual(t, LHS, RHS)

		}

	})

}

func TestChunkLength(t *testing.T) {

	blobVersion := corev2.BlobVersion(0)

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

		chunkLength, err := corev2.GetChunkLength(blobVersion, pair.blobLength)

		assert.NoError(t, err)

		assert.Equal(t, pair.chunkLength, chunkLength)
	}

}

func TestInvalidChunkLength(t *testing.T) {

	blobVersion := corev2.BlobVersion(0)

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

		_, err := corev2.GetChunkLength(blobVersion, length)
		assert.Error(t, err)
	}

}
