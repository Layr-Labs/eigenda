package main

import (
	"context"
	"errors"
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
	"github.com/Layr-Labs/eigenda/disperser/dataapi"
	"github.com/Layr-Labs/eigenda/disperser/dataapi/subgraph"
	"github.com/Layr-Labs/eigenda/tools/ejections"
	"github.com/Layr-Labs/eigenda/tools/ejections/flags"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/urfave/cli"

	gethcommon "github.com/ethereum/go-ethereum/common"
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

	client, err := geth.NewMultiHomingClient(config.EthClientConfig, gethcommon.Address{}, logger)
	if err != nil {
		return err
	}

	tx, err := eth.NewTransactor(logger, client, config.BLSOperatorStateRetrieverAddr, config.EigenDAServiceManagerAddr)
	if err != nil {
		return err
	}

	chainState := eth.NewChainState(tx, client)
	if chainState == nil {
		return errors.New("failed to create chain state")
	}
	subgraphApi := subgraph.NewApi(config.SubgraphEndpoint, config.SubgraphEndpoint)
	subgraphClient := dataapi.NewSubgraphClient(subgraphApi, logger)

	ejections, err := subgraphClient.QueryOperatorEjectionsForTimeWindow(context.Background(), int32(config.Days), config.OperatorId, config.First, config.Skip)
	if err != nil {
		logger.Warn("failed to fetch operator ejections", "operatorId", config.OperatorId, "error", err)
		return errors.New("operator ejections not found")
	}

	rowConfigAutoMerge := table.RowConfig{AutoMerge: true}

	operators := table.NewWriter()
	operators.AppendHeader(table.Row{"Operator Address", "Quorum", "Stake %", "Timestamp", "Txn"}, rowConfigAutoMerge)
	txns := table.NewWriter()
	txns.AppendHeader(table.Row{"Txn", "Timestamp", "Operator Address", "Quorum", "Stake %"}, rowConfigAutoMerge)

	sort.Slice(ejections, func(i, j int) bool {
		return ejections[i].BlockTimestamp > ejections[j].BlockTimestamp
	})

	stateCache := make(map[uint64]*core.OperatorState)
	ejectedOperatorIds := make(map[core.OperatorID]struct{})
	for _, ejection := range ejections {
		previouseBlock := ejection.BlockNumber - 1
		if _, exists := stateCache[previouseBlock]; !exists {
			state, err := chainState.GetOperatorState(context.Background(), uint(previouseBlock), []uint8{0, 1})
			if err != nil {
				return err
			}
			stateCache[previouseBlock] = state
		}

		// construct a set of ejected operator ids for later batch address lookup
		opID, err := core.OperatorIDFromHex(ejection.OperatorId)
		if err != nil {
			return err
		}
		ejectedOperatorIds[opID] = struct{}{}
	}

	// resolve operator id to operator addresses mapping
	operatorIDs := make([]core.OperatorID, 0, len(ejectedOperatorIds))
	for opID := range ejectedOperatorIds {
		operatorIDs = append(operatorIDs, opID)
	}
	operatorAddresses, err := tx.BatchOperatorIDToAddress(context.Background(), operatorIDs)
	if err != nil {
		return err
	}
	operatorIdToAddress := make(map[string]string)
	for i := range operatorAddresses {
		operatorIdToAddress["0x"+operatorIDs[i].Hex()] = strings.ToLower(operatorAddresses[i].Hex())
	}

	for _, ejection := range ejections {
		state := stateCache[ejection.BlockNumber-1]
		opID, err := core.OperatorIDFromHex(ejection.OperatorId)
		if err != nil {
			return err
		}

		stakePercentage := float64(0)
		if stake, ok := state.Operators[ejection.Quorum][opID]; ok {
			totalStake := new(big.Float).SetInt(state.Totals[ejection.Quorum].Stake)
			operatorStake := new(big.Float).SetInt(stake.Stake)
			stakePercentage, _ = new(big.Float).Mul(big.NewFloat(100), new(big.Float).Quo(operatorStake, totalStake)).Float64()
		}

		operatorAddress := operatorIdToAddress[ejection.OperatorId]
		operators.AppendRow(table.Row{operatorAddress, ejection.Quorum, stakePercentage, ejection.BlockTimestamp, ejection.TransactionHash}, rowConfigAutoMerge)
		txns.AppendRow(table.Row{ejection.TransactionHash, ejection.BlockTimestamp, operatorAddress, ejection.Quorum, stakePercentage}, rowConfigAutoMerge)
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
