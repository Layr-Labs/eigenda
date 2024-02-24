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
package fft

import (
	"fmt"

	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

func (fs *FFTSettings) simpleFT(vals []fr.Element, valsOffset uint64, valsStride uint64, rootsOfUnity []fr.Element, rootsOfUnityStride uint64, out []fr.Element) {
	l := uint64(len(out))
	var v fr.Element
	var tmp fr.Element
	var last fr.Element
	for i := uint64(0); i < l; i++ {
		jv := &vals[valsOffset]
		r := &rootsOfUnity[0]
		v.Mul(jv, r)
		last.Set(&v)

		for j := uint64(1); j < l; j++ {
			jv := &vals[valsOffset+j*valsStride]
			r := &rootsOfUnity[((i*j)%l)*rootsOfUnityStride]
			v.Mul(jv, r)
			tmp.Set(&last)
			last.Add(&tmp, &v)

			//bls.MulModFr(&v, jv, r)
			//bls.CopyFr(&tmp, &last)
			//bls.AddModFr(&last, &tmp, &v)
		}
		out[i].Set(&last)
		//bls.CopyFr(&out[i], &last)
	}
}

func (fs *FFTSettings) _fft(vals []fr.Element, valsOffset uint64, valsStride uint64, rootsOfUnity []fr.Element, rootsOfUnityStride uint64, out []fr.Element) {
	if len(out) <= 4 { // if the value count is small, run the unoptimized version instead. // TODO tune threshold.
		fs.simpleFT(vals, valsOffset, valsStride, rootsOfUnity, rootsOfUnityStride, out)
		return
	}

	half := uint64(len(out)) >> 1
	// L will be the left half of out
	fs._fft(vals, valsOffset, valsStride<<1, rootsOfUnity, rootsOfUnityStride<<1, out[:half])
	// R will be the right half of out
	fs._fft(vals, valsOffset+valsStride, valsStride<<1, rootsOfUnity, rootsOfUnityStride<<1, out[half:]) // just take even again

	var yTimesRoot fr.Element
	var x, y fr.Element
	for i := uint64(0); i < half; i++ {
		// temporary copies, so that writing to output doesn't conflict with input
		x.Set(&out[i])
		y.Set(&out[i+half])

		//bls.CopyFr(&x, &out[i])
		//bls.CopyFr(&y, &out[i+half])
		root := &rootsOfUnity[i*rootsOfUnityStride]
		yTimesRoot.Mul(&y, root)
		out[i].Add(&x, &yTimesRoot)
		out[i+half].Sub(&x, &yTimesRoot)

		//bls.MulModFr(&yTimesRoot, &y, root)
		//bls.AddModFr(&out[i], &x, &yTimesRoot)
		//bls.SubModFr(&out[i+half], &x, &yTimesRoot)
	}
}

func (fs *FFTSettings) FFT(vals []fr.Element, inv bool) ([]fr.Element, error) {
	n := uint64(len(vals))
	if n > fs.MaxWidth {
		return nil, fmt.Errorf("got %d values but only have %d roots of unity", n, fs.MaxWidth)
	}
	n = nextPowOf2(n)
	// We make a copy so we can mutate it during the work.
	valsCopy := make([]fr.Element, n)
	for i := 0; i < len(vals); i++ {
		valsCopy[i].Set(&vals[i])
		//bls.CopyFr(&valsCopy[i], &vals[i])
	}
	for i := uint64(len(vals)); i < n; i++ {
		valsCopy[i].SetZero()
		//bls.CopyFr(&valsCopy[i], &bls.ZERO)
	}
	out := make([]fr.Element, n)
	if err := fs.InplaceFFT(valsCopy, out, inv); err != nil {
		return nil, err
	}
	return out, nil
}

func (fs *FFTSettings) InplaceFFT(vals []fr.Element, out []fr.Element, inv bool) error {
	n := uint64(len(vals))
	if n > fs.MaxWidth {
		return fmt.Errorf("got %d values but only have %d roots of unity", n, fs.MaxWidth)
	}
	if !IsPowerOfTwo(n) {
		return fmt.Errorf("got %d values but not a power of two", n)
	}
	if inv {
		var invLen fr.Element
		//bls.AsFr(&invLen, n)
		invLen.SetInt64(int64(n))
		//bls.InvModFr(&invLen, &invLen)
		invLen.Inverse(&invLen)
		rootz := fs.ReverseRootsOfUnity[:fs.MaxWidth]
		stride := fs.MaxWidth / n

		fs._fft(vals, 0, 1, rootz, stride, out)
		var tmp fr.Element
		for i := 0; i < len(out); i++ {
			tmp.Mul(&out[i], &invLen)
			//bls.MulModFr(&tmp, &out[i], &invLen)
			out[i].Set(&tmp)
			//bls.CopyFr(&out[i], &tmp) // TODO: depending on Fr implementation, allow to directly write back to an input
		}
		return nil
	} else {
		rootz := fs.ExpandedRootsOfUnity[:fs.MaxWidth]
		stride := fs.MaxWidth / n
		// Regular FFT
		fs._fft(vals, 0, 1, rootz, stride, out)
		return nil
	}
}

func IsPowerOfTwo(v uint64) bool {
	return v&(v-1) == 0
}
