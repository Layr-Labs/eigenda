package examples

import (
	"context"
	"fmt"
	"time"
)

// This example demonstrates how to use the RelayPayloadRetriever to retrieve a payload from EigenDA, running on
// holesky testnet
//
// The RelayPayloadRetriever retrieves the payload from a network of relays, which exists to serve EigenDA data. This
// is the standard retrieval mechanism. There exists an alternate retrieval mechanism, where payloads are retrieved
// from the EigenDA validator nodes directly (see Example_validatorPayloadRetrieval for an example), but users should
// default to retrieving data from the relays for optimal performance.
func Example_relayPayloadRetrieval() {
	// You must provide a private key that either has a testnet reservation, or you must configure on-demand payments
	// by sending funds to the payment vault.
	privateKey := ""

	ctx := context.Background()

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

	eigenDACert, err := payloadDisperser.SendPayload(ctx, payload)
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
	
	// Retrieve the payload using the certificate
	_, err = payloadRetriever.GetPayload(ctx, eigenDACert)
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
	// CheckDACert is a view-only call to the `EigenDACertVerifier` contract. This call verifies that the provided cert
	// is valid: if this call doesn't return an error, then the EigenDA network has attested to the availability of the
	// dispersed blob.

	err = certVerifier.CheckDACert(verificationCtx, eigenDACert)
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
