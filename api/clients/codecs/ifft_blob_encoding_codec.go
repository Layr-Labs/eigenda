package codecs

import (
	"bytes"
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

// func (v IFFTBlobEncodingCodec) DecodeBlob(encodedData []byte) ([]byte, error) {
// 	// check that the length of the data is a power of two
// 	nextPowerOfTwo := encoding.NextPowerOf2(uint64(len(encodedData)))

// 	if len(encodedData) != int(nextPowerOfTwo) {
// 		return nil, fmt.Errorf("encoded data length is not a power of two, data length: %d", len(encodedData))
// 	}

// 	// perform FFT on data
// 	decodedDataFr, err := rs.ToFrArray(encodedData)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to convert encoded data to fr array")
// 	}

// 	decodedDataFrLen := len(decodedDataFr)

// 	maxScale := uint8(math.Log2(float64(decodedDataFrLen)))

// 	fs := fft.NewFFTSettings(maxScale)
// 	decodedDataFftFr, err := fs.FFT(decodedDataFr, false)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to perform FFT: %w", err)
// 	}

// 	decodedDataBytes := rs.ToByteArray(decodedDataFftFr, uint64(decodedDataFrLen)*encoding.BYTES_PER_SYMBOL)

// 	// decode modulo bn254
// 	decodedData := codec.RemoveEmptyByteFromPaddedBytes(decodedDataBytes)

// 	// Return exact data with buffer removed
// 	reader := bytes.NewReader(decodedData)

// 	versionByte, err := reader.ReadByte()
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to read version byte")
// 	}

// 	if IFFTBlobEncoding != BlobEncodingVersion(versionByte) {
// 		return nil, fmt.Errorf("unsupported blob encoding version: %x", versionByte)
// 	}

// 	// read length uvarint
// 	length, err := binary.ReadUvarint(reader)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to decode length uvarint prefix")
// 	}

// 	rawData := make([]byte, length)
// 	n, err := reader.Read(rawData)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to copy unpadded data into final buffer, length: %d, bytes read: %d", length, n)
// 	}
// 	if uint64(n) != length {
// 		return nil, fmt.Errorf("data length does not match length prefix")
// 	}

// 	return rawData, nil
// }

func (v IFFTBlobEncodingCodec) EncodeBlob(rawData []byte) ([]byte, error) {
	// create the 32 bytes long codec blob header
	codecBlobHeader := EncodeCodecBlobHeader(byte(IFFTBlobEncoding), uint64(len(rawData)))

	// encode modulo bn254
	encodedRawData := codec.ConvertByPaddingEmptyByte(rawData)

	rawDataFr, err := rs.ToFrArray(encodedRawData)
	if err != nil {
		return nil, fmt.Errorf("failed to convert raw data to fr array: %w", err)
	}

	rawDataFrLen := len(rawDataFr)
	rawDataFrLenPow2 := encoding.NextPowerOf2(uint64(rawDataFrLen))
	paddedRawDataFr := make([]fr.Element, rawDataFrLenPow2)

	for i := 0; i < len(paddedRawDataFr); i++ {
		if i < len(rawDataFr) {
			paddedRawDataFr[i].Set(&rawDataFr[i])
		} else {
			paddedRawDataFr[i].SetZero()
		}
	}

	maxScale := uint8(math.Log2(float64(rawDataFrLenPow2)))

	// perform IFFT
	fs := fft.NewFFTSettings(maxScale)
	rawDataIFFTFr, err := fs.FFT(paddedRawDataFr, true)
	if err != nil {
		return nil, fmt.Errorf("failed to perform IFFT: %w", err)
	}

	rawDataIFFTBytes := rs.ToByteArray(rawDataIFFTFr, rawDataFrLenPow2*encoding.BYTES_PER_SYMBOL)

	// append raw data
	encodedData := append(codecBlobHeader, rawDataIFFTBytes...)

	return encodedData, nil
}

func (v IFFTBlobEncodingCodec) DecodeBlob(encodedData []byte) ([]byte, error) {

	versionByte, length, err := DecodeCodecBlobHeader(encodedData[:32])
	if err != nil {
		return nil, fmt.Errorf("failed to decode codec blob header: %w", err)
	}

	if IFFTBlobEncoding != BlobEncodingVersion(versionByte) {
		return nil, fmt.Errorf("unsupported blob encoding version: %x", versionByte)
	}

	paddedIFFTRawData := encodedData[32:]

	nextPowerOfTwo := encoding.NextPowerOf2(uint64(len(paddedIFFTRawData)))

	if len(paddedIFFTRawData) != int(nextPowerOfTwo) {
		return nil, fmt.Errorf("encoded data length is not a power of two, data length: %d", len(encodedData))
	}

	// decode modulo bn254
	rawDataIFFT := codec.RemoveEmptyByteFromPaddedBytes(paddedIFFTRawData)

	rawDataIFFTFr, err := rs.ToFrArray(rawDataIFFT)
	if err != nil {
		return nil, fmt.Errorf("failed to convert raw data IFFT to fr array")
	}

	maxScale := uint8(math.Log2(float64(len(rawDataIFFTFr))))

	fs := fft.NewFFTSettings(maxScale)
	rawDataFr, err := fs.FFT(rawDataIFFTFr, false)
	if err != nil {
		return nil, fmt.Errorf("failed to perform FFT: %w", err)
	}

	rawDataPadded := rs.ToByteArray(rawDataFr, uint64(len(rawDataFr))*encoding.BYTES_PER_SYMBOL)

	reader := bytes.NewReader(rawDataPadded)
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
