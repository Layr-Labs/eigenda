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

// NewPolyTranform takes an input such that l**2 represents the
// max number of data required to represent a field element.
// The max number can be computed as num_byte / 32.
// This number is suggestive, reinitiallization happens when
// encountering data with more bytes,
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
