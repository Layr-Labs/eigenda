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
	// 2^28 points: [1], [tau], [tau^2],..,[tau^(2^28-1)]
	g1SRS     = make([]bn254.G1Affine, len(serializedG1Points)/kzg.G1PointBytes)
	g1SRSOnce sync.Once
)

var (
	// first 16MiB of G2 points: [1], [tau], [tau^2],..,[tau^(2^18-1)]
	g2SRS     = make([]bn254.G2Affine, len(serializedG2Points)/kzg.G2PointBytes)
	g2SRSOnce sync.Once
)

var (
	// trailing 16MiB of G2 points
	g2TrailingSRS     = make([]bn254.G2Affine, len(serializedG2TrailingPoints)/kzg.G2PointBytes)
	g2TrailingSRSOnce sync.Once
)

// GetG1SRS returns the embedded G1 SRS points, deserializing them on the first call.
// Subsequent calls return the already deserialized points.
// This function is safe for concurrent use.
func GetG1SRS() []bn254.G1Affine {
	g1SRSOnce.Do(func() {
		fmt.Println("deserializing embedded g1 srs points...")
		deserializePoints(serializedG1Points, g1SRS, kzg.G1PointBytes)
	})
	return g1SRS
}

// GetG2SRS returns the embedded G2 SRS points, deserializing them on the first call.
// Subsequent calls return the already deserialized points.
// This function is safe for concurrent use.
func GetG2SRS() []bn254.G2Affine {
	g2SRSOnce.Do(func() {
		fmt.Println("deserializing embedded g2 srs points...")
		deserializePoints(serializedG2Points, g2SRS, kzg.G2PointBytes)
	})
	return g2SRS
}

// GetG2TrailingSRS returns the embedded trailing G2 SRS points, deserializing them on the first call.
// Subsequent calls return the already deserialized points.
// This function is safe for concurrent use.
func GetG2TrailingSRS() []bn254.G2Affine {
	g2TrailingSRSOnce.Do(func() {
		fmt.Println("deserializing embedded g2 srs trailing points...")
		deserializePoints(serializedG2TrailingPoints, g2TrailingSRS, kzg.G2PointBytes)
	})
	return g2TrailingSRS
}

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
