package rs

import (
	"encoding/binary"
	"fmt"
	"github.com/Layr-Labs/eigenda/encoding"
)

// BuildBinaryFrame creates a new binary frame from the binary proof and binary coefficients.
// A binary frame can be parsed into an encoding.Frame. It has a MUCH smaller memory footprint
// a parsed encoding.Frame.
//
// The format of the binary frame is as follows:
// [binary proof]       // length SerializedProofLength, as serialized by SerializeFrameProof()
// [coeffs]             // variable length, as serialized by FrameCoeffs.Serialize()
func BuildBinaryFrame(binaryProof []byte, binaryCoeffs []byte) []byte {
	length := len(binaryProof) + len(binaryCoeffs)
	data := make([]byte, length)
	copy(data, binaryProof)
	copy(data[len(binaryProof):], binaryCoeffs)
	return data
}

// ParseBinaryFrame parses a binary frame into an encoding.Frame (as built by BuildBinaryFrame).
func ParseBinaryFrame(data []byte) (*encoding.Frame, error) {
	proof, err := DeserializeFrameProof(data[:SerializedProofLength])
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize proof: %w", err)
	}

	coeffs, _, err := DeserializeFrameCoeffs(data[SerializedProofLength:])
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize coeffs: %w", err)
	}

	return &encoding.Frame{
		Proof:  *proof,
		Coeffs: coeffs,
	}, nil
}

// CombineBinaryFrames combines multiple binary frames into a single byte array.
//
// Syntax:
// [number of frames]        // 4 bytes, unsigned big-endian
// [frame 1 length in bytes] // 4 bytes, unsigned big-endian
// [frame 1]                 // variable length
// [frame 2 length in bytes] // 4 bytes, unsigned big-endian
// [frame 2]                 // variable length
// ...
// [frame n length in bytes] // 4 bytes, unsigned big-endian
// [frame n]                 // variable length
func CombineBinaryFrames(frames [][]byte) []byte {
	length := uint32(4)
	for _, frame := range frames {
		length += uint32(len(frame)) + 4
	}
	data := make([]byte, length)

	binary.BigEndian.PutUint32(data, uint32(len(frames)))
	index := uint32(4)

	for _, frame := range frames {
		binary.BigEndian.PutUint32(data[index:], uint32(len(frame)))
		index += 4
		copy(data[index:], frame)
		index += uint32(len(frame))
	}
	return data
}

// SplitBinaryFrames splits a byte array containing multiple binary frames and splits them into individual binary
// frames. This is the inverse of CombineBinaryFrames.
func SplitBinaryFrames(data []byte) ([][]byte, error) {
	if len(data) < 4 {
		return nil, fmt.Errorf("invalid data size, can not get number of frames: %d", len(data))
	}
	count := binary.BigEndian.Uint32(data)
	index := uint32(4)

	frames := make([][]byte, count)

	for i := 0; i < int(count); i++ {
		if len(data) < int(index+4) {
			return nil, fmt.Errorf("invalid data size, can not read frame length: %d", len(data))
		}
		frameLength := binary.BigEndian.Uint32(data[index:])
		index += 4

		if len(data) < int(index+frameLength) {
			return nil, fmt.Errorf("invalid data size, incomplete frame: %d", len(data))
		}

		frame := make([]byte, frameLength)
		frames[i] = frame
		copy(frame, data[index:])
		index += frameLength
	}
	return frames, nil
}

// DeserializeBinaryFrames parses an array of binary frames into an array of encoding.Frames.
func DeserializeBinaryFrames(data []byte) ([]*encoding.Frame, error) {
	frames, err := SplitBinaryFrames(data)
	if err != nil {
		return nil, fmt.Errorf("failed to split binary frames: %w", err)
	}

	parsedFrames := make([]*encoding.Frame, len(frames))
	for i, frame := range frames {
		parsedFrame, err := ParseBinaryFrame(frame)
		if err != nil {
			return nil, fmt.Errorf("failed to parse binary frame: %w", err)
		}
		parsedFrames[i] = parsedFrame
	}
	return parsedFrames, nil
}
