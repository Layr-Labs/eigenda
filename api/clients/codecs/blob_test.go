package codecs

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestBlobConversion checks that internal blob conversion methods produce consistent results
func FuzzBlobConversion(f *testing.F) {
	for _, seed := range [][]byte{{}, {0x00}, {0xFF}, {0x00, 0x00}, {0xFF, 0xFF}, bytes.Repeat([]byte{0x55}, 1000)} {
		f.Add(seed)
	}

	f.Fuzz(
		func(t *testing.T, originalData []byte) {
			testBlobConversionForForm(t, originalData, PolynomialFormEval)
			testBlobConversionForForm(t, originalData, PolynomialFormCoeff)
		})

}

func testBlobConversionForForm(t *testing.T, payloadBytes []byte, form PolynomialForm) {
	payload := NewPayload(payloadBytes)

	blob, err := payload.ToBlob(form)
	require.NoError(t, err)

	blobBytes := blob.GetBytes()
	blobFromBytes, err := BlobFromBytes(blobBytes, blob.blobLength)
	require.NoError(t, err)

	decodedPayload, err := blobFromBytes.ToPayload(form)
	require.NoError(t, err)

	decodedPayloadBytes := decodedPayload.GetBytes()

	require.Equal(t, payloadBytes, decodedPayloadBytes)
}
