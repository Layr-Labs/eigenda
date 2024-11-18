package verifier_test

import (
	"crypto/rand"
	"fmt"
	"log"
	"os"
	"runtime"
	"testing"

	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/Layr-Labs/eigenda/encoding/kzg/prover"
	"github.com/Layr-Labs/eigenda/encoding/kzg/verifier"
	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	gettysburgAddressBytes = codec.ConvertByPaddingEmptyByte([]byte("Fourscore and seven years ago our fathers brought forth, on this continent, a new nation, conceived in liberty, and dedicated to the proposition that all men are created equal. Now we are engaged in a great civil war, testing whether that nation, or any nation so conceived, and so dedicated, can long endure. We are met on a great battle-field of that war. We have come to dedicate a portion of that field, as a final resting-place for those who here gave their lives, that that nation might live. It is altogether fitting and proper that we should do this. But, in a larger sense, we cannot dedicate, we cannot consecrate—we cannot hallow—this ground. The brave men, living and dead, who struggled here, have consecrated it far above our poor power to add or detract. The world will little note, nor long remember what we say here, but it can never forget what they did here. It is for us the living, rather, to be dedicated here to the unfinished work which they who fought here have thus far so nobly advanced. It is rather for us to be here dedicated to the great task remaining before us—that from these honored dead we take increased devotion to that cause for which they here gave the last full measure of devotion—that we here highly resolve that these dead shall not have died in vain—that this nation, under God, shall have a new birth of freedom, and that government of the people, by the people, for the people, shall not perish from the earth."))
	kzgConfig              *kzg.KzgConfig
	numNode                uint64
	numSys                 uint64
	numPar                 uint64
)

func TestMain(m *testing.M) {
	setup()
	result := m.Run()
	teardown()
	os.Exit(result)
}

func setup() {
	log.Println("Setting up suite")

	kzgConfig = &kzg.KzgConfig{
		G1Path:          "../../../inabox/resources/kzg/g1.point",
		G2Path:          "../../../inabox/resources/kzg/g2.point",
		G2PowerOf2Path:  "../../../inabox/resources/kzg/g2.point.powerOf2",
		CacheDir:        "../../../inabox/resources/kzg/SRSTables",
		SRSOrder:        3000,
		SRSNumberToLoad: 2900,
		NumWorker:       uint64(runtime.GOMAXPROCS(0)),
		LoadG2Points:    true,
	}

	numNode = uint64(4)
	numSys = uint64(3)
	numPar = numNode - numSys

}

func teardown() {
	log.Println("Tearing down")
	os.RemoveAll("./data")
}

// var control interface{ Stop() }

func TestBenchmarkVerifyChunks(t *testing.T) {
	t.Skip("This test is meant to be run manually, not as part of the test suite")
	p, err := prover.NewProver(kzgConfig, nil)
	require.NoError(t, err)

	v, err := verifier.NewVerifier(kzgConfig, nil)
	require.NoError(t, err)

	chunkLengths := []uint64{64, 128, 256, 512, 1024, 2048, 4096, 8192}
	chunkCounts := []int{4, 8, 16}

	file, err := os.Create("benchmark_results.csv")
	if err != nil {
		t.Fatalf("Failed to open file for writing: %v", err)
	}
	defer file.Close()

	fmt.Fprintln(file, "numChunks,chunkLength,ns/op,allocs/op")

	for _, chunkLength := range chunkLengths {

		blobSize := chunkLength * 32 * 2
		params := encoding.EncodingParams{
			ChunkLength: chunkLength,
			NumChunks:   16,
		}
		blob := make([]byte, blobSize)
		_, err = rand.Read(blob)
		assert.NoError(t, err)

		commitments, chunks, err := p.EncodeAndProve(blob, params)
		assert.NoError(t, err)

		indices := make([]encoding.ChunkNumber, params.NumChunks)
		for i := range indices {
			indices[i] = encoding.ChunkNumber(i)
		}

		for _, numChunks := range chunkCounts {

			result := testing.Benchmark(func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					// control = profile.Start(profile.ProfilePath("."))
					err := v.VerifyFrames(chunks[:numChunks], indices[:numChunks], commitments, params)
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
	p, err := prover.NewProver(kzgConfig, nil)
	require.NoError(b, err)

	v, err := verifier.NewVerifier(kzgConfig, nil)
	require.NoError(b, err)

	params := encoding.EncodingParams{
		ChunkLength: 256,
		NumChunks:   8,
	}
	blobSize := 8 * 256
	numSamples := 30
	blobs := make([][]byte, numSamples)
	for i := 0; i < numSamples; i++ {
		blob := make([]byte, blobSize)
		_, _ = rand.Read(blob)
		blobs[i] = blob
	}

	commitments, _, err := p.EncodeAndProve(blobs[0], params)
	assert.NoError(b, err)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err = v.VerifyBlobLength(commitments)
		assert.NoError(b, err)
	}

}
