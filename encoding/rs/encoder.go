package rs

import (
	"math"

	kzg "github.com/Layr-Labs/eigenda/pkg/kzg"
)

type Encoder struct {
	EncodingParams

	Fs *kzg.FFTSettings

	verbose bool
}

// The function creates a high level struct that determines the encoding the a data of a
// specific length under (num systematic node, num parity node) setup. A systematic node
// stores a systematic data chunk that contains part of the original data. A parity node
// stores a parity data chunk which is an encoding of the original data. A receiver that
// collects all systematic chunks can simply stitch data together to reconstruct the
// original data. When some systematic chunks are missing but identical parity chunk are
// available, the receive can go through a Reed Solomon decoding to reconstruct the
// original data.
func NewEncoder(params EncodingParams, verbose bool) (*Encoder, error) {

	err := params.Validate()
	if err != nil {
		return nil, err
	}

	n := uint8(math.Log2(float64(params.NumEvaluations())))
	fs := kzg.NewFFTSettings(n)

	return &Encoder{
		EncodingParams: params,
		Fs:             fs,
		verbose:        verbose,
	}, nil

}
