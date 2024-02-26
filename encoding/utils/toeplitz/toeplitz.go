package toeplitz

import (
	"errors"
	"log"

	"github.com/Layr-Labs/eigenda/encoding/fft"

	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

// V is ordered as (v_0, .., v_6), so it creates a
// matrix below. Slice must be odd
// v_0 v_1 v_2 v_3
// v_6 v_0 v_1 v_2
// v_5 v_6 v_0 v_1
// v_4 v_5 v_6 v_0
type Toeplitz struct {
	V  []fr.Element
	Fs *fft.FFTSettings
}

func NewToeplitz(v []fr.Element, fs *fft.FFTSettings) (*Toeplitz, error) {
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
func (t *Toeplitz) Multiply(x []fr.Element) ([]fr.Element, error) {
	cv := t.ExtendCircularVec()

	rv := t.FromColVToRowV(cv)
	cir := NewCircular(rv, t.Fs)

	xE := make([]fr.Element, len(cv))
	for i := 0; i < len(x); i++ {
		xE[i].Set(&x[i])
	}
	for i := len(x); i < len(cv); i++ {
		xE[i].SetZero()
	}

	product, err := cir.Multiply(xE)
	if err != nil {
		return nil, err
	}

	return product[:len(x)], nil
}

// Take FFT on Toeplitz vector, coefficient is used for computing hadamard product
// but carried with multi scalar multiplication
func (t *Toeplitz) GetFFTCoeff() ([]fr.Element, error) {
	cv := t.ExtendCircularVec()

	rv := t.FromColVToRowV(cv)
	cir := NewCircular(rv, t.Fs)

	return cir.GetFFTCoeff()
}

func (t *Toeplitz) GetCoeff() ([]fr.Element, error) {
	cv := t.ExtendCircularVec()

	rv := t.FromColVToRowV(cv)
	cir := NewCircular(rv, t.Fs)

	return cir.GetCoeff()
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

func (t *Toeplitz) ExtendCircularVec() []fr.Element {
	E := make([]fr.Element, len(t.V)+1) // extra 1 from extended, equal to 2*dimE
	numRow := t.GetMatDim()
	E[0].Set(&t.V[0])

	for i := 1; i < numRow; i++ {

		E[i].Set(&t.V[len(t.V)-i])
	}

	// assign some value to the extra dimension
	E[numRow].SetZero()

	// numRow == numCol
	for i := 1; i < numRow; i++ {
		E[numRow+i].Set(&t.V[numRow-i])
	}

	return E
}

// if   col Vector is [v_0, v_1, v_2, v_3, 0, v_4, v_5, v_6]
// then row Vector is [v_0, v_6, v_5, v_4, 0, v_3, v_2, v_1]
// this operation is involutory. i.e. f(f(v)) = v

func (t *Toeplitz) FromColVToRowV(cv []fr.Element) []fr.Element {
	n := len(cv)
	rv := make([]fr.Element, n)

	rv[0].Set(&cv[0])

	for i := 1; i < n; i++ {

		rv[i].Set(&cv[n-i])
	}

	return rv
}

func (t *Toeplitz) GetMatDim() int {
	return (len(t.V) + 1) / 2
}

// naive implementation for multiplication. Used for testing
func (t *Toeplitz) DirectMultiply(x []fr.Element) []fr.Element {
	numCol := t.GetMatDim()

	n := len(t.V)

	out := make([]fr.Element, numCol)
	for i := 0; i < numCol; i++ {
		var sum fr.Element
		sum.SetZero()

		for j := 0; j < numCol; j++ {
			idx := (j - i + n) % n
			var product fr.Element
			product.Mul(&t.V[idx], &x[j])

			sum.Add(&product, &sum)

		}

		out[i].Set(&sum)
	}

	return out
}
