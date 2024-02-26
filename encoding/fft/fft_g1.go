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

package fft

import (
	"fmt"
	"math/big"

	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

func (fs *FFTSettings) simpleFTG1(vals []bn254.G1Affine, valsOffset uint64, valsStride uint64, rootsOfUnity []fr.Element, rootsOfUnityStride uint64, out []bn254.G1Affine) {
	l := uint64(len(out))
	var v bn254.G1Affine
	var tmp bn254.G1Affine
	var last bn254.G1Affine
	for i := uint64(0); i < l; i++ {
		jv := &vals[valsOffset]
		r := &rootsOfUnity[0]

		var t big.Int
		r.BigInt(&t)
		v.ScalarMultiplication(jv, &t)

		last.Set(&v)

		for j := uint64(1); j < l; j++ {
			jv := &vals[valsOffset+j*valsStride]
			r := &rootsOfUnity[((i*j)%l)*rootsOfUnityStride]

			var t big.Int
			r.BigInt(&t)
			v.ScalarMultiplication(jv, &t)
			tmp.Set(&last)
			last.Add(&tmp, &v)
		}
		out[i].Set(&last)

	}
}

func (fs *FFTSettings) _fftG1(vals []bn254.G1Affine, valsOffset uint64, valsStride uint64, rootsOfUnity []fr.Element, rootsOfUnityStride uint64, out []bn254.G1Affine) {
	if len(out) <= 4 { // if the value count is small, run the unoptimized version instead. // TODO tune threshold. (can be different for G1)
		fs.simpleFTG1(vals, valsOffset, valsStride, rootsOfUnity, rootsOfUnityStride, out)
		return
	}

	half := uint64(len(out)) >> 1
	// L will be the left half of out
	fs._fftG1(vals, valsOffset, valsStride<<1, rootsOfUnity, rootsOfUnityStride<<1, out[:half])
	// R will be the right half of out
	fs._fftG1(vals, valsOffset+valsStride, valsStride<<1, rootsOfUnity, rootsOfUnityStride<<1, out[half:]) // just take even again

	var yTimesRoot bn254.G1Affine
	var x, y bn254.G1Affine
	for i := uint64(0); i < half; i++ {
		// temporary copies, so that writing to output doesn't conflict with input
		x.Set(&out[i])
		y.Set(&out[i+half])

		root := &rootsOfUnity[i*rootsOfUnityStride]

		yTimesRoot.ScalarMultiplication(&y, root.BigInt(new(big.Int)))

		out[i].Add(&x, &yTimesRoot)
		out[i+half].Sub(&x, &yTimesRoot)

	}
}

func (fs *FFTSettings) FFTG1(vals []bn254.G1Affine, inv bool) ([]bn254.G1Affine, error) {
	n := uint64(len(vals))
	if n > fs.MaxWidth {
		return nil, fmt.Errorf("got %d values but only have %d roots of unity", n, fs.MaxWidth)
	}

	if !IsPowerOfTwo(n) {
		return nil, fmt.Errorf("got %d values but not a power of two", n)
	}
	// We make a copy so we can mutate it during the work.
	valsCopy := make([]bn254.G1Affine, n)
	for i := 0; i < len(vals); i++ { // TODO: maybe optimize this away, and write back to original input array?

		valsCopy[i].Set(&vals[i])
	}
	if inv {
		var invLen fr.Element

		invLen.SetUint64(n)

		invLen.Inverse(&invLen)

		rootz := fs.ReverseRootsOfUnity[:fs.MaxWidth]
		stride := fs.MaxWidth / n

		out := make([]bn254.G1Affine, n)
		fs._fftG1(valsCopy, 0, 1, rootz, stride, out)

		for i := 0; i < len(out); i++ {
			out[i].ScalarMultiplication(&out[i], invLen.BigInt(new(big.Int)))
		}
		return out, nil
	} else {
		out := make([]bn254.G1Affine, n)
		rootz := fs.ExpandedRootsOfUnity[:fs.MaxWidth]
		stride := fs.MaxWidth / n
		// Regular FFT
		fs._fftG1(valsCopy, 0, 1, rootz, stride, out)
		return out, nil
	}
}
