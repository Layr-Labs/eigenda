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
	chainState := eth.NewChainState(tx, gethClient)

	logger.Info("Connecting to subgraph", "url", config.ChainStateConfig.Endpoint)
	ics := thegraph.MakeIndexedChainState(config.ChainStateConfig, chainState, logger)

	currentBlock, err := ics.GetCurrentBlockNumber()
	if err != nil {
		return fmt.Errorf("failed to fetch current block number - %s", err)
	}
	operatorState, err := chainState.GetOperatorState(context.Background(), currentBlock, []core.QuorumID{0, 1, 2})
	if err != nil {
		return fmt.Errorf("failed to fetch operator state - %s", err)
	}
	operators, err := ics.GetIndexedOperators(context.Background(), currentBlock)
	if err != nil {
		return fmt.Errorf("failed to fetch indexed operators info - %s", err)
	}
	if config.OperatorId != "" {
		operatorId, err := core.OperatorIDFromHex(config.OperatorId)
		if err != nil {
			return fmt.Errorf("failed to parse operator id %s - %v", config.OperatorId, err)
		}
		for operator := range operators {
			if operator.Hex() != operatorId.Hex() {
				delete(operators, operator)
			}
		}
	}
	logger.Info("Queried operator state", "count", len(operators))

	semvers := semver.ScanOperators(operators, operatorState, config.UseRetrievalClient, config.Workers, config.Timeout, logger)
	for semver, metrics := range semvers {
		logger.Info("Semver Report", "semver", semver, "operators", metrics.Operators, "stake", metrics.QuorumStakePercentage)
	}
	displayResults(semvers)
	return nil
}

func displayResults(results map[string]*semver.SemverMetrics) {
	tw := table.NewWriter()
	rowAutoMerge := table.RowConfig{AutoMerge: true}
	tw.AppendHeader(table.Row{"semver", "install %", "operators", "quorum 0 stake %", "quorum 1 stake %", "quorum 2 stake %"}, rowAutoMerge)
	//tw.AppendHeader(table.Row{"", "", "quorum 0", "quorum 1", "quorum 2"})

	total_operators := 0
	total_semver_pct := 0.0
	total_stake_q0 := 0.0
	total_stake_q1 := 0.0
	total_stake_q2 := 0.0
	for _, metrics := range results {
		total_operators += int(metrics.Operators)
		total_stake_q0 += metrics.QuorumStakePercentage[0]
		total_stake_q1 += metrics.QuorumStakePercentage[1]
		total_stake_q2 += metrics.QuorumStakePercentage[2]
	}
	for semver, metrics := range results {
		semver_pct := 100 * (float64(metrics.Operators) / float64(total_operators))
		total_semver_pct += semver_pct
		tw.AppendRow(table.Row{semver, semver_pct, metrics.Operators, metrics.QuorumStakePercentage[0], metrics.QuorumStakePercentage[1], metrics.QuorumStakePercentage[2]})
	}
	tw.AppendFooter(table.Row{"totals", total_semver_pct, total_operators, total_stake_q0, total_stake_q1, total_stake_q2})
	tw.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, AutoMerge: true},
		{Number: 3, AlignHeader: 2},
		{Number: 4, AlignHeader: 2},
		{Number: 5, AlignHeader: 2},
	})

	fmt.Println(tw.Render())
}
