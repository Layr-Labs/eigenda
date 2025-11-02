//go:build icicle

package icicle

import (
	"context"
	"fmt"
	"slices"
	"sync"
	"time"

	_ "github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/icicle"
	"github.com/Layr-Labs/eigenda/encoding/v2/fft"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"golang.org/x/sync/semaphore"

	"github.com/ingonyama-zk/icicle/v3/wrappers/golang/core"
	iciclebn254 "github.com/ingonyama-zk/icicle/v3/wrappers/golang/curves/bn254"
	"github.com/ingonyama-zk/icicle/v3/wrappers/golang/curves/bn254/ecntt"
	"github.com/ingonyama-zk/icicle/v3/wrappers/golang/curves/bn254/msm"
	"github.com/ingonyama-zk/icicle/v3/wrappers/golang/runtime"
)

const (
	// MAX_NTT_SIZE is the maximum NTT domain size needed to compute FFTs for the
	// largest supported blobs. Assuming a coding ratio of 1/8 and symbol size of 32 bytes:
	// - Encoded size: 2^{MAX_NTT_SIZE} * 32 bytes ≈ 1 GB
	// - Original blob size: 2^{MAX_NTT_SIZE} * 32 / 8 = 2^{MAX_NTT_SIZE + 2} ≈ 128 MB
	MAX_NTT_SIZE = 25
)

type KzgMultiProofBackend struct {
	Logger logging.Logger
	Fs     *fft.FFTSettings
	// TODO(samlaf): we should send the srs table points to the device once in the constructor
	// and keep a deviceSlice pointer to it. This would require a destructor to free the device memory.
	// Also need to account how much memory this would use over all parametrized provers.
	FlatFFTPointsT []iciclebn254.Affine
	NttCfg         core.NTTConfig[[iciclebn254.SCALAR_LIMBS]uint32]
	Device         runtime.Device
	NumWorker      uint64
	// request-weighted semaphore.
	// See [encoding.Config.GPUConcurrentFrameGenerationDangerous] for more details.
	GpuSemaphore *semaphore.Weighted
}

func NewMultiProofBackend(logger logging.Logger,
	fs *fft.FFTSettings, fftPointsT [][]bn254.G1Affine, g1SRS []bn254.G1Affine,
	gpuEnabled bool, numWorker uint64, gpuConcurrentProofs int64,
) (*KzgMultiProofBackend, error) {
	icicleDevice, err := icicle.NewIcicleDevice(icicle.IcicleDeviceConfig{
		Logger:     logger,
		GPUEnable:  gpuEnabled,
		NTTSize:    MAX_NTT_SIZE,
		FFTPointsT: fftPointsT,
		SRSG1:      g1SRS,
	})
	if err != nil {
		return nil, fmt.Errorf("configure icicle device: %w", err)
	}

	// Set up icicle multiproof backend
	return &KzgMultiProofBackend{
		Logger:         logger,
		Fs:             fs,
		FlatFFTPointsT: icicleDevice.FlatFFTPointsT,
		NttCfg:         icicleDevice.NttCfg,
		Device:         icicleDevice.Device,
		GpuSemaphore:   semaphore.NewWeighted(gpuConcurrentProofs),
		NumWorker:      numWorker,
	}, nil
}

type WorkerResult struct {
	err error
}

// This function supports batching over multiple blobs.
// All blobs must have same size and concatenated passed as polyFr
func (p *KzgMultiProofBackend) ComputeMultiFrameProofV2(ctx context.Context, polyFr []fr.Element, numChunks, chunkLen, numWorker uint64) ([]bn254.G1Affine, error) {
	begin := time.Now()

	toeplitzMatrixLen := uint64(len(polyFr)) / chunkLen

	l := chunkLen

	// Pre-processing stage - CPU computations
	flattenCoeffStoreFr, err := p.computeCoeffStore(polyFr, numWorker, l, toeplitzMatrixLen)
	if err != nil {
		return nil, fmt.Errorf("coefficient computation error: %v", err)
	}
	preprocessDone := time.Now()

	var proofs []bn254.G1Affine
	var icicleErr error

	// We acquire a semaphore here to avoid too many concurrent GPU requests,
	// each of which does 1 MSM + 2 NTTs. This is a very unideal and coarse grain solution, but unfortunately
	// icicle doesn't have nice backpressure, and the GPU kernel just panics if RAM is exhausted.
	// We could use a finer-grained semaphore that calculates the RAM usage per request,
	// but we'd have to hardcode some approximation of the RAM usage per MSM/NTT, which feels
	// very hardcoded and hardware dependent. For now opting to keep this simple.
	// TODO(samlaf): rethink this approach.
	p.GpuSemaphore.Acquire(ctx, 1)
	defer p.GpuSemaphore.Release(1)

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

		var projectivePoint iciclebn254.Projective
		var sumVec core.DeviceSlice

		_, mallocErr := sumVec.Malloc(projectivePoint.Size(), int(toeplitzMatrixLen)*2)
		if mallocErr != runtime.Success {
			icicleErr = fmt.Errorf("allocating bytes on device failed: %v", mallocErr.AsString())
			return
		}
		defer sumVec.Free()

		// The msm is computed synchronously (default config value is async=false).
		// We could possibly share the same stream as the ecntt use (c.NttCfg.StreamHandle).
		// TODO(samlaf): rethink how we use streams and async computations in general.
		msmCfg := msm.GetDefaultMSMConfig()
		msmCfg.AreScalarsMontgomeryForm = true
		frsHostOrDeviceSlice := core.HostSliceFromElements(flattenCoeffStoreFr)
		// TODO(samlaf): we could send the srs table points to the device once in the constructor
		// and keep a deviceSlice pointer to it.
		g1PointsHostSlice := core.HostSliceFromElements(p.FlatFFTPointsT)
		msmErr := msm.Msm(frsHostOrDeviceSlice, g1PointsHostSlice, &msmCfg, sumVec)
		if msmErr != runtime.Success {
			icicleErr = fmt.Errorf("msm error: %v", msmErr.AsString())
			return
		}

		msmDone = time.Now()

		// Compute the first ecntt, and set new batch size for ntt
		p.NttCfg.BatchSize = int32(1)
		// run two ecntt in one function, the first and second ecntt operates on the same device slice
		proofs, firstECNttDone, err = p.twoEcnttOnDevice(sumVec, int(numChunks), int(toeplitzMatrixLen))
		if err != nil {
			icicleErr = err
			return
		}
	})

	wg.Wait()

	if icicleErr != nil {
		return nil, icicleErr
	}

	end := time.Now()

	p.Logger.Info("Multiproof Time Decomp",
		"total", end.Sub(begin),
		"preproc", preprocessDone.Sub(begin),
		"msm", msmDone.Sub(preprocessDone),
		"fft1", firstECNttDone.Sub(msmDone),
		"fft2", secondECNttDone.Sub(firstECNttDone),
	)

	return proofs, nil
}

// Modify the function signature to return a flat array
func (p *KzgMultiProofBackend) computeCoeffStore(polyFr []fr.Element, numWorker, l, dimE uint64) ([]fr.Element, error) {
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
func (p *KzgMultiProofBackend) proofWorker(
	polyFr []fr.Element,
	jobChan <-chan uint64,
	l uint64,
	dimE uint64,
	coeffStore []fr.Element,
	results chan<- WorkerResult,
) {
	for j := range jobChan {
		coeffs, err := p.getSlicesCoeff(polyFr, dimE, j, l)
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

// getSlicesCoeff computes step 2 of the FFT trick for computing h,
// in proposition 2 of https://eprint.iacr.org/2023/033.pdf.
// However, given that it's used in the multiple multiproofs scenario,
// the indices used are more complex (eg. (m-j)/l below).
// Those indices are from the matrix in section 3.1.1 of
// https://github.com/khovratovich/Kate/blob/master/Kate_amortized.pdf
// Returned slice has len [2*dimE].
//
// TODO(samlaf): better document/explain/refactor/rename this function,
// to explain how it fits into the overall scheme.
func (p *KzgMultiProofBackend) getSlicesCoeff(polyFr []fr.Element, dimE, j, l uint64) ([]fr.Element, error) {
	toeplitzExtendedVec := make([]fr.Element, 2*dimE)

	m := uint64(len(polyFr)) - 1 // there is a constant term
	dim := (m - j) / l
	for i := range dim {
		toeplitzExtendedVec[i].Set(&polyFr[m-(j+i*l)])
	}
	// We keep the first element as is, and reverse the rest of the slice.
	// This is classic Toeplitz manipulations, as for example describe in
	// https://alinush.github.io/2020/03/19/multiplying-a-vector-by-a-toeplitz-matrix.html
	slices.Reverse(toeplitzExtendedVec[1:])

	out, err := p.Fs.FFT(toeplitzExtendedVec, false)
	if err != nil {
		return nil, fmt.Errorf("fft: %w", err)
	}
	return out, nil
}

// twoEcnttOnDevice takes the first ecntt to generate the kzg proofs. Only the first half of the result
// are considered kzg proof, and it comes from the Toeplitz trick, readers can refer to
// https://alinush.github.io/2020/03/19/multiplying-a-vector-by-a-toeplitz-matrix.html
// Then the kzg proofs are padded with infinity points to the size of numChunks. And this is the vector
// which the second ecntt is taken.
func (c *KzgMultiProofBackend) twoEcnttOnDevice(
	batchPoints core.DeviceSlice,
	numChunks int,
	toeplitzMatrixLen int,
) ([]bn254.G1Affine, time.Time, error) {
	var p iciclebn254.Projective
	// we only allocate one large gpu memory for all operation, so it has to be large enough to cover all cases
	// including the first and the second ECNTT
	var bufferProjectivePointsOnDevice core.DeviceSlice

	numPointsOnDevice := numChunks

	// the size is twice because of the FFT trick on toeplitz matrix
	firstECNTTLen := toeplitzMatrixLen * 2

	// when first ecntt is larger than numChunk, we must allocate enough memory
	// it happens if numChunks is equal of less than toeplitzMatrixLen
	if numChunks < firstECNTTLen {
		numPointsOnDevice = firstECNTTLen
	}
	_, err := bufferProjectivePointsOnDevice.Malloc(p.Size(), numPointsOnDevice)
	if err != runtime.Success {
		return nil, time.Time{}, fmt.Errorf("allocating bytes on device failed: %v", err.AsString())
	}
	// free intermediate GPU memory
	defer bufferProjectivePointsOnDevice.Free()

	// specify device memory sluce for first ecntt
	firstECNTTDeviceSlice := bufferProjectivePointsOnDevice.RangeTo(firstECNTTLen, false)
	err = ecntt.ECNtt(batchPoints, core.KInverse, &c.NttCfg, firstECNTTDeviceSlice)
	if err != runtime.Success {
		return nil, time.Time{}, fmt.Errorf("inverse ecntt failed: %v", err.AsString())
	}

	proofsBatchHost := make(core.HostSlice[iciclebn254.Projective], numChunks)

	// if numChunk is smaller or equal to toeplitzMatrixLen, there is no need to set points to infinity
	// otherwise set all points to infinity
	if numChunks > toeplitzMatrixLen {
		// now only keep the toeplitzMatrixLen elements as they are, set the rest to zero.
		// Zeros are the infinity points for G1Projective points
		// unit in the Range function is measured by element size
		infinityPointsOnDevice := bufferProjectivePointsOnDevice.Range(toeplitzMatrixLen, numChunks, false)
		infinityProjectivePoints := make([]iciclebn254.Projective, numChunks-toeplitzMatrixLen)
		// explicitly sets all value to zero
		// it does not work if we just initialize it, it is most likely due to all members of the struct
		// would be initialized as 0, however, to have a projective point as inifity, Y needs to be 1
		for i := range infinityProjectivePoints {
			infinityProjectivePoints[i].Zero()
		}
		infinityPointsHost := core.HostSliceFromElements(infinityProjectivePoints)
		// copy to device, but don't allocate memory
		infinityPointsHost.CopyToDevice(&infinityPointsOnDevice, false)
	}

	secondECNTTDeviceSlice := bufferProjectivePointsOnDevice.RangeTo(numChunks, false)

	// take the second ecntt
	err = ecntt.ECNtt(secondECNTTDeviceSlice, core.KForward, &c.NttCfg, proofsBatchHost)
	if err != runtime.Success {
		return nil, time.Time{}, fmt.Errorf("forward ecntt failed: %v", err.AsString())
	}

	proofs := icicle.HostSliceIcicleProjectiveToGnarkAffine(proofsBatchHost, int(c.NumWorker))

	return proofs, time.Time{}, nil
}
