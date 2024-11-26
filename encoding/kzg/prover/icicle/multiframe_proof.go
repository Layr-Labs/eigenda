//go:build icicle

package icicle

import (
	"fmt"
	"log/slog"
	"runtime"
	"time"

	"github.com/Layr-Labs/eigenda/encoding/fft"
	"github.com/Layr-Labs/eigenda/encoding/icicle"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/Layr-Labs/eigenda/encoding/utils/toeplitz"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/ingonyama-zk/icicle/v3/wrappers/golang/core"
	iciclebn254 "github.com/ingonyama-zk/icicle/v3/wrappers/golang/curves/bn254"
	icicleRuntime "github.com/ingonyama-zk/icicle/v3/wrappers/golang/runtime"
)

type KzgMultiProofIcicleBackend struct {
	*kzg.KzgConfig
	Fs             *fft.FFTSettings
	FlatFFTPointsT []iciclebn254.Affine
	SRSIcicle      []iciclebn254.Affine
	SFs            *fft.FFTSettings
	Srs            *kzg.SRS
	NttCfg         core.NTTConfig[[iciclebn254.SCALAR_LIMBS]uint32]
	MsmCfg         core.MSMConfig
	Device         icicleRuntime.Device
}

type WorkerResult struct {
	err error
}

// This function supports batching over multiple blobs.
// All blobs must have same size and concatenated passed as polyFr
func (p *KzgMultiProofIcicleBackend) ComputeMultiFrameProof(polyFr []fr.Element, numChunks, chunkLen, numWorker uint64) ([]bn254.G1Affine, error) {
	begin := time.Now()

	dimE := numChunks
	l := chunkLen
	numPoly := uint64(len(polyFr)) / dimE / chunkLen

	// Pre-processing stage - CPU computations
	coeffStore, err := p.computeCoeffStore(polyFr, numWorker, l, dimE)
	if err != nil {
		return nil, fmt.Errorf("coefficient computation error: %v", err)
	}
	preprocessDone := time.Now()

	// Prepare data before GPU operations
	flattenCoeffStoreFr := make([]fr.Element, 0, len(coeffStore)*len(coeffStore[0]))
	for i := 0; i < len(coeffStore); i++ {
		flattenCoeffStoreFr = append(flattenCoeffStoreFr, coeffStore[i]...)
	}

	flattenCoeffStoreSf := icicle.ConvertFrToScalarFieldsBytes(flattenCoeffStoreFr)
	flattenCoeffStoreCopy := core.HostSliceFromElements[iciclebn254.ScalarField](flattenCoeffStoreSf)

	// Lock the OS thread for GPU operations
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	// Set device
	originalDevice, icicleErr := icicleRuntime.GetActiveDevice()
	if icicleErr != icicleRuntime.Success {
		return nil, fmt.Errorf("failed to get active device: %v", err)
	}
	defer icicleRuntime.SetDevice(originalDevice)

	if err := icicleRuntime.SetDevice(&p.Device); err != icicleRuntime.Success {
		return nil, fmt.Errorf("failed to set device: %v", err)
	}

	// Copy data to device
	var flattenStoreCopyToDevice core.DeviceSlice
	flattenCoeffStoreCopy.CopyToDevice(&flattenStoreCopyToDevice, true)
	defer flattenStoreCopyToDevice.Free()

	// 1. MSM Batch Operation
	sumVec, err := p.MsmBatchOnDevice(flattenStoreCopyToDevice, p.FlatFFTPointsT, int(numPoly)*int(dimE)*2)
	if err != nil {
		return nil, fmt.Errorf("msm error: %v", err)
	}
	defer sumVec.Free()
	msmDone := time.Now()

	// 2. First ECNtt
	p.NttCfg.BatchSize = int32(numPoly)
	sumVecInv, err := p.ECNttOnDevice(sumVec, true, int(dimE)*2*int(numPoly))
	if err != nil {
		return nil, fmt.Errorf("first ECNtt error: %v", err)
	}
	defer sumVecInv.Free()
	fft1Done := time.Now()

	// 3. Prune and Second ECNtt
	prunedSumVecInv := sumVecInv.Range(0, int(dimE), false)
	flatProofsBatch, err := p.ECNttToGnarkOnDevice(prunedSumVecInv, false, int(numPoly)*int(dimE))
	if err != nil {
		return nil, fmt.Errorf("second ECNtt error: %v", err)
	}
	defer flatProofsBatch.Free()
	fft2Done := time.Now()

	// 4. Copy results back to host
	flatProofsBatchHost := make(core.HostSlice[iciclebn254.Projective], int(numPoly)*int(dimE))
	flatProofsBatchHost.CopyFromDevice(&flatProofsBatch)

	// Convert to final format
	icicleFFTBatch := icicle.HostSliceIcicleProjectiveToGnarkAffine(flatProofsBatchHost, int(p.NumWorker))

	slog.Info("Multiproof Time Decomp",
		"preproc", preprocessDone.Sub(begin),
		"msm", msmDone.Sub(preprocessDone),
		"fft1", fft1Done.Sub(msmDone),
		"fft2", fft2Done.Sub(fft1Done),
	)

	return icicleFFTBatch, nil
}

// Helper function to handle coefficient computation
func (p *KzgMultiProofIcicleBackend) computeCoeffStore(polyFr []fr.Element, numWorker, l, dimE uint64) ([][]fr.Element, error) {
	jobChan := make(chan uint64, numWorker)
	results := make(chan WorkerResult, numWorker)

	coeffStore := make([][]fr.Element, dimE*2)
	for i := range coeffStore {
		coeffStore[i] = make([]fr.Element, l)
	}

	// Start workers
	for w := uint64(0); w < numWorker; w++ {
		go p.proofWorker(polyFr, jobChan, l, dimE, coeffStore, results)
	}

	// Send jobs
	for j := uint64(0); j < l; j++ {
		jobChan <- j
	}
	close(jobChan)

	// Collect results
	var lastErr error
	for w := uint64(0); w < numWorker; w++ {
		if wr := <-results; wr.err != nil {
			lastErr = wr.err
		}
	}

	if lastErr != nil {
		return nil, fmt.Errorf("proof worker error: %v", lastErr)
	}

	return coeffStore, nil
}

func (p *KzgMultiProofIcicleBackend) proofWorker(
	polyFr []fr.Element,
	jobChan <-chan uint64,
	l uint64,
	dimE uint64,
	coeffStore [][]fr.Element,
	results chan<- WorkerResult,
) {

	for j := range jobChan {
		coeffs, err := p.GetSlicesCoeff(polyFr, dimE, j, l)
		if err != nil {
			results <- WorkerResult{
				err: err,
			}
		} else {
			for i := 0; i < len(coeffs); i++ {
				coeffStore[i][j] = coeffs[i]
			}
		}
	}

	results <- WorkerResult{
		err: nil,
	}
}

// output is in the form see primeField toeplitz
//
// phi ^ (coset size ) = 1
//
// implicitly pad slices to power of 2
func (p *KzgMultiProofIcicleBackend) GetSlicesCoeff(polyFr []fr.Element, dimE, j, l uint64) ([]fr.Element, error) {
	// there is a constant term
	m := uint64(len(polyFr)) - 1
	dim := (m - j) / l

	// maximal number of unique values from a toeplitz matrix
	tDim := 2*dimE - 1

	toeV := make([]fr.Element, tDim)
	for i := uint64(0); i < dim; i++ {

		toeV[i].Set(&polyFr[m-(j+i*l)])
	}

	// use precompute table
	tm, err := toeplitz.NewToeplitz(toeV, p.SFs)
	if err != nil {
		return nil, err
	}
	return tm.GetFFTCoeff()
}
