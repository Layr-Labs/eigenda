package e2e

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda-proxy/eigenda"
	"github.com/Layr-Labs/eigenda-proxy/metrics"
	"github.com/Layr-Labs/eigenda-proxy/server"
	"github.com/Layr-Labs/eigenda-proxy/store"
	"github.com/Layr-Labs/eigenda/api/clients"
	"github.com/ethereum/go-ethereum/log"

	oplog "github.com/ethereum-optimism/optimism/op-service/log"
	opmetrics "github.com/ethereum-optimism/optimism/op-service/metrics"

	"github.com/stretchr/testify/require"
)

const (
	privateKey = "SIGNER_PRIVATE_KEY"
	ethRPC     = "ETHEREUM_RPC"
	transport  = "http"
	svcName    = "eigenda_proxy"
	host       = "127.0.0.1"
	holeskyDA  = "disperser-holesky.eigenda.xyz:443"
)

type TestSuite struct {
	Ctx    context.Context
	Log    log.Logger
	Server *server.Server
}

func CreateTestSuite(t *testing.T, useMemory bool) (TestSuite, func()) {
	ctx := context.Background()

	// load signer key from environment
	pk := os.Getenv(privateKey)
	if pk == "" && !useMemory {
		t.Fatal("SIGNER_PRIVATE_KEY environment variable not set")
	}

	// load node url from environment
	ethRPC := os.Getenv(ethRPC)
	if ethRPC != "" && !useMemory {
		t.Fatal("ETHEREUM_RPC environment variable is not set")
	}

	log := oplog.NewLogger(os.Stdout, oplog.CLIConfig{
		Level:  log.LevelDebug,
		Format: oplog.FormatLogFmt,
		Color:  true,
	}).New("role", svcName)

	eigendaCfg := eigenda.Config{
		ClientConfig: clients.EigenDAClientConfig{
			RPC:                      holeskyDA,
			StatusQueryTimeout:       time.Minute * 45,
			StatusQueryRetryInterval: time.Second * 1,
			DisableTLS:               false,
			SignerPrivateKeyHex:      pk,
		},
		EthRPC:                 ethRPC,
		SvcManagerAddr:         "0xD4A7E1Bd8015057293f0D0A557088c286942e84b", // incompatible with non holeskly networks
		CacheDir:               "../operator-setup/resources/SRSTables",
		G1Path:                 "../operator-setup/resources/g1_abbr.point",
		G2Path:                 "../test/resources/kzg/g2.point", // do we need this?
		MaxBlobLength:          "90kib",
		G2PowerOfTauPath:       "../operator-setup/resources/kzg/g2_abbr.point.powerOf2",
		PutBlobEncodingVersion: 0x00,
	}

	memstoreCfg := store.MemStoreConfig{
		Enabled:        useMemory,
		BlobExpiration: 14 * 24 * time.Hour,
	}

	store, err := server.LoadStore(
		server.CLIConfig{
			EigenDAConfig: eigendaCfg,
			MemStoreCfg:   memstoreCfg,
			MetricsCfg:    opmetrics.CLIConfig{},
		},
		ctx,
		log,
	)
	require.NoError(t, err)
	server := server.NewServer(host, 0, store, log, metrics.NoopMetrics)

	t.Log("Starting proxy server...")
	err = server.Start()
	require.NoError(t, err)

	kill := func() {
		if err := server.Stop(); err != nil {
			panic(err)
		}
	}

	return TestSuite{
		Ctx:    ctx,
		Log:    log,
		Server: server,
	}, kill
}

func (ts *TestSuite) Address() string {
	// read port from listener
	port := ts.Server.Port()

	return fmt.Sprintf("%s://%s:%d", transport, host, port)
}
