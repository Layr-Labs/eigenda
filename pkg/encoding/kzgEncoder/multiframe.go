package kzgEncoder

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"

	rs "github.com/Layr-Labs/eigenda/pkg/encoding/encoder"
	bls "github.com/Layr-Labs/eigenda/pkg/kzg/bn254"
)

type Sample struct {
	Commitment bls.G1Point
	Proof      bls.G1Point
	Row        int
	Coeffs     []bls.Fr
	X          uint // X is int , at which index is evaluated
}

// generate a random value using Fiat Shamir transform
func GenRandomness(params rs.EncodingParams, samples []Sample, m int) (bls.Fr, error) {

	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	err := enc.Encode(samples)
	if err != nil {
		return bls.ZERO, err
	}

	err = enc.Encode(params)
	if err != nil {
		return bls.ZERO, err
	}

	err = enc.Encode(m)
	if err != nil {
		return bls.ZERO, err
	}

	var randomFr bls.Fr

	err = bls.HashToSingleField(&randomFr, buffer.Bytes())
	if err != nil {
		return bls.ZERO, err
	}
	return randomFr, nil
}

// m is number of blob
func (group *KzgEncoderGroup) UniversalVerify(params rs.EncodingParams, samples []Sample, m int) error {
	verifier, _ := group.GetKzgVerifier(params)
	ks := verifier.Ks

	D := params.ChunkLen

	n := len(samples)
	fmt.Printf("Batch verify %v frames of %v symbols out of %v blobs \n", n, params.ChunkLen, m)

	r, err := GenRandomness(params, samples, m)
	if err != nil {
		return err
	}

	randomsFr := make([]bls.Fr, n)
	bls.CopyFr(&randomsFr[0], &r)

	// lhs
	var tmp bls.Fr

	// power of r
	for j := 0; j < n-1; j++ {
		bls.MulModFr(&randomsFr[j+1], &randomsFr[j], &r)
	}

	// array of proofs
	proofs := make([]bls.G1Point, n)
	for i := 0; i < n; i++ {
		bls.CopyG1(&proofs[i], &samples[i].Proof)
	}

	// lhs g1
	lhsG1 := bls.LinCombG1(proofs, randomsFr)

	// lhs g2
	lhsG2 := &ks.Srs.G2[D]

	// rhs g2
	rhsG2 := &bls.GenG2

	// rhs g1
	// get commitments
	commits := make([]bls.G1Point, m)
	// get coeffs
	ftCoeffs := make([]bls.Fr, m)
	for k := 0; k < n; k++ {
		s := samples[k]
		row := s.Row
		bls.AddModFr(&ftCoeffs[row], &ftCoeffs[row], &randomsFr[k])
		bls.CopyG1(&commits[row], &s.Commitment)
	}

	ftG1 := bls.LinCombG1(commits, ftCoeffs)

	// second term
	stCoeffs := make([]bls.Fr, D)
	for k := 0; k < n; k++ {
		coeffs := samples[k].Coeffs

		rk := randomsFr[k]
		for j := uint64(0); j < D; j++ {
			bls.MulModFr(&tmp, &coeffs[j], &rk)
			bls.AddModFr(&stCoeffs[j], &stCoeffs[j], &tmp)
		}
	}
	stG1 := bls.LinCombG1(ks.Srs.G1[:D], stCoeffs)

	// third term
	ttCoeffs := make([]bls.Fr, n)

	// get leading coset powers
	leadingDs := make([]bls.Fr, n)

	for k := 0; k < n; k++ {
		x, err := rs.GetLeadingCosetIndex(
			uint64(samples[k].X),
			params.NumChunks,
		)
		if err != nil {
			return err
		}

		h := ks.ExpandedRootsOfUnity[x]
		var hPow bls.Fr
		bls.CopyFr(&hPow, &bls.ONE)

		for j := uint64(0); j < D; j++ {
			bls.MulModFr(&tmp, &hPow, &h)
			bls.CopyFr(&hPow, &tmp)
		}
		bls.CopyFr(&leadingDs[k], &hPow)
	}

	//
	for k := 0; k < n; k++ {
		rk := randomsFr[k]
		bls.MulModFr(&ttCoeffs[k], &rk, &leadingDs[k])
	}
	ttG1 := bls.LinCombG1(proofs, ttCoeffs)

	var rhsG1 bls.G1Point
	bls.SubG1(&rhsG1, ftG1, stG1)
	bls.AddG1(&rhsG1, &rhsG1, ttG1)

	if bls.PairingsVerify(lhsG1, lhsG2, &rhsG1, rhsG2) {
		return nil
	} else {
		return errors.New("Universal Verify Incorrect paring")
	}
}
