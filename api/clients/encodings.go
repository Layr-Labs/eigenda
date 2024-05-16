package clients

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
)

type BlobEncodingVersion byte

const (
	// This minimal blob encoding includes a version byte, a length varuint, and 31 byte field element mapping. It does not include IFFT padding + IFFT.
	DefaultBlobEncoding BlobEncodingVersion = 0x00
)

type BlobCodec interface {
	DecodeBlob(encodedData []byte) ([]byte, error)
	EncodeBlob(rawData []byte) []byte
}

func EncodingVersionToCodec(version BlobEncodingVersion) (BlobCodec, error) {
	switch version {
	case DefaultBlobEncoding:
		return NoIFFTCodec{}, nil
	default:
		return nil, fmt.Errorf("unsupported blob encoding version: %x", version)
	}
}

type NoIFFTCodec struct{}

var _ BlobCodec = NoIFFTCodec{}

func (v NoIFFTCodec) EncodeBlob(rawData []byte) []byte {
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

	return encodedData
}

func (v NoIFFTCodec) DecodeBlob(encodedData []byte) ([]byte, error) {
	// decode modulo bn254
	decodedData := codec.RemoveEmptyByteFromPaddedBytes(encodedData)

	// Return exact data with buffer removed
	reader := bytes.NewReader(decodedData)

	// read version byte, we will not use it for now since there is only one version
	_, err := reader.ReadByte()
	if err != nil {
		return nil, fmt.Errorf("failed to read version byte")
	}

	// read length uvarint
	length, err := binary.ReadUvarint(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to decode length uvarint prefix")
	}

	rawData := make([]byte, length)
	n, err := reader.Read(rawData)
	if err != nil {
		return nil, fmt.Errorf("failed to copy unpadded data into final buffer")
	}
	if uint64(n) != length {
		return nil, fmt.Errorf("data length does not match length prefix")
	}

	return rawData, nil
}
