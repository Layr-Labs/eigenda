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

type signer struct {
	PrivateKey *ecdsa.PrivateKey
}

func NewSigner(privateKeyHex string) core.BlobRequestSigner {

	privateKeyBytes := common.FromHex(privateKeyHex)
	privateKey, err := crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}

	return &signer{
		PrivateKey: privateKey,
	}
}

func (s *signer) SignBlobRequest(header core.BlobAuthHeader) ([]byte, error) {

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

func (s *signer) GetAccountID() string {

	publicKeyBytes := crypto.FromECDSAPub(&s.PrivateKey.PublicKey)
	return hexutil.Encode(publicKeyBytes)

}
