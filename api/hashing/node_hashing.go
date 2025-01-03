package hashing

import (
	commonv1 "github.com/Layr-Labs/eigenda/api/grpc/common"
	common "github.com/Layr-Labs/eigenda/api/grpc/common/v2"
	grpc "github.com/Layr-Labs/eigenda/api/grpc/node/v2"
	"golang.org/x/crypto/sha3"
	"hash"
)

// This file contains code for hashing gRPC messages that are sent to the DA node.

// HashStoreChunksRequest hashes the given StoreChunksRequest.
func HashStoreChunksRequest(request *grpc.StoreChunksRequest) []byte {
	hasher := sha3.NewLegacyKeccak256()

	hashBatchHeader(hasher, request.Batch.Header)
	for _, blobCertificate := range request.Batch.BlobCertificates {
		hashBlobCertificate(hasher, blobCertificate)
	}
	hashUint32(hasher, request.DisperserID)

	return hasher.Sum(nil)
}

func hashBlobCertificate(hasher hash.Hash, blobCertificate *common.BlobCertificate) {
	hashBlobHeader(hasher, blobCertificate.BlobHeader)
	for _, relayID := range blobCertificate.Relays {
		hashUint32(hasher, relayID)
	}
}

func hashBlobHeader(hasher hash.Hash, header *common.BlobHeader) {
	hashUint32(hasher, header.Version)
	for _, quorum := range header.QuorumNumbers {
		hashUint32(hasher, quorum)
	}
	hashBlobCommitment(hasher, header.Commitment)
	hashPaymentHeader(hasher, header.PaymentHeader)
	hasher.Write(header.Signature)
}

func hashBatchHeader(hasher hash.Hash, header *common.BatchHeader) {
	hasher.Write(header.BatchRoot)
	hashUint64(hasher, header.ReferenceBlockNumber)
}

func hashBlobCommitment(hasher hash.Hash, commitment *commonv1.BlobCommitment) {
	hasher.Write(commitment.Commitment)
	hasher.Write(commitment.LengthCommitment)
	hasher.Write(commitment.LengthProof)
	hashUint32(hasher, commitment.Length)
}

func hashPaymentHeader(hasher hash.Hash, header *commonv1.PaymentHeader) {
	hasher.Write([]byte(header.AccountId))
	hashUint32(hasher, header.ReservationPeriod)
	hasher.Write(header.CumulativePayment)
	hashUint32(hasher, header.Salt)
}
