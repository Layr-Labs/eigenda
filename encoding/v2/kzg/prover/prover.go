package prover

import (
	"errors"
	"fmt"
	gomath "math"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Layr-Labs/eigenda/common/math"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/v2/fft"
	"github.com/Layr-Labs/eigenda/encoding/v2/kzg"
	"github.com/Layr-Labs/eigenda/encoding/v2/kzg/prover/backend"
	"github.com/Layr-Labs/eigenda/encoding/v2/kzg/prover/backend/gnark"
	"github.com/Layr-Labs/eigenda/encoding/v2/kzg/prover/backend/icicle"
	"github.com/Layr-Labs/eigenda/encoding/v2/rs"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	_ "go.uber.org/automaxprocs"
)

// ProvingParams controls the size of matrix multiplication when generating kzg multi-reveal proofs.
// For a blob that is zero appended to BlobLength (equal to power of 2) field elements, two parameters holds the
// relation ChunkLength * ToeplitzMatrixLength = BlobLength, where ChunkLength equals to the same parameters from
// the encoding.EncodingParams. They maps to the Kate Amortized paper, https://eprint.iacr.org/2023/033.pdf,
// proposition 4, where ChunkLength is l, and ToeplitzMatrixLength is r. In the paper, the length of the square
// toeplitz matrix is r-1, but in order to use standard FFT library, we pad the matrix in both dimension with 0;
// and we pad the vector being multiplied with 0. The multiplication result still holds.
type ProvingParams struct {
	ChunkLength uint64
	BlobLength  uint64
}

func (p *ProvingParams) ToeplitzSquareMatrixLength() uint64 {
	return p.BlobLength / p.ChunkLength
}

// blobLength assumes to be power of 2
func BuildProvingParamsFromEncodingParams(params encoding.EncodingParams, blobLength uint64) (ProvingParams, error) {
	if blobLength < params.ChunkLength {
		return ProvingParams{}, fmt.Errorf("blob length should at least equal to the chunk length")
	}

	return ProvingParams{
		ChunkLength: params.ChunkLength,
		BlobLength:  blobLength,
	}, nil
}

func ValidateProvingParams(params ProvingParams, srsOrder uint64) error {
	toeplitzLength := params.ToeplitzSquareMatrixLength()

	if toeplitzLength == 0 {
		return errors.New("size of square toeplitz length must be greater than 0")
	}
	if params.ChunkLength == 0 {
		return errors.New("chunk length must be greater than 0")
	}

	if toeplitzLength > gomath.MaxUint64/params.ChunkLength {
		return fmt.Errorf("multiplication overflow: ChunkLength: %d, NumChunks: %d",
			params.ChunkLength, toeplitzLength)
	}

	if !math.IsPowerOfTwo(params.ChunkLength) || !math.IsPowerOfTwo(toeplitzLength) {
		return fmt.Errorf("proving parameters must be power of 2: ChunkLength: %d, ToeplitzMatrixLength: %d",
			params.ChunkLength, toeplitzLength)
	}

	if params.BlobLength > srsOrder {
		return fmt.Errorf("the supplied encoding parameters are not valid with respect to the SRS. "+
			"BlobLength %d, ChunkLength %d, NumChunks %d, SRSOrder %d",
			params.BlobLength,
			params.ChunkLength,
			toeplitzLength,
			srsOrder,
		)
	}

	return nil
}

// Prover is the main struct that is able to generate frames (chunks and their proofs).
// TODO(samlaf): should we refactor prover to only generate proofs and keep encoding separate?
type Prover struct {
	logger logging.Logger

	KzgConfig *KzgConfig
	G1SRS     kzg.G1SRS

	encoder *rs.Encoder
	Config  *encoding.Config

	// mu protects access to ParametrizedProvers
	mu                  sync.Mutex
	ParametrizedProvers map[ProvingParams]*ParametrizedProver
	SRSTables           map[ProvingParams][][]bn254.G1Affine
}

func NewProver(logger logging.Logger, kzgConfig *KzgConfig, encoderConfig *encoding.Config) (*Prover, error) {
	if encoderConfig == nil {
		encoderConfig = encoding.DefaultConfig()
	}

	if kzgConfig.SRSNumberToLoad > encoding.SRSOrder {
		return nil, errors.New("SRSOrder is less than srsNumberToLoad")
	}

	// read the whole order, and treat it as entire SRS for low degree proof
	g1SRS, err := kzg.ReadG1Points(kzgConfig.G1Path, kzgConfig.SRSNumberToLoad, kzgConfig.NumWorker)
	if err != nil {
		return nil, fmt.Errorf("failed to read G1 points: %w", err)
	}

	rsEncoder := rs.NewEncoder(logger, encoderConfig)

	proverGroup := &Prover{
		logger:              logger,
		Config:              encoderConfig,
		encoder:             rsEncoder,
		KzgConfig:           kzgConfig,
		G1SRS:               g1SRS,
		ParametrizedProvers: make(map[ProvingParams]*ParametrizedProver),
		SRSTables:           make(map[ProvingParams][][]bn254.G1Affine),
	}

	if kzgConfig.PreloadEncoder {
		// create table dir if not exist
		err := os.MkdirAll(kzgConfig.CacheDir, os.ModePerm)
		if err != nil {
			return nil, fmt.Errorf("make cache dir: %w", err)
		}

		err = proverGroup.preloadSRSTableCache()
		if err != nil {
			return nil, fmt.Errorf("preload all provers: %w", err)
		}
	}

	return proverGroup, nil
}

func (e *Prover) GetFrames(inputFr []fr.Element, params encoding.EncodingParams) ([]*encoding.Frame, []uint32, error) {

	blobLength := uint64(math.NextPowOf2u32(uint32(len(inputFr))))
	provingParams, err := BuildProvingParamsFromEncodingParams(params, blobLength)
	if err != nil {
		return nil, nil, fmt.Errorf("get proving params: %w", err)
	}

	prover, err := e.GetKzgProver(params, provingParams)
	if err != nil {
		return nil, nil, fmt.Errorf("get kzg prover: %w", err)
	}

	type encodeChanResult struct {
		chunks   []rs.FrameCoeffs
		indices  []uint32
		duration time.Duration
		err      error
	}
	encodeChan := make(chan encodeChanResult, 1)
	go func() {
		defer close(encodeChan)
		encodeStart := time.Now()
		frames, indices, err := e.encoder.Encode(inputFr, params)
		encodingDuration := time.Since(encodeStart)
		encodeChan <- encodeChanResult{
			chunks:   frames,
			indices:  indices,
			duration: encodingDuration,
			err:      err,
		}
	}()

	getProofsStart := time.Now()
	proofs, err := prover.GetProofs(inputFr)
	getProofsDuration := time.Since(getProofsStart)

	// Wait for both chunks and frames to have finished generating
	encodeResult := <-encodeChan
	if err != nil || encodeResult.err != nil {
		return nil, nil, fmt.Errorf("get frames: %w", errors.Join(err, encodeResult.err))
	}
	if len(encodeResult.chunks) != len(proofs) {
		return nil, nil, fmt.Errorf("number of chunks %v and proofs %v do not match",
			len(encodeResult.chunks), len(proofs))
	}

	e.logger.Info("Frame process details",
		"input_size_bytes", len(inputFr)*encoding.BYTES_PER_SYMBOL,
		"num_chunks", params.NumChunks,
		"chunk_length", params.ChunkLength,
		"rs_encode_duration", encodeResult.duration,
		"multi_proof_duration", getProofsDuration,
	)

	frames := make([]*encoding.Frame, len(proofs))
	for i, index := range encodeResult.indices {
		frames[i] = &encoding.Frame{
			Coeffs: encodeResult.chunks[i],
			// Coeffs are returned according to indices order, but proofs are not
			// TODO(samlaf): we should be consistent about this.
			Proof: proofs[index],
		}
	}
	return frames, encodeResult.indices, nil
}

func (g *Prover) GetKzgProver(
	params encoding.EncodingParams,
	provingParams ProvingParams,
) (*ParametrizedProver, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	enc, ok := g.ParametrizedProvers[provingParams]
	if ok {
		return enc, nil
	}

	enc, err := g.newProver(params, provingParams)
	if err != nil {
		return nil, fmt.Errorf("new prover: %w", err)
	}

	g.ParametrizedProvers[provingParams] = enc
	return enc, nil
}

func (p *Prover) newProver(params encoding.EncodingParams, provingParams ProvingParams) (*ParametrizedProver, error) {
	if err := encoding.ValidateEncodingParams(params, encoding.SRSOrder); err != nil {
		return nil, fmt.Errorf("validate encoding params: %w", err)
	}

	if err := ValidateProvingParams(provingParams, encoding.SRSOrder); err != nil {
		return nil, fmt.Errorf("validate proving params: %w", err)
	}

	// Create FFT settings based on params
	n := uint8(gomath.Log2(float64(params.NumEvaluations())))
	if params.ChunkLength == 1 {
		n = uint8(gomath.Log2(float64(2 * params.NumChunks)))
	}
	fs := fft.NewFFTSettings(n)

	// if SRS already preloaded, don't try to load or generate new ones
	fftPointsT, ok := p.SRSTables[provingParams]
	if !ok {
		var err error
		_, fftPointsT, err = p.setupFFTPoints(provingParams)
		if err != nil {
			return nil, fmt.Errorf("setup fft points: %w", err)
		}
	}

	var multiproofsBackend backend.KzgMultiProofsBackendV2
	switch p.Config.BackendType {
	case encoding.GnarkBackend:
		if p.Config.GPUEnable {
			return nil, errors.New("GPU is not supported in gnark backend")
		}
		multiproofsBackend = gnark.NewMultiProofBackend(p.logger, fs, fftPointsT)
	case encoding.IcicleBackend:
		var err error
		multiproofsBackend, err = icicle.NewMultiProofBackend(
			p.logger, fs, fftPointsT, p.G1SRS, p.Config.GPUEnable, p.KzgConfig.NumWorker)
		if err != nil {
			return nil, fmt.Errorf("create icicle backend prover: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported backend type: %v", p.Config.BackendType)
	}

	return &ParametrizedProver{
		srsNumberToLoad:            p.KzgConfig.SRSNumberToLoad,
		encodingParams:             params,
		computeMultiproofNumWorker: p.KzgConfig.NumWorker,
		kzgMultiProofBackend:       multiproofsBackend,
	}, nil

}

// preload existing SRS tables from the file directory
func (g *Prover) preloadSRSTableCache() error {
	provingParamsAll, err := getAllPrecomputedSrsMap(g.KzgConfig.CacheDir)
	if err != nil {
		return err
	}
	g.logger.Info("Detected SRSTables from cache dir", "NumTables",
		len(provingParamsAll), "TableDetails", provingParamsAll)

	if len(provingParamsAll) == 0 {
		return nil
	}

	// since
	for _, provingParams := range provingParamsAll {
		_, fftPointsT, err := g.setupFFTPoints(provingParams)
		if err != nil {
			return err
		}

		g.SRSTables[provingParams] = fftPointsT
	}

	return nil
}

// Detect the precomputed table from the specified directory
// the file name follow the name convention of
//
//	dimE*.coset&
//
// where the first * specifies the dimension of the matrix which
// equals to the number of chunks
// where the second & specifies the length of each chunk
func getAllPrecomputedSrsMap(tableDir string) ([]ProvingParams, error) {
	files, err := os.ReadDir(tableDir)
	if err != nil {
		return nil, fmt.Errorf("read srs table dir: %w", err)
	}

	tables := make([]ProvingParams, 0)
	for _, file := range files {
		filename := file.Name()

		tokens := strings.Split(filename, ".")
		dimEValue, err := strconv.Atoi(tokens[0][4:])
		if err != nil {
			return nil, fmt.Errorf("parse dimension part of the table: %w", err)
		}
		cosetSizeValue, err := strconv.Atoi(tokens[1][5:])
		if err != nil {
			return nil, fmt.Errorf("parse coset size part of the table: %w", err)
		}

		blobLength := dimEValue * cosetSizeValue

		params := ProvingParams{
			BlobLength:  uint64(blobLength),
			ChunkLength: uint64(cosetSizeValue),
		}
		tables = append(tables, params)
	}
	return tables, nil
}

// Returns SRSTable SRS points, as well as its transpose.
// fftPoints has size [l][2*dimE], and its transpose has size [2*dimE][l]
func (p *Prover) setupFFTPoints(provingParams ProvingParams) ([][]bn254.G1Affine, [][]bn254.G1Affine, error) {
	subTable, err := NewSRSTable(p.logger, p.KzgConfig.CacheDir, p.G1SRS, p.KzgConfig.NumWorker)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create SRS table: %w", err)
	}

	toeplitzLength := provingParams.ToeplitzSquareMatrixLength()

	fftPoints, err := subTable.GetSubTables(toeplitzLength, provingParams.ChunkLength)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get SRS table: %w", err)
	}

	// TODO(samlaf): if we only use the transposed points in MultiProof,
	// why didn't we store the SRSTables in transposed form?
	fftPointsT := make([][]bn254.G1Affine, len(fftPoints[0]))
	for i := range fftPointsT {
		fftPointsT[i] = make([]bn254.G1Affine, len(fftPoints))
		for j := uint64(0); j < provingParams.ChunkLength; j++ {
			fftPointsT[i][j] = fftPoints[j][i]
		}
	}
	return fftPoints, fftPointsT, nil
}
