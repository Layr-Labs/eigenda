package encoding

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/consensys/gnark-crypto/ecc/bn254"
)

// SerializeGob serializes the Frame into a byte slice using gob encoding.
// TODO(samlaf): when do we use gob vs gnark serialization ([Frame.SerializeGnark])?
func (c *Frame) SerializeGob() ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(c)
	if err != nil {
		return nil, fmt.Errorf("gob encode: %w", err)
	}
	return buf.Bytes(), nil
}

// DeserializeGob deserializes the byte slice into a Frame using gob decoding.
func (c *Frame) DeserializeGob(data []byte) (*Frame, error) {
	buf := bytes.NewBuffer(data)
	err := gob.NewDecoder(buf).Decode(c)
	if err != nil {
		return nil, fmt.Errorf("gob decode: %w", err)
	}

	// TODO(samlaf): why do we check this here?
	if !c.Proof.IsInSubGroup() {
		return nil, fmt.Errorf("proof is in not the subgroup")
	}

	return c, nil
}

// SerializeGnark serializes the Frame into a byte slice using gnark encoding.
func (c *Frame) SerializeGnark() ([]byte, error) {
	coded := make([]byte, 0, bn254.SizeOfG1AffineCompressed+BYTES_PER_SYMBOL*len(c.Coeffs))
	// This is compressed format with just 32 bytes.
	proofBytes := c.Proof.Bytes()
	coded = append(coded, proofBytes[:]...)
	for _, coeff := range c.Coeffs {
		coded = append(coded, coeff.Marshal()...)
	}
	return coded, nil
}

// DeserializeGnark deserializes the byte slice into a Frame using gnark decoding.
func (c *Frame) DeserializeGnark(data []byte) (*Frame, error) {
	if len(data) <= bn254.SizeOfG1AffineCompressed {
		return nil, fmt.Errorf("chunk length must be at least %d: %d given", bn254.SizeOfG1AffineCompressed, len(data))
	}
	var f Frame
	buf := data
	err := f.Proof.Unmarshal(buf[:bn254.SizeOfG1AffineCompressed])
	if err != nil {
		return nil, err
	}
	buf = buf[bn254.SizeOfG1AffineCompressed:]
	if len(buf)%BYTES_PER_SYMBOL != 0 {
		return nil, errors.New("invalid chunk length")
	}
	f.Coeffs = make([]Symbol, len(buf)/BYTES_PER_SYMBOL)
	i := 0
	for len(buf) > 0 {
		if len(buf) < BYTES_PER_SYMBOL {
			return nil, errors.New("invalid chunk length")
		}
		f.Coeffs[i].Unmarshal(buf[:BYTES_PER_SYMBOL])
		i++
		buf = buf[BYTES_PER_SYMBOL:]
	}
	return &f, nil
}

const SerializedProofLength = bn254.SizeOfG1AffineCompressed

// SerializeFrameProof serializes an [Proof] to the target byte array.
// Only the first SerializedProofLength bytes of the target array are written to.
func SerializeFrameProof(proof *Proof, target []byte) error {
	if len(target) < SerializedProofLength {
		return fmt.Errorf("target byte array is too short")
	}
	proofBytes := proof.Bytes()
	copy(target, proofBytes[:])

	return nil
}

// SerializeFrameProofs serializes a slice of proofs (as found in [Proof], but without the coefficients)
// into a binary format.
func SerializeFrameProofs(proofs []*Proof) ([]byte, error) {
	bytes := make([]byte, SerializedProofLength*len(proofs))
	for index, proof := range proofs {
		err := SerializeFrameProof(proof, bytes[index*SerializedProofLength:])
		if err != nil {
			return nil, fmt.Errorf("failed to serialize proof: %w", err)
		}
	}
	return bytes, nil
}

// DeserializeFrameProof deserializes an [Proof]. Only the first proof is deserialized
// from the first SerializedProofLength bytes of the input array.
func DeserializeFrameProof(bytes []byte) (*Proof, error) {
	if len(bytes) != SerializedProofLength {
		return nil, fmt.Errorf("unexpected proof length: expected %d, got %d", SerializedProofLength, len(bytes))
	}
	proof := Proof{}
	err := proof.Unmarshal(bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal proof: %w", err)
	}
	return &proof, nil
}

// DeserializeFrameProofs deserializes a slice of proofs (as found in [Proof], but without the coefficients)
// from a binary format. The inverse of SerializeFrameProofs.
func DeserializeFrameProofs(bytes []byte) ([]*Proof, error) {
	if len(bytes)%SerializedProofLength != 0 {
		return nil, fmt.Errorf("input byte array is not a multiple of proof length")
	}

	splitProofs, err := SplitSerializedFrameProofs(bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to split proofs: %w", err)
	}

	return DeserializeSplitFrameProofs(splitProofs), nil
}

// SplitSerializedFrameProofs splits a serialized slice of proofs (as found in [Proof], but without
// the coefficients) into a slice of byte slices, each containing a single serialized proof. Each individual
// serialized proof can be deserialized by [Proof.Unmarshal].
func SplitSerializedFrameProofs(bytes []byte) ([][]byte, error) {
	if len(bytes)%SerializedProofLength != 0 {
		return nil, fmt.Errorf("input byte array is not a multiple of proof length")
	}

	proofCount := len(bytes) / SerializedProofLength
	proofs := make([][]byte, proofCount)

	for i := 0; i < proofCount; i++ {
		proofs[i] = bytes[i*SerializedProofLength : (i+1)*SerializedProofLength]
	}

	return proofs, nil
}

// DeserializeSplitFrameProofs deserializes a slice of byte slices into a slice of Proof objects.
func DeserializeSplitFrameProofs(proofs [][]byte) []*Proof {
	proofsSlice := make([]*Proof, len(proofs))
	for i, proof := range proofs {
		proofsSlice[i], _ = DeserializeFrameProof(proof)
	}
	return proofsSlice
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
