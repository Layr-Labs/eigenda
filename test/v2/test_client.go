package v2

import (
	"context"
	"fmt"
	"github.com/Layr-Labs/eigenda/api/clients/v2"
	commonv2 "github.com/Layr-Labs/eigenda/api/grpc/common/v2"
	v2 "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/core"
	auth "github.com/Layr-Labs/eigenda/core/auth/v2"
	"github.com/Layr-Labs/eigenda/core/eth"
	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
	"github.com/Layr-Labs/eigensdk-go/logging"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

// TestClient encapsulates the various clients necessary for interacting with EigenDA.
type TestClient struct {
	t               *testing.T
	logger          logging.Logger
	disperserClient clients.DisperserClient
	relayClient     clients.RelayClient
}

// TODO pass in args from outer scope

// NewTestClient creates a new TestClient instance.
func NewTestClient(t *testing.T) *TestClient {

	logger, err := common.NewLogger(common.DefaultTextLoggerConfig())
	require.NoError(t, err)

	// Construct the disperser client

	privateKey := os.Getenv("V2_TEST_PRIVATE_KEY")
	if privateKey == "" {
		require.Fail(t, "V2_TEST_PRIVATE_KEY environment variable must be set")
	}

	signer := auth.NewLocalBlobRequestSigner(privateKey)
	signerAccountId, err := signer.GetAccountID()
	require.NoError(t, err)
	accountId := gethcommon.HexToAddress(signerAccountId)
	fmt.Printf("Account ID: %s\n", accountId.String())

	disperserConfig := &clients.DisperserClientConfig{
		Hostname:          "disperser-preprod-holesky.eigenda.xyz",
		Port:              "443",
		UseSecureGrpcFlag: true,
	}
	disperserClient, err := clients.NewDisperserClient(disperserConfig, signer, nil, nil)
	require.NoError(t, err)

	// Construct the relay client

	rpcURLs := []string{"https://ethereum-holesky-rpc.publicnode.com"}
	ethClientConfig := geth.EthClientConfig{
		RPCURLs:          rpcURLs,
		PrivateKeyString: privateKey, // TODO is this correct?
		NumConfirmations: 0,
		NumRetries:       3,
	}
	client, err := geth.NewMultiHomingClient(ethClientConfig, accountId, logger)
	require.NoError(t, err)

	ethReader, err := eth.NewReader(
		logger,
		client,
		"0x93545e3b9013CcaBc31E80898fef7569a4024C0C",
		"0x54A03db2784E3D0aCC08344D05385d0b62d4F432")
	require.NoError(t, err)

	relayURLS, err := ethReader.GetRelayURLs(context.Background())
	require.NoError(t, err)

	relayConfig := &clients.RelayClientConfig{
		Sockets:           relayURLS,
		UseSecureGrpcFlag: true,
	}
	relayClient, err := clients.NewRelayClient(relayConfig, logger)
	require.NoError(t, err)

	return &TestClient{
		t:               t,
		logger:          logger,
		disperserClient: disperserClient,
		relayClient:     relayClient,
	}
}

// TODO break this into helper functions to make more readable

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
	var blobCert *commonv2.BlobCertificate = nil
	for blobCert == nil {
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

				blobCert = reply.BlobVerificationInfo.BlobCertificate
				require.NotNil(c.t, blobCert)
				require.True(c.t, len(blobCert.Relays) >= 1)
				// TODO additional verifications on blob cert
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

	// Read the blob from each relay in custody of the blob.
	for _, relayID := range blobCert.Relays {
		fmt.Printf("Reading blob from relay %d\n", relayID)
		blobFromRelay, err := c.relayClient.GetBlob(ctx, relayID, key)
		require.NoError(c.t, err)

		relayPayload := codec.RemoveEmptyByteFromPaddedBytes(blobFromRelay)
		require.Equal(c.t, payload, relayPayload)
	}

}
