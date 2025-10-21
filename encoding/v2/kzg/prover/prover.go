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

	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/v2/fft"
	"github.com/Layr-Labs/eigenda/encoding/v2/kzg"
	"github.com/Layr-Labs/eigenda/encoding/v2/kzg/prover/backend"
	"github.com/Layr-Labs/eigenda/encoding/v2/kzg/prover/backend/gnark"
	"github.com/Layr-Labs/eigenda/encoding/v2/kzg/prover/backend/icicle"
	"github.com/Layr-Labs/eigenda/encoding/v2/rs"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	_ "go.uber.org/automaxprocs"
)

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
	ParametrizedProvers map[encoding.EncodingParams]*ParametrizedProver
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
		ParametrizedProvers: make(map[encoding.EncodingParams]*ParametrizedProver),
	}

	if kzgConfig.PreloadEncoder {
		// create table dir if not exist
		err := os.MkdirAll(kzgConfig.CacheDir, os.ModePerm)
		if err != nil {
			return nil, fmt.Errorf("make cache dir: %w", err)
		}

		err = proverGroup.preloadProversFromSRSTableCache()
		if err != nil {
			return nil, fmt.Errorf("preload all provers: %w", err)
		}
	}

	return proverGroup, nil
}

func (e *Prover) GetFrames(data []byte, params encoding.EncodingParams) ([]*encoding.Frame, []uint32, error) {
	symbols, err := rs.ToFrArray(data)
	if err != nil {
		return nil, nil, fmt.Errorf("ToFrArray: %w", err)
	}

	prover, err := e.GetKzgProver(params)
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
		frames, indices, err := e.encoder.Encode(symbols, params)
		encodingDuration := time.Since(encodeStart)
		encodeChan <- encodeChanResult{
			chunks:   frames,
			indices:  indices,
			duration: encodingDuration,
			err:      err,
		}
	}()

	getProofsStart := time.Now()
	proofs, err := prover.GetProofs(symbols)
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
		"input_size_bytes", len(symbols)*encoding.BYTES_PER_SYMBOL,
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

func (g *Prover) GetKzgProver(params encoding.EncodingParams) (*ParametrizedProver, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	enc, ok := g.ParametrizedProvers[params]
	if ok {
		return enc, nil
	}

	enc, err := g.newProver(params)
	if err != nil {
		return nil, fmt.Errorf("new prover: %w", err)
	}

	g.ParametrizedProvers[params] = enc
	return enc, nil
}

func (p *Prover) newProver(params encoding.EncodingParams) (*ParametrizedProver, error) {
	if err := encoding.ValidateEncodingParams(params, encoding.SRSOrder); err != nil {
		return nil, fmt.Errorf("validate encoding params: %w", err)
	}

	// Create FFT settings based on params
	n := uint8(gomath.Log2(float64(params.NumEvaluations())))
	if params.ChunkLength == 1 {
		n = uint8(gomath.Log2(float64(2 * params.NumChunks)))
	}
	fs := fft.NewFFTSettings(n)

	_, fftPointsT, err := p.setupFFTPoints(params)
	if err != nil {
		return nil, fmt.Errorf("setup fft points: %w", err)
	}

	var multiproofsBackend backend.KzgMultiProofsBackendV2
	switch p.Config.BackendType {
	case encoding.GnarkBackend:
		if p.Config.GPUEnable {
			return nil, errors.New("GPU is not supported in gnark backend")
		}
		multiproofsBackend = gnark.NewMultiProofBackend(p.logger, fs, fftPointsT)
	case encoding.IcicleBackend:
		multiproofsBackend, err = icicle.NewMultiProofBackend(
			p.logger, fs, fftPointsT, p.Config.GPUEnable, p.KzgConfig.NumWorker)
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

func (g *Prover) preloadProversFromSRSTableCache() error {
	paramsAll, err := getAllPrecomputedSrsMap(g.KzgConfig.CacheDir)
	if err != nil {
		return err
	}
	g.logger.Info("Detected SRSTables from cache dir", "NumTables", len(paramsAll), "TableDetails", paramsAll)

	if len(paramsAll) == 0 {
		return nil
	}

	for _, params := range paramsAll {
		prover, err := g.GetKzgProver(params)
		if err != nil {
			return err
		}
		g.ParametrizedProvers[params] = prover
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
func getAllPrecomputedSrsMap(tableDir string) ([]encoding.EncodingParams, error) {
	files, err := os.ReadDir(tableDir)
	if err != nil {
		return nil, fmt.Errorf("read srs table dir: %w", err)
	}

	tables := make([]encoding.EncodingParams, 0)
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

		params := encoding.EncodingParams{
			NumChunks:   uint64(dimEValue),
			ChunkLength: uint64(cosetSizeValue),
		}
		tables = append(tables, params)
	}
	return tables, nil
}

// Returns SRSTable SRS points, as well as its transpose.
// fftPoints has size [l][2*dimE], and its transpose has size [2*dimE][l]
func (p *Prover) setupFFTPoints(params encoding.EncodingParams) ([][]bn254.G1Affine, [][]bn254.G1Affine, error) {
	subTable, err := NewSRSTable(p.logger, p.KzgConfig.CacheDir, p.G1SRS, p.KzgConfig.NumWorker)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create SRS table: %w", err)
	}

	fftPoints, err := subTable.GetSubTables(params.NumChunks, params.ChunkLength)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get sub tables: %w", err)
	}

	// TODO(samlaf): if we only use the transposed points in MultiProof,
	// why didn't we store the SRSTables in transposed form?
	fftPointsT := make([][]bn254.G1Affine, len(fftPoints[0]))
	for i := range fftPointsT {
		fftPointsT[i] = make([]bn254.G1Affine, len(fftPoints))
		for j := uint64(0); j < params.ChunkLength; j++ {
			fftPointsT[i][j] = fftPoints[j][i]
		}
	}

	return fftPoints, fftPointsT, nil
}
