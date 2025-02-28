package hashing

import (
	"hash"

	commonv1 "github.com/Layr-Labs/eigenda/api/grpc/common"
	common "github.com/Layr-Labs/eigenda/api/grpc/common/v2"
	grpc "github.com/Layr-Labs/eigenda/api/grpc/validator"
	"golang.org/x/crypto/sha3"
)

// This file contains code for hashing gRPC messages that are sent to the DA node.

// HashStoreChunksRequest hashes the given StoreChunksRequest.
func HashStoreChunksRequest(request *grpc.StoreChunksRequest) []byte {
	hasher := sha3.NewLegacyKeccak256()

	hashBatchHeader(hasher, request.GetBatch().GetHeader())
	for _, blobCertificate := range request.GetBatch().GetBlobCertificates() {
		hashBlobCertificate(hasher, blobCertificate)
	}
	hashUint32(hasher, request.GetDisperserID())

	return hasher.Sum(nil)
}

func hashBlobCertificate(hasher hash.Hash, blobCertificate *common.BlobCertificate) {
	hashBlobHeader(hasher, blobCertificate.GetBlobHeader())
	hasher.Write(blobCertificate.GetSignature())
	for _, relayKey := range blobCertificate.GetRelayKeys() {
		hashUint32(hasher, relayKey)
	}
}

func hashBlobHeader(hasher hash.Hash, header *common.BlobHeader) {
	hashUint32(hasher, header.GetVersion())
	for _, quorum := range header.GetQuorumNumbers() {
		hashUint32(hasher, quorum)
	}
	hashBlobCommitment(hasher, header.GetCommitment())
	hashPaymentHeader(hasher, header.GetPaymentHeader())
}

func hashBatchHeader(hasher hash.Hash, header *common.BatchHeader) {
	hasher.Write(header.GetBatchRoot())
	hashUint64(hasher, header.GetReferenceBlockNumber())
}

func hashBlobCommitment(hasher hash.Hash, commitment *commonv1.BlobCommitment) {
	hasher.Write(commitment.GetCommitment())
	hasher.Write(commitment.GetLengthCommitment())
	hasher.Write(commitment.GetLengthProof())
	hashUint32(hasher, commitment.GetLength())
}

func hashPaymentHeader(hasher hash.Hash, header *common.PaymentHeader) {
	hasher.Write([]byte(header.GetAccountId()))
	hashInt64(hasher, header.GetTimestamp())
	hasher.Write(header.GetCumulativePayment())
}
