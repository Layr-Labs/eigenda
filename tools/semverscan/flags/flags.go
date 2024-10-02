package flags

import (
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/core/thegraph"
	"github.com/urfave/cli"
)

const (
	FlagPrefix = ""
	envPrefix  = "SEMVERSCAN"
)

var (
	/* Required Flags*/
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
	/* Optional Flags*/
	TimeoutFlag = cli.DurationFlag{
		Name:     common.PrefixFlag(FlagPrefix, "timeout"),
		Usage:    "seconds to wait for GPRC response",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envPrefix, "TIMEOUT"),
		Value:    3 * time.Second,
	}
	WorkersFlag = cli.UintFlag{
		Name:     common.PrefixFlag(FlagPrefix, "workers"),
		Usage:    "maximum number of concurrent node info requests",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envPrefix, "WORKERS"),
		Value:    10,
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
	BlsOperatorStateRetrieverFlag,
	EigenDAServiceManagerFlag,
}

var optionalFlags = []cli.Flag{
	TimeoutFlag,
	WorkersFlag,
	OperatorIdFlag,
}

// Flags contains the list of configuration options available to the binary.
var Flags []cli.Flag

func init() {
	Flags = append(requiredFlags, optionalFlags...)
	Flags = append(Flags, common.LoggerCLIFlags(envPrefix, FlagPrefix)...)
	Flags = append(Flags, geth.EthClientFlags(envPrefix)...)
	Flags = append(Flags, thegraph.CLIFlags(envPrefix)...)
}
