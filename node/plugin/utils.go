package plugin

import (
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
)

// Returns the decrypted ECDSA private key from the given file.
func GetECDSAPrivateKey(keyFile string, password string) (*keystore.Key, *string, error) {
	keyContents, err := os.ReadFile(keyFile)
	if err != nil {
		return nil, nil, err
	}
	sk, err := keystore.DecryptKey(keyContents, password)
	if err != nil {
		return nil, nil, err
	}
	privateKey := fmt.Sprintf("%x", crypto.FromECDSA(sk.PrivateKey))
	return sk, &privateKey, nil
}
