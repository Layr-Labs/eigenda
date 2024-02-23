package prover

import (
	"fmt"
	"log"
	"math"
	"sync"
	"time"

	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/kzgrs"
	"github.com/Layr-Labs/eigenda/encoding/rs"
	"github.com/Layr-Labs/eigenda/encoding/utils/toeplitz"
	kzg "github.com/Layr-Labs/eigenda/pkg/kzg"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

type ParametrizedProver struct {
	*rs.Encoder

	*kzgrs.KzgConfig
	Srs        *kzg.SRS
	G2Trailing []bn254.G2Affine

	Fs         *kzg.FFTSettings
	Ks         *kzg.KZGSettings
	SFs        *kzg.FFTSettings // fft used for submatrix product helper
	FFTPoints  [][]bn254.G1Affine
	FFTPointsT [][]bn254.G1Affine // transpose of FFTPoints
}

type WorkerResult struct {
	points []bn254.G1Affine
	err    error
}

// just a wrapper to take bytes not Fr Element
func (g *ParametrizedProver) EncodeBytes(inputBytes []byte) (*bn254.G1Affine, *bn254.G2Affine, *bn254.G2Affine, []encoding.Frame, []uint32, error) {
	inputFr := rs.ToFrArray(inputBytes)
	return g.Encode(inputFr)
}

func (g *ParametrizedProver) Encode(inputFr []fr.Element) (*bn254.G1Affine, *bn254.G2Affine, *bn254.G2Affine, []encoding.Frame, []uint32, error) {

	startTime := time.Now()
	poly, frames, indices, err := g.Encoder.Encode(inputFr)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	if len(poly.Coeffs) > int(g.KzgConfig.SRSNumberToLoad) {
		return nil, nil, nil, nil, nil, fmt.Errorf("poly Coeff length %v is greater than Loaded SRS points %v", len(poly.Coeffs), int(g.KzgConfig.SRSNumberToLoad))
	}

	// compute commit for the full poly
	commit := g.Commit(poly.Coeffs)
	//lowDegreeCommitment := bls.LinCombG2(g.Srs.G2[:len(poly.Coeffs)], poly.Coeffs)

	config := ecc.MultiExpConfig{}

	var lowDegreeCommitment bn254.G2Affine
	_, err = lowDegreeCommitment.MultiExp(g.Srs.G2[:len(poly.Coeffs)], poly.Coeffs, config)

	intermediate := time.Now()

	polyDegreePlus1 := uint64(len(inputFr))

	if g.Verbose {
		log.Printf("    Commiting takes  %v\n", time.Since(intermediate))
		intermediate = time.Now()

		log.Printf("shift %v\n", g.SRSOrder-polyDegreePlus1)
		log.Printf("order %v\n", len(g.Srs.G2))
		log.Println("low degree verification info")
	}

	shiftedSecret := g.G2Trailing[g.KzgConfig.SRSNumberToLoad-polyDegreePlus1:]

	//The proof of low degree is commitment of the polynomial shifted to the largest srs degree
	//lowDegreeProof := bls.LinCombG2(shiftedSecret, poly.Coeffs[:polyDegreePlus1])
	var lowDegreeProof bn254.G2Affine
	_, err = lowDegreeProof.MultiExp(shiftedSecret, poly.Coeffs, config)

	//fmt.Println("kzgFFT lowDegreeProof", lowDegreeProof, "poly len ", len(fullCoeffsPoly), "order", len(g.Ks.SecretG2) )
	//ok := VerifyLowDegreeProof(&commit, lowDegreeProof, polyDegreePlus1-1, g.SRSOrder, g.Srs.G2)
	//if !ok {
	//		log.Printf("Kzg FFT Cannot Verify low degree proof %v", lowDegreeProof)
	//		return nil, nil, nil, nil, errors.New("cannot verify low degree proof")
	//	} else {
	//		log.Printf("Kzg FFT Verify low degree proof  PPPASSS %v", lowDegreeProof)
	//	}

	if g.Verbose {
		log.Printf("    Generating Low Degree Proof takes  %v\n", time.Since(intermediate))
		intermediate = time.Now()
	}

	// compute proofs
	paddedCoeffs := make([]fr.Element, g.NumEvaluations())
	copy(paddedCoeffs, poly.Coeffs)

	proofs, err := g.ProveAllCosetThreads(paddedCoeffs, g.NumChunks, g.ChunkLength, g.NumWorker)
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("could not generate proofs: %v", err)
	}

	if g.Verbose {
		log.Printf("    Proving takes    %v\n", time.Since(intermediate))
	}

	kzgFrames := make([]encoding.Frame, len(frames))
	for i, index := range indices {
		kzgFrames[i] = encoding.Frame{
			Proof:  proofs[index],
			Coeffs: frames[i].Coeffs,
		}
	}

	if g.Verbose {
		log.Printf("Total encoding took      %v\n", time.Since(startTime))
	}
	return &commit, &lowDegreeCommitment, &lowDegreeProof, kzgFrames, indices, nil
}

func (g *ParametrizedProver) Commit(polyFr []fr.Element) bn254.G1Affine {
	commit := g.Ks.CommitToPoly(polyFr)
	return *commit
}

func (p *ParametrizedProver) ProveAllCosetThreads(polyFr []fr.Element, numChunks, chunkLen, numWorker uint64) ([]bn254.G1Affine, error) {
	begin := time.Now()
	// Robert: Standardizing this to use the same math used in precomputeSRS
	dimE := numChunks
	l := chunkLen

	sumVec := make([]bn254.G1Affine, dimE*2)

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
			//sumVec[k] = *bls.LinCombG1(p.FFTPointsT[k], coeffStore[k])
			sumVec[k].MultiExp(p.FFTPointsT[k], coeffStore[k], ecc.MultiExpConfig{})
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

func (p *ParametrizedProver) proofWorker(
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
func (p *ParametrizedProver) GetSlicesCoeff(polyFr []fr.Element, dimE, j, l uint64) ([]fr.Element, error) {
	// there is a constant term
	m := uint64(len(polyFr)) - 1
	dim := (m - j) / l

	toeV := make([]fr.Element, 2*dimE-1)
	for i := uint64(0); i < dim; i++ {
		//bls.CopyFr(&toeV[i], &polyFr[m-(j+i*l)])
		toeV[i].Set(&polyFr[m-(j+i*l)])
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
