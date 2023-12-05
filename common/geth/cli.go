package geth

import (
	"github.com/Layr-Labs/eigenda/common"
	"github.com/urfave/cli"
)

var (
	rpcUrlFlagName     = "chain.rpc"
	PrivateKeyFlagName = "chain.private-key"
)

type EthClientConfig struct {
	RPCURL           string
	PrivateKeyString string
}

func EthClientFlags(envPrefix string) []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:     rpcUrlFlagName,
			Usage:    "Chain rpc",
			Required: true,
			EnvVar:   common.PrefixEnvVar(envPrefix, "CHAIN_RPC"),
		},
		cli.StringFlag{
			Name:     PrivateKeyFlagName,
			Usage:    "Ethereum private key for disperser",
			Required: false,
			EnvVar:   common.PrefixEnvVar(envPrefix, "PRIVATE_KEY"),
		},
	}
}

func ReadEthClientConfig(ctx *cli.Context) EthClientConfig {
	cfg := EthClientConfig{}
	cfg.RPCURL = ctx.GlobalString(rpcUrlFlagName)
	pkStr := ctx.GlobalString(PrivateKeyFlagName)
	if len(pkStr) >= 2 && pkStr[:2] == "0x" {
		pkStr = pkStr[2:]
	}
	cfg.PrivateKeyString = pkStr
	return cfg
}

// ReadEthClientConfigRPCOnly doesn't read private key from flag.
// The private key for Node should be read from encrypted key file.
func ReadEthClientConfigRPCOnly(ctx *cli.Context) EthClientConfig {
	cfg := EthClientConfig{}
	cfg.RPCURL = ctx.GlobalString(rpcUrlFlagName)
	return cfg
}
