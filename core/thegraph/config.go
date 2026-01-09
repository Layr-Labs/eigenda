package thegraph

import (
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/urfave/cli"
)

const (
	EndpointFlagName   = "thegraph.endpoint"
	BackoffFlagName    = "thegraph.backoff"
	MaxRetriesFlagName = "thegraph.max_retries"
)

type Config struct {
	// The Graph endpoint
	Endpoint string `docs:"required"`
	// The interval to pull data from The Graph
	PullInterval time.Duration
	// The maximum number of retries to pull data from The Graph
	MaxRetries int
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
	}
}

func ReadCLIConfig(ctx *cli.Context) Config {

	return Config{
		Endpoint:     ctx.String(EndpointFlagName),
		PullInterval: ctx.Duration(BackoffFlagName),
		MaxRetries:   ctx.Int(MaxRetriesFlagName),
	}
}

func DefaultTheGraphConfig() Config {
	return Config{
		PullInterval: 100 * time.Millisecond,
		MaxRetries:   5,
	}
}

func (c *Config) Verify() error {
	if c.Endpoint == "" {
		return fmt.Errorf("thegraph endpoint is required")
	}
	if c.PullInterval <= 0 {
		return fmt.Errorf("pull interval must be positive, got %v", c.PullInterval)
	}
	if c.MaxRetries < 0 {
		return fmt.Errorf("max retries cannot be negative, got %d", c.MaxRetries)
	}
	return nil
}
