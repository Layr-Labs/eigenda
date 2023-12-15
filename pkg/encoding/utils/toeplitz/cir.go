package toeplitz

import (
	"errors"

	kzg "github.com/Layr-Labs/eigenda/pkg/kzg"
	bls "github.com/Layr-Labs/eigenda/pkg/kzg/bn254"
)

type Circular struct {
	V  []bls.Fr
	Fs *kzg.FFTSettings
}

func NewCircular(v []bls.Fr, fs *kzg.FFTSettings) *Circular {
	return &Circular{
		V:  v,
		Fs: fs,
	}
}

// Matrix multiplication between a circular matrix and a vector using FFT
func (c *Circular) Multiply(x []bls.Fr) ([]bls.Fr, error) {
	if len(x) != len(c.V) {
		return nil, errors.New("dimension inconsistent")
	}
	n := len(x)

	colV := make([]bls.Fr, n)
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
	u := make([]bls.Fr, n)
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

// Matrix Multiplication between a circular matrix of Fr element,
// and a vector of G1 points being supplied in the argument. This method uses FFT to
// compute matrix multiplication. On the optimization side, the
// function allows for using precomputed FFT as its input.
func (c *Circular) MultiplyPoints(x []bls.G1Point, inv bool, usePrecompute bool) ([]bls.G1Point, error) {
	if len(x) != len(c.V) {
		return nil, errors.New("dimension inconsistent. Input != vector")
	}
	n := len(x)

	colV := make([]bls.Fr, n)
	for i := 0; i < n; i++ {
		colV[i] = c.V[(n-i)%n]
	}

	y := x

	v, err := c.Fs.FFT(colV, false)
	if err != nil {
		return nil, err
	}
	u := make([]bls.G1Point, n)
	err = HadamardPoints(y, v, u)
	if err != nil {
		return nil, err
	}

	if inv {
		r, err := c.Fs.FFTG1(u, true)
		if err != nil {
			return nil, err
		}
		return r, nil
	} else {
		return u, nil
	}
}

// Taking FFT on the circular matrix vector
func (c *Circular) GetFFTCoeff() ([]bls.Fr, error) {
	n := len(c.V)

	colV := make([]bls.Fr, n)
	for i := 0; i < n; i++ {
		colV[i] = c.V[(n-i)%n]
	}

	out, err := c.Fs.FFT(colV, false)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Hadamard product between 2 vectors, one contains G1 points, the other contains Fr element
func HadamardPoints(a []bls.G1Point, b []bls.Fr, u []bls.G1Point) error {
	if len(a) != len(b) {
		return errors.New("dimension inconsistent. Cannot do Hadamard Product on Points")
	}

	for i := 0; i < len(a); i++ {
		bls.MulG1(&u[i], &a[i], &b[i])
	}
	return nil
}

// Hadamard product between 2 vectors containing Fr elements
func Hadamard(a, b, u []bls.Fr) error {
	if len(a) != len(b) {
		return errors.New("dimension inconsistent. Cannot do Hadamard Product on Fr")
	}

	for i := 0; i < len(a); i++ {
		bls.MulModFr(&u[i], &a[i], &b[i])
	}
	return nil
}

// Naive implementation of a Multiplication between a matrix and vector.
// both contains Fr elements
func (c *Circular) DirectMultiply(x []bls.Fr) []bls.Fr {
	n := len(x)

	out := make([]bls.Fr, n)
	for i := 0; i < n; i++ {
		var sum bls.Fr
		bls.CopyFr(&sum, &bls.ZERO)
		for j := 0; j < n; j++ {
			idx := (j - i + n) % n
			var product bls.Fr
			bls.MulModFr(&product, &c.V[idx], &x[j])
			bls.AddModFr(&sum, &product, &sum)
		}
		bls.CopyFr(&out[i], &sum)
	}
	return out
}
