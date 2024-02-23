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
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

type KZGSettings struct {
	*FFTSettings

	Srs *SRS
	// setup values
}

func NewKZGSettings(fs *FFTSettings, srs *SRS) (*KZGSettings, error) {

	ks := &KZGSettings{
		FFTSettings: fs,
		Srs:         srs,
	}

	return ks, nil
}

// KZG commitment to polynomial in coefficient form
func (ks *KZGSettings) CommitToPoly(coeffs []fr.Element) *bn254.G1Affine {
	var commit bn254.G1Affine
	commit.MultiExp(ks.Srs.G1[:len(coeffs)], coeffs, ecc.MultiExpConfig{})
	return &commit
}

func HashToSingleField(dst *fr.Element, msg []byte) error {
	DST := []byte("-")
	randomFr, err := fr.Hash(msg, DST, 1)
	randomFrBytes := (randomFr[0]).Bytes()
	//FrSetBytes(dst, randomFrBytes[:])
	dst.SetBytes(randomFrBytes[:])
	return err
}
