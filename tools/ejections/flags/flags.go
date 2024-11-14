package flags

import (
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/geth"
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
	BlsOperatorStateRetrieverFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "bls-operator-state-retriever"),
		Usage:    "Address of the BLS Operator State Retriever",
		Required: true,
		EnvVar:   common.PrefixEnvVar(envPrefix, "BLS_OPERATOR_STATE_RETRIVER"),
	}
	EigenDAServiceManagerFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "eigenda-service-manager"),
		Usage:    "Address of the EigenDA Service Manager",
		Required: true,
		EnvVar:   common.PrefixEnvVar(envPrefix, "EIGENDA_SERVICE_MANAGER"),
	}
)

var requiredFlags = []cli.Flag{
	SubgraphEndpointFlag,
	BlsOperatorStateRetrieverFlag,
	EigenDAServiceManagerFlag,
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
	Flags = append(Flags, geth.EthClientFlags(envPrefix)...)
}
