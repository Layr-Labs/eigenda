package verifier_test

import (
	"crypto/rand"
	"fmt"
	"os"
	"testing"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/kzg/committer"
	"github.com/Layr-Labs/eigenda/encoding/kzg/prover/v2"
	"github.com/Layr-Labs/eigenda/encoding/kzg/verifier/v2"
	"github.com/stretchr/testify/require"
)

func TestBenchmarkVerifyChunks(t *testing.T) {
	t.Skip("This test is meant to be run manually, not as part of the test suite")

	harness := getTestHarness()

	p, err := prover.NewProver(harness.proverV2KzgConfig, nil)
	require.NoError(t, err)

	committer, err := committer.NewFromConfig(*harness.committerConfig)
	require.Nil(t, err)

	v, err := verifier.NewVerifier(harness.verifierV2KzgConfig, nil)
	require.NoError(t, err)

	chunkLengths := []uint64{64, 128, 256, 512, 1024, 2048, 4096, 8192}
	chunkCounts := []int{4, 8, 16}

	file, err := os.Create("benchmark_results.csv")
	if err != nil {
		t.Fatalf("Failed to open file for writing: %v", err)
	}
	defer core.CloseLogOnError(file, file.Name(), nil)

	_, _ = fmt.Fprintln(file, "numChunks,chunkLength,ns/op,allocs/op")

	for _, chunkLength := range chunkLengths {

		blobSize := chunkLength * 32 * 2
		params := encoding.EncodingParams{
			ChunkLength: chunkLength,
			NumChunks:   16,
		}
		blob := make([]byte, blobSize)
		_, err = rand.Read(blob)
		require.NoError(t, err)

		commitments, err := committer.GetCommitmentsForPaddedLength(blob)
		require.NoError(t, err)
		frames, err := p.GetFrames(blob, params)
		require.NoError(t, err)

		indices := make([]encoding.ChunkNumber, params.NumChunks)
		for i := range indices {
			indices[i] = encoding.ChunkNumber(i)
		}

		for _, numChunks := range chunkCounts {

			result := testing.Benchmark(func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					// control = profile.Start(profile.ProfilePath("."))
					err := v.VerifyFrames(frames[:numChunks], indices[:numChunks], commitments, params)
					require.NoError(t, err)
					// control.Stop()
				}
			})
			// Print results in CSV format
			_, _ = fmt.Fprintf(file, "%d,%d,%d,%d\n", numChunks, chunkLength, result.NsPerOp(), result.AllocsPerOp())

		}
	}

}

func BenchmarkVerifyBlob(b *testing.B) {
	harness := getTestHarness()

	committer, err := committer.NewFromConfig(*harness.committerConfig)
	require.NoError(b, err)

	v, err := verifier.NewVerifier(harness.verifierV2KzgConfig, nil)
	require.NoError(b, err)

	blobSize := 8 * 256
	numSamples := 30
	blobs := make([][]byte, numSamples)
	for i := 0; i < numSamples; i++ {
		blob := make([]byte, blobSize)
		_, _ = rand.Read(blob)
		blobs[i] = blob
	}

	commitments, err := committer.GetCommitmentsForPaddedLength(blobs[0])
	require.NoError(b, err)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err = v.VerifyBlobLength(commitments)
		require.NoError(b, err)
	}

}
