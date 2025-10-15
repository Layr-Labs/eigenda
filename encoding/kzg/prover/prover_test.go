package prover_test

import (
	cryptorand "crypto/rand"
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime"
	"testing"

	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/Layr-Labs/eigenda/encoding/kzg/prover"
	"github.com/Layr-Labs/eigenda/encoding/kzg/verifier"
	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
	"github.com/Layr-Labs/eigenda/test/random"

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
		G1Path:          "../../../resources/srs/g1.point",
		G2Path:          "../../../resources/srs/g2.point",
		CacheDir:        "../../../resources/srs/SRSTables",
		SRSOrder:        524288,
		SRSNumberToLoad: 524288,
		NumWorker:       uint64(runtime.GOMAXPROCS(0)),
		LoadG2Points:    true,
	}

	numNode = uint64(4)
	numSys = uint64(3)
	numPar = numNode - numSys

}

func teardown() {
	log.Println("Tearing down suite")

	// Some test may want to create a new SRS table so this should clean it up.
	err := os.RemoveAll("./data")
	if err != nil {
		log.Printf("Error removing data directory ./data: %v", err)
	}
}

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
	p, err := prover.NewProver(kzgConfig, nil)
	require.NoError(t, err)

	v, err := verifier.NewVerifier(kzgConfig, nil)
	require.NoError(t, err)

	blobLen := uint64(8192)

	rand := random.NewTestRandom()
	blob := rand.FrElements(blobLen)

	var gettysburgAddressBytes []byte
	for i := 0; i < 8192; i++ {
		b := blob[i].Bytes()
		gettysburgAddressBytes = append(gettysburgAddressBytes, b[:]...)
	}

	params := encoding.ParamsFromMins(8, 8192)
	blobLength := encoding.GetBlobLengthPowerOf2(uint32(len(gettysburgAddressBytes)))
	params.SetBlobLength(uint64(blobLength))
	fmt.Println("params", params)
	fmt.Println("blobLength", blobLength)
	commitments, chunks, err := p.EncodeAndProve(gettysburgAddressBytes, params)
	require.NoError(t, err)

	indices := make([]encoding.ChunkNumber, 8192)
	for i := 0; i < 8192; i++ {
		indices[i] = encoding.ChunkNumber(i)
	}
	err = v.VerifyFrames(chunks, indices, commitments, params)
	require.NoError(t, err)

}

// Ballpark number for 400KiB blob encoding
//
// goos: darwin
// goarch: arm64
// pkg: github.com/Layr-Labs/eigenda/core/encoding
// BenchmarkEncode-12    	       1	2421900583 ns/op
func BenchmarkEncode(b *testing.B) {
	p, err := prover.NewProver(kzgConfig, nil)
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
