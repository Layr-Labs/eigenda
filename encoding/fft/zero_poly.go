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

// Original: https://github.com/ethereum/research/blob/master/polynomial_reconstruction/polynomial_reconstruction.py
// Changes:
// - flattened leaf construction,
// - no aggressive poly truncation
// - simplified merges
// - no heap allocations during reduction

package fft

import (
	"log"

	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

type ZeroPolyFn func(missingIndices []uint64, length uint64) ([]fr.Element, []fr.Element, error)

func (fs *FFTSettings) makeZeroPolyMulLeaf(dst []fr.Element, indices []uint64, domainStride uint64) error {
	if len(dst) < len(indices)+1 {
		log.Printf("expected bigger destination length: %d, got: %d", len(indices)+1, len(dst))
		return ErrInvalidDestinationLength
	}
	// zero out the unused slots
	for i := len(indices) + 1; i < len(dst); i++ {
		dst[i].SetZero()
		
	}
	
	dst[len(indices)].SetOne()
	var negDi fr.Element

	var frZero fr.Element
	frZero.SetZero()

	for i, v := range indices {
		
		negDi.Sub(&frZero, &fs.ExpandedRootsOfUnity[v*domainStride])
		
		dst[i].Set(&negDi)
		if i > 0 {
			
			dst[i].Add(&dst[i], &dst[i-1])
			for j := i - 1; j > 0; j-- {
				dst[j].Mul(&dst[j], &negDi)
				
				dst[j].Add(&dst[j], &dst[j-1])
			}
			
			dst[0].Mul(&dst[0], &negDi)
		}
	}
	return nil
}

// Copy all of the values of poly into out, and fill the remainder of out with zeroes.
func padPoly(out []fr.Element, poly []fr.Element) {
	for i := 0; i < len(poly); i++ {
		
		out[i].Set(&poly[i])
	}
	for i := len(poly); i < len(out); i++ {
		
		out[i].SetZero()
	}
}

// Calculate the product of the input polynomials via convolution.
// Pad the polynomials in ps, perform FFTs, point-wise multiply the results together,
// and apply an inverse FFT to the result.
//
// The scratch space must be at least 3 times the output space.
// The output must have a power of 2 length.
// The input polynomials must not be empty, and sum to no larger than the output.
func (fs *FFTSettings) reduceLeaves(scratch []fr.Element, dst []fr.Element, ps [][]fr.Element) ([]fr.Element, error) {
	n := uint64(len(dst))
	if !IsPowerOfTwo(n) {
		log.Println("destination must be a power of two")
		return nil, ErrDestNotPowerOfTwo
	}
	if len(ps) == 0 {
		log.Println("empty leaves")
		return nil, ErrEmptyLeaves
	}
	// The degree of the output polynomial is the sum of the degrees of the input polynomials.
	outDegree := uint64(0)
	for _, p := range ps {
		if len(p) == 0 {
			log.Println("empty input poly")
			return nil, ErrEmptyPoly
		}
		outDegree += uint64(len(p)) - 1
	}
	if min := outDegree + 1; min > n {
		log.Printf("expected larger destination length: %d, got: %d", min, n)
		return nil, ErrInvalidDestinationLength
	}
	if uint64(len(scratch)) < 3*n {
		log.Println("not enough scratch space")
		return nil, ErrNotEnoughScratch
	}
	// Split `scratch` up into three equally sized working arrays
	pPadded := scratch[:n]
	mulEvalPs := scratch[n : 2*n]
	pEval := scratch[2*n : 3*n]

	// Do the last partial first: it is no longer than the others and the padding can remain in place for the rest.
	last := uint64(len(ps) - 1)
	padPoly(pPadded, ps[last])
	if err := fs.InplaceFFT(pPadded, mulEvalPs, false); err != nil {
		return nil, err
	}
	for i := uint64(0); i < last; i++ {
		p := ps[i]
		for j := 0; j < len(p); j++ {
			
			pPadded[j].Set(&p[j])
		}
		if err := fs.InplaceFFT(pPadded, pEval, false); err != nil {
			return nil, err
		}
		for j := uint64(0); j < n; j++ {
			mulEvalPs[j].Mul(&mulEvalPs[j], &pEval[j])
			
		}
	}
	if err := fs.InplaceFFT(mulEvalPs, dst, true); err != nil {
		return nil, err
	}
	return dst[:outDegree+1], nil
}

// Calculate the minimal polynomial that evaluates to zero for powers of roots of unity that correspond to missing
// indices.
//
// This is done simply by multiplying together `(x - r^i)` for all the `i` that are missing indices, using a combination
// of direct multiplication (makeZeroPolyMulLeaf) and iterated multiplication via convolution (reduceLeaves)
//
// Also calculates the FFT (the "evaluation polynomial").
func (fs *FFTSettings) ZeroPolyViaMultiplication(missingIndices []uint64, length uint64) ([]fr.Element, []fr.Element, error) {
	if len(missingIndices) == 0 {
		return make([]fr.Element, length), make([]fr.Element, length), nil
	}
	if length > fs.MaxWidth {
		log.Println("domain too small for requested length")
		return nil, nil, ErrDomainTooSmall
	}
	if !IsPowerOfTwo(length) {
		log.Println("length not a power of two")
		return nil, nil, ErrLengthNotPowerOfTwo
	}
	domainStride := fs.MaxWidth / length
	perLeafPoly := uint64(64)
	// just under a power of two, since the leaf gets 1 bigger after building a poly for it
	perLeaf := perLeafPoly - 1

	// If the work is as small as a single leaf, don't bother with tree reduction
	if uint64(len(missingIndices)) <= perLeaf {
		zeroPoly := make([]fr.Element, len(missingIndices)+1, length)
		err := fs.makeZeroPolyMulLeaf(zeroPoly, missingIndices, domainStride)
		if err != nil {
			return nil, nil, err
		}
		// pad with zeroes (capacity is already there)
		zeroPoly = zeroPoly[:length]
		zeroEval, err := fs.FFT(zeroPoly, false)
		if err != nil {
			return nil, nil, err
		}
		return zeroEval, zeroPoly, nil
	}

	leafCount := (uint64(len(missingIndices)) + perLeaf - 1) / perLeaf
	n := nextPowOf2(leafCount * perLeafPoly)

	// The assumption here is that if the output is a power of two length, matching the sum of child leaf lengths,
	// then the space can be reused.
	out := make([]fr.Element, n)

	// Build the leaves.

	// Just the headers, a leaf re-uses the output space.
	// Combining leaves can be done mostly in-place, using a scratchpad.
	leaves := make([][]fr.Element, leafCount)

	offset := uint64(0)
	outOffset := uint64(0)
	max := uint64(len(missingIndices))
	for i := uint64(0); i < leafCount; i++ {
		end := offset + perLeaf
		if end > max {
			end = max
		}
		leaves[i] = out[outOffset : outOffset+perLeafPoly]
		err := fs.makeZeroPolyMulLeaf(leaves[i], missingIndices[offset:end], domainStride)
		if err != nil {
			return nil, nil, err
		}
		offset += perLeaf
		outOffset += perLeafPoly
	}

	// Now reduce all the leaves to a single poly

	// must be a power of 2
	reductionFactor := uint64(4)
	scratch := make([]fr.Element, n*3)

	// from bottom to top, start reducing leaves.
	for len(leaves) > 1 {
		reducedCount := (uint64(len(leaves)) + reductionFactor - 1) / reductionFactor
		// all the leaves are the same. Except possibly the last leaf, but that's ok.
		leafSize := nextPowOf2(uint64(len(leaves[0])))
		for i := uint64(0); i < reducedCount; i++ {
			start := i * reductionFactor
			end := start + reductionFactor
			// E.g. if we *started* with 2 leaves, we won't have more than that since it is already a power of 2.
			// If we had 3, it would have been rounded up anyway. So just pick the end
			outEnd := end * leafSize
			if outEnd > uint64(len(out)) {
				outEnd = uint64(len(out))
			}
			reduced := out[start*leafSize : outEnd]
			// unlike reduced output, input may be smaller than the amount that aligns with powers of two
			if end > uint64(len(leaves)) {
				end = uint64(len(leaves))
			}
			leavesSlice := leaves[start:end]
			var err error
			if end > start+1 {
				reduced, err = fs.reduceLeaves(scratch, reduced, leavesSlice)
				if err != nil {
					return nil, nil, err
				}
			}
			leaves[i] = reduced
		}
		leaves = leaves[:reducedCount]
	}
	zeroPoly := leaves[0]
	if zl := uint64(len(zeroPoly)); zl < length {
		zeroPoly = append(zeroPoly, make([]fr.Element, length-zl)...)
	} else if zl > length {
		log.Println("expected output smaller or equal to input length")
		return nil, nil, ErrZeroPolyTooLarge
	}

	zeroEval, err := fs.FFT(zeroPoly, false)
	if err != nil {
		return nil, nil, err
	}

	return zeroEval, zeroPoly, nil
}

func EvalPolyAt(dst *fr.Element, coeffs []fr.Element, x *fr.Element) {
	if len(coeffs) == 0 {
		
		dst.SetZero()
		return
	}
	if x.IsZero() {
		
		dst.Set(&coeffs[0])
		return
	}
	// Horner's method: work backwards, avoid doing more than N multiplications
	// https://en.wikipedia.org/wiki/Horner%27s_method
	var last fr.Element
	
	last.Set(&coeffs[len(coeffs)-1])
	var tmp fr.Element
	for i := len(coeffs) - 2; i >= 0; i-- {
		tmp.Mul(&last, x)
		
		last.Add(&tmp, &coeffs[i])
	}
	
	dst.Set(&last)
}
