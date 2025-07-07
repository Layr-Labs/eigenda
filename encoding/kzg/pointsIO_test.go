package kzg

import (
	"fmt"
	"os"
	"runtime"
	"testing"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/stretchr/testify/require"
)

const (
	G1PointsFilePath         = "../../resources/srs/g1.point"
	G2PointsFilePath         = "../../resources/srs/g2.point"
	G2TrailingPointsFilePath = "../../resources/srs/g2.trailing.point"
)

func TestDeserializePoints(t *testing.T) {
	const testNumPoints = 10000

	// Read G1 points
	g1Points, err := ReadG1Points(G1PointsFilePath, testNumPoints, 1)
	require.NoError(t, err)
	require.Len(t, g1Points, int(testNumPoints))

	// Read G2 points
	g2Points, err := ReadG2Points(G2PointsFilePath, testNumPoints, 1)
	require.NoError(t, err)
	require.Len(t, g2Points, testNumPoints)

	// Read G2 trailing points
	g2TrailingPoints, err := ReadG2Points(G2TrailingPointsFilePath, testNumPoints, 1)
	require.NoError(t, err)
	require.Len(t, g2TrailingPoints, testNumPoints)
}

// Benchmark to test efficacy of parsing G1 and G2 points with different number of goroutines (workers).
func BenchmarkNumWorkers(b *testing.B) {
	workerCounts := []int{1, 2, 4, 8, 16, 32, runtime.GOMAXPROCS(0)}
	const benchNumPoints = 10000

	for _, numWorkers := range workerCounts {
		b.Run(fmt.Sprintf("%d-Workers-G1", numWorkers), func(b *testing.B) {
			for b.Loop() {
				g1Points, err := ReadG1Points(G1PointsFilePath, benchNumPoints, uint64(numWorkers))
				require.NoError(b, err)
				require.Len(b, g1Points, benchNumPoints)
			}
		})
	}

	for _, numWorkers := range workerCounts {
		b.Run(fmt.Sprintf("%d-Workers-G2", numWorkers), func(b *testing.B) {
			for b.Loop() {
				g2Points, err := ReadG2Points(G2PointsFilePath, benchNumPoints, uint64(numWorkers))
				require.NoError(b, err)
				require.Len(b, g2Points, benchNumPoints)
			}
		})
	}
}

// ================== UNCOMPRESSED POINTS FILES  ==================
// We currently store the points in compressed form for smaller file sizes.
// We could store in uncompressed form (double size) for faster binary startup time.
// See https://docs.gnark.consensys.io/HowTo/serialize#compression
// The tests/benchmarks below can be used to compare the performance of reading compressed vs uncompressed points files.
// Results when I ran them on my M1 MacBook Pro were 2x faster parsing at the cost of 2x larger file sizes:
// - G2 points: 32 MiB Compressed (9.5s parsing) vs 64 MiB Uncompressed (4.9s parsing)

const (
	G1PointsUncompressedFilePath         = "../../resources/srs/g1_uncompressed.point"
	G2PointsUncompressedFilePath         = "../../resources/srs/g2_uncompressed.point"
	G2TrailingPointsUncompressedFilePath = "../../resources/srs/g2.trailing_uncompressed.point"
)

// BenchmarkReadG2Points benchmarks the time needed to parse compressed and uncompressed G2 points.
// Reading ~16-64MiB files takes ms so doesn't matter much for the benchmark.
func BenchmarkReadG2PointsCompressedVsUncompressed(b *testing.B) {
	b.Skip("Meant to be run manually, run TestGenerateUncompressedPointFiles first to create uncompressed files")

	numWorkers := uint64(runtime.GOMAXPROCS(0))
	testNumPoints := uint64(16 << 20 / G1PointBytes)

	b.Run("Compressed", func(b *testing.B) {
		for b.Loop() {
			_, err := ReadG2Points(G2PointsFilePath, testNumPoints, numWorkers)
			require.NoError(b, err)
		}
	})

	b.Run("Uncompressed", func(b *testing.B) {
		for b.Loop() {
			_, err := ReadG2PointsUncompressed(G2PointsUncompressedFilePath, testNumPoints, numWorkers)
			require.NoError(b, err)
		}
	})
}

// Used to create the uncompressed points files in the resources/srs directory.
func TestGenerateUncompressedPointFiles(t *testing.T) {
	t.Skip("run manually to create uncompressed srs point files")
	numWorkers := uint64(runtime.GOMAXPROCS(0))

	// 16MiB of compressed G1 points means 16 * 1024 * 1024 / G1PointBytes points
	numPoints := uint64(16 << 20 / G1PointBytes)

	g2Points, err := ReadG2Points(G2PointsFilePath, numPoints, numWorkers)
	require.NoError(t, err)

	err = createUncompressedFile(g2Points, G2PointsUncompressedFilePath)
	require.NoError(t, err)

	g2TrailingPoints, err := ReadG2Points(G2TrailingPointsFilePath, numPoints, numWorkers)
	require.NoError(t, err)
	err = createUncompressedFile(g2TrailingPoints, G2TrailingPointsUncompressedFilePath)
	require.NoError(t, err)

	g1Points, err := ReadG1Points(G1PointsFilePath, numPoints, numWorkers)
	require.NoError(t, err)
	err = createUncompressedFile(g1Points, G1PointsUncompressedFilePath)
	require.NoError(t, err)
}

// TestUncompressedPointsFilesEquivalence tests that the uncompressed points files match the original points
func TestUncompressedPointsFilesEquivalence(t *testing.T) {
	t.Skip("run manually to verify uncompressed points files match original points")
	numWorkers := uint64(runtime.GOMAXPROCS(0))
	numPoints := uint64(16 << 20 / G1PointBytes)

	g2Points, err := ReadG2Points(G2PointsFilePath, numPoints, numWorkers)
	require.NoError(t, err)
	g2PointsUncompressed, err := ReadG2PointsUncompressed(G2PointsUncompressedFilePath, numPoints, numWorkers)
	require.NoError(t, err)

	g2PointsTrailing, err := ReadG2Points(G2TrailingPointsFilePath, numPoints, numWorkers)
	require.NoError(t, err)
	g2PointsTrailingUncompressed, err := ReadG2PointsUncompressed(G2TrailingPointsUncompressedFilePath, numPoints, numWorkers)
	require.NoError(t, err)

	g1Points, err := ReadG1Points(G1PointsFilePath, numPoints, numWorkers)
	require.NoError(t, err)
	g1PointsUncompressed, err := ReadG1PointsUncompressed(G1PointsUncompressedFilePath, numPoints, numWorkers)
	require.NoError(t, err)

	// Verify points are equal
	for i := range numPoints {
		require.Equal(t, g2Points[i], g2PointsUncompressed[i], "G2 point mismatch at index %d", i)
		require.Equal(t, g2PointsTrailing[i], g2PointsTrailingUncompressed[i], "G2 trailing point mismatch at index %d", i)
		require.Equal(t, g1Points[i], g1PointsUncompressed[i], "G1 point mismatch at index %d", i)
	}
}

// createUncompressedFile creates a file with uncompressed G2 points
func createUncompressedFile[T bn254.G1Affine | bn254.G2Affine](points []T, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer core.CloseLogOnError(file, filename, nil)

	for _, point := range points {
		// Uncompressed format using RawBytes
		switch p := any(&point).(type) {
		case *bn254.G1Affine:
			data := p.RawBytes()
			if _, err := file.Write(data[:]); err != nil {
				return err
			}
		case *bn254.G2Affine:
			data := p.RawBytes()
			if _, err := file.Write(data[:]); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unsupported point type: %T", p)
		}
	}
	return nil
}
