package gpu

import (
	"fmt"
	"sync"
	"time"

	"github.com/Layr-Labs/eigenda/encoding/fft"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/Layr-Labs/eigenda/encoding/utils/toeplitz"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/ingonyama-zk/icicle/v2/wrappers/golang/core"
	bn254_icicle "github.com/ingonyama-zk/icicle/v2/wrappers/golang/curves/bn254"
)

type WorkerResult struct {
	points []bn254.G1Affine
	err    error
}

type GpuComputeDevice struct {
	*kzg.KzgConfig
	Fs             *fft.FFTSettings
	FlatFFTPointsT []bn254_icicle.Affine
	SFs            *fft.FFTSettings
	Srs            *kzg.SRS
	G2Trailing     []bn254.G2Affine
	NttCfg         core.NTTConfig[[bn254_icicle.SCALAR_LIMBS]uint32]
	GpuLock        *sync.Mutex // lock whenever gpu is needed,
}

// benchmarks shows cpu commit on 2MB blob only takes 24.165562ms. For now, use cpu
func (p *GpuComputeDevice) ComputeLengthProof(coeffs []fr.Element) (*bn254.G2Affine, error) {
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

// benchmarks shows cpu commit on 2MB blob only takes 11.673738ms. For now, use cpu
func (p *GpuComputeDevice) ComputeCommitment(coeffs []fr.Element) (*bn254.G1Affine, error) {
	// compute commit for the full poly
	config := ecc.MultiExpConfig{}
	var commitment bn254.G1Affine
	_, err := commitment.MultiExp(p.Srs.G1[:len(coeffs)], coeffs, config)
	if err != nil {
		return nil, err
	}
	return &commitment, nil
}

// benchmarks shows cpu commit on 2MB blob only takes 31.318661ms. For now, use cpu
func (p *GpuComputeDevice) ComputeLengthCommitment(coeffs []fr.Element) (*bn254.G2Affine, error) {
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
func (p *GpuComputeDevice) ComputeMultiFrameProof(polyFr []fr.Element, numChunks, chunkLen, numWorker uint64) ([]bn254.G1Affine, error) {
	// Robert: Standardizing this to use the same math used in precomputeSRS
	dimE := numChunks
	l := chunkLen
	numPoly := uint64(len(polyFr)) / dimE / chunkLen
	fmt.Println("numPoly", numPoly)

	begin := time.Now()

	// create storage for intermediate coefficients matrix
	jobChan := make(chan uint64, numWorker)
	results := make(chan WorkerResult, numWorker)

	coeffStore := make([][]fr.Element, l*numPoly)
	for i := range coeffStore {
		coeffStore[i] = make([]fr.Element, dimE*2)
	}

	// Preprpcessing use CPU to compute those coefficients based on polyFr
	for w := uint64(0); w < numWorker; w++ {
		go p.proofWorkerGPU(polyFr, jobChan, l, dimE, coeffStore, results)
	}

	for j := uint64(0); j < l*numPoly; j++ {
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
	// Preprpcessing Completed
	if err != nil {
		return nil, fmt.Errorf("proof worker error: %v", err)
	}
	preprocessDone := time.Now()

	// Start using GPU
	p.GpuLock.Lock()
	defer p.GpuLock.Unlock()

	// Compute NTT on the coeff matrix
	p.NttCfg.BatchSize = int32(l * numPoly)
	coeffStoreFft, e := p.NTT(coeffStore)
	if e != nil {
		return nil, e
	}
	nttDone := time.Now()

	/*
		fmt.Println("after fft")
		vec := gpu_utils.ConvertScalarFieldsToFrBytes(coeffStoreFft)
		for i := 0; i < int(l*numPoly); i++ {
			length := int(dimE) * 2
			for j := 0; j < length; j++ {
				fmt.Printf("%v ", vec[i*length+j].String())
			}
			fmt.Println()
		}
	*/

	// transpose the FFT tranformed matrix
	coeffStoreFftTranspose, err := Transpose(coeffStoreFft, int(l), int(numPoly), int(dimE)*2)
	if err != nil {
		return nil, e
	}
	transposingDone := time.Now()

	// compute msm on each rows of the transposed matrix
	sumVec, err := p.MsmBatch(coeffStoreFftTranspose, p.FlatFFTPointsT, int(numPoly)*int(dimE)*2)
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
	sumVec.Free()

	// extract proofs
	prunedSumVecInv := core.HostSliceWithValue(bn254_icicle.Projective{}, len(sumVecInv)/2)
	k := 0
	for i := 0; i < int(numPoly); i++ {
		for j := 0; j < int(dimE); j++ {
			prunedSumVecInv[k] = sumVecInv[i*int(dimE)*2+j]
			k += 1
		}
	}

	// compute the second ecntt on the reduced size array
	flatProofsBatch, err := p.ECNttToGnark(prunedSumVecInv, false, int(numPoly)*int(dimE))
	if err != nil {
		return nil, fmt.Errorf("second ECNtt error: %w", err)
	}

	/*
		// debug
		for i := 0; i < int(numPoly); i++ {
			for j := 0; j < int(dimE); j++ {
				fmt.Printf("%v ", flatProofsBatch[i*int(dimE)+j].String())
			}
			fmt.Println()
		}
	*/

	secondECNttDone := time.Now()

	fmt.Printf("Multiproof Time Decomp \n\t\ttotal   %-20v \n\t\tpreproc %-20v \n\t\tntt     %-20v \n\t\ttranspose %-20v \n\t\tmsm     %-v \n\t\tfft1    %-v \n\t\tfft2    %-v,\n",
		secondECNttDone.Sub(begin),
		preprocessDone.Sub(begin),
		nttDone.Sub(preprocessDone),
		transposingDone.Sub(nttDone),
		msmDone.Sub(transposingDone),
		firstECNttDone.Sub(msmDone),
		secondECNttDone.Sub(firstECNttDone),
	)

	// only takes the first half
	return flatProofsBatch, nil
}

func (p *GpuComputeDevice) proofWorkerGPU(
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
func (p *GpuComputeDevice) GetSlicesCoeffWithoutFFT(polyFr []fr.Element, dimE, j, l uint64) ([]fr.Element, error) {
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
