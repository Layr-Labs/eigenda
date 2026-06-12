package bn254_test

import (
	"testing"

	"github.com/Layr-Labs/eigenda/core/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/ethereum/go-ethereum/common"
)

func TestMapToCurve(t *testing.T) {
	msg := [32]byte{}
	copy(msg[:], "test message")

	point := bn254.MapToCurve(msg)
	if !point.IsOnCurve() {
		t.Errorf("MapToCurve failed: point is not on the curve")
	}
}

func TestCheckG1AndG2DiscreteLogEquality(t *testing.T) {
	privKey := new(fr.Element).SetUint64(12345) // 示例私钥
	pointG1 := bn254.MulByGeneratorG1(privKey)
	pointG2 := bn254.MulByGeneratorG2(privKey)

	ok, err := bn254.CheckG1AndG2DiscreteLogEquality(pointG1, pointG2)
	if err != nil {
		t.Fatalf("CheckG1AndG2DiscreteLogEquality returned an error: %v", err)
	}
	if !ok {
		t.Errorf("CheckG1AndG2DiscreteLogEquality failed: expected true, got false")
	}
}

func TestGetG1Generator(t *testing.T) {
	gen := bn254.GetG1Generator()
	if !gen.IsOnCurve() {
		t.Errorf("GetG1Generator failed: generator is not on the curve")
	}
}

func TestGetG2Generator(t *testing.T) {
	gen := bn254.GetG2Generator()
	if !gen.IsOnCurve() {
		t.Errorf("GetG2Generator failed: generator is not on the curve")
	}
}

func TestMulByGeneratorG1(t *testing.T) {
	privKey := new(fr.Element).SetUint64(12345) // 示例私钥
	point := bn254.MulByGeneratorG1(privKey)

	if !point.IsOnCurve() {
		t.Errorf("MulByGeneratorG1 failed: point is not on the curve")
	}
}

func TestMulByGeneratorG2(t *testing.T) {
	privKey := new(fr.Element).SetUint64(12345) // 示例私钥
	point := bn254.MulByGeneratorG2(privKey)

	if !point.IsOnCurve() {
		t.Errorf("MulByGeneratorG2 failed: point is not on the curve")
	}
}

func TestMakePubkeyRegistrationData(t *testing.T) {
	privKey := new(fr.Element).SetUint64(12345) // 示例私钥
	operatorAddress := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")

	pubkey := bn254.MakePubkeyRegistrationData(privKey, operatorAddress)

	if !pubkey.IsOnCurve() {
		t.Errorf("MakePubkeyRegistrationData failed: public key is not on the curve")
	}
}
