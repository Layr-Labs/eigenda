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
	// The final result of this is then padded to a multiple of 32
	DefaultBlobEncoding BlobEncodingVersion = 0x0
)

// EncodePayload accepts an arbitrary payload byte array, and encodes it.
//
// The returned bytes may be interpreted as a polynomial in evaluation form, where each contained field element of
// length 32 represents the evaluation of the polynomial at an expanded root of unity
//
// The returned bytes may or may not represent a blob. If the system is configured to distribute blobs in Coeff form,
// then the data returned from this function must be IFFTed to produce the final blob. If the system is configured to
// distribute blobs in Eval form, then the data returned from this function is the final blob representation.
func EncodePayload(payload []byte) []byte {
	codecBlobHeader := make([]byte, 32)
	// first byte is always 0 to ensure the codecBlobHeader is a valid bn254 element
	// encode version byte
	codecBlobHeader[1] = byte(DefaultBlobEncoding)

	// encode length as uint32
	binary.BigEndian.PutUint32(
		codecBlobHeader[2:6],
		uint32(len(payload))) // uint32 should be more than enough to store the length (approx 4gb)

	// encode payload modulo bn254
	rawDataPadded := codec.ConvertByPaddingEmptyByte(payload)

	// append raw data
	encodedData := append(codecBlobHeader, rawDataPadded...)

	return encodedData
}

// DecodePayload accepts bytes representing an encoded payload, and returns the decoded payload
//
// This function expects the parameter bytes to be a polynomial in Eval form. In other words, if blobs in the system
// are being distributed in Coeff form, a blob must be FFTed prior to being passed into the function.
func DecodePayload(encodedPayload []byte) ([]byte, error) {
	if len(encodedPayload) < 32 {
		return nil, fmt.Errorf("blob does not contain 32 header bytes, meaning it is malformed")
	}

	length := binary.BigEndian.Uint32(encodedPayload[2:6])

	// decode raw data modulo bn254
	decodedData := codec.RemoveEmptyByteFromPaddedBytes(encodedPayload[32:])

	// get non blob header data
	reader := bytes.NewReader(decodedData)
	rawData := make([]byte, length)
	n, err := reader.Read(rawData)
	if err != nil {
		return nil, fmt.Errorf("failed to copy unpadded data into final buffer, length: %d, bytes read: %d", length, n)
	}
	if uint32(n) != length {
		return nil, fmt.Errorf("data length does not match length prefix")
	}

	return rawData, nil
}
