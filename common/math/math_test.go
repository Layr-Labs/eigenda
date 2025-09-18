package math

import (
	gomath "math"
	"testing"

	"github.com/stretchr/testify/assert"
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

		assert.Equal(t, expectedResult, result, "IsPowerOfTwo(%d) returned unexpected result '%t'.", i, result)
	}
}
