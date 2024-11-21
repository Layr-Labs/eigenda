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
	commitments, err := p.GetCommitmentsForPaddedLength(data)
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
	// 0x40efb7273649f39590b27550ea06eeb81efd6ae4d719385a302fbd93173a395d verified in solidity
	assert.Equal(t, "40efb7273649f39590b27550ea06eeb81efd6ae4d719385a302fbd93173a395d", blobKey.Hex())
}

func TestBatchHeaderHash(t *testing.T) {
	batchRoot := [32]byte{}
	copy(batchRoot[:], []byte("1"))
	batchHeader := &v2.BatchHeader{
		ReferenceBlockNumber: 1,
		BatchRoot:            batchRoot,
	}

	hash, err := batchHeader.Hash()
	assert.NoError(t, err)
	// 0x891d0936da4627f445ef193aad63afb173409af9e775e292e4e35aff790a45e2 verified in solidity
	assert.Equal(t, "891d0936da4627f445ef193aad63afb173409af9e775e292e4e35aff790a45e2", hex.EncodeToString(hash[:]))
}

func TestBatchHeaderSerialization(t *testing.T) {
	batchRoot := [32]byte{}
	copy(batchRoot[:], []byte("batchRoot"))
	batchHeader := &v2.BatchHeader{
		ReferenceBlockNumber: 1000,
		BatchRoot:            batchRoot,
	}

	serialized, err := batchHeader.Serialize()
	assert.NoError(t, err)
	deserialized, err := v2.DeserializeBatchHeader(serialized)
	assert.NoError(t, err)
	assert.Equal(t, batchHeader, deserialized)
}

func TestBlobCertHash(t *testing.T) {
	data := codec.ConvertByPaddingEmptyByte(GETTYSBURG_ADDRESS_BYTES)
	commitments, err := p.GetCommitmentsForPaddedLength(data)
	if err != nil {
		t.Fatal(err)
	}

	blobCert := &v2.BlobCertificate{
		BlobHeader: &v2.BlobHeader{
			BlobVersion:     0,
			BlobCommitments: commitments,
			QuorumNumbers:   []core.QuorumID{0, 1},
			PaymentMetadata: core.PaymentMetadata{
				AccountID:         "0x123",
				BinIndex:          5,
				CumulativePayment: big.NewInt(100),
			},
			Signature: []byte{1, 2, 3},
		},
		RelayKeys: []v2.RelayKey{4, 5, 6},
	}

	hash, err := blobCert.Hash()
	assert.NoError(t, err)
	// 0x3719a91e2a294feafdd624c1c88a6f1db1a5c79ee0863b352255bc9162f02751 verified in solidity
	assert.Equal(t, "3719a91e2a294feafdd624c1c88a6f1db1a5c79ee0863b352255bc9162f02751", hex.EncodeToString(hash[:]))
}

func TestBlobCertSerialization(t *testing.T) {
	data := codec.ConvertByPaddingEmptyByte(GETTYSBURG_ADDRESS_BYTES)
	commitments, err := p.GetCommitmentsForPaddedLength(data)
	if err != nil {
		t.Fatal(err)
	}

	blobCert := &v2.BlobCertificate{
		BlobHeader: &v2.BlobHeader{
			BlobVersion:     0,
			BlobCommitments: commitments,
			QuorumNumbers:   []core.QuorumID{0, 1},
			PaymentMetadata: core.PaymentMetadata{
				AccountID:         "0x123",
				BinIndex:          5,
				CumulativePayment: big.NewInt(100),
			},
			Signature: []byte{1, 2, 3},
		},
		RelayKeys: []v2.RelayKey{4, 5, 6},
	}

	serialized, err := blobCert.Serialize()
	assert.NoError(t, err)
	deserialized, err := v2.DeserializeBlobCertificate(serialized)
	assert.NoError(t, err)
	assert.Equal(t, blobCert, deserialized)
}
