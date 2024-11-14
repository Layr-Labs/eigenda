package rs

import (
	"fmt"
	"math"
	"runtime"
	"sync"

	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/fft"
	"github.com/Layr-Labs/eigenda/encoding/rs/cpu"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"

	_ "go.uber.org/automaxprocs"
)

// EncoderOption defines the function signature for encoder options
type EncoderOption func(*Encoder)

type Encoder struct {
	Config *encoding.Config

	mu                  sync.Mutex
	ParametrizedEncoder map[encoding.EncodingParams]*ParametrizedEncoder
}

// Proof device represents a device capable of computing reed-solomon operations.
type EncoderDevice interface {
	ExtendPolyEval(coeffs []fr.Element) ([]fr.Element, error)
}

// Default configuration values
const (
	defaultBackend   = encoding.BackendDefault
	defaultEnableGPU = false
	defaultVerbose   = false
	defaultNTTSize   = 25 // Used for NTT setup in Icicle backend
)

// Option Definitions
func WithBackend(backend encoding.BackendType) EncoderOption {
	return func(e *Encoder) {
		e.Config.BackendType = backend
	}
}

func WithGPU(enable bool) EncoderOption {
	return func(e *Encoder) {
		e.Config.EnableGPU = enable
	}
}

func WithNumWorkers(workers uint64) EncoderOption {
	return func(e *Encoder) {
		e.Config.NumWorker = workers
	}
}

func WithVerbose(verbose bool) EncoderOption {
	return func(e *Encoder) {
		e.Config.Verbose = verbose
	}
}

// NewEncoder creates a new encoder with the given options
func NewEncoder(opts ...EncoderOption) (*Encoder, error) {
	e := &Encoder{
		Config: &encoding.Config{
			NumWorker:   uint64(runtime.GOMAXPROCS(0)),
			BackendType: defaultBackend,
			EnableGPU:   defaultEnableGPU,
			Verbose:     defaultVerbose,
		},

		mu:                  sync.Mutex{},
		ParametrizedEncoder: make(map[encoding.EncodingParams]*ParametrizedEncoder),
	}

	for _, opt := range opts {
		opt(e)
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
	case encoding.BackendDefault:
		return e.createDefaultBackendEncoder(params, fs)
	case encoding.BackendIcicle:
		return e.createIcicleBackendEncoder(params, fs)
	default:
		return nil, fmt.Errorf("unsupported backend type: %v", e.Config.BackendType)
	}
}

func (e *Encoder) CreateFFTSettings(params encoding.EncodingParams) *fft.FFTSettings {
	n := uint8(math.Log2(float64(params.NumEvaluations())))
	if params.ChunkLength == 1 {
		n = uint8(math.Log2(float64(2 * params.NumChunks)))
	}
	return fft.NewFFTSettings(n)
}

func (e *Encoder) createDefaultBackendEncoder(params encoding.EncodingParams, fs *fft.FFTSettings) (*ParametrizedEncoder, error) {
	if e.Config.EnableGPU {
		return nil, fmt.Errorf("GPU is not supported in default backend")
	}

	return &ParametrizedEncoder{
		Config:            e.Config,
		EncodingParams:    params,
		Fs:                fs,
		RSEncoderComputer: &cpu.RsDefaultComputeDevice{Fs: fs},
	}, nil
}

func (e *Encoder) createIcicleBackendEncoder(params encoding.EncodingParams, fs *fft.FFTSettings) (*ParametrizedEncoder, error) {
	fmt.Println("CreateIcicleBackendEncoder", e.Config.EnableGPU)
	return CreateIcicleBackendEncoder(e, params, fs)
}
