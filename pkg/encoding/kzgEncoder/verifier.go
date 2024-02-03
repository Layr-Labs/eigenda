package kzgEncoder

import (
	"errors"
	"math"

	rs "github.com/Layr-Labs/eigenda/pkg/encoding/encoder"
	kzg "github.com/Layr-Labs/eigenda/pkg/kzg"
	wbls "github.com/Layr-Labs/eigenda/pkg/kzg/bn254"
)

type KzgVerifier struct {
	*KzgConfig
	Srs *kzg.SRS

	rs.EncodingParams

	Fs *kzg.FFTSettings
	Ks *kzg.KZGSettings
}

func (g *KzgEncoderGroup) GetKzgVerifier(params rs.EncodingParams) (*KzgVerifier, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if err := params.Validate(); err != nil {
		return nil, err
	}

	ver, ok := g.Verifiers[params]
	if ok {
		return ver, nil
	}

	ver, err := g.newKzgVerifier(params)
	if err == nil {
		g.Verifiers[params] = ver
	}

	return ver, err
}

func (g *KzgEncoderGroup) NewKzgVerifier(params rs.EncodingParams) (*KzgVerifier, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.newKzgVerifier(params)
}

func (g *KzgEncoderGroup) newKzgVerifier(params rs.EncodingParams) (*KzgVerifier, error) {

	if err := params.Validate(); err != nil {
		return nil, err
	}

	n := uint8(math.Log2(float64(params.NumEvaluations())))
	fs := kzg.NewFFTSettings(n)
	ks, err := kzg.NewKZGSettings(fs, g.Srs)

	if err != nil {
		return nil, err
	}

	return &KzgVerifier{
		KzgConfig:      g.KzgConfig,
		Srs:            g.Srs,
		EncodingParams: params,
		Fs:             fs,
		Ks:             ks,
	}, nil
}

// VerifyCommit verifies the low degree proof; since it doesn't depend on the encoding parameters
// we leave it as a method of the KzgEncoderGroup
func (v *KzgEncoderGroup) VerifyCommit(lengthCommit *wbls.G2Point, lowDegreeProof *wbls.G2Point, degree uint64) error {

	if !VerifyLowDegreeProof(lengthCommit, lowDegreeProof, degree, v.SRSOrder, v.Srs.G1) {
		return errors.New("low degree proof fails")
	}
	return nil

}

func (v *KzgVerifier) VerifyFrame(commit *wbls.G1Point, f *Frame, index uint64) error {

	j, err := rs.GetLeadingCosetIndex(
		uint64(index),
		v.NumChunks,
	)
	if err != nil {
		return err
	}

	if !f.Verify(v.Ks, commit, &v.Ks.ExpandedRootsOfUnity[j]) {
		return errors.New("multireveal proof fails")
	}

	return nil

}
