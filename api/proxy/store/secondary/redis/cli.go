package redis

import (
	"time"

	"github.com/urfave/cli/v2"
)

var (
	EndpointFlagName = withFlagPrefix("endpoint")
	PasswordFlagName = withFlagPrefix("password")
	DBFlagName       = withFlagPrefix("db")
	EvictionFlagName = withFlagPrefix("eviction")
)

func withFlagPrefix(s string) string {
	return "redis." + s
}

func withEnvPrefix(envPrefix, s string) []string {
	return []string{envPrefix + "_REDIS_" + s}
}

// CLIFlags ... used for Redis backend configuration
// category is used to group the flags in the help output (see https://cli.urfave.org/v2/examples/flags/#grouping)
func CLIFlags(envPrefix, category string) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     EndpointFlagName,
			Usage:    "Redis endpoint",
			EnvVars:  withEnvPrefix(envPrefix, "ENDPOINT"),
			Category: category,
		},
		&cli.StringFlag{
			Name:     PasswordFlagName,
			Usage:    "Redis password",
			EnvVars:  withEnvPrefix(envPrefix, "PASSWORD"),
			Category: category,
		},
		&cli.IntFlag{
			Name:     DBFlagName,
			Usage:    "Redis database",
			Value:    0,
			EnvVars:  withEnvPrefix(envPrefix, "DB"),
			Category: category,
		},
		&cli.DurationFlag{
			Name:     EvictionFlagName,
			Usage:    "Redis eviction time",
			Value:    24 * time.Hour,
			EnvVars:  withEnvPrefix(envPrefix, "EVICTION"),
			Category: category,
		},
	}
}

func ReadConfig(ctx *cli.Context) Config {
	return Config{
		Endpoint: ctx.String(EndpointFlagName),
		Password: ctx.String(PasswordFlagName),
		DB:       ctx.Int(DBFlagName),
		Eviction: ctx.Duration(EvictionFlagName),
	}
}
