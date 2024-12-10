package flags

import (
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/aws"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/core/thegraph"
	"github.com/urfave/cli"
)

const (
	FlagPrefix   = "data-access-api"
	envVarPrefix = "DATA_ACCESS_API"
)

var (
	DynamoTableNameFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "dynamo-table-name"),
		Usage:    "Name of the dynamo table to store blob metadata",
		Required: true,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "DYNAMO_TABLE_NAME"),
	}
	S3BucketNameFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "s3-bucket-name"),
		Usage:    "Name of the bucket to store blobs",
		Required: true,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "S3_BUCKET_NAME"),
	}
	SocketAddrFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "socket-addr"),
		Usage:    "the socket address of the data access api",
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "SOCKET_ADDR"),
		Required: true,
	}
	PrometheusServerURLFlag = cli.StringFlag{
		Name: common.PrefixFlag(FlagPrefix, "prometheus-server-url"),
		//We need the prometheus server url to be able to query the metrics
		Usage:    "the url of the prometheus server",
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "PROMETHEUS_SERVER_URL"),
		Required: true,
	}
	PrometheusServerUsernameFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "prometheus-server-usename"),
		Usage:    "the username for basic auth of the prometheus server",
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "PROMETHEUS_SERVER_USERNAME"),
		Required: true,
	}
	PrometheusServerSecretFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "prometheus-server-secret"),
		Usage:    "the secret for basic auth of the prometheus server",
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "PROMETHEUS_SERVER_SECRET"),
		Required: true,
	}
	PrometheusMetricsClusterLabelFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "prometheus-metrics-cluster-label"),
		Usage:    "the cluster label for metrics in the prometheus",
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "PROMETHEUS_METRICS_CLUSTER_LABEL"),
		Required: true,
	}
	SubgraphApiBatchMetadataAddrFlag = cli.StringFlag{
		Name: common.PrefixFlag(FlagPrefix, "sub-batch-metadata-socket-addr"),
		//We need the socket address of the subgraph batch metadata api to pull the subgraph data from.
		Usage:    "the socket address of the subgraph batch metadata api",
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "SUBGRAPH_BATCH_METADATA_API_SOCKET_ADDR"),
		Required: true,
	}
	SubgraphApiOperatorStateAddrFlag = cli.StringFlag{
		Name: common.PrefixFlag(FlagPrefix, "sub-op-state-socket-addr"),
		//We need the socket address of the subgraph operator state api to pull the subgraph data from.
		Usage:    "the socket address of the subgraph operator state api",
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "SUBGRAPH_OPERATOR_STATE_API_SOCKET_ADDR"),
		Required: true,
	}
	BlsOperatorStateRetrieverFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "bls-operator-state-retriever"),
		Usage:    "Address of the BLS Operator State Retriever",
		Required: true,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "BLS_OPERATOR_STATE_RETRIVER"),
	}
	EigenDAServiceManagerFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "eigenda-service-manager"),
		Usage:    "Address of the EigenDA Service Manager",
		Required: true,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "EIGENDA_SERVICE_MANAGER"),
	}
	ServerModeFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "server-mode"),
		Usage:    "Set the mode of the server (debug, release or test)",
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "SERVER_MODE"),
		Required: false,
		Value:    "debug",
	}
	AllowOriginsFlag = cli.StringSliceFlag{
		Name:     common.PrefixFlag(FlagPrefix, "allow-origins"),
		Usage:    "Set the allowed origins for CORS requests",
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "ALLOW_ORIGINS"),
		Required: true,
	}
	EnableMetricsFlag = cli.BoolFlag{
		Name:     common.PrefixFlag(FlagPrefix, "enable-metrics"),
		Usage:    "start metrics server",
		Required: true,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "ENABLE_METRICS"),
	}
	// EigenDA Disperser and Churner Hostnames to check Server Availability
	// ex:
	// disperser-goerli.eigenda.eigenops.xyz,
	// churner-goerli.eigenda.eigenops.xyz
	DisperserHostnameFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "eigenda-disperser-hostname"),
		Usage:    "HostName of EigenDA Disperser",
		Required: true,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "EIGENDA_DISPERSER_HOSTNAME"),
	}
	ChurnerHostnameFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "eigenda-churner-hostname"),
		Usage:    "HostName of EigenDA Churner",
		Required: true,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "EIGENDA_CHURNER_HOSTNAME"),
	}
	BatcherHealthEndptFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "eigenda-batcher-health-endpoint"),
		Usage:    "Endpt of EigenDA Batcher Health Sidecar",
		Required: true,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "EIGENDA_BATCHER_HEALTH_ENDPOINT"),
	}
	/* Optional Flags*/
	MetricsHTTPPort = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "metrics-http-port"),
		Usage:    "the http port which the metrics prometheus server is listening",
		Required: false,
		Value:    "9100",
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "METRICS_HTTP_PORT"),
	}
	DataApiServerVersionFlag = cli.UintFlag{
		Name:     common.PrefixFlag(FlagPrefix, "dataapi-version"),
		Usage:    "DataApi server version. Options are 1 and 2.",
		Required: false,
		Value:    1,
		EnvVar:   common.PrefixEnvVar(envVarPrefix, "DATA_API_VERSION"),
	}
)

var requiredFlags = []cli.Flag{
	DynamoTableNameFlag,
	SocketAddrFlag,
	S3BucketNameFlag,
	SubgraphApiBatchMetadataAddrFlag,
	SubgraphApiOperatorStateAddrFlag,
	BlsOperatorStateRetrieverFlag,
	EigenDAServiceManagerFlag,
	PrometheusServerURLFlag,
	PrometheusServerUsernameFlag,
	PrometheusServerSecretFlag,
	PrometheusMetricsClusterLabelFlag,
	AllowOriginsFlag,
	EnableMetricsFlag,
	DisperserHostnameFlag,
	ChurnerHostnameFlag,
	BatcherHealthEndptFlag,
}

var optionalFlags = []cli.Flag{
	ServerModeFlag,
	MetricsHTTPPort,
	DataApiServerVersionFlag,
}

// Flags contains the list of configuration options available to the binary.
var Flags []cli.Flag

func init() {
	Flags = append(requiredFlags, optionalFlags...)
	Flags = append(Flags, common.LoggerCLIFlags(envVarPrefix, FlagPrefix)...)
	Flags = append(Flags, geth.EthClientFlags(envVarPrefix)...)
	Flags = append(Flags, aws.ClientFlags(envVarPrefix, FlagPrefix)...)
	Flags = append(Flags, thegraph.CLIFlags(envVarPrefix)...)
}
