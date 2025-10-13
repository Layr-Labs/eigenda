package prover_test

import (
	cryptorand "crypto/rand"
	"math/rand"
	"testing"

	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/kzg/committer"
	"github.com/Layr-Labs/eigenda/encoding/kzg/prover/v2"
	"github.com/Layr-Labs/eigenda/encoding/kzg/verifier/v2"
	"github.com/Layr-Labs/eigenda/encoding/rs"
	"github.com/Layr-Labs/eigenda/encoding/utils/codec"

	"github.com/stretchr/testify/require"
)

func sampleFrames(frames []*encoding.Frame, num uint64) ([]*encoding.Frame, []encoding.ChunkNumber) {
	samples := make([]*encoding.Frame, num)
	indices := rand.Perm(len(frames))
	indices = indices[:num]

	frameIndices := make([]encoding.ChunkNumber, num)
	for i, j := range indices {
		samples[i] = frames[j]
		frameIndices[i] = encoding.ChunkNumber(j)
	}
	return samples, frameIndices
}

func TestEncoder(t *testing.T) {
	harness := getTestHarness()
	p, err := prover.NewProver(harness.proverV2KzgConfig, nil)
	require.NoError(t, err)

	c, err := committer.NewFromConfig(*harness.committerConfig)
	require.NoError(t, err)

	v, err := verifier.NewVerifier(harness.verifierV2KzgConfig)
	require.NoError(t, err)

	encoder := rs.NewEncoder(nil)

	params := encoding.ParamsFromMins(5, 5)
	commitments, err := c.GetCommitmentsForPaddedLength(harness.paddedGettysburgAddressBytes)
	require.NoError(t, err)
	frames, err := p.GetFrames(harness.paddedGettysburgAddressBytes, params)
	require.NoError(t, err)

	indices := []encoding.ChunkNumber{
		0, 1, 2, 3, 4, 5, 6, 7,
	}
	err = v.VerifyFrames(frames, indices, commitments, params)
	require.NoError(t, err)
	err = v.VerifyFrames(frames, []encoding.ChunkNumber{
		7, 6, 5, 4, 3, 2, 1, 0,
	}, commitments, params)
	require.Error(t, err)

	maxInputSize := uint64(len(harness.paddedGettysburgAddressBytes))
	chunks := make([]rs.FrameCoeffs, len(frames))
	for i, f := range frames {
		chunks[i] = f.Coeffs
	}
	decoded, err := encoder.Decode(chunks, indices, maxInputSize, params)
	require.NoError(t, err)
	require.Equal(t, harness.paddedGettysburgAddressBytes, decoded)

	// shuffle frames
	tmp := frames[2]
	frames[2] = frames[5]
	frames[5] = tmp
	indices = []encoding.ChunkNumber{
		0, 1, 5, 3, 4, 2, 6, 7,
	}

	err = v.VerifyFrames(frames, indices, commitments, params)
	require.NoError(t, err)

	chunks = make([]rs.FrameCoeffs, len(frames))
	for i, f := range frames {
		chunks[i] = f.Coeffs
	}
	decoded, err = encoder.Decode(chunks, indices, maxInputSize, params)
	require.NoError(t, err)
	require.Equal(t, harness.paddedGettysburgAddressBytes, decoded)
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
	_, err = p.GetFrames(blobs[0], params)
	require.NoError(b, err)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err = p.GetFrames(blobs[i%numSamples], params)
		require.NoError(b, err)
	}
}

func FuzzOnlySystematic(f *testing.F) {
	harness := getTestHarness()

	f.Add(harness.paddedGettysburgAddressBytes)
	f.Add([]byte("Hello, World!"))
	f.Add([]byte{0})

	f.Fuzz(func(t *testing.T, input []byte) {
		input = codec.ConvertByPaddingEmptyByte(input)
		group, err := prover.NewProver(harness.proverV2KzgConfig, nil)
		require.NoError(t, err)

		params := encoding.ParamsFromSysPar(10, 3, uint64(len(input)))

		//encode the data
		frames, err := group.GetFrames(input, params)
		require.NoError(t, err)

		for _, frame := range frames {
			require.NotEqual(t, len(frame.Coeffs), 0)
		}

		if err != nil {
			t.Errorf("Error Encoding:\n Data:\n %q \n Err: %q", input, err)
		}

		//sample the correct systematic frames
		samples, indices := sampleFrames(frames, uint64(len(frames)))

		encoder := rs.NewEncoder(nil)
		chunks := make([]rs.FrameCoeffs, len(samples))
		for i, f := range samples {
			chunks[i] = f.Coeffs
		}
		data, err := encoder.Decode(chunks, indices, uint64(len(input)), params)
		if err != nil {
			t.Errorf("Error Decoding:\n Data:\n %q \n Err: %q", input, err)
		}
		require.Equal(t, input, data, "Input data was not equal to the decoded data")
	})
}
