package store

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

// All of these flags are deprecated and will be removed in release v2.0.0
// we leave them here with actions that crash the program to ensure they are not used,
// and to make it easier for users to find the new flags (instead of silently crashing late during execution
// because some flag's env var was changed but the user forgot to update it)
var (
	DeprecatedFallbackTargetsFlagName = withDeprecatedFlagPrefix("fallback-targets")
	DeprecatedCacheTargetsFlagName    = withDeprecatedFlagPrefix("cache-targets")
	DeprecatedConcurrentWriteThreads  = withDeprecatedFlagPrefix("concurrent-write-routines")
)

func withDeprecatedFlagPrefix(s string) string {
	return "routing." + s
}

func withDeprecatedEnvPrefix(envPrefix, s string) []string {
	return []string{envPrefix + "_" + s}
}

// CLIFlags ... used for EigenDA client configuration
func DeprecatedCLIFlags(envPrefix, category string) []cli.Flag {
	return []cli.Flag{
		&cli.StringSliceFlag{
			Name:     DeprecatedFallbackTargetsFlagName,
			Usage:    "List of read fallback targets to rollover to if cert can't be read from EigenDA.",
			Value:    cli.NewStringSlice(),
			EnvVars:  withDeprecatedEnvPrefix(envPrefix, "FALLBACK_TARGETS"),
			Category: category,
			Action: func(*cli.Context, []string) error {
				return fmt.Errorf("flag --%s (env var %s) is deprecated, use --%s (env var %s) instead",
					DeprecatedFallbackTargetsFlagName, withDeprecatedEnvPrefix(envPrefix, "FALLBACK_TARGETS"),
					FallbackTargetsFlagName, withEnvPrefix(envPrefix, "FALLBACK_TARGETS"))
			},
		},
		&cli.StringSliceFlag{
			Name:     DeprecatedCacheTargetsFlagName,
			Usage:    "List of caching targets to use fast reads from EigenDA.",
			Value:    cli.NewStringSlice(),
			EnvVars:  withDeprecatedEnvPrefix(envPrefix, "CACHE_TARGETS"),
			Category: category,
			Action: func(*cli.Context, []string) error {
				return fmt.Errorf("flag --%s (env var %s) is deprecated, use --%s (env var %s) instead",
					DeprecatedCacheTargetsFlagName, withDeprecatedEnvPrefix(envPrefix, "CACHE_TARGETS"),
					CacheTargetsFlagName, withEnvPrefix(envPrefix, "CACHE_TARGETS"))
			},
		},
		&cli.IntFlag{
			Name:     DeprecatedConcurrentWriteThreads,
			Usage:    "Number of threads spun-up for async secondary storage insertions. (<=0) denotes single threaded insertions where (>0) indicates decoupled writes.",
			Value:    0,
			EnvVars:  withDeprecatedEnvPrefix(envPrefix, "CONCURRENT_WRITE_THREADS"),
			Category: category,
			Action: func(*cli.Context, int) error {
				return fmt.Errorf("flag --%s (env var %s) is deprecated, use --%s (env var %s) instead",
					DeprecatedCacheTargetsFlagName, withDeprecatedEnvPrefix(envPrefix, "CONCURRENT_WRITE_THREADS"),
					CacheTargetsFlagName, withEnvPrefix(envPrefix, "CONCURRENT_WRITE_THREADS"))
			},
		},
	}
}
