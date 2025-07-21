package coretypes

import (
	"encoding/binary"
	"fmt"

	"github.com/Layr-Labs/eigenda/api/clients/codecs"
	"github.com/Layr-Labs/eigenda/encoding/rs"
	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

// encodedPayload represents a payload that has had an encoding applied to it
//
// Example encoding:
//
//	Encoded Payload header (32 bytes total)                   Encoded Payload Data (len is multiple of 32)
//
// [0x00, version byte, big-endian uint32 len of payload, 0x00, ...] + [0x00, 31 bytes of data, 0x00, 31 bytes of data,...]
type encodedPayload struct {
	// the size of these bytes is guaranteed to be a multiple of 32
	bytes []byte
}

// newEncodedPayload accepts a payload, and performs the PayloadEncodingVersion0 encoding to create an encoded payload
func newEncodedPayload(payload *Payload) (*encodedPayload, error) {
	encodedPayloadHeader := make([]byte, 32)
	// first byte is always 0 to ensure the payloadHeader is a valid bn254 element
	encodedPayloadHeader[1] = byte(codecs.PayloadEncodingVersion0) // encode version byte

	payloadBytes := payload.Serialize()

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
	if len(ep.bytes) < 32 {
		return nil, fmt.Errorf("encoded payload must be at least 32 bytes long, but got %d bytes", len(ep.bytes))
	}
	if ep.bytes[0] != 0x00 {
		return nil, fmt.Errorf("encoded payload header first byte must be 0x00, but got %x", ep.bytes[0])
	}
	if ep.bytes[1] != byte(codecs.PayloadEncodingVersion0) {
		return nil, fmt.Errorf("encoded payload header version byte must be %x, but got %x", codecs.PayloadEncodingVersion0, ep.bytes[1])
	}

	claimedLength := binary.BigEndian.Uint32(ep.bytes[2:6])

	// decode raw data modulo bn254
	unpaddedData, err := codec.RemoveInternalPadding(ep.bytes[32:])
	if err != nil {
		return nil, fmt.Errorf("remove internal padding: %w", err)
	}

	unpaddedDataLength := uint32(len(unpaddedData))

	// data length is checked when constructing an encoded payload. If this error is encountered, that means there
	// must be a flaw in the logic at construction time (or someone was bad and didn't use the proper construction methods)
	if unpaddedDataLength < claimedLength {
		return nil, fmt.Errorf(
			"length of unpadded data %d is less than length claimed in encoded payload header %d. this should never happen",
			unpaddedDataLength, claimedLength)
	}

	// unpadded data length can be slightly bigger than the claimed length, since RemoveInternalPadding doesn't
	// do anything to remove trailing zeros that may have been added when the data was initially padded.
	// however, this extra padding shouldn't exceed 31 bytes, because that's the most that would be added
	// when padding the data length to 32 bytes. If this error occurs, that means there must be a flaw in the logic at
	// construction time (or someone was bad and didn't use the proper construction methods)
	if unpaddedDataLength > claimedLength+31 {
		return nil, fmt.Errorf(
			"length of unpadded data %d is more than 31 bytes longer than claimed length %d. this should never happen",
			unpaddedDataLength, claimedLength)
	}

	return NewPayload(unpaddedData[0:claimedLength]), nil
}

// toFieldElements converts the encoded payload to an array of field elements
func (ep *encodedPayload) toFieldElements() ([]fr.Element, error) {
	fieldElements, err := rs.ToFrArray(ep.bytes)
	if err != nil {
		return nil, fmt.Errorf("deserialize field elements: %w", err)
	}

	return fieldElements, nil
}

// encodedPayloadFromElements accepts an array of field elements, and converts them into an encoded payload
//
// maxPayloadLength is the maximum length in bytes that the contained Payload is permitted to be
func encodedPayloadFromElements(fieldElements []fr.Element, maxPayloadLength uint32) (*encodedPayload, error) {
	polynomialBytes := rs.SerializeFieldElements(fieldElements)
	// this is the payload length in bytes, as claimed by the encoded payload header
	payloadLength := binary.BigEndian.Uint32(polynomialBytes[2:6])

	if payloadLength > maxPayloadLength {
		return nil, fmt.Errorf(
			"payload length claimed in encoded payload header (%d bytes) is larger than the permitted maximum (%d bytes)",
			payloadLength, maxPayloadLength)
	}

	// this is the length you would get if you padded a payload of the length claimed in the encoded payload header
	paddedLength := codec.GetPaddedDataLength(payloadLength)
	// add 32 to the padded data length, since the encoded payload includes an encoded payload header
	encodedPayloadLength := paddedLength + 32

	polynomialByteCount := uint32(len(polynomialBytes))

	// no matter what, this will be a multiple of 32, since the two possible values, encodedPayloadLength and
	// polynomialBytes, are multiples of 32. This is important, since the encoded payload being created is
	// expected to have a byte count that's a multiple of 32.
	lengthToCopy := encodedPayloadLength

	// if encodedPayloadLength is greater than the polynomial bytes, that indicates that the polynomial bytes we have
	// are missing trailing 0 bytes which were originally part of the dispersed blob. For this to happen, it means
	// that whichever source provided us with these bytes truncated the trailing 0s. This probably won't happen in
	// practice, but if it were to happen, it wouldn't be caught when verifying commitments, since trailing 0s don't
	// affect the commitment. This isn't a problem, though: we can handle this edge case here.
	if encodedPayloadLength > polynomialByteCount {
		// we are copying from the polynomialBytes, so make sure that we don't try to copy more data than actually exists
		lengthToCopy = polynomialByteCount
	} else if encodedPayloadLength < polynomialByteCount {
		// we assume that the polynomialBytes might have additional trailing 0s beyond the expected size of the encoded
		// payload. Here, we check the assumption that all trailing values are 0. If there are any non-zero trailing
		// values, something has gone wrong in the data pipeline, and this should produce a loud failure. Either a
		// dispersing client is playing sneaky games, or there's a bug somewhere.
		err := checkTrailingZeros(polynomialBytes, encodedPayloadLength)
		if err != nil {
			return nil, fmt.Errorf("check that trailing values in polynomial are zeros: %w", err)
		}
	}

	encodedPayloadBytes := make([]byte, encodedPayloadLength)
	copy(encodedPayloadBytes, polynomialBytes[:lengthToCopy])

	return &encodedPayload{encodedPayloadBytes}, nil
}

// checkTrailingZeros accepts an array of bytes, and the number of bytes at the front of the array which are permitted
// to be non-zero
//
// This function returns an error if any byte in the array after these permitted non-zero values is found to be non-zero
func checkTrailingZeros(inputBytes []byte, nonZeroLength uint32) error {
	for i := uint32(len(inputBytes)) - 1; i >= nonZeroLength; i-- {
		if inputBytes[i] != 0x0 {
			return fmt.Errorf("byte at index %d was expected to be 0x0, but instead was %x", i, inputBytes[i])
		}
	}

	return nil
}
