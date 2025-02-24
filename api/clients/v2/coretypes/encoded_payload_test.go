package coretypes

import (
	"testing"

	"github.com/Layr-Labs/eigenda/common/testutils/random"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/stretchr/testify/require"
)

// TestDecodeShortBytes checks that an encoded payload with a length less than claimed length fails at decode time
func TestDecodeShortBytes(t *testing.T) {
	testRandom := random.NewTestRandom()
	originalData := testRandom.Bytes(testRandom.Intn(1024) + 33)

	encodedPayload, err := newEncodedPayload(NewPayload(originalData))
	require.NoError(t, err)

	// truncate
	encodedPayload.bytes = encodedPayload.bytes[:len(encodedPayload.bytes) -32]

	payload, err := encodedPayload.decode()
	require.Error(t, err)
	require.Nil(t, payload)
}

// TestDecodeLongBytes checks that an encoded payload with length too much greater than claimed fails at decode
func TestDecodeLongBytes(t *testing.T) {
	testRandom := random.NewTestRandom()
	originalData := testRandom.Bytes(testRandom.Intn(1024) + 1)

	encodedPayload, err := newEncodedPayload(NewPayload(originalData))
	require.NoError(t, err)

	encodedPayload.bytes = append(encodedPayload.bytes, make([]byte, 32)...)
	payload2, err := encodedPayload.decode()
	require.Error(t, err)
	require.Nil(t, payload2)
}

// TestEncodeTooManyElements checks that encodedPayloadFromElements fails at the expect limit, relative to payload
// length and blob length
func TestEncodeTooManyElements(t *testing.T) {
	testRandom := random.NewTestRandom()
	powersOf2 := encoding.GeneratePowersOfTwo(uint32(12))

	for i := 0; i < len(powersOf2); i++ {
		blobLength := powersOf2[i]
		maxPermissiblePayloadLength, err := codec.GetMaxPermissiblePayloadLength(blobLength)
		require.NoError(t, err)

		almostTooLongData := testRandom.Bytes(int(maxPermissiblePayloadLength))
		almostTooLongEncodedPayload, err := newEncodedPayload(NewPayload(almostTooLongData))
		require.NoError(t, err)
		almostTooLongFieldElements, err := almostTooLongEncodedPayload.toFieldElements()
		require.NoError(t, err)
		// there are almost too many field elements for the defined blob length, but not quite
		_, err = encodedPayloadFromElements(almostTooLongFieldElements, maxPermissiblePayloadLength)
		require.NoError(t, err)

		tooLongData := testRandom.Bytes(int(maxPermissiblePayloadLength) + 1)
		tooLongEncodedPayload, err := newEncodedPayload(NewPayload(tooLongData))
		require.NoError(t, err)
		tooLongFieldElements, err := tooLongEncodedPayload.toFieldElements()
		require.NoError(t, err)
		// there is one too many field elements for the defined blob length
		_, err = encodedPayloadFromElements(tooLongFieldElements, maxPermissiblePayloadLength)
		require.Error(t, err)
	}
}

// TestTrailingNonZeros checks that any non-zero values that come after the end of the claimed payload length
// cause an error to be returned.
func TestTrailingNonZeros(t *testing.T) {
	testRandom := random.NewTestRandom()
	originalData := testRandom.Bytes(testRandom.Intn(1024) + 1)

	encodedPayload, err := newEncodedPayload(NewPayload(originalData))
	require.NoError(t, err)

	originalElements, err := encodedPayload.toFieldElements()
	require.NoError(t, err)

	fieldElements1 := make([]fr.Element, len(originalElements))
	copy(fieldElements1, originalElements)

	fieldElements2 := make([]fr.Element, len(originalElements))
	copy(fieldElements2, originalElements)

	// adding a 0 is fine
	fieldElements1 = append(fieldElements1, fr.Element{})
	_, err = encodedPayloadFromElements(fieldElements1, uint32(len(fieldElements1)*encoding.BYTES_PER_SYMBOL))
	require.NoError(t, err)

	// adding a non-0 is non-fine
	fieldElements2 = append(fieldElements2, fr.Element{0,0,0,1})
	_, err = encodedPayloadFromElements(fieldElements2, uint32(len(fieldElements2)*encoding.BYTES_PER_SYMBOL))
	require.Error(t, err)
}

// TestEncodeWithFewerElements tests that having fewer bytes than expected doesn't throw an error
func TestEncodeWithFewerElements(t *testing.T) {
	testRandom := random.NewTestRandom()
	originalData := testRandom.Bytes(testRandom.Intn(1024) + 33)

	encodedPayload, err := newEncodedPayload(NewPayload(originalData))
	require.NoError(t, err)

	originalFieldElements, err := encodedPayload.toFieldElements()
	require.NoError(t, err)

	truncatedFieldElements := make([]fr.Element, len(originalFieldElements)-1)
	// intentionally don't copy all the elements
	copy(truncatedFieldElements, originalFieldElements[:len(originalFieldElements)-1])

	// even though the actual length will be less than the claimed length, we shouldn't see any error
	reconstructedEncodedPayload, err := encodedPayloadFromElements(
		originalFieldElements,
		uint32(len(originalFieldElements))*encoding.BYTES_PER_SYMBOL)
	require.NoError(t, err)
	require.NotNil(t, reconstructedEncodedPayload)
}
