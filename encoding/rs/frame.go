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
func GnarkDecodeFrame(serializedFrame []byte) (*Frame, uint32, error) {
	if len(serializedFrame) < 4 {
		return nil, 0, fmt.Errorf("invalid frame size: %d", len(serializedFrame))
	}

	frameCount := binary.BigEndian.Uint32(serializedFrame)
	index := uint32(4)

	if len(serializedFrame) < int(index+frameCount*encoding.BYTES_PER_SYMBOL) {
		return nil, 0, fmt.Errorf("invalid frame size: %d", len(serializedFrame))
	}

	coeffs := make([]fr.Element, frameCount)
	for i := 0; i < int(frameCount); i++ {
		coeff := fr.Element{}
		coeff.Unmarshal(serializedFrame[index : index+encoding.BYTES_PER_SYMBOL])
		coeffs[i] = coeff
		index += uint32(encoding.BYTES_PER_SYMBOL)
	}

	frame := &Frame{Coeffs: coeffs}

	return frame, index, nil
}
