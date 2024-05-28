package codecs

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
)

type DefaultBlobEncodingCodec struct{}

var _ BlobCodec = DefaultBlobEncodingCodec{}

func (v DefaultBlobEncodingCodec) EncodeBlob(rawData []byte) ([]byte, error) {
	// encode current blob encoding version byte
	encodedData := make([]byte, 0, 1+8+len(rawData))

	// append version byte
	encodedData = append(encodedData, byte(DefaultBlobEncoding))

	// encode data length
	encodedData = append(encodedData, ConvertIntToVarUInt(len(rawData))...)

	// append raw data
	encodedData = append(encodedData, rawData...)

	// encode modulo bn254
	encodedData = codec.ConvertByPaddingEmptyByte(encodedData)

	return encodedData, nil
}

func (v DefaultBlobEncodingCodec) DecodeBlob(encodedData []byte) ([]byte, error) {
	// decode modulo bn254
	decodedData := codec.RemoveEmptyByteFromPaddedBytes(encodedData)

	// Return exact data with buffer removed
	reader := bytes.NewReader(decodedData)

	versionByte, err := reader.ReadByte()
	if err != nil {
		return nil, fmt.Errorf("failed to read version byte")
	}
	if DefaultBlobEncoding != BlobEncodingVersion(versionByte) {
		return nil, fmt.Errorf("unsupported blob encoding version: %x", versionByte)
	}

	// read length uvarint
	length, err := binary.ReadUvarint(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to decode length uvarint prefix")
	}

	rawData := make([]byte, length)
	n, err := reader.Read(rawData)
	if err != nil {
		return nil, fmt.Errorf("failed to copy unpadded data into final buffer, length: %d, bytes read: %d", length, n)
	}
	if uint64(n) != length {
		return nil, fmt.Errorf("data length does not match length prefix")
	}

	return rawData, nil
}
