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
	"fmt"

	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr/fft"
)

type FFTSettings struct {
	// Maximum number of points this FFTSettings can handle
	MaxWidth uint64
	// the generator used to get all roots of unity
	RootOfUnity *fr.Element
	// domain, starting and ending with 1 (duplicate!)
	ExpandedRootsOfUnity []fr.Element
	// reverse domain, same as inverse values of domain. Also starting and ending with 1.
	ReverseRootsOfUnity []fr.Element
	// Used for Fr FFTs using gnark-crypto library
	// TODO: replace FFTSettings entirely with this
	Domain *fft.Domain
}

// NewFFTSettings creates FFTSettings for a given maximum scale (log2 of max width).
// Precomputes the roots of unity for all widths up to 2^maxScale.
// Note that MaxWith is in units of Fr elements, so the actual byte size is 32 * MaxWidth.
// In order to FFT a blob of size 16MiB, you thus need maxScale=19 (2^19 * 32 = 16MiB).
func NewFFTSettings(maxScale uint8) *FFTSettings {
	width := uint64(1) << maxScale
	root := &encoding.Scale2RootOfUnity[maxScale]
	rootz := expandRootOfUnity(maxScale)

	// reverse roots of unity
	rootzReverse := make([]fr.Element, len(rootz))
	copy(rootzReverse, rootz)
	for i, j := uint64(0), uint64(len(rootz)-1); i < j; i, j = i+1, j-1 {
		rootzReverse[i], rootzReverse[j] = rootzReverse[j], rootzReverse[i]
	}
	domain := fft.NewDomain(width)

	return &FFTSettings{
		MaxWidth:             width,
		RootOfUnity:          root,
		ExpandedRootsOfUnity: rootz,
		ReverseRootsOfUnity:  rootzReverse,
		Domain:               domain,
	}
}

// Expands the power circle for a given root of unity to WIDTH+1 values.
// The first entry will be 1, the last entry will also be 1,
// for convenience when reversing the array (useful for inverses)
func expandRootOfUnity(maxScale uint8) []fr.Element {
	rootOfUnity := encoding.Scale2RootOfUnity[maxScale]
	// preallocate with capacity for all roots of unity
	// There are 2^maxScale roots of unity, plus the duplicate 1 at the end.
	rootz := make([]fr.Element, (1<<maxScale)+1)
	rootz[0].SetOne()
	rootz[1] = rootOfUnity

	for i := 2; i < len(rootz); i++ {
		rootz[i].Mul(&rootz[i-1], &rootOfUnity)
	}
	if rootz[len(rootz)-1].Cmp(new(fr.Element).SetOne()) != 0 {
		panic(fmt.Sprintf("last root of unity is not 1, got %v", rootz[len(rootz)-1]))
	}
	return rootz
}
