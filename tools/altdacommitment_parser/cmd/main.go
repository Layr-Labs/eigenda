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
	app.Action = ParseCertFromHex
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

// ParseCertFromHex parses an EigenDA certificate from a hex-encoded RLP string
// and prints a nicely formatted display of its contents to stdout
func ParseCertFromHex(ctx *cli.Context) error {
	hexString := ctx.GlobalString(flags.CertHexFlag.Name)

	// Use the parser library to parse the certificate
	prefix, versionedCert, err := altdacommitment_parser.ParseCertFromHex(hexString)
	if err != nil {
		return fmt.Errorf("failed to parse cert prefix: %w", err)
	}

	// Display the parsed prfix information
	altdacommitment_parser.DisplayPrefixInfo(prefix)

	certV3, err := altdacommitment_parser.ParseCertificateData(versionedCert)
	if err != nil {
		return fmt.Errorf("failed to parse certificate data: %w", err)
	}

	// Display the certificate data
	altdacommitment_parser.DisplayCertificateData(certV3)
	return nil
}
