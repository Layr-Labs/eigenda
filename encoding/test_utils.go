package encoding

import (
	"golang.org/x/exp/constraints"
)

// GeneratePowersOfTwo creates a slice of integers, containing powers of 2 (starting with element == 1), with
// powersToGenerate number of elements
func GeneratePowersOfTwo[T constraints.Integer](powersToGenerate T) []T {
	powers := make([]T, powersToGenerate)
	for i := T(0); i < powersToGenerate; i++ {
		powers[i] = 1 << i
	}

	return powers
}
