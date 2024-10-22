package codecs

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
)

type DefaultBlobCodec struct{}

var _ BlobCodec = DefaultBlobCodec{}

func NewDefaultBlobCodec() DefaultBlobCodec {
	return DefaultBlobCodec{}
}

// EncodeBlob can never return an error, but to maintain the interface it is included
// so that it can be swapped for the IFFTCodec without changing the interface
func (v DefaultBlobCodec) EncodeBlob(rawData []byte) ([]byte, error) {
	codecBlobHeader := make([]byte, 32)
	// first byte is always 0 to ensure the codecBlobHeader is a valid bn254 element
	// encode version byte
	codecBlobHeader[1] = byte(DefaultBlobEncoding)

	// encode length as uint32
	binary.BigEndian.PutUint32(codecBlobHeader[2:6], uint32(len(rawData))) // uint32 should be more than enough to store the length (approx 4gb)

	// encode raw data modulo bn254
	rawDataPadded := codec.ConvertByPaddingEmptyByte(rawData)

	// append raw data
	encodedData := append(codecBlobHeader, rawDataPadded...)

	return encodedData, nil
}

func (v DefaultBlobCodec) DecodeBlob(data []byte) ([]byte, error) {
	if len(data) < 32 {
		return nil, fmt.Errorf("blob does not contain 32 header bytes, meaning it is malformed")
	}

	length := binary.BigEndian.Uint32(data[2:6])

	// decode raw data modulo bn254
	decodedData := codec.RemoveEmptyByteFromPaddedBytes(data[32:])

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
