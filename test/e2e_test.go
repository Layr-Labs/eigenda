package test

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/encoding/kzg"
	plasma "github.com/Layr-Labs/op-plasma-eigenda"
	"github.com/Layr-Labs/op-plasma-eigenda/eigenda"
	"github.com/Layr-Labs/op-plasma-eigenda/store"
	"github.com/Layr-Labs/op-plasma-eigenda/verify"
	"github.com/ethereum-optimism/optimism/op-e2e/e2eutils/wait"
	oplog "github.com/ethereum-optimism/optimism/op-service/log"
	"github.com/ethereum/go-ethereum/log"
	"github.com/stretchr/testify/assert"
)

// Use of single port makes tests incapable of running in parallel
const (
	transport   = "http"
	serviceName = "plasma_test_server"
	testSvrHost = "127.0.0.1"
	testSvrPort = 6969
	holeskyDA   = "disperser-holesky.eigenda.xyz:443"
)

type TestSuite struct {
	ctx    context.Context
	log    log.Logger
	server *plasma.DAServer
}

func createTestSuite(t *testing.T) (TestSuite, func()) {
	ctx := context.Background()

	log := oplog.NewLogger(os.Stdout, oplog.CLIConfig{
		Level:  log.LevelDebug,
		Format: oplog.FormatLogFmt,
		Color:  true,
	}).New("role", serviceName)

	oplog.SetGlobalLogHandler(log.Handler())

	testCfg := eigenda.Config{
		RPC:                      holeskyDA,
		StatusQueryTimeout:       time.Minute * 45,
		StatusQueryRetryInterval: time.Second * 1,
		UseTLS:                   true,
	}

	kzgCfg := &kzg.KzgConfig{
		G1Path:          "../kzg/g1.point",
		G2Path:          "../kzg/g2.point",
		G2PowerOf2Path:  "../kzg/g2.point.powerOf2",
		CacheDir:        "../kzg/SRSTables",
		SRSOrder:        3000,
		SRSNumberToLoad: 3000,
		NumWorker:       uint64(runtime.GOMAXPROCS(0)),
	}

	client := eigenda.NewEigenDAClient(log, testCfg)
	verifier, err := verify.NewVerifier(kzgCfg)
	if err != nil {
		panic(err)
	}

	daStore, err := store.NewEigenDAStore(ctx, client, verifier)
	if err != nil {
		panic(err)
	}

	server := plasma.NewDAServer(testSvrHost, testSvrPort, daStore, log)

	go func() {
		t.Log("Starting test plasma server on separate routine...")
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
	ts, kill := createTestSuite(t)
	defer kill()

	daClient := plasma.NewDAClient(fmt.Sprintf("%s://%s:%d", transport, testSvrHost, testSvrPort), false)
	t.Log("Waiting for client to establish connection with plasma server...")
	// wait for server to come online after starting
	err := wait.For(ts.ctx, 500*time.Millisecond, func() (bool, error) {
		return daClient.Health(), nil
	})

	assert.NoError(t, err)

	// 1 - write arbitrary data to test plasma server

	var testPreimage = []byte("inter-subjective and not objective!")

	t.Log("Setting input data on plasma server...")
	commit, err := daClient.SetInput(ts.ctx, testPreimage)
	assert.NoError(t, err)

	// 2 - fetch pre-image data from test plasma server
	t.Log("Getting input data from plasma server...")
	preimage, err := daClient.GetInput(ts.ctx, commit)
	assert.NoError(t, err)
	assert.Equal(t, testPreimage, preimage)

}
