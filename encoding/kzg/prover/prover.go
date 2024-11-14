package prover

import (
	"errors"
	"fmt"
	"log"
	"log/slog"
	"math"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/fft"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/Layr-Labs/eigenda/encoding/rs"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	_ "go.uber.org/automaxprocs"
)

// ProverOption defines a function that configures a Prover
type ProverOption func(*Prover) error

type Prover struct {
	Config    *encoding.Config
	KzgConfig *kzg.KzgConfig
	*rs.Encoder
	encoding.BackendType
	Srs          *kzg.SRS
	G2Trailing   []bn254.G2Affine
	mu           sync.Mutex
	LoadG2Points bool

	ParametrizedProvers map[encoding.EncodingParams]*ParametrizedProver
}

var _ encoding.Prover = &Prover{}

// Default configuration values
const (
	defaultBackend        = encoding.BackendDefault
	defaultGPUEnable      = false
	defaultLoadG2Points   = true
	defaultPreloadEncoder = false
	defaultNTTSize        = 25 // Used for NTT setup in Icicle backend
	defaultVerbose        = false
)

// WithBackend sets the backend type for the prover
func WithBackend(backend encoding.BackendType) ProverOption {
	return func(p *Prover) error {
		p.Config.BackendType = backend
		return nil
	}
}

// WithGPU enables or disables GPU usage
func WithGPU(enable bool) ProverOption {
	return func(e *Prover) error {
		e.Config.GPUEnable = enable
		return nil
	}
}

// WithKZGConfig sets the KZG configuration
func WithKZGConfig(config *kzg.KzgConfig) ProverOption {
	return func(p *Prover) error {
		if config.SRSNumberToLoad > config.SRSOrder {
			return errors.New("SRSOrder is less than srsNumberToLoad")
		}
		p.KzgConfig = config
		return nil
	}
}

// WithRSEncoder sets a custom RS encoder
func WithRSEncoder(encoder *rs.Encoder) ProverOption {
	return func(p *Prover) error {
		p.Encoder = encoder
		return nil
	}
}

// WithRSEncoderOptions configures the RS encoder with specific options
func WithRSEncoderOptions(opts ...rs.EncoderOption) ProverOption {
	return func(p *Prover) error {
		encoder, err := rs.NewEncoder(opts...)
		if err != nil {
			return fmt.Errorf("failed to create RS encoder: %w", err)
		}
		p.Encoder = encoder
		return nil
	}
}

// WithLoadG2Points enables or disables G2 points loading
func WithLoadG2Points(load bool) ProverOption {
	return func(p *Prover) error {
		p.LoadG2Points = load
		return nil
	}
}

// WithPreloadEncoder enables or disables encoder preloading
func WithPreloadEncoder(preload bool) ProverOption {
	return func(p *Prover) error {
		if !preload {
			return nil
		}

		if p.KzgConfig == nil {
			return errors.New("KZG config must be set before enabling preload encoder")
		}

		// Create table dir if not exist
		err := os.MkdirAll(p.KzgConfig.CacheDir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("cannot make CacheDir: %w", err)
		}

		return p.PreloadAllEncoders()
	}
}

// WithVerbose enables or disables verbose logging
func WithVerbose(verbose bool) ProverOption {
	return func(p *Prover) error {
		p.Config.Verbose = verbose
		return nil
	}
}

func NewProver(opts ...ProverOption) (*Prover, error) {
	p := &Prover{
		Config: &encoding.Config{
			NumWorker:   uint64(runtime.GOMAXPROCS(0)),
			BackendType: defaultBackend,
			GPUEnable:   defaultGPUEnable,
			Verbose:     defaultVerbose,
		},

		LoadG2Points:        defaultLoadG2Points,
		ParametrizedProvers: make(map[encoding.EncodingParams]*ParametrizedProver),
		mu:                  sync.Mutex{},
	}

	// Apply options
	for _, opt := range opts {
		if err := opt(p); err != nil {
			return nil, fmt.Errorf("failed to apply option: %w", err)
		}
	}

	// Validate required configurations
	if p.KzgConfig == nil {
		return nil, errors.New("KZG config is required")
	}

	if err := p.initializeSRS(); err != nil {
		return nil, err
	}

	// Create default RS encoder if none provided
	if p.Encoder == nil {
		encoder, err := rs.NewEncoder()
		if err != nil {
			return nil, fmt.Errorf("failed to create default RS encoder: %w", err)
		}
		p.Encoder = encoder
	}

	return p, nil
}

func (p *Prover) initializeSRS() error {
	startTime := time.Now()
	s1, err := kzg.ReadG1Points(p.KzgConfig.G1Path, p.KzgConfig.SRSNumberToLoad, p.KzgConfig.NumWorker)
	if err != nil {
		return fmt.Errorf("failed to read G1 points: %w", err)
	}
	slog.Info("ReadG1Points", "time", time.Since(startTime), "numPoints", len(s1))

	s2 := make([]bn254.G2Affine, 0)
	g2Trailing := make([]bn254.G2Affine, 0)

	if p.LoadG2Points {
		if err := p.loadG2PointsData(&s2, &g2Trailing); err != nil {
			return err
		}
	} else if len(p.KzgConfig.G2PowerOf2Path) == 0 {
		return errors.New("G2PowerOf2Path is empty but required when loadG2Points is false")
	}

	srs, err := kzg.NewSrs(s1, s2)
	if err != nil {
		return fmt.Errorf("could not create srs: %w", err)
	}

	p.Srs = srs
	p.G2Trailing = g2Trailing

	return nil
}

func (p *Prover) loadG2PointsData(s2 *[]bn254.G2Affine, g2Trailing *[]bn254.G2Affine) error {
	if len(p.KzgConfig.G2Path) == 0 {
		return errors.New("G2Path is empty but required when loadG2Points is true")
	}

	startTime := time.Now()
	points, err := kzg.ReadG2Points(p.KzgConfig.G2Path, p.KzgConfig.SRSNumberToLoad, p.KzgConfig.NumWorker)
	if err != nil {
		return fmt.Errorf("failed to read G2 points: %w", err)
	}
	slog.Info("ReadG2Points", "time", time.Since(startTime), "numPoints", len(points))
	*s2 = points

	startTime = time.Now()
	trailing, err := kzg.ReadG2PointSection(
		p.KzgConfig.G2Path,
		p.KzgConfig.SRSOrder-p.KzgConfig.SRSNumberToLoad,
		p.KzgConfig.SRSOrder,
		p.KzgConfig.NumWorker,
	)
	if err != nil {
		return fmt.Errorf("failed to read G2 point section: %w", err)
	}
	slog.Info("ReadG2PointSection", "time", time.Since(startTime), "numPoints", len(trailing))
	*g2Trailing = trailing

	return nil
}

func (g *Prover) PreloadAllEncoders() error {
	paramsAll, err := GetAllPrecomputedSrsMap(g.KzgConfig.CacheDir)
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

func (e *Prover) EncodeAndProve(data []byte, params encoding.EncodingParams) (encoding.BlobCommitments, []*encoding.Frame, error) {
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
		return encoding.BlobCommitments{}, nil, err
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
		return nil, err
	}

	enc, err := e.GetKzgEncoder(params)
	if err != nil {
		return nil, err
	}

	kzgFrames, _, err := enc.GetFrames(symbols)
	if err != nil {
		return nil, err
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

func (e *Prover) GetCommitmentsForPaddedLength(data []byte) (encoding.BlobCommitments, error) {
	symbols, err := rs.ToFrArray(data)
	if err != nil {
		return encoding.BlobCommitments{}, err
	}

	params := encoding.EncodingParams{
		NumChunks:   2,
		ChunkLength: 2,
	}

	enc, err := e.GetKzgEncoder(params)
	if err != nil {
		return encoding.BlobCommitments{}, err
	}

	length := encoding.NextPowerOf2(uint64(len(symbols)))

	commit, lengthCommit, lengthProof, err := enc.GetCommitments(symbols, length)
	if err != nil {
		return encoding.BlobCommitments{}, err
	}

	commitments := encoding.BlobCommitments{
		Commitment:       (*encoding.G1Commitment)(commit),
		LengthCommitment: (*encoding.G2Commitment)(lengthCommit),
		LengthProof:      (*encoding.G2Commitment)(lengthProof),
		Length:           uint(length),
	}

	return commitments, nil
}

func (e *Prover) GetMultiFrameProofs(data []byte, params encoding.EncodingParams) ([]encoding.Proof, error) {
	symbols, err := rs.ToFrArray(data)
	if err != nil {
		return nil, err
	}

	enc, err := e.GetKzgEncoder(params)
	if err != nil {
		return nil, err
	}

	proofs, err := enc.GetMultiFrameProofs(symbols)
	if err != nil {
		return nil, err
	}

	return proofs, nil
}

func (g *Prover) GetKzgEncoder(params encoding.EncodingParams) (*ParametrizedProver, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	enc, ok := g.ParametrizedProvers[params]
	if ok {
		return enc, nil
	}

	enc, err := g.newProver(params)
	if err == nil {
		g.ParametrizedProvers[params] = enc
	}

	return enc, err
}

func (g *Prover) GetSRSOrder() uint64 {
	return g.SRSOrder
}

// Detect the precomputed table from the specified directory
// the file name follow the name convention of
//
//	dimE*.coset&
//
// where the first * specifies the dimension of the matrix which
// equals to the number of chunks
// where the second & specifies the length of each chunk
func GetAllPrecomputedSrsMap(tableDir string) ([]encoding.EncodingParams, error) {
	files, err := os.ReadDir(tableDir)
	if err != nil {
		log.Println("Error to list SRS Table directory", err)
		return nil, err
	}

	tables := make([]encoding.EncodingParams, 0)
	for _, file := range files {
		filename := file.Name()

		tokens := strings.Split(filename, ".")

		dimEValue, err := strconv.Atoi(tokens[0][4:])
		if err != nil {
			log.Println("Error to parse Dimension part of the Table", err)
			return nil, err
		}
		cosetSizeValue, err := strconv.Atoi(tokens[1][5:])
		if err != nil {
			log.Println("Error to parse Coset size of the Table", err)
			return nil, err
		}

		params := encoding.EncodingParams{
			NumChunks:   uint64(cosetSizeValue),
			ChunkLength: uint64(dimEValue),
		}
		tables = append(tables, params)
	}
	return tables, nil
}

// Decode takes in the chunks, indices, and encoding parameters and returns the decoded blob
// The result is trimmed to the given maxInputSize.
func (p *Prover) Decode(chunks []*encoding.Frame, indices []encoding.ChunkNumber, params encoding.EncodingParams, maxInputSize uint64) ([]byte, error) {
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
	if err := encoding.ValidateEncodingParams(params, p.KzgConfig.SRSOrder); err != nil {
		return nil, err
	}

	// Create FFT settings based on params
	n := uint8(math.Log2(float64(params.NumEvaluations())))
	if params.ChunkLength == 1 {
		n = uint8(math.Log2(float64(2 * params.NumChunks)))
	}
	fs := fft.NewFFTSettings(n)

	// Create base KZG settings
	ks, err := kzg.NewKZGSettings(fs, p.Srs)
	if err != nil {
		return nil, fmt.Errorf("failed to create KZG settings: %w", err)
	}

	switch p.Config.BackendType {

	case encoding.BackendDefault:
		return p.createDefaultBackendProver(params, fs, ks)
	case encoding.BackendIcicle:
		return p.createIcicleBackendProver(params, fs, ks)
	default:
		return nil, fmt.Errorf("unsupported backend type: %v", p.Config.BackendType)
	}

}

func (p *Prover) createDefaultBackendProver(params encoding.EncodingParams, fs *fft.FFTSettings, ks *kzg.KZGSettings) (*ParametrizedProver, error) {
	if p.Config.GPUEnable {
		return nil, fmt.Errorf("GPU is not supported in default backend")
	}

	_, fftPointsT, err := p.SetupFFTPoints(params)
	if err != nil {
		return nil, err
	}

	// Create subgroup FFT settings
	t := uint8(math.Log2(float64(2 * params.NumChunks)))
	sfs := fft.NewFFTSettings(t)

	// Set KZG Prover default backend
	multiproofBackend := &KzgMultiProofDefaultBackend{
		Fs:         fs,
		FFTPointsT: fftPointsT,
		SFs:        sfs,
		KzgConfig:  p.KzgConfig,
	}

	// Set KZG Commitments default backend
	commitmentsBckend := &KzgCommitmentsDefaultBackend{
		Srs:        p.Srs,
		G2Trailing: p.G2Trailing,
		KzgConfig:  p.KzgConfig,
	}

	return &ParametrizedProver{
		Encoder:               p.Encoder,
		EncodingParams:        params,
		KzgConfig:             p.KzgConfig,
		Ks:                    ks,
		KzgMultiProofBackend:  multiproofBackend,
		KzgCommitmentsBackend: commitmentsBckend,
	}, nil
}

func (p *Prover) createIcicleBackendProver(params encoding.EncodingParams, fs *fft.FFTSettings, ks *kzg.KZGSettings) (*ParametrizedProver, error) {
	return CreateIcicleBackendProver(p, params, fs, ks)
}

// Helper methods for setup
func (p *Prover) SetupFFTPoints(params encoding.EncodingParams) ([][]bn254.G1Affine, [][]bn254.G1Affine, error) {
	subTable, err := NewSRSTable(p.KzgConfig.CacheDir, p.Srs.G1, p.Config.NumWorker)
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
