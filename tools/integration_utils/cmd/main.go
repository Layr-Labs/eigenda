package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/tools/integration_utils/altdacommitment_parser"
	"github.com/Layr-Labs/eigenda/tools/integration_utils/calldata_gas_estimator"
	"github.com/Layr-Labs/eigenda/tools/integration_utils/flags"
	"github.com/Layr-Labs/eigenda/tools/integration_utils/gas_exhaustion_cert_meter"
	"github.com/Layr-Labs/eigenda/tools/integration_utils/validate_cert_verifier"

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
			Name:        "calldata-gas-estimator",
			Usage:       "Estimate EVM gas cost to send calldata containing AltDA commitment",
			Description: "Calculate EVM gas costs using EIP-2028 and EIP-7623 pricing models.",
			Flags:       flags.CallDataGasEstimatorFlags,
			Action:      calldata_gas_estimator.RunEstimator,
		},
		{
			Name: "gas-exhaustion-cert-meter",
			Usage: "Estimates gas costs for verifying EigenDA certificates " +
				"when all operators are non-signers (worst case)\n\n",
			Description: "Gas estimation tool for EigenDA certificate verification in worst-case scenarios",
			Flags:       flags.GasExhaustionCertMeterFlags,
			Action:      gas_exhaustion_cert_meter.RunMeterer,
		},
		{
			Name:        "validate-cert-verifier",
			Usage:       "Validate the CertVerifier contract by dispersing a blob and verifying the certificate",
			Description: "Disperses a test blob to EigenDA and validates that the CertVerifier contract correctly verifies the returned certificate using checkDACert",
			Flags:       flags.ValidateCertVerifierFlags,
			Action:      validate_cert_verifier.RunCreateAndValidateCertValidation,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
