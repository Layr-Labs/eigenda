package hashing

import (
	"testing"
	"time"

	commonv1 "github.com/Layr-Labs/eigenda/api/grpc/common"
	commonv2 "github.com/Layr-Labs/eigenda/api/grpc/common/v2"
	grpc "github.com/Layr-Labs/eigenda/api/grpc/validator"
	hashingv2 "github.com/Layr-Labs/eigenda/api/hashing/v2"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/sha3"
)

func TestV2BlobHeaderHashMatchesLegacyHashBlobHeader(t *testing.T) {
	tsNanos := int64(1234567890123456)
	header := &commonv2.BlobHeader{
		Version:       1,
		QuorumNumbers: []uint32{0, 2, 7},
		Commitment: &commonv1.BlobCommitment{
			Commitment:       []byte{0xaa, 0xbb},
			LengthCommitment: []byte{0x10, 0x11, 0x12},
			LengthProof:      []byte{0x20},
			Length:           123,
		},
		PaymentHeader: &commonv2.PaymentHeader{
			AccountId:         "0xabc",
			Timestamp:         tsNanos,
			CumulativePayment: []byte{0x09, 0x08, 0x07},
		},
	}

	// "Old" blob header hash: use the legacy streaming serializer+hash used by node_hashing.go.
	legacyHasher := sha3.NewLegacyKeccak256()
	err := hashBlobHeader(legacyHasher, header)
	require.NoError(t, err)
	legacyHash := legacyHasher.Sum(nil)

	req := &grpc.StoreChunksRequest{
		Batch: &commonv2.Batch{
			Header: &commonv2.BatchHeader{BatchRoot: []byte{0x01}, ReferenceBlockNumber: 1},
			BlobCertificates: []*commonv2.BlobCertificate{{
				BlobHeader: header,
				Signature:  []byte{0x00},
				RelayKeys:  []uint32{1},
			}},
		},
		DisperserID: 1,
		Timestamp:   1,
	}

	got, err := hashingv2.BlobHeadersHashesAndTimestamps(req)
	require.NoError(t, err)
	require.Len(t, got, 1)

	require.Equal(t, legacyHash, got[0].Hash)
	require.True(t, got[0].Timestamp.Equal(time.Unix(0, tsNanos)))
}
