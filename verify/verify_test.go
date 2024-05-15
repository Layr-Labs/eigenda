package verify

import (
	"encoding/hex"
	"runtime"
	"testing"

	"github.com/Layr-Labs/eigenda/api/grpc/common"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/Layr-Labs/op-plasma-eigenda/eigenda"
	"github.com/stretchr/testify/assert"
)

func TestVerification(t *testing.T) {

	var data = []byte("inter-subjective and not objective!")

	x, err := hex.DecodeString("0b187c5351919a9bf83271637be3bcb7b8bbb0abe0b80bb9d632ad8f6e8401e5")
	assert.NoError(t, err)

	y, err := hex.DecodeString("0d41ee143f13cc2526d36189a22538f630ea31398e0af32b5877728c8fe5452e")
	assert.NoError(t, err)

	c := eigenda.Cert{
		BlobCommitment: &common.G1Commitment{
			X: x,
			Y: y,
		},
	}

	kzgConfig := &kzg.KzgConfig{
		G1Path:          "../test/resources/g1.point",
		G2Path:          "../test/resources/g2.point",
		G2PowerOf2Path:  "../test/resources/g2.point.powerOf2",
		CacheDir:        "../test/resources/SRSTables",
		SRSOrder:        3000,
		SRSNumberToLoad: 3000,
		NumWorker:       uint64(runtime.GOMAXPROCS(0)),
	}

	v, err := NewVerifier(kzgConfig)
	assert.NoError(t, err)

	err = v.Verify(c, eigenda.EncodeToBlob(data))
	assert.NoError(t, err)
}
