package kzgEncoder

import (
	"fmt"
	"log"
	"math"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	rs "github.com/Layr-Labs/eigenda/pkg/encoding/encoder"
	"github.com/Layr-Labs/eigenda/pkg/encoding/utils"
	kzg "github.com/Layr-Labs/eigenda/pkg/kzg"

	bls "github.com/Layr-Labs/eigenda/pkg/kzg/bn254"
	_ "go.uber.org/automaxprocs"
)

type KzgConfig struct {
	G1Path         string
	G2Path         string
	CacheDir       string
	NumWorker      uint64
	SRSOrder       uint64 // Order is the total size of SRS
	Verbose        bool
	PreloadEncoder bool
}

type KzgEncoderGroup struct {
	*KzgConfig
	Srs *kzg.SRS
	mu  sync.Mutex

	Encoders  map[rs.EncodingParams]*KzgEncoder
	Verifiers map[rs.EncodingParams]*KzgVerifier
}

type KzgEncoder struct {
	*rs.Encoder

	*KzgConfig
	Srs *kzg.SRS

	Fs         *kzg.FFTSettings
	Ks         *kzg.KZGSettings
	SFs        *kzg.FFTSettings // fft used for submatrix product helper
	FFTPoints  [][]bls.G1Point
	FFTPointsT [][]bls.G1Point // transpose of FFTPoints
}

func NewKzgEncoderGroup(config *KzgConfig) (*KzgEncoderGroup, error) {
	// read the whole order, and treat it as entire SRS for low degree proof
	s1, err := utils.ReadG1Points(config.G1Path, config.SRSOrder, config.NumWorker)
	if err != nil {
		log.Println("failed to read G1 points", err)
		return nil, err
	}
	s2, err := utils.ReadG2Points(config.G2Path, config.SRSOrder, config.NumWorker)
	if err != nil {
		log.Println("failed to read G2 points", err)
		return nil, err
	}

	srs, err := kzg.NewSrs(s1, s2)
	if err != nil {
		log.Println("Could not create srs", err)
		return nil, err
	}

	fmt.Println("numthread", runtime.GOMAXPROCS(0))

	encoderGroup := &KzgEncoderGroup{
		KzgConfig: config,
		Srs:       srs,
		Encoders:  make(map[rs.EncodingParams]*KzgEncoder),
		Verifiers: make(map[rs.EncodingParams]*KzgVerifier),
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

func (g *KzgEncoderGroup) PreloadAllEncoders() error {
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
		g.Encoders[params] = enc
	}

	return nil
}

func (g *KzgEncoderGroup) GetKzgEncoder(params rs.EncodingParams) (*KzgEncoder, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	enc, ok := g.Encoders[params]
	if ok {
		return enc, nil
	}

	enc, err := g.newKzgEncoder(params)
	if err == nil {
		g.Encoders[params] = enc
	}

	return enc, err
}

func (g *KzgEncoderGroup) NewKzgEncoder(params rs.EncodingParams) (*KzgEncoder, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	return g.newKzgEncoder(params)
}

func (g *KzgEncoderGroup) newKzgEncoder(params rs.EncodingParams) (*KzgEncoder, error) {

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

	return &KzgEncoder{
		Encoder:    encoder,
		KzgConfig:  g.KzgConfig,
		Srs:        g.Srs,
		Fs:         fs,
		Ks:         ks,
		SFs:        sfs,
		FFTPoints:  fftPoints,
		FFTPointsT: fftPointsT,
	}, nil
}

// just a wrapper to take bytes not Fr Element
func (g *KzgEncoder) EncodeBytes(inputBytes []byte) (*bls.G1Point, *bls.G1Point, []Frame, []uint32, error) {
	inputFr := rs.ToFrArray(inputBytes)
	return g.Encode(inputFr)
}

func (g *KzgEncoder) Encode(inputFr []bls.Fr) (*bls.G1Point, *bls.G1Point, []Frame, []uint32, error) {

	startTime := time.Now()
	poly, frames, indices, err := g.Encoder.Encode(inputFr)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	// compute commit for the full poly
	commit := g.Commit(poly.Coeffs)

	intermediate := time.Now()

	polyDegreePlus1 := uint64(len(inputFr))

	if g.Verbose {
		log.Printf("    Commiting takes  %v\n", time.Since(intermediate))
		intermediate = time.Now()

		log.Printf("shift %v\n", g.SRSOrder-polyDegreePlus1)
		log.Printf("order %v\n", len(g.Srs.G2))
		log.Println("low degree verification info")
	}

	shiftedSecret := g.Srs.G1[g.SRSOrder-polyDegreePlus1:]

	//The proof of low degree is commitment of the polynomial shifted to the largest srs degree
	lowDegreeProof := bls.LinCombG1(shiftedSecret, poly.Coeffs[:polyDegreePlus1])
	//fmt.Println("kzgFFT lowDegreeProof", lowDegreeProof, "poly len ", len(fullCoeffsPoly), "order", len(g.Ks.SecretG2) )
	//ok := VerifyLowDegreeProof(&commit, lowDegreeProof, polyDegreePlus1-1, g.SRSOrder, g.Srs.G2)
	//if !ok {
	//		log.Printf("Kzg FFT Cannot Verify low degree proof %v", lowDegreeProof)
	//		return nil, nil, nil, nil, errors.New("cannot verify low degree proof")
	//	} else {
	//		log.Printf("Kzg FFT Verify low degree proof  PPPASSS %v", lowDegreeProof)
	//	}

	if g.Verbose {
		log.Printf("    Generating Low Degree Proof takes  %v\n", time.Since(intermediate))
		intermediate = time.Now()
	}

	// compute proofs
	paddedCoeffs := make([]bls.Fr, g.NumEvaluations())
	copy(paddedCoeffs, poly.Coeffs)

	proofs, err := g.ProveAllCosetThreads(paddedCoeffs, g.NumChunks, g.ChunkLen, g.NumWorker)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("could not generate proofs: %v", err)
	}

	if g.Verbose {
		log.Printf("    Proving takes    %v\n", time.Since(intermediate))
	}

	kzgFrames := make([]Frame, len(frames))
	for i, index := range indices {
		kzgFrames[i] = Frame{
			Proof:  proofs[index],
			Coeffs: frames[i].Coeffs,
		}
	}

	if g.Verbose {
		log.Printf("Total encoding took      %v\n", time.Since(startTime))
	}
	return &commit, lowDegreeProof, kzgFrames, indices, nil
}

func (g *KzgEncoder) Commit(polyFr []bls.Fr) bls.G1Point {
	commit := g.Ks.CommitToPoly(polyFr)
	return *commit
}

// The function verify low degree proof against a poly commitment
// We wish to show x^shift poly = shiftedPoly, with
// With shift = SRSOrder-1 - claimedDegree and
// proof = commit(shiftedPoly) on G1
// so we can verify by checking
// e( commit_1, [x^shift]_2) = e( proof_1, G_2 )
func VerifyLowDegreeProof(poly, proof *bls.G1Point, claimedDegree, SRSOrder uint64, srsG2 []bls.G2Point) bool {
	return bls.PairingsVerify(poly, &srsG2[SRSOrder-1-claimedDegree], proof, &bls.GenG2)
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
