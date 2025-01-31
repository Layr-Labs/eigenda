package rs

import (
	"encoding/binary"
	"fmt"
	"github.com/Layr-Labs/eigenda/core"
)

// BinaryFrames is analogous to an []encoding.Frame object, but in a serialized form.
type BinaryFrames struct {
	Frames       [][]byte
	ElementCount uint32
}

// BuildBinaryFrames takes serialized proofs and coefficients and builds a "binary frame", aka a serialized
// version of an encoding.Frame object. The format of the binary frame is the same as a frame within a
// core.Bundle object, but without header information.
//
// The format of a binary frame is simply a binary proof concatenated with binary coefficients.
func BuildBinaryFrames(
	proofs [][]byte,
	elementCount uint32,
	coefficients [][]byte) (*BinaryFrames, error) {

	if len(proofs) != len(coefficients) {
		return nil, fmt.Errorf("proofs and coefficients have different lengths (%d vs %d)",
			len(proofs), len(coefficients))
	}

	binaryFrames := make([][]byte, len(proofs))

	for i := 0; i < len(proofs); i++ {
		binaryFrame := make([]byte, len(proofs[i])+len(coefficients[i]))
		copy(binaryFrame, proofs[i])
		copy(binaryFrame[len(proofs[i]):], coefficients[i])
		binaryFrames[i] = binaryFrame
	}

	return &BinaryFrames{
		Frames:       binaryFrames,
		ElementCount: elementCount,
	}, nil
}

// SerializeAsBundle serializes a BinaryFrames object into a binary bundle format.
//
// Bundle format is as follows:
//
// [serialization protocol version: 1 byte]
// [number of elements per frame: 7 bytes]
// for each frame:
//
//	[proof: 32 bytes]
//	[coefficients: 32 bytes * element count]
func (b *BinaryFrames) SerializeAsBundle() ([]byte, error) {

	if len(b.Frames) == 0 {
		return nil, fmt.Errorf("no binary Frames to serialize")
	}

	lengthPerFrame := len(b.Frames[0])
	dataSize := 8 + len(b.Frames)*lengthPerFrame

	data := make([]byte, dataSize)

	header := core.BinaryBundleHeader(uint64(b.ElementCount))
	binary.LittleEndian.PutUint64(data, header) // LittleEndian... ಠ_ಠ
	index := 8

	for _, frame := range b.Frames {
		copy(data[index:], frame)
		index += lengthPerFrame
	}

	return data, nil
}

// GetApproximateSize returns the approximate size of the BinaryFrames object in memory, in bytes.
func (b *BinaryFrames) GetApproximateSize() uint64 {
	// It is safe to assume that each frame is the same size.
	return uint64(len(b.Frames) * len(b.Frames[0]))
}
