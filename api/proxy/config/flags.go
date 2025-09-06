package config

import (
	"github.com/Layr-Labs/eigenda/api/proxy/config/consts"
	"github.com/Layr-Labs/eigenda/api/proxy/config/eigendaflags"
	eigenda_v2_flags "github.com/Layr-Labs/eigenda/api/proxy/config/v2/eigendaflags"
	"github.com/Layr-Labs/eigenda/api/proxy/server"
	"github.com/Layr-Labs/eigenda/api/proxy/store"
	"github.com/Layr-Labs/eigenda/api/proxy/store/generated_key/eigenda/verify"

	"github.com/Layr-Labs/eigenda/api/proxy/logging"
	"github.com/Layr-Labs/eigenda/api/proxy/metrics"
	"github.com/Layr-Labs/eigenda/api/proxy/store/generated_key/memstore"
	"github.com/Layr-Labs/eigenda/api/proxy/store/secondary/redis"
	"github.com/Layr-Labs/eigenda/api/proxy/store/secondary/s3"
	"github.com/urfave/cli/v2"
)

// Flags contains the list of configuration options available to the binary.
var Flags = []cli.Flag{}

func init() {
	Flags = append(Flags, server.CLIFlags(consts.GlobalEnvVarPrefix, consts.ProxyServerCategory)...)
	Flags = append(Flags, logging.CLIFlags(consts.GlobalEnvVarPrefix, consts.LoggingFlagsCategory)...)
	Flags = append(Flags, metrics.CLIFlags(consts.GlobalEnvVarPrefix, consts.MetricsFlagCategory)...)
	Flags = append(Flags, eigendaflags.CLIFlags(consts.GlobalEnvVarPrefix, consts.EigenDAClientCategory)...)
	Flags = append(Flags, eigenda_v2_flags.CLIFlags(consts.GlobalEnvVarPrefix, consts.EigenDAV2ClientCategory)...)
	Flags = append(Flags, store.CLIFlags(consts.GlobalEnvVarPrefix, consts.StorageFlagsCategory)...)
	Flags = append(Flags, s3.CLIFlags(consts.GlobalEnvVarPrefix, consts.S3Category)...)
	Flags = append(Flags, memstore.CLIFlags(consts.GlobalEnvVarPrefix, consts.MemstoreFlagsCategory)...)
	Flags = append(Flags, verify.VerifierCLIFlags(consts.GlobalEnvVarPrefix, consts.VerifierCategory)...)
	Flags = append(Flags, verify.KZGCLIFlags(consts.GlobalEnvVarPrefix, consts.KZGCategory)...)

	Flags = append(Flags, eigendaflags.DeprecatedCLIFlags(consts.GlobalEnvVarPrefix, consts.EigenDAClientCategory)...)
	Flags = append(Flags, eigenda_v2_flags.DeprecatedCLIFlags(
		consts.GlobalEnvVarPrefix, consts.EigenDAV2ClientCategory)...)
	Flags = append(Flags, verify.DeprecatedCLIFlags(consts.GlobalEnvVarPrefix, consts.VerifierCategory)...)
	Flags = append(Flags, store.DeprecatedCLIFlags(consts.GlobalEnvVarPrefix, consts.StorageFlagsCategory)...)
	Flags = append(Flags, redis.DeprecatedCLIFlags(consts.GlobalEnvVarPrefix, consts.DeprecatedRedisCategory)...)
}
