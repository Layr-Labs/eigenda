package gnark

import (
	"fmt"

	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

type KzgCommitmentsGnarkBackend struct {
	Srs        kzg.SRS
	G2Trailing []bn254.G2Affine
}

func (p *KzgCommitmentsGnarkBackend) ComputeLengthProofV2(coeffs []fr.Element) (*bn254.G2Affine, error) {
	inputLength := uint32(len(coeffs))
	return p.ComputeLengthProofForLengthV2(coeffs, inputLength)
}

func (p *KzgCommitmentsGnarkBackend) ComputeLengthProofForLengthV2(
	coeffs []fr.Element, length uint32,
) (*bn254.G2Affine, error) {
	if length < uint32(len(coeffs)) {
		return nil, fmt.Errorf("length is less than the number of coefficients")
	}

	start := uint32(len(p.G2Trailing)) - length
	shiftedSecret := p.G2Trailing[start : start+uint32(len(coeffs))]
	config := ecc.MultiExpConfig{}

	//The proof of low degree is commitment of the polynomial shifted to the largest srs degree
	var lengthProof bn254.G2Affine
	_, err := lengthProof.MultiExp(shiftedSecret, coeffs, config)
	if err != nil {
		return nil, fmt.Errorf("multi exp: %w", err)
	}

	return &lengthProof, nil
}

func (p *KzgCommitmentsGnarkBackend) ComputeCommitmentV2(coeffs []fr.Element) (*bn254.G1Affine, error) {
	// compute commit for the full poly
	config := ecc.MultiExpConfig{}
	var commitment bn254.G1Affine
	_, err := commitment.MultiExp(p.Srs.G1[:len(coeffs)], coeffs, config)
	if err != nil {
		return nil, fmt.Errorf("multi exp: %w", err)
	}
	return &commitment, nil
}

func (p *KzgCommitmentsGnarkBackend) ComputeLengthCommitmentV2(coeffs []fr.Element) (*bn254.G2Affine, error) {
	config := ecc.MultiExpConfig{}

	var lengthCommitment bn254.G2Affine
	_, err := lengthCommitment.MultiExp(p.Srs.G2[:len(coeffs)], coeffs, config)
	if err != nil {
		return nil, fmt.Errorf("multi exp: %w", err)
	}
	return &lengthCommitment, nil
}
