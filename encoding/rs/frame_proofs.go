package rs

import (
	"fmt"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/consensys/gnark-crypto/ecc/bn254"
)

const SerializedProofLength = bn254.SizeOfG1AffineCompressed

// SerializeFrameProof serializes an encoding.Proof to the target byte array.
// Only the first SerializedProofLength bytes of the target array are written to.
func SerializeFrameProof(proof *encoding.Proof, target []byte) error {
	if len(target) < SerializedProofLength {
		return fmt.Errorf("target byte array is too short")
	}
	proofBytes := proof.Bytes()
	copy(target, proofBytes[:])

	return nil
}

// SerializeFrameProofs serializes a slice of proofs (as found in encoding.Proof, but without the coefficients)
// into a binary format.
func SerializeFrameProofs(proofs []*encoding.Proof) ([]byte, error) {
	bytes := make([]byte, SerializedProofLength*len(proofs))
	for index, proof := range proofs {
		err := SerializeFrameProof(proof, bytes[index*SerializedProofLength:])
		if err != nil {
			return nil, fmt.Errorf("failed to serialize proof: %w", err)
		}
	}
	return bytes, nil
}

// DeserializeFrameProof deserializes an encoding.Proof. Only the first proof is deserialized
// from the first SerializedProofLength bytes of the input array.
func DeserializeFrameProof(bytes []byte) (*encoding.Proof, error) {
	if len(bytes) != SerializedProofLength {
		return nil, fmt.Errorf("unexpected proof length: expected %d, got %d", SerializedProofLength, len(bytes))
	}
	proof := encoding.Proof{}
	err := proof.Unmarshal(bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal proof: %w", err)
	}
	return &proof, nil
}

// DeserializeFrameProofs deserializes a slice of proofs (as found in encoding.Proof, but without the coefficients)
// from a binary format. The inverse of SerializeFrameProofs.
func DeserializeFrameProofs(bytes []byte) ([]*encoding.Proof, error) {
	if len(bytes)%SerializedProofLength != 0 {
		return nil, fmt.Errorf("input byte array is not a multiple of proof length")
	}

	splitProofs, err := SplitSerializedFrameProofs(bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to split proofs: %w", err)
	}

	return DeserializeSplitFrameProofs(splitProofs), nil
}

// SplitSerializedFrameProofs splits a serialized slice of proofs (as found in encoding.Proof, but without
// the coefficients) into a slice of byte slices, each containing a single serialized proof. Each individual
// serialized proof can be deserialized by encoding.Proof.Unmarshal.
func SplitSerializedFrameProofs(bytes []byte) ([][]byte, error) {
	if len(bytes)%SerializedProofLength != 0 {
		return nil, fmt.Errorf("input byte array is not a multiple of proof length")
	}

	proofCount := len(bytes) / SerializedProofLength
	proofs := make([][]byte, proofCount)

	for i := 0; i < proofCount; i++ {
		proof := make([]byte, SerializedProofLength)
		copy(proof, bytes[i*SerializedProofLength:(i+1)*SerializedProofLength])
		proofs[i] = proof
	}

	return proofs, nil
}

// DeserializeSplitFrameProofs deserializes a slice of byte slices into a slice of encoding.Proof objects.
func DeserializeSplitFrameProofs(proofs [][]byte) []*encoding.Proof {
	proofsSlice := make([]*encoding.Proof, len(proofs))
	for i, proof := range proofs {
		proofsSlice[i], _ = DeserializeFrameProof(proof)
	}
	return proofsSlice
}
