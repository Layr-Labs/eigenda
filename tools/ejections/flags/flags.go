package flags

import (
	"github.com/Layr-Labs/eigenda/common"
	"github.com/urfave/cli"
)

const (
	FlagPrefix = ""
	envPrefix  = ""
)

var (
	/* Required Flags*/
	SubgraphEndpointFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "subgraph"),
		Usage:    "Subgraph URL to query operator state",
		Required: true,
		EnvVar:   common.PrefixEnvVar(envPrefix, "SUBGRAPH"),
	}
	OperatorIdFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "operator_id"),
		Usage:    "Query operator id",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envPrefix, "OPERATOR_ID"),
		Value:    "",
	}
	DaysFlag = cli.UintFlag{
		Name:     common.PrefixFlag(FlagPrefix, "days"),
		Usage:    "Lookback days",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envPrefix, "DAYS"),
		Value:    1,
	}
	FirstFlag = cli.UintFlag{
		Name:     common.PrefixFlag(FlagPrefix, "first"),
		Usage:    "Return first n records (default 1000, max 10000)",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envPrefix, "FIRST"),
		Value:    1000,
	}
	SkipFlag = cli.UintFlag{
		Name:     common.PrefixFlag(FlagPrefix, "skip"),
		Usage:    "Skip first n records (default 0, max 1000000)",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envPrefix, "SKIP"),
		Value:    0,
	}
)

var requiredFlags = []cli.Flag{
	SubgraphEndpointFlag,
}

var optionalFlags = []cli.Flag{
	OperatorIdFlag,
	DaysFlag,
	FirstFlag,
	SkipFlag,
}

// Flags contains the list of configuration options available to the binary.
var Flags []cli.Flag

func init() {
	Flags = append(requiredFlags, optionalFlags...)
	Flags = append(Flags, common.LoggerCLIFlags(envPrefix, FlagPrefix)...)
}
