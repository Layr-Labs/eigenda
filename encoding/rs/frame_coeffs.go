package rs

import (
	"encoding/binary"
	"fmt"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

// FrameCoeffs is a slice of coefficients (i.e. an encoding.Frame object without the proofs).
type FrameCoeffs []fr.Element

// SerializeFrameCoeffsSlice serializes a slice FrameCoeffs into a binary format.
// Note that each FrameCoeffs object is required to have the exact same number of coefficients.
// Can be deserialized by DeserializeFrameCoeffsSlice().
//
// [number of elements per FrameCoeffs: 4 byte uint32]
// [coeffs FrameCoeffs 0, element 0][coeffs FrameCoeffs 0, element 1][coeffs FrameCoeffs 0, element 2]...
// [coeffs FrameCoeffs 1, element 0][coeffs FrameCoeffs 1, element 1][coeffs FrameCoeffs 1, element 2]...
// ...
// [coeffs FrameCoeffs n, element 0][coeffs FrameCoeffs n, element 1][coeffs FrameCoeffs n, element 2]...
//
// Where relevant, big endian encoding is used.
func SerializeFrameCoeffsSlice(coeffs []FrameCoeffs) ([]byte, error) {
	if len(coeffs) == 0 {
		return nil, fmt.Errorf("no frame coeffs to serialize")
	}

	elementCount := len(coeffs[0])
	bytesPerFrameCoeffs := encoding.BYTES_PER_SYMBOL * elementCount
	serializedSize := bytesPerFrameCoeffs*len(coeffs) + 4
	serializedBytes := make([]byte, serializedSize)

	binary.BigEndian.PutUint32(serializedBytes, uint32(elementCount))
	index := uint32(4)

	for _, coeff := range coeffs {
		if len(coeff) != elementCount {
			return nil, fmt.Errorf("frame coeffs have different number of elements, expected %d, got %d",
				elementCount, len(coeff))
		}
		for _, element := range coeff {
			serializedCoeff := element.Marshal()
			copy(serializedBytes[index:], serializedCoeff)
			index += encoding.BYTES_PER_SYMBOL
		}
	}

	return serializedBytes, nil
}

// DeserializeFrameCoeffsSlice is the inverse of SerializeFrameCoeffsSlice.
// It deserializes a byte slice into a slice of FrameCoeffs.
func DeserializeFrameCoeffsSlice(serializedData []byte) ([]FrameCoeffs, error) {
	// ElementCount is the number of elements in each FrameCoeffs object
	elementCount := binary.BigEndian.Uint32(serializedData)
	if elementCount == 0 {
		return nil, fmt.Errorf("element count cannot be 0")
	}

	index := uint32(4)

	// coeffsByteSize is the number of bytes required to store all the coefficients in a single FrameCoeffs object
	coeffsByteSize := encoding.BYTES_PER_SYMBOL * int(elementCount)
	remainingSize := len(serializedData[index:])
	if remainingSize%coeffsByteSize != 0 {
		return nil, fmt.Errorf("invalid data size: %d", len(serializedData))
	}
	coeffsCount := len(serializedData[index:]) / coeffsByteSize

	coeffs := make([]FrameCoeffs, coeffsCount)

	for i := 0; i < coeffsCount; i++ {
		coeffs[i] = make(FrameCoeffs, elementCount)
		for j := 0; j < int(elementCount); j++ {
			coeff := fr.Element{}
			coeff.Unmarshal(serializedData[index : index+encoding.BYTES_PER_SYMBOL])
			coeffs[i][j] = coeff
			index += encoding.BYTES_PER_SYMBOL
		}
	}

	return coeffs, nil
}

// SplitSerializedFrameCoeffs splits data as serialized by SerializeFrameCoeffsSlice into a slice of byte slices.
// Each byte slice contains the serialized data for a single FrameCoeffs object as serialized by FrameCoeffs.Serialize.
// Also returns ElementCount, the number of elements in each FrameCoeffs object.
func SplitSerializedFrameCoeffs(serializedData []byte) (elementCount uint32, binaryFrameCoeffs [][]byte, err error) {
	elementCount = binary.BigEndian.Uint32(serializedData)
	index := uint32(4)

	if elementCount == 0 {
		return 0, nil, fmt.Errorf("element count cannot be 0")
	}

	bytesPerFrameCoeffs := encoding.BYTES_PER_SYMBOL * elementCount
	remainingBytes := uint32(len(serializedData[index:]))
	if remainingBytes%bytesPerFrameCoeffs != 0 {
		return 0, nil, fmt.Errorf("invalid data size: %d", len(serializedData))
	}
	frameCoeffCount := uint32(len(serializedData[index:])) / bytesPerFrameCoeffs
	binaryFrameCoeffs = make([][]byte, frameCoeffCount)

	for i := uint32(0); i < frameCoeffCount; i++ {
		binaryFrameCoeffs[i] = make([]byte, bytesPerFrameCoeffs)
		copy(binaryFrameCoeffs[i], serializedData[index:index+bytesPerFrameCoeffs])
		index += bytesPerFrameCoeffs
	}

	return elementCount, binaryFrameCoeffs, nil
}
