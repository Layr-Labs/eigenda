package verifier

import (
	"bufio"
	"log"
	"os"
	"sync"

	"github.com/consensys/gnark-crypto/ecc/bn254"
)

const G1ByteNum = 32
const G2ByteNum = 64

// from is inclusive, to is exclusive
func ReadG1PointSection(filepath string, from, to uint64, numWorker uint64) ([]bn254.G1Affine, error) {
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

	g1r := bufio.NewReaderSize(g1f, int(n*G1ByteNum))

	_, err = g1f.Seek(int64(from*G1ByteNum), 0)
	if err != nil {
		return nil, err
	}

	if n < numWorker {
		numWorker = n
	}

	buf := make([]byte, n*G1ByteNum)
	readN, err := g1r.Read(buf)
	if err != nil {
		return nil, err
	}

	if uint64(readN) != n*G1ByteNum {
		log.Printf("Error. Insufficient G1 points. Only contains %v. Requesting %v\n", len(buf)/G1ByteNum, n)
		log.Println()
		log.Println("ReadG1PointSection.ERR.1", err)
		return nil, err
	}

	s1Outs := make([]bn254.G1Affine, n)

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
		go readG1WorkeGnark(buf, s1Outs, start, end, G1ByteNum, &wg)
	}
	wg.Wait()

	return s1Outs, nil
}
func readG1WorkeGnark(
	buf []byte,
	outs []bn254.G1Affine,
	start uint64, // in element, not in byte
	end uint64,
	step uint64,
	wg *sync.WaitGroup,
) {
	for i := start; i < end; i++ {
		g1 := buf[i*step : (i+1)*step]
		n, err := outs[i].SetBytes(g1[:])
		if err != nil {
			panic(err)
		}
		if n != G1ByteNum {
			panic("cannot read 32 bytes")
		}
	}
	wg.Done()
}

func readG2WorkerGnark(
	buf []byte,
	outs []bn254.G2Affine,
	start uint64, // in element, not in byte
	end uint64,
	step uint64,
	wg *sync.WaitGroup,
) {
	for i := start; i < end; i++ {
		g2 := buf[i*step : (i+1)*step]
		n, err := outs[i].SetBytes(g2[:])
		if err != nil {
			log.Println("Unmarshalling error:", err)
			panic("error")
		}
		if n != G2ByteNum {
			panic("Cannot read 128 bytes")
		}
	}
	wg.Done()
}

func ReadG2PointSection(filepath string, from, to uint64, numWorker uint64) ([]bn254.G2Affine, error) {
	g2f, err := os.Open(filepath)
	if err != nil {
		log.Println("ReadG2PointSection.ERR.0", err)
		return nil, err
	}

	//todo: how to handle?
	defer func() {
		if err := g2f.Close(); err != nil {
			panic(err)
		}
	}()

	n := to - from

	g2r := bufio.NewReaderSize(g2f, int(n*G2ByteNum))

	_, err = g2f.Seek(int64(from*G2ByteNum), 0)
	if err != nil {
		return nil, err
	}

	if n < numWorker {
		numWorker = n
	}

	buf := make([]byte, n*G2ByteNum)
	readN, err := g2r.Read(buf)
	if err != nil {
		return nil, err
	}

	if uint64(readN) != n*G2ByteNum {
		log.Printf("Error. Insufficient G2 points. Only contains %v. Requesting %v\n", len(buf)/G2ByteNum, n)
		log.Println()
		log.Println("ReadG2PointSection.ERR.1", err)
		return nil, err
	}

	s2Outs := make([]bn254.G2Affine, n)

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
		go readG2WorkerGnark(buf, s2Outs, start, end, G2ByteNum, &wg)
	}
	wg.Wait()

	return s2Outs, nil
}
