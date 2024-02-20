package encoder

import (
	"errors"
	"math"

	rb "github.com/Layr-Labs/eigenda/encoding/utils/reverseBits"

	bls "github.com/Layr-Labs/eigenda/pkg/kzg/bn254"
)

func ToFrArray(data []byte) []bls.Fr {
	//numEle := int(math.Ceil(float64(len(data)) / float64(BYTES_PER_COEFFICIENT)))
	numEle := GetNumElement(uint64(len(data)), bls.BYTES_PER_COEFFICIENT)
	eles := make([]bls.Fr, numEle)

	for i := uint64(0); i < numEle; i++ {
		start := i * uint64(bls.BYTES_PER_COEFFICIENT)
		end := (i + 1) * uint64(bls.BYTES_PER_COEFFICIENT)
		if end >= uint64(len(data)) {
			var padded [31]byte
			copy(padded[:], data[start:])
			bls.FrSetBytes(&eles[i], padded[:])
		} else {
			bls.FrSetBytes(&eles[i], data[start:end])
		}
	}

	return eles
}

// ToByteArray converts a list of Fr to a byte array
func ToByteArray(dataFr []bls.Fr, maxDataSize uint64) []byte {
	n := len(dataFr)
	dataSize := int(math.Min(
		float64(n*bls.BYTES_PER_COEFFICIENT),
		float64(maxDataSize),
	))
	data := make([]byte, dataSize)
	for i := 0; i < n; i++ {
		v := bls.FrToBytes(&dataFr[i])

		start := i * bls.BYTES_PER_COEFFICIENT
		end := (i + 1) * bls.BYTES_PER_COEFFICIENT

		if uint64(end) > maxDataSize {
			copy(data[start:maxDataSize], v[1:])
			break
		} else {
			copy(data[start:end], v[1:])
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
