package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"sort"
	"strings"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/eth"
	"github.com/Layr-Labs/eigenda/core/thegraph"
	"github.com/Layr-Labs/eigenda/tools/quorumscan"
	"github.com/Layr-Labs/eigenda/tools/quorumscan/flags"
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
	app.Name = "quorumscan"
	app.Description = "operator quorum scan"
	app.Usage = ""
	app.Flags = flags.Flags
	app.Action = RunScan
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func RunScan(ctx *cli.Context) error {
	config, err := quorumscan.NewConfig(ctx)
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

	var blockNumber uint
	if config.BlockNumber != 0 {
		blockNumber = uint(config.BlockNumber)
	} else {
		blockNumber, err = ics.GetCurrentBlockNumber()
		if err != nil {
			return fmt.Errorf("failed to fetch current block number - %s", err)
		}
	}
	logger.Info("Using block number", "block", blockNumber)

	operatorState, err := chainState.GetOperatorState(context.Background(), blockNumber, config.QuorumIDs)
	if err != nil {
		return fmt.Errorf("failed to fetch operator state - %s", err)
	}
	operators, err := ics.GetIndexedOperators(context.Background(), blockNumber)
	if err != nil {
		return fmt.Errorf("failed to fetch indexed operators info - %s", err)
	}

	logger.Info("Queried operator state", "count", len(operators))

	operatorIDs := make([]core.OperatorID, 0, len(operators))
	for opID := range operators {
		operatorIDs = append(operatorIDs, opID)
	}
	operatorAddresses, err := tx.BatchOperatorIDToAddress(context.Background(), operatorIDs)
	if err != nil {
		return err
	}
	operatorIdToAddress := make(map[string]string)
	for i := range operatorAddresses {
		operatorIdToAddress[operatorIDs[i].Hex()] = strings.ToLower(operatorAddresses[i].Hex())
	}

	quorumMetrics := quorumscan.QuorumScan(operators, operatorState, logger)
	displayResults(quorumMetrics, operatorIdToAddress)
	return nil
}

func displayResults(results map[uint8]*quorumscan.QuorumMetrics, operatorIdToAddress map[string]string) {
	weiToEth := new(big.Float).SetFloat64(1e18)

	// Create sorted list of quorums
	quorums := make([]uint8, 0, len(results))
	for quorum := range results {
		quorums = append(quorums, quorum)
	}
	sort.Slice(quorums, func(i, j int) bool {
		return quorums[i] < quorums[j]
	})

	for _, quorum := range quorums {
		tw := table.NewWriter()
		rowAutoMerge := table.RowConfig{AutoMerge: true}
		tw.AppendHeader(table.Row{"QUORUM", "OPERATOR", "ADDRESS", "STAKE", "STAKE PCT"}, rowAutoMerge)

		total_operators := 0
		total_stake_pct := 0.0
		total_stake := new(big.Float)
		metrics := results[quorum]

		// Create sorted list of operators by stake
		type operatorInfo struct {
			id    string
			stake float64
			pct   float64
		}
		operators := make([]operatorInfo, 0, len(metrics.OperatorStake))
		for op, stake := range metrics.OperatorStake {
			operators = append(operators, operatorInfo{op, stake, metrics.OperatorStakePct[op]})
		}
		sort.Slice(operators, func(i, j int) bool {
			return operators[i].stake > operators[j].stake
		})

		for _, op := range operators {
			stakeInEth := new(big.Float).Quo(new(big.Float).SetFloat64(op.stake), weiToEth)
			total_operators += 1
			total_stake.Add(total_stake, stakeInEth)
			total_stake_pct += op.pct

			tw.AppendRow(table.Row{quorum, op.id, operatorIdToAddress[op.id], stakeInEth, op.pct})
			//fmt.Printf("Quorum %d %s %.2f ETH %.2f%%\n", quorum, op.id, stakeInEth, op.pct)
		}
		tw.AppendFooter(table.Row{"TOTAL", total_operators, total_operators, total_stake, total_stake_pct})
		tw.SetColumnConfigs([]table.ColumnConfig{
			{Number: 1, AutoMerge: true},
		})
		fmt.Println(tw.Render())
	}
}
