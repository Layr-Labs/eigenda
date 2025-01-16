package codecs

import (
	"encoding/binary"
	"fmt"

	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
)

// encodedPayload represents a payload that has had an encoding applied to it
type encodedPayload struct {
	bytes []byte
}

// newEncodedPayload accepts an array of bytes which represent an encodedPayload. It performs the checks necessary
// to guarantee that the bytes are well-formed, and returns a newly constructed object if they are.
//
// Note that this function does not decode the input bytes to perform additional checks, so it is possible to construct
// an encodedPayload, where an attempt to decode will fail.
func newEncodedPayload(encodedPayloadBytes []byte) (*encodedPayload, error) {
	inputLen := len(encodedPayloadBytes)
	if inputLen < 32 {
		return nil, fmt.Errorf(
			"input bytes have length %d, which is smaller than the required 32 header bytes", inputLen)
	}

	return &encodedPayload{
		bytes: encodedPayloadBytes,
	}, nil
}

// decode applies the inverse of DefaultBlobEncoding to an encodedPayload, and returns the decoded Payload
func (ep *encodedPayload) decode() (*Payload, error) {
	claimedLength := binary.BigEndian.Uint32(ep.bytes[2:6])

	// decode raw data modulo bn254
	nonPaddedData, err := codec.RemoveInternalPadding(ep.bytes[32:])
	if err != nil {
		return nil, fmt.Errorf("remove internal padding: %w", err)
	}

	if uint32(len(nonPaddedData)) < claimedLength {
		return nil, fmt.Errorf(
			"data length %d is less than length claimed in payload header %d",
			len(nonPaddedData), claimedLength)
	}

	return NewPayload(nonPaddedData[0:claimedLength]), nil
}

// toEvalPoly converts an encodedPayload into an evalPoly
func (ep *encodedPayload) toEvalPoly() (*evalPoly, error) {
	evalPoly, err := evalPolyFromBytes(ep.bytes)
	if err != nil {
		return nil, fmt.Errorf("new eval poly: %w", err)
	}

	return evalPoly, nil
}

// getBytes returns the raw bytes that underlie the encodedPayload
func (ep *encodedPayload) getBytes() []byte {
	return ep.bytes
}
