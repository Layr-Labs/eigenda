package utils_test

import (
	"fmt"
	"testing"

	"github.com/Layr-Labs/eigenda-proxy/utils"
)

func TestParseByteAmount(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		input    string
		expected uint64
		wantErr  bool
	}{
		{"10 B", 10, false},
		{"15 b", 15, false}, // Case-insensitive
		{"1 KiB", 1024, false},
		{"2 kib", 2048, false},  // Case-insensitive
		{"5 KB", 5000, false},   // Decimal kilobyte
		{"10 kb", 10000, false}, // Decimal kilobyte (case-insensitive)
		{"1 MiB", 1024 * 1024, false},
		{"3 mib", 3 * 1024 * 1024, false},
		{"10 MB", 10 * 1000 * 1000, false},
		{"100 mb", 100 * 1000 * 1000, false},
		{"1 GiB", 1024 * 1024 * 1024, false},
		{"10 gib", 10 * 1024 * 1024 * 1024, false},
		{"10 GB", 10 * 1000 * 1000 * 1000, false},
		{"100 gb", 100 * 1000 * 1000 * 1000, false},
		{"1 TiB", 1024 * 1024 * 1024 * 1024, false},
		{"10 tib", 10 * 1024 * 1024 * 1024 * 1024, false},
		{"1 TB", 1000 * 1000 * 1000 * 1000, false},
		{"100 tb", 100 * 1000 * 1000 * 1000 * 1000, false},

		{"   5 B   ", 5, false}, // Whitespace handling
		{"10", 10, false},       // Default to bytes if no unit

		{"10 XB", 0, true}, // Invalid unit
		{"abc", 0, true},   // Non-numeric value
		{"1.5 KiB", 1536, false},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Input: %s", tc.input), func(t *testing.T) {
			got, err := utils.ParseBytesAmount(tc.input)
			if (err != nil) != tc.wantErr {
				t.Errorf("wantErr: %v, got error: %v", tc.wantErr, err)
			}
			if got != tc.expected {
				t.Errorf("got: %d, expected: %d", got, tc.expected)
			}
		})
	}
}
