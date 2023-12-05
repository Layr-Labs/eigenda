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
	"testing"
)

// setup:
// alloc random application data
// change to reverse bit order
// extend data
// compute commitment over extended data
// func integrationTestSetup(scale uint8, seed int64) (data []byte, extended []bls.Fr, extendedAsPoly []bls.Fr, commit *bls.G1Point, ks *KZGSettings) {
// 	points := 1 << scale
// 	size := points * 31
// 	data = make([]byte, size)
// 	rng := rand.New(rand.NewSource(seed))
// 	rng.Read(data)
// 	for i := 0; i < 100; i++ {
// 		data[i] = 0
// 	}
// 	evenPoints := make([]bls.Fr, points)
// 	// fr nums are set from little-endian ints. The upper byte is always zero for input data.
// 	// 5/8 top bits are unused, other 3 out of range for modulus.
// 	var tmp [32]byte
// 	for i := 0; i < points; i++ {
// 		copy(tmp[:31], data[i*31:(i+1)*31])
// 		bls.FrFrom32(&evenPoints[i], tmp)
// 	}
// 	reverseBitOrderFr(evenPoints)
// 	oddPoints := make([]bls.Fr, points)
// 	for i := 0; i < points; i++ {
// 		bls.CopyFr(&oddPoints[i], &evenPoints[i])
// 	}
// 	// scale is 1 bigger here, since extended data is twice as big
// 	fs := NewFFTSettings(scale + 1)
// 	// convert even points (previous contents of array) to odd points
// 	fs.DASFFTExtension(oddPoints)
// 	extended = make([]bls.Fr, points*2)
// 	for i := 0; i < len(extended); i += 2 {
// 		bls.CopyFr(&extended[i], &evenPoints[i/2])
// 		bls.CopyFr(&extended[i+1], &oddPoints[i/2])
// 	}
// 	s1, s2 := GenerateTestingSetup("1927409816240961209460912649124", uint64(len(extended)))
// 	srs, _ := NewSrs(s1, s2)
// 	ks, _ = NewKZGSettings(fs, srs)
// 	// get coefficient form (half of this is zeroes, but ok)
// 	coeffs, err := ks.FFT(extended, true)
// 	if err != nil {
// 		panic(err)
// 	}
// 	debugFrs("poly", coeffs)
// 	extendedAsPoly = coeffs
// 	// the 2nd half is all zeroes, can ignore it for faster commitment.
// 	commit = ks.CommitToPoly(coeffs[:points])
// 	return
// }

// TODO: make it pass?
// func TestFullDAS(t *testing.T) {
// 	data, extended, extendedAsPoly, commit, ks := integrationTestSetup(10, 1234)
// 	// undo the bit-reverse ordering of the extended data (which was prepared after reverse-bit ordering the input data)
// 	reverseBitOrderFr(extended)
// 	debugFrs("extended data (reordered to original)", extended)

// 	cosetWidth := uint64(128)
// 	fk := NewFK20MultiSettings(ks, ks.MaxWidth, cosetWidth)
// 	// compute proofs for cosets
// 	proofs := fk.FK20MultiDAOptimized(extendedAsPoly)

// 	// package data of cosets with respective proofs
// 	sampleCount := uint64(len(extended)) / cosetWidth
// 	samples := make([]sample, sampleCount, sampleCount)
// 	for i := uint64(0); i < sampleCount; i++ {
// 		sample := &samples[i]

// 		// we can just select it from the original points
// 		sample.sub = make([]bls.Fr, cosetWidth, cosetWidth)
// 		for j := uint64(0); j < cosetWidth; j++ {
// 			bls.CopyFr(&sample.sub[j], &extended[i*cosetWidth+j])
// 		}
// 		debugFrs("sample pre-order", sample.sub)

// 		// construct that same coset from the polynomial form, to make sure we have the correct points.
// 		domainPos := reverseBitsLimited(uint32(sampleCount), uint32(i))

// 		sample.proof = &proofs[domainPos]
// 	}
// 	// skip sample serialization/deserialization, no network to transfer data here.

// 	// verify cosets individually
// 	extSize := sampleCount * cosetWidth
// 	domainStride := ks.MaxWidth / extSize
// 	for i, sample := range samples {
// 		var x bls.Fr
// 		domainPos := uint64(reverseBitsLimited(uint32(sampleCount), uint32(i)))
// 		bls.CopyFr(&x, &ks.ExpandedRootsOfUnity[domainPos*domainStride])
// 		reverseBitOrderFr(sample.sub) // match poly order
// 		val, err := ks.CheckProofMulti(commit, sample.proof, &x, sample.sub)

// 		require.Nil(t, err, "failed to verify proof of sample %d", i)
// 		assert.True(t, val, "failed to verify proof of sample %d", i)

// 		reverseBitOrderFr(sample.sub) // match original data order
// 	}

// 	// make some samples go missing
// 	partialReconstructed := make([]*bls.Fr, extSize, extSize)
// 	rng := rand.New(rand.NewSource(42))
// 	missing := 0
// 	for i, sample := range samples { // samples are already ordered in original data order
// 		// make a random subset (but <= 1/2) go missing.
// 		if rng.Int31n(2) == 0 && missing < len(samples)/2 {
// 			t.Logf("not using sample %d", i)
// 			missing++
// 			continue
// 		}

// 		offset := uint64(i) * cosetWidth
// 		for j := uint64(0); j < cosetWidth; j++ {
// 			partialReconstructed[offset+j] = &sample.sub[j]
// 		}
// 	}
// 	// samples were slices of reverse-bit-ordered data. Undo that order first, then IFFT will match the polynomial.
// 	reverseBitOrderFrPtr(partialReconstructed)
// 	// recover missing data
// 	recovered, err := ks.ErasureCodeRecover(partialReconstructed)
// 	require.Nil(t, err)

// 	// apply reverse bit-ordering again to get original data into first half
// 	reverseBitOrderFr(recovered)
// 	debugFrs("recovered", recovered)

// 	for i := 0; i < len(recovered); i++ {
// 		assert.True(t, bls.EqualFr(&extended[i], &recovered[i]),
// 			"diff %d: %s <> %s", i, bls.FrStr(&extended[i]), bls.FrStr(&recovered[i]))
// 	}
// 	// take first half, convert back to bytes
// 	size := extSize / 2
// 	reconstructedData := make([]byte, size*31, size*31)
// 	for i := uint64(0); i < size; i++ {
// 		p := bls.FrTo32(&recovered[i])
// 		copy(reconstructedData[i*31:(i+1)*31], p[:31])
// 	}

// 	// check that data matches original
// 	assert.Equal(t, data, reconstructedData, "failed to reconstruct original data")
// }

func TestFullUser(t *testing.T) {
	// setup:
	// alloc random application data
	// change to reverse bit order
	// extend data
	// compute commitment over extended data

	// construct application-layer proof for some random points
	// verify application-layer proof
}
