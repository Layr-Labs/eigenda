package hashing

import (
	"testing"

	commonv1 "github.com/Layr-Labs/eigenda/api/grpc/common"
	commonv2 "github.com/Layr-Labs/eigenda/api/grpc/common/v2"
	grpc "github.com/Layr-Labs/eigenda/api/grpc/validator"
	"github.com/Layr-Labs/eigenda/api/hashing"
	"github.com/stretchr/testify/require"
)

func TestHashStoreChunksRequest_CanonicalMatchesHasherImplementation(t *testing.T) {
	req := &grpc.StoreChunksRequest{
		Batch: &commonv2.Batch{
			Header: &commonv2.BatchHeader{
				BatchRoot:            []byte{0x01, 0x02, 0x03, 0x04},
				ReferenceBlockNumber: 42,
			},
			BlobCertificates: []*commonv2.BlobCertificate{
				{
					BlobHeader: &commonv2.BlobHeader{
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
							Timestamp:         999,
							CumulativePayment: []byte{0x09, 0x08, 0x07},
						},
					},
					Signature: []byte{0xde, 0xad, 0xbe, 0xef},
					RelayKeys: []uint32{5, 6},
				},
				{
					BlobHeader: &commonv2.BlobHeader{
						Version:       2,
						QuorumNumbers: []uint32{1},
						Commitment: &commonv1.BlobCommitment{
							Commitment:       []byte{0x01},
							LengthCommitment: []byte{},
							LengthProof:      []byte{0xff, 0xee},
							Length:           0,
						},
						PaymentHeader: &commonv2.PaymentHeader{
							AccountId:         "0xdef",
							Timestamp:         123456789,
							CumulativePayment: []byte{0x01, 0x00},
						},
					},
					Signature: []byte{0x00},
					RelayKeys: []uint32{0},
				},
			},
		},
		DisperserID: 7,
		Timestamp:   55,
	}

	h1, err := hashing.HashStoreChunksRequest(req)
	require.NoError(t, err)

	h2, err := HashStoreChunksRequest_Canonical(req)
	require.NoError(t, err)

	h3, err := HashStoreChunksRequest_V2_Canonical(req)
	require.NoError(t, err)

	require.Equal(t, h1, h2, "canonical (manual) serializer hash must match HashStoreChunksRequest")
	require.Equal(t, h1, h3, "canonical (struc) serializer hash must match HashStoreChunksRequest")
}
