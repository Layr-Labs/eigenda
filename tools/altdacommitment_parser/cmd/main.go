package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Layr-Labs/eigenda/tools/altdacommitment_parser"
	"github.com/Layr-Labs/eigenda/tools/altdacommitment_parser/flags"
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
	app.Name = "altdacommitment_parser"
	app.Description = "Parse and display EigenDA certificates from hex-encoded RLP strings. " +
		"Hex strings can be obtained from eigenda-proxy output or rollup inbox data. For OP rollups, " +
		"remove the '1' prefix byte from calldata before parsing."
	app.Usage = "altdacommitment_parser --hex <hex-encoded-cert-string>"
	app.Flags = flags.Flags
	app.Action = ParseCert
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func ParseCert(ctx *cli.Context) error {
	hexString := ctx.String("hex")
	if hexString == "" {
		return fmt.Errorf("hex string is required")
	}

	return ParseCertFromHex(hexString)
}

// ParseCertFromHex parses an EigenDA certificate from a hex-encoded RLP string
// and prints a nicely formatted display of its contents to stdout
func ParseCertFromHex(hexString string) error {
	// Use the parser library to parse the certificate
	commitment, versionedCert, err := altdacommitment_parser.ParseCertFromHex(hexString)
	if err != nil {
		return fmt.Errorf("failed to parse cert prefix: %w", err)
	}

	// Display the parsed commitment information
	altdacommitment_parser.DisplayCommitmentInfo(commitment)

	certV3, err := altdacommitment_parser.ParseCertificateData(versionedCert)
	if err != nil {
		return fmt.Errorf("failed to parse certificate data: %w", err)
	}

	// Display the certificate data
	altdacommitment_parser.DisplayCertificateData(commitment, certV3)
	return nil
}
