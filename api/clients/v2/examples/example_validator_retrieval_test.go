package examples

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/status-im/keycard-go/hexutils"
)

// This example demonstrates how to use the ValidatorPayloadRetriever to retrieve a payload from EigenDA, running on
// holesky testnet
func Example_validatorPayloadRetrieval() {
	// You must provide a private key that either has a testnet reservation, or you must configure on-demand payments
	// by sending funds to the payment vault.
	privateKey := ""

	// Create a payload disperser and disperse a sample payload to EigenDA
	// This will be the payload we will later retrieve
	payloadDisperser, err := createPayloadDisperser(privateKey)
	if err != nil {
		panic(fmt.Sprintf("create payload disperser: %v", err))
	}
	defer payloadDisperser.Close()

	payload, err := createRandomPayload(4 * 1024) // (4KB of random data)
	if err != nil {
		panic(fmt.Sprintf("create random payload: %v", err))
	}

	dispersalCtx, dispersalCancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer dispersalCancel()
	eigenDACert, err := payloadDisperser.SendPayload(dispersalCtx, payload)
	if err != nil {
		panic(fmt.Sprintf("send payload: %v", err))
	}

	fmt.Printf("Successfully dispersed payload\n")

	// Create a validator payload retriever to retrieve directly from validator nodes
	validatorPayloadRetriever, err := createValidatorPayloadRetriever()
	if err != nil {
		panic(fmt.Sprintf("create validator payload retriever: %v", err))
	}

	// Create a context with timeout for retrieval
	retrievalCtx, retrievalCancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer retrievalCancel()
	// Retrieve the payload using the certificate by fetching from validator nodes
	retrievedPayload, err := validatorPayloadRetriever.GetPayload(retrievalCtx, eigenDACert)
	if err != nil {
		panic(fmt.Sprintf("get payload: %v", err))
	}

	// Verify that the retrieved payload matches the original by comparing bytes
	originalBytes := payload.Serialize()
	retrievedBytes := retrievedPayload.Serialize()
	if !bytes.Equal(originalBytes, retrievedBytes) {
		panic(fmt.Sprintf(
			"retrieved payload doesn't match original payload (original: %s, retrieved: %s)",
			hexutils.BytesToHex(originalBytes), hexutils.BytesToHex(retrievedBytes)))
	}

	fmt.Printf("Successfully retrieved payload\n")

	// Create a cert verifier, to verify the certificate on chain
	certVerifier, err := createCertVerifier()
	if err != nil {
		panic(fmt.Sprintf("create cert verifier: %v", err))
	}

	verificationCtx, verificationCancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer verificationCancel()
	err = certVerifier.VerifyCertV2(verificationCtx, eigenDACert)
	if err != nil {
		panic(fmt.Sprintf("verify cert: %v", err))
	}

	fmt.Printf("Successfully verified eigenDACert")

	// Output is disabled, since tests fail without a valid payment address
	// DisabledOutput: Successfully dispersed payload
	// Successfully retrieved payload
	// Successfully verified eigenDACert
}
