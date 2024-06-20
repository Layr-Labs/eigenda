package verify

import (
	"encoding/hex"
	"runtime"
	"testing"

	"github.com/Layr-Labs/eigenda/api/clients/codecs"
	"github.com/Layr-Labs/eigenda/api/grpc/common"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/stretchr/testify/assert"
)

func TestCommitmentVerification(t *testing.T) {
	t.Parallel()

	var data = []byte("inter-subjective and not objective!")

	x, err := hex.DecodeString("2fc55f968a2d29d22aebf55b382528d1d9401577c166483e162355b19d8bc446")
	assert.NoError(t, err)

	y, err := hex.DecodeString("149e2241c21c391e069b9f317710c7f57f31ee88245a5e61f0d294b11acf9aff")
	assert.NoError(t, err)

	c := &common.G1Commitment{
		X: x,
		Y: y,
	}

	kzgConfig := &kzg.KzgConfig{
		G1Path:          "../e2e/resources/kzg/g1.point",
		G2PowerOf2Path:  "../e2e/resources/kzg/g2.point.powerOf2",
		CacheDir:        "../e2e/resources/kzg/SRSTables",
		SRSOrder:        3000,
		SRSNumberToLoad: 3000,
		NumWorker:       uint64(runtime.GOMAXPROCS(0)),
	}

	cfg := &Config{
		Verify:    false,
		KzgConfig: kzgConfig,
	}

	v, err := NewVerifier(cfg, nil)
	assert.NoError(t, err)

	// Happy path verification
	codec := codecs.NewIFFTCodec(codecs.NewDefaultBlobCodec())
	blob, err := codec.EncodeBlob(data)
	assert.NoError(t, err)
	err = v.VerifyCommitment(c, blob)
	assert.NoError(t, err)

	// failure with wrong data
	fakeData, err := codec.EncodeBlob([]byte("I am an imposter!!"))
	assert.NoError(t, err)
	err = v.VerifyCommitment(c, fakeData)
	assert.Error(t, err)
}
