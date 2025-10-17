package kzg

import (
	"github.com/consensys/gnark-crypto/ecc/bn254"
)

func init() {
	initG1G2()
}

var GenG1 bn254.G1Affine
var GenG2 bn254.G2Affine

var ZeroG1 bn254.G1Affine
var ZeroG2 bn254.G2Affine

func initG1G2() {

	_, _, GenG1, GenG2 = bn254.Generators()

	var g1Jac bn254.G1Jac
	g1Jac.X.SetZero()
	g1Jac.Y.SetOne()
	g1Jac.Z.SetZero()

	var g1Aff bn254.G1Affine
	g1Aff.FromJacobian(&g1Jac)
	ZeroG1 = g1Aff

	var g2Jac bn254.G2Jac
	g2Jac.X.SetZero()
	g2Jac.Y.SetOne()
	g2Jac.Z.SetZero()
	ZeroG2.FromJacobian(&g2Jac)
}
