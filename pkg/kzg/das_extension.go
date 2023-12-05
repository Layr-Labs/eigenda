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

import bls "github.com/Layr-Labs/eigenda/pkg/kzg/bn254"

// warning: the values in `a` are modified in-place to become the outputs.
// Make a deep copy first if you need to use them later.
func (fs *FFTSettings) dASFFTExtension(ab []bls.Fr, domainStride uint64) {
	if len(ab) == 2 {
		aHalf0 := &ab[0]
		aHalf1 := &ab[1]
		var x bls.Fr
		bls.AddModFr(&x, aHalf0, aHalf1)
		var y bls.Fr
		bls.SubModFr(&y, aHalf0, aHalf1)
		var tmp bls.Fr
		bls.MulModFr(&tmp, &y, &fs.ExpandedRootsOfUnity[domainStride])
		bls.AddModFr(&ab[0], &x, &tmp)
		bls.SubModFr(&ab[1], &x, &tmp)
		return
	}

	if len(ab) < 2 {
		panic("bad usage")
	}

	half := uint64(len(ab))
	halfHalf := half >> 1
	abHalf0s := ab[:halfHalf]
	abHalf1s := ab[halfHalf:half]
	// Instead of allocating L0 and L1, just modify a in-place.
	//L0[i] = (((a_half0 + a_half1) % modulus) * inv2) % modulus
	//R0[i] = (((a_half0 - L0[i]) % modulus) * inverse_domain[i * 2]) % modulus
	var tmp1, tmp2 bls.Fr
	for i := uint64(0); i < halfHalf; i++ {
		aHalf0 := &abHalf0s[i]
		aHalf1 := &abHalf1s[i]
		bls.AddModFr(&tmp1, aHalf0, aHalf1)
		bls.SubModFr(&tmp2, aHalf0, aHalf1)
		bls.MulModFr(aHalf1, &tmp2, &fs.ReverseRootsOfUnity[i*2*domainStride])
		bls.CopyFr(aHalf0, &tmp1)
	}

	// L will be the left half of out
	fs.dASFFTExtension(abHalf0s, domainStride<<1)
	// R will be the right half of out
	fs.dASFFTExtension(abHalf1s, domainStride<<1)

	// The odd deduced outputs are written to the output array already, but then updated in-place
	// L1 = b[:halfHalf]
	// R1 = b[halfHalf:]

	// Half the work of a regular FFT: only deal with uneven-index outputs
	var yTimesRoot bls.Fr
	var x, y bls.Fr
	for i := uint64(0); i < halfHalf; i++ {
		// Temporary copies, so that writing to output doesn't conflict with input.
		// Note that one hand is from L1, the other R1
		bls.CopyFr(&x, &abHalf0s[i])
		bls.CopyFr(&y, &abHalf1s[i])
		root := &fs.ExpandedRootsOfUnity[(1+2*i)*domainStride]
		bls.MulModFr(&yTimesRoot, &y, root)
		// write outputs in place, avoid unnecessary list allocations
		bls.AddModFr(&abHalf0s[i], &x, &yTimesRoot)
		bls.SubModFr(&abHalf1s[i], &x, &yTimesRoot)
	}
}

// Takes vals as input, the values of the even indices.
// Then computes the values for the odd indices, which combined would make the right half of coefficients zero.
// Warning: the odd results are written back to the vals slice.
func (fs *FFTSettings) DASFFTExtension(vals []bls.Fr) {
	if uint64(len(vals))*2 > fs.MaxWidth {
		panic("domain too small for extending requested values")
	}
	fs.dASFFTExtension(vals, 1)
	// The above function didn't perform the divide by 2 on every layer.
	// So now do it all at once, by dividing by 2**depth (=length).
	var invLen bls.Fr
	bls.AsFr(&invLen, uint64(len(vals)))
	bls.InvModFr(&invLen, &invLen)
	for i := 0; i < len(vals); i++ {
		bls.MulModFr(&vals[i], &vals[i], &invLen)
	}
}
