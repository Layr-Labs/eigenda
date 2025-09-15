package encoding

import (
	"testing"

	"github.com/Layr-Labs/eigenda/common/math"
	"github.com/stretchr/testify/require"
)

func TestNextPowerOf2(t *testing.T) {
	testHeight := uint64(65536)

	// 2 ^ 16 = 65536
	// i.e., the last element generated here == testHeight
	powers := GeneratePowersOfTwo(uint64(17))

	powerIndex := 0
	for i := uint64(1); i <= testHeight; i++ {
		nextPowerOf2 := math.NextPowOf2u64(i)
		require.Equal(t, nextPowerOf2, powers[powerIndex])

		if i == powers[powerIndex] {
			powerIndex++
		}
	}

	// sanity check the test logic
	require.Equal(t, powerIndex, len(powers))

	// extra sanity check, since we *really* rely on NextPowerOf2 returning
	// the same value, if it's already a power of 2
	require.Equal(t, uint64(16), math.NextPowOf2u64(16))
}
