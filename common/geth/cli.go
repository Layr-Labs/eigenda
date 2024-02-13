/*
 * Copyright 2024 Layr Labs, Inc.
 *
 * Use of this source code is governed by an MIT-style
 * license that can be found in the LICENSE file or at
 * https://opensource.org/licenses/MIT.
 */
package geth

import (
	"github.com/Layr-Labs/eigenda/common"
	"github.com/urfave/cli"
)

var (
	rpcUrlFlagName           = "chain.rpc"
	privateKeyFlagName       = "chain.private-key"
	numConfirmationsFlagName = "chain.num-confirmations"
)

type EthClientConfig struct {
	RPCURL           string
	PrivateKeyString string
	NumConfirmations int
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
	}
}

func ReadEthClientConfig(ctx *cli.Context) EthClientConfig {
	cfg := EthClientConfig{}
	cfg.RPCURL = ctx.GlobalString(rpcUrlFlagName)
	cfg.PrivateKeyString = ctx.GlobalString(privateKeyFlagName)
	cfg.NumConfirmations = ctx.GlobalInt(numConfirmationsFlagName)
	return cfg
}

// ReadEthClientConfigRPCOnly doesn't read private key from flag.
// The private key for Node should be read from encrypted key file.
func ReadEthClientConfigRPCOnly(ctx *cli.Context) EthClientConfig {
	cfg := EthClientConfig{}
	cfg.RPCURL = ctx.GlobalString(rpcUrlFlagName)
	cfg.NumConfirmations = ctx.GlobalInt(numConfirmationsFlagName)
	return cfg
}
