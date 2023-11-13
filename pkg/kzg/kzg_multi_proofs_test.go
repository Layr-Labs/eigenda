//go:build !bignum_pure && !bignum_hol256
// +build !bignum_pure,!bignum_hol256

package kzg

import (
	"testing"

	bls "github.com/Layr-Labs/eigenda/pkg/kzg/bn254"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKZGSettings_CheckProofMulti(t *testing.T) {
	fs := NewFFTSettings(4)
	s1, s2 := GenerateTestingSetup("1927409816240961209460912649124", 16+1)
	srs, _ := NewSrs(s1, s2)
	ks, _ := NewKZGSettings(fs, srs)
	polynomial := testPoly(1, 2, 3, 4, 7, 7, 7, 7, 13, 13, 13, 13, 13, 13, 13, 13)

	commitment := ks.CommitToPoly(polynomial)

	x := uint64(5431)
	var xFr bls.Fr
	bls.AsFr(&xFr, x)
	cosetScale := uint8(3)
	coset := make([]bls.Fr, 1<<cosetScale)
	s1, s2 = GenerateTestingSetup("1927409816240961209460912649124", 8+1)
	srs, _ = NewSrs(s1, s2)
	ks, _ = NewKZGSettings(NewFFTSettings(cosetScale), srs)
	for i := 0; i < len(coset); i++ {
		bls.MulModFr(&coset[i], &xFr, &ks.ExpandedRootsOfUnity[i])
	}
	ys := make([]bls.Fr, len(coset))
	for i := 0; i < len(coset); i++ {
		bls.EvalPolyAt(&ys[i], polynomial, &coset[i])
	}

	proof := ks.ComputeProofMulti(polynomial, x, uint64(len(coset)))
	valid, err := ks.CheckProofMulti(commitment, proof, &xFr, ys)
	require.Nil(t, err, "could not verify proof")
	assert.True(t, valid, "could not verify proof")
}
