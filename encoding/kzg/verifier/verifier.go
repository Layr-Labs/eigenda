package verifier

import (
	"errors"
	"fmt"
	"log"
	"math"
	"math/big"
	"runtime"
	"sync"

	"github.com/Layr-Labs/eigenda/encoding"

	"github.com/Layr-Labs/eigenda/encoding/fft"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/Layr-Labs/eigenda/encoding/rs"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

type Verifier struct {
	*kzg.KzgConfig
	Srs          *kzg.SRS
	G2Trailing   []bn254.G2Affine
	mu           sync.Mutex
	LoadG2Points bool

	ParametrizedVerifiers map[encoding.EncodingParams]*ParametrizedVerifier
}

var _ encoding.Verifier = &Verifier{}

func NewVerifier(config *kzg.KzgConfig, loadG2Points bool) (*Verifier, error) {

	if config.SRSNumberToLoad > config.SRSOrder {
		return nil, errors.New("SRSOrder is less than srsNumberToLoad")
	}

	// read the whole order, and treat it as entire SRS for low degree proof
	s1, err := kzg.ReadG1Points(config.G1Path, config.SRSNumberToLoad, config.NumWorker)
	if err != nil {
		return nil, fmt.Errorf("failed to read %d G1 points from %s: %v", config.SRSNumberToLoad, config.G1Path, err)
	}

	s2 := make([]bn254.G2Affine, 0)
	g2Trailing := make([]bn254.G2Affine, 0)

	// PreloadEncoder is by default not used by operator node, PreloadEncoder
	if loadG2Points {
		if len(config.G2Path) == 0 {
			return nil, errors.New("G2Path is empty. However, object needs to load G2Points")
		}

		s2, err = kzg.ReadG2Points(config.G2Path, config.SRSNumberToLoad, config.NumWorker)
		if err != nil {
			return nil, fmt.Errorf("failed to read %d G2 points from %s: %v", config.SRSNumberToLoad, config.G2Path, err)
		}

		g2Trailing, err = kzg.ReadG2PointSection(
			config.G2Path,
			config.SRSOrder-config.SRSNumberToLoad,
			config.SRSOrder, // last exclusive
			config.NumWorker,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to read trailing G2 points from %s: %v", config.G2Path, err)
		}
	} else {
		if len(config.G2PowerOf2Path) == 0 && len(config.G2Path) == 0 {
			return nil, errors.New("both G2Path and G2PowerOf2Path are empty. However, object needs to load G2Points")
		}

		if len(config.G2PowerOf2Path) != 0 {
			if config.SRSOrder == 0 {
				return nil, errors.New("SRS order cannot be 0")
			}

			maxPower := uint64(math.Log2(float64(config.SRSOrder)))
			_, err := kzg.ReadG2PointSection(config.G2PowerOf2Path, 0, maxPower, 1)
			if err != nil {
				return nil, fmt.Errorf("file located at %v is invalid", config.G2PowerOf2Path)
			}
		} else {
			log.Println("verifier requires accesses to entire g2 points. It is a legacy usage. For most operators, it is likely because G2_POWER_OF_2_PATH is improperly configured.")
		}
	}
	srs, err := kzg.NewSrs(s1, s2)
	if err != nil {
		return nil, fmt.Errorf("failed to create SRS: %v", err)
	}

	fmt.Println("numthread", runtime.GOMAXPROCS(0))

	encoderGroup := &Verifier{
		KzgConfig:             config,
		Srs:                   srs,
		G2Trailing:            g2Trailing,
		ParametrizedVerifiers: make(map[encoding.EncodingParams]*ParametrizedVerifier),
		LoadG2Points:          loadG2Points,
	}

	return encoderGroup, nil

}

type ParametrizedVerifier struct {
	*kzg.KzgConfig
	Srs *kzg.SRS

	*rs.Encoder

	Fs *fft.FFTSettings
	Ks *kzg.KZGSettings
}

func (g *Verifier) GetKzgVerifier(params encoding.EncodingParams) (*ParametrizedVerifier, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if err := params.Validate(); err != nil {
		return nil, err
	}

	ver, ok := g.ParametrizedVerifiers[params]
	if ok {
		return ver, nil
	}

	ver, err := g.newKzgVerifier(params)
	if err == nil {
		g.ParametrizedVerifiers[params] = ver
	}

	return ver, err
}

func (g *Verifier) NewKzgVerifier(params encoding.EncodingParams) (*ParametrizedVerifier, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.newKzgVerifier(params)
}

func (g *Verifier) newKzgVerifier(params encoding.EncodingParams) (*ParametrizedVerifier, error) {

	if err := params.Validate(); err != nil {
		return nil, err
	}

	n := uint8(math.Log2(float64(params.NumEvaluations())))
	fs := fft.NewFFTSettings(n)
	ks, err := kzg.NewKZGSettings(fs, g.Srs)

	if err != nil {
		return nil, err
	}

	encoder, err := rs.NewEncoder(params, g.Verbose)
	if err != nil {
		log.Println("Could not create encoder: ", err)
		return nil, err
	}

	return &ParametrizedVerifier{
		KzgConfig: g.KzgConfig,
		Srs:       g.Srs,
		Encoder:   encoder,
		Fs:        fs,
		Ks:        ks,
	}, nil
}

func (v *Verifier) VerifyBlobLength(commitments encoding.BlobCommitments) error {
	return v.VerifyCommit((*bn254.G2Affine)(commitments.LengthCommitment), (*bn254.G2Affine)(commitments.LengthProof), uint64(commitments.Length))

}

// VerifyCommit verifies the low degree proof; since it doesn't depend on the encoding parameters
// we leave it as a method of the KzgEncoderGroup
func (v *Verifier) VerifyCommit(lengthCommit *bn254.G2Affine, lengthProof *bn254.G2Affine, length uint64) error {

	g1Challenge, err := kzg.ReadG1Point(v.SRSOrder-length, v.KzgConfig)
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
		)

		if err != nil {
			return err
		}
	}

	return nil

}

func (v *ParametrizedVerifier) VerifyFrame(commit *bn254.G1Affine, f *encoding.Frame, index uint64) error {

	j, err := rs.GetLeadingCosetIndex(
		uint64(index),
		v.NumChunks,
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
	encoder, err := v.GetKzgVerifier(params)
	if err != nil {
		return nil, err
	}

	return encoder.Decode(frames, toUint64Array(indices), maxInputSize)
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
