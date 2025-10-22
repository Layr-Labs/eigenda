package config

import (
	"github.com/Layr-Labs/eigenda/api/proxy/config/eigendaflags"
	enabled_apis "github.com/Layr-Labs/eigenda/api/proxy/config/enablement"
	eigenda_v2_flags "github.com/Layr-Labs/eigenda/api/proxy/config/v2/eigendaflags"
	"github.com/Layr-Labs/eigenda/api/proxy/servers/arbitrum_altda"
	"github.com/Layr-Labs/eigenda/api/proxy/servers/rest"
	"github.com/Layr-Labs/eigenda/api/proxy/store"
	"github.com/Layr-Labs/eigenda/api/proxy/store/generated_key/eigenda/verify"

	"github.com/Layr-Labs/eigenda/api/proxy/logging"
	"github.com/Layr-Labs/eigenda/api/proxy/metrics"
	"github.com/Layr-Labs/eigenda/api/proxy/store/generated_key/memstore"
	"github.com/Layr-Labs/eigenda/api/proxy/store/secondary/redis"
	"github.com/Layr-Labs/eigenda/api/proxy/store/secondary/s3"
	"github.com/Layr-Labs/eigenda/api/proxy/telemetry"
	"github.com/urfave/cli/v2"
)

const (
	EnabledAPIsCategory     = "Enabled APIs"
	ProxyRestServerCategory = "Proxy REST API Server (compatible with OP Stack ALT DA and standard commitment clients)"
	ArbCustomDASvrCategory  = "Arbitrum Custom DA JSON RPC Server"

	LoggingFlagsCategory = "Logging"
	MetricsFlagCategory  = "Metrics"
	TelemetryCategory    = "OpenTelemetry Tracing"

	StorageFlagsCategory  = "Storage"
	MemstoreFlagsCategory = "Memstore (for testing purposes - replaces EigenDA backend)"
	S3Category            = "S3 Cache/Fallback"

	EigenDAClientCategory = "EigenDA V1 Client"
	VerifierCategory      = "Cert Verifier (V1 only)"

	EigenDAV2ClientCategory = "EigenDA V2 Client"
	KZGCategory             = "KZG"

	DeprecatedRedisCategory = "Redis Cache/Fallback"
)

// EnvVar prefix added in front of all environment variables accepted by the binary.
// This acts as a namespace to avoid collisions with other binaries.
const GlobalEnvVarPrefix = "EIGENDA_PROXY"

// Flags contains the list of configuration options available to the binary.
var Flags = []cli.Flag{}

func init() {
	Flags = append(Flags, enabled_apis.CLIFlags(EnabledAPIsCategory, GlobalEnvVarPrefix)...)

	Flags = append(Flags, rest.CLIFlags(GlobalEnvVarPrefix, ProxyRestServerCategory)...)
	Flags = append(Flags, arbitrum_altda.CLIFlags(GlobalEnvVarPrefix, ArbCustomDASvrCategory)...)
	Flags = append(Flags, metrics.CLIFlags(GlobalEnvVarPrefix, MetricsFlagCategory)...)
	Flags = append(Flags, telemetry.CLIFlags(GlobalEnvVarPrefix, TelemetryCategory)...)

	Flags = append(Flags, logging.CLIFlags(GlobalEnvVarPrefix, LoggingFlagsCategory)...)
	Flags = append(Flags, eigendaflags.CLIFlags(GlobalEnvVarPrefix, EigenDAClientCategory)...)
	Flags = append(Flags, eigenda_v2_flags.CLIFlags(GlobalEnvVarPrefix, EigenDAV2ClientCategory)...)
	Flags = append(Flags, store.CLIFlags(GlobalEnvVarPrefix, StorageFlagsCategory)...)
	Flags = append(Flags, s3.CLIFlags(GlobalEnvVarPrefix, S3Category)...)
	Flags = append(Flags, memstore.CLIFlags(GlobalEnvVarPrefix, MemstoreFlagsCategory)...)
	Flags = append(Flags, verify.VerifierCLIFlags(GlobalEnvVarPrefix, VerifierCategory)...)
	Flags = append(Flags, verify.KZGCLIFlags(GlobalEnvVarPrefix, KZGCategory)...)

	Flags = append(Flags, metrics.DeprecatedCLIFlags(GlobalEnvVarPrefix, MetricsFlagCategory)...)
	Flags = append(Flags, eigendaflags.DeprecatedCLIFlags(GlobalEnvVarPrefix, EigenDAClientCategory)...)
	Flags = append(Flags, eigenda_v2_flags.DeprecatedCLIFlags(GlobalEnvVarPrefix, EigenDAV2ClientCategory)...)
	Flags = append(Flags, verify.DeprecatedCLIFlags(GlobalEnvVarPrefix, VerifierCategory)...)
	Flags = append(Flags, store.DeprecatedCLIFlags(GlobalEnvVarPrefix, StorageFlagsCategory)...)
	Flags = append(Flags, redis.DeprecatedCLIFlags(GlobalEnvVarPrefix, DeprecatedRedisCategory)...)
}
