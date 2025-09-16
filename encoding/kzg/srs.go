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
	_ "github.com/Layr-Labs/eigenda/encoding"
	"github.com/consensys/gnark-crypto/ecc/bn254"
)

type G1SRS []bn254.G1Affine

type SRS struct {
	// G1 points are used to:
	// 1. On prover (in encoder): generate blob commitments (multiproofs are generated using SRSTables).
	// 2. On prover (in proxy/client): generate blob commitments.
	// 3. On verifier: verify blob multiproofs using initial chunk-length number of G1 points.
	// 4. On verifier: verify length proofs using trailing G1 points.
	//
	// [b.multiply(b.G1, pow(s, i, MODULUS)) for i in range(WIDTH+1)],
	G1 []bn254.G1Affine
	// G2 points are used to:
	// 1. On prover (in encoder): generate length commitments and proofs (see [encoding.BlobCommitments]).
	// 2. On prover (in proxy/client): generate length commitments and length proofs.
	// 3. On verifier: verify blob multiproofs using 28 powerOf2 G2 points.
	//
	// [b.multiply(b.G2, pow(s, i, MODULUS)) for i in range(WIDTH+1)],
	G2 []bn254.G2Affine
}

func NewSrs(G1 []bn254.G1Affine, G2 []bn254.G2Affine) SRS {
	return SRS{G1: G1, G2: G2}
}
