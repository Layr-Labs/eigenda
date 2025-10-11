package bench

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Layr-Labs/eigenda/api/clients/codecs"
	"github.com/Layr-Labs/eigenda/api/clients/v2/coretypes"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/kzg/committer"
	"github.com/Layr-Labs/eigenda/encoding/kzg/prover/v2"
	"github.com/Layr-Labs/eigenda/encoding/rs"
	"github.com/Layr-Labs/eigenda/test/random"
)

// This file contains benchmarks for the high-level math/crypto operations that are
// performed by different actors of the EigenDA network:
// - Clients: PayloadToBlob conversion, Commitment generation
// - Dispersers: Frame generation (RS encoding into chunks + KZG multiproof generation)
// - Validators: Verification of commitments and proofs (TODO: write benchmark for this)

// Before sending their payload to EigenDA, clients need to convert it into a Blob.
// Turning a user payload into a Blob (bn254 Field elements representing coefficients of a polynomial)
// requires encoding the payload into Field Elements, and then possibly doing an IFFT
// if the user interprets his encoded_payload as evaluations instead of coefficients.
func BenchmarkPayloadToBlobConversion(b *testing.B) {
	for _, blobPower := range []uint8{17, 20, 21, 24} {
		b.Run("PayloadToBlob_size_2^"+fmt.Sprint(blobPower)+"_bytes", func(b *testing.B) {
			numSymbols := uint64(1<<blobPower) / 32
			payloadBytesPerSymbols := uint64(encoding.BYTES_PER_SYMBOL - 1)
			payloadBytes := make([]byte, numSymbols*payloadBytesPerSymbols)
			for i := range numSymbols {
				payloadBytes[i*payloadBytesPerSymbols] = byte(i + 1)
			}
			payload := coretypes.Payload(payloadBytes)

			for b.Loop() {
				_, err := payload.ToBlob(codecs.PolynomialFormEval)
				require.NoError(b, err)
			}
		})
	}
}

// Before making a dispersal, clients need to generate commitments for their blob,
// which are included as part of the BlobHeader in the dispersal request.
// This benchmark measures the total time it takes to generate all 3 commitments:
// blob commitment (G1 MSM), blob length commitment (G2 MSM), and blob length proof (G2 MSM).
// The committer package contains benchmarks for each individual commitment,
// since those are private functions that we can't call from here.
func BenchmarkCommittmentGeneration(b *testing.B) {
	blobLen := uint64(1 << 19) // 2^19 = 524,288 field elements = 16 MiB
	config := committer.Config{
		SRSNumberToLoad:   blobLen,
		G1SRSPath:         "../../resources/srs/g1.point",
		G2SRSPath:         "../../resources/srs/g2.point",
		G2TrailingSRSPath: "../../resources/srs/g2.trailing.point",
	}
	committer, err := committer.NewFromConfig(config)
	require.NoError(b, err)

	rand := random.NewTestRandom()
	blob := rand.FrElements(blobLen)

	for b.Loop() {
		_, _, _, err := committer.GetCommitments(blob)
		require.NoError(b, err)
	}
}

// Dispersers need to encode blobs into chunks before dispersing them.
// This entails Reed-Solomon encoding the blob into 8x its size,
// creating 8192 chunks of size 8*blobLen/8192 Field elements each,
// and computing for each chunk the coefficients of the polynomial that
// evaluates to the chunk's data at the chunk's coset indices.
func BenchmarkBlobToChunksEncoding(b *testing.B) {
	cfg := encoding.DefaultConfig()
	enc := rs.NewEncoder(cfg)

	for _, blobPower := range []uint64{17, 20, 24} {
		b.Run("Encode_size_2^"+fmt.Sprint(blobPower)+"_bytes", func(b *testing.B) {
			blobSizeBytes := uint64(1) << blobPower
			params := encoding.EncodingParams{
				NumChunks:   8192,                            // blob_version=0
				ChunkLength: max(1, blobSizeBytes*8/8192/32), // chosen such that numChunks*ChunkLength=blobSize
			}

			rand := random.NewTestRandom()
			blobBytes := rand.Bytes(int(blobSizeBytes))
			for i := 0; i < len(blobBytes); i += 32 {
				blobBytes[i] = 0 // to make them Fr elements
			}
			blob, err := rs.ToFrArray(blobBytes)
			require.Nil(b, err)

			for b.Loop() {
				_, _, err = enc.Encode(blob, params)
				require.Nil(b, err)
			}
		})
	}
}

// The encoder service on the disperser generates a multiproof for each chunk.
// This is the most intensive part of the encoding process.
//
// This Benchmark is very high-level, since GetFrames does many things.
// But the print statements from the Encoder give a breakdown of the different steps. E.g.:
// Multiproof Time Decomp total=9.478006875s preproc=33.987083ms msm=1.496717042s fft1=5.912448708s fft2=2.034854042s
// Where fft1 and fft2 are on G1, and preproc contains an FFT on Fr elements.
func BenchmarkMultiproofFrameGeneration(b *testing.B) {
	proverConfig := prover.KzgConfig{
		G1Path:          "../../resources/srs/g1.point",
		G2Path:          "../../resources/srs/g2.point",
		G2TrailingPath:  "../../resources/srs/g2.trailing.point",
		CacheDir:        "../../resources/srs/SRSTables",
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
		NumChunks:   8192,                       // blob_version=0
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
