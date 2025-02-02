package flags

import (
	"github.com/Layr-Labs/eigenda-proxy/flags/eigendaflags"
	"github.com/Layr-Labs/eigenda-proxy/logging"
	"github.com/Layr-Labs/eigenda-proxy/metrics"
	"github.com/Layr-Labs/eigenda-proxy/store"
	"github.com/Layr-Labs/eigenda-proxy/store/generated_key/memstore"
	"github.com/Layr-Labs/eigenda-proxy/store/precomputed_key/redis"
	"github.com/Layr-Labs/eigenda-proxy/store/precomputed_key/s3"
	"github.com/Layr-Labs/eigenda-proxy/verify"

	"github.com/urfave/cli/v2"

	"github.com/Layr-Labs/eigenda-proxy/common"
)

const (
	EigenDAClientCategory      = "EigenDA Client"
	LoggingFlagsCategory       = "Logging"
	MetricsFlagCategory        = "Metrics"
	EigenDADeprecatedCategory  = "DEPRECATED EIGENDA CLIENT FLAGS -- THESE WILL BE REMOVED IN V2.0.0"
	MemstoreFlagsCategory      = "Memstore (for testing purposes - replaces EigenDA backend)"
	StorageFlagsCategory       = "Storage"
	StorageDeprecatedCategory  = "DEPRECATED STORAGE FLAGS -- THESE WILL BE REMOVED IN V2.0.0"
	RedisCategory              = "Redis Cache/Fallback"
	S3Category                 = "S3 Cache/Fallback"
	VerifierCategory           = "KZG and Cert Verifier"
	VerifierDeprecatedCategory = "DEPRECATED VERIFIER FLAGS -- THESE WILL BE REMOVED IN V2.0.0"
)

const (
	ListenAddrFlagName = "addr"
	PortFlagName       = "port"
)

func CLIFlags() []cli.Flag {
	// TODO: Decompose all flags into constituent parts based on their respective category / usage
	flags := []cli.Flag{
		&cli.StringFlag{
			Name:    ListenAddrFlagName,
			Usage:   "Server listening address",
			Value:   "0.0.0.0",
			EnvVars: common.PrefixEnvVar(common.GlobalPrefix, "ADDR"),
		},
		&cli.IntFlag{
			Name:    PortFlagName,
			Usage:   "Server listening port",
			Value:   3100,
			EnvVars: common.PrefixEnvVar(common.GlobalPrefix, "PORT"),
		},
	}

	return flags
}

// Flags contains the list of configuration options available to the binary.
var Flags = []cli.Flag{}

func init() {
	Flags = CLIFlags()
	Flags = append(Flags, logging.CLIFlags(common.GlobalPrefix, LoggingFlagsCategory)...)
	Flags = append(Flags, metrics.CLIFlags(common.GlobalPrefix, MetricsFlagCategory)...)
	Flags = append(Flags, eigendaflags.CLIFlags(common.GlobalPrefix, EigenDAClientCategory)...)
	Flags = append(Flags, eigendaflags.DeprecatedCLIFlags(common.GlobalPrefix, EigenDADeprecatedCategory)...)
	Flags = append(Flags, store.CLIFlags(common.GlobalPrefix, StorageFlagsCategory)...)
	Flags = append(Flags, store.DeprecatedCLIFlags(common.GlobalPrefix, StorageDeprecatedCategory)...)
	Flags = append(Flags, redis.CLIFlags(common.GlobalPrefix, RedisCategory)...)
	Flags = append(Flags, s3.CLIFlags(common.GlobalPrefix, S3Category)...)
	Flags = append(Flags, memstore.CLIFlags(common.GlobalPrefix, MemstoreFlagsCategory)...)
	Flags = append(Flags, verify.CLIFlags(common.GlobalPrefix, VerifierCategory)...)
	Flags = append(Flags, verify.DeprecatedCLIFlags(common.GlobalPrefix, VerifierDeprecatedCategory)...)
}
