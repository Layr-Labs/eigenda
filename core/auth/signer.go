package auth

import (
	"crypto/ecdsa"
	"encoding/binary"
	"fmt"
	"log"

	"github.com/Layr-Labs/eigenda/core"
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

func (s *LocalBlobRequestSigner) SignBlobRequest(header core.BlobAuthHeader) ([]byte, error) {

	// Message you want to sign
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, header.Nonce)
	hash := crypto.Keccak256(buf)

	// Sign the hash using the private key
	sig, err := crypto.Sign(hash, s.PrivateKey)
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

func NewLocalNoopSigner() *LocalNoopSigner {
	return &LocalNoopSigner{}
}

func (s *LocalNoopSigner) SignBlobRequest(header core.BlobAuthHeader) ([]byte, error) {
	return nil, fmt.Errorf("noop signer cannot sign blob request")
}

func (s *LocalNoopSigner) GetAccountID() (string, error) {
	return "", fmt.Errorf("noop signer cannot get accountID")
}
