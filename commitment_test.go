package plasma

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/Layr-Labs/eigenda/api/grpc/common"
	"github.com/Layr-Labs/op-plasma-eigenda/eigenda"
	"github.com/ethereum/go-ethereum/rlp"
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

	println(fmt.Sprintf("x: %+x", x))
	println(fmt.Sprintf("y: %+x", y))

	b, err := rlp.EncodeToBytes(c)
	assert.NoError(t, err)

	var commit EigenDACommitment = b

	err = commit.Verify(data)
	assert.NoError(t, err)
}
