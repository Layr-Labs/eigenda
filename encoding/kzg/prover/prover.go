package prover

import (
	"errors"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/Layr-Labs/eigenda/encoding/rs"
	"github.com/consensys/gnark-crypto/ecc/bn254"

	_ "go.uber.org/automaxprocs"
)

type Prover struct {
	*kzg.KzgConfig
	Srs          *kzg.SRS
	G2Trailing   []bn254.G2Affine
	mu           sync.Mutex
	LoadG2Points bool

	ParametrizedProvers map[encoding.EncodingParams]*ParametrizedProver
}

var _ encoding.Prover = &Prover{}

func NewProver(config *kzg.KzgConfig, loadG2Points bool) (*Prover, error) {
	if config.SRSNumberToLoad > config.SRSOrder {
		return nil, errors.New("SRSOrder is less than srsNumberToLoad")
	}

	// read the whole order, and treat it as entire SRS for low degree proof
	s1, err := kzg.ReadG1Points(config.G1Path, config.SRSNumberToLoad, config.NumWorker)
	if err != nil {
		log.Println("failed to read G1 points", err)
		return nil, err
	}

	s2 := make([]bn254.G2Affine, 0)
	g2Trailing := make([]bn254.G2Affine, 0)

	// PreloadEncoder is by default not used by operator node, PreloadEncoder
	if loadG2Points {
		if len(config.G2Path) == 0 {
			return nil, errors.New("G2Path is empty. However, object needs to load G2Points")
		}

		s2, err = kzg.ReadG2Points(config.G2Path, config.SRSNumberToLoad, config.NumWorker)
		if err != nil {
			log.Println("failed to read G2 points", err)
			return nil, err
		}

		g2Trailing, err = kzg.ReadG2PointSection(
			config.G2Path,
			config.SRSOrder-config.SRSNumberToLoad,
			config.SRSOrder, // last exclusive
			config.NumWorker,
		)
		if err != nil {
			return nil, err
		}
	} else {
		// todo, there are better ways to handle it
		if len(config.G2PowerOf2Path) == 0 {
			return nil, errors.New("G2PowerOf2Path is empty. However, object needs to load G2Points")
		}
	}

	srs, err := kzg.NewSrs(s1, s2)
	if err != nil {
		log.Println("Could not create srs", err)
		return nil, err
	}

	fmt.Println("numthread", runtime.GOMAXPROCS(0))

	encoderGroup := &Prover{
		KzgConfig:           config,
		Srs:                 srs,
		G2Trailing:          g2Trailing,
		ParametrizedProvers: make(map[encoding.EncodingParams]*ParametrizedProver),
		LoadG2Points:        loadG2Points,
	}

	if config.PreloadEncoder {
		// create table dir if not exist
		err := os.MkdirAll(config.CacheDir, os.ModePerm)
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
	paramsAll, err := GetAllPrecomputedSrsMap(g.CacheDir)
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

func (e *Prover) GetCommitments(data []byte) (encoding.BlobCommitments, error) {
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

	commit, lengthCommit, lengthProof, err := enc.GetCommitments(symbols)
	if err != nil {
		return encoding.BlobCommitments{}, err
	}

	length := uint(len(symbols))
	commitments := encoding.BlobCommitments{
		Commitment:       (*encoding.G1Commitment)(commit),
		LengthCommitment: (*encoding.G2Commitment)(lengthCommit),
		LengthProof:      (*encoding.G2Commitment)(lengthProof),
		Length:           length,
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
