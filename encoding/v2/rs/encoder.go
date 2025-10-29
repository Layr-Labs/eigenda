package rs

import (
	"errors"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/v2/fft"
	"github.com/Layr-Labs/eigenda/encoding/v2/rs/backend"
	"github.com/Layr-Labs/eigenda/encoding/v2/rs/backend/gnark"
	"github.com/Layr-Labs/eigenda/encoding/v2/rs/backend/icicle"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	_ "go.uber.org/automaxprocs"
)

type Encoder struct {
	logger logging.Logger
	Config *encoding.Config

	mu                  sync.Mutex
	ParametrizedEncoder map[encoding.EncodingParams]*ParametrizedEncoder
}

// NewEncoder creates a new encoder with the given options
func NewEncoder(logger logging.Logger, config *encoding.Config) *Encoder {
	if config == nil {
		config = encoding.DefaultConfig()
	}

	e := &Encoder{
		logger:              logger,
		Config:              config,
		mu:                  sync.Mutex{},
		ParametrizedEncoder: make(map[encoding.EncodingParams]*ParametrizedEncoder),
	}

	return e
}

// just a wrapper to take bytes not Fr Element
func (g *Encoder) EncodeBytes(inputBytes []byte, params encoding.EncodingParams) ([]FrameCoeffs, []uint32, error) {
	inputFr, err := ToFrArray(inputBytes)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot convert bytes to field elements, %w", err)
	}
	return g.Encode(inputFr, params)
}

// Encode function takes input in unit of Fr Element and creates a list of FramesCoeffs,
// which each contain a list of multireveal interpolating polynomial coefficients.
// A slice of uint32 is also returned, which corresponds to which leading coset
// root of unity the frame is proving against. This can be deduced from a frame's index.
func (g *Encoder) Encode(inputFr []fr.Element, params encoding.EncodingParams) ([]FrameCoeffs, []uint32, error) {
	start := time.Now()
	intermediate := time.Now()

	// Get RS encoder from params
	encoder, err := g.getRsEncoder(params)
	if err != nil {
		return nil, nil, err
	}

	pdCoeffs, err := encoder.padPolyEval(inputFr)
	if err != nil {
		return nil, nil, err
	}
	paddingDuration := time.Since(intermediate)

	intermediate = time.Now()

	polyEvals, err := encoder.rsEncoderBackend.ExtendPolyEval(pdCoeffs)
	if err != nil {
		return nil, nil, fmt.Errorf("reed-solomon extend poly evals, %w", err)
	}
	extensionDuration := time.Since(intermediate)

	intermediate = time.Now()

	// create Frames to group relevant info
	frames, indices, err := encoder.makeFrames(polyEvals)
	if err != nil {
		return nil, nil, err
	}

	framesDuration := time.Since(intermediate)

	// TODO(samlaf): use an injected logger instead.
	g.logger.Info("RSEncode details",
		"input_size_bytes", len(inputFr)*encoding.BYTES_PER_SYMBOL,
		"num_chunks", encoder.Params.NumChunks,
		"chunk_length", encoder.Params.ChunkLength,
		"padding_duration", paddingDuration,
		"extension_duration", extensionDuration,
		"frames_duration", framesDuration,
		"total_duration", time.Since(start))

	return frames, indices, nil
}

// Decode data when some chunks from systematic nodes are lost. This function implements
// https://ethresear.ch/t/reed-solomon-erasure-code-recovery-in-n-log-2-n-time-with-ffts/3039
//
// It first uses FFT to recover the whole polynomial. Then it extracts only the systematic chunks.
// It takes a list of available frame, and return the original encoded data
// storing the evaluation points, since it is where RS is applied. The input frame contains
// the coefficients of the interpolating polynomial, hence interpolation is needed before
// recovery.
//
// maxInputSize is the upper bound of the original data size. This is needed because
// the Frames and indices don't encode the length of the original data. If maxInputSize
// is smaller than the original input size, decoded data will be trimmed to fit the maxInputSize.
//
// TODO(samlaf): Many call sites have frames and need to convert to FrameCoeffs.
// Would be nice to figure out a Decode interface that doesn't require creating allocations.
// Perhaps Decode could take an iterator that produces one FrameCoeffs at a time?
// That way we could pass either chunks (frameCoeffs) or frames.
func (e *Encoder) Decode(
	frames []FrameCoeffs, indices []encoding.ChunkNumber, maxInputSize uint64, params encoding.EncodingParams,
) ([]byte, error) {
	// Get encoder
	g, err := e.getRsEncoder(params)
	if err != nil {
		return nil, err
	}

	if len(frames) != len(indices) {
		return nil, errors.New("number of frames must equal number of indices")
	}

	// Remove duplicates
	frameMap := make(map[encoding.ChunkNumber]FrameCoeffs, len(indices))
	for i, frameIndex := range indices {
		_, ok := frameMap[frameIndex]
		if !ok {
			frameMap[frameIndex] = frames[i]
		}
	}

	numSys := encoding.GetNumSys(maxInputSize, g.Params.ChunkLength)
	if uint64(len(frameMap)) < numSys {
		return nil, errors.New("number of frame must be sufficient")
	}

	samples := make([]*fr.Element, g.Params.NumEvaluations())
	// copy evals based on frame coeffs into samples
	for d, f := range frameMap {
		e, err := GetLeadingCosetIndex(d, g.Params.NumChunks)
		if err != nil {
			return nil, err
		}

		evals, err := g.getInterpolationPolyEval(f, e)
		if err != nil {
			return nil, err
		}

		// Some pattern i butterfly swap. Find the leading coset, then increment by number of coset
		for j := uint64(0); j < g.Params.ChunkLength; j++ {
			p := j*g.Params.NumChunks + uint64(e)
			samples[p] = new(fr.Element)
			samples[p].Set(&evals[j])
		}
	}

	reconstructedData := make([]fr.Element, g.Params.NumEvaluations())
	missingIndices := false
	for i, s := range samples {
		if s == nil {
			missingIndices = true
			break
		}
		reconstructedData[i] = *s
	}

	if missingIndices {
		var err error
		reconstructedData, err = g.Fs.RecoverPolyFromSamples(
			samples,
			g.Fs.ZeroPolyViaMultiplication,
		)
		if err != nil {
			return nil, fmt.Errorf("recover polynomial from samples: %w", err)
		}
	}

	reconstructedPoly, err := g.Fs.FFT(reconstructedData, true)
	if err != nil {
		return nil, fmt.Errorf("inverse fft on reconstructed data: %w", err)
	}

	data := ToByteArray(reconstructedPoly, maxInputSize)

	return data, nil
}

// getRsEncoder returns a parametrized encoder for the given parameters.
// It caches the encoder for reuse.
func (g *Encoder) getRsEncoder(params encoding.EncodingParams) (*ParametrizedEncoder, error) {
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
		return nil, fmt.Errorf("validate encoding params: %w", err)
	}

	fs := e.createFFTSettings(params)

	var rsEncoderBackend backend.RSEncoderBackend
	switch e.Config.BackendType {
	case encoding.GnarkBackend:
		if e.Config.GPUEnable {
			return nil, errors.New("GPU is not supported in gnark backend")
		}
		rsEncoderBackend = gnark.NewRSBackend(fs)
	case encoding.IcicleBackend:
		rsEncoderBackend, err = icicle.BuildRSBackend(e.logger, e.Config.GPUEnable, e.Config.GPUConcurrentFrameGenerationDangerous)
		if err != nil {
			return nil, fmt.Errorf("build icicle rs backend: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported backend type: %v", e.Config.BackendType)
	}
	return &ParametrizedEncoder{
		Config:           e.Config,
		Params:           params,
		Fs:               fs,
		rsEncoderBackend: rsEncoderBackend,
	}, nil
}

func (e *Encoder) createFFTSettings(params encoding.EncodingParams) *fft.FFTSettings {
	n := uint8(math.Log2(float64(params.NumEvaluations())))
	return fft.NewFFTSettings(n)
}
