package prover

import (
	"fmt"
	"log"
	"math"
	"time"

	"github.com/Layr-Labs/eigenda/encoding"

	"github.com/Layr-Labs/eigenda/encoding/fft"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/Layr-Labs/eigenda/encoding/rs"
	"github.com/Layr-Labs/eigenda/encoding/utils/toeplitz"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

type ParametrizedProver struct {
	*rs.Encoder

	*kzg.KzgConfig
	Srs        *kzg.SRS
	G2Trailing []bn254.G2Affine

	Fs         *fft.FFTSettings
	Ks         *kzg.KZGSettings
	SFs        *fft.FFTSettings   // fft used for submatrix product helper
	FFTPointsT [][]bn254.G1Affine // transpose of FFTPoints
	Holder     chan struct{}
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
	commit, err := g.Commit(poly.Coeffs)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	config := ecc.MultiExpConfig{}

	var lowDegreeCommitment bn254.G2Affine
	_, err = lowDegreeCommitment.MultiExp(g.Srs.G2[:len(poly.Coeffs)], poly.Coeffs, config)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

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
	var lowDegreeProof bn254.G2Affine
	_, err = lowDegreeProof.MultiExp(shiftedSecret, poly.Coeffs, config)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

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

func (g *ParametrizedProver) Commit(polyFr []fr.Element) (bn254.G1Affine, error) {
	commit, err := g.Ks.CommitToPoly(polyFr)
	return *commit, err
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

	t0 := time.Now()

	// compute proof by multi scaler mulplication
	msmErrors := make(chan error, dimE*2)
	for i := uint64(0); i < dimE*2; i++ {

		go func(k uint64) {
			_, err := sumVec[k].MultiExp(p.FFTPointsT[k], coeffStore[k], ecc.MultiExpConfig{})
			// handle error
			msmErrors <- err
		}(i)
	}

	for i := uint64(0); i < dimE*2; i++ {
		err := <-msmErrors
		if err != nil {
			fmt.Println("Error. MSM while adding points", err)
			return nil, err
		}
	}

	//<-p.Holder

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
