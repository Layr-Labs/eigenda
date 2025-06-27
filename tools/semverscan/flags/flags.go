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
	/* Optional Flags*/
	AddressDirectoryFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "address-directory"),
		Usage:    "Address of the EigenDA Directory contract (preferred over individual contract addresses)",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envPrefix, "ADDRESS_DIRECTORY"),
	}
	BlsOperatorStateRetrieverFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "bls-operator-state-retriever"),
		Usage:    "Address of the BLS Operator State Retriever",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envPrefix, "BLS_OPERATOR_STATE_RETRIVER"),
	}
	EigenDAServiceManagerFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "eigenda-service-manager"),
		Usage:    "Address of the EigenDA Service Manager",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envPrefix, "EIGENDA_SERVICE_MANAGER"),
	}
	TimeoutFlag = cli.DurationFlag{
		Name:     common.PrefixFlag(FlagPrefix, "timeout"),
		Usage:    "Seconds to wait for GPRC response",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envPrefix, "TIMEOUT"),
		Value:    3 * time.Second,
	}
	WorkersFlag = cli.UintFlag{
		Name:     common.PrefixFlag(FlagPrefix, "workers"),
		Usage:    "Maximum number of concurrent node info requests",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envPrefix, "WORKERS"),
		Value:    10,
	}
	OperatorIdFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "operator-id"),
		Usage:    "Operator ID to scan",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envPrefix, "OPERATOR_ID"),
		Value:    "",
	}
	UseRetrievalClientFlag = cli.BoolFlag{
		Name:     common.PrefixFlag(FlagPrefix, "use-retrieval-client"),
		Usage:    "Use retrieval client to get operator info (default: false)",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envPrefix, "USE_RETRIEVAL_CLIENT"),
	}
)

var requiredFlags = []cli.Flag{}

var optionalFlags = []cli.Flag{
	TimeoutFlag,
	WorkersFlag,
	OperatorIdFlag,
	UseRetrievalClientFlag,
	AddressDirectoryFlag,
	BlsOperatorStateRetrieverFlag,
	EigenDAServiceManagerFlag,
}

// Flags contains the list of configuration options available to the binary.
var Flags []cli.Flag

func init() {
	Flags = append(requiredFlags, optionalFlags...)
	Flags = append(Flags, common.LoggerCLIFlags(envPrefix, FlagPrefix)...)
	Flags = append(Flags, geth.EthClientFlags(envPrefix)...)
	Flags = append(Flags, thegraph.CLIFlags(envPrefix)...)
}
