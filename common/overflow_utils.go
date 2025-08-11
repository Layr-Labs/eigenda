package common

import (
	"fmt"
	"math"
)

// SafeAddInt64 performs addition with overflow detection for int64.
// Returns the result, or an error if the addition would overflow.
func SafeAddInt64(a, b int64) (int64, error) {
	// positive overflow
	if b > 0 && a > math.MaxInt64-b {
		return 0, fmt.Errorf("positive addition overflow: %d + %d", a, b)
	}

	// negative overflow
	if b < 0 && a < math.MinInt64-b {
		return 0, fmt.Errorf("negative addition overflow: %d + %d", a, b)
	}

	return a + b, nil
}

// SafeSubtractInt64 performs subtraction with overflow detection for int64.
// Returns the result, or an error if the subtraction would overflow.
func SafeSubtractInt64(a, b int64) (int64, error) {
	// positive overflow
	if b < 0 && a > math.MaxInt64+b {
		return 0, fmt.Errorf("positive subtraction overflow: %d - %d", a, b)
	}

	// negative overflow
	if b > 0 && a < math.MinInt64+b {
		return 0, fmt.Errorf("negative subtraction overflow: %d - %d", a, b)
	}

	return a - b, nil
}

// SafeMultiplyInt64 performs multiplication with overflow detection for int64.
// Returns the result, or an error if the multiplication would overflow.
func SafeMultiplyInt64(a, b int64) (int64, error) {
	// handle zero case early to avoid division by zero
	if a == 0 || b == 0 {
		return 0, nil
	}

	// positive overflow
	if (a > 0 && b > 0 && a > math.MaxInt64/b) ||
		(a < 0 && b < 0 && a < math.MaxInt64/b) {
		return 0, fmt.Errorf("positive multiplication overflow: %d * %d", a, b)
	}

	// negative overflow
	if (a > 0 && b < 0 && b < math.MinInt64/a) ||
		(a < 0 && b > 0 && a < math.MinInt64/b) {
		return 0, fmt.Errorf("negative multiplication overflow: %d * %d", a, b)
	}

	return a * b, nil
}
