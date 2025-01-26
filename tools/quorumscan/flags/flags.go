package flags

import (
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/core/thegraph"
	"github.com/urfave/cli"
)

const (
	FlagPrefix = ""
	envPrefix  = "QUORUMSCAN"
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
	BlockNumberFlag = cli.Uint64Flag{
		Name:     common.PrefixFlag(FlagPrefix, "block-number"),
		Usage:    "Block number to query state from (default: latest)",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envPrefix, "BLOCK_NUMBER"),
		Value:    0,
	}
	QuorumIDsFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "quorum-ids"),
		Usage:    "Comma-separated list of quorum IDs to scan (default: 0,1,2)",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envPrefix, "QUORUM_IDS"),
		Value:    "0,1,2",
	}
)

var requiredFlags = []cli.Flag{
	BlsOperatorStateRetrieverFlag,
	EigenDAServiceManagerFlag,
}

var optionalFlags = []cli.Flag{
	BlockNumberFlag,
	QuorumIDsFlag,
}

// Flags contains the list of configuration options available to the binary.
var Flags []cli.Flag

func init() {
	Flags = append(requiredFlags, optionalFlags...)
	Flags = append(Flags, common.LoggerCLIFlags(envPrefix, FlagPrefix)...)
	Flags = append(Flags, geth.EthClientFlags(envPrefix)...)
	Flags = append(Flags, thegraph.CLIFlags(envPrefix)...)
}
