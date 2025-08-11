// nolint: lll // long lines are expected b/c of examples
package coretypes

import (
	"encoding/hex"
	"testing"

	"github.com/Layr-Labs/eigenda/api/clients/codecs"
	"github.com/Layr-Labs/eigenda/common/testutils/random"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/stretchr/testify/require"
)

// TestEncodePayload tests that the encoding of a Payload to an EncodedPayload works as expected.
func TestEncodeDecodePayload(t *testing.T) {

	// map of hex-encoded payloads (inputs) and their expected EncodedPayloads (outputs)
	testCases := map[string]string{
		// empty payload should only have a header symbol
		"": "0000000000000000000000000000000000000000000000000000000000000000",
		"01": "0000000000010000000000000000000000000000000000000000000000000000" + // header with len 1 payload
			"0001000000000000000000000000000000000000000000000000000000000000", // first byte is always 0 due to bn254 encoding
		"0102": "0000000000020000000000000000000000000000000000000000000000000000" +
			"0001020000000000000000000000000000000000000000000000000000000000",
		"01020304050607080910111213141516171819202122232425262728293031": "00000000001f0000000000000000000000000000000000000000000000000000" +
			"0001020304050607080910111213141516171819202122232425262728293031",
	}

	for payloadHex, expectedEncodedPayloadHex := range testCases {
		t.Run("EncodePayload "+payloadHex, func(t *testing.T) {
			payload, err := hex.DecodeString(payloadHex)
			require.NoError(t, err)
			encodedPayload := Payload(payload).ToEncodedPayload()
			require.NoError(t, encodedPayload.checkLenInvariant())
			require.Equal(t, expectedEncodedPayloadHex, hex.EncodeToString(encodedPayload.bytes))
			decodedPayload, err := encodedPayload.Decode()
			require.NoError(t, err)
			require.Equal(t, Payload(payload), decodedPayload)
		})
	}
}

func TestDecodePayloadErrors(t *testing.T) {
	testCases := []struct {
		name       string
		encodedHex string
	}{
		{
			name:       "Insufficient Length Doesn't Contain Header",
			encodedHex: "000000000000",
		},
		{
			name:       "First byte must be 0x00",
			encodedHex: "01000000000000000000000000000000000000000000000000000000",
		},
		{
			name:       "Only version 0x00 is supported",
			encodedHex: "00010000000000000000000000000000000000000000000000000000",
		},
		{
			name:       "Payload length must be a multiple of 32 bytes",
			encodedHex: "00000000000100000000000000000000000000000000000000000000" + "000100",
		},
		{
			name: "wrong payload length: 32 bytes of data, but header says 64",
			encodedHex: "00000000000200000000000000000000000000000000000000000000" +
				"00000000000000000000000000000000000000000000000000000000",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			bytes, err := hex.DecodeString(tc.encodedHex)
			require.NoError(t, err)

			encodedPayload := &EncodedPayload{bytes: bytes}
			_, err = encodedPayload.Decode()
			require.Error(t, err)
		})
	}
}

// TestEncodeWithFewerElements tests that having fewer bytes than expected doesn't throw an error
func TestEncodeWithFewerElements(t *testing.T) {
	testRandom := random.NewTestRandom()
	originalData := testRandom.Bytes(testRandom.Intn(1024) + 33)
	encodedPayload := Payload(originalData).ToEncodedPayload()

	originalBlob, err := encodedPayload.ToBlob(codecs.PolynomialFormCoeff)
	require.NoError(t, err)

	truncatedCoefficients := make([]fr.Element, originalBlob.LenSymbols()-1)
	copy(truncatedCoefficients, originalBlob.coeffPolynomial)
	truncatedBlob := BlobFromCoefficients(truncatedCoefficients)

	reconstructedEncodedPayload := truncatedBlob.ToEncodedPayload(codecs.PolynomialFormCoeff)
	// even though the actual length will be less than the claimed length, we shouldn't see any error
	require.Equal(t, encodedPayload, reconstructedEncodedPayload)
}
