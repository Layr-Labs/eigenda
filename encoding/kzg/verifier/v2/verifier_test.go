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
	"github.com/Layr-Labs/eigenda/encoding/rs"

	"github.com/stretchr/testify/require"
)

func TestVerifyFrames(t *testing.T) {
	harness := getTestHarness()

	params := encoding.ParamsFromSysPar(harness.numSys, harness.numPar, uint64(len(harness.paddedGettysburgAddressBytes)))

	proverGroup, err := prover.NewProver(harness.proverV2KzgConfig, nil)
	require.Nil(t, err)

	committer, err := committer.NewFromConfig(*harness.committerConfig)
	require.Nil(t, err)

	frames, err := proverGroup.GetFrames(harness.paddedGettysburgAddressBytes, params)
	require.Nil(t, err)
	commitments, err := committer.GetCommitmentsForPaddedLength(harness.paddedGettysburgAddressBytes)
	require.Nil(t, err)

	verifierGroup, err := verifier.NewVerifier(harness.verifierV2KzgConfig)
	require.Nil(t, err)

	indices := []encoding.ChunkNumber{}
	for i := range len(frames) {
		indices = append(indices, encoding.ChunkNumber(i))
	}
	err = verifierGroup.VerifyFrames(frames, indices, commitments, params)
	require.Nil(t, err)
}

func TestUniversalVerify(t *testing.T) {
	harness := getTestHarness()

	group, err := prover.NewProver(harness.proverV2KzgConfig, nil)
	require.Nil(t, err)

	committer, err := committer.NewFromConfig(*harness.committerConfig)
	require.Nil(t, err)

	v, err := verifier.NewVerifier(harness.verifierV2KzgConfig)
	require.Nil(t, err)

	params := encoding.ParamsFromSysPar(harness.numSys, harness.numPar, uint64(len(harness.paddedGettysburgAddressBytes)))
	blobLength := uint64(encoding.GetBlobLengthPowerOf2(uint32(len(harness.paddedGettysburgAddressBytes))))
	prover, err := group.GetKzgProver(params, blobLength)
	require.Nil(t, err)

	numBlob := 5
	samples := make([]encoding.Sample, 0)
	for z := 0; z < numBlob; z++ {
		inputFr, err := rs.ToFrArray(harness.paddedGettysburgAddressBytes)
		require.Nil(t, err)

		commit, _, _, err := committer.GetCommitments(inputFr)
		require.Nil(t, err)
		frames, fIndices, err := prover.GetFrames(inputFr)
		require.Nil(t, err)

		// create samples
		for i := 0; i < len(frames); i++ {
			f := frames[i]
			j := fIndices[i]

			q, err := rs.GetLeadingCosetIndex(uint64(i), harness.numSys+harness.numPar)
			require.Nil(t, err)

			require.Equal(t, j, q, "leading coset inconsistency")

			sample := encoding.Sample{
				Commitment:      (*encoding.G1Commitment)(commit),
				Chunk:           &f,
				BlobIndex:       z,
				AssignmentIndex: encoding.ChunkNumber(i),
			}
			samples = append(samples, sample)
		}
	}

	require.True(t, v.UniversalVerifySubBatch(params, samples, numBlob) == nil, "universal batch verification failed\n")
}

func TestUniversalVerifyWithPowerOf2G2(t *testing.T) {
	harness := getTestHarness()
	group, err := prover.NewProver(harness.proverV2KzgConfig, nil)
	require.Nil(t, err)

	committer, err := committer.NewFromConfig(*harness.committerConfig)
	require.Nil(t, err)

	v, err := verifier.NewVerifier(harness.verifierV2KzgConfig)
	require.NoError(t, err)

	params := encoding.ParamsFromSysPar(harness.numSys, harness.numPar, uint64(len(harness.paddedGettysburgAddressBytes)))
	blobLength := uint64(encoding.GetBlobLengthPowerOf2(uint32(len(harness.paddedGettysburgAddressBytes))))
	prover, err := group.GetKzgProver(params, blobLength)
	require.NoError(t, err)

	numBlob := 5
	samples := make([]encoding.Sample, 0)
	for z := 0; z < numBlob; z++ {
		inputFr, err := rs.ToFrArray(harness.paddedGettysburgAddressBytes)
		require.Nil(t, err)

		commit, _, _, err := committer.GetCommitments(inputFr)
		require.Nil(t, err)
		frames, fIndices, err := prover.GetFrames(inputFr)
		require.Nil(t, err)

		// create samples
		for i := 0; i < len(frames); i++ {
			f := frames[i]
			j := fIndices[i]

			q, err := rs.GetLeadingCosetIndex(uint64(i), harness.numSys+harness.numPar)
			require.Nil(t, err)

			require.Equal(t, j, q, "leading coset inconsistency")

			sample := encoding.Sample{
				Commitment:      (*encoding.G1Commitment)(commit),
				Chunk:           &f,
				BlobIndex:       z,
				AssignmentIndex: encoding.ChunkNumber(i),
			}
			samples = append(samples, sample)
		}
	}

	require.True(t, v.UniversalVerifySubBatch(params, samples, numBlob) == nil, "universal batch verification failed\n")
}

func TestBenchmarkVerifyChunks(t *testing.T) {
	t.Skip("This test is meant to be run manually, not as part of the test suite")

	harness := getTestHarness()

	p, err := prover.NewProver(harness.proverV2KzgConfig, nil)
	require.NoError(t, err)

	committer, err := committer.NewFromConfig(*harness.committerConfig)
	require.Nil(t, err)

	v, err := verifier.NewVerifier(harness.verifierV2KzgConfig)
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
