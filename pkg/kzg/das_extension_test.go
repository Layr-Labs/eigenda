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
package kzg

import (
	"fmt"
	"math/rand"
	"testing"

	bls "github.com/Layr-Labs/eigenda/pkg/kzg/bn254"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDASFFTExtension(t *testing.T) {
	fs := NewFFTSettings(4)
	half := fs.MaxWidth / 2
	data := make([]bls.Fr, half)
	for i := uint64(0); i < half; i++ {
		bls.AsFr(&data[i], i)
	}
	fs.DASFFTExtension(data)

	expected := []bls.Fr{
		bls.ToFr("9455244345631016631523862383826656817909262240618707851288319855253023724498"),
		bls.ToFr("10961351032263120273117550959237409754492768732192557560880754261368126052154"),
		bls.ToFr("12432998526208258595130464331726862113180416131685271896346981464741203555841"),
		bls.ToFr("8526477789819225339470181130720054257555649678225561659213114114555273578201"),
		bls.ToFr("12432998526208258595130464331726862113180416131685271896346981464741203555841"),
		bls.ToFr("10961351032263120273117550959237409754492768732192557560880754261368126052154"),
		bls.ToFr("9455244345631016631523862383826656817909262240618707851288319855253023724498"),
		bls.ToFr("13327305889333084549971686500727188725472913714445501098547591469023253739309"),
	}

	for i := range data {
		assert.True(t, bls.EqualFr(&data[i], &expected[i]))
	}
}

func TestParametrizedDASFFTExtension(t *testing.T) {
	testScale := func(seed int64, scale uint8, t *testing.T) {
		fs := NewFFTSettings(scale)
		evenData := make([]bls.Fr, fs.MaxWidth/2)
		rng := rand.New(rand.NewSource(seed))
		for i := uint64(0); i < fs.MaxWidth/2; i++ {
			bls.AsFr(&evenData[i], rng.Uint64()) // TODO could be a full random F_r instead of uint64
		}
		// we don't want to modify the original input, and the inner function would modify it in-place, so make a copy.
		oddData := make([]bls.Fr, fs.MaxWidth/2)
		for i := 0; i < len(oddData); i++ {
			bls.CopyFr(&oddData[i], &evenData[i])
		}
		fs.DASFFTExtension(oddData)

		// reconstruct data
		data := make([]bls.Fr, fs.MaxWidth)
		for i := uint64(0); i < fs.MaxWidth; i += 2 {
			bls.CopyFr(&data[i], &evenData[i>>1])
			bls.CopyFr(&data[i+1], &oddData[i>>1])
		}
		// get coefficients of reconstructed data with inverse FFT
		coeffs, err := fs.FFT(data, true)
		require.Nil(t, err)
		require.NotNil(t, coeffs)

		// second half of all coefficients should be zero
		for i := fs.MaxWidth / 2; i < fs.MaxWidth; i++ {
			assert.True(t, bls.EqualZero(&coeffs[i]), "expected zero coefficient on index %d", i)
		}
	}
	for scale := uint8(4); scale < 10; scale++ {
		for i := int64(0); i < 4; i++ {
			t.Run(fmt.Sprintf("scale_%d_i_%d", scale, i), func(t *testing.T) {
				testScale(i, scale, t)
			})
		}
	}
}
