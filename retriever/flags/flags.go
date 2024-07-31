package flags

import (
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/core/thegraph"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/Layr-Labs/eigenda/indexer"
	"github.com/urfave/cli"
)

const (
	FlagPrefix = "retriever"
	envPrefix  = "RETRIEVER"
)

var (
	/* Required Flags */
	HostnameFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "hostname"),
		Usage:    "Hostname at which retriever service is available",
		Required: true,
		EnvVar:   common.PrefixEnvVar(envPrefix, "HOSTNAME"),
	}
	GrpcPortFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "grpc-port"),
		Usage:    "Port at which a retriever listens for grpc calls",
		Required: true,
		EnvVar:   common.PrefixEnvVar(envPrefix, "GRPC_PORT"),
	}
	TimeoutFlag = cli.DurationFlag{
		Name:     common.PrefixFlag(FlagPrefix, "timeout"),
		Usage:    "Amount of time to wait for GPRC",
		Required: true,
		EnvVar:   common.PrefixEnvVar(envPrefix, "TIMEOUT"),
	}
	BlsOperatorStateRetrieverFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "bls-operator-state-retriever"),
		Usage:    "Address of the BLS Operator State Retriever",
		Required: true,
		EnvVar:   common.PrefixEnvVar(envPrefix, "BLS_OPERATOR_STATE_RETRIVER"),
	}
	EigenDAServiceManagerFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "eigenda-service-manager"),
		Usage:    "Address of the EigenDA Service Manager",
		Required: true,
		EnvVar:   common.PrefixEnvVar(envPrefix, "EIGENDA_SERVICE_MANAGER"),
	}

	/* Optional Flags*/
	NumConnectionsFlag = cli.IntFlag{
		Name:     common.PrefixFlag(FlagPrefix, "num-connections"),
		Usage:    "maximum number of connections to DA nodes (defaults to 20)",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envPrefix, "NUM_CONNECTIONS"),
		Value:    20,
	}
	IndexerDataDirFlag = cli.StringFlag{
		Name:   common.PrefixFlag(FlagPrefix, "indexer-data-dir"),
		Usage:  "the data directory for the indexer",
		EnvVar: common.PrefixEnvVar(envPrefix, "DATA_DIR"),
		Value:  "./data/retriever",
	}
	MetricsHTTPPortFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "metrics-http-port"),
		Usage:    "the http port which the metrics prometheus server is listening",
		Required: false,
		Value:    "9100",
		EnvVar:   common.PrefixEnvVar(envPrefix, "METRICS_HTTP_PORT"),
	}
	UseGraphFlag = cli.BoolFlag{
		Name:     common.PrefixFlag(FlagPrefix, "use-graph"),
		Usage:    "Whether to use the graph node",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envPrefix, "USE_GRAPH"),
	}
)

func RetrieverFlags(envPrefix string) []cli.Flag {
	return []cli.Flag{
		HostnameFlag,
		GrpcPortFlag,
		TimeoutFlag,
		BlsOperatorStateRetrieverFlag,
		EigenDAServiceManagerFlag,
		NumConnectionsFlag,
		IndexerDataDirFlag,
		MetricsHTTPPortFlag,
		UseGraphFlag,
	}
}

// Flags contains the list of configuration options available to the binary.
var Flags []cli.Flag

func init() {
	Flags = append(Flags, RetrieverFlags(envPrefix)...)
	Flags = append(Flags, kzg.CLIFlags(envPrefix)...)
	Flags = append(Flags, geth.EthClientFlags(envPrefix)...)
	Flags = append(Flags, common.LoggerCLIFlags(envPrefix, FlagPrefix)...)
	Flags = append(Flags, indexer.CLIFlags(envPrefix)...)
	Flags = append(Flags, thegraph.CLIFlags(envPrefix)...)
}
