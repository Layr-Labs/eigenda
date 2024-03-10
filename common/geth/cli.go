package geth

import (
	"github.com/Layr-Labs/eigenda/common"
	"github.com/urfave/cli"
)

var (
	rpcUrlFlagName           = "chain.rpc"
	privateKeyFlagName       = "chain.private-key"
	numConfirmationsFlagName = "chain.num-confirmations"
	rpcUrlBackupFlagName     = "chain.rpc-backup"
	privateKeyBackupFlagName = "chain.private-key-backup"
)

type EthClientConfig struct {
	RPCURL           string
	PrivateKeyString string
	NumConfirmations int

	RPCURLBackup           []string
	PrivateKeyStringBackup []string
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
			Name:     privateKeyFlagName,
			Usage:    "Ethereum private key for disperser",
			Required: true,
			EnvVar:   common.PrefixEnvVar(envPrefix, "PRIVATE_KEY"),
		},
		cli.IntFlag{
			Name:     numConfirmationsFlagName,
			Usage:    "Number of confirmations to wait for",
			Required: false,
			Value:    0,
			EnvVar:   common.PrefixEnvVar(envPrefix, "NUM_CONFIRMATIONS"),
		},
		cli.StringSliceFlag{
			Name:     rpcUrlBackupFlagName,
			Usage:    "A list of backup for Chain rpc",
			Required: false,
			EnvVar:   common.PrefixEnvVar(envPrefix, "CHAIN_RPC_BACKUP"),
		},
		cli.StringSliceFlag{
			Name:     privateKeyBackupFlagName,
			Usage:    "A list of backup for Ethereum private key for disperser",
			Required: false,
			EnvVar:   common.PrefixEnvVar(envPrefix, "PRIVATE_KEY_BACKUP"),
		},
	}
}

func ReadEthClientConfig(ctx *cli.Context) EthClientConfig {
	cfg := EthClientConfig{}
	cfg.RPCURL = ctx.GlobalString(rpcUrlFlagName)
	cfg.PrivateKeyString = ctx.GlobalString(privateKeyFlagName)
	cfg.NumConfirmations = ctx.GlobalInt(numConfirmationsFlagName)
	cfg.RPCURLBackup = ctx.GlobalStringSlice(rpcUrlBackupFlagName)
	cfg.PrivateKeyStringBackup = ctx.GlobalStringSlice(privateKeyBackupFlagName)

	return cfg
}

// ReadEthClientConfigRPCOnly doesn't read private key from flag.
// The private key for Node should be read from encrypted key file.
func ReadEthClientConfigRPCOnly(ctx *cli.Context) EthClientConfig {
	cfg := EthClientConfig{}
	cfg.RPCURL = ctx.GlobalString(rpcUrlFlagName)
	cfg.NumConfirmations = ctx.GlobalInt(numConfirmationsFlagName)
	cfg.RPCURLBackup = ctx.GlobalStringSlice(rpcUrlBackupFlagName)
	return cfg
}
