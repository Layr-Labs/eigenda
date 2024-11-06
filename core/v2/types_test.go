package v2_test

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/Layr-Labs/eigenda/core"
	v2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
	"github.com/stretchr/testify/assert"
)

func TestBlobKey(t *testing.T) {
	blobKey := v2.BlobKey([32]byte{1, 2, 3})

	assert.Equal(t, "0102030000000000000000000000000000000000000000000000000000000000", blobKey.Hex())
	bk, err := v2.HexToBlobKey(blobKey.Hex())
	assert.NoError(t, err)
	assert.Equal(t, blobKey, bk)
}

func TestPaymentHash(t *testing.T) {
	pm := core.PaymentMetadata{
		AccountID:         "0x123",
		BinIndex:          5,
		CumulativePayment: big.NewInt(100),
	}
	hash, err := pm.Hash()
	assert.NoError(t, err)
	// 0xf5894a8e9281b5687c0c7757d3d45fb76152bf659e6e61b1062f4c6bcb69c449 verified in solidity
	assert.Equal(t, "f5894a8e9281b5687c0c7757d3d45fb76152bf659e6e61b1062f4c6bcb69c449", hex.EncodeToString(hash[:]))
}

func TestBlobKeyFromHeader(t *testing.T) {
	data := codec.ConvertByPaddingEmptyByte(GETTYSBURG_ADDRESS_BYTES)
	commitments, err := p.GetCommitments(data)
	if err != nil {
		t.Fatal(err)
	}

	bh := v2.BlobHeader{
		BlobVersion:     0,
		BlobCommitments: commitments,
		QuorumNumbers:   []core.QuorumID{0, 1},
		PaymentMetadata: core.PaymentMetadata{
			AccountID:         "0x123",
			BinIndex:          5,
			CumulativePayment: big.NewInt(100),
		},
		Signature: []byte{1, 2, 3},
	}
	blobKey, err := bh.BlobKey()
	assert.NoError(t, err)
	// 0xb19d368345990c79744fe571fe99f427f35787b9383c55089fb5bd6a5c171bbc verified in solidity
	assert.Equal(t, "b19d368345990c79744fe571fe99f427f35787b9383c55089fb5bd6a5c171bbc", blobKey.Hex())
}
