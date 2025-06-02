package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/eth"
	"github.com/Layr-Labs/eigenda/disperser/dataapi"
	"github.com/Layr-Labs/eigenda/disperser/dataapi/subgraph"
	dataapiv2 "github.com/Layr-Labs/eigenda/disperser/dataapi/v2"
	"github.com/Layr-Labs/eigenda/operators/ejector"
	"github.com/Layr-Labs/eigenda/tools/ejections"
	"github.com/Layr-Labs/eigenda/tools/ejections/flags"
	"github.com/Layr-Labs/eigensdk-go/logging"
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

type DataAPIClient interface {
	GetNonSigningRateV1(endTime time.Time, interval int64) (*dataapi.OperatorsNonsigningPercentage, error)
	GetNonSigningRateV2(endTime time.Time, interval int64) (*dataapiv2.OperatorsSigningInfoResponse, error)
}

type dataapiClient struct {
	dataapiURL string
	httpClient *http.Client
	logger     logging.Logger
}

var _ DataAPIClient = (*dataapiClient)(nil)

func NewDataAPIClient(dataapiURL string, httpClient *http.Client, logger logging.Logger) DataAPIClient {
	return &dataapiClient{
		dataapiURL: dataapiURL,
		httpClient: httpClient,
		logger:     logger,
	}
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

	if config.Eval {
		return EvaluateOperators(config)
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

func EvaluateOperators(config *ejections.Config) error {
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

	httpClient := &http.Client{
		Timeout: 20 * time.Second,
	}

	dataapiClient := NewDataAPIClient(config.DataAPIURL, httpClient, logger)

	ctx := context.Background()
	bn, err := tx.GetCurrentBlockNumber(ctx)
	if err != nil {
		return fmt.Errorf("failed to get current block number: %w", err)
	}
	quorumCount, err := tx.GetQuorumCount(ctx, bn)
	if err != nil {
		return fmt.Errorf("failed to get quorum count: %w", err)
	}
	quorumIDs := make(map[core.QuorumID]struct{}, quorumCount)
	for i := uint8(0); i < quorumCount; i++ {
		quorumIDs[core.QuorumID(i)] = struct{}{}
	}
	logger.Info("Using all quorums", "QuorumIDs", quorumIDs)

	e := ejector.NewEjector(nil, client, logger, nil, nil, 0, config.NonsigningRateThreshold)

	nonsigningRateV1, err := dataapiClient.GetNonSigningRateV1(time.Now(), config.EvalInterval)
	if err != nil {
		return fmt.Errorf("failed to get v1 nonsigning rate: %w", err)
	}
	nonSigningMetricsV1Only := make([]*ejector.NonSignerMetric, 0)
	for _, metric := range nonsigningRateV1.Data {
		if _, ok := quorumIDs[core.QuorumID(metric.QuorumId)]; ok {
			nonSigningMetricsV1Only = append(nonSigningMetricsV1Only, &ejector.NonSignerMetric{
				OperatorId:           metric.OperatorId,
				OperatorAddress:      metric.OperatorAddress,
				QuorumId:             metric.QuorumId,
				TotalUnsignedBatches: metric.TotalUnsignedBatches,
				Percentage:           metric.Percentage,
				StakePercentage:      metric.StakePercentage,
			})
		}
	}
	err = e.EvaluateOperatorsForEjection(nonSigningMetricsV1Only)
	if err != nil {
		return fmt.Errorf("failed to evaluate operators for V1 ejection: %w", err)
	}
	sort.Slice(nonSigningMetricsV1Only, func(i, j int) bool {
		return nonSigningMetricsV1Only[i].OperatorAddress < nonSigningMetricsV1Only[j].OperatorAddress
	})
	// Create a new table writer for non-signing metrics
	nonSigningTableV1 := table.NewWriter()
	nonSigningTableV1.AppendHeader(table.Row{"Operator Address", "Quorum ID", "Unsigned Batches", "V1 Non Signing %", "Stake %", "Perf Score", "Violating SLA", "Violating Threshold", "Needs Ejection"}, table.RowConfig{AutoMerge: true})
	// Set the column configuration to merge the Operator Address column
	nonSigningTableV1.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, AutoMerge: true}, // Merging the Operator Address column
		{Number: 2, AutoMerge: false},
		{Number: 3, AutoMerge: false},
		{Number: 4, AutoMerge: false},
		{Number: 5, AutoMerge: false},
		{Number: 6, AutoMerge: false},
		{Number: 7, AutoMerge: false},
		{Number: 8, AutoMerge: false},
		{Number: 9, AutoMerge: false},
	})
	nonSigningTableV1.SetStyle(table.StyleLight)
	nonSigningTableV1.Style().Options.SeparateRows = true
	for _, metric := range nonSigningMetricsV1Only {
		nonSigningTableV1.AppendRow(table.Row{
			metric.OperatorAddress,
			metric.QuorumId,
			metric.TotalUnsignedBatches,
			metric.Percentage,
			metric.StakePercentage,
			metric.PerfScore,
			metric.IsViolatingSLA,
			metric.IsViolatingThreshold,
			metric.NeedsEjection,
		})
	}

	// Print the non-signing metrics table
	fmt.Println(nonSigningTableV1.Render())

	// V2 Only signing metrics
	nonSigningMetricsV2Only := make([]*ejector.NonSignerMetric, 0)
	if config.EvalV2 {
		nonsigningRateV2, err := dataapiClient.GetNonSigningRateV2(time.Now(), config.EvalInterval)
		if err != nil {
			return fmt.Errorf("failed to get v2 nonsigning rate: %w", err)
		}
		for _, metric := range nonsigningRateV2.OperatorSigningInfo {
			if _, ok := quorumIDs[core.QuorumID(metric.QuorumId)]; ok {
				nonSigningMetricsV2Only = append(nonSigningMetricsV2Only, &ejector.NonSignerMetric{
					OperatorId:           metric.OperatorId,
					OperatorAddress:      metric.OperatorAddress,
					QuorumId:             metric.QuorumId,
					TotalUnsignedBatches: metric.TotalUnsignedBatches,
					Percentage:           100 - metric.SigningPercentage,
					StakePercentage:      metric.StakePercentage,
				})
			}
		}
		err = e.EvaluateOperatorsForEjection(nonSigningMetricsV2Only)
		if err != nil {
			return fmt.Errorf("failed to evaluate operators for v2 ejection: %w", err)
		}
		sort.Slice(nonSigningMetricsV2Only, func(i, j int) bool {
			return nonSigningMetricsV2Only[i].OperatorAddress < nonSigningMetricsV2Only[j].OperatorAddress
		})

		// Create a new table writer for non-signing metrics
		nonSigningTableV2 := table.NewWriter()
		nonSigningTableV2.AppendHeader(table.Row{"Operator Address", "Quorum ID", "Unsigned Batches", "V2 Non Signing %", "Stake %", "Perf Score", "Violating SLA", "Violating Threshold", "Needs Ejection"}, table.RowConfig{AutoMerge: true})
		// Set the column configuration to merge the Operator Address column
		nonSigningTableV2.SetColumnConfigs([]table.ColumnConfig{
			{Number: 1, AutoMerge: true}, // Merging the Operator Address column
			{Number: 2, AutoMerge: false},
			{Number: 3, AutoMerge: false},
			{Number: 4, AutoMerge: false},
			{Number: 5, AutoMerge: false},
			{Number: 6, AutoMerge: false},
			{Number: 7, AutoMerge: false},
			{Number: 8, AutoMerge: false},
			{Number: 9, AutoMerge: false},
		})
		nonSigningTableV2.SetStyle(table.StyleLight)
		nonSigningTableV2.Style().Options.SeparateRows = true
		for _, metric := range nonSigningMetricsV2Only {
			nonSigningTableV2.AppendRow(table.Row{
				metric.OperatorAddress,
				metric.QuorumId,
				metric.TotalUnsignedBatches,
				metric.Percentage,
				metric.StakePercentage,
				metric.PerfScore,
				metric.IsViolatingSLA,
				metric.IsViolatingThreshold,
				metric.NeedsEjection,
			})
		}
		fmt.Println(nonSigningTableV2.Render())

		// Merge non-signing metrics from V1 and V2
		nonSigningMetricsMerged := e.MergeNonSigningMetrics(nonsigningRateV1, nonsigningRateV2, quorumIDs, logger)
		err = e.EvaluateOperatorsForEjection(nonSigningMetricsMerged)
		if err != nil {
			return fmt.Errorf("failed to evaluate operators for hybrid ejection: %w", err)
		}
		sort.Slice(nonSigningMetricsMerged, func(i, j int) bool {
			return nonSigningMetricsMerged[i].OperatorAddress < nonSigningMetricsMerged[j].OperatorAddress
		})
		nonSigningTableMerged := table.NewWriter()
		nonSigningTableMerged.AppendHeader(table.Row{"Operator Address", "Quorum ID", "Unsigned Batches", "Hybrid Non Signing %", "Stake %", "Perf Score", "Violating SLA", "Violating Threshold", "Needs Ejection"}, table.RowConfig{AutoMerge: true})
		for _, metric := range nonSigningMetricsMerged {
			nonSigningTableMerged.AppendRow(table.Row{metric.OperatorAddress, metric.QuorumId, metric.TotalUnsignedBatches, metric.Percentage, metric.StakePercentage, metric.PerfScore, metric.IsViolatingSLA, metric.IsViolatingThreshold, metric.NeedsEjection})
		}
		nonSigningTableMerged.SetStyle(table.StyleLight)
		nonSigningTableMerged.Style().Options.SeparateRows = true
		nonSigningTableMerged.SetColumnConfigs([]table.ColumnConfig{
			{Number: 1, AutoMerge: true}, // Merging the Operator Address column
			{Number: 2, AutoMerge: false},
			{Number: 3, AutoMerge: false},
			{Number: 4, AutoMerge: false},
			{Number: 5, AutoMerge: false},
			{Number: 6, AutoMerge: false},
			{Number: 7, AutoMerge: false},
			{Number: 8, AutoMerge: false},
			{Number: 9, AutoMerge: false},
		})
		fmt.Println(nonSigningTableMerged.Render())
	}

	return nil
}

func (c *dataapiClient) GetNonSigningRateV1(endTime time.Time, interval int64) (*dataapi.OperatorsNonsigningPercentage, error) {
	path := "api/v1/metrics/operator-nonsigning-percentage"
	urlStr, err := url.JoinPath(c.dataapiURL, path)
	if err != nil {
		return nil, fmt.Errorf("error joining URL path with %s and %s: %w", c.dataapiURL, path, err)
	}
	url, err := url.Parse(urlStr)
	if err != nil {
		return nil, fmt.Errorf("error parsing URL: %w", err)
	}
	// add query parameters
	q := url.Query()
	// end: datetime formatted in "2006-01-02T15:04:05Z"
	q.Set("end", endTime.UTC().Format("2006-01-02T15:04:05Z"))
	// interval: lookback window in seconds
	q.Set("interval", strconv.Itoa(int(interval)))
	url.RawQuery = q.Encode()
	c.logger.Info("making request to DataAPI", "url", url.String())

	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending HTTP request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}
	c.logger.Debug("Received response", "responseBody", string(respBody))

	if resp.StatusCode != http.StatusOK {
		var errResp dataapi.ErrorResponse
		err = json.Unmarshal(respBody, &errResp)
		if err != nil {
			return nil, fmt.Errorf("error parsing error response: %w", err)
		}
		return nil, fmt.Errorf(
			"error response (%d) from dataapi: %s",
			resp.StatusCode,
			errResp.Error,
		)
	}

	var response dataapi.OperatorsNonsigningPercentage
	err = json.NewDecoder(strings.NewReader(string(respBody))).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("error parsing response body: %w", err)
	}
	return &response, nil
}

func (c *dataapiClient) GetNonSigningRateV2(endTime time.Time, interval int64) (*dataapiv2.OperatorsSigningInfoResponse, error) {
	path := "api/v2/operators/signing-info"
	urlStr, err := url.JoinPath(c.dataapiURL, path)
	if err != nil {
		return nil, fmt.Errorf("error joining URL path with %s and %s: %w", c.dataapiURL, path, err)
	}
	url, err := url.Parse(urlStr)
	if err != nil {
		return nil, fmt.Errorf("error parsing URL: %w", err)
	}
	// add query parameters
	q := url.Query()
	// end: datetime formatted in "2006-01-02T15:04:05Z"
	q.Set("end", endTime.UTC().Format("2006-01-02T15:04:05Z"))
	// interval: lookback window in seconds
	q.Set("interval", strconv.Itoa(int(interval)))
	// nonsigner_only: whether to only return operators with signing rate less than 100%
	q.Set("nonsigner_only", "true")
	url.RawQuery = q.Encode()
	c.logger.Info("making request to DataAPI", "url", url.String())

	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending HTTP request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}
	c.logger.Debug("Received response", "responseBody", string(respBody))

	if resp.StatusCode != http.StatusOK {
		var errResp dataapi.ErrorResponse
		err = json.Unmarshal(respBody, &errResp)
		if err != nil {
			return nil, fmt.Errorf("error parsing error response: %w", err)
		}
		return nil, fmt.Errorf(
			"error response (%d) from dataapi: %s",
			resp.StatusCode,
			errResp.Error,
		)
	}

	var response dataapiv2.OperatorsSigningInfoResponse
	err = json.NewDecoder(strings.NewReader(string(respBody))).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("error parsing response body: %w", err)
	}
	return &response, nil
}
