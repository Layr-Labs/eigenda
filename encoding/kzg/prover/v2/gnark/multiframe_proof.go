package gnark

import (
	"fmt"
	"log/slog"
	"slices"
	"time"

	"github.com/Layr-Labs/eigenda/encoding/fft"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"golang.org/x/sync/errgroup"
)

type KzgMultiProofGnarkBackend struct {
	Fs *fft.FFTSettings
	// FFTPointsT contains the transposed SRSTable points, of size [2*dimE][l]=[2*numChunks][chunkLen].
	// See section 3.1.1 of https://github.com/khovratovich/Kate/blob/master/Kate_amortized.pdf:
	//   "Note that the vector multiplied by the matrix is independent from the polynomial coefficients,
	//   so its Fourier transform can be precomputed"
	FFTPointsT [][]bn254.G1Affine
	SFs        *fft.FFTSettings
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

	g := new(errgroup.Group)
	for i := uint64(0); i < dimE*2; i++ {
		g.Go(func() error {
			// eqn (3) u=y*v
			_, err := sumVec[i].MultiExp(p.FFTPointsT[i], coeffStore[i], ecc.MultiExpConfig{})
			if err != nil {
				return fmt.Errorf("multi exp: %w", err)
			}
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return nil, fmt.Errorf("errgroup: %w", err)
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

	// Now that we have h, we compute C_T = FFT(h).
	// See https://github.com/khovratovich/Kate/blob/master/Kate_amortized.pdf eqn 29.
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
	coeffStore := make([][]fr.Element, dimE*2)
	for i := range coeffStore {
		coeffStore[i] = make([]fr.Element, l)
	}

	// Worker pool to compute each column of coeffStore in parallel
	g := new(errgroup.Group)
	g.SetLimit(int(numWorker))
	for j := range l {
		g.Go(func() error {
			coeffs, err := p.getSlicesCoeff(polyFr, dimE, j, l)
			if err != nil {
				return fmt.Errorf("get slices coeff: %w", err)
			}
			for i := range len(coeffs) {
				// fill in coeffStore column j with coeffs
				coeffStore[i][j] = coeffs[i]
			}
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return nil, fmt.Errorf("errgroup: %w", err)
	}
	return coeffStore, nil
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
func (p *KzgMultiProofGnarkBackend) getSlicesCoeff(polyFr []fr.Element, dimE, j, l uint64) ([]fr.Element, error) {
	toeplitzExtendedVec := make([]fr.Element, 2*dimE)

	m := uint64(len(polyFr)) - 1 // there is a constant term
	dim := (m - j) / l
	for i := range dim {
		toeplitzExtendedVec[i].Set(&polyFr[m-(j+i*l)])
	}
	// Abstracting away the complex indices needed for extracting the multiproof coset,
	// toeplitzExtendedVec here looks like: [f_m,f_{m-1},..., f_0,0,0,...,0] (half zeros)
	// We then reverse it to put it in circulant form: [f_m,0 ,0...,0, f_1,f_1,...,f_{m-1}]
	// This matches Proposition 2 item 2 of https://eprint.iacr.org/2023/033.pdf.
	// Note that this only works because our toeplitz matrix contains many zeros and because
	// we set the extra free diagonal to 0 (alin's blog post uses a_0 for that diagonal).
	// For the generic case, see: https://alinush.github.io/2020/03/19/multiplying-a-vector-by-a-toeplitz-matrix.html
	slices.Reverse(toeplitzExtendedVec[1:])

	out, err := p.Fs.FFT(toeplitzExtendedVec, false)
	if err != nil {
		return nil, fmt.Errorf("fft: %w", err)
	}
	return out, nil
}
