package bench_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/kzg"

	"github.com/Layr-Labs/eigenda/encoding/fft"
	"github.com/Layr-Labs/eigenda/test/random"
)

// This file contains benchmarks for the primitives that we use throughout the codebase.
// Higher level benchmarks for the different EigenDA operations can be found in benchmark_eigenda_test.go.
// Speeding up any of the primitives in this file should lead to speedups in the higher level operations.

// We use FFT in many places:
// 1. RS encoding to generate chunks. Max size of 8*blobLen = 8*16MiB = 128MiB = 2^22 Frs
// 2. Per chunk IFFT to generate chunks. Max size of chunkLen = 8*BlobLen/numChunks = 8*16MiB/8KiB = 16KiB = 2^9 Frs
// 3. KZG multiproof to generate chunk proofs. Max size of 2*numChunks = 2*8192 = 2^14 Frs
// 4. Client side when converting encoded_payloads to blobs. Max size of blobLen = 16MiB = 2^19 Frs
func BenchmarkFFTFr(b *testing.B) {
	for _, numFrsPowerOf2 := range []uint8{9, 14, 19, 22} {
		b.Run(fmt.Sprintf("2^%d_elements", numFrsPowerOf2), func(b *testing.B) {
			fs := fft.NewFFTSettings(numFrsPowerOf2)
			rand := random.NewTestRandom(1337)
			frs := rand.FrElements(fs.MaxWidth)

			for b.Loop() {
				_, err := fs.FFT(frs, false)
				require.NoError(b, err)
			}
		})
	}
}

// We need 2 FFT_G1s when generating KZG multiproofs:
// 1. one in inverse direction of size 2*numChunks = 2*8192 = 2^14 G1 points
// 2. one in forward direction of size numChunks = 8192 = 2^13 G1 points
// Note that we don't need FFT_G2.
func BenchmarkFFTG1(b *testing.B) {
	for _, sizePowOf2 := range []uint8{13, 14} {
		b.Run(fmt.Sprintf("2^%d_Points", sizePowOf2), func(b *testing.B) {
			fs := fft.NewFFTSettings(sizePowOf2)
			rand := random.NewTestRandom(1337)
			g1Points, err := rand.G1Points(fs.MaxWidth)
			require.NoError(b, err)

			for b.Loop() {
				_, err := fs.FFTG1(g1Points, false)
				require.NoError(b, err)
			}
		})
	}
}

func BenchmarkGnarkParallelFFTG1(b *testing.B) {
	for _, sizePowOf2 := range []uint8{13, 14} {
		b.Run(fmt.Sprintf("2^%d_G1Points", sizePowOf2), func(b *testing.B) {
			numPoints := uint64(1) << sizePowOf2
			rand := random.NewTestRandom(1337)
			g1Points, err := rand.G1Points(numPoints)
			require.NoError(b, err)

			for b.Loop() {
				_, err := kzg.ToLagrangeG1(g1Points)
				require.NoError(b, err)
			}
		})
	}
}

// We use G1 MSMs in 2 places:
// 1. KZG commitments. Max size of 16MiB = 2^19 Frs/G1s)
// 2. KZG multiproof generation. Max size of ChunkLen = 8*BlobLen/numChunks = 8*16MiB/8KiB = 16KiB = 2^9 Frs/G1s
func BenchmarkMSMG1(b *testing.B) {
	for _, numG1PointsPowOf2 := range []uint8{12, 15, 19} {
		fs := fft.NewFFTSettings(numG1PointsPowOf2)
		rand := random.NewTestRandom(1337)
		frs := rand.FrElements(fs.MaxWidth)
		g1Points, err := rand.G1Points(fs.MaxWidth)
		require.NoError(b, err)

		b.Run(fmt.Sprintf("2^%d_Points", numG1PointsPowOf2), func(b *testing.B) {
			for b.Loop() {
				_, err := new(bn254.G1Affine).MultiExp(g1Points, frs, ecc.MultiExpConfig{})
				require.NoError(b, err)
			}
		})
	}
}

// We use G2 MSMs in 1 place:
// 1. Length commitment+proof generation. Max size of 2^19 Frs/G2s
func BenchmarkMSMG2(b *testing.B) {
	for _, numG2PointsPowOf2 := range []uint8{12, 15, 19} {
		fs := fft.NewFFTSettings(numG2PointsPowOf2)
		rand := random.NewTestRandom(1337)
		frs := rand.FrElements(fs.MaxWidth)
		g2Points, err := rand.G2Points(fs.MaxWidth)
		require.NoError(b, err)

		b.Run(fmt.Sprintf("2^%d_Points", numG2PointsPowOf2), func(b *testing.B) {
			for b.Loop() {
				_, err := new(bn254.G2Affine).MultiExp(g2Points, frs, ecc.MultiExpConfig{})
				require.NoError(b, err)
			}
		})
	}
}
