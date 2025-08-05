package coretypes

import (
	"github.com/Layr-Labs/eigenda/encoding"
	"testing"

	"github.com/Layr-Labs/eigenda/common/testutils/random"
	"github.com/stretchr/testify/require"
)

func TestPayloadToEncodedPayloadPowerOf2(t *testing.T) {
	testRandom := random.NewTestRandom()
	originalData := testRandom.Bytes(33)
	encodedPayload := Payload(originalData).ToEncodedPayload()

	// Check that the length of the encoded payload is a power of 2
	require.True(t, encoding.IsPowerOfTwo(len(encodedPayload.bytes)), "Encoded payload length should be a power of 2")
}
