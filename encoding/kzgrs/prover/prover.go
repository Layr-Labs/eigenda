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

	"github.com/Layr-Labs/eigenda/encoding/kzgrs"
	"github.com/Layr-Labs/eigenda/encoding/rs"
	kzg "github.com/Layr-Labs/eigenda/pkg/kzg"

	bls "github.com/Layr-Labs/eigenda/pkg/kzg/bn254"
	_ "go.uber.org/automaxprocs"
)

type Prover struct {
	*kzgrs.KzgConfig
	Srs          *kzg.SRS
	G2Trailing   []bls.G2Point
	mu           sync.Mutex
	LoadG2Points bool

	ParametrizedProvers map[rs.EncodingParams]*ParametrizedProver
}

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

	s2 := make([]bls.G2Point, 0)
	g2Trailing := make([]bls.G2Point, 0)

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
		ParametrizedProvers: make(map[rs.EncodingParams]*ParametrizedProver),
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
		fmt.Printf(" %v. NumChunks: %v   ChunkLen: %v\n", i, paramsAll[i].NumChunks, paramsAll[i].ChunkLen)
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

func (g *Prover) GetKzgEncoder(params rs.EncodingParams) (*ParametrizedProver, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	enc, ok := g.ParametrizedProvers[params]
	if ok {
		return enc, nil
	}

	enc, err := g.newKzgEncoder(params)
	if err == nil {
		g.ParametrizedProvers[params] = enc
	}

	return enc, err
}

func (g *Prover) NewKzgEncoder(params rs.EncodingParams) (*ParametrizedProver, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	return g.newKzgEncoder(params)
}

func (g *Prover) newKzgEncoder(params rs.EncodingParams) (*ParametrizedProver, error) {

	// Check that the parameters are valid with respect to the SRS.
	if params.ChunkLen*params.NumChunks >= g.SRSOrder {
		return nil, fmt.Errorf("the supplied encoding parameters are not valid with respect to the SRS. ChunkLength: %d, NumChunks: %d, SRSOrder: %d", params.ChunkLen, params.NumChunks, g.SRSOrder)
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

	fftPoints, err := subTable.GetSubTables(encoder.NumChunks, encoder.ChunkLen)
	if err != nil {
		log.Println("could not get sub tables", err)
		return nil, err
	}

	fftPointsT := make([][]bls.G1Point, len(fftPoints[0]))
	for i := range fftPointsT {
		fftPointsT[i] = make([]bls.G1Point, len(fftPoints))
		for j := uint64(0); j < encoder.ChunkLen; j++ {
			fftPointsT[i][j] = fftPoints[j][i]
		}
	}
	n := uint8(math.Log2(float64(encoder.NumEvaluations())))
	if encoder.ChunkLen == 1 {
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
// func createFiatShamirChallenge(byteArray [][32]byte) *bls.Fr {
// 	alphaBytesTmp := make([]byte, 0)
// 	for i := 0; i < len(byteArray); i++ {
// 		for j := 0; j < len(byteArray[i]); j++ {
// 			alphaBytesTmp = append(alphaBytesTmp, byteArray[i][j])
// 		}
// 	}
// 	alphaBytes := crypto.Keccak256(alphaBytesTmp)
// 	alpha := new(bls.Fr)
// 	bls.FrSetBytes(alpha, alphaBytes)
//
// 	return alpha
// }

// invert the divisor, then multiply
// func polyFactorDiv(dst *bls.Fr, a *bls.Fr, b *bls.Fr) {
// 	// TODO: use divmod instead.
// 	var tmp bls.Fr
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
func GetAllPrecomputedSrsMap(tableDir string) ([]rs.EncodingParams, error) {
	files, err := os.ReadDir(tableDir)
	if err != nil {
		log.Println("Error to list SRS Table directory", err)
		return nil, err
	}

	tables := make([]rs.EncodingParams, 0)
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

		params := rs.EncodingParams{
			NumChunks: uint64(cosetSizeValue),
			ChunkLen:  uint64(dimEValue),
		}
		tables = append(tables, params)
	}
	return tables, nil
}
