package main

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"os"
	"strconv"

	"github.com/Layr-Labs/eigenda/api/hashing"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Printf("Usage: %s <account_id> <private_key> <timestamp>\n", os.Args[0])
		os.Exit(1)
	}

	accountID := os.Args[1]
	privateKeyHex := os.Args[2]
	timestamp, err := strconv.ParseUint(os.Args[3], 10, 64)
	if err != nil {
		fmt.Printf("Error parsing timestamp: %v\n", err)
		os.Exit(1)
	}

	// Convert account ID to address
	accountAddr := common.HexToAddress(accountID)
	fmt.Printf("Account bytes (hex): %x\n", accountAddr.Bytes())
	fmt.Printf("Account bytes length: %d\n", len(accountAddr.Bytes()))

	// Convert private key hex to ECDSA private key
	privateKeyBytes := common.FromHex(privateKeyHex)
	privateKey, err := crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		fmt.Printf("Error creating private key: %v\n", err)
		os.Exit(1)
	}

	// Hash the request using eigenda's hashing package
	requestHash, err := hashing.HashGetPaymentStateRequest(accountAddr, timestamp)
	if err != nil {
		fmt.Printf("Error hashing request: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Keccak256 hash: %x\n", requestHash)

	// Take SHA256 hash of the result
	hash := sha256.Sum256(requestHash)
	fmt.Printf("SHA256 hash: %x\n", hash)

	// Sign the raw hash without Ethereum prefix
	signature, err := crypto.Sign(hash[:], privateKey)
	if err != nil {
		fmt.Printf("Error signing hash: %v\n", err)
		os.Exit(1)
	}

	// Convert signature to base64
	signatureBase64 := base64.StdEncoding.EncodeToString(signature)

	// Print the signature
	fmt.Println(signatureBase64)
}
