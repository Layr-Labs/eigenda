package flags

import (
	"fmt"

	proxycommon "github.com/Layr-Labs/eigenda/api/proxy/common"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/urfave/cli"
)

const (
	FlagPrefix = ""
	envPrefix  = "GAS_EXHAUSTION_CERT_METER"
)

var (
	/* Required Flags*/
	NetworkFlag = cli.StringFlag{
		Name: common.PrefixFlag(FlagPrefix, "eigenda-network"),
		Usage: fmt.Sprintf(`The EigenDA network that is being used. 
See https://github.com/Layr-Labs/eigenda/blob/master/api/proxy/common/eigenda_network.go
for the exact values getting set by this flag. Permitted EigenDANetwork values include 
%s, %s, %s, & %s.`,
			proxycommon.MainnetEigenDANetwork,
			proxycommon.HoleskyTestnetEigenDANetwork,
			proxycommon.HoleskyPreprodEigenDANetwork,
			proxycommon.SepoliaTestnetEigenDANetwork,
		),
		Required: true,
		EnvVar:   common.PrefixEnvVar(envPrefix, "EIGENDA_NETWORK"),
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
	NetworkFlag,
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
