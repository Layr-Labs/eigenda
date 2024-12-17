package auth

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/binary"
	"errors"
	"fmt"
	commonv1 "github.com/Layr-Labs/eigenda/api/grpc/common"
	common "github.com/Layr-Labs/eigenda/api/grpc/common/v2"
	grpc "github.com/Layr-Labs/eigenda/api/grpc/node/v2"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"golang.org/x/crypto/cryptobyte"
	"golang.org/x/crypto/sha3"
	"hash"
	"math/big"
)

type signatureInfo struct {
	R *big.Int
	S *big.Int
}

// AddRecoveryID computes the recovery ID for a given signature and public key and adds it to the signature.
func addRecoveryID(hash []byte, pubKey *ecdsa.PublicKey, partialSignature []byte) error {
	for v := 0; v < 4; v++ {
		partialSignature[64] = byte(v)
		recoveredPubKey, err := secp256k1.RecoverPubkey(hash, partialSignature)
		if err != nil {
			return fmt.Errorf("failed to recover public key: %w", err)
		}

		x, y := elliptic.Unmarshal(secp256k1.S256(), recoveredPubKey)
		if x.Cmp(pubKey.X) == 0 && y.Cmp(pubKey.Y) == 0 {
			return nil
		}
	}

	return fmt.Errorf("no valid recovery ID found")
}

// pad32 pads a byte slice to 32 bytes, inserting zeros at the beginning if necessary.
func pad32(bytes []byte) []byte {
	if len(bytes) == 32 {
		return bytes
	}

	padded := make([]byte, 32)
	copy(padded[32-len(bytes):], bytes)
	return padded
}

// ParseKMSSignature parses a signature (KeySpecEccSecgP256k1) in the format returned by amazon KMS into the
// 65-byte format used by Ethereum.
func ParseKMSSignature(
	publicKey *ecdsa.PublicKey,
	hash []byte,
	signatureBytes []byte) ([]byte, error) {

	si := signatureInfo{}
	rest, err := asn1.Unmarshal(signatureBytes, &si)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal signature: %w", err)
	}
	if len(rest) > 0 {
		return nil, fmt.Errorf("trailing data after signature (%d bytes)", len(rest))
	}

	rBytes := pad32(si.R.Bytes())
	sBytes := pad32(si.S.Bytes())

	result := make([]byte, 65)
	copy(result[0:32], rBytes)
	copy(result[32:64], sBytes)

	err = addRecoveryID(hash, publicKey, result)
	if err != nil {
		return nil, fmt.Errorf("failed to compute recovery ID: %w", err)
	}

	return result, nil
}

type publicKeyInfo struct {
	Raw       asn1.RawContent
	Algorithm pkix.AlgorithmIdentifier
	PublicKey asn1.BitString
}

// ParseKMSPublicKey parses a public key (KeySpecEccSecgP256k1) in the format returned by amazon KMS
// into an ecdsa.PublicKey.
func ParseKMSPublicKey(keyBytes []byte) (*ecdsa.PublicKey, error) {
	pki := publicKeyInfo{}
	rest, err := asn1.Unmarshal(keyBytes, &pki)

	if err != nil {
		return nil, err
	}
	if len(rest) > 0 {
		return nil, fmt.Errorf("trailing data after public key (%d bytes)", len(rest))
	}

	rightAlignedKey := cryptobyte.String(pki.PublicKey.RightAlign())

	x, y := elliptic.Unmarshal(crypto.S256(), rightAlignedKey)
	if x == nil {
		return nil, errors.New("x509: failed to unmarshal elliptic curve point")
	}

	return &ecdsa.PublicKey{
		Curve: crypto.S256(),
		X:     x,
		Y:     y,
	}, nil
}

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
