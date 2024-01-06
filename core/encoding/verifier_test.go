package encoding_test

import (
	"crypto/rand"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/gammazero/workerpool"
	"github.com/stretchr/testify/assert"
	// "github.com/pkg/profile"
)

// var control interface{ Stop() }

func TestBenchmarkVerifyChunks(t *testing.T) {
	t.Skip("This test is meant to be run manually, not as part of the test suite")

	chunkLengths := []int{64, 128, 256, 512, 1024, 2048, 4096, 8192}
	chunkCounts := []int{4, 8, 16}

	file, err := os.Create("benchmark_results.csv")
	if err != nil {
		t.Fatalf("Failed to open file for writing: %v", err)
	}
	defer file.Close()

	fmt.Fprintln(file, "numChunks,chunkLength,ns/op,allocs/op")

	for _, chunkLength := range chunkLengths {

		blobSize := chunkLength * 31 * 2
		params := core.EncodingParams{
			ChunkLength: uint(chunkLength),
			NumChunks:   16,
		}
		blob := make([]byte, blobSize)
		_, err = rand.Read(blob)
		assert.NoError(t, err)

		commitments, chunks, err := enc.Encode(blob, params)
		assert.NoError(t, err)

		indices := make([]core.ChunkNumber, params.NumChunks)
		for i := range indices {
			indices[i] = core.ChunkNumber(i)
		}

		for _, numChunks := range chunkCounts {

			result := testing.Benchmark(func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					// control = profile.Start(profile.ProfilePath("."))
					err := enc.VerifyChunks(chunks[:numChunks], indices[:numChunks], commitments, params)
					assert.NoError(t, err)
					// control.Stop()
				}
			})
			// Print results in CSV format
			fmt.Fprintf(file, "%d,%d,%d,%d\n", numChunks, chunkLength, result.NsPerOp(), result.AllocsPerOp())

		}
	}

}

func BenchmarkVerifyBlob(b *testing.B) {

	params := core.EncodingParams{
		ChunkLength: uint(256),
		NumChunks:   uint(8),
	}
	blobSize := 8 * 256
	numSamples := 30
	blobs := make([][]byte, numSamples)
	for i := 0; i < numSamples; i++ {
		blob := make([]byte, blobSize)
		_, _ = rand.Read(blob)
		blobs[i] = blob
	}

	commitments, _, err := enc.Encode(blobs[0], params)
	assert.NoError(b, err)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err = enc.VerifyBlobLength(commitments)
		assert.NoError(b, err)
	}

}

func TestBenchmarkBatchVerifyChunks(t *testing.T) {
	t.Skip("This test is meant to be run manually, not as part of the test suite")

	chunkLengths := []int{64, 128, 256, 512, 1024, 2048, 4096, 8192}
	chunkCounts := []int{4, 8, 16}
	numBlobs := 5
	// chunkLengths := []int{64}
	// chunkCounts := []int{4}
	// numBlobs := 5

	file, err := os.Create("benchmark_results.csv")
	if err != nil {
		t.Fatalf("Failed to open file for writing: %v", err)
	}
	defer file.Close()

	fmt.Fprintln(file, "numChunks,chunkLength,ns/op,allocs/op")

	for _, chunkLength := range chunkLengths {

		blobSize := chunkLength * 31 * 2
		params := core.EncodingParams{
			ChunkLength: uint(chunkLength),
			NumChunks:   16,
		}

		blobSamples := make([][]core.Sample, numBlobs)

		for i := 0; i < numBlobs; i++ {

			blob := make([]byte, blobSize)
			_, err = rand.Read(blob)
			assert.NoError(t, err)

			commitments, chunks, err := enc.Encode(blob, params)
			assert.NoError(t, err)

			for j := range chunks {

				blobSamples[i] = append(blobSamples[i], core.Sample{
					Commitment:      commitments.Commitment,
					Chunk:           chunks[j],
					AssignmentIndex: core.ChunkNumber(j),
					BlobIndex:       i,
				})

			}
		}

		for _, numChunks := range chunkCounts {

			result := testing.Benchmark(func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					// control = profile.Start(profile.ProfilePath("."))

					samples := make([]core.Sample, 0, numChunks*numBlobs)
					for j := 0; j < numBlobs; j++ {
						samples = append(samples, blobSamples[j][:numChunks]...)
					}

					err := enc.UniversalVerifySubBatch(params, samples, numBlobs)
					assert.NoError(t, err)
					// control.Stop()
				}
			})
			// Print results in CSV format
			fmt.Fprintf(file, "%d,%d,%d,%d\n", numChunks*numBlobs, chunkLength, result.NsPerOp(), result.AllocsPerOp())

		}
	}

}

// Idea:
// - Specify a stake distribution and then one operator within that distribution.
// - Specify a distribution of blob sizes amounting to X encoded throughput over Y interval
// 		- Random distribution
//      - "Worst case" distribution
// - Generate all of the chunks that this operator would see. It's fine to reuse blobs.
// - Time the validation for the selected validator

type testCase struct {
	ChunkLength int
	NumChunks   int
	NumBlobs    int
	Samples     []core.Sample
	Commitments []core.BlobCommitments
	Params      core.EncodingParams
}

func setupTestCase(t *testing.T, c testCase) testCase {

	blobSize := c.ChunkLength * 31 * 2
	params := core.EncodingParams{
		ChunkLength: uint(c.ChunkLength),
		NumChunks:   uint(2 * c.NumChunks),
	}

	samples := make([]core.Sample, 0, c.NumBlobs*c.NumChunks)
	blobCommitments := make([]core.BlobCommitments, 0, c.NumBlobs)

	blob := make([]byte, blobSize)
	_, err := rand.Read(blob)
	assert.NoError(t, err)

	commitments, chunks, err := enc.Encode(blob, params)
	assert.NoError(t, err)

	for i := 0; i < c.NumBlobs; i++ {

		blobCommitments = append(blobCommitments, commitments)

		for j := 0; j < c.NumChunks; j++ {

			samples = append(samples, core.Sample{
				Commitment:      commitments.Commitment,
				Chunk:           chunks[j],
				AssignmentIndex: core.ChunkNumber(j),
				BlobIndex:       i,
			})

		}
	}

	return testCase{
		ChunkLength: c.ChunkLength,
		NumChunks:   c.NumChunks,
		NumBlobs:    c.NumBlobs,
		Samples:     samples,
		Commitments: blobCommitments,
		Params:      params,
	}

}

var smallOperatorTestCases = []testCase{
	{
		ChunkLength: 1,
		NumChunks:   1,
		NumBlobs:    200,
	},
	{
		ChunkLength: 2,
		NumChunks:   1,
		NumBlobs:    200,
	},
	{
		ChunkLength: 4,
		NumChunks:   1,
		NumBlobs:    200,
	},
	{
		ChunkLength: 8,
		NumChunks:   1,
		NumBlobs:    200,
	},
	{
		ChunkLength: 16,
		NumChunks:   1,
		NumBlobs:    200,
	},
	{
		ChunkLength: 32,
		NumChunks:   1,
		NumBlobs:    200,
	},
	{
		ChunkLength: 64,
		NumChunks:   1,
		NumBlobs:    200,
	},
}

var mediumOperatorTestCases = []testCase{
	{
		ChunkLength: 1,
		NumChunks:   50,
		NumBlobs:    200,
	},
	{
		ChunkLength: 2,
		NumChunks:   50,
		NumBlobs:    200,
	},
	{
		ChunkLength: 4,
		NumChunks:   50,
		NumBlobs:    200,
	},
	{
		ChunkLength: 8,
		NumChunks:   50,
		NumBlobs:    200,
	},
	{
		ChunkLength: 16,
		NumChunks:   50,
		NumBlobs:    200,
	},
	{
		ChunkLength: 32,
		NumChunks:   50,
		NumBlobs:    200,
	},
	{
		ChunkLength: 64,
		NumChunks:   50,
		NumBlobs:    200,
	},
}

var largeOperatorTestCases = []testCase{
	{
		ChunkLength: 1,
		NumChunks:   2420,
		NumBlobs:    200,
	},
	{
		ChunkLength: 2,
		NumChunks:   2420,
		NumBlobs:    200,
	},
	{
		ChunkLength: 4,
		NumChunks:   2420,
		NumBlobs:    200,
	},
	{
		ChunkLength: 8,
		NumChunks:   2420,
		NumBlobs:    200,
	},
	{
		ChunkLength: 16,
		NumChunks:   242,
		NumBlobs:    2000,
	},
	{
		ChunkLength: 32,
		NumChunks:   242,
		NumBlobs:    2000,
	},
	{
		ChunkLength: 64,
		NumChunks:   242,
		NumBlobs:    2000,
	},
}

var altLargeOperatorTestCases = []testCase{
	{
		ChunkLength: 32,
		NumChunks:   8,
		NumBlobs:    2000,
	},
	{
		ChunkLength: 64,
		NumChunks:   8,
		NumBlobs:    2000,
	},
	{
		ChunkLength: 128,
		NumChunks:   8,
		NumBlobs:    2000,
	},
	{
		ChunkLength: 256,
		NumChunks:   8,
		NumBlobs:    2000,
	},
	{
		ChunkLength: 512,
		NumChunks:   8,
		NumBlobs:    2000,
	},
	{
		ChunkLength: 1024,
		NumChunks:   8,
		NumBlobs:    2000,
	},
	{
		ChunkLength: 2048,
		NumChunks:   8,
		NumBlobs:    10000,
	},
}

func testBenchmarkCompositeBatchVerifyChunks(t *testing.T, testCases []testCase, nCPU int) {

	pool := workerpool.New(8)

	casesChan := make(chan testCase, len(testCases))

	for i := range testCases {
		pool.Submit(func() {
			casesChan <- setupTestCase(t, testCases[i])
		})
	}

	pool.StopWait()
	close(casesChan)

	testCases = make([]testCase, 0, len(testCases))
	for c := range casesChan {
		testCases = append(testCases, c)
	}

	// Create worker pool
	pool = workerpool.New(nCPU)

	start := time.Now()

	// Send tasks to worker pool
	for _, c := range testCases {
		c := c
		pool.Submit(func() {
			fmt.Printf("Submitting task for chunkLength=%v, numChunks=%v, numBlobs=%v \n", c.ChunkLength, c.NumChunks, c.NumBlobs)
			assert.EqualValues(t, c.NumChunks*c.NumBlobs, len(c.Samples))
			err := enc.UniversalVerifySubBatch(c.Params, c.Samples, c.NumBlobs)

			for _, commitments := range c.Commitments {
				err := enc.VerifyBlobLength(commitments)
				assert.NoError(t, err)
			}

			assert.NoError(t, err)
			fmt.Printf("Completed task for chunkLength=%v, numChunks=%v, numBlobs=%v \n", c.ChunkLength, c.NumChunks, c.NumBlobs)
		})
	}

	// Shutdown worker pool and wait for all tasks to complete
	pool.StopWait()

	elapsed := time.Since(start)
	fmt.Printf("UniversalVerifySubBatch took %s \n", elapsed)

}

func TestBenchmarkCompositeBatchVerifyChunks(t *testing.T) {
	t.Skip("This test is meant to be run manually, not as part of the test suite")

	testBenchmarkCompositeBatchVerifyChunks(t, smallOperatorTestCases, 2)

	testBenchmarkCompositeBatchVerifyChunks(t, mediumOperatorTestCases, 2)

	testBenchmarkCompositeBatchVerifyChunks(t, largeOperatorTestCases, 8)

	testBenchmarkCompositeBatchVerifyChunks(t, altLargeOperatorTestCases, 8)

}
