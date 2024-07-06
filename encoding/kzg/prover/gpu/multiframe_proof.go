package gpu

import (
	"fmt"
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

type GpuComputer struct {
	*kzg.KzgConfig
	Fs         *fft.FFTSettings
	FFTPointsT [][]bn254.G1Affine // transpose of FFTPoints
	SFs        *fft.FFTSettings
	Srs        *kzg.SRS
	G2Trailing []bn254.G2Affine
	cfg        core.NTTConfig[[bn254_icicle.SCALAR_LIMBS]uint32]
}

// benchmarks shows cpu commit on 2MB blob only takes 24.165562ms. For now, use cpu
func (p *GpuComputer) ComputeLengthProof(coeffs []fr.Element) (*bn254.G2Affine, error) {
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
func (p *GpuComputer) ComputeCommitment(coeffs []fr.Element) (*bn254.G1Affine, error) {
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
func (p *GpuComputer) ComputeLengthCommitment(coeffs []fr.Element) (*bn254.G2Affine, error) {
	config := ecc.MultiExpConfig{}

	var lengthCommitment bn254.G2Affine
	_, err := lengthCommitment.MultiExp(p.Srs.G2[:len(coeffs)], coeffs, config)
	if err != nil {
		return nil, err
	}
	return &lengthCommitment, nil
}

func (p *GpuComputer) ComputeMultiFrameProof(polyFr []fr.Element, numChunks, chunkLen, numWorker uint64) ([]bn254.G1Affine, error) {
	// Robert: Standardizing this to use the same math used in precomputeSRS
	dimE := numChunks
	l := chunkLen
	numPoly := uint64(len(polyFr)) / dimE / chunkLen

	p.cfg = SetupNTT()

	begin := time.Now()
	jobChan := make(chan uint64, numWorker)
	results := make(chan WorkerResult, numWorker)

	// create storage for intermediate coefficients
	coeffStore := make([][]fr.Element, l*numPoly)
	for i := range coeffStore {
		coeffStore[i] = make([]fr.Element, dimE*2)
	}

	fmt.Println("numPoly", numPoly)
	fmt.Println("numChunks, ", numChunks)
	fmt.Println("chunkLen", chunkLen)

	for j := 0; j < len(polyFr); j++ {
		fmt.Printf("%v ", polyFr[j].String())
	}
	fmt.Println("len(polyFr)", len(polyFr))

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
	t_prepare := time.Now()

	if err != nil {
		return nil, fmt.Errorf("proof worker error: %v", err)
	}

	fmt.Println("coeffStore")
	for i := 0; i < len(coeffStore); i++ {
		a := coeffStore[i]
		for j := 0; j < len(a); j++ {
			fmt.Printf("%v ", a[j].String())
		}
		fmt.Println()
	}

	fmt.Println("NTT")
	// setup batch size for NTT
	p.cfg.BatchSize = int32(dimE * 2)
	coeffStoreFFT, e := p.NTT(coeffStore)
	if e != nil {
		return nil, e
	}

	// transpose it
	coeffStoreFFTT := make([][]fr.Element, dimE*2*numPoly)
	for i := range coeffStoreFFTT {
		coeffStoreFFTT[i] = make([]fr.Element, l)
	}

	t_ntt := time.Now()

	for k := uint64(0); k < numPoly; k++ {
		step := int(k * dimE * 2)
		for i := 0; i < int(l); i++ {
			vec := coeffStoreFFT[i+int(k*l)]
			for j := 0; j < int(dimE*2); j++ {
				coeffStoreFFTT[j+step][i] = vec[j]
			}
		}
	}

	fmt.Println("Transposed FFT")
	for i := 0; i < len(coeffStoreFFTT); i++ {
		vec := coeffStoreFFTT[i]
		for j := 0; j < len(vec); j++ {
			fmt.Printf("%v ", vec[j].String())
		}
		fmt.Println()
	}

	t0 := time.Now()
	fmt.Println("MsmBatch")
	sumVec, err := p.MsmBatch(coeffStoreFFTT, p.FFTPointsT)
	if err != nil {
		return nil, err
	}

	t1 := time.Now()

	fmt.Println("ECNTT inverse")
	// set new batch size for ntt, this equals to number of blobs
	p.cfg.BatchSize = int32(numPoly)
	sumVecInv, err := p.ECNtt(sumVec, true)
	if err != nil {
		return nil, err
	}

	t2 := time.Now()

	// remove half points per poly
	batchInv := make([]bn254.G1Affine, len(sumVecInv)/2)
	// outputs is out of order - buttefly
	k := 0
	for i := 0; i < int(numPoly); i++ {
		for j := 0; j < int(dimE); j++ {
			batchInv[k] = sumVecInv[i*int(dimE)*2+j]
			k += 1
		}
	}
	fmt.Println("ECNTT last")
	flatProofsBatch, err := p.ECNtt(batchInv, false)
	if err != nil {
		return nil, fmt.Errorf("second ECNtt error: %w", err)
	}

	t3 := time.Now()

	fmt.Println("flatProofsBatch")

	for j := 0; j < len(flatProofsBatch); j++ {
		fmt.Printf("%v ", flatProofsBatch[j].String())
	}

	fmt.Printf("prepare %v, ntt %v,\n", t_prepare.Sub(begin), t_ntt.Sub(t_prepare))
	fmt.Printf("total %v mult-th %v, msm %v,fft1 %v, fft2 %v,\n", t3.Sub(begin), t0.Sub(begin), t1.Sub(t0), t2.Sub(t1), t3.Sub(t2))

	return flatProofsBatch, nil
}

func (p *GpuComputer) proofWorkerGPU(
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

func (p *GpuComputer) GetSlicesCoeffWithoutFFT(polyFr []fr.Element, dimE, j, l uint64) ([]fr.Element, error) {
	// there is a constant term
	m := uint64(len(polyFr)) - 1
	dim := (m - j) / l

	toeV := make([]fr.Element, 2*dimE-1)
	for i := uint64(0); i < dim; i++ {

		toeV[i].Set(&polyFr[m-(j+i*l)])
	}

	// use precompute table
	tm, err := toeplitz.NewToeplitz(toeV, p.SFs)
	if err != nil {
		return nil, err
	}
	return tm.GetCoeff()
}

/*
// capable of batching blobs
func (p *GpuComputer) GetSlicesCoeffWithoutFFT(polyFr []fr.Element, dimE, j, l uint64) ([]fr.Element, error) {
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
*/
