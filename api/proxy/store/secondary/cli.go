package secondary

import (
	"github.com/urfave/cli/v2"
)

const (
	// ErrorOnSecondaryInsertFailureFlagName is the CLI flag name for enabling strict error handling
	// on secondary storage insertion failures.
	ErrorOnSecondaryInsertFailureFlagName = "secondary.error-on-secondary-insert-failure"
)

func withEnvPrefix(envPrefix, s string) []string {
	return []string{envPrefix + "_SECONDARY_" + s}
}

// CLIFlags returns CLI flags for secondary storage configuration.
// category is used to group the flags in the help output (see https://cli.urfave.org/v2/examples/flags/#grouping)
func CLIFlags(envPrefix, category string) []cli.Flag {
	return []cli.Flag{
		&cli.BoolFlag{
			Name: ErrorOnSecondaryInsertFailureFlagName,
			Usage: "Return 500 error if any secondary storage write fails, ensuring all-or-nothing redundancy. " +
				"Cannot be used with async writes (concurrent-write-routines > 0).",
			Value:    false,
			EnvVars:  withEnvPrefix(envPrefix, "ERROR_ON_SECONDARY_INSERT_FAILURE"),
			Category: category,
		},
	}
}

// ReadConfig reads the secondary storage configuration from CLI context.
func ReadConfig(ctx *cli.Context) Config {
	return Config{
		ErrorOnSecondaryInsertFailure: ctx.Bool(ErrorOnSecondaryInsertFailureFlagName),
	}
}

// Config holds configuration for secondary storage behavior.
type Config struct {
	// ErrorOnSecondaryInsertFailure, when true, causes secondary storage write failures
	// to propagate as errors to the client (HTTP 500), rather than being silently logged.
	// This ensures all-or-nothing semantics for data redundancy.
	//
	// IMPORTANT: This must be false when AsyncPutWorkers > 0, as async writes cannot
	// return errors to the client.
	ErrorOnSecondaryInsertFailure bool
}
