package prover

import (
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

// Proof device represents a backend capable of computing KZG multiproofs.
type KzgMultiProofsBackend interface {
	ComputeMultiFrameProof(blobFr []fr.Element, numChunks, chunkLen, numWorker uint64) ([]bn254.G1Affine, error)
}

// CommitmentDevice represents a backend capable of computing various KZG commitments.
type KzgCommitmentsBackend interface {
	ComputeCommitment(coeffs []fr.Element) (*bn254.G1Affine, error)
	ComputeLengthCommitment(coeffs []fr.Element) (*bn254.G2Affine, error)
	ComputeLengthProof(coeffs []fr.Element) (*bn254.G2Affine, error)
	ComputeLengthProofForLength(blobFr []fr.Element, length uint64) (*bn254.G2Affine, error)
}
