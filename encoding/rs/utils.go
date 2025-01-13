package rs

import (
	"errors"
	"math"

	"github.com/Layr-Labs/eigenda/encoding"
	rb "github.com/Layr-Labs/eigenda/encoding/utils/reverseBits"

	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

// ToFrArray deserializes a byte array into a list of bn254 field elements.,
// where each 32-byte chunk needs to be a big-endian serialized bn254 field element.
// The last chunk is allowed to not have 32-bytes, and will be padded with zeroes
// on the right (so make sure that the last partial chunk represents a valid field element
// when padded with zeroes on the right and interpreted as big-endian).
//
// TODO: we should probably just force the data to be a multiple of 32 bytes.
// This would make the API and code simpler to read, and also allow the code
// to be auto-vectorized by the compiler (it probably isn't right now given the if inside the for loop).
func ToFrArray(data []byte) ([]fr.Element, error) {
	numEle := GetNumElement(uint64(len(data)), encoding.BYTES_PER_SYMBOL)
	eles := make([]fr.Element, numEle)

	for i := uint64(0); i < numEle; i++ {
		start := i * uint64(encoding.BYTES_PER_SYMBOL)
		end := (i + 1) * uint64(encoding.BYTES_PER_SYMBOL)
		if end >= uint64(len(data)) {
			padded := make([]byte, encoding.BYTES_PER_SYMBOL)
			copy(padded, data[start:])
			err := eles[i].SetBytesCanonical(padded)
			if err != nil {
				return nil, err
			}
		} else {
			err := eles[i].SetBytesCanonical(data[start:end])
			if err != nil {
				return nil, err
			}
		}
	}

	return eles, nil
}

// ToByteArray serializes a slice of fields elements to a slice of bytes.
// The byte array is created by serializing each Fr element in big-endian format.
// It is the reverse operation of ToFrArray.
func ToByteArray(dataFr []fr.Element, maxDataSize uint64) []byte {
	n := len(dataFr)
	dataSize := int(math.Min(
		float64(n*encoding.BYTES_PER_SYMBOL),
		float64(maxDataSize),
	))
	data := make([]byte, dataSize)
	for i := 0; i < n; i++ {
		v := dataFr[i].Bytes()

		start := i * encoding.BYTES_PER_SYMBOL
		end := (i + 1) * encoding.BYTES_PER_SYMBOL

		if uint64(end) > maxDataSize {
			copy(data[start:maxDataSize], v[:])
			break
		} else {
			copy(data[start:end], v[:])
		}
	}

	return data
}

func GetNumElement(dataLen uint64, CS int) uint64 {
	numEle := int(math.Ceil(float64(dataLen) / float64(CS)))
	return uint64(numEle)
}

// helper function
func RoundUpDivision(a, b uint64) uint64 {
	return uint64(math.Ceil(float64(a) / float64(b)))
}

func NextPowerOf2(d uint64) uint64 {
	nextPower := math.Ceil(math.Log2(float64(d)))
	return uint64(math.Pow(2.0, nextPower))
}

// This function is used by user to get the leading coset for a frame, where i is frame index
func GetLeadingCosetIndex(i uint64, numChunks uint64) (uint32, error) {

	if i < numChunks {
		j := rb.ReverseBitsLimited(uint32(numChunks), uint32(i))
		return j, nil
	} else {
		return 0, errors.New("cannot create number of frame higher than possible")
	}
}
