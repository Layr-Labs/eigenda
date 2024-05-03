package eigenda

import (
	"testing"

	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/assert"
)

func equalSlices[P comparable](s1, s2 []P) bool {
	if len(s1) != len(s2) {
		return false
	}

	for i := 0; i < len(s1); i++ {
		if s1[i] != s2[i] {
			return false
		}
	}

	return true
}

func TestCertEncodingDecoding(t *testing.T) {
	c := Cert{
		BatchHeaderHash:      []byte{0x42, 0x69},
		BlobIndex:            420,
		ReferenceBlockNumber: 80085,
		QuorumIDs:            []uint32{666},
	}

	bytes, err := rlp.EncodeToBytes(c)
	assert.NoError(t, err, "encoding should pass")

	var c2 *Cert
	err = rlp.DecodeBytes(bytes, &c2)
	assert.NoError(t, err, "decoding should pass")

	equal := func() bool {
		return equalSlices(c.BatchHeaderHash, c2.BatchHeaderHash) &&
			c.BlobIndex == c2.BlobIndex &&
			c.ReferenceBlockNumber == c2.ReferenceBlockNumber &&
			equalSlices(c.QuorumIDs, c2.QuorumIDs)
	}

	assert.True(t, equal(), "values shouldn't change")
}
