package verifier

import (
	"errors"
	"fmt"
	"log/slog"
	"math"
	"math/big"
	"runtime"
	"sync"
	"time"

	"github.com/Layr-Labs/eigenda/encoding"

	"github.com/Layr-Labs/eigenda/encoding/fft"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/Layr-Labs/eigenda/encoding/rs"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	_ "go.uber.org/automaxprocs"
)

type Verifier struct {
	config    *encoding.Config
	kzgConfig *kzg.KzgConfig
	*rs.Encoder
	Srs          *kzg.SRS
	G2Trailing   []bn254.G2Affine
	mu           sync.Mutex
	LoadG2Points bool

	ParametrizedVerifiers map[encoding.EncodingParams]*ParametrizedVerifier
}

var _ encoding.Verifier = &Verifier{}

// Default configuration values
const (
	defaultLoadG2Points = true
	defaultVerbose      = false
)

func NewVerifier(opts ...VerifierOption) (*Verifier, error) {
	v := &Verifier{
		config: &encoding.Config{
			NumWorker: uint64(runtime.GOMAXPROCS(0)),
			Verbose:   defaultVerbose,
		},
		LoadG2Points:          defaultLoadG2Points,
		ParametrizedVerifiers: make(map[encoding.EncodingParams]*ParametrizedVerifier),
		mu:                    sync.Mutex{},
	}

	// Apply options
	for _, opt := range opts {
		if err := opt(v); err != nil {
			return nil, fmt.Errorf("failed to apply option: %w", err)
		}
	}

	// Validate required configurations
	if v.kzgConfig == nil {
		return nil, errors.New("KZG config is required")
	}

	if err := v.initializeSRS(); err != nil {
		return nil, err
	}

	// Create default RS encoder if none provided
	if v.Encoder == nil {
		encoder, err := rs.NewEncoder()
		if err != nil {
			return nil, fmt.Errorf("failed to create default RS encoder: %w", err)
		}
		v.Encoder = encoder
	}

	return v, nil
}

func (v *Verifier) initializeSRS() error {
	startTime := time.Now()
	s1, err := kzg.ReadG1Points(v.kzgConfig.G1Path, v.kzgConfig.SRSNumberToLoad, v.kzgConfig.NumWorker)
	if err != nil {
		return fmt.Errorf("failed to read G1 points: %w", err)
	}
	slog.Info("ReadG1Points", "time", time.Since(startTime), "numPoints", len(s1))

	s2 := make([]bn254.G2Affine, 0)
	g2Trailing := make([]bn254.G2Affine, 0)

	if v.LoadG2Points {
		if err := v.loadG2PointsData(&s2, &g2Trailing); err != nil {
			return err
		}
	} else if len(v.kzgConfig.G2PowerOf2Path) == 0 {
		return errors.New("G2PowerOf2Path is empty but required when loadG2Points is false")
	}

	srs, err := kzg.NewSrs(s1, s2)
	if err != nil {
		return fmt.Errorf("could not create srs: %w", err)
	}

	v.Srs = srs
	v.G2Trailing = g2Trailing

	return nil
}

func (v *Verifier) loadG2PointsData(s2 *[]bn254.G2Affine, g2Trailing *[]bn254.G2Affine) error {
	if len(v.kzgConfig.G2Path) == 0 {
		return errors.New("G2Path is empty but required when loadG2Points is true")
	}

	startTime := time.Now()
	points, err := kzg.ReadG2Points(v.kzgConfig.G2Path, v.kzgConfig.SRSNumberToLoad, v.kzgConfig.NumWorker)
	if err != nil {
		return fmt.Errorf("failed to read G2 points: %w", err)
	}
	slog.Info("ReadG2Points", "time", time.Since(startTime), "numPoints", len(points))
	*s2 = points

	startTime = time.Now()
	trailing, err := kzg.ReadG2PointSection(
		v.kzgConfig.G2Path,
		v.kzgConfig.SRSOrder-v.kzgConfig.SRSNumberToLoad,
		v.kzgConfig.SRSOrder,
		v.kzgConfig.NumWorker,
	)
	if err != nil {
		return fmt.Errorf("failed to read G2 point section: %w", err)
	}
	slog.Info("ReadG2PointSection", "time", time.Since(startTime), "numPoints", len(trailing))
	*g2Trailing = trailing

	return nil
}

type ParametrizedVerifier struct {
	*kzg.KzgConfig
	Srs *kzg.SRS

	Fs *fft.FFTSettings
	Ks *kzg.KZGSettings
}

func (v *Verifier) GetKzgVerifier(params encoding.EncodingParams) (*ParametrizedVerifier, error) {
	v.mu.Lock()
	defer v.mu.Unlock()

	if err := encoding.ValidateEncodingParams(params, v.kzgConfig.SRSOrder); err != nil {
		return nil, err
	}

	ver, ok := v.ParametrizedVerifiers[params]
	if ok {
		return ver, nil
	}

	ver, err := v.newKzgVerifier(params)
	if err == nil {
		v.ParametrizedVerifiers[params] = ver
	}

	return ver, err
}

func (g *Verifier) NewKzgVerifier(params encoding.EncodingParams) (*ParametrizedVerifier, error) {
	return g.newKzgVerifier(params)
}

func (v *Verifier) newKzgVerifier(params encoding.EncodingParams) (*ParametrizedVerifier, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	// Create FFT settings based on params
	n := uint8(math.Log2(float64(params.NumEvaluations())))
	fs := fft.NewFFTSettings(n)

	// Create KZG settings
	ks, err := kzg.NewKZGSettings(fs, v.Srs)
	if err != nil {
		return nil, fmt.Errorf("failed to create KZG settings: %w", err)
	}

	return &ParametrizedVerifier{
		KzgConfig: v.kzgConfig,
		Srs:       v.Srs,
		Fs:        fs,
		Ks:        ks,
	}, nil
}

func (v *Verifier) VerifyBlobLength(commitments encoding.BlobCommitments) error {
	return v.VerifyCommit(
		(*bn254.G2Affine)(commitments.LengthCommitment),
		(*bn254.G2Affine)(commitments.LengthProof),
		uint64(commitments.Length),
	)
}

// VerifyCommit verifies the low degree proof; since it doesn't depend on the encoding parameters
// we leave it as a method of the KzgEncoderGroup
func (v *Verifier) VerifyCommit(lengthCommit *bn254.G2Affine, lengthProof *bn254.G2Affine, length uint64) error {

	g1Challenge, err := kzg.ReadG1Point(v.kzgConfig.SRSOrder-length, v.kzgConfig)
	if err != nil {
		return err
	}

	err = VerifyLengthProof(lengthCommit, lengthProof, &g1Challenge)
	if err != nil {
		return fmt.Errorf("%v . %v ", "low degree proof fails", err)
	} else {
		return nil
	}
}

// The function verify low degree proof against a poly commitment
// We wish to show x^shift poly = shiftedPoly, with
// With shift = SRSOrder - length and
// proof = commit(shiftedPoly) on G1
// so we can verify by checking
// e( commit_1, [x^shift]_2) = e( proof_1, G_2 )
func VerifyLengthProof(lengthCommit *bn254.G2Affine, proof *bn254.G2Affine, g1Challenge *bn254.G1Affine) error {
	return PairingsVerify(g1Challenge, lengthCommit, &kzg.GenG1, proof)
}

func (v *Verifier) VerifyFrames(frames []*encoding.Frame, indices []encoding.ChunkNumber, commitments encoding.BlobCommitments, params encoding.EncodingParams) error {

	verifier, err := v.GetKzgVerifier(params)
	if err != nil {
		return err
	}

	for ind := range frames {
		err = verifier.VerifyFrame(
			(*bn254.G1Affine)(commitments.Commitment),
			frames[ind],
			uint64(indices[ind]),
			params.NumChunks,
		)

		if err != nil {
			return err
		}
	}

	return nil

}

func (v *ParametrizedVerifier) VerifyFrame(commit *bn254.G1Affine, f *encoding.Frame, index uint64, numChunks uint64) error {

	j, err := rs.GetLeadingCosetIndex(
		uint64(index),
		numChunks,
	)
	if err != nil {
		return err
	}

	g2Atn, err := kzg.ReadG2Point(uint64(len(f.Coeffs)), v.KzgConfig)
	if err != nil {
		return err
	}

	err = VerifyFrame(f, v.Ks, commit, &v.Ks.ExpandedRootsOfUnity[j], &g2Atn)

	if err != nil {
		return fmt.Errorf("%v . %v ", "VerifyFrame Error", err)
	} else {
		return nil
	}
}

// Verify function assumes the Data stored is coefficients of coset's interpolating poly
func VerifyFrame(f *encoding.Frame, ks *kzg.KZGSettings, commitment *bn254.G1Affine, x *fr.Element, g2Atn *bn254.G2Affine) error {
	var xPow fr.Element
	xPow.SetOne()

	for i := 0; i < len(f.Coeffs); i++ {
		xPow.Mul(&xPow, x)
	}

	var xPowBigInt big.Int

	// [x^n]_2
	var xn2 bn254.G2Affine

	xn2.ScalarMultiplication(&kzg.GenG2, xPow.BigInt(&xPowBigInt))

	// [s^n - x^n]_2
	var xnMinusYn bn254.G2Affine
	xnMinusYn.Sub(g2Atn, &xn2)

	// [interpolation_polynomial(s)]_1
	var is1 bn254.G1Affine
	config := ecc.MultiExpConfig{}
	_, err := is1.MultiExp(ks.Srs.G1[:len(f.Coeffs)], f.Coeffs, config)
	if err != nil {
		return err
	}

	// [commitment - interpolation_polynomial(s)]_1 = [commit]_1 - [interpolation_polynomial(s)]_1
	var commitMinusInterpolation bn254.G1Affine
	commitMinusInterpolation.Sub(commitment, &is1)

	// Verify the pairing equation
	//
	// e([commitment - interpolation_polynomial(s)], [1]) = e([proof],  [s^n - x^n])
	//    equivalent to
	// e([commitment - interpolation_polynomial]^(-1), [1]) * e([proof],  [s^n - x^n]) = 1_T
	//

	return PairingsVerify(&commitMinusInterpolation, &kzg.GenG2, &f.Proof, &xnMinusYn)
}

// Decode takes in the chunks, indices, and encoding parameters and returns the decoded blob
// The result is trimmed to the given maxInputSize.
func (v *Verifier) Decode(chunks []*encoding.Frame, indices []encoding.ChunkNumber, params encoding.EncodingParams, maxInputSize uint64) ([]byte, error) {
	frames := make([]rs.Frame, len(chunks))
	for i := range chunks {
		frames[i] = rs.Frame{
			Coeffs: chunks[i].Coeffs,
		}
	}

	return v.Encoder.Decode(frames, toUint64Array(indices), maxInputSize, params)
}

func toUint64Array(chunkIndices []encoding.ChunkNumber) []uint64 {
	res := make([]uint64, len(chunkIndices))
	for i, d := range chunkIndices {
		res[i] = uint64(d)
	}
	return res
}

func PairingsVerify(a1 *bn254.G1Affine, a2 *bn254.G2Affine, b1 *bn254.G1Affine, b2 *bn254.G2Affine) error {
	var negB1 bn254.G1Affine
	negB1.Neg(b1)

	P := [2]bn254.G1Affine{*a1, negB1}
	Q := [2]bn254.G2Affine{*a2, *b2}

	ok, err := bn254.PairingCheck(P[:], Q[:])
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("PairingCheck pairing not ok. SRS is invalid")
	}

	return nil
}
