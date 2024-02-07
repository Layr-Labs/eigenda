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
	bls "github.com/Layr-Labs/eigenda/pkg/kzg/bn254"
)

type KZGSettings struct {
	*FFTSettings

	Srs *SRS
	// setup values
}

func NewKZGSettings(fs *FFTSettings, srs *SRS) (*KZGSettings, error) {

	ks := &KZGSettings{
		FFTSettings: fs,
		Srs:         srs,
	}

	return ks, nil
}

type FK20SingleSettings struct {
	*KZGSettings
	xExtFFT []bls.G1Point
}

func NewFK20SingleSettings(ks *KZGSettings, n2 uint64) *FK20SingleSettings {
	if n2 > ks.MaxWidth {
		panic("extended size is larger than kzg settings supports")
	}
	if !bls.IsPowerOfTwo(n2) {
		panic("extended size is not a power of two")
	}
	if n2 < 2 {
		panic("extended size is too small")
	}
	n := n2 / 2
	fk := &FK20SingleSettings{
		KZGSettings: ks,
	}
	x := make([]bls.G1Point, n)
	for i, j := uint64(0), n-2; i < n-1; i, j = i+1, j-1 {
		bls.CopyG1(&x[i], &ks.Srs.G1[j])
	}
	bls.CopyG1(&x[n-1], &bls.ZeroG1)
	fk.xExtFFT = fk.toeplitzPart1(x)
	return fk
}

type FK20MultiSettings struct {
	*KZGSettings
	chunkLen uint64
	// chunkLen files, each of size MaxWidth
	xExtFFTFiles [][]bls.G1Point
}

func NewFK20MultiSettings(ks *KZGSettings, n2 uint64, chunkLen uint64) *FK20MultiSettings {
	if n2 > ks.MaxWidth {
		panic("extended size is larger than kzg settings supports")
	}
	if !bls.IsPowerOfTwo(n2) {
		panic("extended size is not a power of two")
	}
	if n2 < 2 {
		panic("extended size is too small")
	}
	if chunkLen > n2/2 {
		panic("chunk length is too large")
	}
	if !bls.IsPowerOfTwo(chunkLen) {
		panic("chunk length must be power of two")
	}
	if chunkLen < 1 {
		panic("chunk length is too small")
	}
	fk := &FK20MultiSettings{
		KZGSettings:  ks,
		chunkLen:     chunkLen,
		xExtFFTFiles: make([][]bls.G1Point, chunkLen),
	}
	// xext_fft = []
	// for i in range(l):
	//   x = setup[0][n - l - 1 - i::-l] + [b.Z1]
	//   xext_fft.append(toeplitz_part1(x))
	n := n2 / 2
	k := n / chunkLen
	xExtFFTPrecompute := func(offset uint64) []bls.G1Point {
		x := make([]bls.G1Point, k)
		start := n - chunkLen - 1 - offset
		for i, j := uint64(0), start; i+1 < k; i, j = i+1, j-chunkLen {
			bls.CopyG1(&x[i], &ks.Srs.G1[j])
		}
		bls.CopyG1(&x[k-1], &bls.ZeroG1)
		return ks.toeplitzPart1(x)
	}
	for i := uint64(0); i < chunkLen; i++ {
		fk.xExtFFTFiles[i] = xExtFFTPrecompute(i)
	}
	return fk
}
