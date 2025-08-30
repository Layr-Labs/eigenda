package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/tools/integration_utils/altdacommitment_parser"
	"github.com/Layr-Labs/eigenda/tools/integration_utils/flags"
	"github.com/Layr-Labs/eigenda/tools/integration_utils/gas_exhaustion_cert_meter"

	"github.com/urfave/cli"
)

var (
	version   = ""
	gitCommit = ""
	gitDate   = ""
)

const (
	FlagPrefix = ""
	envPrefix  = "INTEGRATION_UTILS"
)

func main() {
	app := cli.NewApp()
	app.Version = fmt.Sprintf("%s,%s,%s", version, gitCommit, gitDate)
	app.Name = "integration_utils"
	app.Description = "Integration utilities for EigenDA operations"
	app.Usage = "integration_utils <command> [command options]"
	app.Flags = common.LoggerCLIFlags(envPrefix, FlagPrefix)

	app.Commands = []cli.Command{
		{
			Name:  "parse-altdacommitment",
			Usage: "Parse and display EigenDA certificates from hex-encoded RLP strings",
			Description: "Parse and display EigenDA certificates from hex-encoded RLP strings. " +
				"Hex strings can be obtained from eigenda-proxy output or rollup inbox data. For OP rollups, " +
				"remove the '1' prefix byte from calldata before parsing.",
			Flags:  flags.ParserFlags,
			Action: altdacommitment_parser.DisplayAltDACommitmentFromHex,
		},
		{
			Name: "gas-exhaustion-cert-meter",
			Usage: "Estimates gas costs for verifying EigenDA certificates " +
				"when all operators are non-signers (worst case)\n\n",
			Description: "Gas estimation tool for EigenDA certificate verification in worst-case scenarios",
			Flags:       flags.GasExhaustionCertMeterFlags,
			Action:      gas_exhaustion_cert_meter.RunMeterer,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
