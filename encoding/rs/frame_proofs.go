package rs

import (
	"fmt"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/consensys/gnark-crypto/ecc/bn254"
)

// SerializeFrameProofs serializes a slice of proofs (as found in encoding.Proof, but without the coefficients)
// into a binary format.
func SerializeFrameProofs(proofs []*encoding.Proof) []byte {
	bytes := make([]byte, 0, bn254.SizeOfG1AffineCompressed*len(proofs))
	for _, proof := range proofs {
		proofBytes := proof.Bytes()
		bytes = append(bytes, proofBytes[:]...)
	}
	return bytes
}

// DeserializeFrameProofs deserializes a slice of proofs (as found in encoding.Proof, but without the coefficients)
// from a binary format. The inverse of SerializeFrameProofs.
func DeserializeFrameProofs(bytes []byte) ([]*encoding.Proof, error) {
	proofCount := len(bytes) / bn254.SizeOfG1AffineCompressed
	proofs := make([]*encoding.Proof, proofCount)

	for i := 0; i < proofCount; i++ {
		proof := encoding.Proof{}
		err := proof.Unmarshal(bytes[i*bn254.SizeOfG1AffineCompressed:])
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal proof: %w", err)
		}
		proofs[i] = &proof
	}

	return proofs, nil
}

// SplitSerializedFrameProofs splits a serialized slice of proofs (as found in encoding.Proof, but without
// the coefficients) into a slice of byte slices, each containing a single serialized proof. Each individual
// serialized proof can be deserialized by encoding.Proof.Unmarshal.
func SplitSerializedFrameProofs(bytes []byte) ([][]byte, error) {
	proofCount := len(bytes) / bn254.SizeOfG1AffineCompressed
	proofs := make([][]byte, proofCount)

	for i := 0; i < proofCount; i++ {
		proof := make([]byte, bn254.SizeOfG1AffineCompressed)
		copy(proof, bytes[i*bn254.SizeOfG1AffineCompressed:])
		proofs[i] = proof
	}

	return proofs, nil
}
