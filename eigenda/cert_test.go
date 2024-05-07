package eigenda

import (
	"testing"

	eigen_da_common "github.com/Layr-Labs/eigenda/api/grpc/common"
	"github.com/Layr-Labs/op-plasma-eigenda/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/assert"
)

func TestCertEncodingDecoding(t *testing.T) {
	c := Cert{
		BatchHeaderHash:      []byte{0x42, 0x69},
		BlobIndex:            420,
		ReferenceBlockNumber: 80085,
		QuorumIDs:            []uint32{666},
		BlobCommitment: &eigen_da_common.G1Commitment{
			X: []byte{0x1},
			Y: []byte{0x3},
		},
	}

	bytes, err := rlp.EncodeToBytes(c)
	assert.NoError(t, err, "encoding should pass")

	var c2 *Cert
	err = rlp.DecodeBytes(bytes, &c2)
	assert.NoError(t, err, "decoding should pass")

	equal := func() bool {
		return common.EqualSlices(c.BatchHeaderHash, c2.BatchHeaderHash) &&
			c.BlobIndex == c2.BlobIndex &&
			c.ReferenceBlockNumber == c2.ReferenceBlockNumber &&
			common.EqualSlices(c.QuorumIDs, c2.QuorumIDs) &&
			common.EqualSlices(c.BlobCommitment.X, c2.BlobCommitment.X) &&
			common.EqualSlices(c.BlobCommitment.Y, c2.BlobCommitment.Y)
	}

	assert.True(t, equal(), "values shouldn't change")
}
