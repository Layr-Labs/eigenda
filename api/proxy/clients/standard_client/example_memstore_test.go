package standard_client_test

import (
	"context"
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda-proxy/clients/standard_client"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// This example demonstrates how to use the standard client to
// send/retrieve payloads to/from the proxy running with a memstore backend,
// meaning that it fakes an actual EigenDA Network interaction.
func Example_proxyMemstoreV1() {
	// Start the proxy in memstore mode using testcontainers
	containerCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	proxyContainer, proxyURL := startProxyMemstoreV1(containerCtx)
	defer proxyContainer.Terminate(containerCtx) //nolint: errcheck // no need to check for error

	// ============= EXAMPLE STARTS HERE =================
	payload := []byte("my-eigenda-payload")

	client := standard_client.New(&standard_client.Config{URL: proxyURL})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Submit the payload to the proxy
	certBytes, err := client.SetData(ctx, payload)
	if err != nil {
		panic(err)
	}
	// 0x00 is for eigenda v1
	fmt.Printf("Cert header byte (encodes eigenda version): %x\n", certBytes[:1])

	// Retrieve the payload using the cert
	retrievedPayload, err := client.GetData(ctx, certBytes)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Retrieved payload: %s\n", retrievedPayload)
	// ============= EXAMPLE ENDS HERE =================

	// Output:
	// Cert header byte (encodes eigenda version): 00
	// Retrieved payload: my-eigenda-payload
}

// Start the proxy in memstore mode using testcontainers. This does the equivalent of:
// docker run --rm -p 3100:3100 ghcr.io/layr-labs/eigenda-proxy:latest --memstore.enabled --port 3100
// It returns the URL of the proxy.
func startProxyMemstoreV1(ctx context.Context) (testcontainers.Container, string) {
	req := testcontainers.ContainerRequest{
		Image:        "ghcr.io/layr-labs/eigenda-proxy:latest",
		ExposedPorts: []string{"3100/tcp"},
		WaitingFor:   wait.ForHTTP("/health").WithPort("3100/tcp"),
		Cmd:          []string{"--memstore.enabled", "--port", "3100"},
	}
	proxyContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		panic(err)
	}

	proxyEndpoint, err := proxyContainer.PortEndpoint(ctx, "3100", "http")
	if err != nil {
		panic(err)
	}
	return proxyContainer, proxyEndpoint
}
