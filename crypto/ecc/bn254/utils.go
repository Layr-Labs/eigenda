package bn254

import (
	"math/big"

	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fp"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func VerifySig(sig *bn254.G1Affine, pubkey *bn254.G2Affine, msgBytes [32]byte) (bool, error) {

	g2Gen := GetG2Generator()

	msgPoint := MapToCurve(msgBytes)

	var negSig bn254.G1Affine
	negSig.Neg((*bn254.G1Affine)(sig))

	P := [2]bn254.G1Affine{*msgPoint, negSig}
	Q := [2]bn254.G2Affine{*pubkey, *g2Gen}

	ok, err := bn254.PairingCheck(P[:], Q[:])
	if err != nil {
		return false, nil
	}
	return ok, nil

}

func MapToCurve(digest [32]byte) *bn254.G1Affine {

	one := new(big.Int).SetUint64(1)
	three := new(big.Int).SetUint64(3)
	x := new(big.Int)
	x.SetBytes(digest[:])
	for {
		// y = x^3 + 3
		xP3 := new(big.Int).Exp(x, big.NewInt(3), fp.Modulus())
		y := new(big.Int).Add(xP3, three)
		y.Mod(y, fp.Modulus())

		if y.ModSqrt(y, fp.Modulus()) == nil {
			x.Add(x, one).Mod(x, fp.Modulus())
		} else {
			var fpX, fpY fp.Element
			fpX.SetBigInt(x)
			fpY.SetBigInt(y)
			return &bn254.G1Affine{
				X: fpX,
				Y: fpY,
			}
		}
	}
}

func CheckG1AndG2DiscreteLogEquality(pointG1 *bn254.G1Affine, pointG2 *bn254.G2Affine) (bool, error) {
	negGenG1 := new(bn254.G1Affine).Neg(GetG1Generator())
	return bn254.PairingCheck([]bn254.G1Affine{*pointG1, *negGenG1}, []bn254.G2Affine{*GetG2Generator(), *pointG2})
}

func GetG1Generator() *bn254.G1Affine {
	g1Gen := new(bn254.G1Affine)
	_, err := g1Gen.X.SetString("1")
	if err != nil {
		return nil
	}
	_, err = g1Gen.Y.SetString("2")
	if err != nil {
		return nil
	}
	return g1Gen
}

func GetG2Generator() *bn254.G2Affine {
	g2Gen := new(bn254.G2Affine)
	g2Gen.X.SetString("10857046999023057135944570762232829481370756359578518086990519993285655852781",
		"11559732032986387107991004021392285783925812861821192530917403151452391805634")
	g2Gen.Y.SetString("8495653923123431417604973247489272438418190587263600148770280649306958101930",
		"4082367875863433681332203403145435568316851327593401208105741076214120093531")
	return g2Gen
}

func MulByGeneratorG1(a *fr.Element) *bn254.G1Affine {
	g1Gen := GetG1Generator()
	return new(bn254.G1Affine).ScalarMultiplication(g1Gen, a.BigInt(new(big.Int)))
}

func MulByGeneratorG2(a *fr.Element) *bn254.G2Affine {
	g2Gen := GetG2Generator()
	return new(bn254.G2Affine).ScalarMultiplication(g2Gen, a.BigInt(new(big.Int)))
}

func MakePubkeyRegistrationData(privKey *fr.Element, operatorAddress common.Address) *bn254.G1Affine {
	toHash := make([]byte, 0)
	toHash = append(toHash, crypto.Keccak256([]byte("BN254PubkeyRegistration(address operator)"))...)
	toHash = append(toHash, operatorAddress.Bytes()...)

	msgHash := crypto.Keccak256(toHash)
	// convert to [32]byte
	var msgHash32 [32]byte
	copy(msgHash32[:], msgHash)

	// hash to G1
	hashToSign := MapToCurve(msgHash32)

	return new(bn254.G1Affine).ScalarMultiplication(hashToSign, privKey.BigInt(new(big.Int)))
}
