package v2_test

import (
	"context"
	"fmt"
	"math/big"
	"math/rand"
	"testing"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/mock"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/stretchr/testify/assert"
)

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
	_, err := corev2.GetAssignmentsForBlob(nil, blobParams, []core.QuorumID{0})
	assert.Error(t, err)
}

func TestNonExistentQuorum(t *testing.T) {
	state := dat.GetTotalOperatorState(context.Background(), 0)
	nonExistentQuorum := uint8(99) // Assuming this quorum doesn't exist
	_, err := corev2.GetAssignmentsForBlob(state.OperatorState, blobParams, []core.QuorumID{nonExistentQuorum})
	assert.Error(t, err)
}

func TestNonExistentOperator(t *testing.T) {
	state := dat.GetTotalOperatorState(context.Background(), 0)
	_, err := corev2.GetAssignmentForBlob(state.OperatorState, blobParams, []core.QuorumID{0}, mock.MakeOperatorId(999))
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

	assignments, err := corev2.GetAssignmentsForBlob(state.OperatorState, blobParams, []core.QuorumID{0})
	assert.NoError(t, err)
	assert.Len(t, assignments, 1)
	assignment, exists := assignments[mock.MakeOperatorId(0)]
	assert.True(t, exists)
	assert.Equal(t, blobParams.NumChunks/blobParams.CodingRate, assignment.NumChunks())
}

func TestTwoQuorums(t *testing.T) {

	stakes := map[core.QuorumID]map[core.OperatorID]int{
		0: {
			mock.MakeOperatorId(0): 1,
			mock.MakeOperatorId(1): 10,
			mock.MakeOperatorId(2): 1,
			mock.MakeOperatorId(3): 1,
			mock.MakeOperatorId(4): 3,
		},
		1: {
			mock.MakeOperatorId(0): 2,
			mock.MakeOperatorId(1): 1,
			mock.MakeOperatorId(2): 10,
			mock.MakeOperatorId(3): 1,
		},
	}

	dat, err := mock.NewChainDataMock(stakes)
	assert.NoError(t, err)

	state := dat.GetTotalOperatorState(context.Background(), 0)

	assignmentsBothQuorums, err := corev2.GetAssignmentsForBlob(state.OperatorState, blobParams, []core.QuorumID{0, 1})
	assert.NoError(t, err)
	assert.Len(t, assignmentsBothQuorums, 5)

	assignmentsQuorum0, err := corev2.GetAssignmentsForBlob(state.OperatorState, blobParams, []core.QuorumID{0})
	assert.NoError(t, err)
	assert.Len(t, assignmentsQuorum0, 5)

	assignmentsQuorum1, err := corev2.GetAssignmentsForBlob(state.OperatorState, blobParams, []core.QuorumID{1})
	assert.NoError(t, err)
	assert.Len(t, assignmentsQuorum1, 4)

	// Check that the lenght of the assignment for each operator is equal to the maximum of the assignments for that operator in each quorum
	for id := range assignmentsBothQuorums {

		// Get the bigger assignemnt between the two quorums
		maxChunks := uint32(0)
		assignment, ok := assignmentsQuorum0[id]
		if ok {
			maxChunks = assignment.NumChunks()
		}
		assignment, ok = assignmentsQuorum1[id]
		if ok {
			if assignment.NumChunks() > maxChunks {
				maxChunks = assignment.NumChunks()
			}
		}
		fmt.Println(id, assignmentsBothQuorums[id].NumChunks(), maxChunks)
		assert.LessOrEqual(t, assignmentsBothQuorums[id].NumChunks(), maxChunks)
	}
}

func TestManyQuorums(t *testing.T) {

	testCases := []uint8{1, 2, 3, 4, 5, 10, 15, 20, 50, 100, 200, 255}
	numOperators := 100

	for _, numQuorums := range testCases {
		t.Run("Numer of quorums: "+string(numQuorums), func(t *testing.T) {

			stakes := make(map[core.QuorumID]map[core.OperatorID]int)
			quorumNumbers := make([]core.QuorumID, numQuorums)

			for i := uint8(0); i < numQuorums; i++ {
				quorumNumbers[i] = i
				stakes[i] = make(map[core.OperatorID]int)
				for j := 0; j < numOperators; j++ {
					stakes[i][mock.MakeOperatorId(j)] = rand.Intn(100) + 1
				}
			}

			dat, err := mock.NewChainDataMock(stakes)
			if err != nil {
				t.Fatal(err)
			}

			state := dat.GetTotalOperatorState(context.Background(), 0)

			assignments, err := corev2.GetAssignmentsForBlob(state.OperatorState, blobParams, quorumNumbers)
			assert.NoError(t, err)

			for _, i := range quorumNumbers {

				assignmentForQuorum, err := corev2.GetAssignmentsForBlob(state.OperatorState, blobParams, []core.QuorumID{i})
				assert.NoError(t, err)

				for id := range assignments {
					assert.GreaterOrEqual(t, assignments[id].NumChunks(), assignmentForQuorum[id].NumChunks())
				}
			}

		})
	}
}

func TestValidatorSizes(t *testing.T) {
	thresholdBips := blobParams.GetReconstructionThresholdBips()

	testCases := []struct {
		name              string
		operatorStake     uint32 // Stake for the operator we're testing
		otherStake        uint32 // Stake for the other operator(s) in the quorum
		expectedNumChunks uint32 // Expected number of chunks assigned
	}{
		{
			name:              "Negligible Stake",
			operatorStake:     1,
			otherStake:        1000000, // Large stake to ensure test operator's percentage is negligible
			expectedNumChunks: 1,       // Minimum assignment
		},
		{
			name:              "Exactly Threshold Stake",
			operatorStake:     thresholdBips,
			otherStake:        10000 - thresholdBips, // Ensure we get exactly the threshold percentage
			expectedNumChunks: blobParams.NumChunks / blobParams.CodingRate,
		},
		{
			name:          "Double Threshold Stake",
			operatorStake: thresholdBips * 2,
			otherStake:    10000 - (thresholdBips * 2), // Ensure percentage is double the threshold
			// Capped at the threshold
			expectedNumChunks: blobParams.NumChunks / blobParams.CodingRate,
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
			assignment, err := corev2.GetAssignmentForBlob(state.OperatorState, blobParams, []core.QuorumID{0}, mock.MakeOperatorId(0))
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

func FuzzOperatorAssignmentsV2(f *testing.F) {

	// Add distributions to fuzz

	for i := 1; i < 100; i++ {
		f.Add(i)
	}

	for i := 0; i < 100; i++ {
		f.Add(rand.Intn(int(blobParams.MaxNumOperators)-100) + 100)
	}

	f.Fuzz(func(t *testing.T, numOperators int) {

		// Generate a random slice of integers of length n

		stakes := map[core.QuorumID]map[core.OperatorID]int{
			0: {},
			1: {},
		}
		for i := 0; i < numOperators; i++ {
			stakes[0][mock.MakeOperatorId(i)] = rand.Intn(100) + 1
			stakes[1][mock.MakeOperatorId(i)] = rand.Intn(100) + 10
		}

		dat, err := mock.NewChainDataMock(stakes)
		if err != nil {
			t.Fatal(err)
		}

		state := dat.GetTotalOperatorState(context.Background(), 0)

		assignments, err := corev2.GetAssignmentsForBlob(state.OperatorState, blobParams, []core.QuorumID{0, 1})
		assert.NoError(t, err)

		// Check that the total number of chunks satisfies expected bounds
		if numOperators > 20 {

			totalChunks := uint32(0)
			for _, assignment := range assignments {
				totalChunks += assignment.NumChunks()
			}
			assert.GreaterOrEqual(t, totalChunks, blobParams.NumChunks-blobParams.MaxNumOperators)
		}

		// Sample a random collection of operators whose total stake exceeds the reconstruction threshold and check that they can reconstruct the blob

		// Get the total stake for the quorum
		totalStake := new(big.Int).Set(state.OperatorState.Totals[0].Stake)

		// Calculate the threshold stake required for reconstruction\
		thresholdStake := core.RoundUpDivideBig(new(big.Int).Mul(totalStake, big.NewInt(int64(blobParams.GetReconstructionThresholdBips()))), big.NewInt(10000))

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
		minChunksNeeded := blobParams.NumChunks / blobParams.CodingRate

		// Assert that the sampled operators have enough unique chunks to reconstruct the blob
		assert.GreaterOrEqual(t, uint32(len(uniqueChunkIndices)), minChunksNeeded,
			"Sampled operators should have enough unique chunks to reconstruct the blob")

		if uint32(len(uniqueChunkIndices)) < minChunksNeeded {

			fmt.Println("Quorum: 0")
			for opID, stake := range stakes[0] {
				fmt.Println("Stake: ", stake, "Operator: ", opID.Hex())
			}

			fmt.Println("Quorum: 1")
			for opID, stake := range stakes[1] {
				fmt.Println("Stake: ", stake, "Operator: ", opID.Hex())
			}

			fmt.Println("Sampled operators:")
			for _, opID := range sampledOperators {
				fmt.Println(opID.Hex())
			}

			t.Fatal("Sampled operators should have enough unique chunks to reconstruct the blob")
		}
	})
}
