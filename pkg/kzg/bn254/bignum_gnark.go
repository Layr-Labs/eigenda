package bn254

import (
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

type Fr fr.Element

func init() {
	initGlobals()
	ClearG1(&ZERO_G1)
	ClearG2(&ZERO_G2)
	initG1G2()
}

func (fr *Fr) String() string {
	return FrStr(fr)
}

func SetFr(dst *Fr, v string) {
	// Bad practice to ignore error!
	_, _ = (*fr.Element)(dst).SetString(v)
}

func FrToBytes(src *Fr) [32]byte {
	return (*fr.Element)(src).Bytes()
}

func FrSetBytes(dst *Fr, v []byte) {
	(*fr.Element)(dst).SetBytes(v)
}

func FrFrom32(dst *Fr, v [32]byte) (ok bool) {
	(*fr.Element)(dst).SetBytes(v[:])
	return true
}

func FrTo32(src *Fr) [32]byte {
	return (*fr.Element)(src).Bytes()
}

func CopyFr(dst *Fr, v *Fr) {
	*dst = *v
}

func AsFr(dst *Fr, i uint64) {
	//var data [8]byte
	//binary.BigEndian.PutUint64(data[:], i)
	//(*kbls.Fr)(dst).SetBytes(data[:])
	(*fr.Element)(dst).SetUint64(i)
}

func FrStr(b *Fr) string {
	if b == nil {
		return "<nil>"
	}
	return (*fr.Element)(b).String()
}

func EqualOne(v *Fr) bool {
	return (*fr.Element)(v).IsOne()
}

func EqualZero(v *Fr) bool {
	return (*fr.Element)(v).IsZero()
}

func EqualFr(a *Fr, b *Fr) bool {
	return (*fr.Element)(a).Equal((*fr.Element)(b))
}

func SubModFr(dst *Fr, a, b *Fr) {
	(*fr.Element)(dst).Sub((*fr.Element)(a), (*fr.Element)(b))
}

func AddModFr(dst *Fr, a, b *Fr) {
	(*fr.Element)(dst).Add((*fr.Element)(a), (*fr.Element)(b))
}

func DivModFr(dst *Fr, a, b *Fr) {
	(*fr.Element)(dst).Div((*fr.Element)(a), (*fr.Element)(b))
}

func MulModFr(dst *Fr, a, b *Fr) {
	(*fr.Element)(dst).Mul((*fr.Element)(a), (*fr.Element)(b))
}

func InvModFr(dst *Fr, v *Fr) {
	(*fr.Element)(dst).Inverse((*fr.Element)(v))
}

func EvalPolyAt(dst *Fr, p []Fr, x *Fr) {
	EvalPolyAtUnoptimized(dst, p, x)
}
