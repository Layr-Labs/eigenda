// Original: https://github.com/ethereum/research/blob/master/kzg_data_availability/kzg_proofs.py

//go:build !bignum_pure && !bignum_hol256
// +build !bignum_pure,!bignum_hol256

package kzg

import (
	"log"

	bls "github.com/Layr-Labs/eigenda/pkg/kzg/bn254"
)

// Compute KZG proof for polynomial in coefficient form at positions x * w^y where w is
// an n-th root of unity (this is the proof for one data availability sample, which consists
// of several polynomial evaluations)
func (ks *KZGSettings) ComputeProofMulti(poly []bls.Fr, x uint64, n uint64) *bls.G1Point {
	// divisor = [-pow(x, n, MODULUS)] + [0] * (n - 1) + [1]
	divisor := make([]bls.Fr, n+1)
	var xFr bls.Fr
	bls.AsFr(&xFr, x)
	// TODO: inefficient, could use squaring, or maybe BLS lib offers a power method?
	// TODO: for small ranges, maybe compute pow(x, n, mod) in uint64?
	var xPowN, tmp bls.Fr
	for i := uint64(0); i < n; i++ {
		bls.MulModFr(&tmp, &xPowN, &xFr)
		bls.CopyFr(&xPowN, &tmp)
	}

	// -pow(x, n, MODULUS)
	bls.SubModFr(&divisor[0], &bls.ZERO, &xPowN)
	// [0] * (n - 1)
	for i := uint64(1); i < n; i++ {
		bls.CopyFr(&divisor[i], &bls.ZERO)
	}
	// 1
	bls.CopyFr(&divisor[n], &bls.ONE)

	// quot = poly / divisor
	quotientPolynomial := PolyLongDiv(poly, divisor[:])
	//for i := 0; i < len(quotientPolynomial); i++ {
	//	fmt.Printf("quot poly %d: %s\n", i, FrStr(&quotientPolynomial[i]))
	//}

	// evaluate quotient poly at shared secret, in G1
	return bls.LinCombG1(ks.Srs.G1[:len(quotientPolynomial)], quotientPolynomial)
}

// Check a proof for a KZG commitment for an evaluation f(x w^i) = y_i
// The ys must have a power of 2 length
func (ks *KZGSettings) CheckProofMulti(commitment *bls.G1Point, proof *bls.G1Point, x *bls.Fr, ys []bls.Fr) (bool, error) {
	// Interpolate at a coset. Note because it is a coset, not the subgroup, we have to multiply the
	// polynomial coefficients by x^i
	interpolationPoly, err := ks.FFT(ys, true)
	if err != nil {
		log.Println("ys is bad, cannot compute FFT", err)
		return false, err
	}
	// TODO: can probably be optimized
	// apply div(c, pow(x, i, MODULUS)) to every coeff c in interpolationPoly
	// x^0 at first, then up to x^n
	var xPow bls.Fr
	bls.CopyFr(&xPow, &bls.ONE)

	var tmp, tmp2 bls.Fr
	for i := 0; i < len(interpolationPoly); i++ {
		bls.InvModFr(&tmp, &xPow)
		bls.MulModFr(&tmp2, &interpolationPoly[i], &tmp)
		bls.CopyFr(&interpolationPoly[i], &tmp2)
		bls.MulModFr(&tmp, &xPow, x)
		bls.CopyFr(&xPow, &tmp)
	}
	// [x^n]_2
	var xn2 bls.G2Point
	bls.MulG2(&xn2, &bls.GenG2, &xPow)
	// [s^n - x^n]_2
	var xnMinusYn bls.G2Point
	bls.SubG2(&xnMinusYn, &ks.Srs.G2[len(ys)], &xn2)

	// fmt.Println("CheckMultiProof")
	// for i := 0; i < len(interpolationPoly); i++ {
	// 	fmt.Println(i, interpolationPoly[i].String())
	// }

	// [interpolation_polynomial(s)]_1
	is1 := bls.LinCombG1(ks.Srs.G1[:len(interpolationPoly)], interpolationPoly)
	// [commitment - interpolation_polynomial(s)]_1 = [commit]_1 - [interpolation_polynomial(s)]_1
	var commitMinusInterpolation bls.G1Point
	bls.SubG1(&commitMinusInterpolation, commitment, is1)

	// Verify the pairing equation
	//
	// e([commitment - interpolation_polynomial(s)], [1]) = e([proof],  [s^n - x^n])
	//    equivalent to
	// e([commitment - interpolation_polynomial]^(-1), [1]) * e([proof],  [s^n - x^n]) = 1_T
	//

	return bls.PairingsVerify(&commitMinusInterpolation, &bls.GenG2, proof, &xnMinusYn), nil
}
