package test

import (
	"context"
	"flag"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda-proxy/client"
	"github.com/Layr-Labs/eigenda-proxy/common"
	"github.com/Layr-Labs/eigenda-proxy/eigenda"
	"github.com/Layr-Labs/eigenda-proxy/metrics"
	"github.com/Layr-Labs/eigenda-proxy/server"
	"github.com/Layr-Labs/eigenda-proxy/store"
	"github.com/Layr-Labs/eigenda/api/clients"
	"github.com/Layr-Labs/eigenda/api/clients/codecs"
	"github.com/ethereum-optimism/optimism/op-e2e/e2eutils/wait"
	op_plasma "github.com/ethereum-optimism/optimism/op-plasma"
	opmetrics "github.com/ethereum-optimism/optimism/op-service/metrics"

	oplog "github.com/ethereum-optimism/optimism/op-service/log"
	"github.com/ethereum/go-ethereum/log"
	"github.com/stretchr/testify/require"
)

var runTestnetIntegrationTests bool

func init() {
	flag.BoolVar(&runTestnetIntegrationTests, "testnet-integration", false, "Run testnet-based integration tests")
}

// Use of single port makes tests incapable of running in parallel
const (
	privateKey  = "SIGNER_PRIVATE_KEY"
	transport   = "http"
	serviceName = "eigenda_proxy"
	host        = "127.0.0.1"
	port        = 4200
	holeskyDA   = "disperser-holesky.eigenda.xyz:443"
)

type TestSuite struct {
	ctx    context.Context
	log    log.Logger
	server *server.Server
}

func createTestSuite(t *testing.T, useMemory bool) (TestSuite, func()) {
	ctx := context.Background()

	// load signer key from environment
	pk := os.Getenv(privateKey)
	if pk == "" && !useMemory {
		t.Fatal("SIGNER_PRIVATE_KEY environment variable not set")
	}

	log := oplog.NewLogger(os.Stdout, oplog.CLIConfig{
		Level:  log.LevelDebug,
		Format: oplog.FormatLogFmt,
		Color:  true,
	}).New("role", serviceName)

	oplog.SetGlobalLogHandler(log.Handler())

	eigendaCfg := eigenda.Config{
		ClientConfig: clients.EigenDAClientConfig{
			RPC:                      holeskyDA,
			StatusQueryTimeout:       time.Minute * 45,
			StatusQueryRetryInterval: time.Second * 1,
			DisableTLS:               false,
			SignerPrivateKeyHex:      pk,
		},
		CacheDir:               "../operator-setup/resources/SRSTables",
		G1Path:                 "../operator-setup/resources/g1_abbr.point",
		G2Path:                 "../test/resources/g2.point", // do we need this?
		MaxBlobLength:          "90kib",
		G2PowerOfTauPath:       "../operator-setup/resources/g2_abbr.point.powerOf2",
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
	server := server.NewServer(host, port, store, log, metrics.NoopMetrics)

	t.Log("Starting proxy server...")
	err = server.Start()
	require.NoError(t, err)

	kill := func() {
		if err := server.Stop(); err != nil {
			panic(err)
		}
	}

	return TestSuite{
		ctx:    ctx,
		log:    log,
		server: server,
	}, kill
}

func TestHoleskyWithPlasmaClient(t *testing.T) {
	if !runTestnetIntegrationTests {
		t.Skip("Skipping testnet integration test")
	}

	ts, kill := createTestSuite(t, false)
	defer kill()

	daClient := op_plasma.NewDAClient(fmt.Sprintf("%s://%s:%d", transport, host, port), false, false)
	t.Log("Waiting for client to establish connection with plasma server...")
	// wait for server to come online after starting
	time.Sleep(5 * time.Second)

	// 1 - write arbitrary data to EigenDA

	var testPreimage = []byte("inter-subjective and not objective!")

	t.Log("Setting input data on proxy server...")
	commit, err := daClient.SetInput(ts.ctx, testPreimage)
	require.NoError(t, err)

	// 2 - fetch data from EigenDA for generated commitment key
	t.Log("Getting input data from proxy server...")
	preimage, err := daClient.GetInput(ts.ctx, commit)
	require.NoError(t, err)
	require.Equal(t, testPreimage, preimage)
}

func TestHoleskyWithProxyClient(t *testing.T) {
	if !runTestnetIntegrationTests {
		t.Skip("Skipping testnet integration test")
	}

	ts, kill := createTestSuite(t, false)
	defer kill()

	cfg := &client.Config{
		URL: fmt.Sprintf("%s://%s:%d", transport, host, port),
	}
	daClient := client.New(cfg)
	t.Log("Waiting for client to establish connection with plasma server...")
	// wait for server to come online after starting
	wait.For(ts.ctx, time.Second*1, func() (bool, error) {
		err := daClient.Health()
		if err != nil {
			return false, nil
		}

		return true, nil
	})

	// 1 - write arbitrary data to EigenDA

	var testPreimage = []byte("inter-subjective and not objective!")

	t.Log("Setting input data on proxy server...")
	blobInfo, err := daClient.SetData(ts.ctx, testPreimage)
	require.NoError(t, err)

	// 2 - fetch data from EigenDA for generated commitment key
	t.Log("Getting input data from proxy server...")
	preimage, err := daClient.GetData(ts.ctx, blobInfo, common.BinaryDomain)
	require.NoError(t, err)
	require.Equal(t, testPreimage, preimage)

	// 3 - fetch iFFT representation of preimage
	iFFTPreimage, err := daClient.GetData(ts.ctx, blobInfo, common.PolyDomain)
	require.NoError(t, err)
	require.NotEqual(t, preimage, iFFTPreimage)

	// 4 - Assert domain transformations

	ifftCodec := codecs.NewIFFTCodec(codecs.DefaultBlobCodec{})

	decodedBlob, err := ifftCodec.DecodeBlob(iFFTPreimage)
	require.NoError(t, err)

	require.Equal(t, decodedBlob, preimage)
}

func TestMemStoreWithPlasmaClient(t *testing.T) {
	if runTestnetIntegrationTests {
		t.Skip("Skipping non-testnet integration test")
	}

	ts, kill := createTestSuite(t, true)
	defer kill()

	daClient := op_plasma.NewDAClient(fmt.Sprintf("%s://%s:%d", transport, host, port), false, false)
	t.Log("Waiting for client to establish connection with plasma server...")
	// wait for server to come online after starting
	time.Sleep(5 * time.Second)

	// 1 - write arbitrary data to EigenDA

	var testPreimage = []byte("inter-subjective and not objective!")

	t.Log("Setting input data on proxy server...")
	commit, err := daClient.SetInput(ts.ctx, testPreimage)
	require.NoError(t, err)

	// 2 - fetch data from EigenDA for generated commitment key
	t.Log("Getting input data from proxy server...")
	preimage, err := daClient.GetInput(ts.ctx, commit)
	require.NoError(t, err)
	require.Equal(t, testPreimage, preimage)
}

func TestMemStoreWithProxyClient(t *testing.T) {
	if runTestnetIntegrationTests {
		t.Skip("Skipping non-testnet integration test")
	}

	ts, kill := createTestSuite(t, true)
	defer kill()

	cfg := &client.Config{
		URL: fmt.Sprintf("%s://%s:%d", transport, host, port),
	}
	daClient := client.New(cfg)
	t.Log("Waiting for client to establish connection with plasma server...")
	// wait for server to come online after starting
	wait.For(ts.ctx, time.Second*1, func() (bool, error) {
		err := daClient.Health()
		if err != nil {
			return false, nil
		}

		return true, nil
	})

	// 1 - write arbitrary data to EigenDA

	var testPreimage = []byte("inter-subjective and not objective!")

	t.Log("Setting input data on proxy server...")
	blobInfo, err := daClient.SetData(ts.ctx, testPreimage)
	require.NoError(t, err)

	// 2 - fetch data from EigenDA for generated commitment key
	t.Log("Getting input data from proxy server...")
	preimage, err := daClient.GetData(ts.ctx, blobInfo, common.BinaryDomain)
	require.NoError(t, err)
	require.Equal(t, testPreimage, preimage)

	// 3 - fetch iFFT representation of preimage
	iFFTPreimage, err := daClient.GetData(ts.ctx, blobInfo, common.PolyDomain)
	require.NoError(t, err)
	require.NotEqual(t, preimage, iFFTPreimage)

	// 4 - Assert domain transformations

	ifftCodec := codecs.NewIFFTCodec(codecs.DefaultBlobCodec{})

	decodedBlob, err := ifftCodec.DecodeBlob(iFFTPreimage)
	require.NoError(t, err)

	require.Equal(t, decodedBlob, preimage)
}
