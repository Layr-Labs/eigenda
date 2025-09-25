package math

import (
	"math/bits"

	"golang.org/x/exp/constraints"
)

// IsPowerOfTwo checks if a number is a power of 2
func IsPowerOfTwo[T constraints.Integer](d T) bool {
	return (d != 0) && (d&(d-1) == 0)
}

func RoundUpDivide[T constraints.Integer](a, b T) T {
	return (a + b - 1) / b
}

// NextPowOf2u32 returns the next power of 2 greater than or equal to v
func NextPowOf2u32(v uint32) uint32 {
	if v == 0 {
		return 1
	}
	return uint32(1) << bits.Len32(v-1)
}

// NextPowOf2u64 returns the next power of 2 greater than or equal to v
func NextPowOf2u64(v uint64) uint64 {
	if v == 0 {
		return 1
	}
	return uint64(1) << bits.Len64(v-1)
}
