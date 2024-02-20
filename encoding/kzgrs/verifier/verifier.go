package verifier

import (
	"errors"
	"fmt"
	"log"
	"math"
	"runtime"
	"sync"

	enc "github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/kzgrs"
	"github.com/Layr-Labs/eigenda/encoding/rs"
	kzg "github.com/Layr-Labs/eigenda/pkg/kzg"
	bls "github.com/Layr-Labs/eigenda/pkg/kzg/bn254"
	wbls "github.com/Layr-Labs/eigenda/pkg/kzg/bn254"
)

type Verifier struct {
	*kzgrs.KzgConfig
	Srs          *kzg.SRS
	G2Trailing   []bls.G2Point
	mu           sync.Mutex
	LoadG2Points bool

	ParametrizedVerifiers map[rs.EncodingParams]*ParametrizedVerifier
}

func NewVerifier(config *kzgrs.KzgConfig, loadG2Points bool) (*Verifier, error) {

	if config.SRSNumberToLoad > config.SRSOrder {
		return nil, errors.New("SRSOrder is less than srsNumberToLoad")
	}

	// read the whole order, and treat it as entire SRS for low degree proof
	s1, err := kzgrs.ReadG1Points(config.G1Path, config.SRSNumberToLoad, config.NumWorker)
	if err != nil {
		log.Println("failed to read G1 points", err)
		return nil, err
	}

	s2 := make([]bls.G2Point, 0)
	g2Trailing := make([]bls.G2Point, 0)

	// PreloadEncoder is by default not used by operator node, PreloadEncoder
	if loadG2Points {
		if len(config.G2Path) == 0 {
			return nil, fmt.Errorf("G2Path is empty. However, object needs to load G2Points")
		}

		s2, err = kzgrs.ReadG2Points(config.G2Path, config.SRSNumberToLoad, config.NumWorker)
		if err != nil {
			log.Println("failed to read G2 points", err)
			return nil, err
		}

		g2Trailing, err = kzgrs.ReadG2PointSection(
			config.G2Path,
			config.SRSOrder-config.SRSNumberToLoad,
			config.SRSOrder, // last exclusive
			config.NumWorker,
		)
		if err != nil {
			return nil, err
		}
	} else {
		// todo, there are better ways to handle it
		if len(config.G2PowerOf2Path) == 0 {
			return nil, fmt.Errorf("G2PowerOf2Path is empty. However, object needs to load G2Points")
		}
	}

	srs, err := kzg.NewSrs(s1, s2)
	if err != nil {
		log.Println("Could not create srs", err)
		return nil, err
	}

	fmt.Println("numthread", runtime.GOMAXPROCS(0))

	encoderGroup := &Verifier{
		KzgConfig:             config,
		Srs:                   srs,
		G2Trailing:            g2Trailing,
		ParametrizedVerifiers: make(map[rs.EncodingParams]*ParametrizedVerifier),
		LoadG2Points:          loadG2Points,
	}

	return encoderGroup, nil

}

type ParametrizedVerifier struct {
	*kzgrs.KzgConfig
	Srs *kzg.SRS

	rs.EncodingParams

	Fs *kzg.FFTSettings
	Ks *kzg.KZGSettings
}

func (g *Verifier) GetKzgVerifier(params rs.EncodingParams) (*ParametrizedVerifier, error) {
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

func (g *Verifier) NewKzgVerifier(params rs.EncodingParams) (*ParametrizedVerifier, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.newKzgVerifier(params)
}

func (g *Verifier) newKzgVerifier(params rs.EncodingParams) (*ParametrizedVerifier, error) {

	if err := params.Validate(); err != nil {
		return nil, err
	}

	n := uint8(math.Log2(float64(params.NumEvaluations())))
	fs := kzg.NewFFTSettings(n)
	ks, err := kzg.NewKZGSettings(fs, g.Srs)

	if err != nil {
		return nil, err
	}

	return &ParametrizedVerifier{
		KzgConfig:      g.KzgConfig,
		Srs:            g.Srs,
		EncodingParams: params,
		Fs:             fs,
		Ks:             ks,
	}, nil
}

// VerifyCommit verifies the low degree proof; since it doesn't depend on the encoding parameters
// we leave it as a method of the KzgEncoderGroup
func (v *Verifier) VerifyCommit(lengthCommit *wbls.G2Point, lowDegreeProof *wbls.G2Point, length uint64) error {

	g1Challenge, err := kzgrs.ReadG1Point(v.SRSOrder-length, v.KzgConfig)
	if err != nil {
		return err
	}

	if !VerifyLowDegreeProof(lengthCommit, lowDegreeProof, &g1Challenge) {
		return errors.New("low degree proof fails")
	}
	return nil

}

// The function verify low degree proof against a poly commitment
// We wish to show x^shift poly = shiftedPoly, with
// With shift = SRSOrder-1 - claimedDegree and
// proof = commit(shiftedPoly) on G1
// so we can verify by checking
// e( commit_1, [x^shift]_2) = e( proof_1, G_2 )
func VerifyLowDegreeProof(lengthCommit *bls.G2Point, proof *bls.G2Point, g1Challenge *bls.G1Point) bool {
	return bls.PairingsVerify(g1Challenge, lengthCommit, &bls.GenG1, proof)
}

func (v *ParametrizedVerifier) VerifyFrame(commit *wbls.G1Point, f *enc.Frame, index uint64) error {

	j, err := rs.GetLeadingCosetIndex(
		uint64(index),
		v.NumChunks,
	)
	if err != nil {
		return err
	}

	g2Atn, err := kzgrs.ReadG2Point(uint64(len(f.Coeffs)), v.KzgConfig)
	if err != nil {
		return err
	}

	if !VerifyFrame(f, v.Ks, commit, &v.Ks.ExpandedRootsOfUnity[j], &g2Atn) {
		return errors.New("multireveal proof fails")
	}

	return nil

}

// Verify function assumes the Data stored is coefficients of coset's interpolating poly
func VerifyFrame(f *enc.Frame, ks *kzg.KZGSettings, commitment *bls.G1Point, x *bls.Fr, g2Atn *bls.G2Point) bool {
	var xPow bls.Fr
	bls.CopyFr(&xPow, &bls.ONE)

	var tmp bls.Fr
	for i := 0; i < len(f.Coeffs); i++ {
		bls.MulModFr(&tmp, &xPow, x)
		bls.CopyFr(&xPow, &tmp)
	}

	// [x^n]_2
	var xn2 bls.G2Point
	bls.MulG2(&xn2, &bls.GenG2, &xPow)

	// [s^n - x^n]_2
	var xnMinusYn bls.G2Point

	bls.SubG2(&xnMinusYn, g2Atn, &xn2)

	// [interpolation_polynomial(s)]_1
	is1 := bls.LinCombG1(ks.Srs.G1[:len(f.Coeffs)], f.Coeffs)

	// [commitment - interpolation_polynomial(s)]_1 = [commit]_1 - [interpolation_polynomial(s)]_1
	var commitMinusInterpolation bls.G1Point
	bls.SubG1(&commitMinusInterpolation, commitment, is1)

	// Verify the pairing equation
	//
	// e([commitment - interpolation_polynomial(s)], [1]) = e([proof],  [s^n - x^n])
	//    equivalent to
	// e([commitment - interpolation_polynomial]^(-1), [1]) * e([proof],  [s^n - x^n]) = 1_T
	//

	return bls.PairingsVerify(&commitMinusInterpolation, &bls.GenG2, &f.Proof, &xnMinusYn)
}
