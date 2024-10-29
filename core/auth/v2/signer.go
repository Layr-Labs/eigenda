package v2

import (
	"crypto/ecdsa"
	"fmt"
	"log"

	core "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
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

func (s *LocalBlobRequestSigner) GetAccountID() (string, error) {

	publicKeyBytes := crypto.FromECDSAPub(&s.PrivateKey.PublicKey)
	return hexutil.Encode(publicKeyBytes), nil

}

type LocalNoopSigner struct{}

var _ core.BlobRequestSigner = &LocalNoopSigner{}

func NewLocalNoopSigner() *LocalNoopSigner {
	return &LocalNoopSigner{}
}

func (s *LocalNoopSigner) SignBlobRequest(header *core.BlobHeader) ([]byte, error) {
	return nil, fmt.Errorf("noop signer cannot sign blob request")
}

func (s *LocalNoopSigner) GetAccountID() (string, error) {
	return "", fmt.Errorf("noop signer cannot get accountID")
}
