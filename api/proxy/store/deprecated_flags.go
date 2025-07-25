package store

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

// The flags here are deprecated and will be removed after v3.0.0
// we leave them here with actions that crash the program to ensure they are not used,
// and to make it easier for users to find the new flags (instead of silently crashing late during execution
// because some flag's env var was changed but the user forgot to update it)
var (
	FallbackTargetsFlagName = withFlagPrefix("fallback-targets")
)

// CLIFlags ... used for EigenDA client configuration
func DeprecatedCLIFlags(envPrefix, category string) []cli.Flag {
	return []cli.Flag{
		&cli.StringSliceFlag{
			Name:     FallbackTargetsFlagName,
			Usage:    "List of read fallback targets to rollover to if cert can't be read from EigenDA.",
			Value:    cli.NewStringSlice(),
			EnvVars:  withEnvPrefix(envPrefix, "FALLBACK_TARGETS"),
			Category: category,
			Action: func(*cli.Context, []string) error {
				return fmt.Errorf("Fallback reads are deprecated in favor of cache reads. "+
					"flag --%s (env var %s) is thus deprecated; use --%s (env var %s) instead.",
					FallbackTargetsFlagName, withEnvPrefix(envPrefix, "FALLBACK_TARGETS"),
					CacheTargetsFlagName, withEnvPrefix(envPrefix, "CACHE_TARGETS"))
			},
			Hidden: true,
		},
	}
}
