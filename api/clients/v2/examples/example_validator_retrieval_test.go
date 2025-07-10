package examples

import (
	"context"
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients/v2/coretypes"
)

// This example demonstrates how to use the ValidatorPayloadRetriever to retrieve a payload from EigenDA, running on
// holesky testnet.
//
// The ValidatorPayloadRetriever retrieves the payload from the EigenDA validator nodes directly. This is a fallback
// retrieval mechanism: in normal operation, payloads are retrieved from a network of relays (see
// Example_relayPayloadRetrieval for an example). Retrieval directly from the EigenDA validator nodes is a fallback
// option that provides an additional security guarantee: regardless of whether the relay network is able to serve a
// retrieval request, a user always has the option of retrieving data directly from the nodes which have attested to
// the availability of the data.
func Example_validatorPayloadRetrieval() {
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
	defer payloadDisperser.Close() //nolint:errcheck // just an example, so we ignore the error

	payload, err := createRandomPayload(4 * 1024) // (4KB of random data)
	if err != nil {
		panic(fmt.Sprintf("create random payload: %v", err))
	}

	eigenDACert, err := payloadDisperser.SendPayload(ctx, payload)
	if err != nil {
		panic(fmt.Sprintf("send payload: %v", err))
	}

	fmt.Printf("Successfully dispersed payload\n")

	// Create a validator payload retriever to retrieve directly from validator nodes
	validatorPayloadRetriever, err := createValidatorPayloadRetriever()
	if err != nil {
		panic(fmt.Sprintf("create validator payload retriever: %v", err))
	}

	retrievableCert, ok := eigenDACert.(*coretypes.EigenDACertV3)
	if !ok {
		panic("eigenDACert is not a EigenDACertV3")
	}

	// Retrieve the payload using the certificate by fetching from validator nodes
	_, err = validatorPayloadRetriever.GetPayload(ctx, retrievableCert)
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
