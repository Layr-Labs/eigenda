package prover_test

import (
	cryptorand "crypto/rand"
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
	log.Println("Tearing down suite")

	// Some test may want to create a new SRS table so this should clean it up.
	os.RemoveAll("./data")
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

	params := encoding.ParamsFromMins(5, 5)
	commitments, chunks, err := p.EncodeAndProve(gettysburgAddressBytes, params)
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

	maxInputSize := uint64(len(gettysburgAddressBytes))
	decoded, err := p.Decode(chunks, indices, params, maxInputSize)
	assert.NoError(t, err)
	assert.Equal(t, gettysburgAddressBytes, decoded)

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
	assert.Equal(t, gettysburgAddressBytes, decoded)
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
