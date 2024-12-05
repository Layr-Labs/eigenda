package rs

import (
	"errors"
	"fmt"
	"math"
	"sync"

	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/fft"
	gnarkencoder "github.com/Layr-Labs/eigenda/encoding/rs/gnark"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	_ "go.uber.org/automaxprocs"
)

type Encoder struct {
	Config *encoding.Config

	mu                  sync.Mutex
	ParametrizedEncoder map[encoding.EncodingParams]*ParametrizedEncoder
}

// Proof device represents a device capable of computing reed-solomon operations.
type EncoderDevice interface {
	ExtendPolyEval(coeffs []fr.Element) ([]fr.Element, error)
}

// NewEncoder creates a new encoder with the given options
func NewEncoder(config *encoding.Config) (*Encoder, error) {
	if config == nil {
		config = encoding.DefaultConfig()
	}

	e := &Encoder{
		Config:              config,
		mu:                  sync.Mutex{},
		ParametrizedEncoder: make(map[encoding.EncodingParams]*ParametrizedEncoder),
	}

	return e, nil
}

// GetRsEncoder returns a parametrized encoder for the given parameters.
// It caches the encoder for reuse.
func (g *Encoder) GetRsEncoder(params encoding.EncodingParams) (*ParametrizedEncoder, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	enc, ok := g.ParametrizedEncoder[params]
	if ok {
		return enc, nil
	}

	enc, err := g.newEncoder(params)
	if err == nil {
		g.ParametrizedEncoder[params] = enc
	}

	return enc, err
}

// The function creates a high level struct that determines the encoding the a data of a
// specific length under (num systematic node, num parity node) setup. A systematic node
// stores a systematic data chunk that contains part of the original data. A parity node
// stores a parity data chunk which is an encoding of the original data. A receiver that
// collects all systematic chunks can simply stitch data together to reconstruct the
// original data. When some systematic chunks are missing but identical parity chunk are
// available, the receive can go through a Reed Solomon decoding to reconstruct the
// original data.
func (e *Encoder) newEncoder(params encoding.EncodingParams) (*ParametrizedEncoder, error) {
	err := params.Validate()
	if err != nil {
		return nil, err
	}

	fs := e.CreateFFTSettings(params)

	switch e.Config.BackendType {
	case encoding.GnarkBackend:
		return e.createGnarkBackendEncoder(params, fs)
	case encoding.IcicleBackend:
		return e.createIcicleBackendEncoder(params, fs)
	default:
		return nil, fmt.Errorf("unsupported backend type: %v", e.Config.BackendType)
	}
}

func (e *Encoder) CreateFFTSettings(params encoding.EncodingParams) *fft.FFTSettings {
	n := uint8(math.Log2(float64(params.NumEvaluations())))
	return fft.NewFFTSettings(n)
}

func (e *Encoder) createGnarkBackendEncoder(params encoding.EncodingParams, fs *fft.FFTSettings) (*ParametrizedEncoder, error) {
	if e.Config.GPUEnable {
		return nil, errors.New("GPU is not supported in gnark backend")
	}

	return &ParametrizedEncoder{
		Config:            e.Config,
		EncodingParams:    params,
		Fs:                fs,
		RSEncoderComputer: &gnarkencoder.RsGnarkBackend{Fs: fs},
	}, nil
}

func (e *Encoder) createIcicleBackendEncoder(params encoding.EncodingParams, fs *fft.FFTSettings) (*ParametrizedEncoder, error) {
	return CreateIcicleBackendEncoder(e, params, fs)
}
