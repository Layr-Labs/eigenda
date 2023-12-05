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
	"fmt"

	bls "github.com/Layr-Labs/eigenda/pkg/kzg/bn254"
)

func (fs *FFTSettings) simpleFTG1(vals []bls.G1Point, valsOffset uint64, valsStride uint64, rootsOfUnity []bls.Fr, rootsOfUnityStride uint64, out []bls.G1Point) {
	l := uint64(len(out))
	var v bls.G1Point
	var tmp bls.G1Point
	var last bls.G1Point
	for i := uint64(0); i < l; i++ {
		jv := &vals[valsOffset]
		r := &rootsOfUnity[0]
		bls.MulG1(&v, jv, r)
		bls.CopyG1(&last, &v)

		for j := uint64(1); j < l; j++ {
			jv := &vals[valsOffset+j*valsStride]
			r := &rootsOfUnity[((i*j)%l)*rootsOfUnityStride]
			bls.MulG1(&v, jv, r)
			bls.CopyG1(&tmp, &last)
			bls.AddG1(&last, &tmp, &v)
		}
		bls.CopyG1(&out[i], &last)
	}
}

func (fs *FFTSettings) _fftG1(vals []bls.G1Point, valsOffset uint64, valsStride uint64, rootsOfUnity []bls.Fr, rootsOfUnityStride uint64, out []bls.G1Point) {
	if len(out) <= 4 { // if the value count is small, run the unoptimized version instead. // TODO tune threshold. (can be different for G1)
		fs.simpleFTG1(vals, valsOffset, valsStride, rootsOfUnity, rootsOfUnityStride, out)
		return
	}

	half := uint64(len(out)) >> 1
	// L will be the left half of out
	fs._fftG1(vals, valsOffset, valsStride<<1, rootsOfUnity, rootsOfUnityStride<<1, out[:half])
	// R will be the right half of out
	fs._fftG1(vals, valsOffset+valsStride, valsStride<<1, rootsOfUnity, rootsOfUnityStride<<1, out[half:]) // just take even again

	var yTimesRoot bls.G1Point
	var x, y bls.G1Point
	for i := uint64(0); i < half; i++ {
		// temporary copies, so that writing to output doesn't conflict with input
		bls.CopyG1(&x, &out[i])
		bls.CopyG1(&y, &out[i+half])
		root := &rootsOfUnity[i*rootsOfUnityStride]
		bls.MulG1(&yTimesRoot, &y, root)
		bls.AddG1(&out[i], &x, &yTimesRoot)
		bls.SubG1(&out[i+half], &x, &yTimesRoot)
	}
}

func (fs *FFTSettings) FFTG1(vals []bls.G1Point, inv bool) ([]bls.G1Point, error) {
	n := uint64(len(vals))
	if n > fs.MaxWidth {
		return nil, fmt.Errorf("got %d values but only have %d roots of unity", n, fs.MaxWidth)
	}
	if !bls.IsPowerOfTwo(n) {
		return nil, fmt.Errorf("got %d values but not a power of two", n)
	}
	// We make a copy so we can mutate it during the work.
	valsCopy := make([]bls.G1Point, n)
	for i := 0; i < len(vals); i++ { // TODO: maybe optimize this away, and write back to original input array?
		bls.CopyG1(&valsCopy[i], &vals[i])
	}
	if inv {
		var invLen bls.Fr
		bls.AsFr(&invLen, n)
		bls.InvModFr(&invLen, &invLen)
		rootz := fs.ReverseRootsOfUnity[:fs.MaxWidth]
		stride := fs.MaxWidth / n

		out := make([]bls.G1Point, n)
		fs._fftG1(valsCopy, 0, 1, rootz, stride, out)
		var tmp bls.G1Point
		for i := 0; i < len(out); i++ {
			bls.MulG1(&tmp, &out[i], &invLen)
			bls.CopyG1(&out[i], &tmp)
		}
		return out, nil
	} else {
		out := make([]bls.G1Point, n)
		rootz := fs.ExpandedRootsOfUnity[:fs.MaxWidth]
		stride := fs.MaxWidth / n
		// Regular FFT
		fs._fftG1(valsCopy, 0, 1, rootz, stride, out)
		return out, nil
	}
}

// rearrange G1 elements in reverse bit order. Supports 2**31 max element count.
func reverseBitOrderG1(values []bls.G1Point) error {
	if len(values) > (1 << 31) {
		return ErrG1ListTooLarge
	}
	var tmp bls.G1Point
	reverseBitOrder(uint32(len(values)), func(i, j uint32) {
		bls.CopyG1(&tmp, &values[i])
		bls.CopyG1(&values[i], &values[j])
		bls.CopyG1(&values[j], &tmp)
	})
	return nil
}
