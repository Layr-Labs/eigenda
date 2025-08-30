package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Layr-Labs/eigenda/api/clients/v2/coretypes"
	"github.com/Layr-Labs/eigenda/tools/altdacommitment_parser"
	"github.com/Layr-Labs/eigenda/tools/gas_exhaustion_cert_meter"
	"github.com/Layr-Labs/eigenda/tools/gas_exhaustion_cert_meter/flags"

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
	app.Name = "gas-exhaustion-cert-meter"
	app.Description = "Gas estimation tool for EigenDA V3 certificate verification in worst-case scenarios"
	app.Usage = "Estimates gas costs for verifying EigenDA V3 certificates " +
		"when all operators are non-signers (worst case)\n\n" +
		"REQUIREMENTS:\n" +
		"  - RLP-serialized EigenDA V3 certificate file\n" +
		"  - Example certificates available: ./data/cert_v3.mainnet.rlp, ./data/cert_v3.sepolia.rlp\n\n" +
		"EXAMPLE:\n" +
		"  gas-exhaustion-cert-meter --eigenda-directory 0x... --cert-rlp-file ./data/cert_v3.mainnet.rlp" +
		"  --eth-rpc-url ..."
	app.Flags = flags.Flags
	app.Action = RunMeterer
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func RunMeterer(ctx *cli.Context) error {
	config, err := gas_exhaustion_cert_meter.NewConfig(ctx)
	if err != nil {
		return fmt.Errorf("failed to create config: %w", err)
	}

	// Read and decode the certificate file
	parseResult, err := altdacommitment_parser.ParseCertFromHex(config.CertHexString)
	if err != nil {
		return fmt.Errorf("failed to parse cert hex string: %w", err)
	}

	altdacommitment_parser.DisplayCommitmentInfo(parseResult.Commitment)

	var certV3 coretypes.EigenDACertV3
	if parseResult.CertV2 != nil {
		certV3 = *parseResult.CertV2.ToV3()
	} else {
		certV3 = *parseResult.CertV3
	}

	if err = gas_exhaustion_cert_meter.EstimateGas(config, certV3); err != nil {
		return fmt.Errorf("gas estimation failed: %w", err)
	}

	return nil
}
