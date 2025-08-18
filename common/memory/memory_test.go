package memory

import (
	"fmt"
	"testing"

	"github.com/docker/go-units"
	"github.com/stretchr/testify/require"
)

func TestGetMaximumAvailableMemory(t *testing.T) {
	memory, err := GetMaximumAvailableMemory()
	require.NoError(t, err)

	// Since the outcome of this test depends on the environment, we can only check if the value is greater than 0.
	// This test is mostly intended designed for manual verification, although it does at least verify that the
	// function does not return an error.
	fmt.Printf("Maximum available memory: %dGB\n", memory/units.GiB)
	require.Greater(t, memory, uint64(0), "Memory should be greater than 0")

}
