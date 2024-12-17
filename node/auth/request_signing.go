package auth

import (
	"crypto/ecdsa"
	"encoding/binary"
	"fmt"
	commonv1 "github.com/Layr-Labs/eigenda/api/grpc/common"
	common "github.com/Layr-Labs/eigenda/api/grpc/common/v2"
	grpc "github.com/Layr-Labs/eigenda/api/grpc/node/v2"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/sha3"
	"hash"
)

// SignStoreChunksRequest signs the given StoreChunksRequest with the given private key. Does not
// write the signature into the request.
func SignStoreChunksRequest(key *ecdsa.PrivateKey, request *grpc.StoreChunksRequest) ([]byte, error) {
	requestHash := HashStoreChunksRequest(request)

	signature, err := crypto.Sign(requestHash, key)
	if err != nil {
		return nil, fmt.Errorf("failed to sign request: %w", err)
	}

	return signature, nil
}

// VerifyStoreChunksRequest verifies the given signature of the given StoreChunksRequest with the given
// public key.
func VerifyStoreChunksRequest(key gethcommon.Address, request *grpc.StoreChunksRequest) error {
	requestHash := HashStoreChunksRequest(request)

	signingPublicKey, err := crypto.SigToPub(requestHash, request.Signature)
	if err != nil {
		return fmt.Errorf("failed to recover public key from signature: %w", err)
	}

	fmt.Printf("     public key: %v\n", signingPublicKey)

	signingAddress := crypto.PubkeyToAddress(*signingPublicKey)

	if key.Cmp(signingAddress) != 0 {
		return fmt.Errorf("signature doesn't match with provided public key")
	}
	return nil
}

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

func hashUint32(hasher hash.Hash, value uint32) {
	bytes := make([]byte, 4)
	binary.BigEndian.PutUint32(bytes, value)
	hasher.Write(bytes)
}

func hashUint64(hasher hash.Hash, value uint64) {
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, value)
	hasher.Write(bytes)
}
