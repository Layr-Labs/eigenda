//go:build icicle

package icicle

import (
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/Layr-Labs/eigenda/encoding/fft"
	"github.com/Layr-Labs/eigenda/encoding/icicle"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/Layr-Labs/eigenda/encoding/utils/toeplitz"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/ingonyama-zk/icicle/v3/wrappers/golang/core"
	iciclebn254 "github.com/ingonyama-zk/icicle/v3/wrappers/golang/curves/bn254"
	"github.com/ingonyama-zk/icicle/v3/wrappers/golang/runtime"
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
	Device         runtime.Device
	GpuLock        sync.Mutex
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
	flattenCoeffStoreFr, err := p.computeCoeffStore(polyFr, numWorker, l, dimE)
	if err != nil {
		return nil, fmt.Errorf("coefficient computation error: %v", err)
	}
	preprocessDone := time.Now()

	flattenCoeffStoreSf := icicle.ConvertFrToScalarFieldsBytes(flattenCoeffStoreFr)
	flattenCoeffStoreCopy := core.HostSliceFromElements[iciclebn254.ScalarField](flattenCoeffStoreSf)

	var icicleFFTBatch []bn254.G1Affine
	var icicleErr error

	// GPU operations
	p.GpuLock.Lock()
	defer p.GpuLock.Unlock()

	wg := sync.WaitGroup{}
	wg.Add(1)

	var msmDone, firstECNttDone, secondECNttDone time.Time
	runtime.RunOnDevice(&p.Device, func(args ...any) {
		defer wg.Done()
		defer func() {
			if r := recover(); r != nil {
				icicleErr = fmt.Errorf("GPU operation panic: %v", r)
			}
		}()

		// Copy the flatten coeff to device
		var flattenStoreCopyToDevice core.DeviceSlice
		flattenCoeffStoreCopy.CopyToDevice(&flattenStoreCopyToDevice, true)

		sumVec, err := p.MsmBatchOnDevice(flattenStoreCopyToDevice, p.FlatFFTPointsT, int(numPoly)*int(dimE)*2)
		if err != nil {
			icicleErr = fmt.Errorf("msm error: %w", err)
			return
		}

		// Free the flatten coeff store
		flattenStoreCopyToDevice.Free()

		msmDone = time.Now()

		// Compute the first ecntt, and set new batch size for ntt
		p.NttCfg.BatchSize = int32(numPoly)
		sumVecInv, err := p.ECNttOnDevice(sumVec, true, int(dimE)*2*int(numPoly))
		if err != nil {
			icicleErr = fmt.Errorf("first ECNtt error: %w", err)
			return
		}

		sumVec.Free()

		firstECNttDone = time.Now()

		prunedSumVecInv := sumVecInv.Range(0, int(dimE), false)

		// Compute the second ecntt on the reduced size array
		flatProofsBatch, err := p.ECNttToGnarkOnDevice(prunedSumVecInv, false, int(numPoly)*int(dimE))
		if err != nil {
			icicleErr = fmt.Errorf("second ECNtt error: %w", err)
			return
		}

		prunedSumVecInv.Free()

		secondECNttDone = time.Now()

		flatProofsBatchHost := make(core.HostSlice[iciclebn254.Projective], int(numPoly)*int(dimE))
		flatProofsBatchHost.CopyFromDevice(&flatProofsBatch)
		flatProofsBatch.Free()
		icicleFFTBatch = icicle.HostSliceIcicleProjectiveToGnarkAffine(flatProofsBatchHost, int(p.NumWorker))
	})

	wg.Wait()

	if icicleErr != nil {
		return nil, icicleErr
	}

	end := time.Now()

	slog.Info("Multiproof Time Decomp",
		"total", end.Sub(begin),
		"preproc", preprocessDone.Sub(begin),
		"msm", msmDone.Sub(preprocessDone),
		"fft1", firstECNttDone.Sub(msmDone),
		"fft2", secondECNttDone.Sub(firstECNttDone),
	)

	return icicleFFTBatch, nil
}

// Modify the function signature to return a flat array
func (p *KzgMultiProofIcicleBackend) computeCoeffStore(polyFr []fr.Element, numWorker, l, dimE uint64) ([]fr.Element, error) {
	totalSize := dimE * 2 * l // Total size of the flattened array
	coeffStore := make([]fr.Element, totalSize)

	jobChan := make(chan uint64, numWorker)
	results := make(chan WorkerResult, numWorker)

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

// Modified worker function to write directly to the flat array
func (p *KzgMultiProofIcicleBackend) proofWorker(
	polyFr []fr.Element,
	jobChan <-chan uint64,
	l uint64,
	dimE uint64,
	coeffStore []fr.Element,
	results chan<- WorkerResult,
) {
	for j := range jobChan {
		coeffs, err := p.GetSlicesCoeff(polyFr, dimE, j, l)
		if err != nil {
			results <- WorkerResult{
				err: err,
			}
			return
		}

		// Write directly to the correct positions in the flat array
		// For each j, we need to write to the corresponding position in each block
		for i := uint64(0); i < dimE*2; i++ {
			coeffStore[i*l+j] = coeffs[i]
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
