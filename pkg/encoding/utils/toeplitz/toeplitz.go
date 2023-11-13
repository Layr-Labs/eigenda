package toeplitz

import (
	"errors"
	"log"
	// "fmt"

	kzg "github.com/Layr-Labs/eigenda/pkg/kzg"
	bls "github.com/Layr-Labs/eigenda/pkg/kzg/bn254"
)

// V is ordered as (v_0, .., v_6), so it creates a
// matrix below. Slice must be odd
// v_0 v_1 v_2 v_3
// v_6 v_0 v_1 v_2
// v_5 v_6 v_0 v_1
// v_4 v_5 v_6 v_0
type Toeplitz struct {
	V  []bls.Fr
	Fs *kzg.FFTSettings
}

func NewToeplitz(v []bls.Fr, fs *kzg.FFTSettings) (*Toeplitz, error) {
	if len(v)%2 != 1 {
		log.Println("num diagonal vector must be odd")
		return nil, errors.New("num diagonal vector must be odd")
	}

	return &Toeplitz{
		V:  v,
		Fs: fs,
	}, nil
}

// Mutliplication of a matrix and a vector, both elements are Fr Element
// good info about FFT and Toeplitz,
// https://alinush.github.io/2020/03/19/multiplying-a-vector-by-a-toeplitz-matrix.html
func (t *Toeplitz) Multiply(x []bls.Fr) ([]bls.Fr, error) {
	cv := t.ExtendCircularVec()

	rv := t.FromColVToRowV(cv)
	cir := NewCircular(rv, t.Fs)

	xE := make([]bls.Fr, len(cv))
	for i := 0; i < len(x); i++ {
		bls.CopyFr(&xE[i], &x[i])
	}
	for i := len(x); i < len(cv); i++ {
		bls.CopyFr(&xE[i], &bls.ZERO)
	}

	product, err := cir.Multiply(xE)
	if err != nil {
		return nil, err
	}

	return product[:len(x)], nil
}

// Mutliplication of a matrix and a vector, where the matrix contains Fr elements, vectors are
// G1 Points. It supports precomputed
func (t *Toeplitz) MultiplyPoints(x []bls.G1Point, inv bool, usePrecompute bool) ([]bls.G1Point, error) {
	cv := t.ExtendCircularVec()

	rv := t.FromColVToRowV(cv)
	cir := NewCircular(rv, t.Fs)

	//for i := range x {
	//fmt.Printf("%v", x[i].String())
	//}
	//fmt.Println("x", len(x))
	//fmt.Println("vc", len(cv))

	return cir.MultiplyPoints(x, inv, true)
}

// Take FFT on Toeplitz vector, coefficient is used for computing hadamard product
// but carried with multi scalar multiplication
func (t *Toeplitz) GetFFTCoeff() ([]bls.Fr, error) {
	cv := t.ExtendCircularVec()

	rv := t.FromColVToRowV(cv)
	cir := NewCircular(rv, t.Fs)

	return cir.GetFFTCoeff()
}

// Expand toeplitz matrix into circular matrix
// the outcome is a also concise representation
// if   V is (v_0, v_1, v_2, v_3, v_4, v_5, v_6)
// then E is (v_0, v_6, v_5, v_4, 0,   v_3, v_2, v_1)
// representing
// [v_0, v_6, v_5, v_4, 0  , v_3, v_2, v_1 ]
// [v_1, v_0, v_6, v_5, v_4, 0  , v_3, v_2 ]
// [v_2, v_1, v_0, v_6, v_5, v_4, 0  , v_3 ]
// [v_3, v_2, v_1, v_0, v_6, v_5, v_4, 0   ]
// [0  , v_3, v_2, v_1, v_0, v_6, v_5, v_4 ]
// [v_4, 0  , v_3, v_2, v_1, v_0, v_6, v_5 ]
// [v_5, v_4, 0  , v_3, v_2, v_1, v_0, v_6 ]
// [v_6, v_5, v_4, 0  , v_3, v_2, v_1, v_0 ]

func (t *Toeplitz) ExtendCircularVec() []bls.Fr {
	E := make([]bls.Fr, len(t.V)+1) // extra 1 from extended, equal to 2*dimE
	numRow := t.GetMatDim()

	bls.CopyFr(&E[0], &t.V[0])

	for i := 1; i < numRow; i++ {
		bls.CopyFr(&E[i], &t.V[len(t.V)-i])
	}

	// assign some value to the extra dimension
	bls.CopyFr(&E[numRow], &bls.ZERO)

	// numRow == numCol
	for i := 1; i < numRow; i++ {
		bls.CopyFr(&E[numRow+i], &t.V[numRow-i])
	}

	return E
}

// if   col Vector is [v_0, v_1, v_2, v_3, 0, v_4, v_5, v_6]
// then row Vector is [v_0, v_6, v_5, v_4, 0, v_3, v_2, v_1]
// this operation is involutory. i.e. f(f(v)) = v

func (t *Toeplitz) FromColVToRowV(cv []bls.Fr) []bls.Fr {
	n := len(cv)
	rv := make([]bls.Fr, n)
	bls.CopyFr(&rv[0], &cv[0])

	for i := 1; i < n; i++ {
		bls.CopyFr(&rv[i], &cv[n-i])
	}

	return rv
}

func (t *Toeplitz) GetMatDim() int {
	return (len(t.V) + 1) / 2
}

// naive implementation for multiplication. Used for testing
func (t *Toeplitz) DirectMultiply(x []bls.Fr) []bls.Fr {
	numCol := t.GetMatDim()

	n := len(t.V)

	out := make([]bls.Fr, numCol)
	for i := 0; i < numCol; i++ {
		var sum bls.Fr
		bls.CopyFr(&sum, &bls.ZERO)
		for j := 0; j < numCol; j++ {
			idx := (j - i + n) % n
			var product bls.Fr
			bls.MulModFr(&product, &t.V[idx], &x[j])
			bls.AddModFr(&sum, &product, &sum)
		}
		bls.CopyFr(&out[i], &sum)
	}

	return out
}
