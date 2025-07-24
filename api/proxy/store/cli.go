package store

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Layr-Labs/eigenda/api/proxy/common"
	"github.com/urfave/cli/v2"
)

var (
	BackendsToEnableFlagName = withFlagPrefix("backends-to-enable")
	DispersalBackendFlagName = withFlagPrefix("dispersal-backend")
	FallbackTargetsFlagName  = withFlagPrefix("fallback-targets")
	CacheTargetsFlagName     = withFlagPrefix("cache-targets")
	ConcurrentWriteThreads   = withFlagPrefix("concurrent-write-routines")
	WriteOnCacheMissFlagName = withFlagPrefix("write-on-cache-miss")
)

func withFlagPrefix(s string) string {
	return "storage." + s
}

func withEnvPrefix(envPrefix, s string) []string {
	return []string{envPrefix + "_STORAGE_" + s}
}

// CLIFlags ... used for storage configuration
// category is used to group the flags in the help output (see https://cli.urfave.org/v2/examples/flags/#grouping)
func CLIFlags(envPrefix, category string) []cli.Flag {
	return []cli.Flag{
		&cli.StringSliceFlag{
			Name:     BackendsToEnableFlagName,
			Usage:    "Comma separated list of eigenDA backends to enable (e.g. V1,V2)",
			EnvVars:  withEnvPrefix(envPrefix, "BACKENDS_TO_ENABLE"),
			Value:    cli.NewStringSlice("V1"),
			Category: category,
			Required: false,
		},
		&cli.StringFlag{
			Name:     DispersalBackendFlagName,
			Usage:    "Target EigenDA backend version for blob dispersal (e.g. V1 or V2).",
			EnvVars:  withEnvPrefix(envPrefix, "DISPERSAL_BACKEND"),
			Category: category,
			Required: false,
			Value:    "V1",
		},
		&cli.StringSliceFlag{
			Name:     FallbackTargetsFlagName,
			Usage:    "List of read fallback targets to rollover to if cert can't be read from EigenDA.",
			Value:    cli.NewStringSlice(),
			EnvVars:  withEnvPrefix(envPrefix, "FALLBACK_TARGETS"),
			Category: category,
		},
		&cli.StringSliceFlag{
			Name:     CacheTargetsFlagName,
			Usage:    "List of caching targets to use fast reads from EigenDA.",
			Value:    cli.NewStringSlice(),
			EnvVars:  withEnvPrefix(envPrefix, "CACHE_TARGETS"),
			Category: category,
		},
		&cli.IntFlag{
			Name:     ConcurrentWriteThreads,
			Usage:    "Number of threads spun-up for async secondary storage insertions. (<=0) denotes single threaded insertions where (>0) indicates decoupled writes.",
			Value:    0,
			EnvVars:  withEnvPrefix(envPrefix, "CONCURRENT_WRITE_THREADS"),
			Category: category,
		},
		&cli.BoolFlag{
			Name:     WriteOnCacheMissFlagName,
			Usage:    "While doing a GET, write to the secondary storage if the cert/blob is not found in the cache but is found in EigenDA.",
			Value:    false,
			EnvVars:  withEnvPrefix(envPrefix, "WRITE_ON_CACHE_MISS"),
			Category: category,
		},
	}
}

func ReadConfig(ctx *cli.Context) (Config, error) {
	backendStrings := ctx.StringSlice(BackendsToEnableFlagName)
	if len(backendStrings) == 0 {
		return Config{}, errors.New("backends must not be empty")
	}

	backends := make([]common.EigenDABackend, 0, len(backendStrings))
	for _, backendString := range backendStrings {
		backend, err := common.StringToEigenDABackend(backendString)
		if err != nil {
			return Config{}, fmt.Errorf("string to eigenDA backend: %w", err)
		}
		backends = append(backends, backend)
	}

	dispersalBackend, err := common.StringToEigenDABackend(ctx.String(DispersalBackendFlagName))
	if err != nil {
		return Config{}, fmt.Errorf("string to eigenDA backend: %w", err)
	}

	// We need to filter the cache targets and fallback targets to remove empty strings,
	// since our code downstream doesn't work well with empty strings.
	// Specifically, if the env var is simply set to nothing like `EIGENDA_PROXY_STORAGE_CACHE_TARGETS=`,
	// it will result in an empty string being added to the slice
	// for some reason... seems like a bug in urfave/cli?
	cacheTargets := ctx.StringSlice(CacheTargetsFlagName)
	filteredCacheTargets := make([]string, 0, len(cacheTargets))
	for _, target := range cacheTargets {
		if target != "" {
			filteredCacheTargets = append(filteredCacheTargets, target)
		}
	}

	fallbackTargets := ctx.StringSlice(FallbackTargetsFlagName)
	filteredFallbackTargets := make([]string, 0, len(fallbackTargets))
	for _, target := range fallbackTargets {
		if target != "" {
			filteredFallbackTargets = append(filteredFallbackTargets, target)
		}
	}

	return Config{
		BackendsToEnable: backends,
		DispersalBackend: dispersalBackend,
		AsyncPutWorkers:  ctx.Int(ConcurrentWriteThreads),
		FallbackTargets:  filteredFallbackTargets,
		CacheTargets:     filteredCacheTargets,
		WriteOnCacheMiss: ctx.Bool(WriteOnCacheMissFlagName),
	}, nil
}
