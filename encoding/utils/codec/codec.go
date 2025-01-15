package codec

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/Layr-Labs/eigenda/encoding"
)

type BlobEncodingVersion byte

const (
	// BlobEncodingVersion0 entails a 32 byte header = [0x00, version byte, big-endian uint32 len of payload, 0x00, 0x00,...]
	// followed by the encoded data [0x00, 31 bytes of data, 0x00, 31 bytes of data,...]
	BlobEncodingVersion0 BlobEncodingVersion = 0x0
)

// EncodePayload accepts an arbitrary payload byte array, and encodes it.
//
// The returned bytes shall be interpreted as a polynomial in Eval form, where each contained field element of
// length 32 represents the evaluation of the polynomial at an expanded root of unity
//
// The returned bytes may or may not represent a blob. If the system is configured to distribute blobs in Coeff form,
// then the data returned from this function must be IFFTed to produce the final blob. If the system is configured to
// distribute blobs in Eval form, then the data returned from this function is the final blob representation.
//
// Example encoding:
//                  Payload header (32 bytes total)                                  Encoded Payload Data
// [0x00, version byte, big-endian uint32 len of payload, 0x00, ...] + [0x00, 31 bytes of data, 0x00, 31 bytes of data,...]
func EncodePayload(payload []byte) []byte {
	payloadHeader := make([]byte, 32)
	// first byte is always 0 to ensure the payloadHeader is a valid bn254 element
	payloadHeader[1] = byte(BlobEncodingVersion0) // encode version byte

	// encode payload length as uint32
	binary.BigEndian.PutUint32(
		payloadHeader[2:6],
		uint32(len(payload))) // uint32 should be more than enough to store the length (approx 4gb)

	// encode payload modulo bn254
	// the resulting bytes subsequently may be treated as the evaluation of a polynomial
	polynomialEval := ConvertByPaddingEmptyByte(payload)

	encodedPayload := append(payloadHeader, polynomialEval...)

	return encodedPayload
}

// DecodePayload accepts bytes representing an encoded payload, and returns the decoded payload
//
// This function expects the parameter bytes to be a polynomial in Eval form. In other words, if blobs in the system
// are being distributed in Coeff form, a blob must be FFTed prior to being passed into the function.
func DecodePayload(encodedPayload []byte) ([]byte, error) {
	if len(encodedPayload) < 32 {
		return nil, fmt.Errorf("encoded payload does not contain 32 header bytes, meaning it is malformed")
	}

	payloadLength := binary.BigEndian.Uint32(encodedPayload[2:6])

	// decode raw data modulo bn254
	nonPaddedData := RemoveEmptyByteFromPaddedBytes(encodedPayload[32:])

	reader := bytes.NewReader(nonPaddedData)
	payload := make([]byte, payloadLength)
	readLength, err := reader.Read(payload)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to copy unpadded data into final buffer, length: %d, bytes read: %d",
			payloadLength, readLength)
	}
	if uint32(readLength) != payloadLength {
		return nil, fmt.Errorf("data length does not match length prefix")
	}

	return payload, nil
}

// ConvertByPaddingEmptyByte takes bytes and insert an empty byte at the front of every 31 byte.
// The empty byte is padded at the low address, because we use big endian to interpret a field element.
// This ensures every 32 bytes is within the valid range of a field element for bn254 curve.
// If the input data is not a multiple of 31, the remainder is added to the output by
// inserting a 0 and the remainder. The output is thus not necessarily a multiple of 32.
func ConvertByPaddingEmptyByte(data []byte) []byte {
	dataSize := len(data)
	parseSize := encoding.BYTES_PER_SYMBOL - 1
	putSize := encoding.BYTES_PER_SYMBOL

	dataLen := (dataSize + parseSize - 1) / parseSize

	validData := make([]byte, dataLen*putSize)
	validEnd := len(validData)

	for i := 0; i < dataLen; i++ {
		start := i * parseSize
		end := (i + 1) * parseSize
		if end > len(data) {
			end = len(data)
			// 1 is the empty byte
			validEnd = end - start + 1 + i*putSize
		}

		// with big endian, set first byte is always 0 to ensure data is within valid range of
		validData[i*encoding.BYTES_PER_SYMBOL] = 0x00
		copy(validData[i*encoding.BYTES_PER_SYMBOL+1:(i+1)*encoding.BYTES_PER_SYMBOL], data[start:end])

	}
	return validData[:validEnd]
}

// RemoveEmptyByteFromPaddedBytes takes bytes and remove the first byte from every 32 bytes.
// This reverses the change made by the function ConvertByPaddingEmptyByte.
// The function does not assume the input is a multiple of BYTES_PER_SYMBOL(32 bytes).
// For the reminder of the input, the first byte is taken out, and the rest is appended to
// the output.
func RemoveEmptyByteFromPaddedBytes(data []byte) []byte {
	dataSize := len(data)
	parseSize := encoding.BYTES_PER_SYMBOL
	dataLen := (dataSize + parseSize - 1) / parseSize

	putSize := encoding.BYTES_PER_SYMBOL - 1

	validData := make([]byte, dataLen*putSize)
	validLen := len(validData)

	for i := 0; i < dataLen; i++ {
		// add 1 to leave the first empty byte untouched
		start := i*parseSize + 1
		end := (i + 1) * parseSize

		if end > len(data) {
			end = len(data)
			validLen = end - start + i*putSize
		}

		copy(validData[i*putSize:(i+1)*putSize], data[start:end])
	}
	return validData[:validLen]
}
