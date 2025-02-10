package codecs

import (
	"bytes"
	"testing"

	"github.com/Layr-Labs/eigenda/common/testutils/random"
	"github.com/stretchr/testify/require"
)

// TestConversionConsistency checks that data can be encoded and decoded repeatedly, always getting back the original data
// TODO: we should probably be using fuzzing instead of this kind of ad-hoc random search testing
func TestConversionConsistency(t *testing.T) {
	testRandom := random.NewTestRandom(t)

	iterations := 100

	for i := 0; i < iterations; i++ {
		originalData := testRandom.Bytes(testRandom.Intn(1024) + 1) // ensure it's not length 0

		payload := NewPayload(originalData)

		blob1, err := payload.ToBlob(PolynomialFormEval)
		require.NoError(t, err)

		blob2, err := payload.ToBlob(PolynomialFormCoeff)
		require.NoError(t, err)

		decodedPayload1, err := blob1.ToPayload(PolynomialFormEval)
		require.NoError(t, err)

		decodedPayload2, err := blob2.ToPayload(PolynomialFormCoeff)
		require.NoError(t, err)

		// Compare the original data with the decoded data
		if !bytes.Equal(originalData, decodedPayload1.GetBytes()) {
			t.Fatalf(
				"Iteration %d: original and data decoded from blob1 do not match\nOriginal: %v\nDecoded: %v",
				i,
				originalData,
				decodedPayload1.GetBytes())
		}

		if !bytes.Equal(originalData, decodedPayload2.GetBytes()) {
			t.Fatalf(
				"Iteration %d: original and data decoded from blob2 do not match\nOriginal: %v\nDecoded: %v",
				i,
				originalData,
				decodedPayload1.GetBytes())
		}
	}
}
