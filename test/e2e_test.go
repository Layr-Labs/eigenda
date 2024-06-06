package test

import (
	"context"
	"flag"
	"fmt"
	"os"
	"runtime"
	"testing"
	"time"

	proxy "github.com/Layr-Labs/eigenda-proxy"
	"github.com/Layr-Labs/eigenda-proxy/eigenda"
	"github.com/Layr-Labs/eigenda-proxy/metrics"
	"github.com/Layr-Labs/eigenda-proxy/store"
	"github.com/Layr-Labs/eigenda-proxy/verify"
	"github.com/Layr-Labs/eigenda/api/clients"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	op_plasma "github.com/ethereum-optimism/optimism/op-plasma"

	oplog "github.com/ethereum-optimism/optimism/op-service/log"
	"github.com/ethereum/go-ethereum/log"
	"github.com/stretchr/testify/assert"
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
	port        = 6969
	holeskyDA   = "disperser-holesky.eigenda.xyz:443"
)

type TestSuite struct {
	ctx    context.Context
	log    log.Logger
	server *proxy.DAServer
}

func createTestSuite(t *testing.T) (TestSuite, func()) {
	ctx := context.Background()

	// load signer key from environment
	pk := os.Getenv(privateKey)
	if pk == "" {
		t.Fatal("SIGNER_PRIVATE_KEY environment variable not set")
	}

	log := oplog.NewLogger(os.Stdout, oplog.CLIConfig{
		Level:  log.LevelDebug,
		Format: oplog.FormatLogFmt,
		Color:  true,
	}).New("role", serviceName)

	oplog.SetGlobalLogHandler(log.Handler())

	testCfg := eigenda.Config{
		ClientConfig: clients.EigenDAClientConfig{
			RPC:                      holeskyDA,
			StatusQueryTimeout:       time.Minute * 45,
			StatusQueryRetryInterval: time.Second * 1,
			DisableTLS:               false,
			SignerPrivateKeyHex:      pk,
		},
	}

	// these values can be generated locally by running `make srs`
	kzgCfg := &kzg.KzgConfig{
		G1Path:          "../operator-setup/resources/g1.point",
		G2PowerOf2Path:  "../operator-setup/resources/g2.point.powerOf2",
		CacheDir:        "../operator-setup/resources/SRSTables",
		G2Path:          "../test/resources/g2.point",
		SRSOrder:        3000,
		SRSNumberToLoad: 3000,
		NumWorker:       uint64(runtime.GOMAXPROCS(0)),
	}

	client, err := clients.NewEigenDAClient(log, testCfg.ClientConfig)
	if err != nil {
		panic(err)
	}

	verifier, err := verify.NewVerifier(kzgCfg)
	if err != nil {
		panic(err)
	}

	daStore, err := store.NewEigenDAStore(ctx, client, verifier)
	if err != nil {
		panic(err)
	}

	server := proxy.NewServer(host, port, daStore, log, metrics.NoopMetrics)

	go func() {
		t.Log("Starting proxy server on separate routine...")
		if err := server.Start(); err != nil {
			panic(err)
		}
	}()

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

func TestE2EPutGetLogicForEigenDAStore(t *testing.T) {
	if !runTestnetIntegrationTests {
		t.Skip("Skipping testnet integration test")
	}

	ts, kill := createTestSuite(t)
	defer kill()

	daClient := op_plasma.NewDAClient(fmt.Sprintf("%s://%s:%d", transport, host, port), false, false)
	t.Log("Waiting for client to establish connection with plasma server...")
	// wait for server to come online after starting
	time.Sleep(5 * time.Second)

	// 1 - write arbitrary data to EigenDA

	var testPreimage = []byte("inter-subjective and not objective!")

	t.Log("Setting input data on proxy server...")
	commit, err := daClient.SetInput(ts.ctx, testPreimage)
	assert.NoError(t, err)

	// 2 - fetch data from EigenDA for generated commitment key
	t.Log("Getting input data from proxy server...")
	preimage, err := daClient.GetInput(ts.ctx, commit)
	assert.NoError(t, err)
	assert.Equal(t, testPreimage, preimage)

}
