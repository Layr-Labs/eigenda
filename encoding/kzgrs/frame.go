package kzgrs

import (
	"bytes"
	"encoding/gob"

	kzg "github.com/Layr-Labs/eigenda/pkg/kzg"
	bls "github.com/Layr-Labs/eigenda/pkg/kzg/bn254"
)

// Proof is the multireveal proof
// Coeffs is identical to input data converted into Fr element
type Frame struct {
	Proof  bls.G1Point
	Coeffs []bls.Fr
}

func (f *Frame) Encode() ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(f)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func Decode(b []byte) (Frame, error) {
	var f Frame
	buf := bytes.NewBuffer(b)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&f)
	if err != nil {
		return Frame{}, err
	}
	return f, nil
}

// Verify function assumes the Data stored is coefficients of coset's interpolating poly
func (f *Frame) Verify(ks *kzg.KZGSettings, commitment *bls.G1Point, x *bls.Fr, g2Atn *bls.G2Point) bool {
	var xPow bls.Fr
	bls.CopyFr(&xPow, &bls.ONE)

	var tmp bls.Fr
	for i := 0; i < len(f.Coeffs); i++ {
		bls.MulModFr(&tmp, &xPow, x)
		bls.CopyFr(&xPow, &tmp)
	}

	// [x^n]_2
	var xn2 bls.G2Point
	bls.MulG2(&xn2, &bls.GenG2, &xPow)

	// [s^n - x^n]_2
	var xnMinusYn bls.G2Point

	bls.SubG2(&xnMinusYn, g2Atn, &xn2)

	// [interpolation_polynomial(s)]_1
	is1 := bls.LinCombG1(ks.Srs.G1[:len(f.Coeffs)], f.Coeffs)

	// [commitment - interpolation_polynomial(s)]_1 = [commit]_1 - [interpolation_polynomial(s)]_1
	var commitMinusInterpolation bls.G1Point
	bls.SubG1(&commitMinusInterpolation, commitment, is1)

	// Verify the pairing equation
	//
	// e([commitment - interpolation_polynomial(s)], [1]) = e([proof],  [s^n - x^n])
	//    equivalent to
	// e([commitment - interpolation_polynomial]^(-1), [1]) * e([proof],  [s^n - x^n]) = 1_T
	//

	return bls.PairingsVerify(&commitMinusInterpolation, &bls.GenG2, &f.Proof, &xnMinusYn)
}
