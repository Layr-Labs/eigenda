package encoding

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNextPowerOf2(t *testing.T) {
	testHeight := 65536

	// 2 ^ 16 = 65536
	powersToGenerate := 17
	powers := make([]int, powersToGenerate)
	for i := 0; i < powersToGenerate; i++ {
		powers[i] = int(math.Pow(2, float64(i)))
	}

	powerIndex := 0
	for i := 1; i <= testHeight; i++ {
		nextPowerOf2 := NextPowerOf2(i)
		require.Equal(t, nextPowerOf2, powers[powerIndex])

		if i == powers[powerIndex] {
			powerIndex++
		}
	}

	// sanity check the test logic
	require.Equal(t, powerIndex, len(powers))

	// extra sanity check
	require.Equal(t, 16, NextPowerOf2(16))
}
