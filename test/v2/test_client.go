package v2

import (
	"context"
	"fmt"
	"github.com/Layr-Labs/eigenda/api/clients/v2"
	v2 "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	"github.com/Layr-Labs/eigenda/core"
	auth "github.com/Layr-Labs/eigenda/core/auth/v2"
	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

// TestClient encapsulates the various clients necessary for interacting with EigenDA.
type TestClient struct {
	t               *testing.T
	disperserClient clients.DisperserClient
}

func NewTestClient(t *testing.T) *TestClient {
	disperserConfig := &clients.DisperserClientConfig{
		Hostname:          "disperser-preprod-holesky.eigenda.xyz",
		Port:              "443",
		UseSecureGrpcFlag: true,
	}

	privateKey := os.Getenv("V2_TEST_PRIVATE_KEY")
	if privateKey == "" {
		require.Fail(t, "V2_TEST_PRIVATE_KEY environment variable must be set")
	}

	signer := auth.NewLocalBlobRequestSigner(privateKey)
	signerAccountId, err := signer.GetAccountID()
	require.NoError(t, err)
	accountId := gethcommon.HexToAddress(signerAccountId)
	fmt.Printf("Account ID: %s\n", accountId.String())

	disperserClient, err := clients.NewDisperserClient(disperserConfig, signer, nil, nil)
	require.NoError(t, err)

	return &TestClient{
		t:               t,
		disperserClient: disperserClient,
	}
}

// DispersePayload sends a payload to the disperser. Waits until the payload is confirmed and then reads
// it back from the relays and the validators.
func (c *TestClient) DispersePayload(
	ctx context.Context,
	timeout time.Duration,
	payload []byte,
	quorums []core.QuorumID) {

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Disperse the blob.

	fmt.Printf("Dispersing payload of length %d to quorums %v\n", len(payload), quorums)

	padded := codec.ConvertByPaddingEmptyByte(payload)
	_, key, err := c.disperserClient.DisperseBlob(ctx, padded, 0, quorums, 0)
	require.NoError(c.t, err)
	fmt.Printf("Dispersed blob with key %x\n", key)
	fmt.Printf("Blob Status: %s\n", v2.BlobStatus_QUEUED.String())

	// Wait for the blob to be certified.

	status := v2.BlobStatus_QUEUED
	ticker := time.NewTicker(time.Second)
	start := time.Now()
	statusStart := start
	certified := false
	for !certified {
		select {
		case <-ticker.C:
			reply, err := c.disperserClient.GetBlobStatus(ctx, key)
			require.NoError(c.t, err)

			if reply.Status == v2.BlobStatus_CERTIFIED {
				elapsed := time.Since(statusStart)
				totalElapsed := time.Since(start)
				fmt.Printf(
					"Blob is certified (spent %0.1fs in prior status, total time %0.1fs)\n",
					elapsed.Seconds(),
					totalElapsed.Seconds())
				certified = true
			} else if reply.Status != status {
				elapsed := time.Since(statusStart)
				statusStart = time.Now()
				fmt.Printf("Blob status: %s (spent %0.1fs in prior status)\n",
					reply.Status.String(),
					elapsed.Seconds())
				status = reply.Status

				if reply.Status == v2.BlobStatus_FAILED ||
					reply.Status == v2.BlobStatus_UNKNOWN ||
					reply.Status == v2.BlobStatus_INSUFFICIENT_SIGNATURES {
					require.Fail(
						c.t,
						"Blob status is in a terminal non-successful state.",
						reply.Status.String())
				}
			}
		case <-ctx.Done():
			require.Fail(c.t, "Timed out waiting for blob to be confirmed")
		}
	}

	// TODO read the blob from the relays
	// TODO read the blob from the validators

}
