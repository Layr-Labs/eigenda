package flags

import (
	"github.com/urfave/cli"
)

const (
	FlagPrefix = ""
	envPrefix  = "ALTDACOMMITMENT_PARSER"
)

var (
	CertHexFlag = cli.StringFlag{
		Name:     "cert-hex",
		Usage:    "Hex-encoded RLP string of EigenDA certificate (with or without 0x prefix)",
		Required: true,
		EnvVar:   "ALTDACOMMITMENT_PARSER_CERT_HEX",
	}
)

var Flags = []cli.Flag{
	CertHexFlag,
}
