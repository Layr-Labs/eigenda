package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/disperser/dataapi"
	"github.com/Layr-Labs/eigenda/disperser/dataapi/subgraph"
	"github.com/Layr-Labs/eigenda/tools/ejections"
	"github.com/Layr-Labs/eigenda/tools/ejections/flags"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
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
	logger.Info("ejections", "subgraph", config.SubgraphEndpoint, "days", config.Days, "operatorId", config.OperatorId, "first", config.First, "skip", config.Skip)
	subgraphApi := subgraph.NewApi(config.SubgraphEndpoint, config.SubgraphEndpoint)
	subgraphClient := dataapi.NewSubgraphClient(subgraphApi, logger)

	ejections, err := subgraphClient.QueryOperatorEjectionsForTimeWindow(context.Background(), int32(config.Days), config.OperatorId, config.First, config.Skip)
	if err != nil {
		logger.Warn("failed to fetch operator ejections", "operatorId", config.OperatorId, "error", err)
		return errors.New("operator ejections not found")
	}

	rowConfigAutoMerge := table.RowConfig{AutoMerge: true}

	t := table.NewWriter()
	t.AppendHeader(table.Row{"OperatorID", "Quorum", "Timestamp", "Txn"}, rowConfigAutoMerge)

	sort.Slice(ejections, func(i, j int) bool {
		return ejections[i].OperatorId < ejections[j].OperatorId
	})
	for _, ejection := range ejections {
		logger.Debug("ejection", "ts", ejection.BlockTimestamp, "txn", ejection.TransactionHash, "quorum", ejection.Quorum, "operatorId", ejection.OperatorId)
		var link_prefix string
		if strings.Contains(config.SubgraphEndpoint, "holesky") {
			link_prefix = "https://holesky.etherscan.io/tx/"
		} else {
			link_prefix = "https://etherscan.io/tx/"
		}
		t.AppendRow(table.Row{ejection.OperatorId, ejection.Quorum, ejection.BlockTimestamp, link_prefix + ejection.TransactionHash}, rowConfigAutoMerge)
	}
	t.SetAutoIndex(true)
	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, AutoMerge: true},
		{Number: 2, Align: text.AlignCenter},
	})
	t.SetStyle(table.StyleLight)
	t.Style().Options.SeparateRows = true
	fmt.Println(t.Render())
	return nil
}
