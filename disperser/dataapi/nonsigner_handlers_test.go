package dataapi_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/Layr-Labs/eigenda/disperser/dataapi"
	"github.com/stretchr/testify/assert"
)

func assertEntry(t *testing.T, quorumIntervals dataapi.OperatorQuorumIntervals, operator string, expected map[uint8][]dataapi.BlockInterval) {
	op, ok := quorumIntervals[operator]
	assert.True(t, ok)
	assert.True(t, reflect.DeepEqual(op, expected))
}

func TestCreateOperatorQuorumIntervalsWithInvalidArgs(t *testing.T) {
	addedQuorums := map[string][]*dataapi.OperatorQuorum{}
	removedQuorums := map[string][]*dataapi.OperatorQuorum{}

	// Empty initial quorums
	operatorInitialQuorum := map[string][]uint8{
		"operator-1": {},
		"operator-2": {0x01},
	}
	_, err := dataapi.CreateOperatorQuorumIntervals(10, 25, operatorInitialQuorum, addedQuorums, removedQuorums)
	assert.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "must be in at least one quorum"))

	// StartBlock > EndBlock
	operatorInitialQuorum = map[string][]uint8{
		"operator-1": {0x00},
		"operator-2": {0x00},
	}
	_, err = dataapi.CreateOperatorQuorumIntervals(100, 25, operatorInitialQuorum, addedQuorums, removedQuorums)
	assert.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "startBlock must be no less than endBlock"))

	// Equal block number
	addedQuorums = map[string][]*dataapi.OperatorQuorum{
		"operator-1": []*dataapi.OperatorQuorum{
			{
				Operator:      "operator-1",
				QuorumNumbers: []uint8{0x01},
				BlockNumber:   12,
			},
		},
	}
	removedQuorums = map[string][]*dataapi.OperatorQuorum{
		"operator-1": []*dataapi.OperatorQuorum{
			{
				Operator:      "operator-1",
				QuorumNumbers: []uint8{0x00},
				BlockNumber:   12,
			},
		},
	}
	_, err = dataapi.CreateOperatorQuorumIntervals(10, 25, operatorInitialQuorum, addedQuorums, removedQuorums)
	assert.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "adding and removing quorums at the same block"))

	// Adding existing quorum again
	addedQuorums = map[string][]*dataapi.OperatorQuorum{
		"operator-1": []*dataapi.OperatorQuorum{
			{
				Operator:      "operator-1",
				QuorumNumbers: []uint8{0x00},
				BlockNumber:   11,
			},
		},
	}
	_, err = dataapi.CreateOperatorQuorumIntervals(10, 25, operatorInitialQuorum, addedQuorums, removedQuorums)
	assert.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "it is already in the quorum"))

	// Removing nonexisting quorum
	addedQuorums = map[string][]*dataapi.OperatorQuorum{
		"operator-1": []*dataapi.OperatorQuorum{
			{
				Operator:      "operator-1",
				QuorumNumbers: []uint8{0x02},
				BlockNumber:   12,
			},
		},
	}
	removedQuorums = map[string][]*dataapi.OperatorQuorum{
		"operator-1": []*dataapi.OperatorQuorum{
			{
				Operator:      "operator-1",
				QuorumNumbers: []uint8{0x01},
				BlockNumber:   11,
			},
		},
	}
	_, err = dataapi.CreateOperatorQuorumIntervals(10, 25, operatorInitialQuorum, addedQuorums, removedQuorums)
	assert.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "cannot remove a quorum"))
}

func TestCreateOperatorQuorumIntervalsWithNoQuorumChanges(t *testing.T) {
	addedQuorums := map[string][]*dataapi.OperatorQuorum{}
	removedQuorums := map[string][]*dataapi.OperatorQuorum{}
	operatorInitialQuorum := map[string][]uint8{
		"operator-1": {0x00},
		"operator-2": {0x01},
	}
	quorumIntervals, err := dataapi.CreateOperatorQuorumIntervals(10, 25, operatorInitialQuorum, addedQuorums, removedQuorums)
	assert.NoError(t, err)

	assert.Equal(t, 2, len(quorumIntervals))
	expectedOp1 := map[uint8][]dataapi.BlockInterval{0: []dataapi.BlockInterval{
		{
			StartBlock: 10,
			EndBlock:   25,
		},
	},
	}
	assertEntry(t, quorumIntervals, "operator-1", expectedOp1)
	expectedOp2 := map[uint8][]dataapi.BlockInterval{
		1: []dataapi.BlockInterval{
			{
				StartBlock: 10,
				EndBlock:   25,
			},
		},
	}
	assertEntry(t, quorumIntervals, "operator-2", expectedOp2)
}

func TestCreateOperatorQuorumIntervalsWithOnlyAddOrRemove(t *testing.T) {
	addedQuorums := map[string][]*dataapi.OperatorQuorum{
		"operator-1": []*dataapi.OperatorQuorum{
			{
				Operator:      "operator-1",
				QuorumNumbers: []uint8{0x01},
				BlockNumber:   11,
			},
			{
				Operator:      "operator-1",
				QuorumNumbers: []uint8{0x02, 0x03},
				BlockNumber:   20,
			},
		},
		"operator-2": []*dataapi.OperatorQuorum{
			{
				Operator:      "operator-2",
				QuorumNumbers: []byte{0x01, 0x02},
				BlockNumber:   25,
			},
		},
	}
	removedQuorums := map[string][]*dataapi.OperatorQuorum{}
	operatorInitialQuorum := map[string][]uint8{
		"operator-1": {0x00},
		"operator-2": {0x00},
	}

	quorumIntervals, err := dataapi.CreateOperatorQuorumIntervals(10, 25, operatorInitialQuorum, addedQuorums, removedQuorums)
	assert.NoError(t, err)

	assert.Equal(t, 2, len(quorumIntervals))
	expectedOp1 := map[uint8][]dataapi.BlockInterval{
		0: []dataapi.BlockInterval{
			{
				StartBlock: 10,
				EndBlock:   25,
			},
		},
		1: []dataapi.BlockInterval{
			{
				StartBlock: 11,
				EndBlock:   25,
			},
		},
		2: []dataapi.BlockInterval{
			{
				StartBlock: 20,
				EndBlock:   25,
			},
		},
		3: []dataapi.BlockInterval{
			{
				StartBlock: 20,
				EndBlock:   25,
			},
		},
	}
	assertEntry(t, quorumIntervals, "operator-1", expectedOp1)

	expectedOp2 := map[uint8][]dataapi.BlockInterval{
		0: []dataapi.BlockInterval{
			{
				StartBlock: 10,
				EndBlock:   25,
			},
		},
		1: []dataapi.BlockInterval{
			{
				StartBlock: 25,
				EndBlock:   25,
			},
		},
		2: []dataapi.BlockInterval{
			{
				StartBlock: 25,
				EndBlock:   25,
			},
		},
	}
	assertEntry(t, quorumIntervals, "operator-2", expectedOp2)

	addedQuorums = map[string][]*dataapi.OperatorQuorum{}
	removedQuorums = map[string][]*dataapi.OperatorQuorum{
		"operator-1": []*dataapi.OperatorQuorum{
			{
				Operator:      "operator-1",
				QuorumNumbers: []uint8{0x00},
				BlockNumber:   15,
			},
		},
	}
	quorumIntervals, err = dataapi.CreateOperatorQuorumIntervals(10, 25, operatorInitialQuorum, addedQuorums, removedQuorums)
	assert.NoError(t, err)
	expectedOp3 := map[uint8][]dataapi.BlockInterval{
		0: []dataapi.BlockInterval{
			{
				StartBlock: 10,
				EndBlock:   14,
			},
		},
	}
	assertEntry(t, quorumIntervals, "operator-1", expectedOp3)
}

func TestCreateOperatorQuorumIntervals(t *testing.T) {
	addedQuorums := map[string][]*dataapi.OperatorQuorum{
		"operator-1": []*dataapi.OperatorQuorum{
			{
				Operator:      "operator-1",
				QuorumNumbers: []uint8{0x01},
				BlockNumber:   11,
			},
			{
				Operator:      "operator-1",
				QuorumNumbers: []uint8{0x02, 0x03},
				BlockNumber:   20,
			},
			{
				Operator:      "operator-1",
				QuorumNumbers: []uint8{0x00},
				BlockNumber:   20,
			},
		},
		"operator-2": []*dataapi.OperatorQuorum{
			{
				Operator:      "operator-2",
				QuorumNumbers: []byte{0x02},
				BlockNumber:   15,
			},
			{
				Operator:      "operator-2",
				QuorumNumbers: []byte{0x02},
				BlockNumber:   22,
			},
		},
	}
	removedQuorums := map[string][]*dataapi.OperatorQuorum{
		"operator-1": []*dataapi.OperatorQuorum{
			{
				Operator:      "operator-1",
				QuorumNumbers: []uint8{0x00},
				BlockNumber:   15,
			},
			{
				Operator:      "operator-1",
				QuorumNumbers: []uint8{0x02},
				BlockNumber:   21,
			},
			{
				Operator:      "operator-1",
				QuorumNumbers: []uint8{0x00},
				BlockNumber:   23,
			},
		},
		"operator-2": []*dataapi.OperatorQuorum{
			{
				Operator:      "operator-2",
				QuorumNumbers: []byte{0x01, 0x02},
				BlockNumber:   20,
			},
		},
	}
	operatorInitialQuorum := map[string][]uint8{
		"operator-1": {0x00},
		"operator-2": {0x00, 0x01},
	}

	quorumIntervals, err := dataapi.CreateOperatorQuorumIntervals(10, 25, operatorInitialQuorum, addedQuorums, removedQuorums)
	assert.NoError(t, err)

	assert.Equal(t, 2, len(quorumIntervals))
	expectedOp1 := map[uint8][]dataapi.BlockInterval{
		0: []dataapi.BlockInterval{
			{
				StartBlock: 10,
				EndBlock:   14,
			},
			{
				StartBlock: 20,
				EndBlock:   22,
			},
		},
		1: []dataapi.BlockInterval{
			{
				StartBlock: 11,
				EndBlock:   25,
			},
		},
		2: []dataapi.BlockInterval{
			{
				StartBlock: 20,
				EndBlock:   20,
			},
		},
		3: []dataapi.BlockInterval{
			{
				StartBlock: 20,
				EndBlock:   25,
			},
		},
	}
	assertEntry(t, quorumIntervals, "operator-1", expectedOp1)
	assert.ElementsMatch(t, []uint8{0x00}, quorumIntervals.GetQuorums("operator-1", 10))
	assert.ElementsMatch(t, []uint8{0x00, 0x01}, quorumIntervals.GetQuorums("operator-1", 11))
	assert.ElementsMatch(t, []uint8{0x01}, quorumIntervals.GetQuorums("operator-1", 15))
	assert.ElementsMatch(t, []uint8{0x00, 0x01, 0x02, 0x03}, quorumIntervals.GetQuorums("operator-1", 20))
	assert.ElementsMatch(t, []uint8{0x00, 0x01, 0x03}, quorumIntervals.GetQuorums("operator-1", 22))
	assert.ElementsMatch(t, []uint8{0x01, 0x03}, quorumIntervals.GetQuorums("operator-1", 23))
	assert.ElementsMatch(t, []uint8{0x01, 0x03}, quorumIntervals.GetQuorums("operator-1", 25))

	expectedOp2 := map[uint8][]dataapi.BlockInterval{
		0: []dataapi.BlockInterval{
			{
				StartBlock: 10,
				EndBlock:   25,
			},
		},
		1: []dataapi.BlockInterval{
			{
				StartBlock: 10,
				EndBlock:   19,
			},
		},
		2: []dataapi.BlockInterval{
			{
				StartBlock: 15,
				EndBlock:   19,
			},
			{
				StartBlock: 22,
				EndBlock:   25,
			},
		},
	}
	assertEntry(t, quorumIntervals, "operator-2", expectedOp2)
	assert.ElementsMatch(t, []uint8{0x00, 0x01}, quorumIntervals.GetQuorums("operator-2", 10))
	assert.ElementsMatch(t, []uint8{0x00, 0x01, 0x02}, quorumIntervals.GetQuorums("operator-2", 15))
	assert.ElementsMatch(t, []uint8{0x00}, quorumIntervals.GetQuorums("operator-2", 20))
	assert.ElementsMatch(t, []uint8{0x00, 0x02}, quorumIntervals.GetQuorums("operator-2", 22))
	assert.ElementsMatch(t, []uint8{0x00, 0x02}, quorumIntervals.GetQuorums("operator-2", 25))
}
