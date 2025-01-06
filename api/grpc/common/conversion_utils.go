package common

import (
	"fmt"
	verifierBindings "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDABlobVerifier"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"math/big"
)

// ToBinding converts a BlobCommitment into a contractEigenDABlobVerifier.BlobCommitment
func (c *BlobCommitment) ToBinding() (*verifierBindings.BlobCommitment, error) {
	convertedCommitment, err := BytesToBN254G1Point(c.GetCommitment())
	if err != nil {
		return nil, fmt.Errorf("convert commitment to g1 point: %s", err)
	}

	convertedLengthCommitment, err := BytesToBN254G2Point(c.GetLengthCommitment())
	if err != nil {
		return nil, fmt.Errorf("convert length commitment to g2 point: %s", err)
	}

	convertedLengthProof, err := BytesToBN254G2Point(c.GetLengthProof())
	if err != nil {
		return nil, fmt.Errorf("convert length proof to g2 point: %s", err)
	}

	return &verifierBindings.BlobCommitment{
		Commitment:       *convertedCommitment,
		LengthCommitment: *convertedLengthCommitment,
		LengthProof:      *convertedLengthProof,
		DataLength:       c.GetLength(),
	}, nil
}

// BytesToBN254G1Point accepts a byte array, and converts it into a contractEigenDABlobVerifier.BN254G1Point
func BytesToBN254G1Point(bytes []byte) (*verifierBindings.BN254G1Point, error) {
	var g1Point bn254.G1Affine
	_, err := g1Point.SetBytes(bytes)

	if err != nil {
		return nil, fmt.Errorf("deserialize g1 point: %s", err)
	}

	return &verifierBindings.BN254G1Point{
		X: g1Point.X.BigInt(new(big.Int)),
		Y: g1Point.Y.BigInt(new(big.Int)),
	}, nil
}

// BytesToBN254G2Point accepts a byte array, and converts it into a contractEigenDABlobVerifier.BN254G2Point
func BytesToBN254G2Point(bytes []byte) (*verifierBindings.BN254G2Point, error) {
	var g2Point bn254.G2Affine
	_, err := g2Point.SetBytes(bytes)

	if err != nil {
		return nil, fmt.Errorf("deserialize g2 point: %s", err)
	}

	var x, y [2]*big.Int
	x[0] = g2Point.X.A0.BigInt(new(big.Int))
	x[1] = g2Point.X.A1.BigInt(new(big.Int))

	y[0] = g2Point.Y.A0.BigInt(new(big.Int))
	y[1] = g2Point.Y.A1.BigInt(new(big.Int))

	return &verifierBindings.BN254G2Point{
		X: x,
		Y: y,
	}, nil
}
