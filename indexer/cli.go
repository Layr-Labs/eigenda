package indexer

import (
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/urfave/cli"
)

const (
	PullIntervalFlagName = "indexer-pull-interval"
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
	}
}

func ReadIndexerConfig(ctx *cli.Context) Config {
	return Config{
		PullInterval: ctx.GlobalDuration(PullIntervalFlagName),
	}
}
