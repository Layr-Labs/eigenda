package kzg

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"time"

	"github.com/consensys/gnark-crypto/ecc/bn254"
)

const (
	// Num of bytes per G1 point in serialized format in file.
	G1PointBytes = 32
	// Num of bytes per G2 point in serialized format in file.
	G2PointBytes = 64
)

type EncodeParams struct {
	NumNodeE  uint64
	ChunkLenE uint64
}

// ReadDesiredBytes reads exactly numBytesToRead bytes from the reader and returns
// the result.
func ReadDesiredBytes(reader *bufio.Reader, numBytesToRead uint64) ([]byte, error) {
	buf := make([]byte, numBytesToRead)
	_, err := io.ReadFull(reader, buf)
	// Note that ReadFull() guarantees the bytes read is len(buf) IFF err is nil.
	// See https://pkg.go.dev/io#ReadFull.
	if err != nil {
		return nil, fmt.Errorf("cannot read file %w", err)
	}
	return buf, nil
}

// Read the n-th G1 point from SRS.
func ReadG1Point(n uint64, srsOrder uint64, g1Path string) (bn254.G1Affine, error) {
	if n >= srsOrder {
		return bn254.G1Affine{}, fmt.Errorf("requested power %v is larger than SRSOrder %v", n, srsOrder)
	}

	g1point, err := ReadG1PointSection(g1Path, n, n+1, 1)
	if err != nil {
		return bn254.G1Affine{}, fmt.Errorf("error read g1 point section %w", err)
	}

	return g1point[0], nil
}

// Read the n-th G2 point from SRS.
func ReadG2Point(n uint64, srsOrder uint64, g2Path string) (bn254.G2Affine, error) {
	if n >= srsOrder {
		return bn254.G2Affine{}, fmt.Errorf("requested power %v is larger than SRSOrder %v", n, srsOrder)
	}

	g2point, err := ReadG2PointSection(g2Path, n, n+1, 1)
	if err != nil {
		return bn254.G2Affine{}, fmt.Errorf("error read g2 point section %w", err)
	}
	return g2point[0], nil
}

// Read g2 points from power of 2 file
func ReadG2PointOnPowerOf2(exponent uint64, srsOrder uint64, g2PowerOf2Path string) (bn254.G2Affine, error) {

	// the powerOf2 file, only [tau^exp] are stored.
	// exponent    0,    1,       2,    , ..
	// actual pow [tau],[tau^2],[tau^4],.. (stored in the file)
	// In our convention SRSOrder contains the total number of series of g1, g2 starting with generator
	// i.e. [1] [tau] [tau^2]..
	// So the actual power of tau is SRSOrder - 1
	// The mainnet SRS, the max power is 2^28-1, so the last power in powerOf2 file is [tau^(2^27)]
	// For test case of 3000 SRS, the max power is 2999, so last power in powerOf2 file is [tau^2048]
	// if a actual SRS order is 15, the file will contain four symbols (1,2,4,8) with indices [0,1,2,3]
	// if a actual SRS order is 16, the file will contain five symbols (1,2,4,8,16) with indices [0,1,2,3,4]

	actualPowerOfTau := srsOrder - 1
	largestPowerofSRS := uint64(math.Log2(float64(actualPowerOfTau)))
	if exponent > largestPowerofSRS {
		return bn254.G2Affine{}, fmt.Errorf("requested power %v is larger than largest power of SRS %v",
			uint64(math.Pow(2, float64(exponent))), largestPowerofSRS)
	}

	if len(g2PowerOf2Path) == 0 {
		return bn254.G2Affine{}, errors.New("G2PathPowerOf2 path is empty")
	}

	g2point, err := ReadG2PointSection(g2PowerOf2Path, exponent, exponent+1, 1)
	if err != nil {
		return bn254.G2Affine{}, fmt.Errorf("error read g2 point on power of 2 %w", err)
	}
	return g2point[0], nil
}

func ReadG1Points(filepath string, n uint64, numWorker uint64) ([]bn254.G1Affine, error) {
	g1f, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("error cannot open g1 points file %s: %w", filepath, err)
	}

	defer func() {
		if err := g1f.Close(); err != nil {
			log.Printf("G1 close error %v\n", err)
		}
	}()

	startTimer := time.Now()
	g1r := bufio.NewReaderSize(g1f, int(n*G1PointBytes))

	if n < numWorker {
		numWorker = n
	}

	buf, err := ReadDesiredBytes(g1r, n*G1PointBytes)
	if err != nil {
		return nil, err
	}
	// measure reading time
	t := time.Now()
	elapsed := t.Sub(startTimer)
	log.Printf("    Reading G1 points (%v bytes) takes %v\n", (n * G1PointBytes), elapsed)
	startTimer = time.Now()

	s1Outs := make([]bn254.G1Affine, n)

	start := uint64(0)
	end := uint64(0)
	size := n / numWorker

	results := make(chan error, numWorker)

	for i := uint64(0); i < numWorker; i++ {
		start = i * size

		if i == numWorker-1 {
			end = n
		} else {
			end = (i + 1) * size
		}
		//fmt.Printf("worker %v start %v end %v. size %v\n", i, start, end, end - start)

		go readG1Worker(buf, s1Outs, start, end, G1PointBytes, results)
	}

	for w := uint64(0); w < numWorker; w++ {
		err := <-results
		if err != nil {
			return nil, err
		}
	}

	// measure parsing time
	t = time.Now()
	elapsed = t.Sub(startTimer)
	log.Println("    Parsing takes", elapsed)

	return s1Outs, nil
}

// from is inclusive, to is exclusive
func ReadG1PointSection(filepath string, from, to uint64, numWorker uint64) ([]bn254.G1Affine, error) {
	if to <= from {
		return nil, fmt.Errorf("the range to read is invalid, from: %v, to: %v", from, to)
	}
	g1f, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("error cannot open g1 points file %w", err)
	}

	defer func() {
		if err := g1f.Close(); err != nil {
			log.Printf("g1 close error %v\n", err)
		}
	}()

	n := to - from

	g1r := bufio.NewReaderSize(g1f, int(n*G1PointBytes))

	_, err = g1f.Seek(int64(from)*G1PointBytes, 0)
	if err != nil {
		return nil, err
	}

	if n < numWorker {
		numWorker = n
	}

	buf, err := ReadDesiredBytes(g1r, n*G1PointBytes)
	if err != nil {
		return nil, err
	}

	s1Outs := make([]bn254.G1Affine, n)

	start := uint64(0)
	end := uint64(0)
	size := n / numWorker

	results := make(chan error, numWorker)

	for i := uint64(0); i < numWorker; i++ {
		start = i * size

		if i == numWorker-1 {
			end = n
		} else {
			end = (i + 1) * size
		}

		go readG1Worker(buf, s1Outs, start, end, G1PointBytes, results)
	}

	for w := uint64(0); w < numWorker; w++ {
		err := <-results
		if err != nil {
			return nil, err
		}
	}

	return s1Outs, nil
}

func readG1Worker(
	buf []byte,
	outs []bn254.G1Affine,
	start uint64, // in element, not in byte
	end uint64,
	step uint64,
	results chan<- error,
) {
	for i := start; i < end; i++ {
		g1 := buf[i*step : (i+1)*step]
		_, err := outs[i].SetBytes(g1[:])
		if err != nil {
			results <- err
			return
		}
	}
	results <- nil
}

func readG2Worker(
	buf []byte,
	outs []bn254.G2Affine,
	start uint64, // in element, not in byte
	end uint64,
	step uint64,
	results chan<- error,
) {
	for i := start; i < end; i++ {
		g1 := buf[i*step : (i+1)*step]
		_, err := outs[i].SetBytes(g1[:])
		if err != nil {
			results <- err
			log.Println("Unmarshalling error:", err)
			return
		}
	}
	results <- nil
}

func ReadG2Points(filepath string, n uint64, numWorker uint64) ([]bn254.G2Affine, error) {
	g1f, err := os.Open(filepath)
	if err != nil {
		log.Println("Cannot ReadG2Points", filepath)
		log.Println("ReadG2Points.ERR.0", err)
		return nil, err
	}

	defer func() {
		if err := g1f.Close(); err != nil {
			log.Printf("g2 close error %v\n", err)
		}
	}()

	startTimer := time.Now()
	g1r := bufio.NewReaderSize(g1f, int(n*G2PointBytes))

	if n < numWorker {
		numWorker = n
	}

	buf, err := ReadDesiredBytes(g1r, n*G2PointBytes)
	if err != nil {
		return nil, err
	}

	// measure reading time
	t := time.Now()
	elapsed := t.Sub(startTimer)
	log.Printf("    Reading G2 points (%v bytes) takes %v\n", (n * G2PointBytes), elapsed)

	startTimer = time.Now()

	s2Outs := make([]bn254.G2Affine, n)

	results := make(chan error, numWorker)

	start := uint64(0)
	end := uint64(0)
	size := n / numWorker
	for i := uint64(0); i < numWorker; i++ {
		start = i * size

		if i == numWorker-1 {
			end = n
		} else {
			end = (i + 1) * size
		}

		go readG2Worker(buf, s2Outs, start, end, G2PointBytes, results)
	}

	for w := uint64(0); w < numWorker; w++ {
		err := <-results
		if err != nil {
			return nil, err
		}
	}

	// measure parsing time
	t = time.Now()
	elapsed = t.Sub(startTimer)
	log.Println("    Parsing takes", elapsed)

	return s2Outs, nil
}

// from is inclusive, to is exclusive
func ReadG2PointSection(filepath string, from, to uint64, numWorker uint64) ([]bn254.G2Affine, error) {
	if to <= from {
		return nil, fmt.Errorf("The range to read is invalid, from: %v, to: %v", from, to)
	}
	g2f, err := os.Open(filepath)
	if err != nil {
		log.Println("ReadG2PointSection.ERR.0", err)
		return nil, err
	}

	defer func() {
		if err := g2f.Close(); err != nil {
			log.Printf("error %v\n", err)
		}
	}()

	n := to - from

	g2r := bufio.NewReaderSize(g2f, int(n*G2PointBytes))

	_, err = g2f.Seek(int64(from*G2PointBytes), 0)
	if err != nil {
		return nil, err
	}

	if n < numWorker {
		numWorker = n
	}

	buf, err := ReadDesiredBytes(g2r, n*G2PointBytes)
	if err != nil {
		return nil, err
	}

	s2Outs := make([]bn254.G2Affine, n)

	results := make(chan error, numWorker)

	start := uint64(0)
	end := uint64(0)
	size := n / numWorker

	for i := uint64(0); i < numWorker; i++ {
		start = i * size

		if i == numWorker-1 {
			end = n
		} else {
			end = (i + 1) * size
		}
		//todo: handle error?
		go readG2Worker(buf, s2Outs, start, end, G2PointBytes, results)
	}
	for w := uint64(0); w < numWorker; w++ {
		err := <-results
		if err != nil {
			return nil, err
		}
	}

	return s2Outs, nil
}
