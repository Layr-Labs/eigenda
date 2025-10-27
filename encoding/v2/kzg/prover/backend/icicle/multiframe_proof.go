//go:build icicle

package icicle

import (
	"fmt"
	"slices"
	"sync"
	"time"

	"github.com/Layr-Labs/eigenda/encoding/icicle"
	"github.com/Layr-Labs/eigenda/encoding/v2/fft"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"

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
	Logger         logging.Logger
	Fs             *fft.FFTSettings
	FlatFFTPointsT []iciclebn254.Affine
	NttCfg         core.NTTConfig[[iciclebn254.SCALAR_LIMBS]uint32]
	MsmCfg         core.MSMConfig
	Device         runtime.Device
	GpuLock        sync.Mutex
	NumWorker      uint64
}

func NewMultiProofBackend(logger logging.Logger,
	fs *fft.FFTSettings, fftPointsT [][]bn254.G1Affine, g1SRS []bn254.G1Affine,
	gpuEnabled bool, numWorker uint64,
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
		MsmCfg:         icicleDevice.MsmCfg,
		Device:         icicleDevice.Device,
		GpuLock:        sync.Mutex{},
		NumWorker:      numWorker,
	}, nil
}

type WorkerResult struct {
	err error
}

// This function supports batching over multiple blobs.
// All blobs must have same size and concatenated passed as polyFr
func (p *KzgMultiProofBackend) ComputeMultiFrameProofV2(polyFr []fr.Element, numChunks, chunkLen, numWorker uint64) ([]bn254.G1Affine, error) {
	begin := time.Now()

	l := chunkLen

	toeplitzMatrixLen := uint64(len(polyFr)) / chunkLen

	// Pre-processing stage - CPU computations
	flattenCoeffStoreFr, err := p.computeCoeffStore(polyFr, numWorker, l, toeplitzMatrixLen)
	if err != nil {
		return nil, fmt.Errorf("coefficient computation error: %v", err)
	}
	preprocessDone := time.Now()

	flattenCoeffStoreSf := icicle.ConvertFrToScalarFieldsBytes(flattenCoeffStoreFr)
	flattenCoeffStoreCopy := core.HostSliceFromElements[iciclebn254.ScalarField](flattenCoeffStoreSf)

	var proofs []bn254.G1Affine
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

		sumVec, err := p.msmBatchOnDevice(flattenStoreCopyToDevice, p.FlatFFTPointsT, int(toeplitzMatrixLen)*2)
		if err != nil {
			icicleErr = fmt.Errorf("msm error: %w", err)
			return
		}

		// Free the flatten coeff store
		flattenStoreCopyToDevice.Free()

		msmDone = time.Now()

		// Compute the first ecntt, and set new batch size for ntt
		p.NttCfg.BatchSize = int32(1)

		// run two ecntt in one function, because the second ecntt needs to be larger
		// size, but icicle does not offer device to device copy, so we have to use
		// Range trick. Hence we combine two ecntt into one function allowing us to
		// manage it better
		proofs, firstECNttDone, err = p.twoEcnttOnDevice(sumVec, int(numChunks), int(toeplitzMatrixLen))
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

// MsmBatchOnDevice function supports batch across blobs.
// totalSize is the number of output points, which equals to numPoly * 2 * dimE , dimE is number of chunks
func (c *KzgMultiProofBackend) msmBatchOnDevice(rowsFrIcicleCopy core.DeviceSlice, rowsG1Icicle []iciclebn254.Affine, totalSize int) (core.DeviceSlice, error) {
	rowsG1IcicleCopy := core.HostSliceFromElements[iciclebn254.Affine](rowsG1Icicle)

	var p iciclebn254.Projective
	var out core.DeviceSlice

	_, err := out.Malloc(p.Size(), totalSize)
	if err != runtime.Success {
		return out, fmt.Errorf("allocating bytes on device failed: %v", err.AsString())
	}

	err = msm.Msm(rowsFrIcicleCopy, rowsG1IcicleCopy, &c.MsmCfg, out)
	if err != runtime.Success {
		return out, fmt.Errorf("msm error: %v", err.AsString())
	}

	return out, nil
}

func (c *KzgMultiProofBackend) twoEcnttOnDevice(batchPoints core.DeviceSlice, numChunks int, toeplitzMatrixLen int) ([]bn254.G1Affine, time.Time, error) {
	var p iciclebn254.Projective
	var deviceWithNumChunkElement core.DeviceSlice
	var output core.DeviceSlice

	_, err := deviceWithNumChunkElement.Malloc(p.Size(), numChunks)
	if err != runtime.Success {
		return nil, time.Time{}, fmt.Errorf("allocating bytes on device failed: %v", err.AsString())
	}

	// the size is twice because of the FFT trick on toeplitz matrix
	firstECNTTLen := toeplitzMatrixLen * 2
	// Device memory for first ecntt
	firstECNTTDeviceSlice := deviceWithNumChunkElement.RangeTo(firstECNTTLen, false)
	err = ecntt.ECNtt(batchPoints, core.KInverse, &c.NttCfg, firstECNTTDeviceSlice)
	if err != runtime.Success {
		return nil, time.Time{}, fmt.Errorf("inverse ecntt failed: %v", err.AsString())
	}

	// now only keep the toeplitzMatrixLen elements as they are, set the rest to zero.
	// Zeros are the infinity points for G1Affine, see
	// https://github.com/ingonyama-zk/icicle/blob/f797dae1fbaa30cf1a7133f89f37d333adb590f5/wrappers/golang/curves/bn254/curve.go#L116
	// unit in the Range function is measured by element size
	infinityPoints := deviceWithNumChunkElement.Range(toeplitzMatrixLen, numChunks, false)
	err = runtime.MemSet(
		infinityPoints.AsUnsafePointer(),
		0,
		uint(infinityPoints.Cap()),
	)
	if err != runtime.Success {
		return nil, time.Time{}, fmt.Errorf("setting zeros on inverse ecntt result failed: %v", err.AsString())
	}

	// take the second ecntt
	_, err = deviceWithNumChunkElement.Malloc(p.Size(), numChunks)
	err = ecntt.ECNtt(deviceWithNumChunkElement, core.KForward, &c.NttCfg, output)
	if err != runtime.Success {
		return nil, time.Time{}, fmt.Errorf("forward ecntt failed: %v", err.AsString())
	}

	// free intermediate GPU memory
	deviceWithNumChunkElement.Free()

	proofsBatchHost := make(core.HostSlice[iciclebn254.Projective], numChunks)
	proofsBatchHost.CopyFromDevice(&output)
	// free GPU memory for proof
	output.Free()

	proofs := icicle.HostSliceIcicleProjectiveToGnarkAffine(proofsBatchHost, int(c.NumWorker))

	return proofs, time.Time{}, nil
}
