package core_test

import (
	"context"
	"fmt"
	"math/rand"
	"testing"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/mock"
	"github.com/stretchr/testify/assert"
)

func makeOperatorId(id int) core.OperatorID {
	data := [32]byte{}
	copy(data[:], []byte(fmt.Sprintf("%d", id)))
	return data
}

func TestOperatorAssignments(t *testing.T) {

	state := dat.GetTotalOperatorState(context.Background(), 0)
	operatorState := state.OperatorState
	coordinator := &core.StdAssignmentCoordinator{}

	quorumInfo := &core.BlobQuorumInfo{
		SecurityParam: core.SecurityParam{
			QuorumID:           0,
			AdversaryThreshold: 50,
			QuorumThreshold:    100,
		},
		ChunkLength: 10,
	}

	blobLength := uint(100)

	assignments, info, err := coordinator.GetAssignments(operatorState, blobLength, quorumInfo)
	assert.NoError(t, err)
	expectedAssignments := map[core.OperatorID]core.Assignment{
		makeOperatorId(0): {
			StartIndex: 0,
			NumChunks:  1,
		},
		makeOperatorId(1): {
			StartIndex: 1,
			NumChunks:  1,
		},
		makeOperatorId(2): {
			StartIndex: 2,
			NumChunks:  2,
		},
		makeOperatorId(3): {
			StartIndex: 4,
			NumChunks:  2,
		},
		makeOperatorId(4): {
			StartIndex: 6,
			NumChunks:  2,
		},
		makeOperatorId(5): {
			StartIndex: 8,
			NumChunks:  3,
		},
		makeOperatorId(6): {
			StartIndex: 11,
			NumChunks:  3,
		},
		makeOperatorId(7): {
			StartIndex: 14,
			NumChunks:  3,
		},
		makeOperatorId(8): {
			StartIndex: 17,
			NumChunks:  4,
		},
		makeOperatorId(9): {
			StartIndex: 21,
			NumChunks:  4,
		},
	}
	expectedInfo := core.AssignmentInfo{
		TotalChunks: 25,
	}

	assert.Equal(t, expectedInfo, info)

	for operatorID, assignment := range assignments {
		assert.Equal(t, assignment, expectedAssignments[operatorID])

		header := &core.BlobHeader{
			BlobCommitments: core.BlobCommitments{
				Length: blobLength,
			},
			QuorumInfos: []*core.BlobQuorumInfo{quorumInfo},
		}

		assignment, info, err := coordinator.GetOperatorAssignment(operatorState, header, 0, operatorID)
		assert.NoError(t, err)

		assert.Equal(t, assignment, expectedAssignments[operatorID])
		assert.Equal(t, expectedInfo, info)

	}

}

func FuzzOperatorAssignments(f *testing.F) {

	// Add distributions to fuzz
	asn := &core.StdAssignmentCoordinator{}

	for i := 1; i < 100; i++ {
		f.Add(i, true)
	}

	for i := 1; i < 100; i++ {
		f.Add(i, false)
	}

	for i := 0; i < 100; i++ {
		f.Add(rand.Intn(1000)+1, rand.Intn(2) == 0)
	}

	f.Fuzz(func(t *testing.T, numOperators int, useTargetNumChunks bool) {

		// Generate a random slice of integers of length n

		stakes := make([]int, numOperators)
		for i := range stakes {
			stakes[i] = rand.Intn(100)
		}

		advThreshold := rand.Intn(99)
		quorumThreshold := rand.Intn(100-advThreshold) + advThreshold + 1

		param := &core.SecurityParam{
			QuorumID:           0,
			AdversaryThreshold: uint8(advThreshold),
			QuorumThreshold:    uint8(quorumThreshold),
		}

		dat, err := mock.NewChainDataMock(stakes)
		if err != nil {
			t.Fatal(err)
		}

		state := dat.GetTotalOperatorState(context.Background(), 0)

		blobLength := uint(rand.Intn(100000))

		targetNumChunks := uint(0)
		if useTargetNumChunks {
			targetNumChunks = uint(rand.Intn(1000))
		}

		fmt.Println("advThreshold", advThreshold, "quorumThreshold", quorumThreshold, "numOperators", numOperators, "blobLength", blobLength)

		chunkLength, err := asn.CalculateChunkLength(state.OperatorState, blobLength, targetNumChunks, param)
		assert.NoError(t, err)

		quorumInfo := &core.BlobQuorumInfo{
			SecurityParam: *param,
			ChunkLength:   chunkLength,
		}

		ok, err := asn.ValidateChunkLength(state.OperatorState, blobLength, quorumInfo)
		assert.NoError(t, err)
		assert.True(t, ok)

		assignments, info, err := asn.GetAssignments(state.OperatorState, blobLength, quorumInfo)
		assert.NoError(t, err)

		// fmt.Println("advThreshold", advThreshold, "quorumThreshold", quorumThreshold, "numOperators", numOperators, "chunkLength", chunkLength, "blobLength", blobLength)

		if useTargetNumChunks {

			quorumInfo.ChunkLength = chunkLength * 2
			ok, err := asn.ValidateChunkLength(state.OperatorState, blobLength, quorumInfo)

			// If it's possible to make the chunk larger, then the number of chunks should fall within the target
			if ok && err == nil {
				assert.GreaterOrEqual(t, targetNumChunks, info.TotalChunks)
				assert.Greater(t, info.TotalChunks, targetNumChunks/2)
			}
		}

		// Check that each operator's assignment satisfies the security requirement
		for operatorID, assignment := range assignments {

			totalStake := state.Totals[0].Stake
			myStake := state.Operators[0][operatorID].Stake

			valid := assignment.NumChunks*uint((quorumThreshold-advThreshold))*chunkLength*uint(totalStake.Uint64()) >= blobLength*uint(myStake.Uint64())
			assert.True(t, valid)

		}

	})

}
