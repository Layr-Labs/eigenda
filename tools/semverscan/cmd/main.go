package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/eth"
	"github.com/Layr-Labs/eigenda/core/thegraph"
	"github.com/Layr-Labs/eigenda/disperser/common/semver"
	"github.com/Layr-Labs/eigenda/tools/semverscan"
	"github.com/Layr-Labs/eigenda/tools/semverscan/flags"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/jedib0t/go-pretty/v6/table"
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
	app.Name = "semverscan"
	app.Description = "operator semver scan"
	app.Usage = ""
	app.Flags = flags.Flags
	app.Action = RunScan
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func RunScan(ctx *cli.Context) error {
	config, err := semverscan.NewConfig(ctx)
	if err != nil {
		return err
	}

	logger, err := common.NewLogger(config.LoggerConfig)
	if err != nil {
		return err
	}

	gethClient, err := geth.NewClient(config.EthClientConfig, gethcommon.Address{}, 0, logger)
	if err != nil {
		logger.Error("Cannot create chain.Client", "err", err)
		return err
	}

	tx, err := eth.NewReader(logger, gethClient, config.BLSOperatorStateRetrieverAddr, config.EigenDAServiceManagerAddr)
	if err != nil {
		log.Fatalln("could not start tcp listener", err)
	}
	cs := eth.NewChainState(tx, gethClient)

	logger.Info("Connecting to subgraph", "url", config.ChainStateConfig.Endpoint)
	ics := thegraph.MakeIndexedChainState(config.ChainStateConfig, cs, logger)

	currentBlock, err := ics.GetCurrentBlockNumber()
	if err != nil {
		return fmt.Errorf("failed to fetch current block number - %s", err)
	}
	operatorState, err := ics.GetIndexedOperatorState(context.Background(), currentBlock, []core.QuorumID{0, 1, 2})
	if err != nil {
		return fmt.Errorf("failed to fetch indexed operator state - %s", err)
	}
	logger.Info("Queried operator state", "count", len(operatorState.IndexedOperators))

	semvers := semver.ScanOperators(operatorState.IndexedOperators, config.Workers, config.Timeout, logger)
	displayResults(semvers)
	return nil
}

func displayResults(results map[string]int) {
	tw := table.NewWriter()

	rowHeader := table.Row{"semver", "count"}
	tw.AppendHeader(rowHeader)

	total := 0
	for semver, count := range results {
		tw.AppendRow(table.Row{semver, count})
		total += count
	}
	tw.AppendFooter(table.Row{"total", total})

	fmt.Println(tw.Render())
}
