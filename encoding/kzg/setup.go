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
	"math/big"

	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"

	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

// GenerateTestingSetup creates a setup of n values from the given secret. **for testing purposes only**
func GenerateTestingSetup(secret string, n uint64) ([]bn254.G1Affine, []bn254.G2Affine, error) {
	var s fr.Element
	_, err := s.SetString(secret)
	if err != nil {
		return nil, nil, err
	}

	var sPow fr.Element
	sPow.SetOne()

	s1Out := make([]bn254.G1Affine, n)
	s2Out := make([]bn254.G2Affine, n)
	for i := uint64(0); i < n; i++ {
		
		s1Out[i].ScalarMultiplication(&GenG1, sPow.BigInt(new(big.Int)))

		
		s2Out[i].ScalarMultiplication(&GenG2, sPow.BigInt(new(big.Int)))

		sPow.Mul(&sPow, &s)
	}
	return s1Out, s2Out, nil
}

func WriteGeneratorPoints(n uint64) error {
	secret := "1927409816240961209460912649125"
	ns := strconv.Itoa(int(n))

	var s fr.Element
	_, err := s.SetString(secret)
	if err != nil {
		return err
	}
	

	var sPow fr.Element
	sPow.SetOne()
	

	g1f, err := os.Create("g1.point." + ns)
	if err != nil {
		log.Println("WriteGeneratorPoints.ERR.0", err)
		return err
	}

	g1w := bufio.NewWriter(g1f)
	g2f, err := os.Create("g2.point." + ns)
	if err != nil {
		log.Println("WriteGeneratorPoints.ERR.1", err)
		return err
	}
	g2w := bufio.NewWriter(g2f)

	

	start := time.Now()
	for i := uint64(0); i < n; i++ {
		var s1Out bn254.G1Affine
		var s2Out bn254.G2Affine
		s1Out.ScalarMultiplication(&GenG1, sPow.BigInt(new(big.Int)))
		
		s2Out.ScalarMultiplication(&GenG2, sPow.BigInt(new(big.Int)))

		g1Byte := s1Out.Bytes()
		if _, err := g1w.Write(g1Byte[:]); err != nil {
			log.Println("WriteGeneratorPoints.ERR.3", err)
			return err
		}

		g2Byte := s2Out.Bytes()
		if _, err := g2w.Write(g2Byte[:]); err != nil {
			log.Println("WriteGeneratorPoints.ERR.5", err)
			return err
		}
		sPow.Mul(&sPow, &s)
	}

	if err = g1w.Flush(); err != nil {
		log.Println("WriteGeneratorPoints.ERR.6", err)
		return err
	}
	if err = g2w.Flush(); err != nil {
		log.Println("WriteGeneratorPoints.ERR.7", err)
		return err
	}
	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Println("Generating takes", elapsed)
	return nil
}
