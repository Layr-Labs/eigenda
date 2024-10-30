package rs

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"

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

// EncodeFrames serializes a slice of frames into a byte slice.
func EncodeFrames(frames []*Frame) ([]byte, error) {

	// Serialization format:
	// [number of frames: 4 byte uint32]
	// [size of frame 1: 4 byte uint32][frame 1]
	// [size of frame 2: 4 byte uint32][frame 2]
	// ...
	// [size of frame n: 4 byte uint32][frame n]

	encodedSize := 4
	encodedFrames := make([][]byte, len(frames))

	for i, frame := range frames {
		encodedSize += 4
		encodedFrame, err := frame.Encode()
		if err != nil {
			return nil, err
		}
		encodedFrames[i] = encodedFrame
		encodedSize += len(encodedFrame)
	}

	serializedBytes := make([]byte, encodedSize)
	binary.BigEndian.PutUint32(serializedBytes, uint32(len(frames)))
	index := 4

	for _, frameBytes := range encodedFrames {
		binary.BigEndian.PutUint32(serializedBytes[index:], uint32(len(frameBytes)))
		index += 4
		copy(serializedBytes[index:], frameBytes)
		index += len(frameBytes)
	}

	return serializedBytes, nil
}

// DecodeFrames deserializes a byte slice into a slice of frames.
func DecodeFrames(serializedFrames []byte) ([]*Frame, error) {
	frameCount := binary.BigEndian.Uint32(serializedFrames)
	index := 4

	frames := make([]*Frame, frameCount)

	for i := 0; i < int(frameCount); i++ {
		frameSize := binary.BigEndian.Uint32(serializedFrames[index:])
		index += 4
		frame, err := Decode(serializedFrames[index : index+int(frameSize)])
		if err != nil {
			return nil, fmt.Errorf("failed to decode frame %d: %w", i, err)
		}
		frames[i] = &frame
		index += int(frameSize)
	}

	return frames, nil
}
