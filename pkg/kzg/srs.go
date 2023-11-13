package kzg

import (
	"errors"

	bls "github.com/Layr-Labs/eigenda/pkg/kzg/bn254"
)

type SRS struct {

	// [b.multiply(b.G1, pow(s, i, MODULUS)) for i in range(WIDTH+1)],
	G1 []bls.G1Point
	// [b.multiply(b.G2, pow(s, i, MODULUS)) for i in range(WIDTH+1)],
	G2 []bls.G2Point
}

func NewSrs(G1 []bls.G1Point, G2 []bls.G2Point) (*SRS, error) {

	if len(G1) != len(G2) {
		return nil, errors.New("secret list lengths don't match")
	}
	return &SRS{
		G1: G1,
		G2: G2,
	}, nil
}
