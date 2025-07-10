package kzg

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"os"

	"github.com/consensys/gnark-crypto/ecc/bn254"
)

const (
	// We store the points in compressed form for smaller file sizes.
	// We could store in uncompressed form (double size) for faster binary startup time.
	// See https://docs.gnark.consensys.io/HowTo/serialize#compression
	// and [BenchmarkReadG2PointsCompressedVsUncompressed] for performance comparison.

	// Num of bytes per G1 point in (compressed) serialized format in file.
	G1PointBytes = bn254.SizeOfG1AffineCompressed
	// Num of bytes per G2 point in (compressed) serialized format in file.
	G2PointBytes = bn254.SizeOfG2AffineCompressed
)

// Read the n-th G1 point from SRS.
func ReadG1Point(n uint64, srsOrder uint64, g1Path string) (bn254.G1Affine, error) {
	// TODO: Do we really need to check srsOrder here? Or can we just read the file and let the error propagate if n is out of bounds?
	if n >= srsOrder {
		return bn254.G1Affine{}, fmt.Errorf("requested power %v is larger than SRSOrder %v", n, srsOrder)
	}

	g1point, err := ReadG1PointSection(g1Path, n, n+1, 1)
	if err != nil {
		return bn254.G1Affine{}, fmt.Errorf("error read g1 point section %w", err)
	}

	return g1point[0], nil
}

// Convenience wrapper around [readPointSection] for reading a section of G1 points.
func ReadG1PointSection(filepath string, from, to uint64, numWorker uint64) ([]bn254.G1Affine, error) {
	return readPointSection[bn254.G1Affine](filepath, from, to, G1PointBytes, numWorker)
}

// Convenience wrapper for reading all G1 points from the start of the file.
// n is the number of points to read, numWorker is the number of goroutines to use for parallel parsing.
func ReadG1Points(filepath string, n uint64, numWorker uint64) ([]bn254.G1Affine, error) {
	// ReadG1Points is just ReadG1PointSection starting from 0
	return ReadG1PointSection(filepath, 0, n, numWorker)
}

// Convenience wrapper for reading all G1 points in uncompressed format.
// n is the number of points to read, numWorker is the number of goroutines to use for parallel parsing.
// We don't currently use uncompressed file formats; see [BenchmarkReadG2PointsCompressedVsUncompressed] for performance comparison.
func ReadG1PointsUncompressed(filepath string, n uint64, numWorker uint64) ([]bn254.G1Affine, error) {
	// ReadG1PointsUncompressed is just ReadG1PointSection starting from 0
	result, err := readPointSection[bn254.G1Affine](filepath, 0, n, bn254.SizeOfG1AffineUncompressed, numWorker)
	if err != nil {
		return nil, fmt.Errorf("ReadG1PointsUncompressed: %w", err)
	}

	return result, nil
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

// Convenience wrapper around [readPointSection] for reading G2 points from the start of the file.
// n is the number of points to read, numWorker is the number of goroutines to use for parallel parsing.
func ReadG2Points(filepath string, n uint64, numWorker uint64) ([]bn254.G2Affine, error) {
	result, err := ReadG2PointSection(filepath, 0, n, numWorker)
	if err != nil {
		return nil, fmt.Errorf("ReadG2Points: %w", err)
	}

	return result, nil
}

// Convenience wrapper for reading all G2 points in uncompressed format.
// n is the number of points to read, numWorker is the number of goroutines to use for parallel parsing.
// We don't currently use uncompressed file formats; see [BenchmarkReadG2PointsCompressedVsUncompressed] for performance comparison.
func ReadG2PointsUncompressed(filepath string, n uint64, numWorker uint64) ([]bn254.G2Affine, error) {
	// ReadG2PointsUncompressed is just ReadG2PointSection starting from 0
	result, err := readPointSection[bn254.G2Affine](filepath, 0, n, bn254.SizeOfG2AffineUncompressed, numWorker)
	if err != nil {
		return nil, fmt.Errorf("ReadG2PointsUncompressed: %w", err)
	}

	return result, nil
}

// Convenience wrapper for reading a section of G2 points.
// from and to specify the range of point indices to read (inclusive from, exclusive to).
// numWorker specifies the number of goroutines to use for parallel parsing.
func ReadG2PointSection(filepath string, from, to uint64, numWorker uint64) ([]bn254.G2Affine, error) {
	return readPointSection[bn254.G2Affine](filepath, from, to, G2PointBytes, numWorker)
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

// readPointSection is a generic function for reading a section of points from an SRS file:
//   - `pointsFilePath` is the path to the file containing the points.
//   - `from` and `to` specify the range of point indices to read (inclusive `from`, exclusive `to`).
//   - `pointSizeBytes` is the size of each point in bytes, which can be any of
//     [bn254.SizeOfG1AffineCompressed], [bn254.SizeOfG2AffineCompressed], [bn254.SizeOfG1AffineUncompressed], [bn254.SizeOfG2AffineUncompressed]
//   - `numWorker` specifies the number of goroutines to use for parsing the points in parallel.
func readPointSection[T bn254.G1Affine | bn254.G2Affine](
	pointsFilePath string,
	from, to uint64,
	pointSizeBytes uint64, // TODO: we should probably infer this from the header byte of the first point in the file
	numWorker uint64,
) ([]T, error) {
	if to <= from {
		return nil, fmt.Errorf("to index %v must be greater than from index %v", to, from)
	}
	if numWorker == 0 {
		return nil, fmt.Errorf("numWorker must be greater than 0")
	}

	file, err := os.Open(pointsFilePath)
	if err != nil {
		return nil, fmt.Errorf("error cannot open points file %w", err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("close error %v\n", err)
		}
	}()

	n := to - from
	reader := bufio.NewReaderSize(file, int(n*pointSizeBytes))

	_, err = file.Seek(int64(from)*int64(pointSizeBytes), 0)
	if err != nil {
		return nil, err
	}

	if n < numWorker {
		numWorker = n
	}

	buf, err := readBytes(reader, n*pointSizeBytes)
	if err != nil {
		return nil, fmt.Errorf("readBytes: %w", err)
	}

	points := make([]T, n)
	results := make(chan error, numWorker)
	pointsPerWorker := n / numWorker

	for i := uint64(0); i < numWorker; i++ {
		startPoint := i * pointsPerWorker
		endPoint := startPoint + pointsPerWorker
		if i == numWorker-1 {
			endPoint = n
		}

		go func(startPoint, endPoint uint64) {
			for j := startPoint; j < endPoint; j++ {
				pointData := buf[j*pointSizeBytes : (j+1)*pointSizeBytes]
				switch p := any(&points[j]).(type) {
				case *bn254.G1Affine:
					if _, err := p.SetBytes(pointData); err != nil {
						results <- fmt.Errorf("error setting G1 point bytes: %w", err)
						return
					}
				case *bn254.G2Affine:
					// G2 points are stored as uncompressed points, so we can directly set bytes
					if _, err := p.SetBytes(pointData); err != nil {
						results <- fmt.Errorf("error setting G2 point bytes: %w", err)
						return
					}
				default:
					results <- fmt.Errorf("unsupported point type: %T", p)
					return
				}
			}
			results <- nil
		}(startPoint, endPoint)
	}

	for w := uint64(0); w < numWorker; w++ {
		if err := <-results; err != nil {
			return nil, err
		}
	}

	return points, nil
}

// readBytes reads exactly numBytesToRead bytes from the reader and returns
// the result.
func readBytes(reader *bufio.Reader, numBytesToRead uint64) ([]byte, error) {
	buf := make([]byte, numBytesToRead)
	_, err := io.ReadFull(reader, buf)
	// Note that ReadFull() guarantees the bytes read is len(buf) IFF err is nil.
	if err != nil {
		return nil, fmt.Errorf("cannot read file %w", err)
	}
	return buf, nil
}
