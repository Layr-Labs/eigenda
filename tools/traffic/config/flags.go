package config

import (
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/core/thegraph"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/Layr-Labs/eigenda/indexer"
	"github.com/Layr-Labs/eigenda/retriever/flags"
	"time"

	"github.com/Layr-Labs/eigenda/common"
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

	MetricsBlacklistFlag = cli.StringSliceFlag{
		Name:     common.PrefixFlag(FlagPrefix, "metrics-blacklist"),
		Usage:    "Any metric with a label exactly matching this string will not be sent to the metrics server.",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envPrefix, "METRICS_BLACKLIST"),
	}

	MetricsFuzzyBlacklistFlag = cli.StringSliceFlag{
		Name:     common.PrefixFlag(FlagPrefix, "metrics-fuzzy-blacklist"),
		Usage:    "Any metric that contains any string in this list will not be sent to the metrics server.",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envPrefix, "METRICS_FUZZY_BLACKLIST"),
	}

	/* Configuration for the blob writer. */

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
	UniformBlobsFlag = cli.BoolFlag{
		Name:     common.PrefixFlag(FlagPrefix, "uniform-blobs"),
		Usage:    "If set, do not randomize blobs.",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envPrefix, "UNIFORM_BLOBS"),
	}
	WriteTimeoutFlag = cli.DurationFlag{
		Name:     common.PrefixFlag(FlagPrefix, "write-timeout"),
		Usage:    "Amount of time to wait for a blob to be written.",
		Required: false,
		Value:    10 * time.Second,
		EnvVar:   common.PrefixEnvVar(envPrefix, "WRITE_TIMEOUT"),
	}
	TheGraphUrlFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "the-graph-url"),
		Usage:    "URL of the subgraph instance.",
		Required: true,
		EnvVar:   common.PrefixEnvVar(envPrefix, "THE_GRAPH_URL"),
	}
	TheGraphPullIntervalFlag = cli.DurationFlag{
		Name:     common.PrefixFlag(FlagPrefix, "the-graph-pull-interval"),
		Usage:    "Interval at which to pull data from the subgraph.",
		Required: false,
		Value:    100 * time.Millisecond,
		EnvVar:   common.PrefixEnvVar(envPrefix, "THE_GRAPH_PULL_INTERVAL"),
	}
	TheGraphRetriesFlag = cli.UintFlag{
		Name:     common.PrefixFlag(FlagPrefix, "the-graph-retries"),
		Usage:    "Number of times to retry a subgraph request.",
		Required: false,
		Value:    5,
		EnvVar:   common.PrefixEnvVar(envPrefix, "THE_GRAPH_RETRIES"),
	}
	NodeClientTimeoutFlag = cli.DurationFlag{
		Name:     common.PrefixFlag(FlagPrefix, "node-client-timeout"),
		Usage:    "The timeout for the node client.",
		Required: false,
		Value:    10 * time.Second,
		EnvVar:   common.PrefixEnvVar(envPrefix, "NODE_CLIENT_TIMEOUT"),
	}

	/* Configuration for the blob validator. */

	VerifierIntervalFlag = cli.DurationFlag{
		Name:     common.PrefixFlag(FlagPrefix, "verifier-interval"),
		Usage:    "Amount of time between verifier checks.",
		Required: false,
		Value:    time.Second,
		EnvVar:   common.PrefixEnvVar(envPrefix, "VERIFIER_INTERVAL"),
	}
	GetBlobStatusTimeoutFlag = cli.DurationFlag{
		Name:     common.PrefixFlag(FlagPrefix, "get-blob-status-timeout"),
		Usage:    "Amount of time to wait for a blob status to be fetched.",
		Required: false,
		Value:    5 * time.Second,
		EnvVar:   common.PrefixEnvVar(envPrefix, "GET_BLOB_STATUS_TIMEOUT"),
	}
	VerificationChannelCapacityFlag = cli.UintFlag{
		Name:     common.PrefixFlag(FlagPrefix, "verification-channel-capacity"),
		Usage:    "Size of the channel used to communicate between the writer and verifier.",
		Required: false,
		Value:    1000,
		EnvVar:   common.PrefixEnvVar(envPrefix, "VERIFICATION_CHANNEL_CAPACITY"),
	}

	/* Configuration for the blob reader. */

	NumReadInstancesFlag = cli.UintFlag{
		Name:     common.PrefixFlag(FlagPrefix, "num-read-instances"),
		Usage:    "Number of reader instances producing traffic to run in parallel.",
		Required: false,
		Value:    1,
		EnvVar:   common.PrefixEnvVar(envPrefix, "NUM_READ_INSTANCES"),
	}
	ReadRequestIntervalFlag = cli.DurationFlag{
		Name:     common.PrefixFlag(FlagPrefix, "read-request-interval"),
		Usage:    "Time between read requests.",
		Required: false,
		Value:    time.Second / 5,
		EnvVar:   common.PrefixEnvVar(envPrefix, "READ_REQUEST_INTERVAL"),
	}
	RequiredDownloadsFlag = cli.Float64Flag{
		Name: common.PrefixFlag(FlagPrefix, "required-downloads"),
		Usage: "Number of required downloads. Numbers between 0.0 and 1.0 are treated as probabilities, " +
			"numbers greater than 1.0 are treated as the number of downloads. -1 allows unlimited downloads.",
		Required: false,
		Value:    3.0,
		EnvVar:   common.PrefixEnvVar(envPrefix, "REQUIRED_DOWNLOADS"),
	}
	FetchBatchHeaderTimeoutFlag = cli.DurationFlag{
		Name:     common.PrefixFlag(FlagPrefix, "fetch-batch-header-timeout"),
		Usage:    "Amount of time to wait for a batch header to be fetched.",
		Required: false,
		Value:    5 * time.Second,
		EnvVar:   common.PrefixEnvVar(envPrefix, "FETCH_BATCH_HEADER_TIMEOUT"),
	}
	RetrieveBlobChunksTimeoutFlag = cli.DurationFlag{
		Name:     common.PrefixFlag(FlagPrefix, "retrieve-blob-chunks-timeout"),
		Usage:    "Amount of time to wait for a blob to be retrieved.",
		Required: false,
		Value:    5 * time.Second,
		EnvVar:   common.PrefixEnvVar(envPrefix, "RETRIEVE_BLOB_CHUNKS_TIMEOUT"),
	}
)

var requiredFlags = []cli.Flag{
	HostnameFlag,
	GrpcPortFlag,
	TheGraphUrlFlag,
}

var optionalFlags = []cli.Flag{
	TimeoutFlag,
	UniformBlobsFlag,
	InstanceLaunchIntervalFlag,
	UseSecureGrpcFlag,
	SignerPrivateKeyFlag,
	CustomQuorumNumbersFlag,
	NumWriteInstancesFlag,
	WriteRequestIntervalFlag,
	DataSizeFlag,
	NumReadInstancesFlag,
	ReadRequestIntervalFlag,
	RequiredDownloadsFlag,
	DisableTLSFlag,
	MetricsHTTPPortFlag,
	TheGraphPullIntervalFlag,
	TheGraphRetriesFlag,
	VerifierIntervalFlag,
	NodeClientTimeoutFlag,
	FetchBatchHeaderTimeoutFlag,
	RetrieveBlobChunksTimeoutFlag,
	GetBlobStatusTimeoutFlag,
	WriteTimeoutFlag,
	VerificationChannelCapacityFlag,
	MetricsBlacklistFlag,
	MetricsFuzzyBlacklistFlag,
}

// Flags contains the list of configuration options available to the binary.
var Flags []cli.Flag

func init() {
	Flags = append(requiredFlags, optionalFlags...)
	Flags = append(Flags, flags.RetrieverFlags(envPrefix)...)
	Flags = append(Flags, kzg.CLIFlags(envPrefix)...)
	Flags = append(Flags, common.LoggerCLIFlags(envPrefix, FlagPrefix)...)
	Flags = append(Flags, geth.EthClientFlags(envPrefix)...)
	Flags = append(Flags, indexer.CLIFlags(envPrefix)...)
	Flags = append(Flags, thegraph.CLIFlags(envPrefix)...)
}
