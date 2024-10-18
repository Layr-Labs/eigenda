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
	version   = "1.0.0"
	gitCommit = ""
	gitDate   = ""
)

func main() {
	app := cli.NewApp()
	app.Version = fmt.Sprintf("%s,%s,%s", version, gitCommit, gitDate)
	app.Name = "ejections report"
	app.Description = "operator ejections report"
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

	operators := table.NewWriter()
	operators.AppendHeader(table.Row{"OperatorID", "Quorum", "Timestamp", "Txn"}, rowConfigAutoMerge)
	txns := table.NewWriter()
	txns.AppendHeader(table.Row{"Txn", "Timestamp", "OperatorId", "Quorum"}, rowConfigAutoMerge)

	sort.Slice(ejections, func(i, j int) bool {
		return ejections[i].BlockTimestamp > ejections[j].BlockTimestamp
	})

	for _, ejection := range ejections {
		var link_prefix string
		if strings.Contains(config.SubgraphEndpoint, "holesky") {
			link_prefix = "https://holesky.etherscan.io/tx/"
		} else {
			link_prefix = "https://etherscan.io/tx/"
		}
		operators.AppendRow(table.Row{ejection.OperatorId, ejection.Quorum, ejection.BlockTimestamp, link_prefix + ejection.TransactionHash}, rowConfigAutoMerge)
		txns.AppendRow(table.Row{link_prefix + ejection.TransactionHash, ejection.BlockTimestamp, ejection.OperatorId, ejection.Quorum}, rowConfigAutoMerge)
	}
	operators.SetAutoIndex(true)
	operators.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, AutoMerge: true},
		{Number: 2, Align: text.AlignCenter},
	})
	operators.SetStyle(table.StyleLight)
	operators.Style().Options.SeparateRows = true

	txns.SetAutoIndex(true)
	txns.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, AutoMerge: true},
		{Number: 2, AutoMerge: true},
		{Number: 3, AutoMerge: true},
		{Number: 4, Align: text.AlignCenter},
	})
	txns.SetStyle(table.StyleLight)
	txns.Style().Options.SeparateRows = true

	fmt.Println(operators.Render())
	fmt.Println(txns.Render())
	return nil
}
