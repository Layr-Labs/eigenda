package v2_test

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/Layr-Labs/eigenda/core"
	v2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
	gethcommon "github.com/ethereum/go-ethereum/common"
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
		AccountID:         gethcommon.HexToAddress("0x0000000000000000000000000000000000000123"),
		Timestamp:         5,
		CumulativePayment: big.NewInt(100),
	}
	hash, err := pm.Hash()
	assert.NoError(t, err)
	// 234c3d10881641264afe33cf492000f8ecd505e385050314c63469c3ad2977c9 verified in solidity
	assert.Equal(t, "234c3d10881641264afe33cf492000f8ecd505e385050314c63469c3ad2977c9", hex.EncodeToString(hash[:]))
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
			AccountID:         gethcommon.HexToAddress("0x0000000000000000000000000000000000000123"),
			Timestamp:         5,
			CumulativePayment: big.NewInt(100),
		},
	}
	blobKey, err := bh.BlobKey()
	assert.NoError(t, err)
	// e2fc52cb6213041838c20164eac05a7660b741518d5c14060e47c89ed3dd175b has verified in solidity  with chisel
	assert.Equal(t, "e2fc52cb6213041838c20164eac05a7660b741518d5c14060e47c89ed3dd175b", blobKey.Hex())

	// same blob key should be generated for the blob header with shuffled quorum numbers
	bh2 := v2.BlobHeader{
		BlobVersion:     0,
		BlobCommitments: commitments,
		QuorumNumbers:   []core.QuorumID{1, 0},
		PaymentMetadata: core.PaymentMetadata{
			AccountID:         gethcommon.HexToAddress("0x0000000000000000000000000000000000000123"),
			Timestamp:         5,
			CumulativePayment: big.NewInt(100),
		},
	}

	blobKey2, err := bh2.BlobKey()
	assert.NoError(t, err)
	assert.Equal(t, blobKey2.Hex(), blobKey.Hex())
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
				AccountID:         gethcommon.HexToAddress("0x0000000000000000000000000000000000000123"),
				Timestamp:         5,
				CumulativePayment: big.NewInt(100),
			},
		},
		Signature: []byte{1, 2, 3},
		RelayKeys: []v2.RelayKey{4, 5, 6},
	}

	hash, err := blobCert.Hash()
	assert.NoError(t, err)

	// 932dd5724ce1d7ecd076bd8e7423562005701053b3751cf65a7dd8d25e737484 has verified in solidity with chisel
	assert.Equal(t, "932dd5724ce1d7ecd076bd8e7423562005701053b3751cf65a7dd8d25e737484", hex.EncodeToString(hash[:]))
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
				AccountID:         gethcommon.HexToAddress("0x0000000000000000000000000000000000000123"),
				Timestamp:         5,
				CumulativePayment: big.NewInt(100),
			},
		},
		Signature: []byte{1, 2, 3},
		RelayKeys: []v2.RelayKey{4, 5, 6},
	}

	serialized, err := blobCert.Serialize()
	assert.NoError(t, err)
	deserialized, err := v2.DeserializeBlobCertificate(serialized)
	assert.NoError(t, err)
	assert.Equal(t, blobCert, deserialized)
}
