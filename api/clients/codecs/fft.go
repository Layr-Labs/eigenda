package codecs

import (
	"fmt"
	"math"

	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/fft"
	"github.com/Layr-Labs/eigenda/encoding/rs"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

func FFT(data []byte) ([]byte, error) {
	dataFr, err := rs.ToFrArray(data)
	if err != nil {
		return nil, fmt.Errorf("error converting data to fr.Element: %w", err)
	}
	dataFrLen := uint64(len(dataFr))
	dataFrLenPow2 := encoding.NextPowerOf2(dataFrLen)

	if dataFrLenPow2 != dataFrLen {
		return nil, fmt.Errorf("data length %d is not a power of 2", dataFrLen)
	}

	maxScale := uint8(math.Log2(float64(dataFrLenPow2)))

	fs := fft.NewFFTSettings(maxScale)

	dataFFTFr, err := fs.FFT(dataFr, false)
	if err != nil {
		return nil, fmt.Errorf("failed to perform FFT: %w", err)
	}

	return rs.ToByteArray(dataFFTFr, dataFrLenPow2*encoding.BYTES_PER_SYMBOL), nil
}

func IFFT(data []byte) ([]byte, error) {
	// we now IFFT data regardless of the encoding type
	// convert data to fr.Element
	dataFr, err := rs.ToFrArray(data)
	if err != nil {
		return nil, fmt.Errorf("error converting data to fr.Element: %w", err)
	}

	dataFrLen := len(dataFr)
	dataFrLenPow2 := encoding.NextPowerOf2(uint64(dataFrLen))

	// expand data to the next power of 2
	paddedDataFr := make([]fr.Element, dataFrLenPow2)
	for i := 0; i < len(paddedDataFr); i++ {
		if i < len(dataFr) {
			paddedDataFr[i].Set(&dataFr[i])
		} else {
			paddedDataFr[i].SetZero()
		}
	}

	maxScale := uint8(math.Log2(float64(dataFrLenPow2)))

	// perform IFFT
	fs := fft.NewFFTSettings(maxScale)
	dataIFFTFr, err := fs.FFT(paddedDataFr, true)
	if err != nil {
		return nil, fmt.Errorf("failed to perform IFFT: %w", err)
	}

	return rs.ToByteArray(dataIFFTFr, dataFrLenPow2*encoding.BYTES_PER_SYMBOL), nil
}
