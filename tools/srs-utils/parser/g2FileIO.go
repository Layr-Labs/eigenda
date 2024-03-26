package parser

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fp"
)

func parseG2Point(xA0Bytes, xA1Bytes, yA0Bytes, yA1Bytes []byte) bn254.G2Affine {
	var xA0, xA1 fp.Element
	var yA0, yA1 fp.Element

	xA0.SetBytes(xA0Bytes[:])
	xA1.SetBytes(xA1Bytes[:])
	yA0.SetBytes(yA0Bytes[:])
	yA1.SetBytes(yA1Bytes[:])

	g2Aff := bn254.G2Affine{}
	g2Aff.X.A0 = xA0
	g2Aff.X.A1 = xA1
	g2Aff.Y.A0 = yA0
	g2Aff.Y.A1 = yA1

	if !g2Aff.IsOnCurve() {
		panic("g2Affine is not on curve")
	}
	return g2Aff
}

func readG2Worker(
	buf []byte,
	outs []bn254.G2Affine,
	start uint64, // in element, not in byte
	end uint64,
	step uint64,
	wg *sync.WaitGroup,
) {
	for i := start; i < end; i++ {
		fieldSize := uint64(32)
		xA1 := buf[i*step : i*step+fieldSize]
		xA0 := buf[i*step+fieldSize : i*step+fieldSize*2]
		yA1 := buf[i*step+fieldSize*2 : i*step+fieldSize*3]
		yA0 := buf[i*step+fieldSize*3 : (i+1)*step]

		point := parseG2Point(xA0, xA1, yA0, yA1)
		outs[i] = point
	}
	wg.Done()
}

func ParseG2PointSection(filepath string, params Params, numWorker uint64) ([]bn254.G2Affine, error) {
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

	n := params.NumPoint
	startTimer := time.Now()
	g1r := bufio.NewReaderSize(g1f, int(params.NumPoint*params.G2Size))

	fmt.Println("params.G2StartByte", params.G2StartByte)
	_, err = g1f.Seek(int64(params.G2StartByte), 0)
	if err != nil {
		return nil, err
	}

	if n < numWorker {
		numWorker = n
	}

	numToRead := params.NumPoint * params.G2Size
	buf := make([]byte, numToRead)
	numBytes, err := g1r.Read(buf)

	if err != nil {
		return nil, err
	}

	if uint64(numBytes) != numToRead {
		log.Printf("Error. Insufficient G2 points. Only contains %v. Requesting %v\n", len(buf)/128, params.NumPoint)
		log.Println("numBytes", numBytes, "numToRead", numToRead)
		log.Println("ReadG2PointSection.ERR.1", err)
		return nil, err
	}

	// measure reading time
	t := time.Now()
	elapsed := t.Sub(startTimer)
	log.Printf("    Reading G2 points (%v bytes) takes %v\n", (n * 128), elapsed)
	startTimer = time.Now()

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
		go readG2Worker(buf, s2Outs, start, end, 128, &wg)
	}
	wg.Wait()

	// measure parsing time
	t = time.Now()
	elapsed = t.Sub(startTimer)
	log.Println("    Parsing takes", elapsed)
	return s2Outs, nil
}

func WriteG2PointsForEigenDA(points []bn254.G2Affine, from uint64, to uint64) error {
	n := to - from
	g2f, err := os.OpenFile("g2.point", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		fmt.Printf("Canot write G1 from %v to %v . Error %v\n", from, to, err)
		return err
	}

	g2w := bufio.NewWriter(g2f)

	for i := uint64(0); i < n; i++ {
		pointInBytes := points[i].Bytes()
		numWritten, err := g2w.Write(pointInBytes[:])
		if numWritten != 64 || err != nil {
			fmt.Printf("Cannot write point %v . Error %v\n", from+i, err)
			return err
		}
	}

	if err = g2w.Flush(); err != nil {
		log.Println("Cannot flush points", err)
		return err
	}
	g2f.Close()
	return nil

}
