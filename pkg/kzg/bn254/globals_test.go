package bn254

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO: slow test
// func TestRootOfUnity(t *testing.T) {
// 	one := ToFr("1")

// 	for i := 0; i < len(Scale2RootOfUnity); i++ {
// 		u := Scale2RootOfUnity[i]
// 		var product Fr
// 		CopyFr(&product, &one)

// 		//log.Println("multiply of ", i*i)
// 		order := int(math.Pow(2.0, float64(i)))
// 		for j := 0; j < order; j++ {
// 			MulModFr(&product, &product, &u)
// 		}

// 		if !EqualOne(&product) {
// 			log.Fatalf("%v point is not a root of unity", i)
// 		}
// 	}

// 	log.Printf("all %v root of unity pass", len(Scale2RootOfUnity))
// }

func TestIsPowerOf2(t *testing.T) {
	if !IsPowerOfTwo(2) {
		t.Fatal("2 is not a power of 2")
	}
}

func TestEvalPolyAtUnoptimized(t *testing.T) {
	one := ToFr("1")
	coeffs := []Fr{}
	res := ONE
	EvalPolyAtUnoptimized(&res, coeffs, &one)
	assert.Equal(t, res, ZERO)

	coeffs = []Fr{
		ToFr("1"),
		ToFr("2"),
		ToFr("3"),
	}
	res = ZERO
	EvalPolyAtUnoptimized(&res, coeffs, &one)
	assert.Equal(t, res, ToFr("6"))

	zero := ZERO
	EvalPolyAtUnoptimized(&res, coeffs, &zero)
	assert.Equal(t, res, coeffs[0])
}
