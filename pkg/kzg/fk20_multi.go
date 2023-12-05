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

// Original: https://github.com/ethereum/research/blob/master/kzg_data_availability/fk20_multi.py

//go:build !bignum_pure && !bignum_hol256
// +build !bignum_pure,!bignum_hol256

package kzg

import (
	"fmt"

	bls "github.com/Layr-Labs/eigenda/pkg/kzg/bn254"
)

// FK20 Method to compute all proofs
// Toeplitz multiplication via http://www.netlib.org/utk/people/JackDongarra/etemplates/node384.html
// Multi proof method

// For a polynomial of size n, let w be a n-th root of unity. Then this method will return
// k=n/l KZG proofs for the points
//
//	proof[0]: w^(0*l + 0), w^(0*l + 1), ... w^(0*l + l - 1)
//	proof[0]: w^(0*l + 0), w^(0*l + 1), ... w^(0*l + l - 1)
//	...
//	proof[i]: w^(i*l + 0), w^(i*l + 1), ... w^(i*l + l - 1)
//	...
// func (ks *FK20MultiSettings) FK20Multi(polynomial []bls.Fr) []bls.G1Point {
// 	n := uint64(len(polynomial))
// 	n2 := n * 2
// 	if ks.MaxWidth < n2 {
// 		panic(fmt.Errorf("KZGSettings are set to MaxWidth %d but got half polynomial of length %d",
// 			ks.MaxWidth, n))
// 	}

// 	hExtFFT := make([]bls.G1Point, n2, n2)
// 	for i := uint64(0); i < n2; i++ {
// 		bls.CopyG1(&hExtFFT[i], &bls.ZeroG1)
// 	}

// 	var tmp bls.G1Point
// 	for i := uint64(0); i < ks.chunkLen; i++ {
// 		toeplitzCoeffs := ks.toeplitzCoeffsStepStrided(polynomial, i, ks.chunkLen)
// 		hExtFFTFile := ks.ToeplitzPart2(toeplitzCoeffs, ks.xExtFFTFiles[i])
// 		for j := uint64(0); j < n2; j++ {
// 			bls.AddG1(&tmp, &hExtFFT[j], &hExtFFTFile[j])
// 			bls.CopyG1(&hExtFFT[j], &tmp)
// 		}
// 	}
// 	h := ks.ToeplitzPart3(hExtFFT)

// 	out, err := ks.FFTG1(h, false)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return out
// }

// FK20MultiDAOptimized is FK20 multi-proof method, optimized for dava availability where the top half of polynomial
// coefficients == 0
func (ks *FK20MultiSettings) FK20MultiDAOptimized(polynomial []bls.Fr) []bls.G1Point {
	n2 := uint64(len(polynomial))
	if ks.MaxWidth < n2 {
		panic(fmt.Errorf("KZGSettings are set to MaxWidth %d but got polynomial of length %d",
			ks.MaxWidth, n2))
	}
	n := n2 / 2
	for i := n; i < n2; i++ {
		if !bls.EqualZero(&polynomial[i]) {
			panic("bad input, second half should be zeroed")
		}
	}

	k := n / ks.chunkLen
	k2 := k * 2
	hExtFFT := make([]bls.G1Point, k2)
	for i := uint64(0); i < k2; i++ {
		bls.CopyG1(&hExtFFT[i], &bls.ZeroG1)
	}

	reducedPoly := polynomial[:n]
	var tmp bls.G1Point
	for i := uint64(0); i < ks.chunkLen; i++ {
		toeplitzCoeffs := ks.toeplitzCoeffsStepStrided(reducedPoly, i, ks.chunkLen)
		hExtFFTFile := ks.ToeplitzPart2(toeplitzCoeffs, ks.xExtFFTFiles[i])
		for j := uint64(0); j < k2; j++ {
			bls.AddG1(&tmp, &hExtFFT[j], &hExtFFTFile[j])
			bls.CopyG1(&hExtFFT[j], &tmp)
		}
	}
	h := ks.ToeplitzPart3(hExtFFT)

	// TODO: maybe use a G1 version of the DAS extension FFT to perform the h -> output conversion?

	// Now redo the padding before final step.
	// Instead of copying h into a new extended array, just reuse the old capacity.
	h = h[:k2]
	for i := k; i < k2; i++ {
		bls.CopyG1(&h[i], &bls.ZeroG1)
	}
	out, err := ks.FFTG1(h, false)
	if err != nil {
		panic(err)
	}
	return out
}

// DAUsingFK20Multi computes all the KZG proofs for data availability checks. This involves sampling on the double domain
// and reordering according to reverse bit order
func (ks *FK20MultiSettings) DAUsingFK20Multi(polynomial []bls.Fr) ([]bls.G1Point, error) {
	n := uint64(len(polynomial))
	if n > ks.MaxWidth/2 {
		return nil, ErrInvalidPolyLengthTooLarge
	}
	if !bls.IsPowerOfTwo(n) {
		return nil, ErrInvalidPolyLengthPowerOfTwo
	}
	n2 := n * 2
	extendedPolynomial := make([]bls.Fr, n2)
	for i := uint64(0); i < n; i++ {
		bls.CopyFr(&extendedPolynomial[i], &polynomial[i])
	}
	for i := n; i < n2; i++ {
		bls.CopyFr(&extendedPolynomial[i], &bls.ZERO)
	}
	allProofs := ks.FK20MultiDAOptimized(extendedPolynomial)
	// change to reverse bit order.
	err := reverseBitOrderG1(allProofs)
	if err != nil {
		return nil, err
	}
	return allProofs, nil
}
