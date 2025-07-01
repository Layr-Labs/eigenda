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
	/* Optional Flags*/
	AddressDirectoryFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "address-directory"),
		Usage:    "Address of the EigenDA Directory contract (preferred over individual contract addresses)",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envPrefix, "ADDRESS_DIRECTORY"),
	}
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
	AddressDirectoryFlag,

}

// Flags contains the list of configuration options available to the binary.
var Flags []cli.Flag

func init() {
	Flags = append(requiredFlags, optionalFlags...)
	Flags = append(Flags, common.LoggerCLIFlags(envPrefix, FlagPrefix)...)
	Flags = append(Flags, geth.EthClientFlags(envPrefix)...)
	Flags = append(Flags, thegraph.CLIFlags(envPrefix)...)
}
