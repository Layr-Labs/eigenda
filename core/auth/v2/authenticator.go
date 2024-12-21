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

	// Recover public key from signature
	sigPublicKeyECDSA, err := crypto.SigToPub(blobKey[:], sig)
	if err != nil {
		return fmt.Errorf("failed to recover public key from signature: %v", err)
	}

	accountId := header.PaymentMetadata.AccountID
	accountAddr := common.HexToAddress(accountId)
	pubKeyAddr := crypto.PubkeyToAddress(*sigPublicKeyECDSA)

	if accountAddr.Cmp(pubKeyAddr) != 0 {
		return errors.New("signature doesn't match with provided public key")
	}

	return nil
}

func (*authenticator) AuthenticatePaymentStateRequest(sig []byte, accountId string) error {
	// Ensure the signature is 65 bytes (Recovery ID is the last byte)
	if len(sig) != 65 {
		return fmt.Errorf("signature length is unexpected: %d", len(sig))
	}

	// Verify the signature
	hash := sha256.Sum256([]byte(accountId))
	sigPublicKeyECDSA, err := crypto.SigToPub(hash[:], sig)
	if err != nil {
		return fmt.Errorf("failed to recover public key from signature: %v", err)
	}

	accountAddr := common.HexToAddress(accountId)
	pubKeyAddr := crypto.PubkeyToAddress(*sigPublicKeyECDSA)

	if accountAddr.Cmp(pubKeyAddr) != 0 {
		return errors.New("signature doesn't match with provided public key")
	}

	return nil
}
