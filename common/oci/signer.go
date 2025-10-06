package oci

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/oracle/oci-go-sdk/v65/keymanagement"
)

// OCIKMSSigner implements the signer interface using OCI KMS
type OCIKMSSigner struct {
	ctx              context.Context
	cryptoClient     keymanagement.KmsCryptoClient
	managementClient keymanagement.KmsManagementClient
	publicKey        *ecdsa.PublicKey
	keyOCID          string
	chainID          *big.Int
}

// NewOCIKMSSigner creates a new OCI KMS signer
func NewOCIKMSSigner(
	ctx context.Context,
	cryptoClient keymanagement.KmsCryptoClient,
	managementClient keymanagement.KmsManagementClient,
	publicKey *ecdsa.PublicKey,
	keyOCID string,
	chainID *big.Int) *OCIKMSSigner {

	return &OCIKMSSigner{
		ctx:              ctx,
		cryptoClient:     cryptoClient,
		managementClient: managementClient,
		publicKey:        publicKey,
		keyOCID:          keyOCID,
		chainID:          chainID,
	}
}

// SignMessage signs an arbitrary message using OCI KMS
func (s *OCIKMSSigner) SignMessage(message []byte) ([]byte, error) {
	return SignKMS(s.ctx, s.cryptoClient, s.keyOCID, s.publicKey, message)
}

// GetAddress returns the Ethereum address corresponding to the public key
func (s *OCIKMSSigner) GetAddress() common.Address {
	return crypto.PubkeyToAddress(*s.publicKey)
}

// SignTransaction signs an Ethereum transaction using OCI KMS
func (s *OCIKMSSigner) SignTransaction(tx *types.Transaction) (*types.Transaction, error) {
	signer := types.NewEIP155Signer(s.chainID)
	hash := signer.Hash(tx)

	signature, err := s.SignMessage(hash.Bytes())
	if err != nil {
		return nil, err
	}

	signedTx, err := tx.WithSignature(signer, signature)
	if err != nil {
		return nil, fmt.Errorf("failed to apply signature to transaction: %w", err)
	}
	return signedTx, nil
}

// GetChainID returns the chain ID
func (s *OCIKMSSigner) GetChainID() *big.Int {
	return s.chainID
}

// GetPublicKey returns the ECDSA public key
func (s *OCIKMSSigner) GetPublicKey() *ecdsa.PublicKey {
	return s.publicKey
}

// GetSignerFn returns a bind.SignerFn compatible function for use with walletsdk
func (s *OCIKMSSigner) GetSignerFn() bind.SignerFn {
	return func(address common.Address, tx *types.Transaction) (*types.Transaction, error) {
		return s.SignTransaction(tx)
	}
}
