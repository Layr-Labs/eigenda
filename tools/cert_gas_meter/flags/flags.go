package flags

import (
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/urfave/cli"
)

const (
	FlagPrefix = ""
	envPrefix  = "CERT_GAS_METER"
)

var (
	/* Required Flags*/
	/*
		EigenDADirectoryFlag = cli.StringFlag{
			Name:     common.PrefixFlag(FlagPrefix, "eigenda-directory"),
			Usage:    "Address of the EigenDA directory contract, which points to all other EigenDA contract addresses. This is the only contract entrypoint needed offchain.",
			Required: false,
			EnvVar:   common.PrefixEnvVar(envPrefix, "EIGENDA_DIRECTORY"),
		}
	*/
	EigenDAServiceManagerFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "eigenda-service-manager"),
		Usage:    "[Deprecated: use EigenDADirectory instead] Address of the EigenDA Service Manager",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envPrefix, "EIGENDA_SERVICE_MANAGER"),
	}
	OperatorStateRetrieverFlag = cli.StringFlag{
		Name: common.PrefixFlag(FlagPrefix, "bls-operator-state-retriever"),
		Usage: "[Deprecated: use EigenDADirectory instead] Address of the OperatorStateRetriever contract. " +
			"Note that the contract no longer uses the BLS prefix.",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envPrefix, "BLS_OPERATOR_STATE_RETRIVER"),
	}
	CertFileFlag = cli.StringFlag{
		Name:     "cert-rlp-file",
		Usage:    "Path to the RLP-encoded EigenDA certificate file",
		Required: true,
	}
)

var requiredFlags = []cli.Flag{}

var optionalFlags = []cli.Flag{
	CertFileFlag,
	OperatorStateRetrieverFlag,
	EigenDAServiceManagerFlag,
}

// Flags contains the list of configuration options available to the binary.
var Flags []cli.Flag

func init() {
	Flags = append(requiredFlags, optionalFlags...)
	Flags = append(Flags, common.LoggerCLIFlags(envPrefix, FlagPrefix)...)
	Flags = append(Flags, geth.EthClientFlags(envPrefix)...)
	//Flags = append(Flags, thegraph.CLIFlags(envPrefix)...)
}
