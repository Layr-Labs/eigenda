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
)

var requiredFlags = []cli.Flag{
	SubgraphEndpointFlag,
}

var optionalFlags = []cli.Flag{
	OperatorIdFlag,
	DaysFlag,
}

// Flags contains the list of configuration options available to the binary.
var Flags []cli.Flag

func init() {
	Flags = append(requiredFlags, optionalFlags...)
	Flags = append(Flags, common.LoggerCLIFlags(envPrefix, FlagPrefix)...)
}
