package kzgEncoder

import (
	"fmt"
	"log"

	rs "github.com/Layr-Labs/eigenda/pkg/encoding/encoder"
	kzg "github.com/Layr-Labs/eigenda/pkg/kzg"
	bls "github.com/Layr-Labs/eigenda/pkg/kzg/bn254"
)

type Sample struct {
	Commitment bls.G1Point
	Proof      bls.G1Point
	Row        int
	Coeffs     []bls.Fr
	X          uint // X is int , at which index is evaluated
}

// m is number of blob
func (group *KzgEncoderGroup) UniversalVerify(params rs.EncodingParams, samples []Sample, m int) bool {
	verifier, _ := group.GetKzgVerifier(params)
	ks := verifier.Ks

	for ind, s := range samples {
		q, err := rs.GetLeadingCosetIndex(
			uint64(s.X),
			params.NumChunks,
		)
		if err != nil {
			return false
		}

		lc := ks.FFTSettings.ExpandedRootsOfUnity[uint64(q)]

		ok := SingleVerify(ks, &s.Commitment, &lc, s.Coeffs, s.Proof)
		if !ok {
			fmt.Println("proof", s.Proof.String())
			fmt.Println("commitment", s.Commitment.String())

			for i := 0; i < len(s.Coeffs); i++ {
				fmt.Printf("%v ", s.Coeffs[i].String())
			}
			fmt.Println("q", q, lc.String())

			log.Fatalf("Proof %v failed\n", ind)
		} else {

			fmt.Println("&&&&&&&&&&&&&&&&&&tested frame and pass", ind)
		}
	}

	D := len(samples[0].Coeffs) // chunkLen

	n := len(samples)

	rInt := uint64(22894)
	var r bls.Fr
	bls.AsFr(&r, rInt)

	randomsFr := make([]bls.Fr, n)
	bls.AsFr(&randomsFr[0], rInt)

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
		for j := 0; j < D; j++ {
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
			return false
		}

		h := ks.ExpandedRootsOfUnity[x]
		var hPow bls.Fr
		bls.CopyFr(&hPow, &bls.ONE)

		for j := 0; j < D; j++ {
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

	return bls.PairingsVerify(lhsG1, lhsG2, &rhsG1, rhsG2)
}

func SingleVerify(ks *kzg.KZGSettings, commitment *bls.G1Point, x *bls.Fr, coeffs []bls.Fr, proof bls.G1Point) bool {
	var xPow bls.Fr
	bls.CopyFr(&xPow, &bls.ONE)

	var tmp bls.Fr
	for i := 0; i < len(coeffs); i++ {
		bls.MulModFr(&tmp, &xPow, x)
		bls.CopyFr(&xPow, &tmp)
	}

	// [x^n]_2
	var xn2 bls.G2Point
	bls.MulG2(&xn2, &bls.GenG2, &xPow)

	// [s^n - x^n]_2
	var xnMinusYn bls.G2Point
	bls.SubG2(&xnMinusYn, &ks.Srs.G2[len(coeffs)], &xn2)

	// [interpolation_polynomial(s)]_1
	is1 := bls.LinCombG1(ks.Srs.G1[:len(coeffs)], coeffs)
	// [commitment - interpolation_polynomial(s)]_1 = [commit]_1 - [interpolation_polynomial(s)]_1
	var commitMinusInterpolation bls.G1Point
	bls.SubG1(&commitMinusInterpolation, commitment, is1)

	// Verify the pairing equation
	//
	// e([commitment - interpolation_polynomial(s)], [1]) = e([proof],  [s^n - x^n])
	//    equivalent to
	// e([commitment - interpolation_polynomial]^(-1), [1]) * e([proof],  [s^n - x^n]) = 1_T
	//

	return bls.PairingsVerify(&commitMinusInterpolation, &bls.GenG2, &proof, &xnMinusYn)
}
