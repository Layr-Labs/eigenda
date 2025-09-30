// Note that all functions are suffixed with V2 to avoid passing a V1 backend to a V2 prover.
// The main difference is that the V2 prover requires blobs to be of power-of-two size.
package prover

import (
	"github.com/Layr-Labs/eigenda/encoding/kzg/prover/v2/gnark"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

// Proof device represents a backend capable of computing KZG multiproofs.
type KzgMultiProofsBackendV2 interface {
	ComputeMultiFrameProofV2(blobFr []fr.Element, numChunks, chunkLen, numWorker uint64) ([]bn254.G1Affine, error)
}

var _ KzgMultiProofsBackendV2 = &gnark.KzgMultiProofGnarkBackend{}

// CommitmentDevice represents a backend capable of computing various KZG commitments.
type KzgCommitmentsBackendV2 interface {
	ComputeCommitmentV2(coeffs []fr.Element) (*bn254.G1Affine, error)
	ComputeLengthCommitmentV2(coeffs []fr.Element) (*bn254.G2Affine, error)
	ComputeLengthProofV2(coeffs []fr.Element) (*bn254.G2Affine, error)
	ComputeLengthProofForLengthV2(blobFr []fr.Element, length uint32) (*bn254.G2Affine, error)
}

var _ KzgCommitmentsBackendV2 = &gnark.KzgCommitmentsGnarkBackend{}
