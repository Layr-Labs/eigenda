package verifier

import (
	"errors"
	"fmt"
	"math"
	"sync"

	"github.com/Layr-Labs/eigenda/encoding"

	"github.com/Layr-Labs/eigenda/encoding/fft"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/Layr-Labs/eigenda/encoding/rs"

	"github.com/consensys/gnark-crypto/ecc/bn254"
	_ "go.uber.org/automaxprocs"
)

type Verifier struct {
	kzgConfig *kzg.KzgConfig
	encoder   *rs.Encoder

	G1SRS kzg.G1SRS

	// mu protects access to ParametrizedVerifiers
	mu                    sync.Mutex
	ParametrizedVerifiers map[encoding.EncodingParams]*ParametrizedVerifier
}

var _ encoding.Verifier = &Verifier{}

func NewVerifier(config *kzg.KzgConfig, encoderConfig *encoding.Config) (*Verifier, error) {
	if config.SRSNumberToLoad > config.SRSOrder {
		return nil, errors.New("SRSOrder is less than srsNumberToLoad")
	}

	// read the whole order, and treat it as entire SRS for low degree proof
	g1SRS, err := kzg.ReadG1Points(config.G1Path, config.SRSNumberToLoad, config.NumWorker)
	if err != nil {
		return nil, fmt.Errorf("failed to read %d G1 points from %s: %v", config.SRSNumberToLoad, config.G1Path, err)
	}

	encoder, err := rs.NewEncoder(encoderConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create encoder: %v", err)
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
	if err := encoding.ValidateEncodingParams(params, v.kzgConfig.SRSOrder); err != nil {
		return nil, err
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
	n := uint8(math.Log2(float64(params.NumEvaluations())))
	fs := fft.NewFFTSettings(n)

	return &ParametrizedVerifier{
		KzgConfig: v.kzgConfig,
		g1SRS:     v.G1SRS,
		Fs:        fs,
	}, nil
}

func (v *Verifier) VerifyBlobLength(commitments encoding.BlobCommitments) error {
	return v.VerifyCommit(
		(*bn254.G2Affine)(commitments.LengthCommitment),
		(*bn254.G2Affine)(commitments.LengthProof),
		uint64(commitments.Length),
	)
}

// VerifyCommit verifies the low degree proof; since it doesn't depend on the encoding parameters
// we leave it as a method of the KzgEncoderGroup
func (v *Verifier) VerifyCommit(lengthCommit *bn254.G2Affine, lengthProof *bn254.G2Affine, length uint64) error {

	g1Challenge, err := kzg.ReadG1Point(v.kzgConfig.SRSOrder-length, v.kzgConfig.SRSOrder, v.kzgConfig.G1Path)
	if err != nil {
		return fmt.Errorf("read g1 point: %w", err)
	}

	err = VerifyLengthProof(lengthCommit, lengthProof, &g1Challenge)
	if err != nil {
		return fmt.Errorf("low degree proof: %w", err)
	}
	return nil
}

// The function verify low degree proof against a poly commitment
// We wish to show x^shift poly = shiftedPoly, with
// With shift = SRSOrder - length and
// proof = commit(shiftedPoly) on G1
// so we can verify by checking
// e( commit_1, [x^shift]_2) = e( proof_1, G_2 )
func VerifyLengthProof(lengthCommit *bn254.G2Affine, proof *bn254.G2Affine, g1Challenge *bn254.G1Affine) error {
	return PairingsVerify(g1Challenge, lengthCommit, &kzg.GenG1, proof)
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
func (v *Verifier) Decode(chunks []*encoding.Frame, indices []encoding.ChunkNumber, params encoding.EncodingParams, maxInputSize uint64) ([]byte, error) {
	frames := make([]rs.FrameCoeffs, len(chunks))
	for i := range chunks {
		frames[i] = chunks[i].Coeffs
	}

	return v.encoder.Decode(frames, toUint64Array(indices), maxInputSize, params)
}

func toUint64Array(chunkIndices []encoding.ChunkNumber) []uint64 {
	res := make([]uint64, len(chunkIndices))
	for i, d := range chunkIndices {
		res[i] = uint64(d)
	}
	return res
}

func PairingsVerify(a1 *bn254.G1Affine, a2 *bn254.G2Affine, b1 *bn254.G1Affine, b2 *bn254.G2Affine) error {
	var negB1 bn254.G1Affine
	negB1.Neg(b1)

	P := [2]bn254.G1Affine{*a1, negB1}
	Q := [2]bn254.G2Affine{*a2, *b2}

	ok, err := bn254.PairingCheck(P[:], Q[:])
	if err != nil {
		return fmt.Errorf("PairingCheck: %w", err)
	}
	if !ok {
		return errors.New("PairingCheck pairing not ok. SRS is invalid")
	}

	return nil
}
