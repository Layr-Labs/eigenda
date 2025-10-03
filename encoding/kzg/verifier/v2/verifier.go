package verifier

import (
	"errors"
	"fmt"
	gomath "math"
	"sync"

	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/fft"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/Layr-Labs/eigenda/encoding/rs"

	"github.com/consensys/gnark-crypto/ecc/bn254"
	_ "go.uber.org/automaxprocs"
)

type Verifier struct {
	kzgConfig *KzgConfig
	encoder   *rs.Encoder

	G1SRS kzg.G1SRS

	// mu protects access to ParametrizedVerifiers
	mu                    sync.Mutex
	ParametrizedVerifiers map[encoding.EncodingParams]*ParametrizedVerifier
}

func NewVerifier(config *KzgConfig, encoderConfig *encoding.Config) (*Verifier, error) {
	if config.SRSNumberToLoad > encoding.SRSOrder {
		return nil, errors.New("SRSOrder is less than srsNumberToLoad")
	}

	// read the whole order, and treat it as entire SRS for low degree proof
	g1SRS, err := kzg.ReadG1Points(config.G1Path, config.SRSNumberToLoad, config.NumWorker)
	if err != nil {
		return nil, fmt.Errorf("failed to read %d G1 points from %s: %w", config.SRSNumberToLoad, config.G1Path, err)
	}

	encoder, err := rs.NewEncoder(encoderConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create encoder: %w", err)
	}

	encoderGroup := &Verifier{
		kzgConfig:             config,
		encoder:               encoder,
		G1SRS:                 g1SRS,
		ParametrizedVerifiers: make(map[encoding.EncodingParams]*ParametrizedVerifier),
	}

	return encoderGroup, nil
}

func (v *Verifier) GetKzgVerifier(params encoding.EncodingParams) (*ParametrizedVerifier, error) {
	if err := encoding.ValidateEncodingParams(params, encoding.SRSOrder); err != nil {
		return nil, fmt.Errorf("validate encoding params: %w", err)
	}

	// protect access to ParametrizedVerifiers
	v.mu.Lock()
	defer v.mu.Unlock()

	ver, ok := v.ParametrizedVerifiers[params]
	if ok {
		return ver, nil
	}

	ver, err := v.newKzgVerifier(params)
	if err != nil {
		return nil, fmt.Errorf("new KZG verifier: %w", err)
	}

	v.ParametrizedVerifiers[params] = ver
	return ver, nil
}

func (v *Verifier) newKzgVerifier(params encoding.EncodingParams) (*ParametrizedVerifier, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("invalid encoding params: %w", err)
	}

	// Create FFT settings based on params
	n := uint8(gomath.Log2(float64(params.NumEvaluations())))
	fs := fft.NewFFTSettings(n)

	return &ParametrizedVerifier{
		g1SRS: v.G1SRS,
		Fs:    fs,
	}, nil
}

// VerifyFrame verifies a single frame against a commitment.
// If needing to verify multiple frames of the same chunk length, prefer [Verifier.UniversalVerify].
//
// This function is only used in the v1 and v2 validator (distributed) retrievers.
// TODO(samlaf): replace these with UniversalVerify, and consider deleting this function.
func (v *Verifier) VerifyFrames(
	frames []*encoding.Frame,
	indices []encoding.ChunkNumber,
	commitments encoding.BlobCommitments,
	params encoding.EncodingParams) error {

	if len(frames) != len(indices) {
		return fmt.Errorf("invalid number of frames and indices: %d != %d", len(frames), len(indices))
	}

	verifier, err := v.GetKzgVerifier(params)
	if err != nil {
		return err
	}

	for ind := range frames {
		err = verifier.VerifyFrame(
			frames[ind],
			uint64(indices[ind]),
			(*bn254.G1Affine)(commitments.Commitment),
			params.NumChunks,
		)

		if err != nil {
			return err
		}
	}

	return nil
}

// Decode takes in the chunks, indices, and encoding parameters and returns the decoded blob
// The result is trimmed to the given maxInputSize.
func (v *Verifier) Decode(
	chunks []*encoding.Frame, indices []encoding.ChunkNumber, params encoding.EncodingParams, maxInputSize uint64,
) ([]byte, error) {
	frames := make([]rs.FrameCoeffs, len(chunks))
	for i := range chunks {
		frames[i] = chunks[i].Coeffs
	}

	return v.encoder.Decode(frames, toUint64Array(indices), maxInputSize, params) //nolint:wrapcheck
}

func toUint64Array(chunkIndices []encoding.ChunkNumber) []uint64 {
	res := make([]uint64, len(chunkIndices))
	for i, d := range chunkIndices {
		res[i] = uint64(d)
	}
	return res
}
