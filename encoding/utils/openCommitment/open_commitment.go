package openCommitment

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

// Implement https://github.com/ethereum/consensus-specs/blob/017a8495f7671f5fff2075a9bfc9238c1a0982f8/specs/deneb/polynomial-commitments.md#compute_kzg_proof_impl
func ComputeKzgProof(
	evalFr []fr.Element,
	index int,
	G1srsLagrange []bn254.G1Affine,
	rootOfUnities []fr.Element,
) (*bn254.G1Affine, *fr.Element, error) {
	if len(evalFr) != len(rootOfUnities) {
		return nil, nil, fmt.Errorf("inconsistent length between blob and root of unities")
	}
	if index < 0 || index >= len(evalFr) {
		return nil, nil, fmt.Errorf("the function only opens points within a blob")
	}

	polyShift := make([]fr.Element, len(evalFr))

	valueFr := evalFr[index]

	zFr := rootOfUnities[index]

	for i := 0; i < len(polyShift); i++ {
		polyShift[i].Sub(&evalFr[i], &valueFr)
	}

	denomPoly := make([]fr.Element, len(rootOfUnities))

	for i := 0; i < len(evalFr); i++ {
		denomPoly[i].Sub(&rootOfUnities[i], &zFr)
	}

	quotientPoly := make([]fr.Element, len(rootOfUnities))
	for i := 0; i < len(quotientPoly); i++ {
		if denomPoly[i].IsZero() {
			quotientPoly[i] = computeQuotientEvalOnDomain(zFr, evalFr, valueFr, rootOfUnities)
		} else {
			quotientPoly[i].Div(&polyShift[i], &denomPoly[i])
		}
	}

	config := ecc.MultiExpConfig{}

	var proof bn254.G1Affine
	_, err := proof.MultiExp(G1srsLagrange, quotientPoly, config)
	if err != nil {
		return nil, nil, err
	}

	return &proof, &valueFr, nil
}

func VerifyKzgProof(G1Gen, commitment, proof bn254.G1Affine, G2Gen, G2tau bn254.G2Affine, valueFr, zFr fr.Element) error {

	var valueG1 bn254.G1Affine
	var valueBig big.Int
	valueG1.ScalarMultiplication(&G1Gen, valueFr.BigInt(&valueBig))

	var commitMinusValue bn254.G1Affine
	commitMinusValue.Sub(&commitment, &valueG1)

	var zG2 bn254.G2Affine
	zG2.ScalarMultiplication(&G2Gen, zFr.BigInt(&valueBig))

	var xMinusZ bn254.G2Affine
	xMinusZ.Sub(&G2tau, &zG2)

	return PairingsVerify(&commitMinusValue, &G2Gen, &proof, &xMinusZ)
}

func PairingsVerify(a1 *bn254.G1Affine, a2 *bn254.G2Affine, b1 *bn254.G1Affine, b2 *bn254.G2Affine) error {
	var negB1 bn254.G1Affine
	negB1.Neg(b1)

	P := [2]bn254.G1Affine{*a1, negB1}
	Q := [2]bn254.G2Affine{*a2, *b2}

	ok, err := bn254.PairingCheck(P[:], Q[:])
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("pairingCheck pairing not ok")
	}

	return nil
}

func CommitInLagrange(evalFr []fr.Element, G1srsLagrange []bn254.G1Affine) (*bn254.G1Affine, error) {
	config := ecc.MultiExpConfig{}

	var proof bn254.G1Affine
	_, err := proof.MultiExp(G1srsLagrange, evalFr, config)
	if err != nil {
		return nil, err
	}
	return &proof, nil
}

// Implement https://github.com/ethereum/consensus-specs/blob/017a8495f7671f5fff2075a9bfc9238c1a0982f8/specs/deneb/polynomial-commitments.md#compute_quotient_eval_within_domain
func computeQuotientEvalOnDomain(zFr fr.Element, evalFr []fr.Element, valueFr fr.Element, rootOfunities []fr.Element) fr.Element {
	var quotient fr.Element
	var f_i, numerator, denominator, temp fr.Element

	for i := 0; i < len(rootOfunities); i++ {
		omega_i := rootOfunities[i]
		if omega_i.Equal(&zFr) {
			continue
		}

		f_i.Sub(&evalFr[i], &valueFr)
		numerator.Mul(&f_i, &omega_i)

		denominator.Sub(&zFr, &omega_i)
		denominator.Mul(&denominator, &zFr)

		numerator.Mul(&f_i, &omega_i)
		temp.Div(&numerator, &denominator)

		quotient.Add(&quotient, &temp)

	}
	return quotient
}
