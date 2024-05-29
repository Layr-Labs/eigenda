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

	x, err := hex.DecodeString("07c23d7720de3f10064c8f48774d8f59207964c482419063246a67e1c454a886")
	assert.NoError(t, err)

	y, err := hex.DecodeString("0f747070e6fdb4e1346fec54dbc3d2d61a2c9ad2cb6b1744fa7f47072ad13370")
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
