package flags

import "github.com/urfave/cli"

var (
	CertHexFlag = cli.StringFlag{
		Name:     "hex",
		Usage:    "Hex-encoded RLP altda commitment string to parse (can include 0x prefix)",
		Required: true,
	}
)

var ParserFlags = []cli.Flag{
	CertHexFlag,
}
