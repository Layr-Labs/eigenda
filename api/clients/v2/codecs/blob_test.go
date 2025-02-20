package codecs

import (
	"bytes"
	"testing"

	"github.com/Layr-Labs/eigenda/api/clients/codecs"
	"github.com/stretchr/testify/require"
)

// TestBlobConversion checks that internal blob conversion methods produce consistent results
func FuzzBlobConversion(f *testing.F) {
	for _, seed := range [][]byte{{}, {0x00}, {0xFF}, {0x00, 0x00}, {0xFF, 0xFF}, bytes.Repeat([]byte{0x55}, 1000)} {
		f.Add(seed)
	}

	f.Fuzz(
		func(t *testing.T, originalData []byte) {
			testBlobConversionForForm(t, originalData, codecs.PolynomialFormEval)
			testBlobConversionForForm(t, originalData, codecs.PolynomialFormCoeff)
		})

}

func testBlobConversionForForm(t *testing.T, payloadBytes []byte, payloadForm codecs.PolynomialForm) {
	blob, err := NewPayload(payloadBytes).ToBlob(payloadForm)
	require.NoError(t, err)

	blobDeserialized, err := DeserializeBlob(blob.Serialize(), blob.blobLengthSymbols)
	require.NoError(t, err)

	payloadFromBlob, err := blob.ToPayload(payloadForm)
	require.NoError(t, err)

	payloadFromDeserializedBlob, err := blobDeserialized.ToPayload(payloadForm)
	require.NoError(t, err)

	require.Equal(t, payloadFromBlob.Serialize(), payloadFromDeserializedBlob.Serialize())
	require.Equal(t, payloadBytes, payloadFromBlob.Serialize())
}
