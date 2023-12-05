// MIT License
//
// Copyright (c) 2020 @protolambda
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

//go:build !bignum_pure && !bignum_hol256
// +build !bignum_pure,!bignum_hol256

package kzg

import (
	"testing"

	bls "github.com/Layr-Labs/eigenda/pkg/kzg/bn254"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKZGSettings_DAUsingFK20Multi(t *testing.T) {
	fs := NewFFTSettings(4 + 5 + 1)
	chunkLen := uint64(16)
	chunkCount := uint64(32)
	n := chunkLen * chunkCount
	s1, s2 := GenerateTestingSetup("1927409816240961209460912649124", chunkLen*chunkCount*2)
	srs, _ := NewSrs(s1, s2)
	ks, _ := NewKZGSettings(fs, srs)
	fk := NewFK20MultiSettings(ks, n*2, chunkLen)

	// replicate same polynomial as in python test
	polynomial := make([]bls.Fr, n)
	var tmp134 bls.Fr
	bls.AsFr(&tmp134, 134)
	for i := uint64(0); i < chunkCount; i++ {
		// Note: different contents from older python test, make each section different,
		// to cover toeplitz coefficient edge cases better.
		for j, v := range []uint64{1, 2, 3, 4 + i, 7, 8 + i*i, 9, 10, 13, 14, 1, 15, 0, 1000, 0, 33} {
			bls.AsFr(&polynomial[i*chunkLen+uint64(j)], v)
		}
		bls.SubModFr(&polynomial[i*chunkLen+12], &bls.ZERO, &bls.ONE) // "MODULUS - 1"
		bls.SubModFr(&polynomial[i*chunkLen+14], &bls.ZERO, &tmp134)  // "MODULUS - 134"
	}

	commitment := ks.CommitToPoly(polynomial)

	allProofs, err := fk.DAUsingFK20Multi(polynomial)
	require.Nil(t, err, "could not compute proof")
	require.NotNil(t, allProofs)

	// We have the data in polynomial form already,
	// no need to use the DAS FFT (which extends data directly, not coeffs).
	extendedCoeffs := make([]bls.Fr, n*2)
	for i := uint64(0); i < n; i++ {
		bls.CopyFr(&extendedCoeffs[i], &polynomial[i])
	}
	for i := n; i < n*2; i++ {
		bls.CopyFr(&extendedCoeffs[i], &bls.ZERO)
	}
	extendedData, err := ks.FFT(extendedCoeffs, false)
	require.Nil(t, err)
	require.NotNil(t, extendedData)

	err = reverseBitOrderFr(extendedData)
	assert.Nil(t, err)

	n2 := n * 2
	domainStride := fk.MaxWidth / n2
	for pos := uint64(0); pos < 2*chunkCount; pos++ {
		domainPos := reverseBitsLimited(uint32(2*chunkCount), uint32(pos))
		var x bls.Fr
		bls.CopyFr(&x, &ks.ExpandedRootsOfUnity[uint64(domainPos)*domainStride])
		ys := extendedData[chunkLen*pos : chunkLen*(pos+1)]
		// ys, but constructed by evaluating the polynomial in the sub-domain range
		ys2 := make([]bls.Fr, chunkLen)
		// don't recompute the subgroup domain, just select it from the bigger domain by applying a stride
		stride := ks.MaxWidth / chunkLen
		coset := make([]bls.Fr, chunkLen)
		for i := uint64(0); i < chunkLen; i++ {
			var z bls.Fr // a value of the coset list
			bls.MulModFr(&z, &x, &ks.ExpandedRootsOfUnity[i*stride])
			bls.CopyFr(&coset[i], &z)
			bls.EvalPolyAt(&ys2[i], polynomial, &z)
		}
		// permanently change order of ys values
		err := reverseBitOrderFr(ys)
		assert.Nil(t, err)
		for i := 0; i < len(ys); i++ {
			assert.True(t, bls.EqualFr(&ys[i], &ys2[i]), "failed to reproduce matching y values for subgroup")
		}

		proof := &allProofs[pos]
		val, err := ks.CheckProofMulti(commitment, proof, &x, ys)
		require.Nil(t, err, "could not verify proof")
		assert.True(t, val, "could not verify proof")

		t.Logf("Data availability check %d passed", pos)
	}
}
