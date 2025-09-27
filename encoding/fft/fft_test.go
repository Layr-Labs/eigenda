package fft_test

import (
	"testing"

	"github.com/Layr-Labs/eigenda/encoding/fft"
	"github.com/stretchr/testify/require"
)

const (
	// Change this to benchmark different maxScales.
	maxScale = uint8(22) // 2^22 * 32 = 128MiB
)

// BenchmarkFFTSettings benchmarks the creation of FFTSettings for a given maxScale.
// This maxScale of 22 allows FFTs of up to 128MiB (2^22 * 32 bytes).
// This in turn allows blobs of up to 16MiB, given that our RS encoding uses a 8x expansion
// for blob version 0.
//
// The main thing we are interested in here is the memory allocation,
// to make sure that we smartly allocate the arrays for the roots of unity.
// See [TestFFTSettingsBytesAllocation] below.
func BenchmarkFFTSettings(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		_ = fft.NewFFTSettings(maxScale)
	}
}

// TestFFTSettingsBytesAllocation tests that the FFTSettings creation
// allocates a reasonable amount of memory, given the maxScale.
// We expect at least 2 arrays of size 2^maxScale * 32 bytes (roots of unity and reverse roots of unity).
// We allow an extra 5MiB for overhead.
func TestFFTSettingsBytesAllocation(t *testing.T) {
	numElements := int64(1 << maxScale)
	numBytes := numElements * 32
	// 2 arrays of size numBytes (roots of unity and reverse roots of unity)
	minExpectedAllocBytes := 2 * numBytes
	fiveMiB := int64(5 << 20)
	// We allow an extra 5MiB for overhead.
	maxExpectedAllocBytes := minExpectedAllocBytes + fiveMiB

	result := testing.Benchmark(BenchmarkFFTSettings)
	allocatedBytes := result.AllocedBytesPerOp()
	require.GreaterOrEqual(t, allocatedBytes, minExpectedAllocBytes,
		"expected at least %d bytes allocated, got %d", minExpectedAllocBytes, allocatedBytes)
	require.Less(t, allocatedBytes, maxExpectedAllocBytes,
		"expected less than %d bytes allocated, got %d", maxExpectedAllocBytes, allocatedBytes)
}
