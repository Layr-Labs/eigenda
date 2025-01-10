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
	CustomQuorumNumbersFlag = cli.IntSliceFlag{
		Name:     common.PrefixFlag(FlagPrefix, "custom-quorum-numbers"),
		Usage:    "Custom quorum numbers to use for the traffic generator.",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envPrefix, "CUSTOM_QUORUM_NUMBERS"),
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

	/* Common Configuration. */

	InstanceLaunchIntervalFlag = cli.DurationFlag{
		Name:     common.PrefixFlag(FlagPrefix, "instance-launch-interva"),
		Usage:    "Duration between generator instance launches.",
		Required: false,
		Value:    1 * time.Second,
		EnvVar:   common.PrefixEnvVar(envPrefix, "INSTANCE_LAUNCH_INTERVAL"),
	}

	/* Configuration for the blob writers. */
	NumWriteInstancesFlag = cli.UintFlag{
		Name:     common.PrefixFlag(FlagPrefix, "num-write-instances"),
		Usage:    "Number of writer instances producing traffic to run in parallel.",
		Required: false,
		Value:    1,
		EnvVar:   common.PrefixEnvVar(envPrefix, "NUM_WRITE_INSTANCES"),
	}
	WriteRequestIntervalFlag = cli.DurationFlag{
		Name:     common.PrefixFlag(FlagPrefix, "write-request-interval"),
		Usage:    "Time between write requests.",
		Required: false,
		Value:    30 * time.Second,
		EnvVar:   common.PrefixEnvVar(envPrefix, "WRITE_REQUEST_INTERVAL"),
	}
	DataSizeFlag = cli.Uint64Flag{
		Name:     common.PrefixFlag(FlagPrefix, "data-size"),
		Usage:    "Size of the data blob.",
		Required: false,
		Value:    1024,
		EnvVar:   common.PrefixEnvVar(envPrefix, "DATA_SIZE"),
	}
	RandomizeBlobsFlag = cli.BoolFlag{
		Name:     common.PrefixFlag(FlagPrefix, "randomize-blobs"),
		Usage:    "If set, do not randomize blobs.",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envPrefix, "RANDOMIZE_BLOBS"),
	}
	WriteTimeoutFlag = cli.DurationFlag{
		Name:     common.PrefixFlag(FlagPrefix, "write-timeout"),
		Usage:    "Amount of time to wait for a blob to be written.",
		Required: false,
		Value:    10 * time.Second,
		EnvVar:   common.PrefixEnvVar(envPrefix, "WRITE_TIMEOUT"),
	}
	NodeClientTimeoutFlag = cli.DurationFlag{
		Name:     common.PrefixFlag(FlagPrefix, "node-client-timeout"),
		Usage:    "The timeout for the node client.",
		Required: false,
		Value:    10 * time.Second,
		EnvVar:   common.PrefixEnvVar(envPrefix, "NODE_CLIENT_TIMEOUT"),
	}
)

var requiredFlags = []cli.Flag{
	HostnameFlag,
	GrpcPortFlag,
}

var optionalFlags = []cli.Flag{
	TimeoutFlag,
	RandomizeBlobsFlag,
	InstanceLaunchIntervalFlag,
	UseSecureGrpcFlag,
	SignerPrivateKeyFlag,
	CustomQuorumNumbersFlag,
	NumWriteInstancesFlag,
	WriteRequestIntervalFlag,
	DataSizeFlag,
	DisableTLSFlag,
	MetricsHTTPPortFlag,
	NodeClientTimeoutFlag,
	WriteTimeoutFlag,
}

// Flags contains the list of configuration options available to the binary.
var Flags []cli.Flag

func init() {
	Flags = append(requiredFlags, optionalFlags...)
	Flags = append(Flags, common.LoggerCLIFlags(envPrefix, FlagPrefix)...)
	Flags = append(Flags, indexer.CLIFlags(envPrefix)...)
}
