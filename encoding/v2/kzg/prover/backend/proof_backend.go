// Note that all functions are suffixed with V2 to avoid passing a V1 backend to a V2 prover.
// The main difference is that the V2 prover requires blobs to be of power-of-two size.
package backend

import (
	"context"

	"github.com/Layr-Labs/eigenda/encoding/v2/kzg/prover/backend/gnark"
	"github.com/Layr-Labs/eigenda/encoding/v2/kzg/prover/backend/icicle"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

// Proof device represents a backend capable of computing KZG multiproofs.
type KzgMultiProofsBackendV2 interface {
	// the length of blobFr must be power of 2
	ComputeMultiFrameProofV2(
		ctx context.Context, blobFr []fr.Element, numChunks, chunkLen, numWorker uint64,
	) ([]bn254.G1Affine, error)
}

// We implement two backends: gnark and icicle.
//   - Gnark uses the gnark library and is the default CPU-based backend, and is always available.
//   - Icicle uses the icicle library and can leverage GPU acceleration, but requires building with the icicle tag.
//     Building with the icicle tag will inject the dynamic libraries required to use icicle.
//
// Both backends implement a NewMultiProofBackend constructor, which in the case of icicle
// will return an error if the icicle build tag was not used.
var _ KzgMultiProofsBackendV2 = &gnark.KzgMultiProofBackend{}
var _ KzgMultiProofsBackendV2 = &icicle.KzgMultiProofBackend{}
