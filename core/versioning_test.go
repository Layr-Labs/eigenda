package core_test

import (
	"math/big"
	"testing"

	pbvalidator "github.com/Layr-Labs/eigenda/api/grpc/validator"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/stretchr/testify/assert"
)

func TestCalculateQuorumRolloutReadiness_Verbose(t *testing.T) {
	q0 := core.QuorumID(0)
	q1 := core.QuorumID(1)
	requiredVersion := "1.2.3"
	threshold := 0.8

	// Helper to build OperatorInfoVerbose
	makeOp := func(id core.OperatorID, stake int64, semver *string) core.OperatorInfoVerbose {
		var nodeInfo *pbvalidator.GetNodeInfoReply
		if semver != nil {
			nodeInfo = &pbvalidator.GetNodeInfoReply{Semver: *semver}
		}
		return core.OperatorInfoVerbose{
			OperatorID: id,
			Stake:      big.NewInt(stake),
			NodeInfo:   nodeInfo,
		}
	}

	// 3 operators: 2 in q0, 2 in q1, 1 in both
	opA := core.OperatorID{0x01}
	opB := core.OperatorID{0x02}
	opC := core.OperatorID{0x03}
	ops := core.OperatorStateVerbose{
		q0: {
			0: makeOp(opA, 100, &requiredVersion),
			1: makeOp(opB, 100, strPtr("old")),
		},
		q1: {
			0: makeOp(opB, 100, strPtr("old")),
			1: makeOp(opC, 200, &requiredVersion),
		},
	}
	pctByQuorum, readyByQuorum := core.CalculateQuorumRolloutReadiness(ops, requiredVersion, threshold)
	assert.Equal(t, 0.5, pctByQuorum[q0])
	assert.False(t, readyByQuorum[q0])
	assert.InDelta(t, 2.0/3.0, pctByQuorum[q1], 1e-6)
	assert.False(t, readyByQuorum[q1])

	// Now upgrade opB
	ops[q0][1] = makeOp(opB, 100, &requiredVersion)
	ops[q1][0] = makeOp(opB, 100, &requiredVersion)
	pctByQuorum, readyByQuorum = core.CalculateQuorumRolloutReadiness(ops, requiredVersion, threshold)
	assert.Equal(t, 1.0, pctByQuorum[q0])
	assert.True(t, readyByQuorum[q0])
	assert.Equal(t, 1.0, pctByQuorum[q1])
	assert.True(t, readyByQuorum[q1])

	// Edge case 1: No operators
	pctByQuorum, readyByQuorum = core.CalculateQuorumRolloutReadiness(core.OperatorStateVerbose{}, requiredVersion, threshold)
	assert.Equal(t, 0, len(pctByQuorum))
	assert.Equal(t, 0, len(readyByQuorum))

	// Edge case 2: Operator with zero stake
	opD := core.OperatorID{0x04}
	ops2 := core.OperatorStateVerbose{
		q0: {0: makeOp(opD, 0, &requiredVersion)},
	}
	pctByQuorum, readyByQuorum = core.CalculateQuorumRolloutReadiness(ops2, requiredVersion, threshold)
	assert.Equal(t, 0.0, pctByQuorum[q0])
	assert.False(t, readyByQuorum[q0])

	// Edge case 3: Operator with nil NodeInfo
	opE := core.OperatorID{0x05}
	ops3 := core.OperatorStateVerbose{
		q0: {0: makeOp(opE, 100, nil)},
	}
	pctByQuorum, readyByQuorum = core.CalculateQuorumRolloutReadiness(ops3, requiredVersion, threshold)
	assert.Equal(t, 0.0, pctByQuorum[q0])
	assert.False(t, readyByQuorum[q0])

	// Edge case 4: Threshold exactly met
	opF := core.OperatorID{0x06}
	opG := core.OperatorID{0x07}
	ops4 := core.OperatorStateVerbose{
		q0: {
			0: makeOp(opF, 80, &requiredVersion),
			1: makeOp(opG, 20, strPtr("old")),
		},
	}
	pctByQuorum, readyByQuorum = core.CalculateQuorumRolloutReadiness(ops4, requiredVersion, 0.8)
	assert.Equal(t, 0.8, pctByQuorum[q0])
	assert.True(t, readyByQuorum[q0])

	// Edge case 5: All upgraded except one with tiny stake
	opH := core.OperatorID{0x08}
	ops5 := core.OperatorStateVerbose{
		q0: {
			0: makeOp(opA, 100, &requiredVersion),
			1: makeOp(opB, 100, &requiredVersion),
			2: makeOp(opC, 200, &requiredVersion),
			3: makeOp(opH, 1, strPtr("old")),
		},
	}
	pctByQuorum, readyByQuorum = core.CalculateQuorumRolloutReadiness(ops5, requiredVersion, 0.999)
	assert.True(t, pctByQuorum[q0] < 0.999)
	assert.False(t, readyByQuorum[q0])

	// Edge case 6: Operator in multiple quorums, only upgraded in some
	opI := core.OperatorID{0x09}
	ops6 := core.OperatorStateVerbose{
		q0: {
			0: makeOp(opI, 100, &requiredVersion),
			1: makeOp(opB, 100, strPtr("old")),
		},
		q1: {
			0: makeOp(opB, 100, strPtr("old")),
			1: makeOp(opC, 200, &requiredVersion),
		},
	}
	pctByQuorum, readyByQuorum = core.CalculateQuorumRolloutReadiness(ops6, requiredVersion, threshold)
	assert.Equal(t, 0.5, pctByQuorum[q0])
	assert.InDelta(t, 2.0/3.0, pctByQuorum[q1], 1e-6)
	assert.False(t, readyByQuorum[q0])
	assert.False(t, readyByQuorum[q1])
}

func strPtr(s string) *string { return &s }
