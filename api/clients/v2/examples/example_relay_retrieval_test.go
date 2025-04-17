package examples

import (
	"context"
	"fmt"
	"time"
)

// This example demonstrates how to use the RelayPayloadRetriever to retrieve a payload from EigenDA, running on
// holesky testnet
func Example_relayPayloadRetrieval() {
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

	dispersalCtx, dispersalCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer dispersalCancel()
	eigenDACert, err := payloadDisperser.SendPayload(dispersalCtx, payload)
	if err != nil {
		panic(fmt.Sprintf("send payload: %v", err))
	}

	fmt.Printf("Successfully dispersed payload\n")

	// Create a payload retriever to retrieve from EigenDA relays
	payloadRetriever, err := createRelayPayloadRetriever()
	if err != nil {
		panic(fmt.Sprintf("create relay payload retriever: %v", err))
	}
	defer payloadRetriever.Close()

	retrievalCtx, retrievalCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer retrievalCancel()
	// Retrieve the payload using the certificate
	_, err = payloadRetriever.GetPayload(retrievalCtx, eigenDACert)
	if err != nil {
		panic(fmt.Sprintf("get payload: %v", err))
	}

	fmt.Printf("Successfully retrieved payload\n")

	// Create a cert verifier, to verify the certificate on chain
	certVerifier, err := createCertVerifier()
	if err != nil {
		panic(fmt.Sprintf("create cert verifier: %v", err))
	}

	verificationCtx, verificationCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer verificationCancel()
	// VerifyCertV2 is a view-only call to the `EigenDACertVerifier` contract. This call verifies that the provided cert
	// is valid: if this call doesn't return an error, then the eigenDA network has attested to the availability of the
	// dispersed blob.
	err = certVerifier.VerifyCertV2(verificationCtx, eigenDACert)
	if err != nil {
		panic(fmt.Sprintf("verify cert: %v", err))
	}

	fmt.Printf("Successfully verified eigenDACert")

	// Output is disabled, since tests fail without a valid payment address. To enable the test, delete this comment,
	// and change `DisabledOutput` to `Output`.

	// DisabledOutput: Successfully dispersed payload
	// Successfully retrieved payload
	// Successfully verified eigenDACert
}
