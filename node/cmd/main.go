package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/common/memory"
	coreeth "github.com/Layr-Labs/eigenda/core/eth"
	rpccalls "github.com/Layr-Labs/eigensdk-go/metrics/collectors/rpc_calls"
	"github.com/docker/go-units"

	"github.com/Layr-Labs/eigenda/common/pubip"
	"github.com/Layr-Labs/eigenda/common/ratelimit"
	"github.com/Layr-Labs/eigenda/common/store"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/urfave/cli"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/node"
	"github.com/Layr-Labs/eigenda/node/flags"
	nodegrpc "github.com/Layr-Labs/eigenda/node/grpc"
)

var (
	bucketStoreSize          = 10000
	bucketMultiplier float32 = 2
	bucketDuration           = 450 * time.Second
)

func main() {
	app := cli.NewApp()
	app.Flags = flags.Flags
	app.Version = fmt.Sprintf("%s %s %s", node.SemVer, node.GitCommit, node.GitDate)
	app.Name = node.AppName
	app.Usage = "EigenDA Node"
	app.Description = "Service for receiving and storing encoded blobs from disperser"

	app.Action = NodeMain
	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf("application failed: %v", err)
	}

	select {}
}

func NodeMain(ctx *cli.Context) error {
	log.Println("Initializing Node")
	config, err := node.NewConfig(ctx)
	if err != nil {
		return err
	}

	logger, err := common.NewLogger(&config.LoggerConfig)
	if err != nil {
		return err
	}

	if config.GCSafetyBufferSizeGB > 0 {
		safetyBuffer := uint64(config.GCSafetyBufferSizeGB * float64(units.GiB))
		err = memory.SetGCMemorySafetyBuffer(logger, safetyBuffer)
		if err != nil {
			return fmt.Errorf("failed to set memory limit: %w", err)
		}
	}

	pubIPProvider := pubip.ProviderOrDefault(logger, config.PubIPProviders...)

	// Rate limiter
	reg := prometheus.NewRegistry()
	globalParams := common.GlobalRateParams{
		BucketSizes: []time.Duration{bucketDuration},
		Multipliers: []float32{bucketMultiplier},
		CountFailed: true,
	}

	bucketStore, err := store.NewLocalParamStore[common.RateBucketParams](bucketStoreSize)
	if err != nil {
		return err
	}

	ratelimiter := ratelimit.NewRateLimiter(reg, globalParams, bucketStore, logger)

	rpcCallsCollector := rpccalls.NewCollector(node.AppName, reg)
	client, err := geth.NewInstrumentedEthClient(config.EthClientConfig, rpcCallsCollector, logger)
	if err != nil {
		return fmt.Errorf("cannot create chain.Client: %w", err)
	}

	reader, err := coreeth.NewReader(
		logger,
		client,
		config.BLSOperatorStateRetrieverAddr,
		config.EigenDAServiceManagerAddr)
	if err != nil {
		return fmt.Errorf("cannot create eth.Reader: %w", err)
	}

	// Create the node.
	node, err := node.NewNode(reg, config, pubIPProvider, client, logger)
	if err != nil {
		return err
	}

	err = node.Start(context.Background())
	if err != nil {
		node.Logger.Error("could not start node", "error", err)
		return err
	}

	// Creates the GRPC server.

	// TODO(cody-littley): the metrics server is currently started by eigenmetrics, which is in another repo.
	//  When we fully remove v1 support, we need to start the metrics server inside the v2 metrics code.
	server := nodegrpc.NewServer(config, node, logger, ratelimiter)

	var serverV2 *nodegrpc.ServerV2
	if config.EnableV2 {
		serverV2, err = nodegrpc.NewServerV2(context.Background(), config, node, logger, ratelimiter, reg, reader)
		if err != nil {
			return fmt.Errorf("failed to create server v2: %v", err)
		}
	}
	err = nodegrpc.RunServers(server, serverV2, config, logger)

	return err
}
