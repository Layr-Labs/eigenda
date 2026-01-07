package verifier

import (
	"fmt"
	"math"
	"math/big"

	eigenbn254 "github.com/Layr-Labs/eigenda/crypto/ecc/bn254"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/v2/fft"
	"github.com/Layr-Labs/eigenda/encoding/v2/kzg"
	"github.com/Layr-Labs/eigenda/encoding/v2/rs"
	"github.com/Layr-Labs/eigenda/resources/srs"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

type ParametrizedVerifier struct {
	g1SRS kzg.G1SRS
	Fs    *fft.FFTSettings
}

// VerifyFrame verifies a single frame against a commitment.
// If needing to verify multiple frames of the same chunk length, prefer [Verifier.UniversalVerify].
func (v *ParametrizedVerifier) verifyFrame(
	frame *encoding.Frame, frameIndex uint64, commitment *bn254.G1Affine, numChunks uint64,
) error {
	j, err := rs.GetLeadingCosetIndex(frameIndex, numChunks)
	if err != nil {
		return fmt.Errorf("GetLeadingCosetIndex: %w", err)
	}

	exponent := uint64(math.Log2(float64(len(frame.Coeffs))))
	G2atD := srs.G2PowerOf2SRS[exponent]

	err = verifyFrame(frame, v.g1SRS, commitment, &v.Fs.ExpandedRootsOfUnity[j], &G2atD)
	if err != nil {
		return fmt.Errorf("VerifyFrame: %w", err)
	}
	return nil
}

// Verify function assumes the Data stored is coefficients of coset's interpolating poly
func verifyFrame(
	frame *encoding.Frame, g1SRS kzg.G1SRS, commitment *bn254.G1Affine, x *fr.Element, g2Atn *bn254.G2Affine,
) error {
	var xPow fr.Element
	xPow.SetOne()

	for i := 0; i < len(frame.Coeffs); i++ {
		xPow.Mul(&xPow, x)
	}

	var xPowBigInt big.Int

	// [x^n]_2
	var xn2 bn254.G2Affine

	xn2.ScalarMultiplication(&kzg.GenG2, xPow.BigInt(&xPowBigInt))

	// [s^n - x^n]_2
	var xnMinusYn bn254.G2Affine
	xnMinusYn.Sub(g2Atn, &xn2)

	// [interpolation_polynomial(s)]_1
	var is1 bn254.G1Affine
	config := ecc.MultiExpConfig{}
	_, err := is1.MultiExp(g1SRS[:len(frame.Coeffs)], frame.Coeffs, config)
	if err != nil {
		return fmt.Errorf("MultiExp: %w", err)
	}

	// [commitment - interpolation_polynomial(s)]_1 = [commit]_1 - [interpolation_polynomial(s)]_1
	var commitMinusInterpolation bn254.G1Affine
	commitMinusInterpolation.Sub(commitment, &is1)

	// Verify the pairing equation
	//
	// e([commitment - interpolation_polynomial(s)], [1]) = e([proof],  [s^n - x^n])
	//    equivalent to
	// e([commitment - interpolation_polynomial]^(-1), [1]) * e([proof],  [s^n - x^n]) = 1_T
	//

	err = eigenbn254.PairingsVerify(&commitMinusInterpolation, &kzg.GenG2, &frame.Proof, &xnMinusYn)
	if err != nil {
		return fmt.Errorf("verify pairing: %w", err)
	}
	return nil
}
