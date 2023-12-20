// This code is sourced from the go-kzg Repository by protolambda.
// Original code: https://github.com/protolambda/go-kzg
// MIT License
//
// Copyright (c) 2020 @protolambda
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

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

func HashToSingleField(dst *Fr, msg []byte) error {
	DST := []byte("-")
	randomFr, err := fr.Hash(msg, DST, 1)
	randomFrBytes := (randomFr[0]).Bytes()
	FrSetBytes(dst, randomFrBytes[:])
	return err
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
