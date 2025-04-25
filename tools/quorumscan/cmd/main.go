package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"sort"
	"strconv"
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
		log.Fatalln("could not start eth.NewReader", err)
	}
	chainState := eth.NewChainState(tx, gethClient)

	logger.Info("Connecting to subgraph", "url", config.ChainStateConfig.Endpoint)
	ics := thegraph.MakeIndexedChainState(config.ChainStateConfig, chainState, logger)

	var blockNumber uint
	if config.BlockNumber != 0 {
		blockNumber = uint(config.BlockNumber)
	} else {
		blockNumber, err = ics.GetCurrentBlockNumber(context.Background())
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

	// Handle file output if specified
	if config.OutputFile != "" {
		file, err := os.Create(config.OutputFile)
		if err != nil {
			return fmt.Errorf("failed to create output file: %v", err)
		}
		defer file.Close()

		err = displayResultsToWriter(quorumMetrics, operatorIdToAddress, config.TopN, config.OutputFormat, bufio.NewWriter(file))
		if err != nil {
			return fmt.Errorf("failed to write to output file: %v", err)
		}

		logger.Info("Output written to file", "path", config.OutputFile)
	} else {
		// Display to stdout
		displayResults(quorumMetrics, operatorIdToAddress, config.TopN, config.OutputFormat)
	}

	return nil
}

func humanizeEth(value *big.Float) string {
	v, _ := value.Float64()
	switch {
	case v >= 1_000_000:
		return fmt.Sprintf("%.2fM", v/1_000_000)
	case v >= 1_000:
		return fmt.Sprintf("%.2fK", v/1_000)
	default:
		return fmt.Sprintf("%.2f", v)
	}
}

// displayResults outputs to stdout
func displayResults(results map[uint8]*quorumscan.QuorumMetrics, operatorIdToAddress map[string]string, topN uint, outputFormat string) {
	// Use standard output
	writer := bufio.NewWriter(os.Stdout)
	err := displayResultsToWriter(results, operatorIdToAddress, topN, outputFormat, writer)
	if err != nil {
		log.Fatalf("Error writing to stdout: %v", err)
	}
}

// displayResultsToWriter outputs to the provided writer
func displayResultsToWriter(results map[uint8]*quorumscan.QuorumMetrics, operatorIdToAddress map[string]string, topN uint, outputFormat string, writer *bufio.Writer) error {
	weiToEth := new(big.Float).SetFloat64(1e18)

	// Create sorted list of quorums
	quorums := make([]uint8, 0, len(results))
	for quorum := range results {
		quorums = append(quorums, quorum)
	}
	sort.Slice(quorums, func(i, j int) bool {
		return quorums[i] < quorums[j]
	})

	// Get block number from the first quorum's metrics
	var blockNumber uint
	if len(results) > 0 {
		blockNumber = results[quorums[0]].BlockNumber
	}

	// Display block number at the top
	if outputFormat == "table" {
		_, err := writer.WriteString(fmt.Sprintf("Block Number: %d\n\n", blockNumber))
		if err != nil {
			return err
		}
	} else if outputFormat == "csv" {
		// Print CSV header with block number in first row
		_, err := writer.WriteString(fmt.Sprintf("BLOCK_NUMBER,%d\n", blockNumber))
		if err != nil {
			return err
		}
		_, err = writer.WriteString("QUORUM,OPERATOR,ADDRESS,SOCKET,STAKE,STAKE_PERCENTAGE\n")
		if err != nil {
			return err
		}
	} else {
		// For any other format, still display the block number
		_, err := writer.WriteString(fmt.Sprintf("Block Number: %d\n\n", blockNumber))
		if err != nil {
			return err
		}
	}

	for _, quorum := range quorums {
		var tw table.Writer
		if outputFormat == "table" {
			tw = table.NewWriter()
			rowAutoMerge := table.RowConfig{AutoMerge: true}
			operatorHeader := "OPERATOR"
			if topN > 0 {
				operatorHeader = "TOP " + strconv.Itoa(int(topN)) + " OPERATORS"
			}
			tw.AppendHeader(table.Row{"QUORUM", operatorHeader, "ADDRESS", "SOCKET", "STAKE", "STAKE"}, rowAutoMerge)
		}

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
			if topN > 0 && uint(total_operators) >= topN {
				break
			}
			stakeInEth := new(big.Float).Quo(new(big.Float).SetFloat64(op.stake), weiToEth)
			stakeInEth.SetPrec(64)
			total_operators += 1
			total_stake.Add(total_stake, stakeInEth)
			total_stake_pct += op.pct

			socket := metrics.OperatorSocket[op.id]
			if socket == "" {
				socket = "N/A"
			}

			if outputFormat == "csv" {
				_, err := writer.WriteString(fmt.Sprintf("%d,%s,%s,%s,%s,%.2f%%\n",
					quorum,
					op.id,
					operatorIdToAddress[op.id],
					socket,
					humanizeEth(stakeInEth),
					op.pct))
				if err != nil {
					return err
				}
			} else {
				tw.AppendRow(table.Row{quorum, op.id, operatorIdToAddress[op.id], socket, humanizeEth(stakeInEth), fmt.Sprintf("%.2f%%", op.pct)})
			}
		}

		if outputFormat == "table" {
			total_stake.SetPrec(64)
			tw.AppendFooter(table.Row{"TOTAL", total_operators, total_operators, total_operators, humanizeEth(total_stake), fmt.Sprintf("%.2f%%", total_stake_pct)})
			tw.SetColumnConfigs([]table.ColumnConfig{
				{Number: 1, AutoMerge: true},
			})
			_, err := writer.WriteString(tw.Render() + "\n")
			if err != nil {
				return err
			}
		} else if outputFormat == "csv" && total_operators > 0 {
			// Add total row for CSV
			_, err := writer.WriteString(fmt.Sprintf("TOTAL,%d,%d,%d,%s,%.2f%%\n",
				total_operators,
				total_operators,
				total_operators,
				humanizeEth(total_stake),
				total_stake_pct))
			if err != nil {
				return err
			}
		}
	}

	// Make sure to flush the writer to ensure all data is written
	return writer.Flush()
}
