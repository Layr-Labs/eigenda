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
)

func TestFFTSettings_reduceLeaves(t *testing.T) {
	fs := NewFFTSettings(4)

	var fromTreeReduction []fr.Element
	{
		// prepare some leaves
		leaves := [][]fr.Element{make([]fr.Element, 3), make([]fr.Element, 3), make([]fr.Element, 3), make([]fr.Element, 3)}
		leafIndices := [][]uint64{{1, 3}, {7, 8}, {9, 10}, {12, 13}}
		for i := 0; i < 4; i++ {
			err := fs.makeZeroPolyMulLeaf(leaves[i], leafIndices[i], 1)
			assert.Nil(t, err)
		}

		dst := make([]fr.Element, 16)
		scratch := make([]fr.Element, 16*3)
		_, err := fs.reduceLeaves(scratch, dst, leaves)
		if err != nil {
			assert.Nil(t, err)
		}
		fromTreeReduction = dst[:2*4+1]
	}

	var fromDirect []fr.Element
	{
		dst := make([]fr.Element, 9)
		indices := []uint64{1, 3, 7, 8, 9, 10, 12, 13}
		err := fs.makeZeroPolyMulLeaf(dst, indices, 1)
		if err != nil {
			assert.Nil(t, err)
		}
		fromDirect = dst
	}
	assert.Equal(t, len(fromDirect), len(fromTreeReduction), "length mismatch")

	for i := 0; i < len(fromDirect); i++ {
		a, b := &fromDirect[i], &fromTreeReduction[i]
		if !a.Equal(b) {
			t.Errorf("zero poly coeff %d is different. direct: %s, tree: %s", i, a.String(), b.String())
		}
		assert.True(t, a.Equal(b),
			"zero poly coeff %d is different. direct: %s, tree: %s", i, a.String(), b.String())
	}
}

func TestFFTSettings_reduceLeaves_parametrized(t *testing.T) {
	ratios := []float64{0.01, 0.1, 0.2, 0.4, 0.5, 0.7, 0.9, 0.99}
	for scale := uint8(5); scale < 13; scale++ {
		t.Run(fmt.Sprintf("scale_%d", scale), func(t *testing.T) {
			for i, ratio := range ratios {
				t.Run(fmt.Sprintf("ratio_%.3f", ratio), func(t *testing.T) {
					seed := int64(1000*int(scale) + i)
					testReduceLeaves(scale, ratio, seed, t)
				})
			}
		})
	}
}

func testReduceLeaves(scale uint8, missingRatio float64, seed int64, t *testing.T) {
	fs := NewFFTSettings(scale)
	rng := rand.New(rand.NewSource(seed))
	pointCount := uint64(1) << scale
	missingCount := uint64(int(float64(pointCount) * missingRatio))
	if missingCount == 0 {
		return // nothing missing
	}

	// select the missing points
	missing := make([]uint64, pointCount)
	for i := uint64(0); i < pointCount; i++ {
		missing[i] = i
	}
	rng.Shuffle(int(pointCount), func(i, j int) {
		missing[i], missing[j] = missing[j], missing[i]
	})
	missing = missing[:missingCount]

	// build the leaves
	pointsPerLeaf := uint64(63)
	leafCount := (missingCount + pointsPerLeaf - 1) / pointsPerLeaf
	leaves := make([][]fr.Element, leafCount)
	for i := uint64(0); i < leafCount; i++ {
		start := i * pointsPerLeaf
		end := start + pointsPerLeaf
		if end > missingCount {
			end = missingCount
		}
		leafSize := end - start
		leaf := make([]fr.Element, leafSize+1)
		indices := make([]uint64, leafSize)
		for j := uint64(0); j < leafSize; j++ {
			indices[j] = missing[i*pointsPerLeaf+j]
		}
		err := fs.makeZeroPolyMulLeaf(leaf, indices, 1)
		assert.Nil(t, err)
		leaves[i] = leaf
	}

	var fromTreeReduction []fr.Element
	{
		dst := make([]fr.Element, pointCount)
		scratch := make([]fr.Element, pointCount*3)
		_, err := fs.reduceLeaves(scratch, dst, leaves)
		if err != nil {
			assert.Nil(t, err)
		}
		fromTreeReduction = dst[:missingCount+1]
	}

	var fromDirect []fr.Element
	{
		dst := make([]fr.Element, missingCount+1)
		err := fs.makeZeroPolyMulLeaf(dst, missing, fs.MaxWidth/pointCount)
		assert.Nil(t, err)
		fromDirect = dst
	}
	assert.Equal(t, len(fromDirect), len(fromTreeReduction), "length mismatch")

	for i := 0; i < len(fromDirect); i++ {
		a, b := &fromDirect[i], &fromTreeReduction[i]
		assert.True(t, a.Equal(b),
			"zero poly coeff %d is different. direct: %s, tree: %s", i, a.String(), b.String())
	}
}

// TODO: Make pass
// func TestFFTSettings_ZeroPolyViaMultiplication_Python(t *testing.T) {
// 	fs := NewFFTSettings(4)

// 	exists := []bool{
// 		true, false, false, true,
// 		false, true, true, false,
// 		false, false, true, true,
// 		false, true, false, true,
// 	}
// 	var missingIndices []uint64
// 	for i, v := range exists {
// 		if !v {
// 			missingIndices = append(missingIndices, uint64(i))
// 		}
// 	}

// 	zeroEval, zeroPoly, _ := fs.ZeroPolyViaMultiplication(missingIndices, uint64(len(exists)))

// 	// produced from python implementation, check it's exactly correct.
// 	expectedEval := []fr.Element{
// 		bls.ToFr("40868503138626303263713448452028063093974861640573380501185290423282553381059"),
// 		bls.ToFr("0"),
// 		bls.ToFr("0"),
// 		bls.ToFr("9059493333851894280622930192031068801018187410981018272280547403745554404951"),
// 		bls.ToFr("0"),
// 		bls.ToFr("589052107338478098858761185551735055781651813398303959420821217298541933174"),
// 		bls.ToFr("1980700778768058987161339158728243463014673552245301202287722613196911807966"),
// 		bls.ToFr("0"),
// 		bls.ToFr("0"),
// 		bls.ToFr("0"),
// 		bls.ToFr("48588946696503834689243119316363329218956542308951664733900338765742108388091"),
// 		bls.ToFr("17462668815085674001076443909983570919844170615339489499875900337907893054793"),
// 		bls.ToFr("0"),
// 		bls.ToFr("32986316229085390499922301497961243665601583888595873281538162159212447231217"),
// 		bls.ToFr("0"),
// 		bls.ToFr("31340620128536760059637470141592017333700483773455661424257920684057136952965"),
// 	}

// 	for i := range zeroEval {
// 		fmt.Println(expectedEval[i])
// 		assert.True(t, bls.EqualFr(&expectedEval[i], &zeroEval[i]),
// 			"at eval %d, expected: %s, got: %s", i, fr.ElementStr(&expectedEval[i]), fr.ElementStr(&zeroEval[i]))
// 	}

// 	expectedPoly := []fr.Element{
// 		bls.ToFr("37647706414300369857238608619982937390838535937985112215973498325246987289395"),
// 		bls.ToFr("2249310547870908874251949653552971443359134481191188461034956129255788965773"),
// 		bls.ToFr("14214218681578879810156974734536988864583938194339599855352132142401756507144"),
// 		bls.ToFr("11562429031388751544281783289945994468702719673309534612868555280828261838388"),
// 		bls.ToFr("38114263339263944057999429128256535679768370097817780187577397655496877536510"),
// 		bls.ToFr("21076784030567214561538347586500535789557219054084066119912281151549494675620"),
// 		bls.ToFr("9111875896859243625633322505516518368332415340935654725595105138403527134249"),
// 		bls.ToFr("11763665547049371891508513950107512764213633861965719968078681999977021803005"),
// 		bls.ToFr("1"),
// 		bls.ToFr("0"),
// 		bls.ToFr("0"),
// 		bls.ToFr("0"),
// 		bls.ToFr("0"),
// 		bls.ToFr("0"),
// 		bls.ToFr("0"),
// 		bls.ToFr("0"),
// 	}

// 	for i := range zeroPoly {
// 		assert.True(t, bls.EqualFr(&expectedPoly[i], &zeroPoly[i]),
// 			"at poly %d, expected: %s, got: %s", i, fr.ElementStr(&expectedPoly[i]), fr.ElementStr(&zeroPoly[i]))
// 	}
// }

func testZeroPoly(t *testing.T, scale uint8, seed int64) {
	fs := NewFFTSettings(scale)

	rng := rand.New(rand.NewSource(seed))

	exists := make([]bool, fs.MaxWidth)
	var missingIndices []uint64
	missingStr := ""
	for i := 0; i < len(exists); i++ {
		if rng.Intn(2) == 0 {
			exists[i] = true
		} else {
			missingIndices = append(missingIndices, uint64(i))
			missingStr += fmt.Sprintf(" %d", i)
		}
	}

	zeroEval, zeroPoly, _ := fs.ZeroPolyViaMultiplication(missingIndices, uint64(len(exists)))

	for i, v := range exists {
		if !v {
			var at fr.Element
			//xbls.CopyFr(&at, &fs.ExpandedRootsOfUnity[i])
			at.Set(&fs.ExpandedRootsOfUnity[i])
			var out fr.Element
			EvalPolyAt(&out, zeroPoly, &at)
			if !out.IsZero() {
				t.Errorf("expected zero at %d, but got: %s", i, out.String())
			}
		}
	}

	p, err := fs.FFT(zeroEval, true)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < len(zeroPoly); i++ {
		if !p[i].Equal(&zeroPoly[i]) {
			t.Errorf("fft not correct, i: %v, a: %s, b: %s", i, p[i].String(), zeroPoly[i].String())
		}
	}
	for i := len(zeroPoly); i < len(p); i++ {
		if !p[i].IsZero() {
			t.Errorf("fft not correct, i: %v, a: %s, b: 0", i, p[i].String())
		}
	}
}

func TestFFTSettings_ZeroPolyViaMultiplication_Parametrized(t *testing.T) {
	for i := uint8(3); i < 12; i++ {
		t.Run(fmt.Sprintf("scale_%d", i), func(t *testing.T) {
			for j := int64(0); j < 3; j++ {
				t.Run(fmt.Sprintf("case_%d", j), func(t *testing.T) {
					testZeroPoly(t, i, int64(i)*1000+j)
				})
			}
		})
	}
}
