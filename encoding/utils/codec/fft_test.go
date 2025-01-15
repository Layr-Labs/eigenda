package codec

import (
	"bytes"
	"testing"

	"github.com/Layr-Labs/eigenda/common/testutils/random"
	"github.com/stretchr/testify/require"
)

// TestFFT checks that data can be IFFTed and FFTed repeatedly, always getting back the original data
func TestFFT(t *testing.T) {
	testRandom := random.NewTestRandom(t)

	// Number of test iterations
	iterations := testRandom.Intn(100) + 1

	for i := 0; i < iterations; i++ {
		originalData := testRandom.Bytes(testRandom.Intn(1024) + 1) // ensure it's not length 0

		encodedData := EncodePayload(originalData)
		coeffPoly, err := IFFT(encodedData)
		require.NoError(t, err)

		evalPoly, err := FFT(coeffPoly)
		require.NoError(t, err)

		// Decode the encoded data
		decodedData, err := DecodePayload(evalPoly)
		if err != nil {
			t.Fatalf("Iteration %d: failed to decode blob: %v", i, err)
		}

		// Compare the original data with the decoded data
		if !bytes.Equal(originalData, decodedData) {
			t.Fatalf("Iteration %d: original and decoded data do not match\nOriginal: %v\nDecoded: %v", i, originalData, decodedData)
		}
	}
}
