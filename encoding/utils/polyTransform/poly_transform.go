package polyTranform

import (
	"fmt"
	"math"

	"github.com/Layr-Labs/eigenda/encoding/fft"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

type PolyTranform struct {
	transformer *fft.FFTSettings
}

// NewPolyTranform takes an input uint8(l). It is a suggestive parameter that 2**l is the
// largest possible field elements to be used on the object. If larger data (number of field element)
// is used with methods of this object, the object reinitializes its capacity.
// If data size is unknown when creating this object, use a small number like 4
func NewPolyTranform(l uint8) (*PolyTranform, error) {
	// the max number data supported by bn254 curve is 2**28
	if l >= 28 {
		return nil, fmt.Errorf("requested size is greater than capacity of the bn254 curve")
	}

	transformer := fft.NewFFTSettings(l)

	return &PolyTranform{
		transformer: transformer,
	}, nil
}

func (p *PolyTranform) ConvertEvalsToCoeffs(inputFr []fr.Element) ([]fr.Element, error) {
	n := len(inputFr)

	// the max number data supported by bn254 curve is 2**28
	if n > 268435456 {
		return nil, fmt.Errorf("requested size is greater than capacity of the bn254 curve")
	}

	if n > int(p.transformer.MaxWidth) {
		l := uint8(math.Ceil(math.Log2(float64(n))))
		// reinitialize fft size
		p.transformer = fft.NewFFTSettings(l)
	}
	return p.transformer.ConvertEvalsToCoeffs(inputFr)
}

func (p *PolyTranform) ConvertCoeffsToEvals(inputFr []fr.Element) ([]fr.Element, error) {
	n := len(inputFr)

	// the max number data supported by bn254 curve is 2**28
	if n > 268435456 {
		return nil, fmt.Errorf("requested data size is greater than capacity of the bn254 curve")
	}

	if n > int(p.transformer.MaxWidth) {
		l := uint8(math.Ceil(math.Log2(float64(n))))
		// reinitialize fft size
		p.transformer = fft.NewFFTSettings(l)
	}
	return p.transformer.ConvertCoeffsToEvals(inputFr)
}
