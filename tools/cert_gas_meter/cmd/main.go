package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Layr-Labs/eigenda/api/clients/v2/coretypes"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/tools/cert_gas_meter"
	"github.com/Layr-Labs/eigenda/tools/cert_gas_meter/flags"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
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
	app.Name = "cert-gas-meter"
	app.Description = "a worst case gas meter for a DA cert"
	app.Usage = ""
	app.Flags = flags.Flags
	app.Action = RunMeterer
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func RunMeterer(ctx *cli.Context) error {
	config, err := cert_gas_meter.NewConfig(ctx)
	if err != nil {
		return err
	}

	// Read the file
	data, err := os.ReadFile(config.CertPath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", config.CertPath, err)
	}
	// Try to parse as V3 cert first
	var certV3 coretypes.EigenDACertV3
	err = rlp.DecodeBytes(data, &certV3)

	fmt.Printf("certV3 %v", certV3)

	logger, err := common.NewLogger(&config.LoggerConfig)
	if err != nil {
		return err
	}

	gethClient, err := geth.NewClient(config.EthClientConfig, gethcommon.Address{}, 0, logger)
	if err != nil {
		logger.Error("Cannot create chain.Client", "err", err)
		return err
	}

	blockNumber := certV3.BatchHeader.ReferenceBlockNumber
	quorumIDsBytes := certV3.SignedQuorumNumbers
	fmt.Println("quorumIDsBytes", quorumIDsBytes)

	// Copied code below

	_, err = cert_gas_meter.OperatorIDsScan(
		quorumIDsBytes,
		blockNumber,
		gethClient,
		logger,
		certV3,
	)
	//indexedOperatorState, err := ics.GetIndexedOperatorState(context.Background(), blockNumber, quorumIDs)
	if err != nil {
		return err
	}

	return nil
}
