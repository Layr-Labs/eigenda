package auth

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"os"
)

// ReadPublicECDSAKeyFile reads a public ECDSA key from a .pem file.
func ReadPublicECDSAKeyFile(publicKeyFile string) (*ecdsa.PublicKey, error) {
	file, err := os.Open(publicKeyFile)
	if err != nil {
		return nil, fmt.Errorf("error opening public key file: %w", err)
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("error reading public key file: %w", err)
	}

	block, _ := pem.Decode(bytes)
	if block == nil {
		return nil, fmt.Errorf("no PEM data found in public key file")
	}

	if block.Type != "PUBLIC KEY" {
		return nil, fmt.Errorf("unexpected block type: %s", block.Type)
	}

	genericPublicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("error parsing public key: %w", err)
	}

	publicKey := genericPublicKey.(*ecdsa.PublicKey)

	return publicKey, nil
}

// ReadPrivateECDSAKeyFile reads a private ECDSA key from a .pem file.
func ReadPrivateECDSAKeyFile(privateKeyFile string) (*ecdsa.PrivateKey, error) {
	file, err := os.Open(privateKeyFile)
	if err != nil {
		return nil, fmt.Errorf("error opening private key file: %w", err)
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("error reading private key file: %w", err)
	}

	block, _ := pem.Decode(bytes)
	if block == nil {
		return nil, fmt.Errorf("no PEM data found in private key file")
	}

	if block.Type != "EC PRIVATE KEY" {
		return nil, fmt.Errorf("unexpected block type: %s", block.Type)
	}

	privateKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("error parsing private key: %w", err)
	}

	return privateKey, nil
}
