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

func ParseG1PointSection(filepath string, params Params, numWorker uint64) ([]bn254.G1Affine, error) {
	fmt.Printf("Start to read %v points from Byte pos at %v to at %v",
		params.NumPoint,
		params.G1StartByte,
		params.GetG1EndBytePos(),
	)

	g1f, err := os.Open(filepath)
	if err != nil {
		log.Println("ReadG1PointSection.ERR.0", err)
		return nil, err
	}

	defer func() {
		if err := g1f.Close(); err != nil {
			panic(err)
		}
	}()

	n := params.NumPoint
	startTimer := time.Now()
	g1r := bufio.NewReaderSize(g1f, int(params.NumPoint*params.G1Size))

	_, err = g1f.Seek(int64(params.G1StartByte), 0)
	if err != nil {
		return nil, err
	}

	if n < numWorker {
		numWorker = n
	}

	numToRead := params.NumPoint * params.G1Size
	buf := make([]byte, numToRead)
	numBytes, err := g1r.Read(buf)

	if err != nil {
		return nil, err
	}

	if uint64(numBytes) != numToRead {
		log.Printf("Error. Insufficient G1 points. Only contains %v. Requesting %v, NumByte %v\n", len(buf)/64, params.NumPoint, numBytes)
		log.Println("numBytes", numBytes, "numToRead", numToRead)
		log.Println("ReadG1PointSection.ERR.1", err)
		return nil, err
	}

	// measure reading time
	t := time.Now()
	elapsed := t.Sub(startTimer)
	log.Printf("    Reading G1 points (%v bytes) takes %v\n", (n * 64), elapsed)
	startTimer = time.Now()

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
		go readG1Worker(buf, s1Outs, start, end, 64, &wg)
	}
	wg.Wait()

	t = time.Now()
	elapsed = t.Sub(startTimer)
	fmt.Println("Finish Parsing takes", elapsed)
	return s1Outs, nil
}

func WriteG1PointsForEigenDA(points []bn254.G1Affine, from uint64, to uint64) error {
	n := to - from
	g1f, err := os.OpenFile("g1.point", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Canot write G1 from %v to %v . Error %v\n", from, to, err)
		return err
	}

	g1w := bufio.NewWriter(g1f)

	for i := uint64(0); i < n; i++ {
		pointInBytes := points[i].Bytes()
		numWritten, err := g1w.Write(pointInBytes[:])
		if numWritten != 32 || err != nil {
			fmt.Printf("Cannot write point %v . Error %v\n", from+i, err)
			return err
		}
	}

	if err = g1w.Flush(); err != nil {
		log.Println("Cannot flush points", err)
		return err
	}
	g1f.Close()
	return nil

}

func readG1Worker(
	buf []byte,
	outs []bn254.G1Affine,
	start uint64, // in element, not in byte
	end uint64,
	step uint64,
	wg *sync.WaitGroup,
) {
	for i := start; i < end; i++ {
		fieldSize := step / uint64(2)
		g1x := buf[i*step : (i)*step+fieldSize]
		g1y := buf[i*step+fieldSize : (i+1)*step]

		point := parseG1Point(g1x, g1y)
		outs[i] = *point
	}
	wg.Done()
}

func parseG1Point(xBytes, yBytes []byte) *bn254.G1Affine {
	var x fp.Element
	var y fp.Element

	x.SetBytes(xBytes[:])
	y.SetBytes(yBytes[:])

	g1Aff := bn254.G1Affine{}
	g1Aff.X = x
	g1Aff.Y = y

	if !g1Aff.IsOnCurve() {
		panic("g1Affine is not on curve")
	}
	return &g1Aff
}
