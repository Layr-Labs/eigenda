package codecs

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

// FuzzConversionConsistency checks that data can be encoded and decoded repeatedly, always getting back the original data
func FuzzConversionConsistency(f *testing.F) {
	for _, seed := range [][]byte{{}, {0x00}, {0xFF}, {0x00, 0x00}, {0xFF, 0xFF}, bytes.Repeat([]byte{0x55}, 1000)} {
		f.Add(seed)
	}

	f.Fuzz(
		func(t *testing.T, originalData []byte) {
			payload := NewPayload(originalData)

			blob1, err := payload.ToBlob(PolynomialFormEval)
			require.NoError(t, err)

			blob2, err := payload.ToBlob(PolynomialFormCoeff)
			require.NoError(t, err)

			decodedPayload1, err := blob1.ToPayload(PolynomialFormEval)
			require.NoError(t, err)

			decodedPayload2, err := blob2.ToPayload(PolynomialFormCoeff)
			require.NoError(t, err)

			require.Equal(t, originalData, decodedPayload1.GetBytes())
			require.Equal(t, originalData, decodedPayload2.GetBytes())
		})
}
