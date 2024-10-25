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

type EjectionTransaction struct {
	BlockNumber           uint64            `json:"block_number"`
	BlockTimestamp        string            `json:"block_timestamp"`
	TransactionHash       string            `json:"transaction_hash"`
	QuorumStakePercentage map[uint8]float64 `json:"stake_percentage"`
	QuorumEjections       map[uint8]uint8   `json:"ejections"`
}

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

	tx, err := eth.NewReader(logger, client, config.BLSOperatorStateRetrieverAddr, config.EigenDAServiceManagerAddr)
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

	sort.Slice(ejections, func(i, j int) bool {
		return ejections[i].BlockTimestamp > ejections[j].BlockTimestamp
	})

	// Create a sorted slice from the set of quorums
	quorumSet := make(map[uint8]struct{})
	for _, ejection := range ejections {
		quorumSet[ejection.Quorum] = struct{}{}
	}
	quorums := make([]uint8, 0, len(quorumSet))
	for quorum := range quorumSet {
		quorums = append(quorums, quorum)
	}
	sort.Slice(quorums, func(i, j int) bool {
		return quorums[i] < quorums[j]
	})

	stateCache := make(map[uint64]*core.OperatorState)
	ejectedOperatorIds := make(map[core.OperatorID]struct{})
	for _, ejection := range ejections {
		previouseBlock := ejection.BlockNumber - 1
		if _, exists := stateCache[previouseBlock]; !exists {
			state, err := chainState.GetOperatorState(context.Background(), uint(previouseBlock), quorums)
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

	rowConfigAutoMerge := table.RowConfig{AutoMerge: true}
	rowConfigNoAutoMerge := table.RowConfig{AutoMerge: false}
	operators := table.NewWriter()
	operators.AppendHeader(table.Row{"Operator Address", "Quorum", "Stake %", "Timestamp", "Txn"}, rowConfigAutoMerge)
	txns := table.NewWriter()
	txns.AppendHeader(table.Row{"Txn", "Timestamp", "Operator Address", "Quorum", "Stake %"}, rowConfigAutoMerge)
	txnQuorums := table.NewWriter()
	txnQuorums.AppendHeader(table.Row{"Txn", "Timestamp", "Quorum", "Stake %", "Operators"}, rowConfigNoAutoMerge)

	ejectionTransactions := make(map[string]*EjectionTransaction)
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

		if _, exists := ejectionTransactions[ejection.TransactionHash]; !exists {
			ejectionTransactions[ejection.TransactionHash] = &EjectionTransaction{
				BlockNumber:           ejection.BlockNumber,
				BlockTimestamp:        ejection.BlockTimestamp,
				TransactionHash:       ejection.TransactionHash,
				QuorumStakePercentage: make(map[uint8]float64),
				QuorumEjections:       make(map[uint8]uint8),
			}
			ejectionTransactions[ejection.TransactionHash].QuorumStakePercentage[ejection.Quorum] = stakePercentage
			ejectionTransactions[ejection.TransactionHash].QuorumEjections[ejection.Quorum] = 1
		} else {
			ejectionTransactions[ejection.TransactionHash].QuorumStakePercentage[ejection.Quorum] += stakePercentage
			ejectionTransactions[ejection.TransactionHash].QuorumEjections[ejection.Quorum] += 1
		}

		operatorAddress := operatorIdToAddress[ejection.OperatorId]
		operators.AppendRow(table.Row{operatorAddress, ejection.Quorum, stakePercentage, ejection.BlockTimestamp, ejection.TransactionHash}, rowConfigAutoMerge)
		txns.AppendRow(table.Row{ejection.TransactionHash, ejection.BlockTimestamp, operatorAddress, ejection.Quorum, stakePercentage}, rowConfigAutoMerge)
	}

	orderedEjectionTransactions := make([]*EjectionTransaction, 0, len(ejectionTransactions))
	for _, txn := range ejectionTransactions {
		orderedEjectionTransactions = append(orderedEjectionTransactions, txn)
	}
	sort.Slice(orderedEjectionTransactions, func(i, j int) bool {
		return orderedEjectionTransactions[i].BlockNumber > orderedEjectionTransactions[j].BlockNumber
	})
	for _, txn := range orderedEjectionTransactions {
		for _, quorum := range quorums {
			if _, exists := txn.QuorumEjections[quorum]; exists {
				txnQuorums.AppendRow(table.Row{txn.TransactionHash, txn.BlockTimestamp, quorum, txn.QuorumStakePercentage[quorum], txn.QuorumEjections[quorum]}, rowConfigAutoMerge)
			}
		}
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

	txnQuorums.SetAutoIndex(true)
	txnQuorums.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, AutoMerge: true},
		{Number: 2, AutoMerge: true, Align: text.AlignCenter},
		{Number: 3, Align: text.AlignCenter},
		{Number: 5, Align: text.AlignCenter},
	})
	txnQuorums.SetStyle(table.StyleLight)
	txnQuorums.Style().Options.SeparateRows = true

	fmt.Println(operators.Render())
	fmt.Println(txns.Render())
	fmt.Println(txnQuorums.Render())
	return nil
}
