package prover_test

import (
	cryptorand "crypto/rand"
	"math/rand"
	"testing"

	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/kzg/prover/v2"
	"github.com/Layr-Labs/eigenda/encoding/kzg/verifier/v2"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func sampleFrames(frames []encoding.Frame, num uint64) ([]encoding.Frame, []uint64) {
	samples := make([]encoding.Frame, num)
	indices := rand.Perm(len(frames))
	indices = indices[:num]

	frameIndices := make([]uint64, num)
	for i, j := range indices {
		samples[i] = frames[j]
		frameIndices[i] = uint64(j)
	}
	return samples, frameIndices
}

func TestEncoder(t *testing.T) {
	harness := getTestHarness()
	p, err := prover.NewProver(harness.proverV2KzgConfig, nil)
	require.NoError(t, err)

	v, err := verifier.NewVerifier(harness.verifierV2KzgConfig, nil)
	require.NoError(t, err)

	params := encoding.ParamsFromMins(5, 5)
	commitments, chunks, err := p.EncodeAndProve(harness.paddedGettysburgAddressBytes, params)
	assert.NoError(t, err)

	indices := []encoding.ChunkNumber{
		0, 1, 2, 3, 4, 5, 6, 7,
	}
	err = v.VerifyFrames(chunks, indices, commitments, params)
	assert.NoError(t, err)
	err = v.VerifyFrames(chunks, []encoding.ChunkNumber{
		7, 6, 5, 4, 3, 2, 1, 0,
	}, commitments, params)
	assert.Error(t, err)

	maxInputSize := uint64(len(harness.paddedGettysburgAddressBytes))
	decoded, err := p.Decode(chunks, indices, params, maxInputSize)
	assert.NoError(t, err)
	assert.Equal(t, harness.paddedGettysburgAddressBytes, decoded)

	// shuffle chunks
	tmp := chunks[2]
	chunks[2] = chunks[5]
	chunks[5] = tmp
	indices = []encoding.ChunkNumber{
		0, 1, 5, 3, 4, 2, 6, 7,
	}

	err = v.VerifyFrames(chunks, indices, commitments, params)
	assert.NoError(t, err)

	decoded, err = p.Decode(chunks, indices, params, maxInputSize)
	assert.NoError(t, err)
	assert.Equal(t, harness.paddedGettysburgAddressBytes, decoded)
}

// Ballpark number for 400KiB blob encoding
//
// goos: darwin
// goarch: arm64
// pkg: github.com/Layr-Labs/eigenda/core/encoding
// BenchmarkEncode-12    	       1	2421900583 ns/op
func BenchmarkEncode(b *testing.B) {
	harness := getTestHarness()
	p, err := prover.NewProver(harness.proverV2KzgConfig, nil)
	require.NoError(b, err)

	params := encoding.EncodingParams{
		ChunkLength: 512,
		NumChunks:   256,
	}
	blobSize := 400 * 1024
	numSamples := 30
	blobs := make([][]byte, numSamples)
	for i := 0; i < numSamples; i++ {
		blob := make([]byte, blobSize)
		_, _ = cryptorand.Read(blob)
		blobs[i] = blob
	}

	// Warm up the encoder: ensures that all SRS tables are loaded so these aren't included in the benchmark.
	_, _, _ = p.EncodeAndProve(blobs[0], params)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _, _ = p.EncodeAndProve(blobs[i%numSamples], params)
	}
}
