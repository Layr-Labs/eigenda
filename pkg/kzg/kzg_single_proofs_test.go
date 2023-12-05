// This code is sourced from the go-kzg Repository by protolambda.
// Original code: https://github.com/protolambda/go-kzg
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

func TestKZGSettings_CommitToEvalPoly(t *testing.T) {
	fs := NewFFTSettings(4)
	s1, s2 := GenerateTestingSetup("1927409816240961209460912649124", 16+1)
	srs, _ := NewSrs(s1, s2)
	ks, _ := NewKZGSettings(fs, srs)
	polynomial := testPoly(1, 2, 3, 4, 7, 7, 7, 7, 13, 13, 13, 13, 13, 13, 13, 13)
	evalPoly, err := fs.FFT(polynomial, false)
	require.Nil(t, err)

	secretG1IFFT, err := fs.FFTG1(ks.Srs.G1[:16], true)
	require.Nil(t, err)

	commitmentByCoeffs := ks.CommitToPoly(polynomial)
	commitmentByEval := CommitToEvalPoly(secretG1IFFT, evalPoly)
	assert.True(t, bls.EqualG1(commitmentByEval, commitmentByCoeffs),
		"expected commitments to be equal, but got:\nby eval: %s\nby coeffs: %s", commitmentByEval, commitmentByCoeffs)
}

func TestKZGSettings_CheckProofSingle(t *testing.T) {
	fs := NewFFTSettings(4)
	s1, s2 := GenerateTestingSetup("1927409816240961209460912649124", 16+1)
	srs, _ := NewSrs(s1, s2)
	ks, _ := NewKZGSettings(fs, srs)
	polynomial := testPoly(1, 2, 3, 4, 7, 7, 7, 7, 13, 13, 13, 13, 13, 13, 13, 13)
	commitment := ks.CommitToPoly(polynomial)
	proof := ks.ComputeProofSingle(polynomial, 17)

	var x bls.Fr
	bls.AsFr(&x, 17)
	var value bls.Fr
	bls.EvalPolyAt(&value, polynomial, &x)

	assert.True(t, ks.CheckProofSingle(commitment, proof, &x, &value), "could not verify proof")
}

func testPoly(polynomial ...uint64) []bls.Fr {
	n := len(polynomial)
	polynomialFr := make([]bls.Fr, n)
	for i := 0; i < n; i++ {
		bls.AsFr(&polynomialFr[i], polynomial[i])
	}
	return polynomialFr
}
