package config

import (
	"github.com/Layr-Labs/eigenda-proxy/config/eigendaflags"
	eigenda_v2_flags "github.com/Layr-Labs/eigenda-proxy/config/v2/eigendaflags"
	"github.com/Layr-Labs/eigenda-proxy/store"
	"github.com/Layr-Labs/eigenda-proxy/verify"

	"github.com/Layr-Labs/eigenda-proxy/logging"
	"github.com/Layr-Labs/eigenda-proxy/metrics"
	"github.com/Layr-Labs/eigenda-proxy/store/generated_key/memstore"
	"github.com/Layr-Labs/eigenda-proxy/store/precomputed_key/redis"
	"github.com/Layr-Labs/eigenda-proxy/store/precomputed_key/s3"
	"github.com/urfave/cli/v2"

	"github.com/Layr-Labs/eigenda-proxy/common"
)

const (
	EigenDAClientCategory   = "EigenDA V1 Client"
	EigenDAV2ClientCategory = "EigenDA V2 Client"
	LoggingFlagsCategory    = "Logging"
	MetricsFlagCategory     = "Metrics"
	MemstoreFlagsCategory   = "Memstore (for testing purposes - replaces EigenDA backend)"
	StorageFlagsCategory    = "Storage"
	RedisCategory           = "Redis Cache/Fallback"
	S3Category              = "S3 Cache/Fallback"
	VerifierCategory        = "Cert Verifier (V1 only)"
	KZGCategory             = "KZG"
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
	Flags = append(Flags, eigenda_v2_flags.CLIFlags(common.GlobalPrefix, EigenDAV2ClientCategory)...)
	Flags = append(Flags, store.CLIFlags(common.GlobalPrefix, StorageFlagsCategory)...)
	Flags = append(Flags, redis.CLIFlags(common.GlobalPrefix, RedisCategory)...)
	Flags = append(Flags, s3.CLIFlags(common.GlobalPrefix, S3Category)...)
	Flags = append(Flags, memstore.CLIFlags(common.GlobalPrefix, MemstoreFlagsCategory)...)
	Flags = append(Flags, verify.VerifierCLIFlags(common.GlobalPrefix, VerifierCategory)...)
	Flags = append(Flags, verify.KZGCLIFlags(common.GlobalPrefix, KZGCategory)...)

	Flags = append(Flags, eigendaflags.DeprecatedCLIFlags(common.GlobalPrefix, EigenDAClientCategory)...)
	Flags = append(Flags, verify.DeprecatedCLIFlags(common.GlobalPrefix, VerifierCategory)...)
	Flags = append(Flags, store.DeprecatedCLIFlags(common.GlobalPrefix, StorageFlagsCategory)...)
}
