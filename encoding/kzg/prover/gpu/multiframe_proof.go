//go:build gpu
// +build gpu

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
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/ingonyama-zk/icicle/v3/wrappers/golang/core"
	bn254_icicle "github.com/ingonyama-zk/icicle/v3/wrappers/golang/curves/bn254"
	icicle_bn254 "github.com/ingonyama-zk/icicle/v3/wrappers/golang/curves/bn254"
	bn254_icicle_g2 "github.com/ingonyama-zk/icicle/v3/wrappers/golang/curves/bn254/g2"
	"github.com/ingonyama-zk/icicle/v3/wrappers/golang/runtime"
)

type WorkerResult struct {
	points []bn254.G1Affine
	err    error
}

type KzgGpuProofDevice struct {
	*kzg.KzgConfig
	Fs             *fft.FFTSettings
	FlatFFTPointsT []bn254_icicle.Affine
	SRSIcicle      []bn254_icicle.Affine
	SFs            *fft.FFTSettings
	Srs            *kzg.SRS
	G2Trailing     []bn254.G2Affine
	NttCfg         core.NTTConfig[[bn254_icicle.SCALAR_LIMBS]uint32]
	MsmCfg         core.MSMConfig
	MsmCfgG2       core.MSMConfig
	GpuLock        *sync.Mutex // lock whenever gpu is needed,
	HeadsG2        []bn254_icicle_g2.G2Affine
	TrailsG2       []bn254_icicle_g2.G2Affine
	Stream         *runtime.Stream
	Device         runtime.Device
}

func (p *KzgGpuProofDevice) ComputeLengthProof(coeffs []fr.Element) (*bn254.G2Affine, error) {
	inputLength := uint64(len(coeffs))
	shiftedSecret := p.G2Trailing[p.KzgConfig.SRSNumberToLoad-inputLength:]
	config := ecc.MultiExpConfig{}
	//The proof of low degree is commitment of the polynomial shifted to the largest srs degree
	var lengthProof bn254.G2Affine
	_, err := lengthProof.MultiExp(shiftedSecret, coeffs, config)
	if err != nil {
		return nil, err
	}
	return &lengthProof, nil
}

func (p *KzgGpuProofDevice) ComputeCommitment(coeffs []fr.Element) (*bn254.G1Affine, error) {
	// compute commit for the full poly
	config := ecc.MultiExpConfig{}
	var commitment bn254.G1Affine
	_, err := commitment.MultiExp(p.Srs.G1[:len(coeffs)], coeffs, config)
	if err != nil {
		return nil, err
	}
	return &commitment, nil
}

func (p *KzgGpuProofDevice) ComputeLengthCommitment(coeffs []fr.Element) (*bn254.G2Affine, error) {
	config := ecc.MultiExpConfig{}

	var lengthCommitment bn254.G2Affine
	_, err := lengthCommitment.MultiExp(p.Srs.G2[:len(coeffs)], coeffs, config)
	if err != nil {
		return nil, err
	}
	return &lengthCommitment, nil
}

// This function supports batching over multiple blobs.
// All blobs must have same size and concatenated passed as polyFr
func (p *KzgGpuProofDevice) ComputeMultiFrameProof(polyFr []fr.Element, numChunks, chunkLen, numWorker uint64) ([]bn254.G1Affine, error) {
	slog.Debug("Starting ComputeMultiFrameProof")

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
	flattenCoeffStoreCopy := core.HostSliceFromElements[bn254_icicle.ScalarField](flattenCoeffStoreSf)

	// Start using GPU
	var gpuFFTBatch []bn254.G1Affine
	var gpuErr error
	wg := sync.WaitGroup{}
	wg.Add(1)

	p.GpuLock.Lock()
	defer p.GpuLock.Unlock()
	var msmDone, firstECNttDone, secondECNttDone time.Time
	runtime.RunOnDevice(&p.Device, func(args ...any) {
		defer wg.Done()

		// Copy the flatten coeff to device
		var flattenStoreCopyToDevice core.DeviceSlice
		flattenCoeffStoreCopy.CopyToDevice(&flattenStoreCopyToDevice, true)

		sumVec, err := p.MsmBatchOnDevice(flattenStoreCopyToDevice, p.FlatFFTPointsT, int(numPoly)*int(dimE)*2, p.Device)
		if err != nil {
			gpuErr = fmt.Errorf("msm error: %w", err)
		}
		msmDone = time.Now()

		// Free the flatten coeff store
		flattenStoreCopyToDevice.Free()

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

		secondECNttDone = time.Now()

		prunedSumVecInv.Free()

		flatProofsBatchHost := make(core.HostSlice[icicle_bn254.Projective], int(numPoly)*int(dimE))
		flatProofsBatchHost.CopyFromDeviceAsync(&flatProofsBatch, *p.Stream)
		flatProofsBatch.Free()
		gpuFFTBatch = gpu_utils.HostSliceIcicleProjectiveToGnarkAffine(flatProofsBatchHost, int(p.NumWorker))
	})

	wg.Wait()

	if gpuErr != nil {
		return nil, gpuErr
	}

	slog.Info("Multiproof Time Decomp",
		"total", secondECNttDone.Sub(begin),
		"preproc", preprocessDone.Sub(begin),
		"msm", msmDone.Sub(preprocessDone),
		"fft1", firstECNttDone.Sub(msmDone),
		"fft2", secondECNttDone.Sub(firstECNttDone),
	)

	return gpuFFTBatch, nil
}

func (p *KzgGpuProofDevice) proofWorkerGPU(
	polyFr []fr.Element,
	jobChan <-chan uint64,
	l uint64,
	dimE uint64,
	coeffStore [][]fr.Element,
	results chan<- WorkerResult,
) {

	for j := range jobChan {
		coeffs, err := p.GetSlicesCoeffWithoutFFT(polyFr, dimE, j, l)

		if err != nil {
			results <- WorkerResult{
				points: nil,
				err:    err,
			}
		} else {
			coeffStore[j] = coeffs
		}
	}

	results <- WorkerResult{
		err: nil,
	}
}

// capable of batching blobs
func (p *KzgGpuProofDevice) GetSlicesCoeffWithoutFFT(polyFr []fr.Element, dimE, j, l uint64) ([]fr.Element, error) {
	// there is a constant term
	m := uint64(dimE*l) - 1
	dim := (m - j%l) / l
	k := j % l
	q := j / l

	toeV := make([]fr.Element, 2*dimE-1)
	for i := uint64(0); i < dim; i++ {
		toeV[i].Set(&polyFr[m+dimE*l*q-(k+i*l)])
	}

	// use precompute table
	tm, err := toeplitz.NewToeplitz(toeV, p.SFs)
	if err != nil {
		return nil, err
	}
	return tm.GetCoeff()
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
