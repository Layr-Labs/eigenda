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

// unshift poly, in-place. Multiplies each coeff with 1/shift_factor**i
func (fs *FFTSettings) ShiftPoly(poly []fr.Element) {
	var shiftFactor fr.Element
	shiftFactor.SetInt64(int64(5))
	var factorPower fr.Element
	factorPower.SetOne()
	var invFactor fr.Element
	invFactor.Inverse(&shiftFactor)
	var tmp fr.Element
	for i := 0; i < len(poly); i++ {

		tmp.Set(&poly[i])

		poly[i].Mul(&tmp, &factorPower)

		// TODO: pre-compute all these shift scalars

		tmp.Set(&factorPower)

		factorPower.Mul(&tmp, &invFactor)
	}
}

// unshift poly, in-place. Multiplies each coeff with shift_factor**i
func (fs *FFTSettings) UnshiftPoly(poly []fr.Element) {
	var shiftFactor fr.Element

	shiftFactor.SetInt64(int64(5))
	var factorPower fr.Element
	factorPower.SetOne()

	var tmp fr.Element
	for i := 0; i < len(poly); i++ {
		tmp.Set(&poly[i])

		poly[i].Mul(&tmp, &factorPower)

		// TODO: pre-compute all these shift scalars

		tmp.Set(&factorPower)

		factorPower.Mul(&tmp, &shiftFactor)
	}
}

func (fs *FFTSettings) RecoverPolyFromSamples(samples []*fr.Element, zeroPolyFn ZeroPolyFn) ([]fr.Element, error) {
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
		if (s == nil) != zeroEval[i].IsZero() {
			return nil, errors.New("bad zero eval")
		}
	}

	polyEvaluationsWithZero := make([]fr.Element, len(samples))
	for i, s := range samples {
		if s == nil {

			polyEvaluationsWithZero[i].SetZero()
		} else {

			polyEvaluationsWithZero[i].Mul(s, &zeroEval[i])
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

		evalShiftedReconstructedPoly[i].Div(&evalShiftedPolyWithZero[i], &evalShiftedZeroPoly[i])
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
		if s != nil && !reconstructedData[i].Equal(s) {
			return nil, fmt.Errorf("failed to reconstruct data correctly, changed value at index %d. Expected: %s, got: %s", i, s.String(), reconstructedData[i].String())
		}
	}
	return reconstructedData, nil
}
