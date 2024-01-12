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

package kzg

import (
	"fmt"

	bls "github.com/Layr-Labs/eigenda/pkg/kzg/bn254"
	//"github.com/protolambda/go-kzg/bls"
)

// unshift poly, in-place. Multiplies each coeff with 1/shift_factor**i
func (fs *FFTSettings) ShiftPoly(poly []bls.Fr) {
	var shiftFactor bls.Fr
	bls.AsFr(&shiftFactor, 5) // primitive root of unity
	var factorPower bls.Fr
	bls.CopyFr(&factorPower, &bls.ONE)
	var invFactor bls.Fr
	bls.InvModFr(&invFactor, &shiftFactor)
	var tmp bls.Fr
	for i := 0; i < len(poly); i++ {
		bls.CopyFr(&tmp, &poly[i])
		bls.MulModFr(&poly[i], &tmp, &factorPower)
		// TODO: pre-compute all these shift scalars
		bls.CopyFr(&tmp, &factorPower)
		bls.MulModFr(&factorPower, &tmp, &invFactor)
	}
}

// unshift poly, in-place. Multiplies each coeff with shift_factor**i
func (fs *FFTSettings) UnshiftPoly(poly []bls.Fr) {
	var shiftFactor bls.Fr
	bls.AsFr(&shiftFactor, 5) // primitive root of unity
	var factorPower bls.Fr
	bls.CopyFr(&factorPower, &bls.ONE)
	var tmp bls.Fr
	for i := 0; i < len(poly); i++ {
		bls.CopyFr(&tmp, &poly[i])
		bls.MulModFr(&poly[i], &tmp, &factorPower)
		// TODO: pre-compute all these shift scalars
		bls.CopyFr(&tmp, &factorPower)
		bls.MulModFr(&factorPower, &tmp, &shiftFactor)
	}
}

func (fs *FFTSettings) RecoverPolyFromSamples(samples []*bls.Fr, zeroPolyFn ZeroPolyFn) ([]bls.Fr, error) {
	// TODO: using a single additional temporary array, all the FFTs can run in-place.

	missingIndices := make([]uint64, 0, len(samples))
	for i, s := range samples {
		if s == nil {
			missingIndices = append(missingIndices, uint64(i))
		}
	}

	zeroEval, zeroPoly, err := zeroPolyFn(missingIndices, uint64(len(samples)))
	if err != nil {
		return nil, err
	}

	for i, s := range samples {
		if (s == nil) != bls.EqualZero(&zeroEval[i]) {
			panic("bad zero eval")
		}
	}

	polyEvaluationsWithZero := make([]bls.Fr, len(samples))
	for i, s := range samples {
		if s == nil {
			bls.CopyFr(&polyEvaluationsWithZero[i], &bls.ZERO)
		} else {
			bls.MulModFr(&polyEvaluationsWithZero[i], s, &zeroEval[i])
		}
	}
	polyWithZero, err := fs.FFT(polyEvaluationsWithZero, true)
	if err != nil {
		return nil, err
	}
	// shift in-place
	fs.ShiftPoly(polyWithZero)
	shiftedPolyWithZero := polyWithZero

	fs.ShiftPoly(zeroPoly)
	shiftedZeroPoly := zeroPoly

	evalShiftedPolyWithZero, err := fs.FFT(shiftedPolyWithZero, false)
	if err != nil {
		return nil, err
	}
	evalShiftedZeroPoly, err := fs.FFT(shiftedZeroPoly, false)
	if err != nil {
		return nil, err
	}

	evalShiftedReconstructedPoly := evalShiftedPolyWithZero
	for i := 0; i < len(evalShiftedReconstructedPoly); i++ {
		bls.DivModFr(&evalShiftedReconstructedPoly[i], &evalShiftedPolyWithZero[i], &evalShiftedZeroPoly[i])
	}
	shiftedReconstructedPoly, err := fs.FFT(evalShiftedReconstructedPoly, true)
	if err != nil {
		return nil, err
	}
	fs.UnshiftPoly(shiftedReconstructedPoly)
	reconstructedPoly := shiftedReconstructedPoly

	reconstructedData, err := fs.FFT(reconstructedPoly, false)
	if err != nil {
		return nil, err
	}
	for i, s := range samples {
		if s != nil && !bls.EqualFr(&reconstructedData[i], s) {
			return nil, fmt.Errorf("failed to reconstruct data correctly, changed value at index %d. Expected: %s, got: %s", i, bls.FrStr(s), bls.FrStr(&reconstructedData[i]))
		}
	}
	return reconstructedData, nil
}
