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
					Stake:  big.NewInt(12),
					Index:  uint(2),
					Socket: core.NewOperatorSocket("192.168.1.100", "8080", "", "", ""),
				},
				[32]byte{1}: &core.OperatorInfo{
					Stake:  big.NewInt(23),
					Index:  uint(3),
					Socket: core.NewOperatorSocket("127.0.0.1", "3000", "", "", ""),
				},
			},
			1: {
				[32]byte{1}: &core.OperatorInfo{
					Stake:  big.NewInt(23),
					Index:  uint(3),
					Socket: core.NewOperatorSocket("127.0.0.1", "3000", "", "", ""),
				},
				[32]byte{2}: &core.OperatorInfo{
					Stake:  big.NewInt(34),
					Index:  uint(4),
					Socket: core.NewOperatorSocket("192.168.1.100", "8080", "", "", ""),
				},
			},
		},
		Totals: map[core.QuorumID]*core.OperatorInfo{
			0: {
				Stake:  big.NewInt(35),
				Index:  uint(2),
				Socket: core.OperatorSocket{},
			},
			1: {
				Stake:  big.NewInt(57),
				Index:  uint(2),
				Socket: core.OperatorSocket{},
			},
		},
		BlockNumber: uint(123),
	}

	hash1, err := s1.Hash()
	assert.NoError(t, err)
	q0 := hash1[0]
	q1 := hash1[1]
	assert.Equal(t, "7f227566a0f077dd9bba7613a40de260", hex.EncodeToString(q0[:]))
	assert.Equal(t, "aa1da3df03ce4c71fae05d5e6c957748", hex.EncodeToString(q1[:]))

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
	assert.Equal(t, "a382adf83f330be719af7e2cdaad0ebe", hex.EncodeToString(q0[:]))
	assert.Equal(t, "aa1da3df03ce4c71fae05d5e6c957748", hex.EncodeToString(q1[:]))
}
