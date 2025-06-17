package memstore

import (
	"fmt"
	"os"
	"time"

	"github.com/Layr-Labs/eigenda-proxy/store/generated_key/eigenda/verify"
	"github.com/Layr-Labs/eigenda-proxy/store/generated_key/memstore/memconfig"
	"github.com/urfave/cli/v2"
)

var (
	EnabledFlagName                 = withFlagPrefix("enabled")
	ExpirationFlagName              = withFlagPrefix("expiration")
	PutLatencyFlagName              = withFlagPrefix("put-latency")
	GetLatencyFlagName              = withFlagPrefix("get-latency")
	PutReturnsFailoverErrorFlagName = withFlagPrefix("put-returns-failover-error")
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
			Action: func(ctx *cli.Context, enabled bool) error {
				if _, ok := os.LookupEnv(withDeprecatedEnvPrefix(envPrefix, "ENABLED")); ok {
					return fmt.Errorf("env var %s is deprecated for flag %s, use %s instead",
						withDeprecatedEnvPrefix(envPrefix, "ENABLED"),
						EnabledFlagName,
						withEnvPrefix(envPrefix, "ENABLED"))
				}
				if enabled {
					// If memstore is enabled, we disable cert verification,
					// because memstore generates some meaningless certs.
					err := ctx.Set(verify.CertVerificationDisabledFlagName, "true")
					if err != nil {
						return fmt.Errorf("failed to set %s: %w", verify.CertVerificationDisabledFlagName, err)
					}
				}
				return nil
			},
		},
		&cli.DurationFlag{
			Name:  ExpirationFlagName,
			Usage: "Duration that a memstore blob/commitment pair is allowed to live. Setting to (0) results in no expiration.",
			Value: 25 * time.Minute,
			EnvVars: []string{
				withEnvPrefix(envPrefix, "EXPIRATION"),
				withDeprecatedEnvPrefix(envPrefix, "EXPIRATION"),
			},
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
		&cli.BoolFlag{
			Name: PutReturnsFailoverErrorFlagName,
			Usage: fmt.Sprintf(
				"When true, Put requests will return a failover error, after sleeping for --%s duration.",
				PutLatencyFlagName,
			),
			Value:    false,
			EnvVars:  []string{withEnvPrefix(envPrefix, "PUT_RETURNS_FAILOVER_ERROR")},
			Category: category,
		},
	}
}

func ReadConfig(ctx *cli.Context, maxBlobSizeBytes uint64) (*memconfig.SafeConfig, error) {
	return memconfig.NewSafeConfig(
		memconfig.Config{
			MaxBlobSizeBytes:        maxBlobSizeBytes,
			BlobExpiration:          ctx.Duration(ExpirationFlagName),
			PutLatency:              ctx.Duration(PutLatencyFlagName),
			GetLatency:              ctx.Duration(GetLatencyFlagName),
			PutReturnsFailoverError: ctx.Bool(PutReturnsFailoverErrorFlagName),
		}), nil
}
