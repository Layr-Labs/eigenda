//go:build gpu
// +build gpu

package gpu

import (
	"fmt"
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

	// Start using GPU
	p.GpuLock.Lock()
	defer p.GpuLock.Unlock()

	// Compute NTT on the coeff matrix
	// p.NttCfg.BatchSize = int32(l * numPoly)
	// coeffStoreFft, e := p.NTT(coeffStore)
	// if e != nil {
	// 	return nil, e
	// }
	// nttDone := time.Now()

	// // transpose the FFT tranformed matrix
	// coeffStoreFftTranspose, err := Transpose(coeffStoreFft, int(l), int(numPoly), int(dimE)*2)
	// if err != nil {
	// 	return nil, e
	// }
	// transposingDone := time.Now()
	// numSymbol := len(coeffStore[0])
	// batchSize := len(coeffStore)
	// totalSize := numSymbol * batchSize

	// prepare scaler fields
	flattenCoeffStoreFr := make([]fr.Element, 0)
	for i := 0; i < len(coeffStore); i++ {
		flattenCoeffStoreFr = append(flattenCoeffStoreFr, coeffStore[i]...)
	}

	flattenCoeffStoreSf := gpu_utils.ConvertFrToScalarFieldsBytes(flattenCoeffStoreFr)
	flattenCoeffStoreCopy := core.HostSliceFromElements[bn254_icicle.ScalarField](flattenCoeffStoreSf)

	var flattenStoreCopyToDevice core.DeviceSlice
	flattenCoeffStoreCopy.CopyToDeviceAsync(&flattenStoreCopyToDevice, *p.Stream, true)

	// compute msm on each rows of the transposed matrix
	fmt.Println("numPoly", numPoly)
	fmt.Println("dimE", dimE)
	fmt.Println("row", p.FlatFFTPointsT[0].Size())
	fmt.Println("col", len(p.FlatFFTPointsT))

	sumVec, err := p.MsmBatchOnDevice(flattenStoreCopyToDevice, p.FlatFFTPointsT, int(numPoly)*int(dimE)*2)
	if err != nil {
		return nil, err
	}
	msmDone := time.Now()

	// Free the flatten coeff store
	flattenStoreCopyToDevice.FreeAsync(*p.Stream)

	// compute the first ecntt, and set new batch size for ntt
	p.NttCfg.BatchSize = int32(numPoly)
	sumVecInv, err := p.ECNttOnDevice(sumVec, true, int(dimE)*2*int(numPoly))
	if err != nil {
		return nil, err
	}

	// Free sumVec
	sumVec.FreeAsync(*p.Stream)

	firstECNttDone := time.Now()

	prunedSumVecInv := sumVecInv.Range(0, int(dimE), false)
	// compute the second ecntt on the reduced size array
	flatProofsBatch, err := p.ECNttToGnarkOnDevice(prunedSumVecInv, false, int(numPoly)*int(dimE))
	if err != nil {
		return nil, fmt.Errorf("second ECNtt error: %w", err)
	}

	prunedSumVecInv.FreeAsync(*p.Stream)

	flatProofsBatchHost := make(core.HostSlice[icicle_bn254.Projective], int(numPoly)*int(dimE))
	flatProofsBatchHost.CopyFromDeviceAsync(&flatProofsBatch, *p.Stream)
	flatProofsBatch.FreeAsync(*p.Stream)
	gpuFFTBatch := gpu_utils.HostSliceIcicleProjectiveToGnarkAffine(flatProofsBatchHost, int(p.NumWorker))

	secondECNttDone := time.Now()

	fmt.Printf("Multiproof Time Decomp \n\t\ttotal   %-20s \n\t\tpreproc %-20s \n\t\tmsm     %-20s \n\t\tfft1    %-20s \n\t\tfft2    %-20s\n",
		secondECNttDone.Sub(begin).String(),
		preprocessDone.Sub(begin).String(),
		msmDone.Sub(preprocessDone).String(),
		firstECNttDone.Sub(msmDone).String(),
		secondECNttDone.Sub(firstECNttDone).String(),
	)

	runtime.SynchronizeStream(*p.Stream)

	// only takes the first half
	return gpuFFTBatch, nil
}

// This function supports batching over multiple blobs.
// All blobs must have same size and concatenated passed as polyFr
func (p *KzgGpuProofDevice) ComputeMultiFrameProofHost(polyFr []fr.Element, numChunks, chunkLen, numWorker uint64) ([]bn254.G1Affine, error) {
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

	// Start using GPU
	p.GpuLock.Lock()
	defer p.GpuLock.Unlock()

	// Compute NTT on the coeff matrix
	// p.NttCfg.BatchSize = int32(l * numPoly)
	// coeffStoreFft, e := p.NTT(coeffStore)
	// if e != nil {
	// 	return nil, e
	// }
	// nttDone := time.Now()

	// // transpose the FFT tranformed matrix
	// coeffStoreFftTranspose, err := Transpose(coeffStoreFft, int(l), int(numPoly), int(dimE)*2)
	// if err != nil {
	// 	return nil, e
	// }
	// transposingDone := time.Now()
	// numSymbol := len(coeffStore[0])
	// batchSize := len(coeffStore)
	// totalSize := numSymbol * batchSize

	// prepare scaler fields
	flattenCoeffStoreFr := make([]fr.Element, 0)
	for i := 0; i < len(coeffStore); i++ {
		flattenCoeffStoreFr = append(flattenCoeffStoreFr, coeffStore[i]...)
	}
	flattenCoeffStoreSf := gpu_utils.ConvertFrToScalarFieldsBytes(flattenCoeffStoreFr)
	flattenCoeffStoreCopy := core.HostSliceFromElements[bn254_icicle.ScalarField](flattenCoeffStoreSf)

	// compute msm on each rows of the transposed matrix
	sumVec, err := p.MsmBatch(flattenCoeffStoreCopy, p.FlatFFTPointsT, int(numPoly)*int(dimE)*2)
	if err != nil {
		return nil, err
	}
	msmDone := time.Now()

	// compute the first ecntt, and set new batch size for ntt
	p.NttCfg.BatchSize = int32(numPoly)
	sumVecInv, err := p.ECNtt(sumVec, true, int(dimE)*2*int(numPoly))
	if err != nil {
		return nil, err
	}

	firstECNttDone := time.Now()

	prunedSumVecInv := sumVecInv[:int(dimE)]

	// compute the second ecntt on the reduced size array
	flatProofsBatch, err := p.ECNttToGnark(prunedSumVecInv, false, int(numPoly)*int(dimE))
	if err != nil {
		return nil, fmt.Errorf("second ECNtt error: %w", err)
	}

	secondECNttDone := time.Now()

	fmt.Printf("Multiproof Time Decomp \n\t\ttotal   %-20s \n\t\tpreproc %-20s \n\t\tmsm     %-20s \n\t\tfft1    %-20s \n\t\tfft2    %-20s\n",
		secondECNttDone.Sub(begin).String(),
		preprocessDone.Sub(begin).String(),
		msmDone.Sub(preprocessDone).String(),
		firstECNttDone.Sub(msmDone).String(),
		secondECNttDone.Sub(firstECNttDone).String(),
	)

	// only takes the first half
	return flatProofsBatch, nil
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
