package flags

import (
	"github.com/Layr-Labs/eigenda/common"
	"github.com/urfave/cli"
)

const (
	FlagPrefix = ""
	envPrefix  = "OPSCAN"
)

var (
	/* Required Flags*/
	SubgraphEndpointFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "subgraph-endpoint"),
		Usage:    "Subgraph endpoint to query operator state",
		Required: true,
		EnvVar:   common.PrefixEnvVar(envPrefix, "SUBGRAPH_ENDPOINT"),
	}
	/* Optional Flags*/
	TimeoutFlag = cli.DurationFlag{
		Name:     common.PrefixFlag(FlagPrefix, "timeout"),
		Usage:    "seconds to wait for GPRC response",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envPrefix, "TIMEOUT"),
		Value:    3,
	}
	MaxConnectionsFlag = cli.IntFlag{
		Name:     common.PrefixFlag(FlagPrefix, "max-connections"),
		Usage:    "maximum number of connections to DA nodes (defaults to 20)",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envPrefix, "MAX_CONNECTIONS"),
		Value:    30,
	}
	OperatorIdFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "operator-id"),
		Usage:    "operator id to scan",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envPrefix, "OPERATOR_ID"),
		Value:    "",
	}
)

var requiredFlags = []cli.Flag{
	SubgraphEndpointFlag,
}

var optionalFlags = []cli.Flag{
	TimeoutFlag,
	MaxConnectionsFlag,
	OperatorIdFlag,
}

// Flags contains the list of configuration options available to the binary.
var Flags []cli.Flag

func init() {
	Flags = append(requiredFlags, optionalFlags...)
	Flags = append(Flags, common.LoggerCLIFlags(envPrefix, FlagPrefix)...)
}
