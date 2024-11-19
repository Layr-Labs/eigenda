package v2_test

import (
	"math/big"
	"testing"

	"github.com/Layr-Labs/eigenda/core"
	v2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
	"github.com/stretchr/testify/assert"
)

func TestConvertBatchToFromProtobuf(t *testing.T) {
	data := codec.ConvertByPaddingEmptyByte(GETTYSBURG_ADDRESS_BYTES)
	commitments, err := p.GetCommitments(data)
	if err != nil {
		t.Fatal(err)
	}

	bh0 := &v2.BlobHeader{
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
	bh1 := &v2.BlobHeader{
		BlobVersion:     0,
		BlobCommitments: commitments,
		QuorumNumbers:   []core.QuorumID{0, 1},
		PaymentMetadata: core.PaymentMetadata{
			AccountID:         "0x456",
			BinIndex:          6,
			CumulativePayment: big.NewInt(200),
		},
		Signature: []byte{1, 2, 3},
	}

	blobCert0 := &v2.BlobCertificate{
		BlobHeader: bh0,
		RelayKeys:  []v2.RelayKey{0, 1},
	}
	blobCert1 := &v2.BlobCertificate{
		BlobHeader: bh1,
		RelayKeys:  []v2.RelayKey{2, 3},
	}

	batch := &v2.Batch{
		BatchHeader: &v2.BatchHeader{
			BatchRoot:            [32]byte{1, 1, 1},
			ReferenceBlockNumber: 100,
		},
		BlobCertificates: []*v2.BlobCertificate{blobCert0, blobCert1},
	}

	pb, err := batch.ToProtobuf()
	assert.NoError(t, err)

	newBatch, err := v2.BatchFromProtobuf(pb)
	assert.NoError(t, err)

	assert.Equal(t, batch, newBatch)
}

func TestConvertBlobHeaderToFromProtobuf(t *testing.T) {
	data := codec.ConvertByPaddingEmptyByte(GETTYSBURG_ADDRESS_BYTES)
	commitments, err := p.GetCommitments(data)
	if err != nil {
		t.Fatal(err)
	}

	bh := &v2.BlobHeader{
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

	pb, err := bh.ToProtobuf()
	assert.NoError(t, err)

	newBH, err := v2.BlobHeaderFromProtobuf(pb)
	assert.NoError(t, err)

	assert.Equal(t, bh, newBH)
}

func TestConvertBlobCertToFromProtobuf(t *testing.T) {
	data := codec.ConvertByPaddingEmptyByte(GETTYSBURG_ADDRESS_BYTES)
	commitments, err := p.GetCommitments(data)
	if err != nil {
		t.Fatal(err)
	}

	bh := &v2.BlobHeader{
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

	blobCert := &v2.BlobCertificate{
		BlobHeader: bh,
		RelayKeys:  []v2.RelayKey{0, 1},
	}

	pb, err := blobCert.ToProtobuf()
	assert.NoError(t, err)

	newBlobCert, err := v2.BlobCertificateFromProtobuf(pb)
	assert.NoError(t, err)

	assert.Equal(t, blobCert, newBlobCert)
}
