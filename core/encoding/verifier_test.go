package encoding_test

// import (
// 	"crypto/rand"
// 	"fmt"
// 	"os"
// 	"testing"

// 	"github.com/Layr-Labs/eigenda/encoding"
// 	"github.com/stretchr/testify/assert"
// 	// "github.com/pkg/profile"
// )

// // var control interface{ Stop() }

// func TestBenchmarkVerifyChunks(t *testing.T) {
// 	t.Skip("This test is meant to be run manually, not as part of the test suite")

// 	chunkLengths := []int{64, 128, 256, 512, 1024, 2048, 4096, 8192}
// 	chunkCounts := []int{4, 8, 16}

// 	file, err := os.Create("benchmark_results.csv")
// 	if err != nil {
// 		t.Fatalf("Failed to open file for writing: %v", err)
// 	}
// 	defer file.Close()

// 	fmt.Fprintln(file, "numChunks,chunkLength,ns/op,allocs/op")

// 	for _, chunkLength := range chunkLengths {

// 		blobSize := chunkLength * 31 * 2
// 		params := encoding.EncodingParams{
// 			ChunkLength: uint(chunkLength),
// 			NumChunks:   16,
// 		}
// 		blob := make([]byte, blobSize)
// 		_, err = rand.Read(blob)
// 		assert.NoError(t, err)

// 		commitments, chunks, err := enc.Encode(blob, params)
// 		assert.NoError(t, err)

// 		indices := make([]encoding.ChunkNumber, params.NumChunks)
// 		for i := range indices {
// 			indices[i] = encoding.ChunkNumber(i)
// 		}

// 		for _, numChunks := range chunkCounts {

// 			result := testing.Benchmark(func(b *testing.B) {
// 				for i := 0; i < b.N; i++ {
// 					// control = profile.Start(profile.ProfilePath("."))
// 					err := enc.VerifyChunks(chunks[:numChunks], indices[:numChunks], commitments, params)
// 					assert.NoError(t, err)
// 					// control.Stop()
// 				}
// 			})
// 			// Print results in CSV format
// 			fmt.Fprintf(file, "%d,%d,%d,%d\n", numChunks, chunkLength, result.NsPerOp(), result.AllocsPerOp())

// 		}
// 	}

// }

// func BenchmarkVerifyBlob(b *testing.B) {

// 	params := encoding.EncodingParams{
// 		ChunkLength: uint(256),
// 		NumChunks:   uint(8),
// 	}
// 	blobSize := 8 * 256
// 	numSamples := 30
// 	blobs := make([][]byte, numSamples)
// 	for i := 0; i < numSamples; i++ {
// 		blob := make([]byte, blobSize)
// 		_, _ = rand.Read(blob)
// 		blobs[i] = blob
// 	}

// 	commitments, _, err := enc.Encode(blobs[0], params)
// 	assert.NoError(b, err)

// 	b.ResetTimer()

// 	for i := 0; i < b.N; i++ {
// 		err = enc.VerifyBlobLength(commitments)
// 		assert.NoError(b, err)
// 	}

// }
