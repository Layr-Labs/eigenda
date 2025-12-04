package verifier

import (
	"errors"
	"fmt"
	"math"
	"math/big"
	"sync"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"

	eigenbn254 "github.com/Layr-Labs/eigenda/crypto/ecc/bn254"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/v2/fft"
	"github.com/Layr-Labs/eigenda/encoding/v2/kzg"
	"github.com/Layr-Labs/eigenda/encoding/v2/rs"
	"github.com/Layr-Labs/eigenda/resources/srs"

	_ "go.uber.org/automaxprocs"
)

type Verifier struct {
	G1SRS kzg.G1SRS

	// mu protects access to ParametrizedVerifiers
	mu                    sync.Mutex
	ParametrizedVerifiers map[encoding.EncodingParams]*ParametrizedVerifier
}

func NewVerifier(config *Config) (*Verifier, error) {
	if config.SRSNumberToLoad > encoding.SRSOrder {
		return nil, errors.New("SRSOrder is less than srsNumberToLoad")
	}

	// read the whole order, and treat it as entire SRS for low degree proof
	g1SRS, err := kzg.ReadG1Points(config.G1Path, config.SRSNumberToLoad, config.NumWorker)
	if err != nil {
		return nil, fmt.Errorf("failed to read %d G1 points from %s: %w", config.SRSNumberToLoad, config.G1Path, err)
	}

	encoderGroup := &Verifier{
		G1SRS:                 g1SRS,
		ParametrizedVerifiers: make(map[encoding.EncodingParams]*ParametrizedVerifier),
	}

	return encoderGroup, nil
}

func (v *Verifier) getKzgVerifier(params encoding.EncodingParams) (*ParametrizedVerifier, error) {
	if err := encoding.ValidateEncodingParams(params, encoding.SRSOrder); err != nil {
		return nil, fmt.Errorf("validate encoding params: %w", err)
	}

	// protect access to ParametrizedVerifiers
	v.mu.Lock()
	defer v.mu.Unlock()

	ver, ok := v.ParametrizedVerifiers[params]
	if ok {
		return ver, nil
	}

	ver, err := v.newKzgVerifier(params)
	if err != nil {
		return nil, fmt.Errorf("new KZG verifier: %w", err)
	}

	v.ParametrizedVerifiers[params] = ver
	return ver, nil
}

func (v *Verifier) newKzgVerifier(params encoding.EncodingParams) (*ParametrizedVerifier, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("invalid encoding params: %w", err)
	}

	// Create FFT settings based on params
	n := uint8(math.Log2(float64(params.NumEvaluations())))
	fs := fft.NewFFTSettings(n)

	return &ParametrizedVerifier{
		g1SRS: v.G1SRS,
		Fs:    fs,
	}, nil
}

// VerifyFrame verifies a single frame against a commitment.
// If needing to verify multiple frames of the same chunk length, prefer [Verifier.UniversalVerify].
//
// This function is only used in the v1 and v2 validator (distributed) retrievers.
// TODO(samlaf): replace with UniversalVerifySubBatch, and consider deleting this function.
func (v *Verifier) VerifyFrames(
	frames []*encoding.Frame,
	indices []encoding.ChunkNumber,
	commitments encoding.BlobCommitments,
	params encoding.EncodingParams) error {

	if len(frames) != len(indices) {
		return fmt.Errorf("invalid number of frames and indices: %d != %d", len(frames), len(indices))
	}

	verifier, err := v.getKzgVerifier(params)
	if err != nil {
		return err
	}

	for ind := range frames {
		err = verifier.verifyFrame(
			frames[ind],
			uint64(indices[ind]),
			(*bn254.G1Affine)(commitments.Commitment),
			params.NumChunks,
		)

		if err != nil {
			return err
		}
	}

	return nil
}

// TODO(mooselumph): Cleanup this function
func (v *Verifier) UniversalVerifySubBatch(
	params encoding.EncodingParams, samplesCore []encoding.Sample, numBlobs int,
) error {

	samples := make([]Sample, len(samplesCore))

	for i, sc := range samplesCore {
		x, err := rs.GetLeadingCosetIndex(
			uint64(sc.AssignmentIndex),
			params.NumChunks,
		)
		if err != nil {
			return fmt.Errorf("get leading coset index: %w", err)
		}

		sample := Sample{
			Commitment: (bn254.G1Affine)(*sc.Commitment),
			Proof:      sc.Chunk.Proof,
			RowIndex:   sc.BlobIndex,
			Coeffs:     sc.Chunk.Coeffs,
			X:          uint(x),
		}
		samples[i] = sample
	}

	return v.universalVerify(params, samples, numBlobs)
}

// Sample is the basic unit for a verification
// A blob may contain multiple Samples
type Sample struct {
	Commitment bn254.G1Affine
	Proof      bn254.G1Affine
	RowIndex   int // corresponds to a row in the verification matrix
	Coeffs     []fr.Element
	X          uint // X is the evaluating index which corresponds to the leading coset
}

// the rhsG1 consists of three terms, see
// https://ethresear.ch/t/a-universal-verification-equation-for-data-availability-sampling/13240/1
func genRhsG1(
	samples []Sample, randomsFr []fr.Element, m int,
	params encoding.EncodingParams, fftSettings *fft.FFTSettings, g1SRS kzg.G1SRS, proofs []bn254.G1Affine,
) (*bn254.G1Affine, error) {
	n := len(samples)
	commits := make([]bn254.G1Affine, m)
	D := params.ChunkLength

	var tmp fr.Element

	// first term
	// get coeffs to compute the aggregated commitment
	// note the coeff is affected by how many chunks are validated per blob
	// if x chunks are sampled from one blob, we need to compute the sum of all
	// x random field element corresponding to each sample
	aggCommitCoeffs := make([]fr.Element, m)
	setCommit := make([]bool, m)
	for k := 0; k < n; k++ {
		s := samples[k]
		row := s.RowIndex

		aggCommitCoeffs[row].Add(&aggCommitCoeffs[row], &randomsFr[k])

		if !setCommit[row] {
			commits[row].Set(&s.Commitment)

			setCommit[row] = true
		} else {

			if !commits[row].Equal(&s.Commitment) {
				return nil, errors.New("samples of the same row has different commitments")
			}
		}
	}

	var aggCommit bn254.G1Affine
	_, err := aggCommit.MultiExp(commits, aggCommitCoeffs, ecc.MultiExpConfig{})
	if err != nil {
		return nil, fmt.Errorf("compute aggregated commitment G1: %w", err)
	}

	// second term
	// compute the aggregated interpolation polynomial
	aggPolyCoeffs := make([]fr.Element, D)

	// we sum over the weighted coefficients (by the random field element) over all D monomial in all n samples
	for k := 0; k < n; k++ {
		coeffs := samples[k].Coeffs

		rk := randomsFr[k]
		// for each monomial in a given polynomial, multiply its coefficient with the corresponding random field,
		// then sum it with others. Given ChunkLen (D) is identical for all samples in a subBatch.
		// The operation is always valid.
		for j := uint64(0); j < D; j++ {
			tmp.Mul(&coeffs[j], &rk)
			//bls.MulModFr(&tmp, &coeffs[j], &rk)
			//bls.AddModFr(&aggPolyCoeffs[j], &aggPolyCoeffs[j], &tmp)
			aggPolyCoeffs[j].Add(&aggPolyCoeffs[j], &tmp)
		}
	}

	// All samples in a subBatch has identical chunkLen
	var aggPolyG1 bn254.G1Affine
	_, err = aggPolyG1.MultiExp(g1SRS[:D], aggPolyCoeffs, ecc.MultiExpConfig{})
	if err != nil {
		return nil, fmt.Errorf("failed to compute aggregated polynomial G1: %w", err)
	}

	// third term
	// leading coset is an evaluation index, here we compute the weighted leading coset evaluation by random fields
	lcCoeffs := make([]fr.Element, n)

	// get leading coset powers
	leadingDs := make([]fr.Element, n)
	bigD := big.NewInt(int64(D))

	for k := 0; k < n; k++ {

		// got the leading coset field element
		h := fftSettings.ExpandedRootsOfUnity[samples[k].X]
		var hPow fr.Element
		hPow.Exp(h, bigD)
		leadingDs[k].Set(&hPow)
	}

	// applying the random weights to leading coset elements
	for k := 0; k < n; k++ {
		rk := randomsFr[k]

		lcCoeffs[k].Mul(&rk, &leadingDs[k])
	}

	var offsetG1 bn254.G1Affine
	_, err = offsetG1.MultiExp(proofs, lcCoeffs, ecc.MultiExpConfig{})
	if err != nil {
		return nil, fmt.Errorf("failed to compute offset G1: %w", err)
	}

	var rhsG1 bn254.G1Affine

	rhsG1.Sub(&aggCommit, &aggPolyG1)

	rhsG1.Add(&rhsG1, &offsetG1)
	return &rhsG1, nil
}

// UniversalVerify implements batch verification on a set of chunks given the same chunk dimension (chunkLen, numChunk).
// The details is given in Ethereum Research post whose authors are George Kadianakis, Ansgar Dietrichs, Dankrad Feist
// https://ethresear.ch/t/a-universal-verification-equation-for-data-availability-sampling/13240
//
// samples is a list of chunks. The order of samples do not matter.
// Each sample need not have unique row, it is possible that multiple chunks of the same blob are validated altogether
func (v *Verifier) universalVerify(params encoding.EncodingParams, samples []Sample, numBlobs int) error {
	// precheck
	for _, s := range samples {
		if s.RowIndex >= numBlobs {
			return fmt.Errorf(
				"sample.RowIndex and numBlob are inconsistent: sample has %d rows, but there are only %d blobs",
				s.RowIndex, numBlobs)
		}
	}

	verifier, err := v.getKzgVerifier(params)
	if err != nil {
		return err
	}

	D := params.ChunkLength

	if D > uint64(len(v.G1SRS)) {
		return fmt.Errorf("requested chunkLen %v is larger than Loaded G1SRS points %v", D, len(v.G1SRS))
	}

	n := len(samples)
	if n == 0 {
		return errors.New("the number of samples (i.e. chunks) must not be empty")
	}

	// generate random field elements to aggregate equality check
	randomsFr, err := eigenbn254.RandomFrs(n)
	if err != nil {
		return fmt.Errorf("create randomness vector: %w", err)
	}

	// array of proofs
	proofs := make([]bn254.G1Affine, n)
	for i := 0; i < n; i++ {
		proofs[i].Set(&samples[i].Proof)
	}

	// lhs g1
	var lhsG1 bn254.G1Affine
	_, err = lhsG1.MultiExp(proofs, randomsFr, ecc.MultiExpConfig{})
	if err != nil {
		return fmt.Errorf("compute lhsG1: %w", err)
	}

	// lhs g2
	exponent := uint64(math.Log2(float64(D)))
	G2atD := srs.G2PowerOf2SRS[exponent]
	lhsG2 := &G2atD

	// rhs g2
	rhsG2 := &kzg.GenG2

	// rhs g1
	rhsG1, err := genRhsG1(
		samples,
		randomsFr,
		numBlobs,
		params,
		verifier.Fs,
		verifier.g1SRS,
		proofs,
	)
	if err != nil {
		return fmt.Errorf("generate rhsG1: %w", err)
	}

	err = eigenbn254.PairingsVerify(&lhsG1, lhsG2, rhsG1, rhsG2)
	if err != nil {
		return fmt.Errorf("verify pairing: %w", err)
	}
	return nil
}
