package client_test

import (
	"context"
	"flag"
	"os"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/api/client"
	"github.com/ethereum/go-ethereum/log"
	"github.com/stretchr/testify/assert"
)

var runTestnetIntegrationTests bool

func init() {
	flag.BoolVar(&runTestnetIntegrationTests, "testnet-integration", false, "Run testnet-based integration tests")
}

func TestClient(t *testing.T) {
	if !runTestnetIntegrationTests {
		t.Skip("Skipping testnet integration test")
	}
	logger := log.NewLogger(log.NewTerminalHandler(os.Stderr, true))
	client, err := client.NewEigenDAClient(logger, client.Config{
		RPC:                      "disperser-holesky.eigenda.xyz:443",
		StatusQueryTimeout:       25 * time.Minute,
		StatusQueryRetryInterval: 5 * time.Second,
		CustomQuorumIDs:          []uint{},
		SignerPrivateKeyHex:      "2d23e142a9e86a9175b9dfa213f20ea01f6c1731e09fa6edf895f70fe279cbb1",
	})
	data := "hello world!"
	assert.NoError(t, err)
	cert, err := client.PutBlob(context.Background(), []byte(data))
	assert.NoError(t, err)
	blob, err := client.GetBlob(context.Background(), cert.BatchHeaderHash, cert.BlobIndex)
	assert.NoError(t, err)
	assert.Equal(t, data, string(blob))
}
