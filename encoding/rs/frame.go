package rs

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"

	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

// Proof is the multireveal proof
// Coeffs is identical to input data converted into Fr element
type Frame struct {
	Coeffs []fr.Element
}

// Encode serializes the frame into a byte slice.
func (f *Frame) Encode() ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(f)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Decode deserializes a byte slice into a frame.
func Decode(b []byte) (Frame, error) {
	var f Frame
	buf := bytes.NewBuffer(b)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&f)
	if err != nil {
		return Frame{}, err
	}
	return f, nil
}

// GnarkEncodeFrames serializes a slice of frames into a byte slice.
//
// Serialization format:
// [number of frames: 4 byte uint32]
// [size of frame 1: 4 byte uint32][frame 1]
// [size of frame 2: 4 byte uint32][frame 2]
// ...
// [size of frame n: 4 byte uint32][frame n]
//
// Where relevant, big endian encoding is used.
func GnarkEncodeFrames(frames []*Frame) ([]byte, error) {

	// Count the number of bytes.
	encodedSize := uint32(4) // stores the number of frames
	for _, frame := range frames {
		encodedSize += 4                     // stores the size of the frame
		encodedSize += GnarkFrameSize(frame) // size of the frame
	}

	serializedBytes := make([]byte, encodedSize)
	binary.BigEndian.PutUint32(serializedBytes, uint32(len(frames)))
	index := uint32(4)

	for _, frame := range frames {
		index += GnarkEncodeFrame(frame, serializedBytes[index:])
	}

	if index != encodedSize {
		// Sanity check, this should never happen.
		return nil, fmt.Errorf("encoded size mismatch: expected %d, got %d", encodedSize, index)
	}

	return serializedBytes, nil
}

// GnarkEncodeFrame serializes a frame into a target byte slice. Returns the number of bytes written.
func GnarkEncodeFrame(frame *Frame, target []byte) uint32 {
	binary.BigEndian.PutUint32(target, uint32(len(frame.Coeffs)))
	index := uint32(4)

	for _, coeff := range frame.Coeffs {
		serializedCoeff := coeff.Marshal()
		copy(target[index:], serializedCoeff)
		index += uint32(len(serializedCoeff))
	}

	return index
}

// GnarkFrameSize returns the size of a frame in bytes.
func GnarkFrameSize(frame *Frame) uint32 {
	return uint32(encoding.BYTES_PER_SYMBOL * len(frame.Coeffs))
}

// GnarkDecodeFrames deserializes a byte slice into a slice of frames.
func GnarkDecodeFrames(serializedFrames []byte) ([]*Frame, error) {
	frameCount := binary.BigEndian.Uint32(serializedFrames)
	index := uint32(4)

	frames := make([]*Frame, frameCount)

	for i := 0; i < int(frameCount); i++ {
		frame, bytesRead, err := GnarkDecodeFrame(serializedFrames[index:])

		if err != nil {
			return nil, fmt.Errorf("failed to decode frame %d: %w", i, err)
		}

		frames[i] = frame
		index += bytesRead
	}

	if index != uint32(len(serializedFrames)) {
		return nil, fmt.Errorf("decoded size mismatch: expected %d, got %d", len(serializedFrames), index)
	}

	return frames, nil
}

// GnarkDecodeFrame deserializes a byte slice into a frame. Returns the frame and the number of bytes read.
// If passed a byte array that contains multiple frames, it will only decode the first frame. The uint32
// returned is the number of bytes read from the input slice.
func GnarkDecodeFrame(serializedFrame []byte) (*Frame, uint32, error) {
	if len(serializedFrame) < 4 {
		return nil, 0, fmt.Errorf("invalid frame size: %d", len(serializedFrame))
	}

	symbolCount := binary.BigEndian.Uint32(serializedFrame)
	index := uint32(4)

	if len(serializedFrame) < int(index+symbolCount*encoding.BYTES_PER_SYMBOL) {
		return nil, 0, fmt.Errorf("invalid frame size: %d", len(serializedFrame))
	}

	coeffs := make([]fr.Element, symbolCount)
	for i := 0; i < int(symbolCount); i++ {
		coeff := fr.Element{}
		coeff.Unmarshal(serializedFrame[index : index+encoding.BYTES_PER_SYMBOL])
		coeffs[i] = coeff
		index += uint32(encoding.BYTES_PER_SYMBOL)
	}

	frame := &Frame{Coeffs: coeffs}

	return frame, index, nil
}

// GnarkSplitBinaryFrames deserializes a serialized list of frames into slice of individually serialized frames
// (which can be individually deserialized by GnarkDecodeFrame).
func GnarkSplitBinaryFrames(serializedFrames []byte) ([][]byte, error) {
	frameCount := binary.BigEndian.Uint32(serializedFrames)
	index := uint32(4)

	frameBytes := make([][]byte, frameCount)
	for i := 0; i < int(frameCount); i++ {
		symbolCount := binary.BigEndian.Uint32(serializedFrames[index:])
		frameLength := 4 + symbolCount*encoding.BYTES_PER_SYMBOL
		if len(serializedFrames) < int(index+frameLength) {
			return nil, fmt.Errorf("invalid frame size: %d", len(serializedFrames))
		}
		frameBytes[i] = serializedFrames[index : index+frameLength]

		index += frameLength
	}

	if index != uint32(len(serializedFrames)) {
		return nil, fmt.Errorf("decoded size mismatch: expected %d, got %d", len(serializedFrames), index)
	}

	return frameBytes, nil
}

// CombineBinaryFrames combines a slice of serialized frames into a single serialized byte slice. This is the inverse
// of GnarkSplitBinaryFrames, and produces bytes that can be deserialized by GnarkDecodeFrames.
func CombineBinaryFrames(frameBytes [][]byte) []byte {
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
