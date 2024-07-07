package prover

import (
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

type ProofComputeDevice interface {
	// blobFr are coefficients
	ComputeCommitment(blobFr []fr.Element) (*bn254.G1Affine, error)
	ComputeMultiFrameProof(blobFr []fr.Element, numChunks, chunkLen, numWorker uint64) ([]bn254.G1Affine, error)
	ComputeLengthCommitment(blobFr []fr.Element) (*bn254.G2Affine, error)
	ComputeLengthProof(blobFr []fr.Element) (*bn254.G2Affine, error)
}
