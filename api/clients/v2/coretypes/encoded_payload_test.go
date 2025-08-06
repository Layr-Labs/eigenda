package coretypes

import (
	"encoding/binary"
	"testing"

	"github.com/Layr-Labs/eigenda/api/clients/codecs"
	"github.com/Layr-Labs/eigenda/common/testutils/random"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/stretchr/testify/require"
)

// TestDecodeShortBytes checks that an encoded payload with a length less than claimed length fails at decode time
func TestDecodeInvalidPayloadLen(t *testing.T) {
	testRandom := random.NewTestRandom()
	originalData := testRandom.Bytes(testRandom.Intn(1024) + 33)
	encodedPayload := Payload(originalData).ToEncodedPayload()

	// Changed the header payload length to be longer than the actual encodedPayload length.
	// This way the claimed payload clearly doesn't fit in the encoded payload.
	binary.BigEndian.PutUint32(encodedPayload.bytes[2:6], uint32(len(encodedPayload.bytes)+1))

	payload, err := encodedPayload.Decode()
	require.Error(t, err)
	require.Nil(t, payload)
}

// TestDecodeLongBytes checks that an encoded payload with length too much greater than claimed fails at decode
func TestDecodeLongBytes(t *testing.T) {
	testRandom := random.NewTestRandom()
	originalData := testRandom.Bytes(testRandom.Intn(1024) + 1)
	encodedPayload := Payload(originalData).ToEncodedPayload()

	// appending 33 bytes to the encoded payload guarantees that, after removing padding, the unpadded bytes will be
	// at least 32 bytes longer than the expected length, which is the error case we're trying to trigger here
	encodedPayload.bytes = append(encodedPayload.bytes, make([]byte, 33)...)
	payload2, err := encodedPayload.Decode()
	require.Error(t, err)
	require.Nil(t, payload2)
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

	reconstructedEncodedPayload, err := truncatedBlob.ToEncodedPayload(codecs.PolynomialFormCoeff)
	require.NoError(t, err)
	// even though the actual length will be less than the claimed length, we shouldn't see any error
	require.Equal(t, encodedPayload, reconstructedEncodedPayload)
}
