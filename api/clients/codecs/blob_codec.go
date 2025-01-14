package codecs

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
)

type BlobEncodingVersion byte

const (
	// DefaultBlobEncoding entails a 32 byte header = [0x00, version byte, uint32 len of data, 0x00, 0x00,...]
	// followed by the encoded data [0x00, 31 bytes of data, 0x00, 31 bytes of data,...]
	DefaultBlobEncoding BlobEncodingVersion = 0x0
)

// EncodePayload accepts an arbitrary payload byte array, and encodes it.
//
// The returned bytes may be interpreted as a polynomial in Eval form, where each contained field element of
// length 32 represents the evaluation of the polynomial at an expanded root of unity
//
// The returned bytes may or may not represent a blob. If the system is configured to distribute blobs in Coeff form,
// then the data returned from this function must be IFFTed to produce the final blob. If the system is configured to
// distribute blobs in Eval form, then the data returned from this function is the final blob representation.
func EncodePayload(payload []byte) []byte {
	payloadHeader := make([]byte, 32)
	// first byte is always 0 to ensure the payloadHeader is a valid bn254 element
	// encode version byte
	payloadHeader[1] = byte(DefaultBlobEncoding)

	// encode payload length as uint32
	binary.BigEndian.PutUint32(
		payloadHeader[2:6],
		uint32(len(payload))) // uint32 should be more than enough to store the length (approx 4gb)

	// encode payload modulo bn254
	// the resulting bytes are subsequently treated as the evaluation of a polynomial
	polynomialEval := codec.ConvertByPaddingEmptyByte(payload)

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
	decodedData := codec.RemoveEmptyByteFromPaddedBytes(encodedPayload[32:])

	// get non blob header data
	reader := bytes.NewReader(decodedData)
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
