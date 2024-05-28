package clients

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/fft"
	"github.com/Layr-Labs/eigenda/encoding/rs"
	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

type BlobEncodingVersion byte

const (
	// This minimal blob encoding includes a version byte, a length varuint, and 31 byte field element mapping. It does not include IFFT padding + IFFT.
	DefaultBlobEncoding BlobEncodingVersion = 0x0
	IFFTBlobEncoding    BlobEncodingVersion = 0x01
)

type BlobCodec interface {
	DecodeBlob(encodedData []byte) ([]byte, error)
	EncodeBlob(rawData []byte) []byte
}

func BlobEncodingVersionToCodec(version BlobEncodingVersion) (BlobCodec, error) {
	switch version {
	case DefaultBlobEncoding:
		return DefaultBlobEncodingCodec{}, nil
	case IFFTBlobEncoding:
		return IFFTBlobEncodingCodec{}, nil
	default:
		return nil, fmt.Errorf("unsupported blob encoding version: %x", version)
	}
}

type DefaultBlobEncodingCodec struct{}
type IFFTBlobEncodingCodec struct{}

var _ BlobCodec = DefaultBlobEncodingCodec{}
var _ IFFTBlobCodec = IFFTBlobEncodingCodec{}

func (v DefaultBlobEncodingCodec) EncodeBlob(rawData []byte) []byte {
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

func (v DefaultBlobEncodingCodec) DecodeBlob(encodedData []byte) ([]byte, error) {
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
		return nil, fmt.Errorf("failed to copy unpadded data into final buffer, length: %d, bytes read: %d", length, n)
	}
	if uint64(n) != length {
		return nil, fmt.Errorf("data length does not match length prefix")
	}

	return rawData, nil
}

func (v DefaultBlobEncodingCodec) EncodeBlobIFFT(rawData []byte) ([]byte, error) {
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

	// expand to next power of 2
	encodedDataFr, err := rs.ToFrArray(encodedData)
	if err != nil {
		return nil, fmt.Errorf("failed to convert encoded data to fr array")
	}

	encodedDataFrLen := uint64(len(encodedDataFr))
	encodedDataFrLenPow2 := encoding.NextPowerOf2(encodedDataFrLen)
	paddedEncodedDataFr := make([]fr.Element, encodedDataFrLenPow2)

	for i := 0; i < len(paddedEncodedDataFr); i++ {
		if i < len(encodedDataFr) {
			paddedEncodedDataFr[i].Set(&encodedDataFr[i])
		} else {
			paddedEncodedDataFr[i].SetZero()
		}
	}

	// perform IFFT
	fs := fft.NewFFTSettings(uint8(encodedDataFrLen))
	encodedDataIfftFr, err := fs.FFT(paddedEncodedDataFr, true)
	if err != nil {
		return nil, fmt.Errorf("failed to perform IFFT: %w", err)
	}

	encodedDataIfft := rs.ToByteArray(encodedDataIfftFr, encodedDataFrLenPow2*encoding.BYTES_PER_SYMBOL)

	return encodedDataIfft, nil
}

func (v DefaultBlobEncodingCodec) DecodeBlobIFFT(encodedData []byte) ([]byte, error) {
	// perform FFT on data
	decodedDataFr, err := rs.ToFrArray(encodedData)
	if err != nil {
		return nil, fmt.Errorf("failed to convert encoded data to fr array")
	}

	decodedDataFrLen := uint64(len(decodedDataFr))
	decodedDataFrLenPow2 := encoding.NextPowerOf2(decodedDataFrLen)
}
