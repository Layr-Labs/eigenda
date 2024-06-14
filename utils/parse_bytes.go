package utils

import (
	"fmt"
	"strconv"
	"strings"
)

// Helper utility functions //

func EqualSlices[P comparable](s1, s2 []P) bool {
	if len(s1) != len(s2) {
		return false
	}

	for i := 0; i < len(s1); i++ {
		if s1[i] != s2[i] {
			return false
		}
	}

	return true
}

func ParseBytesAmount(s string) (uint64, error) {
	s = strings.TrimSpace(s)

	// Extract numeric part and unit
	numStr := s
	unit := ""
	for i, r := range s {
		if !('0' <= r && r <= '9' || r == '.') {
			numStr = s[:i]
			unit = s[i:]
			break
		}
	}

	// Convert numeric part to float64
	num, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid numeric value: %v", err)
	}

	unit = strings.ToLower(strings.TrimSpace(unit))

	// Convert to uint64 based on the unit (case-insensitive)
	switch unit {
	case "b", "":
		return uint64(num), nil
	case "kib":
		return uint64(num * 1024), nil
	case "kb":
		return uint64(num * 1000), nil // Decimal kilobyte
	case "mib":
		return uint64(num * 1024 * 1024), nil
	case "mb":
		return uint64(num * 1000 * 1000), nil // Decimal megabyte
	case "gib":
		return uint64(num * 1024 * 1024 * 1024), nil
	case "gb":
		return uint64(num * 1000 * 1000 * 1000), nil // Decimal gigabyte
	case "tib":
		return uint64(num * 1024 * 1024 * 1024 * 1024), nil
	case "tb":
		return uint64(num * 1000 * 1000 * 1000 * 1000), nil // Decimal terabyte
	default:
		return 0, fmt.Errorf("unsupported unit: %s", unit)
	}
}
