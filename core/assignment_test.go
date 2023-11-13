package core_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Layr-Labs/eigenda/core"
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

	assignments, info, err := coordinator.GetAssignments(operatorState, 0, uint(2))
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

		assignment, info, err := coordinator.GetOperatorAssignment(operatorState, 0, uint(2), operatorID)
		assert.NoError(t, err)

		assert.Equal(t, assignment, expectedAssignments[operatorID])
		assert.Equal(t, expectedInfo, info)

	}

}
