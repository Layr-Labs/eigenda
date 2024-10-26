package flags

import (
	"github.com/Layr-Labs/eigenda-proxy/flags/eigendaflags"
	"github.com/Layr-Labs/eigenda-proxy/store/generated_key/memstore"
	"github.com/Layr-Labs/eigenda-proxy/store/precomputed_key/redis"
	"github.com/Layr-Labs/eigenda-proxy/store/precomputed_key/s3"
	"github.com/Layr-Labs/eigenda-proxy/verify"
	"github.com/urfave/cli/v2"

	opservice "github.com/ethereum-optimism/optimism/op-service"
	oplog "github.com/ethereum-optimism/optimism/op-service/log"
	opmetrics "github.com/ethereum-optimism/optimism/op-service/metrics"
)

const (
	EigenDAClientCategory      = "EigenDA Client"
	EigenDADeprecatedCategory  = "DEPRECATED EIGENDA CLIENT FLAGS -- THESE WILL BE REMOVED IN V2.0.0"
	MemstoreFlagsCategory      = "Memstore (for testing purposes - replaces EigenDA backend)"
	RedisCategory              = "Redis Cache/Fallback"
	S3Category                 = "S3 Cache/Fallback"
	VerifierCategory           = "KZG and Cert Verifier"
	VerifierDeprecatedCategory = "DEPRECATED VERIFIER FLAGS -- THESE WILL BE REMOVED IN V2.0.0"
)

const (
	ListenAddrFlagName = "addr"
	PortFlagName       = "port"

	// routing flags
	// TODO: change "routing" --> "secondary"
	FallbackTargetsFlagName = "routing.fallback-targets"
	CacheTargetsFlagName    = "routing.cache-targets"
	ConcurrentWriteThreads  = "routing.concurrent-write-routines"
)

const EnvVarPrefix = "EIGENDA_PROXY"

func prefixEnvVars(name string) []string {
	return opservice.PrefixEnvVar(EnvVarPrefix, name)
}

func CLIFlags() []cli.Flag {
	// TODO: Decompose all flags into constituent parts based on their respective category / usage
	flags := []cli.Flag{
		&cli.StringFlag{
			Name:    ListenAddrFlagName,
			Usage:   "Server listening address",
			Value:   "0.0.0.0",
			EnvVars: prefixEnvVars("ADDR"),
		},
		&cli.IntFlag{
			Name:    PortFlagName,
			Usage:   "Server listening port",
			Value:   3100,
			EnvVars: prefixEnvVars("PORT"),
		},
		&cli.StringSliceFlag{
			Name:    FallbackTargetsFlagName,
			Usage:   "List of read fallback targets to rollover to if cert can't be read from EigenDA.",
			Value:   cli.NewStringSlice(),
			EnvVars: prefixEnvVars("FALLBACK_TARGETS"),
		},
		&cli.StringSliceFlag{
			Name:    CacheTargetsFlagName,
			Usage:   "List of caching targets to use fast reads from EigenDA.",
			Value:   cli.NewStringSlice(),
			EnvVars: prefixEnvVars("CACHE_TARGETS"),
		},
		&cli.IntFlag{
			Name:    ConcurrentWriteThreads,
			Usage:   "Number of threads spun-up for async secondary storage insertions. (<=0) denotes single threaded insertions where (>0) indicates decoupled writes.",
			Value:   0,
			EnvVars: prefixEnvVars("CONCURRENT_WRITE_THREADS"),
		},
	}

	return flags
}

// Flags contains the list of configuration options available to the binary.
var Flags = []cli.Flag{}

func init() {
	Flags = CLIFlags()
	Flags = append(Flags, oplog.CLIFlags(EnvVarPrefix)...)
	Flags = append(Flags, opmetrics.CLIFlags(EnvVarPrefix)...)
	Flags = append(Flags, eigendaflags.CLIFlags(EnvVarPrefix, EigenDAClientCategory)...)
	Flags = append(Flags, eigendaflags.DeprecatedCLIFlags(EnvVarPrefix, EigenDADeprecatedCategory)...)
	Flags = append(Flags, redis.CLIFlags(EnvVarPrefix, RedisCategory)...)
	Flags = append(Flags, s3.CLIFlags(EnvVarPrefix, S3Category)...)
	Flags = append(Flags, memstore.CLIFlags(EnvVarPrefix, MemstoreFlagsCategory)...)
	Flags = append(Flags, verify.CLIFlags(EnvVarPrefix, VerifierCategory)...)
	Flags = append(Flags, verify.DeprecatedCLIFlags(EnvVarPrefix, VerifierDeprecatedCategory)...)
}
