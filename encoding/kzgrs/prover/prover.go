package prover

import (
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/kzgrs"
	"github.com/Layr-Labs/eigenda/encoding/rs"
	kzg "github.com/Layr-Labs/eigenda/pkg/kzg"
	"github.com/consensys/gnark-crypto/ecc/bn254"

	_ "go.uber.org/automaxprocs"
)

type Prover struct {
	*kzgrs.KzgConfig
	Srs          *kzg.SRS
	G2Trailing   []bn254.G2Affine
	mu           sync.Mutex
	LoadG2Points bool

	ParametrizedProvers map[encoding.EncodingParams]*ParametrizedProver
}

var _ encoding.Prover = &Prover{}

func NewProver(config *kzgrs.KzgConfig, loadG2Points bool) (*Prover, error) {

	if config.SRSNumberToLoad > config.SRSOrder {
		return nil, errors.New("SRSOrder is less than srsNumberToLoad")
	}

	// read the whole order, and treat it as entire SRS for low degree proof
	s1, err := kzgrs.ReadG1Points(config.G1Path, config.SRSNumberToLoad, config.NumWorker)
	if err != nil {
		log.Println("failed to read G1 points", err)
		return nil, err
	}

	s2 := make([]bn254.G2Affine, 0)
	g2Trailing := make([]bn254.G2Affine, 0)

	// PreloadEncoder is by default not used by operator node, PreloadEncoder
	if loadG2Points {
		if len(config.G2Path) == 0 {
			return nil, fmt.Errorf("G2Path is empty. However, object needs to load G2Points")
		}

		s2, err = kzgrs.ReadG2Points(config.G2Path, config.SRSNumberToLoad, config.NumWorker)
		if err != nil {
			log.Println("failed to read G2 points", err)
			return nil, err
		}

		g2Trailing, err = kzgrs.ReadG2PointSection(
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
			return nil, fmt.Errorf("G2PowerOf2Path is empty. However, object needs to load G2Points")
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
		return nil
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

	commit, lowDegreeCommit, lowDegreeProof, kzgFrames, _, err := enc.EncodeBytes(data)
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

	length := uint(len(rs.ToFrArray(data)))
	commitments := encoding.BlobCommitments{
		Commitment:       (*encoding.G1Commitment)(commit),
		LengthCommitment: (*encoding.G2Commitment)(lowDegreeCommit),
		LengthProof:      (*encoding.G2Commitment)(lowDegreeProof),
		Length:           length,
	}

	return commitments, chunks, nil
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

func (g *Prover) newProver(params encoding.EncodingParams) (*ParametrizedProver, error) {

	// Check that the parameters are valid with respect to the SRS.
	if params.ChunkLength*params.NumChunks >= g.SRSOrder {
		return nil, fmt.Errorf("the supplied encoding parameters are not valid with respect to the SRS. ChunkLength: %d, NumChunks: %d, SRSOrder: %d", params.ChunkLength, params.NumChunks, g.SRSOrder)
	}

	encoder, err := rs.NewEncoder(params, g.Verbose)
	if err != nil {
		log.Println("Could not create encoder: ", err)
		return nil, err
	}

	subTable, err := NewSRSTable(g.CacheDir, g.Srs.G1, g.NumWorker)
	if err != nil {
		log.Println("Could not create srs table:", err)
		return nil, err
	}

	fftPoints, err := subTable.GetSubTables(encoder.NumChunks, encoder.ChunkLength)
	if err != nil {
		log.Println("could not get sub tables", err)
		return nil, err
	}

	fftPointsT := make([][]bn254.G1Affine, len(fftPoints[0]))
	for i := range fftPointsT {
		fftPointsT[i] = make([]bn254.G1Affine, len(fftPoints))
		for j := uint64(0); j < encoder.ChunkLength; j++ {
			fftPointsT[i][j] = fftPoints[j][i]
		}
	}
	n := uint8(math.Log2(float64(encoder.NumEvaluations())))
	if encoder.ChunkLength == 1 {
		n = uint8(math.Log2(float64(2 * encoder.NumChunks)))
	}
	fs := kzg.NewFFTSettings(n)

	ks, err := kzg.NewKZGSettings(fs, g.Srs)
	if err != nil {
		return nil, err
	}

	t := uint8(math.Log2(float64(2 * encoder.NumChunks)))
	sfs := kzg.NewFFTSettings(t)

	return &ParametrizedProver{
		Encoder:    encoder,
		KzgConfig:  g.KzgConfig,
		Srs:        g.Srs,
		G2Trailing: g.G2Trailing,
		Fs:         fs,
		Ks:         ks,
		SFs:        sfs,
		FFTPoints:  fftPoints,
		FFTPointsT: fftPointsT,
	}, nil
}

// get Fiat-Shamir challenge
// func createFiatShamirChallenge(byteArray [][32]byte) *fr.Element {
// 	alphaBytesTmp := make([]byte, 0)
// 	for i := 0; i < len(byteArray); i++ {
// 		for j := 0; j < len(byteArray[i]); j++ {
// 			alphaBytesTmp = append(alphaBytesTmp, byteArray[i][j])
// 		}
// 	}
// 	alphaBytes := crypto.Keccak256(alphaBytesTmp)
// 	alpha := new(fr.Element)
// 	fr.ElementSetBytes(alpha, alphaBytes)
//
// 	return alpha
// }

// invert the divisor, then multiply
// func polyFactorDiv(dst *fr.Element, a *fr.Element, b *fr.Element) {
// 	// TODO: use divmod instead.
// 	var tmp fr.Element
// 	bls.InvModFr(&tmp, b)
// 	bls.MulModFr(dst, &tmp, a)
// }

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
