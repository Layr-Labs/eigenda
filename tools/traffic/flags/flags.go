package flags

import (
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
		EnvVar:   common.PrefixEnvVar(envPrefix, "HOSTNAME"),
	}
	GrpcPortFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "disperser-port"),
		Usage:    "Port at which a disperser listens for grpc calls.",
		Required: true,
		EnvVar:   common.PrefixEnvVar(envPrefix, "GRPC_PORT"),
	}
	TimeoutFlag = cli.DurationFlag{
		Name:     common.PrefixFlag(FlagPrefix, "timeout"),
		Usage:    "Amount of time to wait for grpc.",
		Required: false,
		Value:    10 * time.Second,
		EnvVar:   common.PrefixEnvVar(envPrefix, "TIMEOUT"),
	}
	UseSecureGrpcFlag = cli.BoolFlag{
		Name:     common.PrefixFlag(FlagPrefix, "use-secure-grpc"),
		Usage:    "Whether to use secure grpc.",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envPrefix, "USE_SECURE_GRPC"),
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
	EthClientHostnameFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "eth-client-hostname"),
		Usage:    "Hostname at which the Ethereum client is available.",
		Required: true,
		EnvVar:   common.PrefixEnvVar(envPrefix, "ETH_CLIENT_HOSTNAME"),
	}
	EthClientPortFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "eth-client-port"),
		Usage:    "Port at which the Ethereum client is available.",
		Required: true,
		EnvVar:   common.PrefixEnvVar(envPrefix, "ETH_CLIENT_PORT"),
	}
	BLSOperatorStateRetrieverFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "bls-operator-state-retriever"),
		Usage:    "Hex address of the BLS operator state retriever contract.",
		Required: true,
		EnvVar:   common.PrefixEnvVar(envPrefix, "BLS_OPERATOR_STATE_RETRIEVER"),
	}
	EigenDAServiceManagerFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "eigenda-service-manager"),
		Usage:    "Hex address of the EigenDA service manager contract.",
		Required: true,
		EnvVar:   common.PrefixEnvVar(envPrefix, "EIGENDA_SERVICE_MANAGER"),
	}
	EthClientRetriesFlag = cli.UintFlag{
		Name:     common.PrefixFlag(FlagPrefix, "eth-client-retries"),
		Usage:    "Number of times to retry an Ethereum client request.",
		Required: false,
		Value:    2,
		EnvVar:   common.PrefixEnvVar(envPrefix, "ETH_CLIENT_RETRIES"),
	}

	/* Common Configuration. */

	InstanceLaunchIntervalFlag = cli.DurationFlag{
		Name:     common.PrefixFlag(FlagPrefix, "instance-launch-interva"),
		Usage:    "Duration between generator instance launches.",
		Required: false,
		Value:    1 * time.Second,
		EnvVar:   common.PrefixEnvVar(envPrefix, "INSTANCE_LAUNCH_INTERVAL"),
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
	EncoderG1PathFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "encoder-g1-path"),
		Usage:    "Path to the encoder G1 binary.",
		Required: true,
		EnvVar:   common.PrefixEnvVar(envPrefix, "ENCODER_G1_PATH"),
	}
	EncoderG2PathFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "encoder-g2-path"),
		Usage:    "Path to the encoder G2 binary.",
		Required: true,
		EnvVar:   common.PrefixEnvVar(envPrefix, "ENCODER_G2_PATH"),
	}
	EncoderCacheDirFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "encoder-cache-dir"),
		Usage:    "Path to the encoder cache directory.",
		Required: true,
		EnvVar:   common.PrefixEnvVar(envPrefix, "ENCODER_CACHE_DIR"),
	}
	EncoderSRSOrderFlag = cli.UintFlag{
		Name:     common.PrefixFlag(FlagPrefix, "encoder-srs-order"),
		Usage:    "The SRS order to use for the encoder.",
		Required: false,
		Value:    3000,
		EnvVar:   common.PrefixEnvVar(envPrefix, "ENCODER_SRS_ORDER"),
	}
	EncoderSRSNumberToLoadFlag = cli.UintFlag{
		Name:     common.PrefixFlag(FlagPrefix, "encoder-srs-number-to-load"),
		Usage:    "The SRS number to load for the encoder.",
		Required: false,
		Value:    3000,
		EnvVar:   common.PrefixEnvVar(envPrefix, "ENCODER_SRS_NUMBER_TO_LOAD"),
	}
	EncoderNumWorkersFlag = cli.UintFlag{
		Name:     common.PrefixFlag(FlagPrefix, "encoder-num-workers"),
		Usage:    "The number of worker threads to use for the encoder.",
		Required: false,
		Value:    4,
	}

	/* Configuration for the blob validator. */

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
		Value:    time.Second / 3,
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
)

var requiredFlags = []cli.Flag{
	HostnameFlag,
	GrpcPortFlag,
	EthClientHostnameFlag,
	EthClientPortFlag,
	BLSOperatorStateRetrieverFlag,
	EigenDAServiceManagerFlag,
	TheGraphUrlFlag,
	EncoderG1PathFlag,
	EncoderG2PathFlag,
	EncoderCacheDirFlag,
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
	EthClientRetriesFlag,
	TheGraphPullIntervalFlag,
	TheGraphRetriesFlag,
	EncoderSRSOrderFlag,
	EncoderSRSNumberToLoadFlag,
	EncoderNumWorkersFlag,
}

// Flags contains the list of configuration options available to the binary.
var Flags []cli.Flag

func init() {
	Flags = append(requiredFlags, optionalFlags...)
	Flags = append(Flags, common.LoggerCLIFlags(envPrefix, FlagPrefix)...)
}
