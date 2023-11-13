package utils

import (
	"bufio"
	"log"
	"os"
	"sync"
	"time"

	bls "github.com/Layr-Labs/eigenda/pkg/kzg/bn254"
)

type EncodeParams struct {
	NumNodeE  uint64
	ChunkLenE uint64
}

func ReadG1Points(filepath string, n uint64, numWorker uint64) ([]bls.G1Point, error) {
	g1f, err := os.Open(filepath)
	if err != nil {
		log.Println("Cannot ReadG1Points", filepath, err)
		return nil, err
	}

	//todo: resolve panic
	defer func() {
		if err := g1f.Close(); err != nil {
			panic(err)
		}
	}()

	startTimer := time.Now()
	g1r := bufio.NewReaderSize(g1f, int(n*64))

	if n < numWorker {
		numWorker = n
	}

	buf, _, err := g1r.ReadLine()
	if err != nil {
		return nil, err
	}

	if uint64(len(buf)) < 64*n {
		log.Printf("Error. Insufficient G1 points. Only contains %v. Requesting %v\n", len(buf)/64, n)
		log.Println()
		log.Println("ReadG1Points.ERR.1", err)
		return nil, err
	}

	// measure reading time
	t := time.Now()
	elapsed := t.Sub(startTimer)
	log.Printf("    Reading G1 points (%v bytes) takes %v\n", (n * 64), elapsed)
	startTimer = time.Now()

	s1Outs := make([]bls.G1Point, n)

	var wg sync.WaitGroup
	wg.Add(int(numWorker))

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
		//fmt.Printf("worker %v start %v end %v. size %v\n", i, start, end, end - start)
		//todo: handle error?
		go readG1Worker(buf, s1Outs, start, end, 64, &wg)
	}
	wg.Wait()

	// measure parsing time
	t = time.Now()
	elapsed = t.Sub(startTimer)
	log.Println("    Parsing takes", elapsed)
	return s1Outs, nil
}

func ReadG1PointSection(filepath string, from, to uint64, numWorker uint64) ([]bls.G1Point, error) {
	g1f, err := os.Open(filepath)
	if err != nil {
		log.Println("ReadG1PointSection.ERR.0", err)
		return nil, err
	}

	//todo: how to handle?
	defer func() {
		if err := g1f.Close(); err != nil {
			panic(err)
		}
	}()

	n := to - from

	startTimer := time.Now()
	g1r := bufio.NewReaderSize(g1f, int(to*64))

	_, err = g1f.Seek(int64(from*64), 0)
	if err != nil {
		return nil, err
	}

	if n < numWorker {
		numWorker = n
	}

	buf, _, err := g1r.ReadLine()
	if err != nil {
		return nil, err
	}

	if uint64(len(buf)) < 64*n {
		log.Printf("Error. Insufficient G1 points. Only contains %v. Requesting %v\n", len(buf)/64, n)
		log.Println()
		log.Println("ReadG1PointSection.ERR.1", err)
		return nil, err
	}

	// measure reading time
	t := time.Now()
	elapsed := t.Sub(startTimer)
	log.Printf("    Reading G1 points (%v bytes) takes %v\n", (n * 64), elapsed)
	startTimer = time.Now()

	s1Outs := make([]bls.G1Point, n)

	var wg sync.WaitGroup
	wg.Add(int(numWorker))

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
		go readG1Worker(buf, s1Outs, start, end, 64, &wg)
	}
	wg.Wait()

	// measure parsing time
	t = time.Now()
	elapsed = t.Sub(startTimer)
	log.Println("    Parsing takes", elapsed)
	return s1Outs, nil
}

func readG1Worker(
	buf []byte,
	outs []bls.G1Point,
	start uint64, // in element, not in byte
	end uint64,
	step uint64,
	wg *sync.WaitGroup,
) {
	for i := start; i < end; i++ {
		g1 := buf[i*step : (i+1)*step]
		err := outs[i].UnmarshalText(g1[:])
		if err != nil {
			panic(err)
		}
	}
	wg.Done()
}

func readG2Worker(
	buf []byte,
	outs []bls.G2Point,
	start uint64, // in element, not in byte
	end uint64,
	step uint64,
	wg *sync.WaitGroup,
) {
	for i := start; i < end; i++ {
		g1 := buf[i*step : (i+1)*step]
		err := outs[i].UnmarshalText(g1[:])
		if err != nil {
			log.Println("Unmarshalling error:", err)
		}
	}
	wg.Done()
}

func ReadG2Points(filepath string, n uint64, numWorker uint64) ([]bls.G2Point, error) {
	g1f, err := os.Open(filepath)
	if err != nil {
		log.Println("Cannot ReadG2Points", filepath)
		log.Println("ReadG2Points.ERR.0", err)
		return nil, err
	}
	//todo: resolve panic
	defer func() {
		if err := g1f.Close(); err != nil {
			panic(err)
		}
	}()

	startTimer := time.Now()
	g1r := bufio.NewReaderSize(g1f, int(n*128))

	if n < numWorker {
		numWorker = n
	}

	buf, _, err := g1r.ReadLine()
	if err != nil {
		return nil, err
	}

	if uint64(len(buf)) < 128*n {
		log.Printf("Error. Insufficient G1 points. Only contains %v. Requesting %v\n", len(buf)/128, n)
		log.Println()
		log.Println("ReadG2Points.ERR.1", err)
		return nil, err
	}

	// measure reading time
	t := time.Now()
	elapsed := t.Sub(startTimer)
	log.Printf("    Reading G2 points (%v bytes) takes %v\n", (n * 128), elapsed)

	startTimer = time.Now()

	s2Outs := make([]bls.G2Point, n)

	var wg sync.WaitGroup
	wg.Add(int(numWorker))

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
		go readG2Worker(buf, s2Outs, start, end, 128, &wg)
	}
	wg.Wait()

	// measure parsing time
	t = time.Now()
	elapsed = t.Sub(startTimer)
	log.Println("    Parsing takes", elapsed)
	return s2Outs, nil
}
