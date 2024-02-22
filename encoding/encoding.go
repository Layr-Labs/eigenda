package encoding

import (
	"bytes"
	"encoding/gob"

	"github.com/Layr-Labs/eigenda/pkg/kzg/bn254"
)

// Proof is the multireveal proof
// Coeffs is identical to input data converted into Fr element
type Frame struct {
	Proof  bn254.G1Point
	Coeffs []bn254.Fr
}

func (f *Frame) Encode() ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(f)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

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
