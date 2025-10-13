package clients

import (
	"context"
	"crypto/ecdsa"
	"fmt"

	grpc "github.com/Layr-Labs/eigenda/api/grpc/validator"
	"github.com/Layr-Labs/eigenda/node/auth"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var _ DispersalRequestSigner = &localRequestSigner{}

// localRequestSigner implements DispersalRequestSigner using a local private key
type localRequestSigner struct {
	privateKey *ecdsa.PrivateKey
}

// NewLocalDispersalRequestSigner creates a new DispersalRequestSigner using a local private key.
// This signer uses secp256k1 curve for Ethereum compatibility.
func NewLocalDispersalRequestSigner(privateKey *ecdsa.PrivateKey) DispersalRequestSigner {
	return &localRequestSigner{
		privateKey: privateKey,
	}
}

// NewLocalDispersalRequestSignerFromHex creates a new DispersalRequestSigner from a hex-encoded private key.
// The private key should be in hex format (with or without 0x prefix).
func NewLocalDispersalRequestSignerFromHex(privateKeyHex string) (DispersalRequestSigner, error) {
	privateKeyBytes := common.FromHex(privateKeyHex)
	privateKey, err := crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	return &localRequestSigner{
		privateKey: privateKey,
	}, nil
}

func (s *localRequestSigner) SignStoreChunksRequest(
	ctx context.Context, request *grpc.StoreChunksRequest) ([]byte, error) {

	// Use the existing auth package function for signing
	return auth.SignStoreChunksRequest(s.privateKey, request)
}
