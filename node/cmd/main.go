package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/common/pubip"

	"github.com/urfave/cli"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/logging"
	"github.com/Layr-Labs/eigenda/common/ratelimit"
	"github.com/Layr-Labs/eigenda/common/store"
	"github.com/Layr-Labs/eigenda/node"
	"github.com/Layr-Labs/eigenda/node/flags"
	"github.com/Layr-Labs/eigenda/node/grpc"
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

func ExtraFlagsValidation(ctx *cli.Context) error {
	testMode := ctx.GlobalBool(flags.EnableTestModeFlag.Name)

	blsKeyFile := ctx.GlobalString(flags.BlsKeyFileFlag.Name)
	blsKeyPassword := ctx.GlobalString(flags.BlsKeyPasswordFlag.Name)
	blsKey := ctx.GlobalString(flags.TestPrivateBlsFlag.Name)

	if !testMode && blsKey != "" {
		_ = cli.ShowAppHelp(ctx)
		return errors.New("may not pass BLS private key in plaintext in production mode")
	}

	if blsKey == "" && (blsKeyFile == "" || blsKeyPassword == "") {
		_ = cli.ShowAppHelp(ctx)
		if testMode {
			return errors.New("in test mode, must pass either a BLS private key OR a BLS encrypted private key file and the password to that file")
		} else {
			return errors.New("in prod mode, must pass a BLS encrypted private key file and the password to that file")
		}
	}

	ecdsaKeyFile := ctx.GlobalString(flags.EcdsaKeyFileFlag.Name)
	ecdsaKeyPassword := ctx.GlobalString(flags.EcdsaKeyPasswordFlag.Name)
	ecdsaKey := ctx.GlobalString(geth.PrivateKeyFlagName)

	if !testMode && ecdsaKey != "" {
		_ = cli.ShowAppHelp(ctx)
		return errors.New("may not pass ECDSA private key in plaintext in production mode")
	}

	if ecdsaKey == "" && (ecdsaKeyFile == "" || ecdsaKeyPassword == "") {
		_ = cli.ShowAppHelp(ctx)
		if testMode {
			return errors.New("in test mode, must pass either a ECDSA private key OR a ECDSA encrypted private key file and the password to that file")
		} else {
			return errors.New("in prod mode, must pass a ECDSA encrypted private key file and the password to that file")
		}
	}

	return nil
}

func NodeMain(ctx *cli.Context) error {
	log.Println("Initializing Node")
	err := ExtraFlagsValidation(ctx)
	if err != nil {
		return err
	}

	config, err := node.NewConfig(ctx)
	if err != nil {
		return err
	}

	logger, err := logging.GetLogger(config.LoggingConfig)
	if err != nil {
		return err
	}

	pubIPProvider := pubip.ProviderOrDefault(config.PubIPProvider)

	// Create the node.
	node, err := node.NewNode(config, pubIPProvider, logger)
	if err != nil {
		return err
	}

	err = node.Start(context.Background())
	if err != nil {
		node.Logger.Error("could not start node", "error", err)
		return err
	}

	globalParams := common.GlobalRateParams{
		BucketSizes: []time.Duration{bucketDuration},
		Multipliers: []float32{bucketMultiplier},
		CountFailed: true,
	}

	bucketStore, err := store.NewLocalParamStore[common.RateBucketParams](bucketStoreSize)
	if err != nil {
		return err
	}

	ratelimiter := ratelimit.NewRateLimiter(globalParams, bucketStore, logger)

	// Creates the GRPC server.
	server := grpc.NewServer(config, node, logger, ratelimiter)
	server.Start()

	return nil
}
