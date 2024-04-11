package thegraph

import (
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/urfave/cli"
)

const (
	EndpointFlagName     = "thegraph.endpoint"
	PullIntervalFlagName = "thegraph.pull_interval"
	MaxRetriesFlagName   = "thegraph.max_retries"
)

type Config struct {
	Endpoint     string        // The Graph endpoint
	PullInterval time.Duration // The interval to pull data from The Graph
	MaxRetries   int           // The maximum number of retries to pull data from The Graph
}

func CLIFlags(envPrefix string) []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:     EndpointFlagName,
			Usage:    "The Graph endpoint",
			Required: true,
			EnvVar:   common.PrefixEnvVar(envPrefix, "GRAPH_URL"),
		},
		cli.DurationFlag{
			Name:   PullIntervalFlagName,
			Usage:  "Pull interval for retries",
			Value:  100 * time.Millisecond,
			EnvVar: common.PrefixEnvVar(envPrefix, "GRAPH_PULL_INTERVAL"),
		},
		cli.UintFlag{
			Name:   MaxRetriesFlagName,
			Usage:  "The maximum number of retries",
			Value:  5,
			EnvVar: common.PrefixEnvVar(envPrefix, "GRAPH_MAX_RETRIES"),
		},
	}
}

func ReadCLIConfig(ctx *cli.Context) Config {

	return Config{
		Endpoint:     ctx.String(EndpointFlagName),
		PullInterval: ctx.Duration(PullIntervalFlagName),
		MaxRetries:   ctx.Int(MaxRetriesFlagName),
	}

}
