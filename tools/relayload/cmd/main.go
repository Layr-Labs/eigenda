package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/tools/relayload"
	"github.com/Layr-Labs/eigenda/tools/relayload/flags"
	"github.com/urfave/cli"
)

var (
	version   = ""
	gitCommit = ""
	gitDate   = ""
)

func main() {
	app := cli.NewApp()
	app.Version = fmt.Sprintf("%s,%s,%s", version, gitCommit, gitDate)
	app.Name = "relayload"
	app.Description = "relay load testing tool"
	app.Usage = ""
	app.Flags = flags.Flags
	app.Action = Run
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func Run(cliCtx *cli.Context) error {
	config, err := relayload.NewConfig(cliCtx)
	if err != nil {
		return err
	}

	logger, err := common.NewLogger(config.LoggerConfig)
	if err != nil {
		return err
	}

	relayClient := relayload.NewRelayLoad(config.RelayUrl, config.DataApiUrl, config.RangeSizes, config.RequestSizes, config.NumThreads, logger)

	// Create a context that can be cancelled
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle Ctrl+C
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	go func() {
		<-sigChan
		cancel()
	}()

	// Run parallel requests
	return relayClient.RunParallel(ctx, config.OperatorId, config.NumThreads)
}
