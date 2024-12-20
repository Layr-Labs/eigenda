package auth

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/Layr-Labs/eigenda/core"

	"github.com/ethereum/go-ethereum/crypto"
)

type authenticator struct{}

var _ core.BlobRequestAuthenticator = &authenticator{}

func NewAuthenticator() core.BlobRequestAuthenticator {
	return &authenticator{}
}

func (*authenticator) AuthenticateBlobRequest(header core.BlobAuthHeader) error {
	sig := header.AuthenticationData

	// Ensure the signature is 65 bytes (Recovery ID is the last byte)
	if len(sig) != 65 {
		return fmt.Errorf("signature length is unexpected: %d", len(sig))
	}

	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, header.Nonce)
	hash := crypto.Keccak256(buf)
	// Verify the signature
	sigPublicKeyECDSA, err := crypto.SigToPub(hash, sig)
	if err != nil {
		return fmt.Errorf("failed to recover public key from signature: %v", err)
	}

	pubKey := crypto.PubkeyToAddress(*sigPublicKeyECDSA).Hex()

	if pubKey != header.AccountID {
		return errors.New("signature doesn't match with provided public key")
	}

	return nil

}
