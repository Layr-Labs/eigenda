package gnark

import (
	"fmt"
	"log/slog"
	"math"
	"time"

	"github.com/Layr-Labs/eigenda/encoding/fft"
	"github.com/Layr-Labs/eigenda/encoding/utils/toeplitz"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

type KzgMultiProofGnarkBackend struct {
	Fs *fft.FFTSettings
	// FFTPointsT contains the transposed SRSTable points.
	// See section 3.1 of https://eprint.iacr.org/2023/033.pdf:
	//   "Note that the vector multiplied by the matrix is independent from the polynomial coefficients,
	//   so its Fourier transform can be precomputed"
	FFTPointsT [][]bn254.G1Affine // transpose of FFTPoints
	SFs        *fft.FFTSettings
}

type WorkerResult struct {
	err error
}

// Computes a KZG multi-reveal proof for chunks containing in each frame.
//
// Each RS encoded blob contains numChunks*chunkLen field elements (symbols).
// For each chunk, we generate a multiproof opening for the chunkLen field elements
// belonging to that chunk.
// There are thus 2 levels of acceleration:
// 1. multiproof generates a single proof per chunk, revealing all field elements contained in that chunk.
// 2. each of the numChunks multiproofs are generated in parallel
//
// This algorithm is described in the "Fast Amortized KZG/Kate Proofs" papers. For background, read:
// 1. https://dankradfeist.de/ethereum/2020/06/16/kate-polynomial-commitments.html (single multiproof theory)
// 2. https://eprint.iacr.org/2023/033.pdf (how to compute the single multiproof fast)
// 3. https://github.com/khovratovich/Kate/blob/master/Kate_amortized.pdf (fast multiple multiproofs)
func (p *KzgMultiProofGnarkBackend) ComputeMultiFrameProofV2(
	polyFr []fr.Element, numChunks, chunkLen, numWorker uint64,
) ([]bn254.G1Affine, error) {
	// We describe the steps in the computation by following section 2.2 of
	// https://eprint.iacr.org/2023/033.pdf, generalized to the multiple multiproofs case.
	// eqn (1) DFT_2d(s^) is already precomputed and stored in [p.FFTPointsT].

	begin := time.Now()
	// Robert: Standardizing this to use the same math used in precomputeSRS
	dimE := numChunks
	l := chunkLen

	// eqn (2) DFT_2d(c^)
	coeffStore, err := p.computeCoeffStore(polyFr, numWorker, l, dimE)
	if err != nil {
		return nil, fmt.Errorf("coefficient computation error: %w", err)
	}
	preprocessDone := time.Now()

	// compute proof by multi scaler multiplication
	sumVec := make([]bn254.G1Affine, dimE*2)
	msmErrors := make(chan error, dimE*2)
	for i := uint64(0); i < dimE*2; i++ {

		go func(k uint64) {
			// eqn (3) u=y*v
			_, err := sumVec[k].MultiExp(p.FFTPointsT[k], coeffStore[k], ecc.MultiExpConfig{})
			// handle error
			msmErrors <- err
		}(i)
	}

	for i := uint64(0); i < dimE*2; i++ {
		err := <-msmErrors
		if err != nil {
			return nil, fmt.Errorf("msm error: %w", err)
		}
	}

	msmDone := time.Now()

	// eqn (4) h^ = iDFT_2d(u)
	sumVecInv, err := p.Fs.FFTG1(sumVec, true)
	if err != nil {
		return nil, fmt.Errorf("fft error: %w", err)
	}

	firstECNttDone := time.Now()

	// last step (5) "take first d elements of h^ as h"
	h := sumVecInv[:dimE]

	// Now that we have h, we compute C_T = FFT(h), from section 2.1.
	// Also see https://github.com/khovratovich/Kate/blob/master/Kate_amortized.pdf section 2
	// for more explanation as to why we take the FFT.
	// outputs is out of order - butterfly
	proofs, err := p.Fs.FFTG1(h, false)
	if err != nil {
		return nil, fmt.Errorf("fft error: %w", err)
	}

	secondECNttDone := time.Now()

	slog.Info("Multiproof Time Decomp",
		"total", secondECNttDone.Sub(begin),
		"preproc", preprocessDone.Sub(begin),
		"msm", msmDone.Sub(preprocessDone),
		"fft1", firstECNttDone.Sub(msmDone),
		"fft2", secondECNttDone.Sub(firstECNttDone),
	)

	return proofs, nil
}

// Helper function to handle coefficient computation.
// Returns a [2*dimE][l] slice.
func (p *KzgMultiProofGnarkBackend) computeCoeffStore(
	polyFr []fr.Element, numWorker, l, dimE uint64,
) ([][]fr.Element, error) {
	jobChan := make(chan uint64, numWorker)
	results := make(chan WorkerResult, numWorker)

	coeffStore := make([][]fr.Element, dimE*2)
	for i := range coeffStore {
		coeffStore[i] = make([]fr.Element, l)
	}

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
		return nil, fmt.Errorf("proof worker error: %w", lastErr)
	}

	return coeffStore, nil
}

func (p *KzgMultiProofGnarkBackend) proofWorker(
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
func (p *KzgMultiProofGnarkBackend) GetSlicesCoeff(polyFr []fr.Element, dimE, j, l uint64) ([]fr.Element, error) {
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
		return nil, fmt.Errorf("toeplitz new: %w", err)
	}
	e, err := tm.GetFFTCoeff()
	if err != nil {
		return nil, fmt.Errorf("toeplitz get fft coeff: %w", err)
	}
	return e, nil
}

/*
returns the power of 2 which is immediately bigger than the input
*/
func CeilIntPowerOf2Num(d uint64) uint64 {
	nextPower := math.Ceil(math.Log2(float64(d)))
	return uint64(math.Pow(2.0, nextPower))
}
