package codecs

import (
	"encoding/binary"
	"fmt"

	"github.com/Layr-Labs/eigenda/encoding/rs"
	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

// encodedPayload represents a payload that has had an encoding applied to it
//
// Example encoding:
//
//              Encoded Payload header (32 bytes total)                   Encoded Payload Data (len is multiple of 32)
// [0x00, version byte, big-endian uint32 len of payload, 0x00, ...] + [0x00, 31 bytes of data, 0x00, 31 bytes of data,...]
type encodedPayload struct {
	// the size of these bytes is guaranteed to be a multiple of 32
	bytes []byte
}

// newEncodedPayload accepts a payload, and performs the PayloadEncodingVersion0 encoding to create an encoded payload
func newEncodedPayload(payload *Payload) (*encodedPayload, error) {
	encodedPayloadHeader := make([]byte, 32)
	// first byte is always 0 to ensure the payloadHeader is a valid bn254 element
	encodedPayloadHeader[1] = byte(PayloadEncodingVersion0) // encode version byte

	payloadBytes := payload.GetBytes()

	// encode payload length as uint32
	binary.BigEndian.PutUint32(
		encodedPayloadHeader[2:6],
		uint32(len(payloadBytes))) // uint32 should be more than enough to store the length (approx 4gb)

	// encode payload modulo bn254, and align to 32 bytes
	encodedData := codec.PadPayload(payloadBytes)
	encodedPayloadBytes := append(encodedPayloadHeader, encodedData...)

	return &encodedPayload{encodedPayloadBytes}, nil
}

// decode applies the inverse of PayloadEncodingVersion0 to an encodedPayload, and returns the decoded Payload
func (ep *encodedPayload) decode() (*Payload, error) {
	claimedLength := binary.BigEndian.Uint32(ep.bytes[2:6])

	// decode raw data modulo bn254
	unpaddedData, err := codec.RemoveInternalPadding(ep.bytes[32:])
	if err != nil {
		return nil, fmt.Errorf("remove internal padding: %w", err)
	}

	if uint32(len(unpaddedData)) < claimedLength {
		return nil, fmt.Errorf(
			"length of unpadded data %d is less than length claimed in encoded payload header %d",
			len(unpaddedData), claimedLength)
	}

	return NewPayload(unpaddedData[0:claimedLength]), nil
}

// toEvalPoly converts the encoded payload to a polynomial in evaluation form
func (ep *encodedPayload) toEvalPoly() (*evalPoly, error) {
	fieldElements, err := rs.BytesToFieldElements(ep.bytes)
	if err != nil {
		return nil, fmt.Errorf("deserialize field elements: %w", err)
	}

	return evalPolyFromElements(fieldElements), nil
}

// toCoeffPoly converts the encoded payload to a polynomial in coefficient form
func (ep *encodedPayload) toCoeffPoly() (*coeffPoly, error) {
	fieldElements, err := rs.BytesToFieldElements(ep.bytes)
	if err != nil {
		return nil, fmt.Errorf("deserialize field elements: %w", err)
	}

	return coeffPolyFromElements(fieldElements), nil
}

// encodedPayloadFromElements accepts an array of field elements, and converts them into an encoded payload
func encodedPayloadFromElements(fieldElements []fr.Element) (*encodedPayload, error) {
	polynomialBytes := rs.FieldElementsToBytes(fieldElements)

	// this is the payload length, as claimed by the encoded payload header
	payloadLength := binary.BigEndian.Uint32(polynomialBytes[2:6])
	// add 32 to the padded data length, since the encoded payload includes an encoded payload header
	encodedPayloadLength := codec.GetPaddedDataLength(payloadLength) + 32

	// no matter what, this will be a multiple of 32, since both encodedPayloadLength and polynomialBytes are a multiple of 32
	// we can't just copy to length of encodedPayloadLength, since it's possible that the polynomial truncated 0s that
	// are counted in the length of the encoded data payload.
	lengthToCopy := min(encodedPayloadLength, uint32(len(polynomialBytes)))

	// TODO: we need to check the claimed payload length before creating this, to make sure it doesn't exceed the max blob size?
	//  Otherwise, I think there is an attack vector to maliciously say a payload is super huge, and OOM clients
	encodedPayloadBytes := make([]byte, encodedPayloadLength)
	copy(encodedPayloadBytes, polynomialBytes[:lengthToCopy])

	return &encodedPayload{encodedPayloadBytes}, nil
}
