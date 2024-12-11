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

	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/fft"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	gnarkprover "github.com/Layr-Labs/eigenda/encoding/kzg/prover/gnark"
	"github.com/Layr-Labs/eigenda/encoding/rs"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	_ "go.uber.org/automaxprocs"
)

type Prover struct {
	Config    *encoding.Config
	KzgConfig *kzg.KzgConfig
	encoder   *rs.Encoder
	encoding.BackendType
	Srs        *kzg.SRS
	G2Trailing []bn254.G2Affine
	mu         sync.Mutex

	ParametrizedProvers map[encoding.EncodingParams]*ParametrizedProver
}

var _ encoding.Prover = &Prover{}

func NewProver(kzgConfig *kzg.KzgConfig, encoderConfig *encoding.Config) (*Prover, error) {
	if encoderConfig == nil {
		encoderConfig = encoding.DefaultConfig()
	}

	if kzgConfig.SRSNumberToLoad > kzgConfig.SRSOrder {
		return nil, errors.New("SRSOrder is less than srsNumberToLoad")
	}

	// read the whole order, and treat it as entire SRS for low degree proof
	s1, err := kzg.ReadG1Points(kzgConfig.G1Path, kzgConfig.SRSNumberToLoad, kzgConfig.NumWorker)
	if err != nil {
		log.Println("failed to read G1 points", err)
		return nil, err
	}

	s2 := make([]bn254.G2Affine, 0)
	g2Trailing := make([]bn254.G2Affine, 0)

	// PreloadEncoder is by default not used by operator node, PreloadEncoder
	if kzgConfig.LoadG2Points {
		if len(kzgConfig.G2Path) == 0 {
			return nil, errors.New("G2Path is empty. However, object needs to load G2Points")
		}

		s2, err = kzg.ReadG2Points(kzgConfig.G2Path, kzgConfig.SRSNumberToLoad, kzgConfig.NumWorker)
		if err != nil {
			log.Println("failed to read G2 points", err)
			return nil, err
		}

		g2Trailing, err = kzg.ReadG2PointSection(
			kzgConfig.G2Path,
			kzgConfig.SRSOrder-kzgConfig.SRSNumberToLoad,
			kzgConfig.SRSOrder, // last exclusive
			kzgConfig.NumWorker,
		)
		if err != nil {
			return nil, err
		}
	} else {
		// todo, there are better ways to handle it
		if len(kzgConfig.G2PowerOf2Path) == 0 {
			return nil, errors.New("G2PowerOf2Path is empty. However, object needs to load G2Points")
		}
	}

	srs, err := kzg.NewSrs(s1, s2)
	if err != nil {
		log.Println("Could not create srs", err)
		return nil, err
	}

	fmt.Println("numthread", runtime.GOMAXPROCS(0))

	// Create RS encoder
	rsEncoder, err := rs.NewEncoder(encoderConfig)
	if err != nil {
		slog.Error("Could not create RS encoder", "err", err)
		return nil, err
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
			log.Println("Cannot make CacheDir", err)
			return nil, err
		}

		err = encoderGroup.PreloadAllEncoders()
		if err != nil {
			return nil, err
		}
	}

	return encoderGroup, nil
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
	return g.KzgConfig.SRSOrder
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
	case encoding.GnarkBackend:
		return p.createGnarkBackendProver(params, fs, ks)
	case encoding.IcicleBackend:
		return p.createIcicleBackendProver(params, fs, ks)
	default:
		return nil, fmt.Errorf("unsupported backend type: %v", p.Config.BackendType)
	}

}

func (p *Prover) createGnarkBackendProver(params encoding.EncodingParams, fs *fft.FFTSettings, ks *kzg.KZGSettings) (*ParametrizedProver, error) {
	if p.Config.GPUEnable {
		return nil, errors.New("GPU is not supported in gnark backend")
	}

	_, fftPointsT, err := p.SetupFFTPoints(params)
	if err != nil {
		return nil, err
	}

	// Create subgroup FFT settings
	t := uint8(math.Log2(float64(2 * params.NumChunks)))
	sfs := fft.NewFFTSettings(t)

	// Set KZG Prover gnark backend
	multiproofBackend := &gnarkprover.KzgMultiProofGnarkBackend{
		Fs:         fs,
		FFTPointsT: fftPointsT,
		SFs:        sfs,
		KzgConfig:  p.KzgConfig,
	}

	// Set KZG Commitments gnark backend
	commitmentsBackend := &gnarkprover.KzgCommitmentsGnarkBackend{
		Srs:        p.Srs,
		G2Trailing: p.G2Trailing,
		KzgConfig:  p.KzgConfig,
	}

	return &ParametrizedProver{
		Encoder:               p.encoder,
		EncodingParams:        params,
		KzgConfig:             p.KzgConfig,
		Ks:                    ks,
		KzgMultiProofBackend:  multiproofBackend,
		KzgCommitmentsBackend: commitmentsBackend,
	}, nil
}

func (p *Prover) createIcicleBackendProver(params encoding.EncodingParams, fs *fft.FFTSettings, ks *kzg.KZGSettings) (*ParametrizedProver, error) {
	return CreateIcicleBackendProver(p, params, fs, ks)
}

// Helper methods for setup
func (p *Prover) SetupFFTPoints(params encoding.EncodingParams) ([][]bn254.G1Affine, [][]bn254.G1Affine, error) {
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
