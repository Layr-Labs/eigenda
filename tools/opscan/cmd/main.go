package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/disperser/common/semver"
	"github.com/Layr-Labs/eigenda/disperser/dataapi"
	"github.com/Layr-Labs/eigenda/disperser/dataapi/subgraph"
	"github.com/Layr-Labs/eigenda/tools/opscan"
	"github.com/Layr-Labs/eigenda/tools/opscan/flags"
	"github.com/Layr-Labs/eigensdk-go/logging"
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
	app.Name = "opscan"
	app.Description = "operator network scanner"
	app.Usage = ""
	app.Flags = flags.Flags
	app.Action = RunScan
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func RunScan(ctx *cli.Context) error {
	config, err := opscan.NewConfig(ctx)
	if err != nil {
		return err
	}

	logger, err := common.NewLogger(config.LoggerConfig)
	if err != nil {
		return err
	}

	subgraphApi := subgraph.NewApi(config.SubgraphEndpoint, config.SubgraphEndpoint)
	subgraphClient := dataapi.NewSubgraphClient(subgraphApi, logger)

	activeOperators := make([]string, 0)
	if config.OperatorId != "" {
		activeOperators = append(activeOperators, config.OperatorId)
	} else {
		registrations, err := subgraphApi.QueryOperators(context.Background(), 10000)
		if err != nil {
			return fmt.Errorf("failed to fetch indexed operator state - %s", err)
		}
		deregistrations, err := subgraphApi.QueryOperatorsDeregistered(context.Background(), 10000)
		if err != nil {
			return fmt.Errorf("failed to fetch indexed operator state - %s", err)
		}

		// Count registrations
		operators := make(map[string]int)
		for _, registration := range registrations {
			logger.Info("Operator", "operatorId", string(registration.OperatorId), "info", registration)
			operators[string(registration.OperatorId)]++
		}

		// Count deregistrations
		for _, deregistration := range deregistrations {
			operators[string(deregistration.OperatorId)]--
		}

		for operatorId, count := range operators {
			if count > 0 {
				activeOperators = append(activeOperators, operatorId)
			}
		}
		logger.Info("Active operators", "count", len(activeOperators))
	}

	semvers := scanOperators(subgraphClient, activeOperators, config, logger)
	displayResults(semvers)
	return nil
}

func getOperatorInfo(subgraphClient dataapi.SubgraphClient, operatorId string, logger logging.Logger) (*core.IndexedOperatorInfo, error) {
	operatorInfo, err := subgraphClient.QueryOperatorInfoByOperatorId(context.Background(), operatorId)
	if err != nil {
		logger.Warn("failed to fetch operator info", "operatorId", operatorId, "error", err)
		return nil, fmt.Errorf("operator info not found for operatorId %s", operatorId)
	}
	return operatorInfo, nil
}

func scanOperators(subgraphClient dataapi.SubgraphClient, operatorIds []string, config *opscan.Config, logger logging.Logger) map[string]int {
	var wg sync.WaitGroup
	var mu sync.Mutex
	semvers := make(map[string]int)
	operatorChan := make(chan string, len(operatorIds))
	numWorkers := 10 // Adjust the number of workers as needed
	worker := func() {
		for operatorId := range operatorChan {
			operatorInfo, err := getOperatorInfo(subgraphClient, operatorId, logger)
			if err != nil {
				mu.Lock()
				semvers["not-found"]++
				mu.Unlock()
				continue
			}
			operatorSocket := core.OperatorSocket(operatorInfo.Socket)
			retrievalSocket := operatorSocket.GetRetrievalSocket()
			semver := semver.GetSemverInfo(context.Background(), retrievalSocket, operatorId, true, logger)

			mu.Lock()
			semvers[semver]++
			mu.Unlock()
		}
		wg.Done()
	}

	// Launch worker goroutines
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker()
	}

	// Send operator IDs to the channel
	for _, operatorId := range operatorIds {
		operatorChan <- operatorId
	}
	close(operatorChan)

	// Wait for all workers to finish
	wg.Wait()
	return semvers
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
