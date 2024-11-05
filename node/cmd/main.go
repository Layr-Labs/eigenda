package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

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
	app.Version = fmt.Sprintf("%s-%s-%s", node.SemVer, node.GitCommit, node.GitDate)
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

	logger, err := common.NewLogger(config.LoggerConfig)
	if err != nil {
		return err
	}

	pubIPProvider := pubip.ProviderOrDefault(config.PubIPProvider)

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

	// Create the node.
	node, err := node.NewNode(reg, config, pubIPProvider, logger)
	if err != nil {
		return err
	}

	err = node.Start(context.Background())
	if err != nil {
		node.Logger.Error("could not start node", "error", err)
		return err
	}

	// Creates the GRPC server.
	server := nodegrpc.NewServer(config, node, logger, ratelimiter)
	serverV2 := nodegrpc.NewServerV2(config, node, logger, ratelimiter)
	err = nodegrpc.RunServers(server, serverV2, config, logger)

	return err
}
