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
	"errors"
	"fmt"

	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

type FFTInputNotPowerOfTwoError struct {
	inputLen uint64
}

func (e *FFTInputNotPowerOfTwoError) Error() string {
	return fmt.Sprintf("(I)FFT input length %d is not a power of two", e.inputLen)
}

func (e *FFTInputNotPowerOfTwoError) Is(target error) bool {
	if _, ok := target.(*FFTInputNotPowerOfTwoError); ok {
		return true
	}
	return false
}

func NewFFTInputNotPowerOfTwoError(inputLen uint64) *FFTInputNotPowerOfTwoError {
	return &FFTInputNotPowerOfTwoError{
		inputLen: inputLen,
	}
}

var (
	// Sentinel error that can be used to check if an error is an FFTInputNotPowerOfTwoError
	// by calling errors.Is(err, ErrNotPowerOfTwo)
	ErrNotPowerOfTwo = &FFTInputNotPowerOfTwoError{}
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
		}
		out[i].Set(&last)
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

		root := &rootsOfUnity[i*rootsOfUnityStride]
		yTimesRoot.Mul(&y, root)
		out[i].Add(&x, &yTimesRoot)
		out[i+half].Sub(&x, &yTimesRoot)
	}
}

// FFT performs a fast Fourier transform on the provided values, using the roots of unity
// provided in the FFTSettings.
//
// The input values does not have to be a power of two, because we pad them to the next power of two.
//
// It outputs a newly allocated slice of field elements, which is the transformed values.
// To perform the FFT in-place, use [FFTSettings.InplaceFFT] instead.
//
// The only error returned is if the FFTSettings does not have enough roots of unity to perform the FFT on the input values.
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

	}
	for i := uint64(len(vals)); i < n; i++ {
		valsCopy[i].SetZero()
	}
	out := make([]fr.Element, n)
	if err := fs.InplaceFFT(valsCopy, out, inv); err != nil {
		if errors.Is(err, ErrNotPowerOfTwo) {
			panic("bug: we passed a non-power of two to FFT, which is not possible because we called nextPowOf2 on the input above")
		}
		panic(fmt.Sprintf("bug: InplaceFFT doesn't contain enough roots of unity to perform the computation, "+
			"which is impossible because we already checked it above: %v", err))
	}
	return out, nil
}

func (fs *FFTSettings) InplaceFFT(vals []fr.Element, out []fr.Element, inv bool) error {
	n := uint64(len(vals))
	if n > fs.MaxWidth {
		return fmt.Errorf("got %d values but only have %d roots of unity", n, fs.MaxWidth)
	}
	if !IsPowerOfTwo(n) {
		return NewFFTInputNotPowerOfTwoError(n)
	}
	if inv {
		var invLen fr.Element

		invLen.SetInt64(int64(n))

		invLen.Inverse(&invLen)
		rootz := fs.ReverseRootsOfUnity[:fs.MaxWidth]
		stride := fs.MaxWidth / n

		fs._fft(vals, 0, 1, rootz, stride, out)
		var tmp fr.Element
		for i := 0; i < len(out); i++ {
			tmp.Mul(&out[i], &invLen)
			out[i].Set(&tmp)
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

// IsPowerOfTwo returns true if the provided integer v is a power of 2.
func IsPowerOfTwo(v uint64) bool {
	return (v&(v-1) == 0) && (v != 0)
}
