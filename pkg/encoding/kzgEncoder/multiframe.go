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

	//rInt := uint64(22894)
	//var r bls.Fr
	//bls.AsFr(&r, rInt)

	r, err := GenRandomness(params, samples, m)
	if err != nil {
		return err
	}

	randomsFr := make([]bls.Fr, n)
	//bls.AsFr(&randomsFr[0], rInt)
	bls.CopyFr(&randomsFr[0], &r)

	fmt.Println("random", r.String())

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

	fmt.Printf("len proof %v len ran %v\n", len(proofs), len(randomsFr))
	// lhs g1
	lhsG1 := bls.LinCombG1(proofs, randomsFr)

	// lhs g2
	lhsG2 := &ks.Srs.G2[D]

	// rhs g2
	rhsG2 := &bls.GenG2

	// rhs g1
	// get commitments
	commits := make([]bls.G1Point, m)
	//for k := 0 ; k < n ; k++ {
	// commits[k] = samples[k].Commitment
	//}
	// get coeffs
	ftCoeffs := make([]bls.Fr, m)
	for k := 0; k < n; k++ {
		s := samples[k]
		row := s.Row
		bls.AddModFr(&ftCoeffs[row], &ftCoeffs[row], &randomsFr[k])
		bls.CopyG1(&commits[row], &s.Commitment)
	}
	fmt.Printf("len commit %v len coeff %v\n", len(commits), len(ftCoeffs))

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

//func SingleVerify(ks *kzg.KZGSettings, commitment *bls.G1Point, x *bls.Fr, coeffs []bls.Fr, proof bls.G1Point) bool {
//	var xPow bls.Fr
//	bls.CopyFr(&xPow, &bls.ONE)

//	var tmp bls.Fr
//	for i := 0; i < len(coeffs); i++ {
//		bls.MulModFr(&tmp, &xPow, x)
//		bls.CopyFr(&xPow, &tmp)
//	}

// [x^n]_2
//	var xn2 bls.G2Point
//	bls.MulG2(&xn2, &bls.GenG2, &xPow)

// [s^n - x^n]_2
//	var xnMinusYn bls.G2Point
//	bls.SubG2(&xnMinusYn, &ks.Srs.G2[len(coeffs)], &xn2)

// [interpolation_polynomial(s)]_1
//	is1 := bls.LinCombG1(ks.Srs.G1[:len(coeffs)], coeffs)
// [commitment - interpolation_polynomial(s)]_1 = [commit]_1 - [interpolation_polynomial(s)]_1
//	var commitMinusInterpolation bls.G1Point
//	bls.SubG1(&commitMinusInterpolation, commitment, is1)

//	return bls.PairingsVerify(&commitMinusInterpolation, &bls.GenG2, &proof, &xnMinusYn)
//}
