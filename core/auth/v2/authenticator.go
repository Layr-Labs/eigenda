package v2

import (
	"bytes"
	"errors"
	"fmt"

	core "github.com/Layr-Labs/eigenda/core/v2"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

type authenticator struct{}

func NewAuthenticator() *authenticator {
	return &authenticator{}
}

var _ core.BlobRequestAuthenticator = &authenticator{}

func (*authenticator) AuthenticateBlobRequest(header *core.BlobHeader) error {
	sig := header.Signature

	// Ensure the signature is 65 bytes (Recovery ID is the last byte)
	if len(sig) != 65 {
		return fmt.Errorf("signature length is unexpected: %d", len(sig))
	}

	blobKey, err := header.BlobKey()
	if err != nil {
		return fmt.Errorf("failed to get blob key: %v", err)
	}

	publicKeyBytes, err := hexutil.Decode(header.PaymentMetadata.AccountID)
	if err != nil {
		return fmt.Errorf("failed to decode public key (%v): %v", header.PaymentMetadata.AccountID, err)
	}

	// Decode public key
	pubKey, err := crypto.UnmarshalPubkey(publicKeyBytes)
	if err != nil {
		return fmt.Errorf("failed to decode public key (%v): %v", header.PaymentMetadata.AccountID, err)
	}

	// Verify the signature
	sigPublicKeyECDSA, err := crypto.SigToPub(blobKey[:], sig)
	if err != nil {
		return fmt.Errorf("failed to recover public key from signature: %v", err)
	}

	if !bytes.Equal(pubKey.X.Bytes(), sigPublicKeyECDSA.X.Bytes()) || !bytes.Equal(pubKey.Y.Bytes(), sigPublicKeyECDSA.Y.Bytes()) {
		return errors.New("signature doesn't match with provided public key")
	}

	return nil
}
