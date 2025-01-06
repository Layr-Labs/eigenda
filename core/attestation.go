package core

import (
	"crypto/rand"
	"math/big"

	bn254utils "github.com/Layr-Labs/eigenda/core/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fp"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
)

type G1Point struct {
	*bn254.G1Affine
}

func newFpElement(x *big.Int) fp.Element {
	var p fp.Element
	p.SetBigInt(x)
	return p
}

func NewG1Point(x, y *big.Int) *G1Point {
	return &G1Point{
		&bn254.G1Affine{
			X: newFpElement(x),
			Y: newFpElement(y),
		},
	}
}

// Add another G1 point to this one
func (p *G1Point) Add(p2 *G1Point) {
	p.G1Affine.Add(p.G1Affine, p2.G1Affine)
}

// Sub another G1 point from this one
func (p *G1Point) Sub(p2 *G1Point) {
	p.G1Affine.Sub(p.G1Affine, p2.G1Affine)
}

// VerifyEquivalence verifies G1Point is equivalent the G2Point
func (p *G1Point) VerifyEquivalence(p2 *G2Point) (bool, error) {
	return bn254utils.CheckG1AndG2DiscreteLogEquality(p.G1Affine, p2.G2Affine)
}

func (p *G1Point) SerializeCompressed() [32]byte {
	return p.Bytes()
}

func (p *G1Point) Serialize() []byte {
	res := p.RawBytes()
	return res[:]
}

func (p *G1Point) Deserialize(data []byte) (*G1Point, error) {
	var point bn254.G1Affine
	_, err := point.SetBytes(data)
	if err != nil {
		return nil, err
	}
	return &G1Point{&point}, nil
}

func (p *G1Point) Clone() *G1Point {
	return &G1Point{&bn254.G1Affine{
		X: newFpElement(p.X.BigInt(new(big.Int))),
		Y: newFpElement(p.Y.BigInt(new(big.Int))),
	}}
}

func (p *G1Point) Hash() [32]byte {
	return crypto.Keccak256Hash(p.Serialize())
}

type G2Point struct {
	*bn254.G2Affine
}

// Add another G2 point to this one
func (p *G2Point) Add(p2 *G2Point) {
	p.G2Affine.Add(p.G2Affine, p2.G2Affine)
}

// Sub another G2 point from this one
func (p *G2Point) Sub(p2 *G2Point) {
	p.G2Affine.Sub(p.G2Affine, p2.G2Affine)
}

func (p *G2Point) Serialize() []byte {
	res := p.RawBytes()
	return res[:]
}

func (p *G2Point) Deserialize(data []byte) (*G2Point, error) {
	var point bn254.G2Affine
	_, err := point.SetBytes(data)
	if err != nil {
		return nil, err
	}
	return &G2Point{&point}, nil
}

func (p *G2Point) Clone() *G2Point {
	return &G2Point{&bn254.G2Affine{
		X: struct {
			A0, A1 fp.Element
		}{
			A0: newFpElement(p.X.A0.BigInt(new(big.Int))),
			A1: newFpElement(p.X.A1.BigInt(new(big.Int))),
		},
		Y: struct {
			A0, A1 fp.Element
		}{
			A0: newFpElement(p.Y.A0.BigInt(new(big.Int))),
			A1: newFpElement(p.Y.A1.BigInt(new(big.Int))),
		},
	}}
}

type Signature struct {
	*G1Point
}

// Verify a message against a G2 public key
func (s *Signature) Verify(pubkey *G2Point, message [32]byte) bool {
	ok, err := bn254utils.VerifySig(s.G1Affine, pubkey.G2Affine, message)
	if err != nil {
		return false
	}
	return ok
}

// GetOperatorID hashes the G1Point (public key of an operator) to generate the operator ID.
// It does it to match how it's hashed in solidity: `keccak256(abi.encodePacked(pk.X, pk.Y))`
// Ref: https://github.com/Layr-Labs/eigenlayer-contracts/blob/avs-unstable/src/contracts/libraries/BN254.sol#L285
func (p *G1Point) GetOperatorID() OperatorID {
	x := p.X.BigInt(new(big.Int))
	y := p.Y.BigInt(new(big.Int))
	return OperatorID(crypto.Keccak256Hash(append(math.U256Bytes(x), math.U256Bytes(y)...)))
}

type PrivateKey = fr.Element

type KeyPair struct {
	PrivKey *PrivateKey
	PubKey  *G1Point
}

func MakeKeyPair(sk *PrivateKey) *KeyPair {
	pk := bn254utils.MulByGeneratorG1(sk)
	return &KeyPair{sk, &G1Point{pk}}
}

func MakeKeyPairFromString(sk string) (*KeyPair, error) {
	ele, err := new(fr.Element).SetString(sk)
	if err != nil {
		return nil, err
	}
	return MakeKeyPair(ele), nil
}

func GenRandomBlsKeys() (*KeyPair, error) {

	//Max random value is order of the curve
	max := new(big.Int)
	max.SetString(fr.Modulus().String(), 10)

	//Generate cryptographically strong pseudo-random between 0 - max
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return nil, err
	}

	sk := new(PrivateKey).SetBigInt(n)
	return MakeKeyPair(sk), nil
}

func (k *KeyPair) SignMessage(message [32]byte) *Signature {
	H := bn254utils.MapToCurve(message)
	sig := new(bn254.G1Affine).ScalarMultiplication(H, k.PrivKey.BigInt(new(big.Int)))
	return &Signature{&G1Point{sig}}
}

func (k *KeyPair) SignHashedToCurveMessage(g1HashedMsg *G1Point) *Signature {
	sig := new(bn254.G1Affine).ScalarMultiplication(g1HashedMsg.G1Affine, k.PrivKey.BigInt(new(big.Int)))
	return &Signature{&G1Point{sig}}
}

func (k *KeyPair) GetPubKeyG2() *G2Point {
	return &G2Point{bn254utils.MulByGeneratorG2(k.PrivKey)}
}

func (k *KeyPair) GetPubKeyG1() *G1Point {
	return k.PubKey
}

// MakePubkeyRegistrationData returns the data that should be sent to the pubkey compendium smart contract to register the public key.
// The values returned constitute a proof that the operator knows the secret key corresponding to the public key, and prevents the operator
// from attacking the signature protocol by registering a public key that is derived from other public keys.
// (e.g., see https://medium.com/@coolcottontail/rogue-key-attack-in-bls-signature-and-harmony-security-eac1ea2370ee)
func (k *KeyPair) MakePubkeyRegistrationData(operatorAddress common.Address) *G1Point {
	return &G1Point{bn254utils.MakePubkeyRegistrationData(k.PrivKey, operatorAddress)}

}
