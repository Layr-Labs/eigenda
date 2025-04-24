package thegraph

import (
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/urfave/cli"
)

const (
	EndpointFlagName       = "thegraph.endpoint"
	BackoffFlagName        = "thegraph.backoff"
	MaxRetriesFlagName     = "thegraph.max_retries"
	OperatorStateCacheSize = "thegraph.operator_state_cache_size"
)

type Config struct {
	Endpoint               string        // The Graph endpoint
	PullInterval           time.Duration // The interval to pull data from The Graph
	MaxRetries             int           // The maximum number of retries to pull data from The Graph
	OperatorStateCacheSize int           // The size of the cache
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
			Name:   BackoffFlagName,
			Usage:  "Backoff for retries",
			Value:  100 * time.Millisecond,
			EnvVar: common.PrefixEnvVar(envPrefix, "GRAPH_BACKOFF"),
		},
		cli.UintFlag{
			Name:   MaxRetriesFlagName,
			Usage:  "The maximum number of retries",
			Value:  5,
			EnvVar: common.PrefixEnvVar(envPrefix, "GRAPH_MAX_RETRIES"),
		},
		cli.IntFlag{
			Name:   OperatorStateCacheSize,
			Usage:  "The size of the operator state cache in elements (default 32)",
			Value:  32,
			EnvVar: common.PrefixEnvVar(envPrefix, "GRAPH_OPERATOR_STATE_CACHE_SIZE"),
		},
	}
}

func ReadCLIConfig(ctx *cli.Context) Config {

	return Config{
		Endpoint:               ctx.String(EndpointFlagName),
		PullInterval:           ctx.Duration(BackoffFlagName),
		MaxRetries:             ctx.Int(MaxRetriesFlagName),
		OperatorStateCacheSize: ctx.Int(OperatorStateCacheSize),
	}

}
