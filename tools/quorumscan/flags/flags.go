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
	EigenDADirectoryFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "eigenda-directory"),
		Usage:    "Address of the EigenDA directory contract, which points to all other EigenDA contract addresses. This is the only contract entrypoint needed offchain.",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envPrefix, "EIGENDA_DIRECTORY"),
	}
	BlsOperatorStateRetrieverFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "bls-operator-state-retriever"),
		Usage:    "[Deprecated: use EigenDADirectory instead] Address of the BLS operator state Retriever",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envPrefix, "BLS_OPERATOR_STATE_RETRIVER"),
	}
	EigenDAServiceManagerFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "eigenda-service-manager"),
		Usage:    "[Deprecated: use EigenDADirectory instead] Address of the EigenDA Service Manager",
		Required: false,
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
		Usage:    "Comma-separated list of quorum IDs to scan (default: all)",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envPrefix, "QUORUM_IDS"),
		Value:    "",
	}
	TopNFlag = cli.UintFlag{
		Name:     common.PrefixFlag(FlagPrefix, "top"),
		Usage:    "Show only top N operators by stake",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envPrefix, "TOP"),
		Value:    0,
	}
	OutputFormatFlag = cli.StringFlag{
		Name:     "output-format",
		Usage:    "Output format (table/csv)",
		Value:    "table",
		Required: false,
	}
	OutputFileFlag = cli.StringFlag{
		Name:     "output-file",
		Usage:    "Write output to a file instead of stdout",
		Required: false,
	}
)

var requiredFlags = []cli.Flag{}

var optionalFlags = []cli.Flag{
	BlockNumberFlag,
	QuorumIDsFlag,
	TopNFlag,
	OutputFormatFlag,
	OutputFileFlag,
	EigenDADirectoryFlag,
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
