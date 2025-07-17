package common_test

import (
	"fmt"
	"testing"

	"github.com/Layr-Labs/eigenda/api/proxy/common"
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
		{"2 kib", 2048, false}, // Case-insensitive
		{"1 MiB", 1024 * 1024, false},
		{"3 mib", 3 * 1024 * 1024, false},

		{"   5 B   ", 5, false}, // Whitespace handling
		{"10", 10, false},       // Default to bytes if no unit

		{"10 XB", 0, true}, // Invalid unit
		{"abc", 0, true},   // Non-numeric value
		{"1.5 KiB", 1536, false},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Input: %s", tc.input), func(t *testing.T) {
			t.Parallel()

			got, err := common.ParseBytesAmount(tc.input)
			if (err != nil) != tc.wantErr {
				t.Errorf("wantErr: %v, got error: %v", tc.wantErr, err)
			}
			if got != tc.expected {
				t.Errorf("got: %d, expected: %d", got, tc.expected)
			}
		})
	}
}
