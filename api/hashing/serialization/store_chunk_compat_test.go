package serialization

import (
	"bytes"
	"testing"

	commonv1 "github.com/Layr-Labs/eigenda/api/grpc/common"
	commonv2 "github.com/Layr-Labs/eigenda/api/grpc/common/v2"
	grpc "github.com/Layr-Labs/eigenda/api/grpc/validator"
)

func TestSerializeStoreChunksRequest_V1MatchesV2(t *testing.T) {
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

	b1, err := SerializeStoreChunksRequest(req)
	if err != nil {
		t.Fatalf("SerializeStoreChunksRequest: %v", err)
	}
	b2, err := SerializeStoreChunksRequestV2(req)
	if err != nil {
		t.Fatalf("SerializeStoreChunksRequestV2: %v", err)
	}

	if !bytes.Equal(b1, b2) {
		t.Fatalf("serialization mismatch\nv1=%x\nv2=%x", b1, b2)
	}
}

func TestSerializeBlobHeader_V1MatchesV2(t *testing.T) {
	hdr := &commonv2.BlobHeader{
		Version:       3,
		QuorumNumbers: []uint32{9, 8},
		Commitment: &commonv1.BlobCommitment{
			Commitment:       []byte{0x01, 0x02, 0x03},
			LengthCommitment: []byte{0x99},
			LengthProof:      []byte{0x88},
			Length:           777,
		},
		PaymentHeader: &commonv2.PaymentHeader{
			AccountId:         "0x1234",
			Timestamp:         123,
			CumulativePayment: []byte{0x01},
		},
	}

	b1, err := SerializeBlobHeader(hdr)
	if err != nil {
		t.Fatalf("SerializeBlobHeader: %v", err)
	}
	b2, err := SerializeBlobHeaderV2(hdr)
	if err != nil {
		t.Fatalf("SerializeBlobHeaderV2: %v", err)
	}

	if !bytes.Equal(b1, b2) {
		t.Fatalf("blob header serialization mismatch\nv1=%x\nv2=%x", b1, b2)
	}
}
