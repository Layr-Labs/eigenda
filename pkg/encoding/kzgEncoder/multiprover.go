package kzgEncoder

import (
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/Layr-Labs/eigenda/pkg/encoding/utils/toeplitz"
	bls "github.com/Layr-Labs/eigenda/pkg/kzg/bn254"
)

type WorkerResult struct {
	points []bls.G1Point
	err    error
}

func (p *KzgEncoder) ProveAllCosetThreads(polyFr []bls.Fr, numChunks, chunkLen, numWorker uint64) ([]bls.G1Point, error) {
	begin := time.Now()
	// Robert: Standardizing this to use the same math used in precomputeSRS
	dimE := numChunks
	l := chunkLen

	sumVec := make([]bls.G1Point, dimE*2)

	jobChan := make(chan uint64, numWorker)
	results := make(chan WorkerResult, numWorker)

	// create storage for intermediate fft outputs
	coeffStore := make([][]bls.Fr, dimE*2)
	for i := range coeffStore {
		coeffStore[i] = make([]bls.Fr, l)
	}

	for w := uint64(0); w < numWorker; w++ {
		go p.proofWorker(polyFr, jobChan, l, dimE, coeffStore, results)
	}

	for j := uint64(0); j < l; j++ {
		jobChan <- j
	}
	close(jobChan)

	// return only first error
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

	t0 := time.Now()

	// compute proof by multi scaler mulplication
	var wg sync.WaitGroup
	for i := uint64(0); i < dimE*2; i++ {
		wg.Add(1)
		go func(k uint64) {
			defer wg.Done()
			sumVec[k] = *bls.LinCombG1(p.FFTPointsT[k], coeffStore[k])
		}(i)
	}

	wg.Wait()

	t1 := time.Now()

	// only 1 ifft is needed
	sumVecInv, err := p.Fs.FFTG1(sumVec, true)
	if err != nil {
		return nil, fmt.Errorf("fft error: %v", err)
	}

	t2 := time.Now()

	// outputs is out of order - buttefly
	proofs, err := p.Fs.FFTG1(sumVecInv[:dimE], false)
	if err != nil {
		return nil, err
	}

	t3 := time.Now()

	fmt.Printf("mult-th %v, msm %v,fft1 %v, fft2 %v,\n", t0.Sub(begin), t1.Sub(t0), t2.Sub(t1), t3.Sub(t2))

	//rb.ReverseBitOrderG1Point(proofs)
	return proofs, nil
}

func (p *KzgEncoder) proofWorker(
	polyFr []bls.Fr,
	jobChan <-chan uint64,
	l uint64,
	dimE uint64,
	coeffStore [][]bls.Fr,
	results chan<- WorkerResult,
) {

	for j := range jobChan {
		coeffs, err := p.GetSlicesCoeff(polyFr, dimE, j, l)
		if err != nil {
			results <- WorkerResult{
				points: nil,
				err:    err,
			}
		}

		for i := 0; i < len(coeffs); i++ {
			coeffStore[i][j] = coeffs[i]
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
func (p *KzgEncoder) GetSlicesCoeff(polyFr []bls.Fr, dimE, j, l uint64) ([]bls.Fr, error) {
	// there is a constant term
	m := uint64(len(polyFr)) - 1
	dim := (m - j) / l

	toeV := make([]bls.Fr, 2*dimE-1)
	for i := uint64(0); i < dim; i++ {
		bls.CopyFr(&toeV[i], &polyFr[m-(j+i*l)])
	}

	// use precompute table
	tm, err := toeplitz.NewToeplitz(toeV, p.SFs)
	if err != nil {
		return nil, err
	}
	return tm.GetFFTCoeff()
}

/*
returns the power of 2 which is immediately bigger than the input
*/
func CeilIntPowerOf2Num(d uint64) uint64 {
	nextPower := math.Ceil(math.Log2(float64(d)))
	return uint64(math.Pow(2.0, nextPower))
}
