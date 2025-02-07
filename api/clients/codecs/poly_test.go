package codecs

import (
	"bytes"
	"testing"

	"github.com/Layr-Labs/eigenda/common/testutils/random"
	"github.com/stretchr/testify/require"
)

// TestFftEncode checks that data can be IfftEncoded and FftEncoded repeatedly, always getting back the original data
// TODO: we should probably be using fuzzing instead of this kind of ad-hoc random search testing
func TestFftEncode(t *testing.T) {
	testRandom := random.NewTestRandom(t)

	// Number of test iterations
	iterations := 100

	for i := 0; i < iterations; i++ {
		originalData := testRandom.Bytes(testRandom.Intn(1024) + 1) // ensure it's not length 0

		payload := NewPayload(originalData)
		encodedPayload, err := newEncodedPayload(payload)
		require.NoError(t, err)

		evalPoly, err := encodedPayload.toEvalPoly()
		require.NoError(t, err)

		coeffPoly, err := evalPoly.toCoeffPoly()
		require.NoError(t, err)

		convertedEvalPoly, err := coeffPoly.toEvalPoly()
		require.NoError(t, err)

		convertedEncodedPayload, err := convertedEvalPoly.toEncodedPayload()
		require.NoError(t, err)

		decodedPayload, err := convertedEncodedPayload.decode()
		require.NoError(t, err)

		// Compare the original data with the decoded data
		if !bytes.Equal(originalData, decodedPayload.GetBytes()) {
			t.Fatalf(
				"Iteration %d: original and decoded data do not match\nOriginal: %v\nDecoded: %v",
				i,
				originalData,
				decodedPayload.GetBytes())
		}
	}
}
