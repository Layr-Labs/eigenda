package store

import (
	"github.com/Layr-Labs/eigenda-proxy/store/precomputed_key/redis"
	"github.com/Layr-Labs/eigenda-proxy/store/precomputed_key/s3"
	"github.com/urfave/cli/v2"
)

var (
	FallbackTargetsFlagName = withFlagPrefix("fallback-targets")
	CacheTargetsFlagName    = withFlagPrefix("cache-targets")
	ConcurrentWriteThreads  = withFlagPrefix("concurrent-write-routines")
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
	}
}

func ReadConfig(ctx *cli.Context) Config {
	return Config{
		AsyncPutWorkers: ctx.Int(ConcurrentWriteThreads),
		FallbackTargets: ctx.StringSlice(FallbackTargetsFlagName),
		CacheTargets:    ctx.StringSlice(CacheTargetsFlagName),
		RedisConfig:     redis.ReadConfig(ctx),
		S3Config:        s3.ReadConfig(ctx),
	}
}
