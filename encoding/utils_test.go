package encoding

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNextPowerOf2(t *testing.T) {
	testHeight := 65536

	// 2 ^ 16 = 65536
	// i.e., the last element generated here == testHeight
	powers := GeneratePowersOfTwo(17)

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

	// extra sanity check, since we *really* rely on NextPowerOf2 returning
	// the same value, if it's already a power of 2
	require.Equal(t, 16, NextPowerOf2(16))
}
