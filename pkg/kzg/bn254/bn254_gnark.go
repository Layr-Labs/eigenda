package bn254

import (
	"log"
	"math/big"

	"github.com/consensys/gnark-crypto/ecc"
	bn "github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

func (p *G1Point) String() string {
	return StrG1(p)
}

func (p *G2Point) String() string {
	return StrG2(p)
}

var ZERO_G1 G1Point
var ZERO_G2 G2Point

var GenG1 G1Point
var GenG2 G2Point

var ZeroG1 G1Point
var ZeroG2 G2Point

//var InfG1 G1Point
//var InfG2 G2Point

func initG1G2() {

	_, _, genG1, genG2 := bn.Generators()

	GenG1 = *(*G1Point)(&genG1)
	GenG2 = *(*G2Point)(&genG2)

	//InfG1.X.SetOne()
	//InfG1.Y.SetOne()

	//ZeroG1 = G1Point(*kbls.NewG1().Zero())
	//ZeroG2 = G2Point(*kbls.NewG2().Zero())

	var g1Jac bn.G1Jac
	g1Jac.X.SetZero()
	g1Jac.Y.SetOne()
	g1Jac.Z.SetZero()

	var g1Aff bn.G1Affine
	g1Aff.FromJacobian(&g1Jac)
	ZeroG1 = *(*G1Point)(&g1Aff)

	var g2Jac bn.G2Jac
	g2Jac.X.SetZero()
	g2Jac.Y.SetOne()
	g2Jac.Z.SetZero()
	var g2Aff bn.G2Affine
	g2Aff.FromJacobian(&g2Jac)
	ZeroG2 = *(*G2Point)(&g2Aff)
}

type G1Point bn.G1Affine

// zeroes the point (like herumi BLS does with theirs). This is not co-factor clearing.
func ClearG1(x *G1Point) {
	(*bn.G1Affine)(x).Sub((*bn.G1Affine)(x), (*bn.G1Affine)(x))
}

func CopyG1(dst *G1Point, v *G1Point) {
	*dst = *v
}

func MulG1(dst *G1Point, a *G1Point, b *Fr) {
	//tmp := (kbls.Fr)(*b) // copy, we want to leave the original in mont-red form
	//(&tmp).FromRed()
	//kbls.NewG1().MulScalar((*kbls.PointG1)(dst), (*kbls.PointG1)(a), &tmp)

	var t big.Int
	(*fr.Element)(b).BigInt(&t)
	(*bn.G1Affine)(dst).ScalarMultiplication((*bn.G1Affine)(a), &t)
}

func AddG1(dst *G1Point, a *G1Point, b *G1Point) {
	var sum, bJac bn.G1Jac
	sum.FromAffine((*bn.G1Affine)(a))
	bJac.FromAffine((*bn.G1Affine)(b))
	sum.AddAssign(&bJac)
	(*bn.G1Affine)(dst).FromJacobian(&sum)
}

func SubG1(dst *G1Point, a *G1Point, b *G1Point) {
	var diff, bJac bn.G1Jac
	diff.FromAffine((*bn.G1Affine)(a))
	bJac.FromAffine((*bn.G1Affine)(b))
	diff.SubAssign(&bJac)
	(*bn.G1Affine)(dst).FromJacobian(&diff)
}

func StrG1(v *G1Point) string {
	return (*bn.G1Affine)(v).String() + "\n"
}

func NegG1(dst *G1Point) {
	// in-place should be safe here (TODO double check)
	(*bn.G1Affine)(dst).Neg((*bn.G1Affine)(dst))
}

type G2Point bn.G2Affine

func ClearG2(x *G2Point) {
	(*bn.G2Affine)(x).Sub((*bn.G2Affine)(x), (*bn.G2Affine)(x))
}

func CopyG2(dst *G2Point, v *G2Point) {
	*dst = *v
}

func MulG2(dst *G2Point, a *G2Point, b *Fr) {
	//tmp := (kbls.Fr)(*b) // copy, we want to leave the original in mont-red form
	//(&tmp).FromRed()
	//kbls.NewG2().MulScalar((*kbls.PointG2)(dst), (*kbls.PointG2)(a), &tmp)
	var t big.Int
	(*fr.Element)(b).BigInt(&t)
	(*bn.G2Affine)(dst).ScalarMultiplication((*bn.G2Affine)(a), &t)
}

func AddG2(dst *G2Point, a *G2Point, b *G2Point) {
	var sum, bJac bn.G2Jac
	sum.FromAffine((*bn.G2Affine)(a))
	bJac.FromAffine((*bn.G2Affine)(b))
	sum.AddAssign(&bJac)
	(*bn.G2Affine)(dst).FromJacobian(&sum)
}

func SubG2(dst *G2Point, a *G2Point, b *G2Point) {
	var diff, bJac bn.G2Jac
	diff.FromAffine((*bn.G2Affine)(a))
	bJac.FromAffine((*bn.G2Affine)(b))
	diff.SubAssign(&bJac)
	(*bn.G2Affine)(dst).FromJacobian(&diff)
}

func StrG2(v *G2Point) string {
	return (*bn.G2Affine)(v).String()
}

func NegG2(dst *G2Point) {
	// in-place should be safe here (TODO double check)
	(*bn.G2Affine)(dst).Neg((*bn.G2Affine)(dst))
}

func EqualG1(a *G1Point, b *G1Point) bool {
	return (*bn.G1Affine)(a).Equal((*bn.G1Affine)(b))
}

func EqualG2(a *G2Point, b *G2Point) bool {
	return (*bn.G2Affine)(a).Equal((*bn.G2Affine)(b))
}

func ToCompressedG1(p *G1Point) []byte {
	d := (*bn.G1Affine)(p).Bytes()
	return d[:]
}

func FromCompressedG1(v []byte) (*G1Point, error) {

	p := new(bn.G1Affine)
	_, err := p.SetBytes(v)

	return (*G1Point)(p), err
}

func ToCompressedG2(p *G2Point) []byte {
	//return hbls.CastToSign((*hbls.G2)(p)).Serialize()
	d := (*bn.G2Affine)(p).Bytes()
	return d[:]
}

func FromCompressedG2(v []byte) (*G2Point, error) {
	//p, err := kbls.NewG1().FromCompressed(v)

	p := new(bn.G2Affine)
	_, err := p.SetBytes(v)

	return (*G2Point)(p), err
}

func LinCombG1(numbers []G1Point, factors []Fr) *G1Point {
	out := new(G1Point)
	CopyG1(out, &ZeroG1)
	out1 := *(*bn.G1Affine)(out)

	np := make([]bn.G1Affine, len(numbers))
	fs := make([]fr.Element, len(factors))

	for i := 0; i < len(np); i++ {
		np[i] = *(*bn.G1Affine)(&numbers[i])
	}
	for i := 0; i < len(fs); i++ {
		fs[i] = *(*fr.Element)(&factors[i])
	}

	config := ecc.MultiExpConfig{}
	_, err := out1.MultiExp(np, fs, config)
	if err != nil {
		log.Fatal(err)
	}
	return (*G1Point)(&out1)
}

func LinCombG2(numbers []G2Point, factors []Fr) *G2Point {
	out := new(G2Point)
	CopyG2(out, &ZeroG2)
	out1 := *(*bn.G2Affine)(out)

	np := make([]bn.G2Affine, len(numbers))
	fs := make([]fr.Element, len(factors))

	for i := 0; i < len(np); i++ {
		np[i] = *(*bn.G2Affine)(&numbers[i])
	}
	for i := 0; i < len(fs); i++ {
		fs[i] = *(*fr.Element)(&factors[i])
	}

	config := ecc.MultiExpConfig{}
	_, err := out1.MultiExp(np, fs, config)
	if err != nil {
		log.Fatal(err)
	}
	return (*G2Point)(&out1)
}

func PairingsVerify(a1 *G1Point, a2 *G2Point, b1 *G1Point, b2 *G2Point) bool {
	//var tmp hbls.GT
	//hbls.Pairing(&tmp, (*hbls.G1)(a1), (*hbls.G2)(a2))
	////fmt.Println("tmp", tmp.GetString(10))
	//var tmp2 hbls.GT
	//hbls.Pairing(&tmp2, (*hbls.G1)(b1), (*hbls.G2)(b2))

	//// invert left pairing
	//var tmp3 hbls.GT
	//hbls.GTInv(&tmp3, &tmp)

	//// multiply the two
	//var tmp4 hbls.GT
	//hbls.GTMul(&tmp4, &tmp3, &tmp2)

	//// final exp.
	//var tmp5 hbls.GT
	//hbls.FinalExp(&tmp5, &tmp4)

	// = 1_T
	//return tmp5.IsOne()

	// TODO, alternatively use the equal check (faster or slower?):
	////fmt.Println("tmp2", tmp2.GetString(10))
	//return tmp.IsEqual(&tmp2)

	var negB1 bn.G1Affine
	negB1.Neg((*bn.G1Affine)(b1))

	P := [2]bn.G1Affine{*(*bn.G1Affine)(a1), negB1}
	Q := [2]bn.G2Affine{*(*bn.G2Affine)(a2), *(*bn.G2Affine)(b2)}

	ok, err := bn.PairingCheck(P[:], Q[:])
	if err != nil {
		log.Fatal(err)
	}
	return ok
}

func Generators() (G1Point, G2Point) {
	_, _, g1, g2 := bn.Generators()
	return *(*G1Point)(&g1), *(*G2Point)(&g2)
}
