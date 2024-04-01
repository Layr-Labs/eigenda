package encoding

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"

	"github.com/consensys/gnark-crypto/ecc/bn254"
)

func (c *Frame) Serialize() ([]byte, error) {
	return encode(c)
}

func (c *Frame) Deserialize(data []byte) (*Frame, error) {
	err := decode(data, c)
	if !c.Proof.IsInSubGroup() {
		return nil, fmt.Errorf("proof is in not the subgroup")
	}

	return c, err
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

func (c *G1Commitment) Serialize() ([]byte, error) {
	res := (*bn254.G1Affine)(c).Bytes()
	return res[:], nil
}

func (c *G1Commitment) Deserialize(data []byte) (*G1Commitment, error) {
	_, err := (*bn254.G1Affine)(c).SetBytes(data)
	if err != nil {
		return nil, err
	}
	return c, err
}

func (c *G1Commitment) UnmarshalJSON(data []byte) error {
	var g1Point bn254.G1Affine
	err := json.Unmarshal(data, &g1Point)
	if err != nil {
		return err
	}
	c.X = g1Point.X
	c.Y = g1Point.Y

	if !(*bn254.G1Affine)(c).IsInSubGroup() {
		return fmt.Errorf("G1Commitment not in the subgroup")
	}

	return nil
}

func (c *G2Commitment) Serialize() ([]byte, error) {
	res := (*bn254.G2Affine)(c).Bytes()
	return res[:], nil
}

func (c *G2Commitment) Deserialize(data []byte) (*G2Commitment, error) {
	_, err := (*bn254.G2Affine)(c).SetBytes(data)
	if err != nil {
		return nil, err
	}

	return c, err
}

func (c *G2Commitment) UnmarshalJSON(data []byte) error {
	var g2Point bn254.G2Affine
	err := json.Unmarshal(data, &g2Point)
	if err != nil {
		return err
	}
	c.X = g2Point.X
	c.Y = g2Point.Y

	if !(*bn254.G2Affine)(c).IsInSubGroup() {
		return fmt.Errorf("G2Commitment not in the subgroup")
	}
	return nil
}

func encode(obj any) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(obj)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func decode(data []byte, obj any) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(obj)
	if err != nil {
		return err
	}
	return nil
}
