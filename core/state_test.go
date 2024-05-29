package core_test

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/stretchr/testify/assert"
)

func TestOperatorStateHash(t *testing.T) {
	s1 := core.OperatorState{
		Operators: map[core.QuorumID]map[core.OperatorID]*core.OperatorInfo{
			0: {
				[32]byte{0}: &core.OperatorInfo{
					Stake: big.NewInt(12),
					Index: uint(2),
				},
				[32]byte{1}: &core.OperatorInfo{
					Stake: big.NewInt(23),
					Index: uint(3),
				},
			},
			1: {
				[32]byte{1}: &core.OperatorInfo{
					Stake: big.NewInt(23),
					Index: uint(3),
				},
				[32]byte{2}: &core.OperatorInfo{
					Stake: big.NewInt(34),
					Index: uint(4),
				},
			},
		},
		Totals: map[core.QuorumID]*core.OperatorInfo{
			0: {
				Stake: big.NewInt(35),
				Index: uint(2),
			},
			1: {
				Stake: big.NewInt(57),
				Index: uint(2),
			},
		},
		BlockNumber: uint(123),
	}

	hash1, err := s1.Hash()
	assert.NoError(t, err)
	q0 := hash1[0]
	q1 := hash1[1]
	assert.Equal(t, "3805338f34f77ff1fa23bbc23b1e86c4", hex.EncodeToString(q0[:]))
	assert.Equal(t, "2f110a29f2bdd8a19c2d87d05736be0a", hex.EncodeToString(q1[:]))

	s2 := core.OperatorState{
		Operators: map[core.QuorumID]map[core.OperatorID]*core.OperatorInfo{
			0: {
				[32]byte{0}: &core.OperatorInfo{
					Stake: big.NewInt(12),
					Index: uint(3), // different from s1
				},
				[32]byte{1}: &core.OperatorInfo{
					Stake: big.NewInt(23),
					Index: uint(3),
				},
			},
			1: {
				[32]byte{1}: &core.OperatorInfo{
					Stake: big.NewInt(23),
					Index: uint(3),
				},
				[32]byte{2}: &core.OperatorInfo{
					Stake: big.NewInt(34),
					Index: uint(4),
				},
			},
		},
		Totals: map[core.QuorumID]*core.OperatorInfo{
			0: {
				Stake: big.NewInt(35),
				Index: uint(2),
			},
			1: {
				Stake: big.NewInt(57),
				Index: uint(2),
			},
		},
		BlockNumber: uint(123),
	}

	hash2, err := s2.Hash()
	assert.NoError(t, err)
	q0 = hash2[0]
	q1 = hash2[1]
	assert.Equal(t, "1836448b57ae79decdcb77157cf31698", hex.EncodeToString(q0[:]))
	assert.Equal(t, "2f110a29f2bdd8a19c2d87d05736be0a", hex.EncodeToString(q1[:]))
}
