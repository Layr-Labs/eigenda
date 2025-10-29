//go:build !icicle

package icicle

import (
	"context"
	"errors"

	"github.com/Layr-Labs/eigenda/encoding/v2/fft"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

// KzgMultiProofBackend cannot be constructed without icicle build tag.
// We still define the struct and methods to satisfy the interface,
// just to make it clear that this backend could exist but is not available in this build.
type KzgMultiProofBackend struct{}

func (*KzgMultiProofBackend) ComputeMultiFrameProofV2(
	_ context.Context, blobFr []fr.Element, numChunks, chunkLen, numWorker uint64,
) ([]bn254.G1Affine, error) {
	// Not supported
	return nil, errors.New("icicle backend called without icicle build tag")
}

func NewMultiProofBackend(logger logging.Logger,
	fs *fft.FFTSettings, fftPointsT [][]bn254.G1Affine, g1SRS []bn254.G1Affine,
	gpuEnabled bool, numWorker uint64, gpuConcurrentTasks int64,
) (*KzgMultiProofBackend, error) {
	// Not supported
	return nil, errors.New("icicle backend called without icicle build tag")
}
