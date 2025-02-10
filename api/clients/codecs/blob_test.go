package codecs

import (
	"testing"

	"github.com/Layr-Labs/eigenda/common/testutils/random"
	"github.com/stretchr/testify/require"
)

// TestBlobConversion checks that internal blob conversion methods produce consistent results
func TestBlobConversion(t *testing.T) {
	testRandom := random.NewTestRandom(t)

	iterations := 100

	for i := 0; i < iterations; i++ {
		originalData := testRandom.Bytes(testRandom.Intn(1024) + 1)
		testBlobConversionForForm(t, originalData, PolynomialFormEval)
		testBlobConversionForForm(t, originalData, PolynomialFormCoeff)
	}
}

func testBlobConversionForForm(t *testing.T, payloadBytes []byte, form PolynomialForm) {
	payload := NewPayload(payloadBytes)

	blob, err := payload.ToBlob(form)
	require.NoError(t, err)

	blobBytes := blob.GetBytes()
	blobFromBytes, err := BlobFromBytes(blobBytes)
	require.NoError(t, err)

	decodedPayload, err := blobFromBytes.ToPayload(form)
	require.NoError(t, err)

	decodedPayloadBytes := decodedPayload.GetBytes()

	require.Equal(t, payloadBytes, decodedPayloadBytes)
}
