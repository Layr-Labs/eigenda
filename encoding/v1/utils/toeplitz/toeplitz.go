// toeplitz package is outdated, and only kept around for v1 prover.
// prover v2 replaces this implementation with an inlined version
// that does a lot less needless allocations and copies.
// See getSlicesCoeff in encoding/v2/kzg/prover/gnark/multiframe_proof.go
package toeplitz

import (
	"errors"
	"fmt"
	"log"

	"github.com/Layr-Labs/eigenda/encoding/v1/fft"

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

// Take FFT on Toeplitz vector, coefficient is used for computing hadamard product
// but carried with multi scalar multiplication
// Returns a slice of size 2*dimE
func (t *Toeplitz) GetFFTCoeff() ([]fr.Element, error) {
	cv := t.extendCirculantVec()
	// TODO(samlaf): why do we convert to row if inside getFFTCoeff we convert back to col?
	rv := t.fromColVToRowV(cv)
	return t.getFFTCoeff(rv)
}

// Expand toeplitz matrix into circulant matrix
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
func (t *Toeplitz) extendCirculantVec() []fr.Element {
	E := make([]fr.Element, len(t.V)+1) // extra 1 from extended, equal to 2*dimE
	E[0].Set(&t.V[0])

	numRow := t.getMatDim()
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
func (t *Toeplitz) fromColVToRowV(cv []fr.Element) []fr.Element {
	n := len(cv)
	rv := make([]fr.Element, n)

	rv[0].Set(&cv[0])

	for i := 1; i < n; i++ {
		rv[i].Set(&cv[n-i])
	}

	return rv
}

// Taking FFT on the circulant matrix vector
func (t *Toeplitz) getFFTCoeff(rowV []fr.Element) ([]fr.Element, error) {
	n := len(rowV)
	colV := make([]fr.Element, n)
	for i := 0; i < n; i++ {
		colV[i] = rowV[(n-i)%n]
	}

	out, err := t.Fs.FFT(colV, false)
	if err != nil {
		return nil, fmt.Errorf("fft: %w", err)
	}
	return out, nil
}

func (t *Toeplitz) getMatDim() int {
	return (len(t.V) + 1) / 2
}
