package auth

import (
	"crypto/ecdsa"
	"fmt"

	grpc "github.com/Layr-Labs/eigenda/api/grpc/validator"
	"github.com/Layr-Labs/eigenda/api/hashing"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// SignStoreChunksRequest signs the given StoreChunksRequest with the given private key. Does not
// write the signature into the request.
func SignStoreChunksRequest(key *ecdsa.PrivateKey, request *grpc.StoreChunksRequest) ([]byte, error) {
	requestHash, err := hashing.HashStoreChunksRequest(request)
	if err != nil {
		return nil, fmt.Errorf("failed to hash request: %w", err)
	}

	signature, err := crypto.Sign(requestHash, key)
	if err != nil {
		return nil, fmt.Errorf("failed to sign request: %w", err)
	}

	return signature, nil
}

// VerifyStoreChunksRequestWithKeys verifies the given signature against any of the provided keys.
// Returns the hash of the request if valid, error if signature doesn't match any key.
func VerifyStoreChunksRequestWithKeys(keys []gethcommon.Address, request *grpc.StoreChunksRequest) ([]byte, error) {
	requestHash, err := hashing.HashStoreChunksRequest(request)
	if err != nil {
		return nil, fmt.Errorf("failed to hash request: %w", err)
	}

	signingPublicKey, err := crypto.SigToPub(requestHash, request.GetSignature())
	if err != nil {
		return nil, fmt.Errorf("failed to recover public key from signature %x: %w", request.GetSignature(), err)
	}

	signingAddress := crypto.PubkeyToAddress(*signingPublicKey)

	for _, key := range keys {
		if key.Cmp(signingAddress) == 0 {
			return requestHash, nil
		}
	}

	return nil, fmt.Errorf("signature public key %x doesn't match any registered public key", signingAddress.Hex())
}
