package rs

import (
	"encoding/binary"
	"fmt"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

// TODO all "Frame" references in function names and docs
// TODO should this live in the same package as encoding.Frame?

// FrameCoeffs is a slice of coefficients (i.e. an encoding.Frame object without the proofs).
type FrameCoeffs []fr.Element

// SerializeFrameCoeffsSlice serializes a slice FrameCoeffs into a binary format.
// Can be deserialized by DeserializeFrameCoeffsSlice().
//
// Serialization format:
// [number of coeffs: 4 byte uint32]
// [size of coeffs 1: 4 byte uint32][coeffs 1]
// [size of coeffs 2: 4 byte uint32][coeffs 2]
// ...
// [size of coeffs n: 4 byte uint32][coeffs n]
//
// Where relevant, big endian encoding is used.
func SerializeFrameCoeffsSlice(coeffs []FrameCoeffs) ([]byte, error) {

	// Count the number of bytes.
	encodedSize := uint32(4) // stores the number of coeffs
	for _, coeff := range coeffs {
		encodedSize += 4                // stores the size of the coeff
		encodedSize += CoeffSize(coeff) // size of the coeff
	}

	serializedBytes := make([]byte, encodedSize)
	binary.BigEndian.PutUint32(serializedBytes, uint32(len(coeffs)))
	index := uint32(4)

	for _, frame := range coeffs {
		index += frame.Serialize(serializedBytes[index:])
	}

	if index != encodedSize {
		// Sanity check, this should never happen.
		return nil, fmt.Errorf("encoded size mismatch: expected %d, got %d", encodedSize, index)
	}

	return serializedBytes, nil
}

// Serialize serializes a FrameCoeffs object into a byte slice.
func (c FrameCoeffs) Serialize(target []byte) uint32 {
	binary.BigEndian.PutUint32(target, uint32(len(c)))
	index := uint32(4)

	for _, coeff := range c {
		serializedCoeff := coeff.Marshal()
		copy(target[index:], serializedCoeff)
		index += uint32(len(serializedCoeff))
	}

	return index
}

// CoeffSize returns the size of a frame in bytes.
func CoeffSize(coeffs FrameCoeffs) uint32 { // TODO don't export this!
	return uint32(encoding.BYTES_PER_SYMBOL * len(coeffs))
}

// DeserializeFrameCoeffsSlice is the inverse of SerializeFrameCoeffsSlice.
// It deserializes a byte slice into a slice of FrameCoeffs.
func DeserializeFrameCoeffsSlice(serializedData []byte) ([]FrameCoeffs, error) {
	frameCount := binary.BigEndian.Uint32(serializedData)
	index := uint32(4)

	coeffs := make([]FrameCoeffs, frameCount)

	for i := 0; i < int(frameCount); i++ {
		coeff, bytesRead, err := DeserializeFrameCoeffs(serializedData[index:])

		if err != nil {
			return nil, fmt.Errorf("failed to decode coeff %d: %w", i, err)
		}

		coeffs[i] = coeff
		index += bytesRead
	}

	if index != uint32(len(serializedData)) {
		return nil, fmt.Errorf("decoded size mismatch: expected %d, got %d", len(serializedData), index)
	}

	return coeffs, nil
}

// DeserializeFrameCoeffs is the inverse of FrameCoeffs.Serialize(). It deserializes a byte slice into a
// FrameCoeffs object. If passed a byte array that contains multiple serialized FrameCoeffs, it will only
// deserialize the first one. The uint32 returned is the number of bytes read from the input slice.
func DeserializeFrameCoeffs(serializedData []byte) (FrameCoeffs, uint32, error) {
	if len(serializedData) < 4 {
		return nil, 0, fmt.Errorf("invalid data size: %d", len(serializedData))
	}

	symbolCount := binary.BigEndian.Uint32(serializedData)
	index := uint32(4)

	if len(serializedData) < int(index+symbolCount*encoding.BYTES_PER_SYMBOL) {
		return nil, 0, fmt.Errorf("invalid data size: %d", len(serializedData))
	}

	coeffs := make([]fr.Element, symbolCount)
	for i := 0; i < int(symbolCount); i++ {
		coeff := fr.Element{}
		coeff.Unmarshal(serializedData[index : index+encoding.BYTES_PER_SYMBOL])
		coeffs[i] = coeff
		index += uint32(encoding.BYTES_PER_SYMBOL)
	}

	return coeffs, index, nil
}

// SplitSerializedFrameCoeffs splits data as serialized by SerializeFrameCoeffsSlice into a slice of byte slices.
// Each byte slice contains the serialized data for a single FrameCoeffs object as serialized by FrameCoeffs.Serialize.
func SplitSerializedFrameCoeffs(serializedData []byte) ([][]byte, error) {
	objectCount := binary.BigEndian.Uint32(serializedData)
	index := uint32(4)

	coeffsBytes := make([][]byte, objectCount)
	for i := 0; i < int(objectCount); i++ {
		symbolCount := binary.BigEndian.Uint32(serializedData[index:])
		coeffsLength := 4 + symbolCount*encoding.BYTES_PER_SYMBOL
		if len(serializedData) < int(index+coeffsLength) {
			return nil, fmt.Errorf("invalid frame size: %d", len(serializedData))
		}
		coeffsBytes[i] = serializedData[index : index+coeffsLength]

		index += coeffsLength
	}

	if index != uint32(len(serializedData)) {
		return nil, fmt.Errorf("decoded size mismatch: expected %d, got %d", len(serializedData), index)
	}

	return coeffsBytes, nil
}

// TODO this method may not be needed

// CombineSerializedFrameCoeffs combines a slice of serialized FrameCoeffs into a single serialized byte slice.
// This is the inverse of SplitSerializedFrameCoeffs, and produces bytes that can be deserialized by
// DeserializeFrameCoeffsSlice.
func CombineSerializedFrameCoeffs(frameBytes [][]byte) []byte {
	length := uint32(4)
	for _, frame := range frameBytes {
		length += uint32(len(frame))
	}

	result := make([]byte, length)

	// first four bytes are the number of frames
	binary.BigEndian.PutUint32(result, uint32(len(frameBytes)))
	index := uint32(4)
	for _, frame := range frameBytes {
		copy(result[index:], frame)
		index += uint32(len(frame))
	}

	return result
}
