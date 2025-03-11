package v2

import (
	"crypto/sha256"
	"errors"
	"fmt"

	core "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type authenticator struct{}

func NewAuthenticator() *authenticator {
	return &authenticator{}
}

var _ core.BlobRequestAuthenticator = &authenticator{}

func (*authenticator) AuthenticateBlobRequest(header *core.BlobHeader, signature []byte) error {
	// Ensure the signature is 65 bytes (Recovery ID is the last byte)
	if len(signature) != 65 {
		return fmt.Errorf("signature length is unexpected: %d", len(signature))
	}

	blobKey, err := header.BlobKey()
	if err != nil {
		return fmt.Errorf("failed to get blob key: %v", err)
	}

	// Recover public key from signature
	sigPublicKeyECDSA, err := crypto.SigToPub(blobKey[:], signature)
	if err != nil {
		return fmt.Errorf("failed to recover public key from signature: %v", err)
	}

	accountAddr := header.PaymentMetadata.AccountID
	pubKeyAddr := crypto.PubkeyToAddress(*sigPublicKeyECDSA)

	if accountAddr.Cmp(pubKeyAddr) != 0 {
		return errors.New("signature doesn't match with provided public key")
	}

	return nil
}

// AuthenticatePaymentStateRequest verifies the signature of the payment state request
// The signature is signed over the byte representation of the account ID
// See implementation of BlobRequestSigner.SignPaymentStateRequest for more details
func (*authenticator) AuthenticatePaymentStateRequest(sig []byte, accountAddr common.Address) error {
	// Ensure the signature is 65 bytes (Recovery ID is the last byte)
	if len(sig) != 65 {
		return fmt.Errorf("signature length is unexpected: %d", len(sig))
	}

	// Verify the signature
	hash := sha256.Sum256(accountAddr.Bytes())
	sigPublicKeyECDSA, err := crypto.SigToPub(hash[:], sig)
	if err != nil {
		return fmt.Errorf("failed to recover public key from signature: %v", err)
	}

	pubKeyAddr := crypto.PubkeyToAddress(*sigPublicKeyECDSA)

	if accountAddr.Cmp(pubKeyAddr) != 0 {
		return errors.New("signature doesn't match with provided public key")
	}

	return nil
}
