package utils

import (
	"math"
	"testing"

	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/stretchr/testify/require"
)

// TestIsPowerOf2 checks that the IsPowerOfTwo utility is working as expected
func TestIsPowerOf2(t *testing.T) {
	// Test the special case
	is0PowerOf2 := encoding.IsPowerOfTwo(0)
	require.False(t, is0PowerOf2)

	testValue := uint32(1)
	require.True(t, encoding.IsPowerOfTwo(testValue), "expected %d to be a valid power of 2", testValue)

	for testValue < math.MaxUint32 / 2 {
		testValue = testValue * 2
		require.True(t, encoding.IsPowerOfTwo(testValue), "expected %d to be a valid power of 2", testValue)
	}
}
