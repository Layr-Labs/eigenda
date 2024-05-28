package verify

import (
	"encoding/hex"
	"runtime"
	"testing"

	"github.com/Layr-Labs/eigenda-proxy/eigenda"
	"github.com/Layr-Labs/eigenda/api/clients/codecs"
	"github.com/Layr-Labs/eigenda/api/grpc/common"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/stretchr/testify/assert"
)

func TestVerification(t *testing.T) {

	var data = []byte("inter-subjective and not objective!")

	x, err := hex.DecodeString("0184B47F64FBA17D6F49CDFED20434B1015A2A369AB203256EC4CD00C324E83B")
	assert.NoError(t, err)

	y, err := hex.DecodeString("122CD859CC5CDD048B482C50721821CB413C151BA7AF10285C1D2483F2A88085")
	assert.NoError(t, err)

	c := eigenda.Cert{
		BlobCommitment: &common.G1Commitment{
			X: x,
			Y: y,
		},
	}

	kzgConfig := &kzg.KzgConfig{
		G1Path:          "../test/resources/g1.point",
		G2PowerOf2Path:  "../test/resources/g2.point.powerOf2",
		CacheDir:        "../test/resources/SRSTables",
		SRSOrder:        3000,
		SRSNumberToLoad: 3000,
		NumWorker:       uint64(runtime.GOMAXPROCS(0)),
	}

	v, err := NewVerifier(kzgConfig)
	assert.NoError(t, err)

	// Happy path verification

	// TODO: Update this test to use the IFFT codec
	codec := codecs.DefaultBlobEncodingCodec{}
	blob, err := codec.EncodeBlob(data)
	assert.NoError(t, err)
	err = v.Verify(c, blob)
	assert.NoError(t, err)

	// failure with wrong data
	fakeData, err := codec.EncodeBlob([]byte("I am an imposter!!"))
	assert.NoError(t, err)
	err = v.Verify(c, fakeData)
	assert.Error(t, err)
}
