package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/disperser/dataapi"
	"github.com/Layr-Labs/eigenda/disperser/dataapi/subgraph"
	"github.com/Layr-Labs/eigenda/tools/ejections"
	"github.com/Layr-Labs/eigenda/tools/ejections/flags"
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
	app.Name = "ejections"
	app.Description = "operator ejections scanner"
	app.Usage = ""
	app.Flags = flags.Flags
	app.Action = RunScan
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func RunScan(ctx *cli.Context) error {
	config, err := ejections.NewConfig(ctx)
	if err != nil {
		return err
	}

	logger, err := common.NewLogger(config.LoggerConfig)
	if err != nil {
		return err
	}
	logger.Info("starting ejections scanner", "subgraph", config.SubgraphEndpoint)
	subgraphApi := subgraph.NewApi(config.SubgraphEndpoint, config.SubgraphEndpoint)
	subgraphClient := dataapi.NewSubgraphClient(subgraphApi, logger)

	ejections, err := subgraphClient.QueryOperatorEjectionsForTimeWindow(context.Background(), int32(config.Days), config.OperatorId)
	if err != nil {
		logger.Warn("failed to fetch operator ejections", "operatorId", config.OperatorId, "error", err)
		return errors.New("operator ejections not found")
	}
	for _, ejection := range ejections {
		logger.Info("ejection", "operatorId", ejection.OperatorId, "timestamp", ejection.BlockTimestamp, "txnHash", ejection.TransactionHash, "quorum", ejection.Quorum)
	}

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
