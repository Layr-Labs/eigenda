package utils

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	bls "github.com/Layr-Labs/eigenda/pkg/kzg/bn254"
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
		return nil, err
	}
	return buf, nil
}

func ReadG1Points(filepath string, n uint64, numWorker uint64) ([]bls.G1Point, error) {
	g1f, err := os.Open(filepath)
	if err != nil {
		log.Println("Cannot ReadG1Points", filepath, err)
		return nil, err
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

	s1Outs := make([]bls.G1Point, n)

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
func ReadG1PointSection(filepath string, from, to uint64, numWorker uint64) ([]bls.G1Point, error) {
	if to <= from {
		return nil, fmt.Errorf("The range to read is invalid, from: %v, to: %v", from, to)
	}
	g1f, err := os.Open(filepath)
	if err != nil {
		log.Println("ReadG1PointSection.ERR.0", err)
		return nil, err
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

	s1Outs := make([]bls.G1Point, n)

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
	outs []bls.G1Point,
	start uint64, // in element, not in byte
	end uint64,
	step uint64,
	results chan<- error,
) {
	for i := start; i < end; i++ {
		g1 := buf[i*step : (i+1)*step]
		err := outs[i].UnmarshalText(g1[:])
		if err != nil {
			results <- err
			panic(err)
		}
	}
	results <- nil
}

func readG2Worker(
	buf []byte,
	outs []bls.G2Point,
	start uint64, // in element, not in byte
	end uint64,
	step uint64,
	results chan<- error,
) {
	for i := start; i < end; i++ {
		g1 := buf[i*step : (i+1)*step]
		err := outs[i].UnmarshalText(g1[:])
		if err != nil {
			results <- err
			log.Println("Unmarshalling error:", err)
			panic(err)
		}
	}
	results <- nil
}

func ReadG2Points(filepath string, n uint64, numWorker uint64) ([]bls.G2Point, error) {
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

	s2Outs := make([]bls.G2Point, n)

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
func ReadG2PointSection(filepath string, from, to uint64, numWorker uint64) ([]bls.G2Point, error) {
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

	startTimer := time.Now()
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

	// measure reading time
	t := time.Now()
	elapsed := t.Sub(startTimer)
	log.Printf("    Reading G2 points (%v bytes) takes %v\n", (n * G2PointBytes), elapsed)
	startTimer = time.Now()

	s2Outs := make([]bls.G2Point, n)

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

	// measure parsing time
	t = time.Now()
	elapsed = t.Sub(startTimer)
	log.Println("    Parsing takes", elapsed)
	return s2Outs, nil
}
