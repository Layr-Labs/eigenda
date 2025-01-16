package codecs

import (
	"bytes"
	"testing"

	"github.com/Layr-Labs/eigenda/common/testutils/random"
	"github.com/stretchr/testify/require"
)

// TestCodec tests the encoding and decoding of random byte streams
func TestPayloadEncoding(t *testing.T) {
	testRandom := random.NewTestRandom(t)

	// Number of test iterations
	const iterations = 100

	for i := 0; i < iterations; i++ {
		payload := NewPayload(testRandom.Bytes(testRandom.Intn(1024) + 1))
		encodedPayload, err := payload.encode()
		require.NoError(t, err)

		// Decode the encoded data
		decodedPayload, err := encodedPayload.decode()
		require.NoError(t, err)

		if err != nil {
			t.Fatalf("Iteration %d: failed to decode blob: %v", i, err)
		}

		// Compare the original data with the decoded data
		if !bytes.Equal(payload.GetBytes(), decodedPayload.GetBytes()) {
			t.Fatalf(
				"Iteration %d: original and decoded data do not match\nOriginal: %v\nDecoded: %v",
				i,
				payload.GetBytes(),
				decodedPayload.GetBytes())
		}
	}
}
