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

	// StartBlock > EndBlock
	operatorInitialQuorum := map[string][]uint8{
		"operator-1": {0x00},
		"operator-2": {0x00},
	}
	_, err := dataapi.CreateOperatorQuorumIntervals(100, 25, operatorInitialQuorum, addedQuorums, removedQuorums)
	assert.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "endBlock must be no less than startBlock"))

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
	assert.True(t, strings.Contains(err.Error(), "operator is already in the quorum"))

	// addedQuurums not in ascending order of block number
	addedQuorums = map[string][]*dataapi.OperatorQuorum{
		"operator-1": []*dataapi.OperatorQuorum{
			{
				Operator:      "operator-1",
				QuorumNumbers: []uint8{0x01},
				BlockNumber:   15,
			},
			{
				Operator:      "operator-1",
				QuorumNumbers: []uint8{0x03},
				BlockNumber:   11,
			},
		},
	}
	_, err = dataapi.CreateOperatorQuorumIntervals(10, 25, operatorInitialQuorum, addedQuorums, removedQuorums)
	assert.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "must be in ascending order by block number"))

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

func TestComputeNumBatches(t *testing.T) {
	quorumBatches := &dataapi.QuorumBatches{
		NumBatches:  []*dataapi.NumBatchesAtBlock{},
		AccuBatches: []int{},
	}
	assert.Equal(t, 0, dataapi.ComputeNumBatches(quorumBatches, 1, 4))

	numBatches := []*dataapi.NumBatchesAtBlock{
		{
			BlockNumber: 5,
			NumBatches:  2,
		},
	}
	quorumBatches = &dataapi.QuorumBatches{
		NumBatches:  numBatches,
		AccuBatches: []int{2},
	}
	assert.Equal(t, 0, dataapi.ComputeNumBatches(quorumBatches, 1, 4))
	assert.Equal(t, 2, dataapi.ComputeNumBatches(quorumBatches, 1, 5))
	assert.Equal(t, 2, dataapi.ComputeNumBatches(quorumBatches, 5, 5))
	assert.Equal(t, 2, dataapi.ComputeNumBatches(quorumBatches, 5, 6))

	numBatches = []*dataapi.NumBatchesAtBlock{
		{
			BlockNumber: 5,
			NumBatches:  2,
		},
		{
			BlockNumber: 10,
			NumBatches:  2,
		},
		{
			BlockNumber: 15,
			NumBatches:  2,
		},
		{
			BlockNumber: 20,
			NumBatches:  2,
		},
	}
	quorumBatches = &dataapi.QuorumBatches{
		NumBatches:  numBatches,
		AccuBatches: []int{2, 4, 6, 8},
	}

	assert.Equal(t, 0, dataapi.ComputeNumBatches(quorumBatches, 1, 4))
	assert.Equal(t, 0, dataapi.ComputeNumBatches(quorumBatches, 21, 22))
	assert.Equal(t, 2, dataapi.ComputeNumBatches(quorumBatches, 1, 5))
	assert.Equal(t, 2, dataapi.ComputeNumBatches(quorumBatches, 5, 5))
	assert.Equal(t, 2, dataapi.ComputeNumBatches(quorumBatches, 5, 9))
	assert.Equal(t, 4, dataapi.ComputeNumBatches(quorumBatches, 5, 10))
	assert.Equal(t, 2, dataapi.ComputeNumBatches(quorumBatches, 6, 10))
	assert.Equal(t, 4, dataapi.ComputeNumBatches(quorumBatches, 5, 14))
	assert.Equal(t, 2, dataapi.ComputeNumBatches(quorumBatches, 6, 14))
	assert.Equal(t, 6, dataapi.ComputeNumBatches(quorumBatches, 5, 15))
	assert.Equal(t, 8, dataapi.ComputeNumBatches(quorumBatches, 5, 20))
	assert.Equal(t, 8, dataapi.ComputeNumBatches(quorumBatches, 5, 22))
	assert.Equal(t, 8, dataapi.ComputeNumBatches(quorumBatches, 1, 22))
	assert.Equal(t, 6, dataapi.ComputeNumBatches(quorumBatches, 6, 22))
	assert.Equal(t, 4, dataapi.ComputeNumBatches(quorumBatches, 11, 22))
	assert.Equal(t, 2, dataapi.ComputeNumBatches(quorumBatches, 16, 22))
}

func TestCreatQuorumBatches(t *testing.T) {
	// The nonsigning info for a list of batches.
	batchNonSigningInfo := []*dataapi.BatchNonSigningInfo{
		{
			QuorumNumbers:        []uint8{0, 1},
			ReferenceBlockNumber: 2,
		},
		{
			QuorumNumbers:        []uint8{0},
			ReferenceBlockNumber: 2,
		},
		{
			QuorumNumbers:        []uint8{1, 2},
			ReferenceBlockNumber: 4,
		},
	}

	quorumBatches := dataapi.CreatQuorumBatches(batchNonSigningInfo)

	assert.Equal(t, 3, len(quorumBatches))

	q0, ok := quorumBatches[0]
	assert.True(t, ok)
	assert.Equal(t, 1, len(q0.NumBatches))
	assert.Equal(t, uint32(2), q0.NumBatches[0].BlockNumber)
	assert.Equal(t, 2, q0.AccuBatches[0])

	q1, ok := quorumBatches[1]
	assert.True(t, ok)
	assert.Equal(t, 2, len(q1.NumBatches))
	assert.Equal(t, uint32(2), q1.NumBatches[0].BlockNumber)
	assert.Equal(t, 1, q1.AccuBatches[0])
	assert.Equal(t, uint32(4), q1.NumBatches[1].BlockNumber)
	assert.Equal(t, 2, q1.AccuBatches[1])

	q2, ok := quorumBatches[2]
	assert.True(t, ok)
	assert.Equal(t, 1, len(q2.NumBatches))
	assert.Equal(t, uint32(4), q2.NumBatches[0].BlockNumber)
	assert.Equal(t, 1, q2.AccuBatches[0])
}
