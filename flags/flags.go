package flags

import (
	"github.com/Layr-Labs/eigenda-proxy/flags/eigendaflags"
	"github.com/Layr-Labs/eigenda-proxy/store/generated_key/memstore"
	"github.com/Layr-Labs/eigenda-proxy/store/precomputed_key/redis"
	"github.com/Layr-Labs/eigenda-proxy/store/precomputed_key/s3"
	"github.com/urfave/cli/v2"

	opservice "github.com/ethereum-optimism/optimism/op-service"
	oplog "github.com/ethereum-optimism/optimism/op-service/log"
	opmetrics "github.com/ethereum-optimism/optimism/op-service/metrics"
)

const (
	EigenDAClientCategory = "EigenDA Client"
	MemstoreFlagsCategory = "Memstore (replaces EigenDA when enabled)"
	RedisCategory         = "Redis Cache/Fallback"
	S3Category            = "S3 Cache/Fallback"
)

const (
	ListenAddrFlagName = "addr"
	PortFlagName       = "port"

	// routing flags
	FallbackTargetsFlagName = "routing.fallback-targets"
	CacheTargetsFlagName    = "routing.cache-targets"
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
			Usage:   "server listening address",
			Value:   "0.0.0.0",
			EnvVars: prefixEnvVars("ADDR"),
		},
		&cli.IntFlag{
			Name:    PortFlagName,
			Usage:   "server listening port",
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
	Flags = append(Flags, redis.CLIFlags(EnvVarPrefix, RedisCategory)...)
	Flags = append(Flags, s3.CLIFlags(EnvVarPrefix, S3Category)...)
	Flags = append(Flags, memstore.CLIFlags(EnvVarPrefix, MemstoreFlagsCategory)...)
}
