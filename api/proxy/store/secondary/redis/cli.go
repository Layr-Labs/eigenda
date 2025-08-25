package redis

import (
	"fmt"
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

// DeprecatedCLIFlags ... used for Redis backend configuration
// category is used to group the flags in the help output (see https://cli.urfave.org/v2/examples/flags/#grouping)
func DeprecatedCLIFlags(envPrefix, category string) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     EndpointFlagName,
			Usage:    "Redis endpoint",
			EnvVars:  withEnvPrefix(envPrefix, "ENDPOINT"),
			Category: category,
			Hidden:   true,
			Action: func(ctx *cli.Context, s string) error {
				return fmt.Errorf("redis secondary store is no longer supported: flag --%s is deprecated", EndpointFlagName)
			},
		},
		&cli.StringFlag{
			Name:     PasswordFlagName,
			Usage:    "Redis password",
			EnvVars:  withEnvPrefix(envPrefix, "PASSWORD"),
			Category: category,
			Hidden:   true,
			Action: func(ctx *cli.Context, s string) error {
				return fmt.Errorf("redis secondary store is no longer supported: flag --%s is deprecated", PasswordFlagName)
			},
		},
		&cli.IntFlag{
			Name:     DBFlagName,
			Usage:    "Redis database",
			Value:    0,
			EnvVars:  withEnvPrefix(envPrefix, "DB"),
			Category: category,
			Hidden:   true,
			Action: func(ctx *cli.Context, _ int) error {
				return fmt.Errorf("redis secondary store is no longer supported: flag --%s is deprecated", DBFlagName)
			},
		},
		&cli.DurationFlag{
			Name:     EvictionFlagName,
			Usage:    "Redis eviction time",
			Value:    24 * time.Hour,
			EnvVars:  withEnvPrefix(envPrefix, "EVICTION"),
			Category: category,
			Hidden:   true,
			Action: func(ctx *cli.Context, _ time.Duration) error {
				return fmt.Errorf("redis secondary store is no longer supported: flag --%s is deprecated", EvictionFlagName)
			},
		},
	}
}
