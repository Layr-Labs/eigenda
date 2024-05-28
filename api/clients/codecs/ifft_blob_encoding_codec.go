package codecs

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"

	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/fft"
	"github.com/Layr-Labs/eigenda/encoding/rs"
	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

type IFFTBlobEncodingCodec struct{}

var _ BlobCodec = IFFTBlobEncodingCodec{}

func (v IFFTBlobEncodingCodec) EncodeBlob(rawData []byte) ([]byte, error) {
	// encode current blob encoding version byte
	encodedData := make([]byte, 0, 1+8+len(rawData))

	// append version byte
	encodedData = append(encodedData, byte(IFFTBlobEncoding))

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

	maxScale := uint8(math.Log2(float64(encodedDataFrLenPow2)))

	// perform IFFT
	fs := fft.NewFFTSettings(maxScale)
	encodedDataIfftFr, err := fs.FFT(paddedEncodedDataFr, true)
	if err != nil {
		return nil, fmt.Errorf("failed to perform IFFT: %w", err)
	}

	encodedDataIfft := rs.ToByteArray(encodedDataIfftFr, encodedDataFrLenPow2*encoding.BYTES_PER_SYMBOL)

	return encodedDataIfft, nil
}

func (v IFFTBlobEncodingCodec) DecodeBlob(encodedData []byte) ([]byte, error) {
	// check that the length of the data is a power of two
	nextPowerOfTwo := encoding.NextPowerOf2(uint64(len(encodedData)))

	if len(encodedData) != int(nextPowerOfTwo) {
		return nil, fmt.Errorf("encoded data length is not a power of two, data length: %d", len(encodedData))
	}

	// perform FFT on data
	decodedDataFr, err := rs.ToFrArray(encodedData)
	if err != nil {
		return nil, fmt.Errorf("failed to convert encoded data to fr array")
	}

	decodedDataFrLen := len(decodedDataFr)

	maxScale := uint8(math.Log2(float64(decodedDataFrLen)))

	fs := fft.NewFFTSettings(maxScale)
	decodedDataFftFr, err := fs.FFT(decodedDataFr, false)
	if err != nil {
		return nil, fmt.Errorf("failed to perform FFT: %w", err)
	}

	decodedDataBytes := rs.ToByteArray(decodedDataFftFr, uint64(decodedDataFrLen)*encoding.BYTES_PER_SYMBOL)

	// decode modulo bn254
	decodedData := codec.RemoveEmptyByteFromPaddedBytes(decodedDataBytes)

	// Return exact data with buffer removed
	reader := bytes.NewReader(decodedData)

	versionByte, err := reader.ReadByte()
	if err != nil {
		return nil, fmt.Errorf("failed to read version byte")
	}

	if IFFTBlobEncoding != BlobEncodingVersion(versionByte) {
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
