// nolint: lll // long lines are expected b/c of examples
package coretypes

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestEncodePayload tests that the encoding of a Payload to an EncodedPayload works as expected.
func TestEncodeDecodePayload(t *testing.T) {

	// map of hex-encoded payloads (inputs) and their expected EncodedPayloads (outputs).
	// The encoded payloads are broken into 32 byte chunks so as to make them more easily understandable.
	// For example, the first string is always the header.
	testCases := []struct {
		name                      string
		payloadHex                string
		expectedEncodedPayloadHex string
	}{
		{
			name:       "Empty Payload -> header-only (single FE) encodedPayload",
			payloadHex: "",
			// Empty payload encodes to an all zero header (because version=0 and payloadlen=0)
			expectedEncodedPayloadHex: "0000000000000000000000000000000000000000000000000000000000000000",
		},
		// The 3 below cases are all very similar; their payload doesn't matter, we just
		// check that they are contained in the EncodedPayload FE and that the header has the correct length.
		{
			name:       "1 Byte Payload -> 2 FE EncodedPayload",
			payloadHex: "01",
			expectedEncodedPayloadHex: "0000000000010000000000000000000000000000000000000000000000000000" + // header with len 1 payload
				"0001000000000000000000000000000000000000000000000000000000000000", // first byte is always 0 due to bn254 encoding
		},
		{
			name:       "2 Byte Payload -> 2 FE EncodedPayload",
			payloadHex: "0102",
			expectedEncodedPayloadHex: "0000000000020000000000000000000000000000000000000000000000000000" +
				"0001020000000000000000000000000000000000000000000000000000000000",
		},
		{
			name:       "31 Byte Payload -> 2 FE EncodedPayload",
			payloadHex: "01020304050607080910111213141516171819202122232425262728293031",
			expectedEncodedPayloadHex: "00000000001f0000000000000000000000000000000000000000000000000000" +
				"0001020304050607080910111213141516171819202122232425262728293031",
		},
		{
			// Each 31 bytes of payload get encoded into a single FE, so we need 2 FEs to contain the payload,
			// which with the header leads to 3 FEs. Since EncodedPayload have to have a power of 2 number of FEs,
			// the result is a 4 FE encodedPayload.
			name:       "32 Byte Payload -> 4 FE EncodedPayload (EncodedPayload is always power of 2 FE)",
			payloadHex: "0102030405060708091011121314151617181920212223242526272829303132",
			expectedEncodedPayloadHex: "0000000000200000000000000000000000000000000000000000000000000000" +
				"0001020304050607080910111213141516171819202122232425262728293031" +
				"0032000000000000000000000000000000000000000000000000000000000000" +
				"0000000000000000000000000000000000000000000000000000000000000000",
		},
	}
	for _, tc := range testCases {
		t.Run("EncodePayload "+tc.payloadHex, func(t *testing.T) {
			payload, err := hex.DecodeString(tc.payloadHex)
			require.NoError(t, err)
			encodedPayload := Payload(payload).ToEncodedPayload()
			// Run this here even though its called in Decode() in order to catch encoding bugs early.
			require.NoError(t, encodedPayload.checkLenInvariant())
			require.Equal(t, tc.expectedEncodedPayloadHex, hex.EncodeToString(encodedPayload.bytes))
			decodedPayload, err := encodedPayload.Decode()
			require.NoError(t, err)
			require.Equal(t, Payload(payload), decodedPayload)
		})
	}
}

func TestDecodePayloadErrors(t *testing.T) {
	// The encodedHex payloads are broken into 32 byte chunks so as to make them more easily understandable.
	// For example, the first string is always the header.
	testCases := []struct {
		name              string
		encodedPayloadHex string
	}{
		{
			name:              "Insufficient Length Doesn't Contain Header",
			encodedPayloadHex: "000000000000",
		},
		{
			name:              "First byte must be 0x00",
			encodedPayloadHex: "0100000000000000000000000000000000000000000000000000000000000000",
		},
		{
			name:              "Only version 0x00 is supported",
			encodedPayloadHex: "0001000000000000000000000000000000000000000000000000000000000000",
		},
		{
			name:              "Payload length must be a multiple of 32 bytes",
			encodedPayloadHex: "0000000000010000000000000000000000000000000000000000000000000000" + "000100",
		},
		{
			name: "wrong payload length: 32 bytes of data, but header says 64",
			encodedPayloadHex: "0000000000400000000000000000000000000000000000000000000000000000" +
				"0000000000000000000000000000000000000000000000000000000000000000",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			bytes, err := hex.DecodeString(tc.encodedPayloadHex)
			require.NoError(t, err)

			encodedPayload := DeserializeEncodedPayloadUnchecked(bytes)
			_, err = encodedPayload.Decode()
			require.Error(t, err)
		})
	}
}
