package toeplitz

import (
	"errors"

	kzg "github.com/Layr-Labs/eigenda/pkg/kzg"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

type Circular struct {
	V  []fr.Element
	Fs *kzg.FFTSettings
}

func NewCircular(v []fr.Element, fs *kzg.FFTSettings) *Circular {
	return &Circular{
		V:  v,
		Fs: fs,
	}
}

// Matrix multiplication between a circular matrix and a vector using FFT
func (c *Circular) Multiply(x []fr.Element) ([]fr.Element, error) {
	if len(x) != len(c.V) {
		return nil, errors.New("dimension inconsistent")
	}
	n := len(x)

	colV := make([]fr.Element, n)
	for i := 0; i < n; i++ {
		colV[i] = c.V[(n-i)%n]
	}

	y, err := c.Fs.FFT(x, false)
	if err != nil {
		return nil, err
	}
	v, err := c.Fs.FFT(colV, false)
	if err != nil {
		return nil, err
	}
	u := make([]fr.Element, n)
	err = Hadamard(y, v, u)
	if err != nil {
		return nil, err
	}

	r, err := c.Fs.FFT(u, true)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// Taking FFT on the circular matrix vector
func (c *Circular) GetFFTCoeff() ([]fr.Element, error) {
	n := len(c.V)

	colV := make([]fr.Element, n)
	for i := 0; i < n; i++ {
		colV[i] = c.V[(n-i)%n]
	}

	out, err := c.Fs.FFT(colV, false)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Taking FFT on the circular matrix vector
func (c *Circular) GetCoeff() ([]fr.Element, error) {
	n := len(c.V)

	colV := make([]fr.Element, n)
	for i := 0; i < n; i++ {
		colV[i] = c.V[(n-i)%n]
	}
	return colV, nil
}

// Hadamard product between 2 vectors containing Fr elements
func Hadamard(a, b, u []fr.Element) error {
	if len(a) != len(b) {
		return errors.New("dimension inconsistent. Cannot do Hadamard Product on Fr")
	}

	for i := 0; i < len(a); i++ {
		u[i].Mul(&a[i], &b[i])
	}
	return nil
}
