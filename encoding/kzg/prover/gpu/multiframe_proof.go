package gpu

import (
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/Layr-Labs/eigenda/encoding/fft"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/Layr-Labs/eigenda/encoding/utils/gpu_utils"
	"github.com/Layr-Labs/eigenda/encoding/utils/toeplitz"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/ingonyama-zk/icicle/v3/wrappers/golang/core"
	icicle_bn254 "github.com/ingonyama-zk/icicle/v3/wrappers/golang/curves/bn254"
	"github.com/ingonyama-zk/icicle/v3/wrappers/golang/runtime"
)

type KzgGpuProofDevice struct {
	*kzg.KzgConfig
	Fs             *fft.FFTSettings
	FlatFFTPointsT []icicle_bn254.Affine
	SRSIcicle      []icicle_bn254.Affine
	SFs            *fft.FFTSettings
	Srs            *kzg.SRS
	NttCfg         core.NTTConfig[[icicle_bn254.SCALAR_LIMBS]uint32]
	MsmCfg         core.MSMConfig
	Device         runtime.Device
}

type WorkerResult struct {
	err error
}

// This function supports batching over multiple blobs.
// All blobs must have same size and concatenated passed as polyFr
func (p *KzgGpuProofDevice) ComputeMultiFrameProof(polyFr []fr.Element, numChunks, chunkLen, numWorker uint64) ([]bn254.G1Affine, error) {
	begin := time.Now()

	// Robert: Standardizing this to use the same math used in precomputeSRS
	dimE := numChunks
	l := chunkLen
	numPoly := uint64(len(polyFr)) / dimE / chunkLen

	jobChan := make(chan uint64, numWorker)
	results := make(chan WorkerResult, numWorker)

	// create storage for intermediate fft outputs
	coeffStore := make([][]fr.Element, dimE*2)
	for i := range coeffStore {
		coeffStore[i] = make([]fr.Element, l)
	}

	for w := uint64(0); w < numWorker; w++ {
		go p.proofWorker(polyFr, jobChan, l, dimE, coeffStore, results)
	}

	for j := uint64(0); j < l; j++ {
		jobChan <- j
	}
	close(jobChan)

	// return last error
	var err error
	for w := uint64(0); w < numWorker; w++ {
		wr := <-results
		if wr.err != nil {
			err = wr.err
		}
	}

	if err != nil {
		return nil, fmt.Errorf("proof worker error: %v", err)
	}

	preprocessDone := time.Now()

	flattenCoeffStoreFr := make([]fr.Element, 0)
	for i := 0; i < len(coeffStore); i++ {
		flattenCoeffStoreFr = append(flattenCoeffStoreFr, coeffStore[i]...)
	}

	flattenCoeffStoreSf := gpu_utils.ConvertFrToScalarFieldsBytes(flattenCoeffStoreFr)
	flattenCoeffStoreCopy := core.HostSliceFromElements[icicle_bn254.ScalarField](flattenCoeffStoreSf)

	// Start using GPU
	var gpuFFTBatch []bn254.G1Affine
	var gpuErr error
	wg := sync.WaitGroup{}
	wg.Add(1)

	var msmDone, firstECNttDone, secondECNttDone time.Time
	runtime.RunOnDevice(&p.Device, func(args ...any) {
		defer wg.Done()

		// Copy the flatten coeff to device
		var flattenStoreCopyToDevice core.DeviceSlice
		flattenCoeffStoreCopy.CopyToDevice(&flattenStoreCopyToDevice, true)

		sumVec, err := p.MsmBatchOnDevice(flattenStoreCopyToDevice, p.FlatFFTPointsT, int(numPoly)*int(dimE)*2)
		if err != nil {
			gpuErr = fmt.Errorf("msm error: %w", err)
		}

		// Free the flatten coeff store
		flattenStoreCopyToDevice.Free()

		msmDone = time.Now()

		// Compute the first ecntt, and set new batch size for ntt
		p.NttCfg.BatchSize = int32(numPoly)
		sumVecInv, err := p.ECNttOnDevice(sumVec, true, int(dimE)*2*int(numPoly))
		if err != nil {
			gpuErr = fmt.Errorf("first ECNtt error: %w", err)
		}

		sumVec.Free()

		firstECNttDone = time.Now()

		prunedSumVecInv := sumVecInv.Range(0, int(dimE), false)

		// Compute the second ecntt on the reduced size array
		flatProofsBatch, err := p.ECNttToGnarkOnDevice(prunedSumVecInv, false, int(numPoly)*int(dimE))
		if err != nil {
			gpuErr = fmt.Errorf("second ECNtt error: %w", err)
		}

		prunedSumVecInv.Free()

		secondECNttDone = time.Now()

		flatProofsBatchHost := make(core.HostSlice[icicle_bn254.Projective], int(numPoly)*int(dimE))
		flatProofsBatchHost.CopyFromDevice(&flatProofsBatch)
		flatProofsBatch.Free()
		gpuFFTBatch = gpu_utils.HostSliceIcicleProjectiveToGnarkAffine(flatProofsBatchHost, int(p.NumWorker))
	})

	wg.Wait()

	if gpuErr != nil {
		return nil, gpuErr
	}

	end := time.Now()

	slog.Info("Multiproof Time Decomp",
		"total", end.Sub(begin),
		"preproc", preprocessDone.Sub(begin),
		"msm", msmDone.Sub(preprocessDone),
		"fft1", firstECNttDone.Sub(msmDone),
		"fft2", secondECNttDone.Sub(firstECNttDone),
	)

	return gpuFFTBatch, nil
}

func (p *KzgGpuProofDevice) proofWorker(
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
func (p *KzgGpuProofDevice) GetSlicesCoeff(polyFr []fr.Element, dimE, j, l uint64) ([]fr.Element, error) {
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
