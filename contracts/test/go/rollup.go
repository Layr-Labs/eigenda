package main

import (
	"fmt"
	"math/big"

	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

func main() {
	g2Gen := GetG2Generator()

	//powers of tau where s = 2
	tau := make([]*bn254.G2Affine, 5)
	tau[0] = new(bn254.G2Affine).ScalarMultiplication(g2Gen, big.NewInt(1))
	tau[1] = new(bn254.G2Affine).ScalarMultiplication(g2Gen, big.NewInt(2))
	tau[2] = new(bn254.G2Affine).ScalarMultiplication(g2Gen, big.NewInt(4))
	tau[3] = new(bn254.G2Affine).ScalarMultiplication(g2Gen, big.NewInt(8))
	tau[4] = new(bn254.G2Affine).ScalarMultiplication(g2Gen, big.NewInt(16))

	//polynomial coefficients for x = 6
	poly := make([]*big.Int, 4)
	poly[0] = new(fr.Element).SetInt64(259).BigInt(new(big.Int))
	poly[1] = new(fr.Element).SetInt64(43).BigInt(new(big.Int))
	poly[2] = new(fr.Element).SetInt64(7).BigInt(new(big.Int))
	poly[3] = new(fr.Element).SetInt64(1).BigInt(new(big.Int))

	//G2 proof for illegal value at x = 6
	proof := new(bn254.G2Affine)
	for i := 0; i < 4; i++ {
		proof = new(bn254.G2Affine).Add(proof, new(bn254.G2Affine).ScalarMultiplication(tau[i], poly[i]))
	}

	fmt.Println(proof)
}

func GetG2Generator() *bn254.G2Affine {
	g2Gen := new(bn254.G2Affine)
	g2Gen.X.SetString("10857046999023057135944570762232829481370756359578518086990519993285655852781",
		"11559732032986387107991004021392285783925812861821192530917403151452391805634")
	g2Gen.Y.SetString("8495653923123431417604973247489272438418190587263600148770280649306958101930",
		"4082367875863433681332203403145435568316851327593401208105741076214120093531")
	return g2Gen
}
