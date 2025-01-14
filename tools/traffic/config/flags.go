package config

import (
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/indexer"
	"github.com/urfave/cli"
)

const (
	FlagPrefix = "traffic-generator"
	envPrefix  = "TRAFFIC_GENERATOR"
)

var (
	/* Configuration for DA clients. */

	HostnameFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "disperser-hostname"),
		Usage:    "Hostname at which disperser service is available.",
		Required: true,
		EnvVar:   common.PrefixEnvVar(envPrefix, "DISPERSER_HOSTNAME"),
	}
	GrpcPortFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "disperser-port"),
		Usage:    "Port at which a disperser listens for grpc calls.",
		Required: true,
		EnvVar:   common.PrefixEnvVar(envPrefix, "DISPERSER_PORT"),
	}
	TimeoutFlag = cli.DurationFlag{
		Name:     common.PrefixFlag(FlagPrefix, "timeout"),
		Usage:    "Amount of time to wait for grpc.",
		Required: false,
		Value:    10 * time.Second,
		EnvVar:   common.PrefixEnvVar(envPrefix, "DISPERSER_TIMEOUT"),
	}
	UseSecureGrpcFlag = cli.BoolFlag{
		Name:     common.PrefixFlag(FlagPrefix, "use-secure-grpc"),
		Usage:    "Whether to use secure grpc.",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envPrefix, "DISPERSER_USE_SECURE_GRPC"),
	}
	SignerPrivateKeyFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "signer-private-key-hex"),
		Usage:    "Private key to use for signing requests.",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envPrefix, "SIGNER_PRIVATE_KEY_HEX"),
	}
	DisableTLSFlag = cli.BoolFlag{
		Name:     common.PrefixFlag(FlagPrefix, "disable-tls"),
		Usage:    "Whether to disable TLS for an insecure connection.",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envPrefix, "DISABLE_TLS"),
	}
	MetricsHTTPPortFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "metrics-http-port"),
		Usage:    "Port at which to expose metrics.",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envPrefix, "METRICS_HTTP_PORT"),
	}
	NodeClientTimeoutFlag = cli.DurationFlag{
		Name:     common.PrefixFlag(FlagPrefix, "node-client-timeout"),
		Usage:    "The timeout for the node client.",
		Required: false,
		Value:    10 * time.Second,
		EnvVar:   common.PrefixEnvVar(envPrefix, "NODE_CLIENT_TIMEOUT"),
	}
	RuntimeConfigPathFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "runtime-config-path"),
		Usage:    "Path to the runtime configuration file that defines writer groups.",
		Required: true,
		EnvVar:   common.PrefixEnvVar(envPrefix, "RUNTIME_CONFIG_PATH"),
	}
)

var requiredFlags = []cli.Flag{
	HostnameFlag,
	GrpcPortFlag,
	RuntimeConfigPathFlag,
}

var optionalFlags = []cli.Flag{
	TimeoutFlag,
	UseSecureGrpcFlag,
	SignerPrivateKeyFlag,
	DisableTLSFlag,
	MetricsHTTPPortFlag,
	NodeClientTimeoutFlag,
}

// Flags contains the list of configuration options available to the binary.
var Flags []cli.Flag

func init() {
	Flags = append(requiredFlags, optionalFlags...)
	Flags = append(Flags, common.LoggerCLIFlags(envPrefix, FlagPrefix)...)
	Flags = append(Flags, indexer.CLIFlags(envPrefix)...)
}
