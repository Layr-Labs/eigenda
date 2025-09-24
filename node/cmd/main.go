package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Layr-Labs/eigenda/common/enforce"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/common/pubip"
	"github.com/Layr-Labs/eigenda/common/ratelimit"
	"github.com/Layr-Labs/eigenda/common/store"
	commonversion "github.com/Layr-Labs/eigenda/common/version"
	coreeth "github.com/Layr-Labs/eigenda/core/eth"
	"github.com/Layr-Labs/eigenda/core/eth/directory"
	rpccalls "github.com/Layr-Labs/eigensdk-go/metrics/collectors/rpc_calls"
	gethcommon "github.com/ethereum/go-ethereum/common"
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

// Set via go build -ldflags="-X main.version=${SEMVER}...
var version string

func main() {

	err := commonversion.SetVersion(version)
	if err != nil {
		log.Printf("Version string \"%s\" is invalid, valling back to hard coded version", version)
	}

	semver, err := commonversion.CurrentVersion()
	enforce.NilError(err, "invalid current version")

	log.Printf("Starting EigenDA Validator, version %s", semver)

	app := cli.NewApp()
	app.Flags = flags.Flags

	app.Version = semver.String()
	app.Name = node.AppName
	app.Usage = "EigenDA Node"
	app.Description = "Service for receiving and storing encoded blobs from disperser"

	app.Action = NodeMain
	err = app.Run(os.Args)
	if err != nil {
		log.Fatalf("application failed: %v", err)
	}

	select {}
}

func NodeMain(ctx *cli.Context) error {

	// TODO (cody.littley): pull all business logic in this function into the NewNode() constructor.

	log.Println("Initializing Node")
	config, err := node.NewConfig(ctx)
	if err != nil {
		return err
	}

	logger, err := common.NewLogger(&config.LoggerConfig)
	if err != nil {
		return err
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

	contractDirectory, err := directory.NewContractDirectory(
		context.Background(),
		logger,
		client,
		gethcommon.HexToAddress(config.EigenDADirectory))
	if err != nil {
		return fmt.Errorf("failed to create contract directory: %w", err)
	}

	operatorStateRetrieverAddress, err :=
		contractDirectory.GetContractAddress(context.Background(), directory.OperatorStateRetriever)
	if err != nil {
		return fmt.Errorf("failed to get OperatorStateRetriever address: %w", err)
	}

	eigenDAServiceManagerAddress, err :=
		contractDirectory.GetContractAddress(context.Background(), directory.ServiceManager)
	if err != nil {
		return fmt.Errorf("failed to get ServiceManager address: %w", err)
	}

	// Create and start the node.
	node, err := node.NewNode(context.Background(), reg, config, contractDirectory, pubIPProvider, client, logger)
	if err != nil {
		return err
	}

	// TODO(cody-littley): the metrics server is currently started by eigenmetrics, which is in another repo.
	//  When we fully remove v1 support, we need to start the metrics server inside the v2 metrics code.
	server := nodegrpc.NewServer(config, node, logger, ratelimiter)

	reader, err := coreeth.NewReader(
		logger,
		client,
		operatorStateRetrieverAddress.Hex(),
		eigenDAServiceManagerAddress.Hex())
	if err != nil {
		return fmt.Errorf("cannot create eth.Reader: %w", err)
	}

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
