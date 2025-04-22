package flags

import (
	"github.com/Layr-Labs/eigenda/common"
	"github.com/urfave/cli"
)

const (
	FlagPrefix = ""
	envPrefix  = "RELAYLOAD"
)

var (
	/* Required Flags*/
	RelayUrlFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "relay-url"),
		Usage:    "Relay to query",
		Required: true,
		EnvVar:   common.PrefixEnvVar(envPrefix, "RELAY_URL"),
	}
	OperatorIdFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "operator-id"),
		Usage:    "Operator ID to query",
		Required: true,
		EnvVar:   common.PrefixEnvVar(envPrefix, "OPERATOR_ID"),
	}
	OperatorPKeyFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "operator-pkey"),
		Usage:    "Operator private key to query",
		Required: true,
		EnvVar:   common.PrefixEnvVar(envPrefix, "OPERATOR_PKEY"),
	}
	DataApiUrlFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "data-api-url"),
		Usage:    "Data API URL to query",
		Required: true,
		EnvVar:   common.PrefixEnvVar(envPrefix, "DATA_API_URL"),
	}
	NumThreadsFlag = cli.IntFlag{
		Name:     common.PrefixFlag(FlagPrefix, "num-threads"),
		Usage:    "Number of threads to run",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envPrefix, "NUM_THREADS"),
		Value:    1,
	}
	RangeSizesFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "range-sizes"),
		Usage:    "Range sizes to select from",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envPrefix, "RANGE_SIZES"),
		Value:    "10,25,100",
	}
	RequestSizesFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "request-sizes"),
		Usage:    "Request sizes to select from",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envPrefix, "REQUEST_SIZES"),
		Value:    "15,25,35,50,64,100",
	}
)

var requiredFlags = []cli.Flag{
	RelayUrlFlag,
	OperatorIdFlag,
	OperatorPKeyFlag,
	DataApiUrlFlag,
	NumThreadsFlag,
	RangeSizesFlag,
	RequestSizesFlag,
}

var optionalFlags = []cli.Flag{}

// Flags contains the list of configuration options available to the binary.
var Flags []cli.Flag

func init() {
	Flags = append(requiredFlags, optionalFlags...)
	Flags = append(Flags, common.LoggerCLIFlags(envPrefix, FlagPrefix)...)
}
