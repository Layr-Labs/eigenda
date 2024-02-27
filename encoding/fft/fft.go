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

// Original: https://github.com/ethereum/research/blob/master/mimc_stark/fft.py

package fft

import (
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"

	"math/bits"
)

// if not already a power of 2, return the next power of 2
func nextPowOf2(v uint64) uint64 {
	if v == 0 {
		return 1
	}
	return uint64(1) << bits.Len64(v-1)
}

// Expands the power circle for a given root of unity to WIDTH+1 values.
// The first entry will be 1, the last entry will also be 1,
// for convenience when reversing the array (useful for inverses)
func expandRootOfUnity(rootOfUnity *fr.Element) []fr.Element {
	rootz := make([]fr.Element, 2)
	rootz[0].SetOne() // some unused number in py code
	rootz[1] = *rootOfUnity

	for i := 1; !rootz[i].IsOne(); {
		rootz = append(rootz, fr.Element{})
		this := &rootz[i]
		i++
		rootz[i].Mul(this, rootOfUnity)
		//bls.MulModFr(&rootz[i], this, rootOfUnity)
	}
	return rootz
}

type FFTSettings struct {
	MaxWidth uint64
	// the generator used to get all roots of unity
	RootOfUnity *fr.Element
	// domain, starting and ending with 1 (duplicate!)
	ExpandedRootsOfUnity []fr.Element
	// reverse domain, same as inverse values of domain. Also starting and ending with 1.
	ReverseRootsOfUnity []fr.Element
}

func NewFFTSettings(maxScale uint8) *FFTSettings {
	width := uint64(1) << maxScale
	root := &encoding.Scale2RootOfUnity[maxScale]
	rootz := expandRootOfUnity(&encoding.Scale2RootOfUnity[maxScale])

	// reverse roots of unity
	rootzReverse := make([]fr.Element, len(rootz))
	copy(rootzReverse, rootz)
	for i, j := uint64(0), uint64(len(rootz)-1); i < j; i, j = i+1, j-1 {
		rootzReverse[i], rootzReverse[j] = rootzReverse[j], rootzReverse[i]
	}

	return &FFTSettings{
		MaxWidth:             width,
		RootOfUnity:          root,
		ExpandedRootsOfUnity: rootz,
		ReverseRootsOfUnity:  rootzReverse,
	}
}
