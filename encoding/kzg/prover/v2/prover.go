package prover

import (
	"errors"
	"fmt"
	gomath "math"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/Layr-Labs/eigenda/common/math"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/fft"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	gnarkprover "github.com/Layr-Labs/eigenda/encoding/kzg/prover/v2/gnark"
	"github.com/Layr-Labs/eigenda/encoding/rs"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	_ "go.uber.org/automaxprocs"
)

type Prover struct {
	Config     *encoding.Config
	KzgConfig  *KzgConfig
	encoder    *rs.Encoder
	Srs        kzg.SRS
	G2Trailing []bn254.G2Affine

	// mu protects access to ParametrizedProvers
	mu                  sync.Mutex
	ParametrizedProvers map[encoding.EncodingParams]*ParametrizedProver
}

func NewProver(kzgConfig *KzgConfig, encoderConfig *encoding.Config) (*Prover, error) {
	if encoderConfig == nil {
		encoderConfig = encoding.DefaultConfig()
	}

	if kzgConfig.SRSNumberToLoad > encoding.SRSOrder {
		return nil, errors.New("SRSOrder is less than srsNumberToLoad")
	}

	// read the whole order, and treat it as entire SRS for low degree proof
	s1, err := kzg.ReadG1Points(kzgConfig.G1Path, kzgConfig.SRSNumberToLoad, kzgConfig.NumWorker)
	if err != nil {
		return nil, fmt.Errorf("failed to read G1 points: %w", err)
	}

	s2 := make([]bn254.G2Affine, 0)
	g2Trailing := make([]bn254.G2Affine, 0)

	// PreloadEncoder is by default not used by operator node, PreloadEncoder
	if kzgConfig.LoadG2Points { //nolint: nestif
		if len(kzgConfig.G2Path) == 0 {
			return nil, errors.New("G2Path is empty. However, object needs to load G2Points")
		}

		s2, err = kzg.ReadG2Points(kzgConfig.G2Path, kzgConfig.SRSNumberToLoad, kzgConfig.NumWorker)
		if err != nil {
			return nil, fmt.Errorf("failed to read G2 points: %w", err)
		}

		hasG2TrailingFile := len(kzgConfig.G2TrailingPath) != 0
		if hasG2TrailingFile {
			// TODO(samlaf): this function/check should probably be done in ReadG2PointSection
			numG2point, err := kzg.NumberOfPointsInSRSFile(kzgConfig.G2TrailingPath, kzg.G2PointBytes)
			if err != nil {
				return nil, fmt.Errorf("number of points in srs file %v: %w", kzgConfig.G2TrailingPath, err)
			}
			if numG2point < kzgConfig.SRSNumberToLoad {
				return nil, fmt.Errorf("kzgConfig.G2TrailingPath=%v contains %v G2 Points, "+
					"which is < kzgConfig.SRSNumberToLoad=%v",
					kzgConfig.G2TrailingPath, numG2point, kzgConfig.SRSNumberToLoad)
			}

			// use g2 trailing file
			g2Trailing, err = kzg.ReadG2PointSection(
				kzgConfig.G2TrailingPath,
				numG2point-kzgConfig.SRSNumberToLoad,
				numG2point, // last exclusive
				kzgConfig.NumWorker,
			)
			if err != nil {
				return nil, fmt.Errorf("failed to read G2 trailing points (%v to %v) from file %v: %w",
					numG2point-kzgConfig.SRSNumberToLoad, numG2point, kzgConfig.G2TrailingPath, err)
			}
		} else {
			// require entire g2 srs be available on disk
			numG2point, err := kzg.NumberOfPointsInSRSFile(kzgConfig.G2Path, kzg.G2PointBytes)
			if err != nil {
				return nil, fmt.Errorf("number of points in srs file: %w", err)
			}
			if numG2point < encoding.SRSOrder {
				return nil, fmt.Errorf("no kzgConfig.G2TrailingPath was passed, yet the G2 SRS file %v is incomplete: contains %v < 2^28 G2 Points", kzgConfig.G2Path, numG2point)
			}
			g2Trailing, err = kzg.ReadG2PointSection(
				kzgConfig.G2Path,
				encoding.SRSOrder-kzgConfig.SRSNumberToLoad,
				encoding.SRSOrder, // last exclusive
				kzgConfig.NumWorker,
			)
			if err != nil {
				return nil, fmt.Errorf("failed to read G2 points (%v to %v) from file %v: %w",
					encoding.SRSOrder-kzgConfig.SRSNumberToLoad, encoding.SRSOrder, kzgConfig.G2Path, err)
			}
		}
	}

	srs := kzg.NewSrs(s1, s2)

	// Create RS encoder
	rsEncoder, err := rs.NewEncoder(encoderConfig)
	if err != nil {
		return nil, fmt.Errorf("create rs encoder: %w", err)
	}

	encoderGroup := &Prover{
		Config:              encoderConfig,
		encoder:             rsEncoder,
		KzgConfig:           kzgConfig,
		Srs:                 srs,
		G2Trailing:          g2Trailing,
		ParametrizedProvers: make(map[encoding.EncodingParams]*ParametrizedProver),
	}

	if kzgConfig.PreloadEncoder {
		// create table dir if not exist
		err := os.MkdirAll(kzgConfig.CacheDir, os.ModePerm)
		if err != nil {
			return nil, fmt.Errorf("make cache dir: %w", err)
		}

		err = encoderGroup.preloadAllEncoders()
		if err != nil {
			return nil, fmt.Errorf("preload all encoders: %w", err)
		}
	}

	return encoderGroup, nil
}

func (e *Prover) EncodeAndProve(
	data []byte, params encoding.EncodingParams,
) (encoding.BlobCommitments, []*encoding.Frame, error) {
	enc, err := e.GetKzgEncoder(params)
	if err != nil {
		return encoding.BlobCommitments{}, nil, err
	}

	commit, lengthCommit, lengthProof, kzgFrames, _, err := enc.EncodeBytes(data)
	if err != nil {
		return encoding.BlobCommitments{}, nil, err
	}

	chunks := make([]*encoding.Frame, len(kzgFrames))
	for ind, frame := range kzgFrames {

		chunks[ind] = &encoding.Frame{
			Coeffs: frame.Coeffs,
			Proof:  frame.Proof,
		}
	}

	symbols, err := rs.ToFrArray(data)
	if err != nil {
		return encoding.BlobCommitments{}, nil, fmt.Errorf("ToFrArray: %w", err)
	}

	length := uint(len(symbols))
	commitments := encoding.BlobCommitments{
		Commitment:       (*encoding.G1Commitment)(commit),
		LengthCommitment: (*encoding.G2Commitment)(lengthCommit),
		LengthProof:      (*encoding.G2Commitment)(lengthProof),
		Length:           length,
	}

	return commitments, chunks, nil
}

func (e *Prover) GetFrames(data []byte, params encoding.EncodingParams) ([]*encoding.Frame, error) {
	symbols, err := rs.ToFrArray(data)
	if err != nil {
		return nil, fmt.Errorf("ToFrArray: %w", err)
	}

	enc, err := e.GetKzgEncoder(params)
	if err != nil {
		return nil, fmt.Errorf("get kzg encoder: %w", err)
	}

	kzgFrames, _, err := enc.GetFrames(symbols)
	if err != nil {
		return nil, fmt.Errorf("get frames: %w", err)
	}

	chunks := make([]*encoding.Frame, len(kzgFrames))
	for ind, frame := range kzgFrames {
		chunks[ind] = &encoding.Frame{
			Coeffs: frame.Coeffs,
			Proof:  frame.Proof,
		}
	}

	return chunks, nil
}

// GetCommitmentsForPaddedLength takes in a byte slice representing a list of bn254
// field elements (32 bytes each, except potentially the last element),
// pads the (potentially incomplete) last element with zeroes, and returns the commitments for the padded list.
func (e *Prover) GetCommitmentsForPaddedLength(data []byte) (encoding.BlobCommitments, error) {
	symbols, err := rs.ToFrArray(data)
	if err != nil {
		return encoding.BlobCommitments{}, fmt.Errorf("ToFrArray: %w", err)
	}

	params := encoding.EncodingParams{
		NumChunks:   2,
		ChunkLength: 2,
	}

	enc, err := e.GetKzgEncoder(params)
	if err != nil {
		return encoding.BlobCommitments{}, fmt.Errorf("get kzg encoder: %w", err)
	}

	commit, lengthCommit, lengthProof, err := enc.GetCommitments(symbols)
	if err != nil {
		return encoding.BlobCommitments{}, fmt.Errorf("get commitments: %w", err)
	}

	commitments := encoding.BlobCommitments{
		Commitment:       (*encoding.G1Commitment)(commit),
		LengthCommitment: (*encoding.G2Commitment)(lengthCommit),
		LengthProof:      (*encoding.G2Commitment)(lengthProof),
		Length:           uint(math.NextPowOf2u32(uint32(len(symbols)))),
	}

	return commitments, nil
}

func (g *Prover) GetKzgEncoder(params encoding.EncodingParams) (*ParametrizedProver, error) {
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

// Decode takes in the chunks, indices, and encoding parameters and returns the decoded blob
// The result is trimmed to the given maxInputSize.
func (p *Prover) Decode(
	chunks []*encoding.Frame, indices []encoding.ChunkNumber, params encoding.EncodingParams, maxInputSize uint64,
) ([]byte, error) {
	frames := make([]encoding.Frame, len(chunks))
	for i := range chunks {
		frames[i] = encoding.Frame{
			Proof:  chunks[i].Proof,
			Coeffs: chunks[i].Coeffs,
		}
	}

	encoder, err := p.GetKzgEncoder(params)
	if err != nil {
		return nil, err
	}

	return encoder.Decode(frames, toUint64Array(indices), maxInputSize)
}

func toUint64Array(chunkIndices []encoding.ChunkNumber) []uint64 {
	res := make([]uint64, len(chunkIndices))
	for i, d := range chunkIndices {
		res[i] = uint64(d)
	}
	return res
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

	switch p.Config.BackendType {
	case encoding.GnarkBackend:
		return p.createGnarkBackendProver(params, fs)
	case encoding.IcicleBackend:
		return p.createIcicleBackendProver(params, fs)
	default:
		return nil, fmt.Errorf("unsupported backend type: %v", p.Config.BackendType)
	}

}

func (p *Prover) createGnarkBackendProver(
	params encoding.EncodingParams, fs *fft.FFTSettings,
) (*ParametrizedProver, error) {
	if p.Config.GPUEnable {
		return nil, errors.New("GPU is not supported in gnark backend")
	}

	_, fftPointsT, err := p.setupFFTPoints(params)
	if err != nil {
		return nil, err
	}

	// Create subgroup FFT settings
	t := uint8(gomath.Log2(float64(2 * params.NumChunks)))
	sfs := fft.NewFFTSettings(t)

	// Set KZG Prover gnark backend
	multiproofBackend := &gnarkprover.KzgMultiProofGnarkBackend{
		Fs:         fs,
		FFTPointsT: fftPointsT,
		SFs:        sfs,
	}

	// Set KZG Commitments gnark backend
	commitmentsBackend := &gnarkprover.KzgCommitmentsGnarkBackend{
		Srs:        p.Srs,
		G2Trailing: p.G2Trailing,
	}

	return &ParametrizedProver{
		srsNumberToLoad:           p.KzgConfig.SRSNumberToLoad,
		encoder:                   p.encoder,
		encodingParams:            params,
		computeMultiprooNumWorker: p.KzgConfig.NumWorker,
		kzgMultiProofBackend:      multiproofBackend,
		kzgCommitmentsBackend:     commitmentsBackend,
	}, nil
}

func (p *Prover) createIcicleBackendProver(
	params encoding.EncodingParams, fs *fft.FFTSettings,
) (*ParametrizedProver, error) {
	return CreateIcicleBackendProver(p, params, fs)
}

func (g *Prover) preloadAllEncoders() error {
	paramsAll, err := getAllPrecomputedSrsMap(g.KzgConfig.CacheDir)
	if err != nil {
		return err
	}
	fmt.Printf("detect %v srs maps\n", len(paramsAll))
	for i := 0; i < len(paramsAll); i++ {
		fmt.Printf(" %v. NumChunks: %v   ChunkLength: %v\n", i, paramsAll[i].NumChunks, paramsAll[i].ChunkLength)
	}

	if len(paramsAll) == 0 {
		return nil
	}

	for _, params := range paramsAll {
		// get those encoders and store them
		enc, err := g.GetKzgEncoder(params)
		if err != nil {
			return err
		}
		g.ParametrizedProvers[params] = enc
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
			NumChunks:   uint64(cosetSizeValue),
			ChunkLength: uint64(dimEValue),
		}
		tables = append(tables, params)
	}
	return tables, nil
}

// Helper methods for setup
func (p *Prover) setupFFTPoints(params encoding.EncodingParams) ([][]bn254.G1Affine, [][]bn254.G1Affine, error) {
	subTable, err := NewSRSTable(p.KzgConfig.CacheDir, p.Srs.G1, p.KzgConfig.NumWorker)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create SRS table: %w", err)
	}

	fftPoints, err := subTable.GetSubTables(params.NumChunks, params.ChunkLength)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get sub tables: %w", err)
	}

	fftPointsT := make([][]bn254.G1Affine, len(fftPoints[0]))
	for i := range fftPointsT {
		fftPointsT[i] = make([]bn254.G1Affine, len(fftPoints))
		for j := uint64(0); j < params.ChunkLength; j++ {
			fftPointsT[i][j] = fftPoints[j][i]
		}
	}

	return fftPoints, fftPointsT, nil
}
