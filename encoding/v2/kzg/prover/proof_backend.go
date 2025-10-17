// Note that all functions are suffixed with V2 to avoid passing a V1 backend to a V2 prover.
// The main difference is that the V2 prover requires blobs to be of power-of-two size.
package prover

import (
	"github.com/Layr-Labs/eigenda/encoding/v2/kzg/prover/gnark"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

// Proof device represents a backend capable of computing KZG multiproofs.
type KzgMultiProofsBackendV2 interface {
	ComputeMultiFrameProofV2(blobFr []fr.Element, numChunks, chunkLen, numWorker uint64) ([]bn254.G1Affine, error)
}

var _ KzgMultiProofsBackendV2 = &gnark.KzgMultiProofGnarkBackend{}
