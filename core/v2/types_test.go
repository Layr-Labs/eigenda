package v2_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/core"
	v2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConvertBatchToFromProtobuf(t *testing.T) {
	data := codec.ConvertByPaddingEmptyByte(GETTYSBURG_ADDRESS_BYTES)
	commitments, err := p.GetCommitmentsForPaddedLength(data)
	require.NoError(t, err)

	bh0 := &v2.BlobHeader{
		BlobVersion:     0,
		BlobCommitments: commitments,
		QuorumNumbers:   []core.QuorumID{0, 1},
		PaymentMetadata: core.PaymentMetadata{
			AccountID:         gethcommon.HexToAddress("0x123"),
			Timestamp:         5,
			CumulativePayment: big.NewInt(100),
		},
	}
	bh1 := &v2.BlobHeader{
		BlobVersion:     0,
		BlobCommitments: commitments,
		QuorumNumbers:   []core.QuorumID{0, 1},
		PaymentMetadata: core.PaymentMetadata{
			AccountID:         gethcommon.HexToAddress("0x456"),
			Timestamp:         6,
			CumulativePayment: big.NewInt(200),
		},
	}

	blobCert0 := &v2.BlobCertificate{
		BlobHeader: bh0,
		Signature:  []byte{1, 2, 3},
		RelayKeys:  []v2.RelayKey{0, 1},
	}
	blobCert1 := &v2.BlobCertificate{
		BlobHeader: bh1,
		Signature:  []byte{1, 2, 3},
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
	commitments, err := p.GetCommitmentsForPaddedLength(data)
	require.NoError(t, err)

	bh := &v2.BlobHeader{
		BlobVersion:     0,
		BlobCommitments: commitments,
		QuorumNumbers:   []core.QuorumID{0, 1},
		PaymentMetadata: core.PaymentMetadata{
			AccountID:         gethcommon.HexToAddress("0x123"),
			Timestamp:         5,
			CumulativePayment: big.NewInt(100),
		},
	}

	pb, err := bh.ToProtobuf()
	assert.NoError(t, err)

	newBH, err := v2.BlobHeaderFromProtobuf(pb)
	assert.NoError(t, err)

	assert.Equal(t, bh, newBH)
}

func TestConvertBlobCertToFromProtobuf(t *testing.T) {
	data := codec.ConvertByPaddingEmptyByte(GETTYSBURG_ADDRESS_BYTES)
	commitments, err := p.GetCommitmentsForPaddedLength(data)
	require.NoError(t, err)

	bh := &v2.BlobHeader{
		BlobVersion:     0,
		BlobCommitments: commitments,
		QuorumNumbers:   []core.QuorumID{0, 1},
		PaymentMetadata: core.PaymentMetadata{
			AccountID:         gethcommon.HexToAddress("0x123"),
			Timestamp:         5,
			CumulativePayment: big.NewInt(100),
		},
	}

	blobCert := &v2.BlobCertificate{
		BlobHeader: bh,
		Signature:  []byte{1, 2, 3},
		RelayKeys:  []v2.RelayKey{0, 1},
	}

	pb, err := blobCert.ToProtobuf()
	assert.NoError(t, err)

	newBlobCert, err := v2.BlobCertificateFromProtobuf(pb)
	assert.NoError(t, err)

	assert.Equal(t, blobCert, newBlobCert)
}

func TestAttestationToProtobuf(t *testing.T) {
	zeroAttestation := &v2.Attestation{
		BatchHeader: &v2.BatchHeader{
			BatchRoot:            [32]byte{1, 1, 1},
			ReferenceBlockNumber: 100,
		},
		AttestedAt:       uint64(time.Now().UnixNano()),
		NonSignerPubKeys: nil,
		APKG2:            nil,
		QuorumAPKs:       nil,
		Sigma:            nil,
		QuorumNumbers:    nil,
		QuorumResults:    nil,
	}
	attestationProto, err := zeroAttestation.ToProtobuf()
	assert.NoError(t, err)
	assert.Empty(t, attestationProto.NonSignerPubkeys)
	assert.Empty(t, attestationProto.ApkG2)
	assert.Empty(t, attestationProto.QuorumApks)
	assert.Empty(t, attestationProto.Sigma)
	assert.Empty(t, attestationProto.QuorumNumbers)
	assert.Empty(t, attestationProto.QuorumSignedPercentages)
}
