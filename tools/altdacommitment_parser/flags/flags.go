package flags

import (
	"github.com/urfave/cli"
)

var (
	CertHexFlag = cli.StringFlag{
		Name:  "hex",
		Usage: "Hex-encoded RLP string of EigenDA certificate (with or without 0x prefix)",
		Required: true,
	}
)

var Flags = []cli.Flag{
	CertHexFlag,
}
