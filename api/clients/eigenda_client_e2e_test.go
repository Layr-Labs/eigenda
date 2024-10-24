package clients_test

import (
	"context"
	"flag"
	"os"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients"
	"github.com/ethereum/go-ethereum/log"
	"github.com/stretchr/testify/assert"
)

var runTestnetIntegrationTests bool

func init() {
	flag.BoolVar(&runTestnetIntegrationTests, "testnet-integration", false, "Run testnet-based integration tests")
}

func TestClientUsingTestnet(t *testing.T) {
	if !runTestnetIntegrationTests {
		t.Skip("Skipping testnet integration test")
	}
	logger := log.NewLogger(log.NewTerminalHandler(os.Stderr, true))
	client, err := clients.NewEigenDAClient(logger, clients.EigenDAClientConfig{
		RPC:                      "disperser-holesky.eigenda.xyz:443",
		StatusQueryTimeout:       25 * time.Minute,
		StatusQueryRetryInterval: 5 * time.Second,
		CustomQuorumIDs:          []uint{},
		SignerPrivateKeyHex:      "2d23e142a9e86a9175b9dfa213f20ea01f6c1731e09fa6edf895f70fe279cbb1",
		// Waiting for finality adds 12 minutes to the test, and is not necessary
		// because we already test for this correct behavior in the unit tests using a mock disperser
		// which is much faster.
		WaitForFinalization: false,
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
}
