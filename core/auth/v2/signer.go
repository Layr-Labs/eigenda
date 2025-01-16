package v2

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"fmt"
	"log"

	core "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type LocalBlobRequestSigner struct {
	PrivateKey *ecdsa.PrivateKey
}

var _ core.BlobRequestSigner = &LocalBlobRequestSigner{}

func NewLocalBlobRequestSigner(privateKeyHex string) *LocalBlobRequestSigner {
	privateKeyBytes := common.FromHex(privateKeyHex)
	privateKey, err := crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}

	return &LocalBlobRequestSigner{
		PrivateKey: privateKey,
	}
}

func (s *LocalBlobRequestSigner) SignBlobRequest(header *core.BlobHeader) ([]byte, error) {
	blobKey, err := header.BlobKey()
	if err != nil {
		return nil, fmt.Errorf("failed to get blob key: %v", err)
	}

	// Sign the blob key using the private key
	sig, err := crypto.Sign(blobKey[:], s.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign hash: %v", err)
	}

	return sig, nil
}

func (s *LocalBlobRequestSigner) SignPaymentStateRequest() ([]byte, error) {
	accountId, err := s.GetAccountID()
	if err != nil {
		return nil, fmt.Errorf("failed to get account ID: %v", err)
	}

	hash := sha256.Sum256([]byte(accountId))
	// Sign the account ID using the private key
	sig, err := crypto.Sign(hash[:], s.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign hash: %v", err)
	}

	return sig, nil
}

func (s *LocalBlobRequestSigner) GetAccountID() (string, error) {

	accountId := crypto.PubkeyToAddress(s.PrivateKey.PublicKey).Hex()
	return accountId, nil
}

type LocalNoopSigner struct{}

var _ core.BlobRequestSigner = &LocalNoopSigner{}

func NewLocalNoopSigner() *LocalNoopSigner {
	return &LocalNoopSigner{}
}

func (s *LocalNoopSigner) SignBlobRequest(header *core.BlobHeader) ([]byte, error) {
	return nil, fmt.Errorf("noop signer cannot sign blob request")
}

func (s *LocalNoopSigner) SignPaymentStateRequest() ([]byte, error) {
	return nil, fmt.Errorf("noop signer cannot sign payment state request")
}

func (s *LocalNoopSigner) GetAccountID() (string, error) {
	return "", fmt.Errorf("noop signer cannot get accountID")
}
