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
	"github.com/Layr-Labs/eigenda/core/thegraph"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/Layr-Labs/eigenda/encoding/kzg/verifier"
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
	t                 *testing.T
	logger            logging.Logger
	disperserClient   clients.DisperserClient
	relayClient       clients.RelayClient
	indexedChainState core.IndexedChainState
	retrievalClient   clients.RetrievalClient
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
	ethClient, err := geth.NewMultiHomingClient(ethClientConfig, accountId, logger)
	require.NoError(t, err)

	ethReader, err := eth.NewReader(
		logger,
		ethClient,
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

	// Construct the retrieval client

	chainState := eth.NewChainState(ethReader, ethClient)
	icsConfig := thegraph.Config{
		Endpoint:     "https://subgraph.satsuma-prod.com/51caed8fa9cb/eigenlabs/eigenda-operator-state-preprod-holesky/version/v0.7.0/api",
		PullInterval: 100 * time.Millisecond,
		MaxRetries:   5,
	}
	indexedChainState := thegraph.MakeIndexedChainState(icsConfig, chainState, logger)

	//TRAFFIC_GENERATOR_G1_PATH=../../inabox/resources/kzg/g1.point \
	//TRAFFIC_GENERATOR_G2_PATH=../../inabox/resources/kzg/g2.point \
	//TRAFFIC_GENERATOR_CACHE_PATH=../../inabox/resources/kzg/SRSTables \
	//TRAFFIC_GENERATOR_SRS_ORDER=3000 \
	//TRAFFIC_GENERATOR_SRS_LOAD=3000 \

	kzgConfig := &kzg.KzgConfig{
		LoadG2Points:    true,
		G1Path:          "/Users/cody/ws/srs/g1.point",
		G2Path:          "/Users/cody/ws/srs/g2.point",
		G2PowerOf2Path:  "/Users/cody/ws/srs/g2.point.powerOf2",
		CacheDir:        "/Users/cody/ws/srs/SRSTables",
		SRSOrder:        268435456,
		SRSNumberToLoad: 2097152,
		NumWorker:       32,
	}
	blobVerifier, err := verifier.NewVerifier(kzgConfig, nil)
	require.NoError(t, err)

	retrievalClient := clients.NewRetrievalClient(
		logger,
		ethReader,
		indexedChainState,
		blobVerifier,
		20)

	return &TestClient{
		t:                 t,
		logger:            logger,
		disperserClient:   disperserClient,
		relayClient:       relayClient,
		indexedChainState: indexedChainState,
		retrievalClient:   retrievalClient,
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

	currentBlockNumber, err := c.indexedChainState.GetCurrentBlockNumber()
	require.NoError(c.t, err)

	// Read the blob from the validators from each quorum.
	for _, quorumID := range quorums {
		fmt.Printf("Reading blob from validators for quorum %d\n", quorumID)
		header, err := corev2.BlobHeaderFromProtobuf(blobCert.BlobHeader)
		require.NoError(c.t, err)

		retrievedBlob, err := c.retrievalClient.GetBlob(ctx, header, uint64(currentBlockNumber), quorumID)
		require.NoError(c.t, err)

		retrievedPayload := codec.RemoveEmptyByteFromPaddedBytes(retrievedBlob)

		// The payload may have a bunch of 0s appended at the end. Remove them.
		require.True(c.t, len(retrievedPayload) >= len(payload))
		truncatedPayload := retrievedPayload[:len(payload)]

		// Only 0s should be appended at the end.
		for i := len(payload); i < len(retrievedPayload); i++ {
			require.Equal(c.t, byte(0), retrievedPayload[i])
		}

		require.Equal(c.t, payload, truncatedPayload)
	}
}
