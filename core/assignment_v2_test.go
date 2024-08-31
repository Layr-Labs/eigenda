package core_test

import (
	"context"
	"math/rand"
	"testing"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/mock"
	"github.com/stretchr/testify/assert"
)

func TestOperatorAssignmentsV2(t *testing.T) {

	state := dat.GetTotalOperatorState(context.Background(), 0)
	operatorState := state.OperatorState
	coordinator := &core.StdAssignmentCoordinatorV2{}

	blobVersion := byte(0)

	assignments, err := coordinator.GetAssignments(operatorState, blobVersion, 0)
	assert.NoError(t, err)
	expectedAssignments := map[core.OperatorID]core.Assignment{
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

		assignment, err := coordinator.GetAssignment(operatorState, blobVersion, 0, operatorID)
		assert.NoError(t, err)

		assert.Equal(t, assignment, expectedAssignments[operatorID])

	}

}

func FuzzOperatorAssignmentsV2(f *testing.F) {

	// Add distributions to fuzz
	asn := &core.StdAssignmentCoordinatorV2{}

	for i := 1; i < 100; i++ {
		f.Add(i)
	}

	for i := 0; i < 100; i++ {
		f.Add(rand.Intn(2048) + 100)
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

		blobVersion := byte(0)

		assignments, err := asn.GetAssignments(state.OperatorState, blobVersion, 0)
		assert.NoError(t, err)

		// Check that the total number of chunks is correct
		totalChunks := uint(0)
		for _, assignment := range assignments {
			totalChunks += assignment.NumChunks
		}
		assert.Equal(t, totalChunks, core.ParametersMap[blobVersion].NumChunks)

		// Check that each operator's assignment satisfies the security requirement
		for operatorID, assignment := range assignments {

			totalStake := uint(state.Totals[0].Stake.Uint64())
			myStake := uint(state.Operators[0][operatorID].Stake.Uint64())

			LHS := assignment.NumChunks * totalStake * core.ParametersMap[blobVersion].CodingRate * uint(core.ParametersMap[blobVersion].ReconstructionThreshold*100)
			RHS := 100 * myStake * core.ParametersMap[blobVersion].NumChunks

			assert.GreaterOrEqual(t, LHS, RHS)

		}

	})

}
