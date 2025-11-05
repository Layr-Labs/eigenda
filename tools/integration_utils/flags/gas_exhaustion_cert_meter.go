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
%s, %s, & %s.`,
			proxycommon.MainnetEigenDANetwork,
			proxycommon.SepoliaTestnetEigenDANetwork,
			proxycommon.HoodiTestnetEigenDANetwork,
		),
		Required: true,
		EnvVar:   common.PrefixEnvVar(envPrefix, "EIGENDA_NETWORK"),
	}
	EthRpcUrlFlag = &cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "eth-rpc-url"),
		Usage:    "Ethereum RPC URL",
		EnvVar:   common.PrefixEnvVar(envPrefix, "ETH_RPC_URL"),
		Required: true,
	}

	CertVerifierAddrFlag = &cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "cert-verifier-addr"),
		Usage:    "immutable cert verifier address",
		EnvVar:   common.PrefixEnvVar(envPrefix, "CERT_VERIFIER_ADDR"),
		Required: true,
	}
)

var requiredFlags = []cli.Flag{
	NetworkFlag,
	EthRpcUrlFlag,
	CertHexFlag,
	CertVerifierAddrFlag,
}

var optionalFlags = []cli.Flag{}

var GasExhaustionCertMeterFlags []cli.Flag

func init() {
	GasExhaustionCertMeterFlags = append(requiredFlags, optionalFlags...)
}
