package memstore

import (
	"fmt"
	"os"
	"time"

	"github.com/Layr-Labs/eigenda-proxy/verify"
	"github.com/urfave/cli/v2"
)

var (
	EnabledFlagName    = withFlagPrefix("enabled")
	ExpirationFlagName = withFlagPrefix("expiration")
	PutLatencyFlagName = withFlagPrefix("put-latency")
	GetLatencyFlagName = withFlagPrefix("get-latency")
)

func withFlagPrefix(s string) string {
	return "memstore." + s
}

func withEnvPrefix(envPrefix, s string) string {
	return envPrefix + "_MEMSTORE_" + s
}

// if these deprecated env vars are used, we force the user to update their config
// in the flags' actions
func withDeprecatedEnvPrefix(_, s string) string {
	return "MEMSTORE_" + s
}

// CLIFlags ... used for memstore backend configuration
// category is used to group the flags in the help output (see https://cli.urfave.org/v2/examples/flags/#grouping)
func CLIFlags(envPrefix, category string) []cli.Flag {
	return []cli.Flag{
		&cli.BoolFlag{
			Name:     EnabledFlagName,
			Usage:    "Whether to use memstore for DA logic.",
			EnvVars:  []string{withEnvPrefix(envPrefix, "ENABLED"), withDeprecatedEnvPrefix(envPrefix, "ENABLED")},
			Category: category,
			Action: func(_ *cli.Context, _ bool) error {
				if _, ok := os.LookupEnv(withDeprecatedEnvPrefix(envPrefix, "ENABLED")); ok {
					return fmt.Errorf("env var %s is deprecated for flag %s, use %s instead",
						withDeprecatedEnvPrefix(envPrefix, "ENABLED"),
						EnabledFlagName,
						withEnvPrefix(envPrefix, "ENABLED"))
				}
				return nil
			},
		},
		&cli.DurationFlag{
			Name:     ExpirationFlagName,
			Usage:    "Duration that a memstore blob/commitment pair is allowed to live.",
			Value:    25 * time.Minute,
			EnvVars:  []string{withEnvPrefix(envPrefix, "EXPIRATION"), withDeprecatedEnvPrefix(envPrefix, "EXPIRATION")},
			Category: category,
			Action: func(_ *cli.Context, _ time.Duration) error {
				if _, ok := os.LookupEnv(withDeprecatedEnvPrefix(envPrefix, "EXPIRATION")); ok {
					return fmt.Errorf("env var %s is deprecated for flag %s, use %s instead",
						withDeprecatedEnvPrefix(envPrefix, "EXPIRATION"),
						ExpirationFlagName,
						withEnvPrefix(envPrefix, "EXPIRATION"))
				}
				return nil
			},
		},
		&cli.DurationFlag{
			Name:     PutLatencyFlagName,
			Usage:    "Artificial latency added for memstore backend to mimic EigenDA's dispersal latency.",
			Value:    0,
			EnvVars:  []string{withEnvPrefix(envPrefix, "PUT_LATENCY")},
			Category: category,
		},
		&cli.DurationFlag{
			Name:     GetLatencyFlagName,
			Usage:    "Artificial latency added for memstore backend to mimic EigenDA's retrieval latency.",
			Value:    0,
			EnvVars:  []string{withEnvPrefix(envPrefix, "GET_LATENCY")},
			Category: category,
		},
	}
}

func ReadConfig(ctx *cli.Context) Config {
	return Config{
		// TODO: there has to be a better way to get MaxBlobLengthBytes
		// right now we get it from the verifier cli, but there's probably a way to share flags more nicely?
		// maybe use a duplicate but hidden flag in memstore category, and set it using the action by reading
		// from the other flag?
		MaxBlobSizeBytes: verify.MaxBlobLengthBytes,
		BlobExpiration:   ctx.Duration(ExpirationFlagName),
		PutLatency:       ctx.Duration(PutLatencyFlagName),
		GetLatency:       ctx.Duration(GetLatencyFlagName),
	}
}
