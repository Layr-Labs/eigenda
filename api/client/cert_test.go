package client

import (
	"testing"

	"github.com/Layr-Labs/eigenda/api/grpc/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/assert"
)

func TestCertEncodingDecoding(t *testing.T) {
	c := Cert{
		BatchHeaderHash:      []byte{0x42, 0x69},
		BlobIndex:            420,
		ReferenceBlockNumber: 80085,
		QuorumIDs:            []uint32{666},
		BlobCommitment: &common.G1Commitment{
			X: []byte{0x1},
			Y: []byte{0x3},
		},
	}

	bytes, err := rlp.EncodeToBytes(c)
	assert.NoError(t, err, "encoding should pass")

	var c2 *Cert
	err = rlp.DecodeBytes(bytes, &c2)
	assert.NoError(t, err, "decoding should pass")

	assert.Equal(t, c.BatchHeaderHash, c2.BatchHeaderHash)
	assert.Equal(t, c.BlobIndex, c2.BlobIndex)
	assert.Equal(t, c.ReferenceBlockNumber, c2.ReferenceBlockNumber)
	assert.Equal(t, c.QuorumIDs, c2.QuorumIDs)
	assert.Equal(t, c.BlobCommitment.X, c2.BlobCommitment.X)
	assert.Equal(t, c.BlobCommitment.Y, c2.BlobCommitment.Y)
}
