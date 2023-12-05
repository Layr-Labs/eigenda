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

// Original: https://github.com/ethereum/research/blob/master/kzg_data_availability/fk20_single.py

//go:build !bignum_pure && !bignum_hol256
// +build !bignum_pure,!bignum_hol256

package kzg

import (
	"fmt"

	bls "github.com/Layr-Labs/eigenda/pkg/kzg/bn254"
)

// FK20 Method to compute all proofs
// Toeplitz multiplication via http://www.netlib.org/utk/people/JackDongarra/etemplates/node384.html
// Single proof method

// A Toeplitz matrix is of the form
//
// t_0     t_(-1) t_(-2) ... t_(1-n)
// t_1     t_0    t_(-1) ... t_(2-n)
// t_2     t_1               .
// .              .          .
// .                 .       .
// .                    .    t(-1)
// t_(n-1)   ...       t_1   t_0
//
// The vector [t_0, t_1, ..., t_(n-2), t_(n-1), 0, t_(1-n), t_(2-n), ..., t_(-2), t_(-1)]
// completely determines the Toeplitz matrix and is called the "toeplitz_coefficients" below

// The composition toeplitz_part3(toeplitz_part2(toeplitz_coefficients, toeplitz_part1(x)))
// compute the matrix-vector multiplication T * x
//
// The algorithm here is written under the assumption x = G1 elements, T scalars
//
// For clarity, vectors in "Fourier space" are written with _fft. So for example, the vector
// xext is the extended x vector (padded with zero), and xext_fft is its Fourier transform.

// Performs the first part of the Toeplitz matrix multiplication algorithm, which is a Fourier
// transform of the vector x extended
func (ks *KZGSettings) toeplitzPart1(x []bls.G1Point) []bls.G1Point {
	n := uint64(len(x))
	n2 := n * 2
	// Extend x with zeros (neutral element of G1)
	xExt := make([]bls.G1Point, n2)
	for i := uint64(0); i < n; i++ {
		bls.CopyG1(&xExt[i], &x[i])
	}
	for i := n; i < n2; i++ {
		bls.CopyG1(&xExt[i], &bls.ZeroG1)
	}
	xExtFFT, err := ks.FFTG1(xExt, false)
	if err != nil {
		panic(fmt.Errorf("FFT G1 failed in toeplitz part 1: %v", err))
	}
	return xExtFFT
}

// Performs the second part of the Toeplitz matrix multiplication algorithm
func (ks *KZGSettings) ToeplitzPart2(toeplitzCoeffs []bls.Fr, xExtFFT []bls.G1Point) (hExtFFT []bls.G1Point) {
	if uint64(len(toeplitzCoeffs)) != uint64(len(xExtFFT)) {
		panic("expected toeplitz coeffs to match xExtFFT length")
	}
	toeplitzCoeffsFFT, err := ks.FFT(toeplitzCoeffs, false)
	if err != nil {
		panic(fmt.Errorf("FFT failed in toeplitz part 2: %v", err))
	}
	n := uint64(len(toeplitzCoeffsFFT))
	hExtFFT = make([]bls.G1Point, n)
	for i := uint64(0); i < n; i++ {
		bls.MulG1(&hExtFFT[i], &xExtFFT[i], &toeplitzCoeffsFFT[i])
	}
	return hExtFFT
}

// Transform back and return the first half of the vector
func (ks *KZGSettings) ToeplitzPart3(hExtFFT []bls.G1Point) []bls.G1Point {
	out, err := ks.FFTG1(hExtFFT, true)
	if err != nil {
		panic(fmt.Errorf("toeplitz part 3 err: %v", err))
	}
	// Only the top half is the Toeplitz product, the rest is padding
	return out[:len(out)/2]
}

func (ks *KZGSettings) toeplitzCoeffsStepStrided(polynomial []bls.Fr, offset uint64, stride uint64) []bls.Fr {
	n := uint64(len(polynomial))
	k := n / stride
	k2 := k * 2
	// [last poly item] + [0]*(n+1) + [poly items except first and last]
	toeplitzCoeffs := make([]bls.Fr, k2)
	bls.CopyFr(&toeplitzCoeffs[0], &polynomial[n-1-offset])
	for i := uint64(1); i <= k+1; i++ {
		bls.CopyFr(&toeplitzCoeffs[i], &bls.ZERO)
	}
	for i, j := k+2, 2*stride-offset-1; i < k2; i, j = i+1, j+stride {
		bls.CopyFr(&toeplitzCoeffs[i], &polynomial[j])
	}
	return toeplitzCoeffs
}

// TODO: call above with offset=0, stride=1
func (ks *KZGSettings) toeplitzCoeffsStep(polynomial []bls.Fr) []bls.Fr {
	n := uint64(len(polynomial))
	n2 := n * 2
	// [last poly item] + [0]*(n+1) + [poly items except first and last]
	toeplitzCoeffs := make([]bls.Fr, n2)
	bls.CopyFr(&toeplitzCoeffs[0], &polynomial[n-1])
	for i := uint64(1); i <= n+1; i++ {
		bls.CopyFr(&toeplitzCoeffs[i], &bls.ZERO)
	}
	for i, j := n+2, 1; i < n2; i, j = i+1, j+1 {
		bls.CopyFr(&toeplitzCoeffs[i], &polynomial[j])
	}
	return toeplitzCoeffs
}

// Compute all n (single) proofs according to FK20 method
// func (fk *FK20SingleSettings) FK20Single(polynomial []bls.Fr) []bls.G1Point {
// 	toeplitzCoeffs := fk.toeplitzCoeffsStep(polynomial)
// 	// Compute the vector h from the paper using a Toeplitz matrix multiplication
// 	hExtFFT := fk.ToeplitzPart2(toeplitzCoeffs, fk.xExtFFT)
// 	h := fk.ToeplitzPart3(hExtFFT)

// 	// TODO: correct? It will pad up implicitly again, but
// 	out, err := fk.FFTG1(h, false)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return out
// }

// Special version of the FK20 for the situation of data availability checks:
// The upper half of the polynomial coefficients is always 0, so we do not need to extend to twice the size
// for Toeplitz matrix multiplication
func (fk *FK20SingleSettings) FK20SingleDAOptimized(polynomial []bls.Fr) []bls.G1Point {
	if uint64(len(polynomial)) > fk.MaxWidth {
		panic(fmt.Errorf(
			"expected input of length %d (incl half of zeroes) to not exceed precomputed settings length %d",
			len(polynomial), fk.MaxWidth))
	}
	n2 := uint64(len(polynomial))
	if !bls.IsPowerOfTwo(n2) {
		panic(fmt.Errorf("expected input length to be power of two, got %d", n2))
	}
	n := n2 / 2
	for i := n; i < n2; i++ {
		if !bls.EqualZero(&polynomial[i]) {
			panic("bad input, second half should be zeroed")
		}
	}
	reducedPoly := polynomial[:n]
	toeplitzCoeffs := fk.toeplitzCoeffsStep(reducedPoly)
	// Compute the vector h from the paper using a Toeplitz matrix multiplication
	hExtFFT := fk.ToeplitzPart2(toeplitzCoeffs, fk.xExtFFT)
	h := fk.ToeplitzPart3(hExtFFT)

	// Now redo the padding before final step.
	// Instead of copying h into a new extended array, just reuse the old capacity.
	h = h[:n2]
	for i := n; i < n2; i++ {
		bls.CopyG1(&h[i], &bls.ZeroG1)
	}
	out, err := fk.FFTG1(h, false)
	if err != nil {
		panic(err)
	}
	return out
}

// Computes all the KZG proofs for data availability checks. This involves sampling on the double domain
// and reordering according to reverse bit order
func (fk *FK20SingleSettings) DAUsingFK20(polynomial []bls.Fr) ([]bls.G1Point, error) {
	n := uint64(len(polynomial))
	if n > fk.MaxWidth/2 {
		panic("expected poly contents not bigger than half the size of the FK20-single settings")
	}
	if !bls.IsPowerOfTwo(n) {
		panic("expected poly length to be power of two")
	}
	n2 := n * 2
	extendedPolynomial := make([]bls.Fr, n2)
	for i := uint64(0); i < n; i++ {
		bls.CopyFr(&extendedPolynomial[i], &polynomial[i])
	}
	for i := n; i < n2; i++ {
		bls.CopyFr(&extendedPolynomial[i], &bls.ZERO)
	}
	allProofs := fk.FK20SingleDAOptimized(extendedPolynomial)
	// change to reverse bit order.
	err := reverseBitOrderG1(allProofs)
	return allProofs, err
}
