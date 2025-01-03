package auth

import (
	"crypto/ecdsa"
	"fmt"
	grpc "github.com/Layr-Labs/eigenda/api/grpc/node/v2"
	"github.com/Layr-Labs/eigenda/api/hashing"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// SignStoreChunksRequest signs the given StoreChunksRequest with the given private key. Does not
// write the signature into the request.
func SignStoreChunksRequest(key *ecdsa.PrivateKey, request *grpc.StoreChunksRequest) ([]byte, error) {
	requestHash := hashing.HashStoreChunksRequest(request)

	signature, err := crypto.Sign(requestHash, key)
	if err != nil {
		return nil, fmt.Errorf("failed to sign request: %w", err)
	}

	return signature, nil
}

// VerifyStoreChunksRequest verifies the given signature of the given StoreChunksRequest with the given
// public key.
func VerifyStoreChunksRequest(key gethcommon.Address, request *grpc.StoreChunksRequest) error {
	requestHash := hashing.HashStoreChunksRequest(request)

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
