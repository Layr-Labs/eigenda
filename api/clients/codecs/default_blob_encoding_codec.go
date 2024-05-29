package codecs

import (
	"bytes"
	"fmt"

	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
)

type DefaultBlobEncodingCodec struct{}

var _ BlobCodec = DefaultBlobEncodingCodec{}

func (v DefaultBlobEncodingCodec) EncodeBlob(rawData []byte) ([]byte, error) {
	// encode blob encoding version byte
	codecBlobHeader := EncodeCodecBlobHeader(byte(DefaultBlobEncoding), uint32(len(rawData)))

	// encode raw data modulo bn254
	rawDataPadded := codec.ConvertByPaddingEmptyByte(rawData)

	// append raw data
	encodedData := append(codecBlobHeader, rawDataPadded...)

	return encodedData, nil
}

func (v DefaultBlobEncodingCodec) DecodeBlob(encodedData []byte) ([]byte, error) {
	versionByte, length, err := DecodeCodecBlobHeader(encodedData[:32])
	if err != nil {
		return nil, err
	}

	if DefaultBlobEncoding != BlobEncodingVersion(versionByte) {
		return nil, fmt.Errorf("unsupported blob encoding version: %x", versionByte)
	}

	// get non blob header data
	rawDataPadded := encodedData[32:]

	// decode raw data modulo bn254
	decodedRawData := codec.RemoveEmptyByteFromPaddedBytes(rawDataPadded)

	reader := bytes.NewReader(decodedRawData)
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
