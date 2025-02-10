package auth

import (
	"github.com/Layr-Labs/eigenda/api/grpc/common"
	v2 "github.com/Layr-Labs/eigenda/api/grpc/common/v2"
	grpc "github.com/Layr-Labs/eigenda/api/grpc/validator"
	"github.com/Layr-Labs/eigenda/common/testutils/random"
)

func RandomStoreChunksRequest(rand *random.TestRandom) *grpc.StoreChunksRequest {
	certificateCount := rand.Intn(10) + 1
	blobCertificates := make([]*v2.BlobCertificate, certificateCount)
	for i := 0; i < certificateCount; i++ {

		relayCount := rand.Intn(10) + 1
		relays := make([]uint32, relayCount)
		for j := 0; j < relayCount; j++ {
			relays[j] = rand.Uint32()
		}

		quorumCount := rand.Intn(10) + 1
		quorumNumbers := make([]uint32, quorumCount)
		for j := 0; j < quorumCount; j++ {
			quorumNumbers[j] = rand.Uint32()
		}

		blobCertificates[i] = &v2.BlobCertificate{
			BlobHeader: &v2.BlobHeader{
				Version:       rand.Uint32(),
				QuorumNumbers: quorumNumbers,
				Commitment: &common.BlobCommitment{
					Commitment:       rand.Bytes(32),
					LengthCommitment: rand.Bytes(32),
					LengthProof:      rand.Bytes(32),
					Length:           rand.Uint32(),
				},
				PaymentHeader: &v2.PaymentHeader{
					AccountId:         rand.String(32),
					Timestamp:         rand.Uint64(),
					CumulativePayment: rand.Bytes(32),
				},
			},
			Signature: rand.Bytes(32),
			RelayKeys: relays,
		}
	}

	return &grpc.StoreChunksRequest{
		Batch: &v2.Batch{
			Header: &v2.BatchHeader{
				BatchRoot:            rand.Bytes(32),
				ReferenceBlockNumber: rand.Uint64(),
			},
			BlobCertificates: blobCertificates,
		},
		DisperserID: rand.Uint32(),
		Signature:   rand.Bytes(32),
	}
}
