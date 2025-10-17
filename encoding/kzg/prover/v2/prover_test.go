package prover_test

import (
	"math/rand"
	"testing"

	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/codec"
	"github.com/Layr-Labs/eigenda/encoding/kzg/prover/v2"
	"github.com/Layr-Labs/eigenda/encoding/kzg/verifier/v2"
	"github.com/Layr-Labs/eigenda/encoding/rs"
	"github.com/Layr-Labs/eigenda/encoding/v2/kzg/committer"

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
	harness := getTestHarness(t)
	p, err := prover.NewProver(harness.logger, harness.proverV2KzgConfig, nil)
	require.NoError(t, err)

	c, err := committer.NewFromConfig(*harness.committerConfig)
	require.NoError(t, err)

	v, err := verifier.NewVerifier(harness.verifierV2KzgConfig)
	require.NoError(t, err)

	encoder := rs.NewEncoder(harness.logger, nil)

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

func FuzzOnlySystematic(f *testing.F) {
	harness := getTestHarness(f)

	f.Add(harness.paddedGettysburgAddressBytes)
	f.Add([]byte("Hello, World!"))
	f.Add([]byte{0})

	f.Fuzz(func(t *testing.T, input []byte) {
		input = codec.ConvertByPaddingEmptyByte(input)
		group, err := prover.NewProver(harness.logger, harness.proverV2KzgConfig, nil)
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

		encoder := rs.NewEncoder(harness.logger, nil)
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
