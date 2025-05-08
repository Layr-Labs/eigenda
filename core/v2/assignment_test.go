package v2_test

import (
	"context"
	"math"
	"math/big"
	"math/rand"
	"testing"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/mock"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/stretchr/testify/assert"
)

var blobKey1 = []byte("blobKey1")
var blobKey2 = []byte("blobKey2")

// func TestOperatorAssignmentsV2(t *testing.T) {

// 	state := dat.GetTotalOperatorState(context.Background(), 0)
// 	operatorState := state.OperatorState

// 	assignments, err := corev2.GetAssignments(operatorState, blobParams, 0)
// 	assert.NoError(t, err)
// 	expectedAssignments := map[core.OperatorID]corev2.Assignment{
// 		mock.MakeOperatorId(0): {
// 			StartIndex: 7802,
// 			NumChunks:  390,
// 		},
// 		mock.MakeOperatorId(1): {
// 			StartIndex: 7022,
// 			NumChunks:  780,
// 		},
// 		mock.MakeOperatorId(2): {
// 			StartIndex: 5852,
// 			NumChunks:  1170,
// 		},
// 		mock.MakeOperatorId(3): {
// 			StartIndex: 4291,
// 			NumChunks:  1561,
// 		},
// 		mock.MakeOperatorId(4): {
// 			StartIndex: 2340,
// 			NumChunks:  1951,
// 		},
// 		mock.MakeOperatorId(5): {
// 			StartIndex: 0,
// 			NumChunks:  2340,
// 		},
// 	}

// 	for operatorID, assignment := range assignments {

// 		assert.Equal(t, assignment, expectedAssignments[operatorID])

// 		assignment, err := corev2.GetAssignment(operatorState, blobParams, 0, operatorID)
// 		assert.NoError(t, err)

// 		assert.Equal(t, assignment, expectedAssignments[operatorID])

// 	}

// }

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
	_, err := corev2.GetAssignments(nil, blobParams, []core.QuorumID{0}, blobKey1[:])
	assert.Error(t, err)
}

// func TestNonExistentQuorum(t *testing.T) {
// 	state := dat.GetTotalOperatorState(context.Background(), 0)
// 	nonExistentQuorum := uint8(99) // Assuming this quorum doesn't exist
// 	_, err := corev2.GetAssignments(state.OperatorState, blobParams, nonExistentQuorum)
// 	assert.Error(t, err)
// }

func TestNonExistentOperator(t *testing.T) {
	state := dat.GetTotalOperatorState(context.Background(), 0)
	_, err := corev2.GetAssignment(state.OperatorState, blobParams, []core.QuorumID{0}, blobKey1[:], mock.MakeOperatorId(999))
	assert.Error(t, err, corev2.ErrNotFound)
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

	assignments, err := corev2.GetAssignments(state.OperatorState, blobParams, []core.QuorumID{0}, blobKey1[:])
	assert.NoError(t, err)
	assert.Len(t, assignments, 1)
	assignment, exists := assignments[mock.MakeOperatorId(0)]
	assert.True(t, exists)
	// assert.Equal(t, blobParams.NumChunks, assignment.NumChunks())
	_ = assignment
}

func TestValidatorSizes(t *testing.T) {
	thresholdBips := blobParams.ReconstructionThresholdBips
	thresholdPercentage := float64(thresholdBips) / 10000.0

	testCases := []struct {
		name              string
		operatorStake     uint32 // Stake for the operator we're testing
		otherStake        uint32 // Stake for the other operator(s) in the quorum
		expectedNumChunks uint32 // Expected number of chunks assigned
	}{
		{
			name:              "Negligible Stake",
			operatorStake:     1,
			otherStake:        1000000,                   // Large stake to ensure test operator's percentage is negligible
			expectedNumChunks: blobParams.SamplesPerUnit, // Minimum assignment
		},
		{
			name:              "Exactly Threshold Stake",
			operatorStake:     thresholdBips,
			otherStake:        10000 - thresholdBips, // Ensure we get exactly the threshold percentage
			expectedNumChunks: blobParams.SamplesPerUnit * uint32(math.Ceil(thresholdPercentage*float64(blobParams.NumUnits))),
		},
		{
			name:          "Double Threshold Stake",
			operatorStake: thresholdBips * 2,
			otherStake:    10000 - (thresholdBips * 2), // Ensure percentage is double the threshold
			// Capped at the threshold
			expectedNumChunks: blobParams.SamplesPerUnit * uint32(math.Ceil(thresholdPercentage*float64(blobParams.NumUnits))),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create stakes for this test case
			stakes := map[core.QuorumID]map[core.OperatorID]int{
				0: {
					mock.MakeOperatorId(0): int(tc.operatorStake),
					mock.MakeOperatorId(1): int(tc.otherStake),
				},
			}

			dat, err := mock.NewChainDataMock(stakes)
			assert.NoError(t, err)

			state := dat.GetTotalOperatorState(context.Background(), 0)

			// Get assignment for the test operator
			assignment, err := corev2.GetAssignment(state.OperatorState, blobParams, []core.QuorumID{0}, blobKey1[:], mock.MakeOperatorId(0))
			assert.NoError(t, err)

			// Verify the assignment has the expected number of chunks
			assert.Equal(t, tc.expectedNumChunks, assignment.NumChunks(),
				"Expected %d chunks assigned, got %d", tc.expectedNumChunks, assignment.NumChunks())

			// Verify all indices are unique
			uniqueIndices := make(map[uint32]struct{})
			for _, idx := range assignment.GetIndices() {
				uniqueIndices[idx] = struct{}{}
			}
			assert.Equal(t, int(assignment.NumChunks()), len(uniqueIndices),
				"All assigned indices should be unique")

			// Verify all indices are within the valid range
			for _, idx := range assignment.GetIndices() {
				assert.Less(t, idx, blobParams.NumChunks,
					"Index %d is out of valid range [0, %d)", idx, blobParams.NumChunks)
			}
		})
	}
}

func TestDeterministicAssignment(t *testing.T) {
	state := dat.GetTotalOperatorState(context.Background(), 0)
	operatorState := state.OperatorState

	// Get assignments for the same operator with identical headers
	assignment1, err := corev2.GetAssignment(operatorState, blobParams, []core.QuorumID{0}, blobKey1[:], mock.MakeOperatorId(0))
	assert.NoError(t, err)

	assignment2, err := corev2.GetAssignment(operatorState, blobParams, []core.QuorumID{0}, blobKey1[:], mock.MakeOperatorId(0))
	assert.NoError(t, err)

	// Assignments should be identical
	assert.Equal(t, assignment1, assignment2)

	// Get assignments for different operators with the same header
	assignment3, err := corev2.GetAssignment(operatorState, blobParams, []core.QuorumID{0}, blobKey1[:], mock.MakeOperatorId(1))
	assert.NoError(t, err)

	// Assignments should be different for different operators
	assert.NotEqual(t, assignment1, assignment3)

	// Get assignment for the same operator but different header
	assignment4, err := corev2.GetAssignment(operatorState, blobParams, []core.QuorumID{0}, blobKey2[:], mock.MakeOperatorId(0))
	assert.NoError(t, err)

	// Assignments should be different for different headers
	assert.NotEqual(t, assignment1, assignment4)
}

func FuzzOperatorAssignmentsV2(f *testing.F) {

	// Add distributions to fuzz

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

		assignments, err := corev2.GetAssignments(state.OperatorState, blobParams, []core.QuorumID{0}, blobKey2[:])
		assert.NoError(t, err)

		// Check that the total number of chunks satisfies expected bounds
		totalChunks := uint32(0)
		for _, assignment := range assignments {
			totalChunks += assignment.NumChunks()
		}
		assert.GreaterOrEqual(t, totalChunks, blobParams.NumUnits*blobParams.SamplesPerUnit)

		// Sample a random collection of operators whose total stake exceeds the reconstruction threshold and check that they can reconstruct the blob

		// Get the total stake for the quorum
		totalStake := new(big.Int).Set(state.OperatorState.Totals[0].Stake)

		// Calculate the threshold stake required for reconstruction
		thresholdPercentage := float64(blobParams.ReconstructionThresholdBips) / 10000.0 // Convert from basis points to percentage
		thresholdStake := new(big.Int).Mul(totalStake, big.NewInt(int64(thresholdPercentage*10000)))
		thresholdStake.Div(thresholdStake, big.NewInt(10000))

		// Create a slice of operator IDs to randomly sample from
		operatorIDs := make([]core.OperatorID, 0, len(stakes[0]))
		for opID := range stakes[0] {
			operatorIDs = append(operatorIDs, opID)
		}

		// Shuffle the operators for random sampling
		rand.Shuffle(len(operatorIDs), func(i, j int) {
			operatorIDs[i], operatorIDs[j] = operatorIDs[j], operatorIDs[i]
		})

		// Sample operators until we exceed the threshold
		sampledOperators := make([]core.OperatorID, 0)
		currentStake := big.NewInt(0)

		for _, opID := range operatorIDs {
			sampledOperators = append(sampledOperators, opID)
			currentStake.Add(currentStake, state.OperatorState.Operators[0][opID].Stake)

			if currentStake.Cmp(thresholdStake) >= 0 {
				break
			}
		}

		// Verify that the sampled operators' total stake exceeds the threshold
		assert.True(t, currentStake.Cmp(thresholdStake) >= 0,
			"Sampled operators' stake (%s) should exceed threshold stake (%s)",
			currentStake.String(), thresholdStake.String())

		// Collect all unique chunk indices from the sampled operators
		uniqueChunkIndices := make(map[uint32]struct{})
		for _, opID := range sampledOperators {
			assignment, exists := assignments[opID]
			assert.True(t, exists, "Assignment should exist for sampled operator %s", opID.Hex())

			// Add each chunk index to the set of unique indices
			for _, index := range assignment.GetIndices() {
				uniqueChunkIndices[index] = struct{}{}
			}
		}

		// Calculate the minimum required unique chunks for reconstruction
		minChunksNeeded := blobParams.NumUnits / blobParams.CodingRate

		// Assert that the sampled operators have enough unique chunks to reconstruct the blob
		assert.GreaterOrEqual(t, uint32(len(uniqueChunkIndices)), minChunksNeeded,
			"Sampled operators should have enough unique chunks to reconstruct the blob")
	})
}
