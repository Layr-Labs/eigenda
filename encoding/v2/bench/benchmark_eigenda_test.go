package bench_test

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Layr-Labs/eigenda/api/clients/codecs"
	"github.com/Layr-Labs/eigenda/api/clients/v2/coretypes"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/icicle"
	"github.com/Layr-Labs/eigenda/encoding/v2/fft"
	"github.com/Layr-Labs/eigenda/encoding/v2/kzg/committer"
	"github.com/Layr-Labs/eigenda/encoding/v2/kzg/prover"
	"github.com/Layr-Labs/eigenda/encoding/v2/rs"
	"github.com/Layr-Labs/eigenda/encoding/v2/rs/backend"
	"github.com/Layr-Labs/eigenda/encoding/v2/rs/backend/gnark"
	rsicicle "github.com/Layr-Labs/eigenda/encoding/v2/rs/backend/icicle"
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
	config := committer.Config{
		SRSNumberToLoad:   1 << 19, // 2^19 = 524,288 field elements = 16 MiB
		G1SRSPath:         "../../../resources/srs/g1.point",
		G2SRSPath:         "../../../resources/srs/g2.point",
		G2TrailingSRSPath: "../../../resources/srs/g2.trailing.point",
	}
	committer, err := committer.NewFromConfig(config)
	require.NoError(b, err)

	for _, blobPower := range []uint8{17, 20, 21, 24} {
		b.Run("Commitments_size_2^"+fmt.Sprint(blobPower)+"_bytes", func(b *testing.B) {
			blobLen := uint64(1 << blobPower / encoding.BYTES_PER_SYMBOL)
			rand := random.NewTestRandomNoPrint(1337)
			blob := rand.FrElements(blobLen)

			for b.Loop() {
				_, _, _, err := committer.GetCommitments(blob)
				require.NoError(b, err)
			}
		})
	}
}

// TODO(samlaf): maybe move this to benchmark_icicle_test.go file?
// That file is currently metal only, we should generalize it.
func BenchmarkRSBackendIcicle(b *testing.B) {
	if !icicle.IsAvailable {
		b.Skip("code compiled without the icicle build tag")
	}
	icicleBackend, err := rsicicle.BuildRSBackend(common.SilentLogger(), true)
	require.NoError(b, err)
	benchmarkRSBackend(b, icicleBackend)
}

func BenchmarkRSBackendGnark(b *testing.B) {
	fs := fft.NewFFTSettings(24)
	gnarkBackend := gnark.NewRSBackend(fs)
	benchmarkRSBackend(b, gnarkBackend)
}

func benchmarkRSBackend(b *testing.B, rsBackend backend.RSEncoderBackend) {
	rand := random.NewTestRandomNoPrint(1337)
	blobCoeffs := rand.FrElements(1 << 22) // max size we benchmark below: 24+3-5=22
	for _, blobPowerBytes := range []uint8{17, 20, 21, 24} {
		// Reed-Solomon encoding with 8x redundancy: 2^3 = 8
		rsExtendedBlobPowerBytes := blobPowerBytes + 3
		rsExtendedBlobPowerFrs := rsExtendedBlobPowerBytes - 5 // 32 bytes per Fr element
		b.Run("2^"+fmt.Sprint(rsExtendedBlobPowerFrs)+"_Frs", func(b *testing.B) {
			numFrs := uint64(1) << rsExtendedBlobPowerFrs
			// run multiple goroutines in parallel to better utilize the GPU
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					_, err := rsBackend.ExtendPolyEval(blobCoeffs[:numFrs])
					require.NoError(b, err)
				}
			})
		})
	}
}

// Dispersers need to encode blobs into chunks before dispersing them.
// This entails Reed-Solomon encoding the blob into 8x its size,
// creating 8192 chunks of size 8*blobLen/8192 Field elements each,
// and computing for each chunk the coefficients of the polynomial that
// evaluates to the chunk's data at the chunk's coset indices.
func BenchmarkBlobToChunksEncoding(b *testing.B) {
	cfg := encoding.DefaultConfig()
	enc := rs.NewEncoder(common.SilentLogger(), cfg)

	for _, blobPower := range []uint64{17, 20, 21, 24} {
		b.Run("Encode_size_2^"+fmt.Sprint(blobPower)+"_bytes", func(b *testing.B) {
			blobSizeBytes := uint64(1) << blobPower
			params := encoding.EncodingParams{
				NumChunks:   8192,                            // blob_version=0
				ChunkLength: max(1, blobSizeBytes*8/8192/32), // chosen such that numChunks*ChunkLength=blobSize
			}

			rand := random.NewTestRandomNoPrint(1337)
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

// TODO(samlaf): maybe move this to benchmark_icicle_test.go file?
// That file is currently metal only, we should generalize it.
func BenchmarkMultiproofGenerationIcicle(b *testing.B) {
	if !icicle.IsAvailable {
		b.Skip("code compiled without the icicle build tag")
	}
	encodingConfig := encoding.Config{
		NumWorker:   uint64(runtime.NumCPU()),
		BackendType: encoding.IcicleBackend,
		GPUEnable:   true,
	}
	benchmarkMultiproofGeneration(b, encodingConfig)
}

func BenchmarkMultiproofGenerationGnark(b *testing.B) {
	encodingConfig := encoding.Config{
		NumWorker:   uint64(runtime.NumCPU()),
		BackendType: encoding.GnarkBackend,
		GPUEnable:   false,
	}
	benchmarkMultiproofGeneration(b, encodingConfig)
}

// The encoder service on the disperser generates a multiproof for each chunk.
// This is the most intensive part of the encoding process.
//
// The benchmark uses a silent logger, but you can switch to a normal logger to see
// the log lines giving a breakdown of the different proof steps. E.g.:
// Multiproof Time Decomp total=9.478006875s preproc=33.987083ms msm=1.496717042s fft1=5.912448708s fft2=2.034854042s
// Where fft1 and fft2 are on G1, and preproc contains an FFT on Fr elements.
func benchmarkMultiproofGeneration(b *testing.B, encodingConfig encoding.Config) {
	proverConfig := prover.KzgConfig{
		// The loaded G1 point is not used because we require the SRSTables to be preloaded for the benchmark.
		// We don't have enough SRS points in resourcs/srs/g1.point to compute the largest SRSTables anyways.
		// Note that we can't input 0 here because the prover checks that at least 1 point is loaded.
		// TODO(samlaf): fix this. We should be able to not load any G1 points if we are preloading the SRSTables.
		SRSNumberToLoad: 1,
		G1Path:          "../../../resources/srs/g1.point",
		// make sure to run `make download_srs_tables` to have the SRSTables available here.
		PreloadEncoder: true,
		CacheDir:       "../../../resources/srs/SRSTables",
		NumWorker:      uint64(runtime.GOMAXPROCS(0)),
	}
	b.Log("Reading precomputed SRSTables, this may take a while...")
	// use a non-silent logger to see the "Multiproof Time Decomp" log lines.
	p, err := prover.NewProver(common.TestLogger(b), &proverConfig, &encodingConfig)
	require.NoError(b, err)

	rand := random.NewTestRandomNoPrint(1337)
	maxSizeBlobCoeffs := rand.FrElements(1 << 22)

	for _, blobPowerBytes := range []uint64{17, 20, 21, 24} {
		b.Run("Multiproof_size_2^"+fmt.Sprint(blobPowerBytes)+"_bytes", func(b *testing.B) {
			// Reed-Solomon encoding with 8x redundancy: 2^3 = 8
			rsExtendedBlobPowerBytes := blobPowerBytes + 3
			rsExtendedBlobPowerFrs := rsExtendedBlobPowerBytes - 5 // 32 bytes per Fr element
			rsExtendedBlobFrs := uint64(1) << rsExtendedBlobPowerFrs
			params := encoding.EncodingParams{
				NumChunks:   8192,                           // blob_version=0
				ChunkLength: max(1, rsExtendedBlobFrs/8192), // chosen such that numChunks*ChunkLength=rsExtendedBlobFrs
			}
			provingParams := prover.ProvingParams{
				BlobLength:  rsExtendedBlobFrs,
				ChunkLength: max(1, rsExtendedBlobFrs/8192), // chosen such that numChunks*ChunkLength=rsExtendedBlobFrs
			}
			parametrizedProver, err := p.GetKzgProver(params, provingParams)
			require.NoError(b, err)

			for b.Loop() {
				_, err = parametrizedProver.GetProofs(maxSizeBlobCoeffs[:rsExtendedBlobFrs], provingParams)
				require.NoError(b, err)
			}
		})
	}
}
