package srs

import (
	_ "embed"
	"fmt"
	"runtime"
	"sync"

	"github.com/Layr-Labs/eigenda/encoding/v2/kzg"
	"github.com/consensys/gnark-crypto/ecc/bn254"
)

//go:embed g1.point
var serializedG1Points []byte

//go:embed g2.point
var serializedG2Points []byte

//go:embed g2.trailing.point
var serializedG2TrailingPoints []byte

var (
	// Deserializes embedded G1 SRS points on first call. Safe for concurrent use.
	// Points represent [1], [tau], [tau^2],...,[tau^(n-1)] where n is determined by the embedded file size.
	GetG1SRS = sync.OnceValue(func() []bn254.G1Affine {
		fmt.Println("deserializing embedded g1 srs points...")
		points := make([]bn254.G1Affine, len(serializedG1Points)/kzg.G1PointBytes)
		deserializePoints(serializedG1Points, points, kzg.G1PointBytes)
		return points
	})

	// Deserializes embedded G2 SRS points on first call. Safe for concurrent use.
	// Points represent [1], [tau], [tau^2],...,[tau^(n-1)] where n is determined by the embedded file size.
	GetG2SRS = sync.OnceValue(func() []bn254.G2Affine {
		fmt.Println("deserializing embedded g2 srs points...")
		points := make([]bn254.G2Affine, len(serializedG2Points)/kzg.G2PointBytes)
		deserializePoints(serializedG2Points, points, kzg.G2PointBytes)
		return points
	})

	// Deserializes embedded G2 trailing SRS points on first call. Safe for concurrent use.
	// Points represent [tau^(2^28 - n)], [tau^(2^28 - n +1)],...,[tau^(2^28 -1)],
	// where n is determined by the embedded file size.
	GetG2TrailingSRS = sync.OnceValue(func() []bn254.G2Affine {
		fmt.Println("deserializing embedded g2 srs trailing points...")
		points := make([]bn254.G2Affine, len(serializedG2TrailingPoints)/kzg.G2PointBytes)
		deserializePoints(serializedG2TrailingPoints, points, kzg.G2PointBytes)
		return points
	})
)

// deserializes the serializedPoints into the points slice using multiple goroutines.
func deserializePoints[T bn254.G1Affine | bn254.G2Affine](serializedPoints []byte, points []T, pointSizeBytes uint64) {
	n := len(points)
	numWorkers := runtime.GOMAXPROCS(0)
	results := make(chan error, numWorkers)
	pointsPerWorker := n / numWorkers

	for workerIndex := 0; workerIndex < numWorkers; workerIndex++ {
		startPoint := workerIndex * pointsPerWorker
		endPoint := startPoint + pointsPerWorker
		if workerIndex == numWorkers-1 {
			endPoint = n
		}

		go kzg.DeserializePointsInRange(serializedPoints, points,
			uint64(startPoint), uint64(endPoint), pointSizeBytes, results)
	}

	for w := 0; w < numWorkers; w++ {
		if err := <-results; err != nil {
			panic(err)
		}
	}
}
