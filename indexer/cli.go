package indexer

import (
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/urfave/cli"
)

const (
	PullIntervalFlagName            = "indexer-pull-interval"
	ContractDeploymentBlockFlagName = "indexer-contract-deployment-block"
)

func CLIFlags(envPrefix string) []cli.Flag {
	return []cli.Flag{
		cli.DurationFlag{
			Name:     PullIntervalFlagName,
			Usage:    "Interval at which to pull and index new blocks and events from chain",
			Required: false,
			EnvVar:   common.PrefixEnvVar(envPrefix, "INDEXER_PULL_INTERVAL"),
			Value:    1 * time.Second,
		},
		cli.Uint64Flag{
			Name:     ContractDeploymentBlockFlagName,
			Usage:    "Block number at which the contract was deployed (used as starting point when no headers exist)",
			Required: false,
			EnvVar:   common.PrefixEnvVar(envPrefix, "INDEXER_CONTRACT_DEPLOYMENT_BLOCK"),
			Value:    0,
		},
	}
}

func ReadIndexerConfig(ctx *cli.Context) Config {
	return Config{
		PullInterval:            ctx.GlobalDuration(PullIntervalFlagName),
		ContractDeploymentBlock: ctx.GlobalUint64(ContractDeploymentBlockFlagName),
	}
}
