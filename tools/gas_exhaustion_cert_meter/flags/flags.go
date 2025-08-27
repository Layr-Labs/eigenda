package flags

import (
	"github.com/Layr-Labs/eigenda/common"
	"github.com/urfave/cli"
)

const (
	FlagPrefix = ""
	envPrefix  = "GAS_EXHAUSTION_CERT_METER"
)

var (
	/* Required Flags*/
	EigenDADirectoryFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "eigenda-directory"),
		Usage:    "Address of the EigenDA directory contract, which points to all other EigenDA contract addresses.",
		Required: true,
		EnvVar:   common.PrefixEnvVar(envPrefix, "EIGENDA_DIRECTORY"),
	}
	CertRlpFileFlag = cli.StringFlag{
		Name: common.PrefixFlag(FlagPrefix, "cert-rlp-file"),
		Usage: "Path to the RLP-encoded EigenDA V3 certificate file. " +
			"Examples: ./data/cert_v3.mainnet.rlp, ./data/cert_v3.sepolia.rlp",
		Required: true,
		EnvVar:   common.PrefixEnvVar(envPrefix, "CERT_RLP_FILE"),
	}
	EthRpcUrlFlag = &cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "eth-rpc-url"),
		Usage:    "Ethereum RPC URL",
		EnvVar:   common.PrefixEnvVar(envPrefix, "ETH_RPC_URL"),
		Required: true,
	}
)

var requiredFlags = []cli.Flag{
	EigenDADirectoryFlag,
	CertRlpFileFlag,
	EthRpcUrlFlag,
}

var optionalFlags = []cli.Flag{}

// Flags contains the list of configuration options available to the binary.
var Flags []cli.Flag

func init() {
	Flags = append(requiredFlags, optionalFlags...)
	Flags = append(Flags, common.LoggerCLIFlags(envPrefix, FlagPrefix)...)
}
