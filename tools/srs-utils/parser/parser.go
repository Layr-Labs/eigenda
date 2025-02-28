package parser

import (
	"fmt"
	"math"
	"os"
	"time"
)

type Config struct {
	PtauPath  string
	NumBatch  uint64
	NumPoint  uint64
	NumWorker int
}

// format https://github.com/iden3/snarkjs/blob/master/src/powersoftau_new.js
const (
	totalPoints     = uint64(268435456) // 2^28, starting from generator
	numTotalG1Point = totalPoints*2 - 1
	g1Size          = uint64(64)
	g2Size          = uint64(128)
	OffsetToG1      = uint64(64)
)

func ParsePtauChallenge(config Config) {
	numPoint := config.NumPoint
	numBatch := config.NumBatch

	batchSize := uint64(math.Ceil(float64(numPoint) / float64(numBatch)))

	// Truncate file at beginning
	g1f, err := os.Create("g1.point")
	if err != nil {
		panic(err)
	}
	g1f.Close()
	g2f, err := os.Create("g2.point")
	if err != nil {
		panic(err)
	}
	g2f.Close()

	begin := time.Now()
	for i := uint64(0); i < numBatch; i++ {
		batchBegin := time.Now()
		from := i * batchSize
		to := (i + 1) * batchSize
		if to > numPoint {
			to = numPoint
		}

		fmt.Println("to", to, numPoint)
		actualPoint := to - from
		fmt.Println("actual points", actualPoint)
		p := Params{
			NumPoint:         actualPoint,
			NumTotalG1Points: numTotalG1Point,
			G1Size:           g1Size,
			G2Size:           g2Size,
		}
		p.SetG1StartBytePos(from)
		p.SetG2StartBytePos(from)

		g1Points, err := ParseG1PointSection(config.PtauPath, p, 1)
		if err != nil {
			fmt.Println("main err", err)
		}

		err = WriteG1PointsForEigenDA(g1Points, from, to)
		if err != nil {
			fmt.Println("main err", err)
		}

		g2Points, err := ParseG2PointSection(config.PtauPath, p, 1)
		if err != nil {
			fmt.Println("main err", err)
		}
		err = WriteG2PointsForEigenDA(g2Points, from, to)
		if err != nil {
			fmt.Println("main err", err)
		}
		fmt.Printf("Batch %v takes %v\n", i, time.Since(batchBegin))
	}

	fmt.Println("entire parsing take", time.Since(begin))
}
