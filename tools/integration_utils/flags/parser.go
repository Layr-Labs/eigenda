package flags

import "github.com/urfave/cli"

var (
	CertHexFlag = cli.StringFlag{
		Name:     "cert-hex",
		Usage:    "Hex-encoded RLP certificate string to parse (can include 0x prefix)",
		Required: true,
	}
)

var ParserFlags = []cli.Flag{
	CertHexFlag,
}
