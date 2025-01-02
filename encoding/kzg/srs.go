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
	"errors"
	"fmt"
	"math"

	"github.com/consensys/gnark-crypto/ecc/bn254"
)

type SRS struct {
	// [b.multiply(b.G1, pow(s, i, MODULUS)) for i in range(WIDTH+1)],
	G1 []bn254.G1Affine
	// [b.multiply(b.G2, pow(s, i, MODULUS)) for i in range(WIDTH+1)],
	G2 []bn254.G2Affine
	// [b.multiply(b.G2, pow(s, i, MODULUS)) for i in range(SRS_ORDER - WIDTH -1, SRS_ORDER)],
	G2Trailing []bn254.G2Affine
}

// NewSrs initializes the SRS struct using the configuration specified in a KzgConfig
func NewSrs(kzgConfig *KzgConfig) (*SRS, error) {
	if kzgConfig.SRSNumberToLoad > kzgConfig.SRSOrder {
		return nil, fmt.Errorf(
			"SRSOrder (%d) is less than srsNumberToLoad (%d)",
			kzgConfig.SRSOrder,
			kzgConfig.SRSNumberToLoad)
	}

	// read the whole order, and treat it as entire SRS for low degree proof
	s1, err := ReadG1Points(kzgConfig.G1Path, kzgConfig.SRSNumberToLoad, kzgConfig.NumWorker)
	if err != nil {
		return nil, fmt.Errorf("failed to read %d G1 points from %s: %v", kzgConfig.SRSNumberToLoad, kzgConfig.G1Path, err)
	}

	s2 := make([]bn254.G2Affine, 0)
	g2Trailing := make([]bn254.G2Affine, 0)

	if kzgConfig.LoadG2Points {
		if len(kzgConfig.G2Path) == 0 {
			return nil, errors.New("G2Path is empty. However, object needs to load G2Points")
		}

		s2, err = ReadG2Points(kzgConfig.G2Path, kzgConfig.SRSNumberToLoad, kzgConfig.NumWorker)
		if err != nil {
			return nil, fmt.Errorf("failed to read %d G2 points from %s: %v", kzgConfig.SRSNumberToLoad, kzgConfig.G2Path, err)
		}

		g2Trailing, err = ReadG2PointSection(
			kzgConfig.G2Path,
			kzgConfig.SRSOrder-kzgConfig.SRSNumberToLoad,
			kzgConfig.SRSOrder, // last exclusive
			kzgConfig.NumWorker,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to read trailing G2 points from %s: %v", kzgConfig.G2Path, err)
		}
	} else {
		if len(kzgConfig.G2PowerOf2Path) == 0 && len(kzgConfig.G2Path) == 0 {
			return nil, errors.New("both G2Path and G2PowerOf2Path are empty. However, object needs to load G2Points")
		}

		if len(kzgConfig.G2PowerOf2Path) != 0 {
			if kzgConfig.SRSOrder == 0 {
				return nil, errors.New("SRS order cannot be 0")
			}

			maxPower := uint64(math.Log2(float64(kzgConfig.SRSOrder)))
			_, err := ReadG2PointSection(kzgConfig.G2PowerOf2Path, 0, maxPower, 1)
			if err != nil {
				return nil, fmt.Errorf("file located at %v is invalid", kzgConfig.G2PowerOf2Path)
			}
		} else {
			return nil, fmt.Errorf(`G2PowerOf2Path is empty. However, object needs to load G2Points.
						For most operators, this likely indicates that G2_POWER_OF_2_PATH is improperly configured.`)
		}
	}

	return &SRS{
		G1:         s1,
		G2:         s2,
		G2Trailing: g2Trailing,
	}, nil
}
