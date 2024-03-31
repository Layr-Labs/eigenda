package encoding

import (
	"golang.org/x/exp/constraints"

	"math"
)

func PadToPowerOf2Frames(data []byte) []byte {

	// Pad data to the next power of 2
	length := roundUpDivide[int](len(data), BYTES_PER_COEFFICIENT)

	// Pad data to the next power of 2
	padSize := NextPowerOf2(uint64(length)) * BYTES_PER_COEFFICIENT

	paddedData := make([]byte, padSize)
	copy(paddedData, data)

	return paddedData
}

// GetBlobLength converts from blob size in bytes to blob size in symbols
func GetBlobLength(blobSize uint) uint {
	symSize := uint(BYTES_PER_COEFFICIENT)
	return (blobSize + symSize - 1) / symSize
}

// GetBlobSize converts from blob length in symbols to blob size in bytes. This is not an exact conversion.
func GetBlobSize(blobLength uint) uint {
	return blobLength * BYTES_PER_COEFFICIENT
}

// GetBlobLength converts from blob size in bytes to blob size in symbols
func GetEncodedBlobLength(blobLength uint, quorumThreshold, advThreshold uint8) uint {
	return roundUpDivide(blobLength*100, uint(quorumThreshold-advThreshold))
}

func NextPowerOf2(d uint64) uint64 {
	nextPower := math.Ceil(math.Log2(float64(d)))
	return uint64(math.Pow(2.0, nextPower))
}

func roundUpDivide[T constraints.Integer](a, b T) T {
	return (a + b - 1) / b

}
