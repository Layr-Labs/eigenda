package clients

import (
	"context"
	"flag"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/log"
	"github.com/stretchr/testify/assert"
)

var runTestnetIntegrationTests bool

func init() {
	flag.BoolVar(&runTestnetIntegrationTests, "testnet-integration", false, "Run testnet-based integration tests")
}

// TestClientUsingTestnet tests the eigenda client against holesky testnet disperser.
// We don't test waiting for finality because that adds 12 minutes to the test, and is not necessary
// because we already test for this in the unit tests using a mock disperser which is much faster.
func TestClientUsingTestnet(t *testing.T) {
	if !runTestnetIntegrationTests {
		t.Skip("Skipping testnet integration test")
	}

	t.Run("PutBlobWaitForConfirmationDepth0AndGetBlob", func(t *testing.T) {
		t.Parallel()
		logger := log.NewLogger(log.NewTerminalHandler(os.Stdout, true))
		client, err := NewEigenDAClient(logger, EigenDAClientConfig{
			RPC: "disperser-holesky.eigenda.xyz:443",
			// Should need way less than 20 minutes, but we set it to 20 minutes to be safe
			// In worst case we had 10 min batching interval + some time for the tx to land onchain,
			// plus wait for 3 blocks of confirmation.
			StatusQueryTimeout:       20 * time.Minute,
			StatusQueryRetryInterval: 5 * time.Second,
			CustomQuorumIDs:          []uint{},
			SignerPrivateKeyHex:      "2d23e142a9e86a9175b9dfa213f20ea01f6c1731e09fa6edf895f70fe279cbb1",
			WaitForFinalization:      false,
			WaitForConfirmationDepth: 0,
			SvcManagerAddr:           "0xD4A7E1Bd8015057293f0D0A557088c286942e84b",
			EthRpcUrl:                "https://1rpc.io/holesky",
		})
		data := "hello world!"
		assert.NoError(t, err)
		blobInfo, err := client.PutBlob(context.Background(), []byte(data))
		assert.NoError(t, err)
		batchHeaderHash := blobInfo.BlobVerificationProof.BatchMetadata.BatchHeaderHash
		blobIndex := blobInfo.BlobVerificationProof.BlobIndex
		blob, err := client.GetBlob(context.Background(), batchHeaderHash, blobIndex)
		assert.NoError(t, err)
		assert.Equal(t, data, string(blob))
	})

	t.Run("PutBlobWaitForConfirmationDepth3AndGetBlob", func(t *testing.T) {
		t.Parallel()
		confDepth := uint64(3)
		logger := log.NewLogger(log.NewTerminalHandler(os.Stdout, true))
		client, err := NewEigenDAClient(logger, EigenDAClientConfig{
			RPC: "disperser-holesky.eigenda.xyz:443",
			// Should need way less than 20 minutes, but we set it to 20 minutes to be safe
			// In worst case we had 10 min batching interval + some time for the tx to land onchain,
			// plus wait for 3 blocks of confirmation.
			StatusQueryTimeout:       20 * time.Minute,
			StatusQueryRetryInterval: 5 * time.Second,
			CustomQuorumIDs:          []uint{},
			SignerPrivateKeyHex:      "2d23e142a9e86a9175b9dfa213f20ea01f6c1731e09fa6edf895f70fe279cbb1",
			WaitForFinalization:      false,
			WaitForConfirmationDepth: confDepth,
			SvcManagerAddr:           "0xD4A7E1Bd8015057293f0D0A557088c286942e84b",
			EthRpcUrl:                "https://1rpc.io/holesky",
		})
		data := "hello world!"
		assert.NoError(t, err)
		blobInfo, err := client.PutBlob(context.Background(), []byte(data))
		assert.NoError(t, err)
		batchHeaderHash := blobInfo.BlobVerificationProof.BatchMetadata.BatchHeaderHash
		blobIndex := blobInfo.BlobVerificationProof.BlobIndex
		blob, err := client.GetBlob(context.Background(), batchHeaderHash, blobIndex)
		assert.NoError(t, err)
		assert.Equal(t, data, string(blob))

		// assert confirmation depth by making sure the batch metadata hash was registered onchain
		// at least confDepth blocks ago
		blockNumCur, err := client.ethClient.BlockNumber(context.Background())
		assert.NoError(t, err)
		blockNumAtDepth := new(big.Int).SetUint64(blockNumCur - confDepth)
		batchId := blobInfo.BlobVerificationProof.GetBatchId()
		onchainBatchMetadataHash, err := client.edasmCaller.BatchIdToBatchMetadataHash(&bind.CallOpts{BlockNumber: blockNumAtDepth}, batchId)
		assert.NoError(t, err)
		assert.NotEqual(t, onchainBatchMetadataHash, make([]byte, 32))
	})
}
