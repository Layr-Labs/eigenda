package math

import (
	gomath "math"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/exp/constraints"
)

func TestIsPowerOfTwo(t *testing.T) {
	var i uint64
	for i = 0; i <= 1024; i++ {
		result := IsPowerOfTwo(i)

		var expectedResult bool
		if i == 0 {
			// Special case: gomath.Log2() is undefined for 0
			expectedResult = false
		} else {
			// If a number is not a power of two then the log base 2 of that number will not be a whole integer.
			logBase2 := gomath.Log2(float64(i))
			truncatedLogBase2 := float64(uint64(logBase2))
			expectedResult = logBase2 == truncatedLogBase2
		}

		require.Equal(t, expectedResult, result, "IsPowerOfTwo(%d) returned unexpected result '%t'.", i, result)
	}
}

func TestNextPowerOf2(t *testing.T) {
	testHeight := uint64(65536)

	// 2 ^ 16 = 65536
	// i.e., the last element generated here == testHeight
	powers := generatePowersOfTwo(uint64(17))

	powerIndex := 0
	for i := uint64(1); i <= testHeight; i++ {
		nextPowerOf2 := NextPowOf2u64(i)
		require.Equal(t, nextPowerOf2, powers[powerIndex])

		if i == powers[powerIndex] {
			powerIndex++
		}
	}

	// sanity check the test logic
	require.Equal(t, powerIndex, len(powers))

	// extra sanity check, since we *really* rely on NextPowerOf2 returning
	// the same value, if it's already a power of 2
	require.Equal(t, uint64(16), NextPowOf2u64(16))
}

// GeneratePowersOfTwo creates a slice of integers, containing powers of 2 (starting with element == 1), with
// powersToGenerate number of elements
func generatePowersOfTwo[T constraints.Integer](powersToGenerate T) []T {
	powers := make([]T, powersToGenerate)
	for i := T(0); i < powersToGenerate; i++ {
		powers[i] = 1 << i
	}

	return powers
}
