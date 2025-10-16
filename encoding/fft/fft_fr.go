package fft

import (
	"fmt"

	"github.com/Layr-Labs/eigenda/common/math"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr/fft"
)

// InputNotPowerOfTwoError is an error that indicates that the input to the FFT is not a power of two.
type InputNotPowerOfTwoError struct {
	inputLen uint64
}

func (e *InputNotPowerOfTwoError) Error() string {
	return fmt.Sprintf("(I)FFT input length %d is not a power of two", e.inputLen)
}

// Is checks if the error is an InputNotPowerOfTwoError.
// It is implemented to allow errors.Is to work with this error type,
// so that we can use the sentinel as errors.Is(err, ErrNotPowerOfTwo) to check for this error type.
func (e *InputNotPowerOfTwoError) Is(target error) bool {
	if _, ok := target.(*InputNotPowerOfTwoError); ok {
		return true
	}
	return false
}

// NewFFTInputNotPowerOfTwoError creates a new FFTInputNotPowerOfTwoError with the given input length.
func NewFFTInputNotPowerOfTwoError(inputLen uint64) *InputNotPowerOfTwoError {
	return &InputNotPowerOfTwoError{
		inputLen: inputLen,
	}
}

var (
	// ErrNotPowerOfTwo is a sentinel error that can be used to check if an error is an [FFTInputNotPowerOfTwoError].
	// by calling errors.Is(err, ErrNotPowerOfTwo)
	ErrNotPowerOfTwo = &InputNotPowerOfTwoError{inputLen: 0}
)

// FFT performs a fast Fourier transform on the provided values, using the roots of unity
// provided in the FFTSettings.
//
// The input values does not have to be a power of two, because we pad them to the next power of two.
// It's power of two must be equal to the max width of the FFTSettings.
//
// It outputs a newly allocated slice of field elements, which is the transformed values.
// To perform the FFT in-place, use [FFTSettings.InplaceFFT] instead.
//
// The only error returned is if the FFTSettings does not have enough roots of unity
// to perform the FFT on the input values.
func (fs *FFTSettings) FFT(vals []fr.Element, inv bool) ([]fr.Element, error) {
	n := uint64(len(vals))
	if n != fs.MaxWidth {
		return nil, fmt.Errorf("FFT input length %d is not equal to max width %d", n, fs.MaxWidth)
	}
	n = math.NextPowOf2u64(n)
	// We make a copy so we can mutate it during the work.
	valsCopy := make([]fr.Element, n)
	for i := 0; i < len(vals); i++ {
		valsCopy[i].Set(&vals[i])
	}
	for i := uint64(len(vals)); i < n; i++ {
		// Otherwise like this we change the commitment wrt the original polynomial.
		valsCopy[i].SetZero()
	}
	if inv {
		fs.Domain.FFTInverse(valsCopy, fft.DIF)
	} else {
		fs.Domain.FFT(valsCopy, fft.DIF)
	}
	fft.BitReverse(valsCopy)
	return valsCopy, nil
}

func (fs *FFTSettings) InplaceFFT(vals []fr.Element, out []fr.Element, inv bool) error {
	n := uint64(len(vals))
	if n != fs.MaxWidth {
		return fmt.Errorf("FFT input length %d is not equal to max width %d", n, fs.MaxWidth)
	}
	if !math.IsPowerOfTwo(n) {
		return NewFFTInputNotPowerOfTwoError(n)
	}
	for i, val := range vals {
		out[i].Set(&val)
	}
	if inv {
		fs.Domain.FFTInverse(out, fft.DIF)
	} else {
		fs.Domain.FFT(out, fft.DIF)
	}
	fft.BitReverse(out)
	return nil
}
