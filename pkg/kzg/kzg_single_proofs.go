// Original: https://github.com/ethereum/research/blob/master/kzg_data_availability/kzg_proofs.py

//go:build !bignum_pure && !bignum_hol256
// +build !bignum_pure,!bignum_hol256

package kzg

import (
	bls "github.com/Layr-Labs/eigenda/pkg/kzg/bn254"
)

// KZG commitment to polynomial in evaluation form, i.e. eval = FFT(coeffs).
// The eval length must match the prepared KZG settings width.
func CommitToEvalPoly(secretG1IFFT []bls.G1Point, eval []bls.Fr) *bls.G1Point {
	return bls.LinCombG1(secretG1IFFT, eval)
}

// KZG commitment to polynomial in coefficient form
func (ks *KZGSettings) CommitToPoly(coeffs []bls.Fr) *bls.G1Point {
	return bls.LinCombG1(ks.Srs.G1[:len(coeffs)], coeffs)
}

// KZG commitment to polynomial in coefficient form, unoptimized version
// func (ks *KZGSettings) CommitToPolyUnoptimized(coeffs []bls.Fr) *bls.G1Point {
// 	// Do so by computing the linear combination with the shared secret.
// 	var out bls.G1Point
// 	bls.ClearG1(&out)
// 	var tmp, tmp2 bls.G1Point
// 	for i := 0; i < len(coeffs); i++ {
// 		bls.MulG1(&tmp, &ks.Srs.G1[i], &coeffs[i])
// 		bls.AddG1(&tmp2, &out, &tmp)
// 		bls.CopyG1(&out, &tmp2)
// 	}
// 	return &out
// }

// Compute KZG proof for polynomial in coefficient form at position x
func (ks *KZGSettings) ComputeProofSingle(poly []bls.Fr, x uint64) *bls.G1Point {
	// divisor = [-x, 1]
	divisor := [2]bls.Fr{}
	var tmp bls.Fr
	bls.AsFr(&tmp, x)
	bls.SubModFr(&divisor[0], &bls.ZERO, &tmp)
	bls.CopyFr(&divisor[1], &bls.ONE)
	//for i := 0; i < 2; i++ {
	//	fmt.Printf("div poly %d: %s\n", i, FrStr(&divisor[i]))
	//}
	// quot = poly / divisor
	quotientPolynomial := PolyLongDiv(poly, divisor[:])
	//for i := 0; i < len(quotientPolynomial); i++ {
	//	fmt.Printf("quot poly %d: %s\n", i, FrStr(&quotientPolynomial[i]))
	//}

	// evaluate quotient poly at shared secret, in G1
	return bls.LinCombG1(ks.Srs.G1[:len(quotientPolynomial)], quotientPolynomial)
}

// Compute KZG proof for polynomial in coefficient form at position x
func (ks *KZGSettings) ComputeProofSingleAtFr(poly []bls.Fr, x bls.Fr) *bls.G1Point {
	// divisor = [-x, 1]
	divisor := [2]bls.Fr{}
	bls.SubModFr(&divisor[0], &bls.ZERO, &x)
	bls.CopyFr(&divisor[1], &bls.ONE)
	//for i := 0; i < 2; i++ {
	//	fmt.Printf("div poly %d: %s\n", i, FrStr(&divisor[i]))
	//}
	// quot = poly / divisor
	quotientPolynomial := PolyLongDiv(poly, divisor[:])
	//for i := 0; i < len(quotientPolynomial); i++ {
	//	fmt.Printf("quot poly %d: %s\n", i, FrStr(&quotientPolynomial[i]))
	//}

	// evaluate quotient poly at shared secret, in G1
	return bls.LinCombG1(ks.Srs.G1[:len(quotientPolynomial)], quotientPolynomial)
}

// To prove a list of field elements are a polynomial with commitment C,
// we open C up at hash(hash(poly), C). then the verifier verifies the opening proof
// and evaluates the polynomial at that point. This requires 2 pairings and 1 G1
// subtraction and n field multiplications and additions as opposed to n G1 additions
// func (ks *KZGSettings) ComputePolynomialEquivalenceProofInG2(poly []bls.Fr, commit bls.G1Point) *bls.G2Point {
// 	r := GetEquivalenceProofChallenge(poly, commit)
// 	rFr := new(bls.Fr)
// 	bls.FrSetBytes(rFr, r)
// 	// eval := new(bls.Fr)
// 	// bls.EvalPolyAt(eval, poly, rFr)
// 	// fmt.Println(eval)
// 	return ks.ComputeProofSingleAtFrInG2(poly, *rFr)
// }

// To prove a list of field elements are a polynomial with commitment C,
// we open C up at hash(hash(poly), C). then the verifier verifies the opening proof
// and evaluates the polynomial at that point. This requires 2 pairings and 1 G1
// subtraction and n field multiplications and additions as opposed to n G1 additions
// technique described in section 3.1 of https://eprint.iacr.org/2019/953.pdf
// func (ks *KZGSettings) ComputeBatchPolynomialEquivalenceProofInG2(polys [][]bls.Fr, commits []bls.G1Point) *bls.G2Point {
// 	rs := make([]byte, 0)
// 	for i := 0; i < len(polys); i++ {
// 		rs = append(rs, GetEquivalenceProofChallenge(polys[i], commits[i])...)
// 	}

// 	r := crypto.Keccak256(rs)
// 	rFr := new(bls.Fr)
// 	bls.FrSetBytes(rFr, r)

// 	gammaBytes := crypto.Keccak256(append(rs, bytes32Zero...))
// 	gamma := new(bls.Fr)
// 	bls.FrSetBytes(gamma, gammaBytes)
// 	gammaPower := new(bls.Fr)
// 	bls.CopyFr(gammaPower, gamma)

// 	gammaShiftedEquivalenceProofs := ks.ComputeProofSingleAtFrInG2(polys[0], *rFr)
// 	// gammaShiftedCommits := commits[0]
// 	// gammaShiftedEvals := new(bls.Fr)
// 	// bls.EvalPolyAt(gammaShiftedEvals, polys[0], rFr)

// 	for i := 1; i < len(polys); i++ {
// 		equivalenceProof := ks.ComputeProofSingleAtFrInG2(polys[i], *rFr)
// 		bls.MulG2(equivalenceProof, equivalenceProof, gammaPower)
// 		bls.AddG2(gammaShiftedEquivalenceProofs, gammaShiftedEquivalenceProofs, equivalenceProof)

// 		// commitTmp := new(bls.G1Point)
// 		// bls.MulG1(commitTmp, &commits[i], gammaPower)
// 		// bls.AddG1(&gammaShiftedCommits, &gammaShiftedCommits, commitTmp)

// 		// eval := new(bls.Fr)
// 		// bls.EvalPolyAt(eval, polys[i], rFr)
// 		// bls.MulModFr(eval, eval, gammaPower)
// 		// bls.AddModFr(gammaShiftedEvals, gammaShiftedEvals, eval)
// 		// fmt.Println("gammashifted", i, gamma, gammaPower)
// 		bls.MulModFr(gammaPower, gammaPower, gamma)
// 	}
// 	// fmt.Println("gammashifted", len(polys), gammaPower)

// 	// fmt.Println(ks.CheckProofSingleProofInG2(&gammaShiftedCommits, gammaShiftedEquivalenceProofs, rFr, gammaShiftedEvals))

// 	return gammaShiftedEquivalenceProofs
// }

// func GetEquivalenceProofChallenge(poly []bls.Fr, commit bls.G1Point) []byte {
// 	polyBytes := make([]byte, 32*len(poly))
// 	for i := 0; i < len(poly); i++ {
// 		coeffBytes := bls.FrToBytes(&poly[i])
// 		for j := 0; j < 32; j++ {
// 			polyBytes[32*i+j] = coeffBytes[j]
// 		}
// 	}
// 	r := crypto.Keccak256(polyBytes)
// 	commitX := commit.X.Bytes()
// 	commitY := commit.Y.Bytes()
// 	for i := 0; i < 32; i++ {
// 		r = append(r, commitX[i])
// 	}
// 	for i := 0; i < 32; i++ {
// 		r = append(r, commitY[i])
// 	}
// 	r = crypto.Keccak256(r)
// 	return r
// }

// Compute KZG proof for polynomial in coefficient form at position x
// func (ks *KZGSettings) ComputeProofSingleInG2(poly []bls.Fr, x uint64) *bls.G2Point {
// 	// divisor = [-x, 1]
// 	divisor := [2]bls.Fr{}
// 	var tmp bls.Fr
// 	bls.AsFr(&tmp, x)
// 	bls.SubModFr(&divisor[0], &bls.ZERO, &tmp)
// 	bls.CopyFr(&divisor[1], &bls.ONE)

// 	quotientPolynomial := polyLongDiv(poly, divisor[:])

// 	// evaluate quotient poly at shared secret, in G2
// 	return bls.LinCombG2(ks.Srs.G2[:len(quotientPolynomial)], quotientPolynomial)
// }

// Compute KZG proof for polynomial in coefficient form at position x
// func (ks *KZGSettings) ComputeProofSingleAtFrInG2(poly []bls.Fr, x bls.Fr) *bls.G2Point {
// 	// divisor = [-x, 1]
// 	divisor := [2]bls.Fr{}
// 	bls.SubModFr(&divisor[0], &bls.ZERO, &x)
// 	bls.CopyFr(&divisor[1], &bls.ONE)

// 	quotientPolynomial := polyLongDiv(poly, divisor[:])

// 	// fmt.Println("divisor")
// 	// fmt.Println(ks.CommitToPoly(divisor[:]))
// 	// for i := 0; i < len(divisor); i++ {
// 	// 	fmt.Println(divisor[i].String())
// 	// }

// 	// fmt.Println("quotient")
// 	// for i := 0; i < len(quotientPolynomial); i++ {
// 	// 	fmt.Println(quotientPolynomial[i].String())
// 	// }

// 	// evaluate quotient poly at shared secret, in G1
// 	pi := bls.LinCombG2(ks.Srs.G2[:len(quotientPolynomial)], quotientPolynomial)
// 	// y := new(bls.Fr)
// 	// bls.EvalPolyAt(y, poly, &x)

// 	// ok := ks.CheckProofSingleProofInG2(ks.CommitToPoly(poly), pi, &x, y)

// 	// fmt.Println(ok)

// 	return pi
// }

// Check a proof for a KZG commitment for an evaluation f(x) = y
func (ks *KZGSettings) CheckProofSingle(commitment *bls.G1Point, proof *bls.G1Point, x *bls.Fr, y *bls.Fr) bool {
	// Verify the pairing equation
	var xG2 bls.G2Point
	bls.MulG2(&xG2, &bls.GenG2, x)
	var sMinuxX bls.G2Point
	bls.SubG2(&sMinuxX, &ks.Srs.G2[1], &xG2)
	var yG1 bls.G1Point
	bls.MulG1(&yG1, &bls.GenG1, y)
	var commitmentMinusY bls.G1Point
	bls.SubG1(&commitmentMinusY, commitment, &yG1)

	// This trick may be applied in the BLS-lib specific code:
	//
	// e([commitment - y], [1]) = e([proof],  [s - x])
	//    equivalent to
	// e([commitment - y]^(-1), [1]) * e([proof],  [s - x]) = 1_T
	//
	return bls.PairingsVerify(&commitmentMinusY, &bls.GenG2, proof, &sMinuxX)
}

// Check a proof for a KZG commitment for an evaluation f(x) = y
// func (ks *KZGSettings) CheckProofSingleProofInG2(commitment *bls.G1Point, proof *bls.G2Point, x *bls.Fr, y *bls.Fr) bool {
// 	// Verify the pairing equation
// 	var xG1 bls.G1Point
// 	bls.MulG1(&xG1, &bls.GenG1, x)
// 	var sMinuxX bls.G1Point
// 	bls.SubG1(&sMinuxX, &ks.Srs.G1[1], &xG1)
// 	var yG1 bls.G1Point
// 	bls.MulG1(&yG1, &bls.GenG1, y)
// 	var commitmentMinusY bls.G1Point
// 	bls.SubG1(&commitmentMinusY, commitment, &yG1)

// 	fmt.Println("fminusfaG1Aff")
// 	fmt.Println(commitmentMinusY.String())

// 	fmt.Println("xminusaG1Aff")
// 	fmt.Println(sMinuxX.String())

// 	fmt.Println("proof")
// 	fmt.Println(proof.String())

// 	// This trick may be applied in the BLS-lib specific code:
// 	//
// 	// e([commitment - y], [1]) = e([proof],  [s - x])
// 	//    equivalent to
// 	// e([commitment - y]^(-1), [1]) * e([proof],  [s - x]) = 1_T
// 	//
// 	return bls.PairingsVerify(&commitmentMinusY, &bls.GenG2, &sMinuxX, proof)
// }
