package core

import (
	"errors"
	"fmt"

	geth "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type BlobRequestAuthenticator interface {
	AuthenticateBlobRequest(header BlobAuthHeader) error
}

type BlobRequestSigner interface {
	SignBlobRequest(header BlobAuthHeader) ([]byte, error)
	GetAccountID() (string, error)
}

func VerifySignature(message []byte, accountAddr geth.Address, sig []byte) error {
	// Ensure the signature is 65 bytes (Recovery ID is the last byte)
	if len(sig) != 65 {
		return fmt.Errorf("signature length is unexpected: %d", len(sig))
	}

	// Verify the signature
	sigPublicKeyECDSA, err := crypto.SigToPub(message, sig)
	if err != nil {
		return fmt.Errorf("failed to recover public key from signature: %v", err)
	}

	pubKey := crypto.PubkeyToAddress(*sigPublicKeyECDSA).Hex()
	if pubKey != accountAddr.Hex() {
		return errors.New("signature doesn't match with provided public key")
	}

	return nil
}

type PaymentSigner interface {
	SignBlobPayment(header *PaymentMetadata) ([]byte, error)
	GetAccountID() string
}
