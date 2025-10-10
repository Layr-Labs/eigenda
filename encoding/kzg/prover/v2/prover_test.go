package prover_test

import (
	"math/rand"
	"runtime"
	"testing"

	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/kzg/committer"
	"github.com/Layr-Labs/eigenda/encoding/kzg/prover/v2"
	"github.com/Layr-Labs/eigenda/encoding/kzg/verifier/v2"
	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
	"github.com/Layr-Labs/eigenda/test/random"

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

	c, err := committer.NewFromConfig(committer.Config{
		SRSNumberToLoad:   harness.proverV2KzgConfig.SRSNumberToLoad,
		G1SRSPath:         harness.proverV2KzgConfig.G1Path,
		G2SRSPath:         harness.proverV2KzgConfig.G2Path,
		G2TrailingSRSPath: harness.proverV2KzgConfig.G2TrailingPath,
	})
	require.NoError(t, err)

	v, err := verifier.NewVerifier(harness.verifierV2KzgConfig, nil)
	require.NoError(t, err)

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
	decoded, err := v.Decode(frames, indices, params, maxInputSize)
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

	decoded, err = v.Decode(frames, indices, params, maxInputSize)
	require.NoError(t, err)
	require.Equal(t, harness.paddedGettysburgAddressBytes, decoded)
}

// This Benchmark is very high-level, since GetFrames does many things.
// The benchmark itself is roughly always ~8-10seconds on M4 Macbook Pro.
// But the print statements from the Encoder give a breakdown of the different steps:
// eg: Multiproof Time Decomp total=9.478006875s preproc=33.987083ms msm=1.496717042s fft1=5.912448708s fft2=2.034854042s
// Where fft1 and fft2 are on G1, and preproc contains an FFT on Fr elements.
func BenchmarkGetFrame(b *testing.B) {
	proverConfig := prover.KzgConfig{
		G1Path:          "../../../../resources/srs/g1.point",
		G2Path:          "../../../../resources/srs/g2.point",
		G2TrailingPath:  "../../../../resources/srs/g2.trailing.point",
		CacheDir:        "../../../../resources/srs/SRSTables",
		SRSNumberToLoad: 1 << 19,
		NumWorker:       uint64(runtime.GOMAXPROCS(0)),
		LoadG2Points:    true,
	}
	p, err := prover.NewProver(&proverConfig, nil)
	require.NoError(b, err)

	// We only have 16MiBs of SRS points. Since we use blob_version=0's 8x coding 
	// ratio, we create a blob of size 2MiB and 8x rs encode it up to 16MiB.
	blobSize := uint64(1) << 21 // 2 MiB
	params := encoding.EncodingParams{
		NumChunks:   8192,                     // blob_version=0
		ChunkLength: max(1, blobSize*8/8192/32), // chosen such that numChunks*ChunkLength=blobSize
	}

	rand := random.NewTestRandom()
	blobBytes := rand.Bytes(int(blobSize))
	for i := 0; i < len(blobBytes); i += 32 {
		blobBytes[i] = 0 // to make them Fr elements
	}

	// Warm up the encoder: ensures that all SRS tables are loaded so these aren't included in the benchmark.
	_, err = p.GetFrames(blobBytes, params)
	require.NoError(b, err)

	for b.Loop() {
		_, err = p.GetFrames(blobBytes, params)
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

		v, err := verifier.NewVerifier(harness.verifierV2KzgConfig, nil)
		require.NoError(t, err)
		data, err := v.Decode(samples, indices, params, uint64(len(input)))
		if err != nil {
			t.Errorf("Error Decoding:\n Data:\n %q \n Err: %q", input, err)
		}
		require.Equal(t, input, data, "Input data was not equal to the decoded data")
	})
}
