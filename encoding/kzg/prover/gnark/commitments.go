package gnark

import (
	"fmt"

	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

type KzgCommitmentsGnarkBackend struct {
	KzgConfig  *kzg.KzgConfig
	Srs        *kzg.SRS
	G2Trailing []bn254.G2Affine
}

func (p *KzgCommitmentsGnarkBackend) ComputeLengthProof(coeffs []fr.Element) (*bn254.G2Affine, error) {
	inputLength := uint64(len(coeffs))
	return p.ComputeLengthProofForLength(coeffs, inputLength)
}

func (p *KzgCommitmentsGnarkBackend) ComputeLengthProofForLength(coeffs []fr.Element, length uint64) (*bn254.G2Affine, error) {
	if length < uint64(len(coeffs)) {
		return nil, fmt.Errorf("length is less than the number of coefficients")
	}

	start := p.KzgConfig.SRSNumberToLoad - length
	shiftedSecret := p.G2Trailing[start : start+uint64(len(coeffs))]
	config := ecc.MultiExpConfig{}

	//The proof of low degree is commitment of the polynomial shifted to the largest srs degree
	var lengthProof bn254.G2Affine
	_, err := lengthProof.MultiExp(shiftedSecret, coeffs, config)
	if err != nil {
		return nil, err
	}

	return &lengthProof, nil
}

func (p *KzgCommitmentsGnarkBackend) ComputeCommitment(coeffs []fr.Element) (*bn254.G1Affine, error) {
	// compute commit for the full poly
	config := ecc.MultiExpConfig{}
	var commitment bn254.G1Affine
	_, err := commitment.MultiExp(p.Srs.G1[:len(coeffs)], coeffs, config)
	if err != nil {
		return nil, err
	}
	return &commitment, nil
}

func (p *KzgCommitmentsGnarkBackend) ComputeLengthCommitment(coeffs []fr.Element) (*bn254.G2Affine, error) {
	config := ecc.MultiExpConfig{}

	var lengthCommitment bn254.G2Affine
	_, err := lengthCommitment.MultiExp(p.Srs.G2[:len(coeffs)], coeffs, config)
	if err != nil {
		return nil, err
	}
	return &lengthCommitment, nil
}
