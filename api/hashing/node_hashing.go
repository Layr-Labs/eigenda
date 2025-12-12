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

// ValidatorStoreChunksRequestDomain is the domain for hashing StoreChunksRequest messages (i.e. this string
// is added to the digest before hashing the message). This makes it difficult for an attacker to create a
// different type of object that has the same hash as a StoreChunksRequest.
const ValidatorStoreChunksRequestDomain = "validator.StoreChunksRequest"

// HashStoreChunksRequest hashes the given StoreChunksRequest.
func HashStoreChunksRequest(request *grpc.StoreChunksRequest) ([]byte, error) {
	hasher := sha3.New256()

	hasher.Write([]byte(ValidatorStoreChunksRequestDomain))

	err := hashBatchHeader(hasher, request.GetBatch().GetHeader())
	if err != nil {
		return nil, fmt.Errorf("failed to hash batch header: %w", err)
	}
	err = hashLength(hasher, request.GetBatch().GetBlobCertificates())
	if err != nil {
		return nil, fmt.Errorf("failed to hash BlobCertificates length: %w", err)
	}
	for _, blobCertificate := range request.GetBatch().GetBlobCertificates() {
		err = hashBlobCertificate(hasher, blobCertificate)
		if err != nil {
			return nil, fmt.Errorf("failed to hash blob certificate: %w", err)
		}
	}
	hashUint32(hasher, request.GetDisperserID())
	hashUint32(hasher, request.GetTimestamp())

	return hasher.Sum(nil), nil
}

// HashStoreChunksRequestBlobHeaders returns a list of per-BlobHeader hashes (one per BlobCertificate).
func HashStoreChunksRequestBlobHeaders(request *grpc.StoreChunksRequest) ([][]byte, error) {
	certs := request.GetBatch().GetBlobCertificates()
	blobHeaderHashes := make([][]byte, len(certs))
	for i, cert := range certs {
		if cert == nil {
			return nil, fmt.Errorf("nil BlobCertificate at index %d", i)
		}
		header := cert.GetBlobHeader()
		if header == nil {
			return nil, fmt.Errorf("nil BlobHeader at index %d", i)
		}

		h := sha3.New256()
		if err := hashBlobHeader(h, header); err != nil {
			return nil, fmt.Errorf("failed to hash blob header at index %d: %w", i, err)
		}
		blobHeaderHashes[i] = h.Sum(nil)
	}

	return blobHeaderHashes, nil
}

func hashBlobCertificate(hasher hash.Hash, blobCertificate *common.BlobCertificate) error {
	err := hashBlobHeader(hasher, blobCertificate.GetBlobHeader())
	if err != nil {
		return fmt.Errorf("failed to hash blob header: %w", err)
	}
	err = hashByteArray(hasher, blobCertificate.GetSignature())
	if err != nil {
		return fmt.Errorf("failed to hash signature: %w", err)
	}
	err = hashUint32Array(hasher, blobCertificate.GetRelayKeys())
	if err != nil {
		return fmt.Errorf("failed to hash RelayKeys: %w", err)
	}
	return nil
}

func hashBlobHeader(hasher hash.Hash, header *common.BlobHeader) error {
	hashUint32(hasher, header.GetVersion())
	hashUint32(hasher, uint32(len(header.GetQuorumNumbers())))

	err := hashUint32Array(hasher, header.GetQuorumNumbers())
	if err != nil {
		return fmt.Errorf("failed to hash QuorumNumbers: %w", err)
	}

	err = hashBlobCommitment(hasher, header.GetCommitment())
	if err != nil {
		return fmt.Errorf("failed to hash commitment: %w", err)
	}

	err = hashPaymentHeader(hasher, header.GetPaymentHeader())
	if err != nil {
		return fmt.Errorf("failed to hash payment header: %w", err)
	}

	return nil
}

func hashBatchHeader(hasher hash.Hash, header *common.BatchHeader) error {
	err := hashByteArray(hasher, header.GetBatchRoot())
	if err != nil {
		return fmt.Errorf("failed to hash BatchRoot: %w", err)
	}
	hashUint64(hasher, header.GetReferenceBlockNumber())

	return nil
}

func hashBlobCommitment(hasher hash.Hash, commitment *commonv1.BlobCommitment) error {
	err := hashByteArray(hasher, commitment.GetCommitment())
	if err != nil {
		return fmt.Errorf("failed to hash commitment: %w", err)
	}

	err = hashByteArray(hasher, commitment.GetLengthCommitment())
	if err != nil {
		return fmt.Errorf("failed to hash LengthCommitment: %w", err)
	}

	err = hashByteArray(hasher, commitment.GetLengthProof())
	if err != nil {
		return fmt.Errorf("failed to hash LengthProof: %w", err)
	}

	hashUint32(hasher, commitment.GetLength())

	return nil
}

func hashPaymentHeader(hasher hash.Hash, header *common.PaymentHeader) error {
	err := hashByteArray(hasher, []byte(header.GetAccountId()))
	if err != nil {
		return fmt.Errorf("failed to hash AccountId: %w", err)
	}

	hashInt64(hasher, header.GetTimestamp())

	err = hashByteArray(hasher, header.GetCumulativePayment())
	if err != nil {
		return fmt.Errorf("failed to hash CumulativePayment: %w", err)
	}

	return nil
}
