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
	"math/rand"
	"testing"

	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFFTSettings_RecoverPolyFromSamples_Simple(t *testing.T) {
	// Create some random data, with padding...
	fs := NewFFTSettings(2)
	poly := make([]fr.Element, fs.MaxWidth)
	for i := uint64(0); i < fs.MaxWidth/2; i++ {
		poly[i].SetInt64(int64(i))
	}
	for i := fs.MaxWidth / 2; i < fs.MaxWidth; i++ {
		poly[i].SetZero()
	}

	// Get data for polynomial SLOW_INDICES
	data, err := fs.FFT(poly, false)
	require.Nil(t, err)

	subset := make([]*fr.Element, fs.MaxWidth)
	subset[0] = &data[0]
	subset[3] = &data[3]

	recovered, err := fs.RecoverPolyFromSamples(subset, fs.ZeroPolyViaMultiplication)
	require.Nil(t, err)

	for i := range recovered {
		assert.True(t, recovered[i].Equal(&data[i]),
			"recovery at index %d got %s but expected %s", i, recovered[i].String(), data[i].String())
	}

	// And recover the original coeffs for good measure
	back, err := fs.FFT(recovered, true)
	require.Nil(t, err)

	for i := uint64(0); i < fs.MaxWidth/2; i++ {
		assert.True(t, back[i].Equal(&poly[i]),
			"coeff at index %d got %s but expected %s", i, back[i].String(), poly[i].String())
	}

	for i := fs.MaxWidth / 2; i < fs.MaxWidth; i++ {
		assert.True(t, back[i].IsZero(),
			"expected zero padding in index %d", i)
	}
}

func TestFFTSettings_RecoverPolyFromSamples(t *testing.T) {
	// Create some random poly, with padding so we get redundant data
	fs := NewFFTSettings(10)
	poly := make([]fr.Element, fs.MaxWidth)
	for i := uint64(0); i < fs.MaxWidth/2; i++ {
		poly[i].SetInt64(int64(i))
	}
	for i := fs.MaxWidth / 2; i < fs.MaxWidth; i++ {
		poly[i].SetZero()
	}

	// Get coefficients for polynomial SLOW_INDICES
	data, err := fs.FFT(poly, false)
	require.Nil(t, err)

	// Util to pick a random subnet of the values
	randomSubset := func(known uint64, rngSeed uint64) []*fr.Element {
		withMissingValues := make([]*fr.Element, fs.MaxWidth)
		for i := range data {
			withMissingValues[i] = &data[i]
		}
		rng := rand.New(rand.NewSource(int64(rngSeed)))
		missing := fs.MaxWidth - known
		pruned := rng.Perm(int(fs.MaxWidth))[:missing]
		for _, i := range pruned {
			withMissingValues[i] = nil
		}
		return withMissingValues
	}

	// Try different amounts of known indices, and try it in multiple random ways
	var lastKnown uint64 = 0
	for knownRatio := 0.7; knownRatio < 1.0; knownRatio += 0.05 {
		known := uint64(float64(fs.MaxWidth) * knownRatio)
		if known == lastKnown {
			continue
		}
		lastKnown = known
		for i := 0; i < 3; i++ {
			t.Run(fmt.Sprintf("random_subset_%d_known_%d", i, known), func(t *testing.T) {
				subset := randomSubset(known, uint64(i))

				recovered, err := fs.RecoverPolyFromSamples(subset, fs.ZeroPolyViaMultiplication)
				require.Nil(t, err)

				for i := range recovered {
					assert.True(t, recovered[i].Equal(&data[i]),
						"recovery at index %d got %s but expected %s", i, recovered[i].String(), data[i].String())
				}

				// And recover the original coeffs for good measure
				back, err := fs.FFT(recovered, true)
				require.Nil(t, err)

				half := uint64(len(back)) / 2
				for i := uint64(0); i < half; i++ {
					assert.True(t, back[i].Equal(&poly[i]),
						"coeff at index %d got %s but expected %s", i, back[i].String(), poly[i].String())
				}
				for i := half; i < fs.MaxWidth; i++ {
					assert.True(t, back[i].IsZero(),
						"expected zero padding in index %d", i)
				}
			})
		}
	}
}
