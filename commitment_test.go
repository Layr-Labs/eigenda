package plasma

import (
	"encoding/binary"
	"encoding/hex"
	"testing"

	"github.com/Layr-Labs/eigenda/api/grpc/common"
	"github.com/Layr-Labs/op-plasma-eigenda/eigenda"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/assert"
)

// TODO - Add commitment verification logic
func TestVerifyCommitment(t *testing.T) {
	x, err := hex.DecodeString("0b187c5351919a9bf83271637be3bcb7b8bbb0abe0b80bb9d632ad8f6e8401e5")
	if err != nil {
		panic(err)
	}

	y, err := hex.DecodeString("0d41ee143f13cc2526d36189a22538f630ea31398e0af32b5877728c8fe5452e")
	if err != nil {
		panic(err)
	}

	testCert := eigenda.Cert{
		BlobCommitment: &common.G1Commitment{
			X: x,
			Y: y,
		},
	}

	bytes, err := rlp.EncodeToBytes(testCert)
	assert.NoError(t, err)

	var testCommit EigenDACommitment = bytes

	var data = []byte("inter-subjective and not objective!")

	err = testCommit.Verify(data)
	assert.NoError(t, err)
}

func ConvertIntToVarUInt(v int) []byte {
	buf := make([]byte, binary.MaxVarintLen64)
	n := binary.PutUvarint(buf, uint64(v))
	return buf[:n]

}
