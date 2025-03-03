package hashing

import (
	"fmt"
	"hash"

	commonv1 "github.com/Layr-Labs/eigenda/api/grpc/common"
	common "github.com/Layr-Labs/eigenda/api/grpc/common/v2"
	grpc "github.com/Layr-Labs/eigenda/api/grpc/validator"
	"golang.org/x/crypto/sha3"
)

// This file contains code for hashing gRPC messages that are sent to the DA node.

// HashStoreChunksRequest hashes the given StoreChunksRequest.
func HashStoreChunksRequest(request *grpc.StoreChunksRequest) ([]byte, error) {
	hasher := sha3.NewLegacyKeccak256()

	err := hashBatchHeader(hasher, request.Batch.Header)
	if err != nil {
		return nil, fmt.Errorf("failed to hash batch header: %w", err)
	}
	err = hashLength(hasher, request.Batch.BlobCertificates)
	if err != nil {
		return nil, fmt.Errorf("failed to hash BlobCertificates length: %w", err)
	}
	for _, blobCertificate := range request.Batch.BlobCertificates {
		err = hashBlobCertificate(hasher, blobCertificate)
		if err != nil {
			return nil, fmt.Errorf("failed to hash blob certificate: %w", err)
		}
	}
	hashUint32(hasher, request.DisperserID)

	return hasher.Sum(nil), nil
}

func hashBlobCertificate(hasher hash.Hash, blobCertificate *common.BlobCertificate) error {
	err := hashBlobHeader(hasher, blobCertificate.BlobHeader)
	if err != nil {
		return fmt.Errorf("failed to hash blob header: %w", err)
	}
	err = hashByteArray(hasher, blobCertificate.Signature)
	if err != nil {
		return fmt.Errorf("failed to hash signature: %w", err)
	}
	err = hashUint32Array(hasher, blobCertificate.RelayKeys)
	if err != nil {
		return fmt.Errorf("failed to hash RelayKeys: %w", err)
	}
	return nil
}

func hashBlobHeader(hasher hash.Hash, header *common.BlobHeader) error {
	hashUint32(hasher, header.Version)
	hashUint32(hasher, uint32(len(header.QuorumNumbers)))

	err := hashUint32Array(hasher, header.QuorumNumbers)
	if err != nil {
		return fmt.Errorf("failed to hash QuorumNumbers: %w", err)
	}

	err = hashBlobCommitment(hasher, header.Commitment)
	if err != nil {
		return fmt.Errorf("failed to hash commitment: %w", err)
	}

	err = hashPaymentHeader(hasher, header.PaymentHeader)
	if err != nil {
		return fmt.Errorf("failed to hash payment header: %w", err)
	}

	return nil
}

func hashBatchHeader(hasher hash.Hash, header *common.BatchHeader) error {
	err := hashByteArray(hasher, header.BatchRoot)
	if err != nil {
		return fmt.Errorf("failed to hash BatchRoot: %w", err)
	}
	hashUint64(hasher, header.ReferenceBlockNumber)

	return nil
}

func hashBlobCommitment(hasher hash.Hash, commitment *commonv1.BlobCommitment) error {
	err := hashByteArray(hasher, commitment.Commitment)
	if err != nil {
		return fmt.Errorf("failed to hash commitment: %w", err)
	}

	err = hashByteArray(hasher, commitment.LengthCommitment)
	if err != nil {
		return fmt.Errorf("failed to hash LengthCommitment: %w", err)
	}

	err = hashByteArray(hasher, commitment.LengthProof)
	if err != nil {
		return fmt.Errorf("failed to hash LengthProof: %w", err)
	}

	hashUint32(hasher, commitment.Length)

	return nil
}

func hashPaymentHeader(hasher hash.Hash, header *common.PaymentHeader) error {
	err := hashByteArray(hasher, []byte(header.AccountId))
	if err != nil {
		return fmt.Errorf("failed to hash AccountId: %w", err)
	}

	hashInt64(hasher, header.Timestamp)

	err = hashByteArray(hasher, header.CumulativePayment)
	if err != nil {
		return fmt.Errorf("failed to hash CumulativePayment: %w", err)
	}

	return nil
}
