package node_test

import (
	"testing"

	"github.com/Layr-Labs/eigenda/node"
	"github.com/stretchr/testify/assert"
)

func TestDecodeHashSlice(t *testing.T) {
	hash0 := [32]byte{0, 1}
	hash1 := [32]byte{0, 1, 2, 3, 4}
	hash2 := [32]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31}

	input := make([]byte, 0)
	input = append(input, hash0[:]...)
	input = append(input, hash1[:]...)
	input = append(input, hash2[:]...)

	hashes, err := node.DecodeHashSlice(input)
	assert.NoError(t, err)
	assert.Len(t, hashes, 3)
	assert.Equal(t, hash0, hashes[0])
	assert.Equal(t, hash1, hashes[1])
	assert.Equal(t, hash2, hashes[2])
}
