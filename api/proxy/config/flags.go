package config

import (
	"github.com/Layr-Labs/eigenda-proxy/config/eigendaflags"
	eigenda_v2_flags "github.com/Layr-Labs/eigenda-proxy/config/v2/eigendaflags"
	"github.com/Layr-Labs/eigenda-proxy/server"
	"github.com/Layr-Labs/eigenda-proxy/store"
	"github.com/Layr-Labs/eigenda-proxy/store/generated_key/eigenda/verify"

	"github.com/Layr-Labs/eigenda-proxy/logging"
	"github.com/Layr-Labs/eigenda-proxy/metrics"
	"github.com/Layr-Labs/eigenda-proxy/store/generated_key/memstore"
	"github.com/Layr-Labs/eigenda-proxy/store/secondary/redis"
	"github.com/Layr-Labs/eigenda-proxy/store/secondary/s3"
	"github.com/urfave/cli/v2"
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
	ProxyServerCategory     = "Proxy Server"
)

// EnvVar prefix added in front of all environment variables accepted by the binary.
// This acts as a namespace to avoid collisions with other binaries.
const GlobalEnvVarPrefix = "EIGENDA_PROXY"

// Flags contains the list of configuration options available to the binary.
var Flags = []cli.Flag{}

func init() {
	Flags = append(Flags, server.CLIFlags(GlobalEnvVarPrefix, ProxyServerCategory)...)
	Flags = append(Flags, logging.CLIFlags(GlobalEnvVarPrefix, LoggingFlagsCategory)...)
	Flags = append(Flags, metrics.CLIFlags(GlobalEnvVarPrefix, MetricsFlagCategory)...)
	Flags = append(Flags, eigendaflags.CLIFlags(GlobalEnvVarPrefix, EigenDAClientCategory)...)
	Flags = append(Flags, eigenda_v2_flags.CLIFlags(GlobalEnvVarPrefix, EigenDAV2ClientCategory)...)
	Flags = append(Flags, store.CLIFlags(GlobalEnvVarPrefix, StorageFlagsCategory)...)
	Flags = append(Flags, redis.CLIFlags(GlobalEnvVarPrefix, RedisCategory)...)
	Flags = append(Flags, s3.CLIFlags(GlobalEnvVarPrefix, S3Category)...)
	Flags = append(Flags, memstore.CLIFlags(GlobalEnvVarPrefix, MemstoreFlagsCategory)...)
	Flags = append(Flags, verify.VerifierCLIFlags(GlobalEnvVarPrefix, VerifierCategory)...)
	Flags = append(Flags, verify.KZGCLIFlags(GlobalEnvVarPrefix, KZGCategory)...)

	Flags = append(Flags, eigendaflags.DeprecatedCLIFlags(GlobalEnvVarPrefix, EigenDAClientCategory)...)
	Flags = append(Flags, verify.DeprecatedCLIFlags(GlobalEnvVarPrefix, VerifierCategory)...)
	Flags = append(Flags, store.DeprecatedCLIFlags(GlobalEnvVarPrefix, StorageFlagsCategory)...)
}
