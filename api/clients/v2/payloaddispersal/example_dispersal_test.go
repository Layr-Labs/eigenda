package payloaddispersal_test

import (
	"context"
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients/v2/exampleutils"
)

// This example demonstrates how to use the PayloadDisperser to disperse a payload to EigenDA.
func Example_payloadDispersal() {
	// You must provide a private key that either has a testnet reservation, or you must configure on-demand payments
	// by sending funds to the payment vault.
	privateKey := "aaaaaaa"

	payloadDisperser, err := exampleutils.CreatePayloadDisperser(privateKey)
	if err != nil {
		panic(fmt.Sprintf("failed to create payload disperser: %v", err))
	}
	defer payloadDisperser.Close()

	// Create a sample payload (4KB of random data)
	payload, err := exampleutils.CreateRandomPayload(4 * 1024)
	if err != nil {
		panic(fmt.Sprintf("failed to generate random payload: %v", err))
	}

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	// Send the payload to EigenDA
	// This handles the entire dispersal workflow:
	// 1. Converting payload to blob
	// 2. Dispersing the blob
	// 3. Polling until the blob is certified
	// 4. Building and verifying the EigenDA certificate
	_, err = payloadDisperser.SendPayload(ctx, payload)
	if err != nil {
		panic(fmt.Sprintf("failed to disperse payload: %v", err))
	}

	// Note: In a real implementation, you would typically:
	// 1. Store the EigenDA certificate for later verification
	// 2. Use the certificate with a PayloadRetriever to retrieve the payload
	//    from EigenDA when needed

	fmt.Printf("Successfully dispersed payload")
}
