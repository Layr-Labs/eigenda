package flags

import (
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/core/thegraph"
	"github.com/Layr-Labs/eigenda/indexer"
	"github.com/urfave/cli"
)

const (
	FlagPrefix = "churner"
	envPrefix  = "CHURNER"
)

var (
	/* Required Flags */
	// TODO(robert): This flag is not used in the churner code; it is only used in the deployment code
	// to determine the hostname of the churner service. We should update the deployment code with a different
	// method of setting the churner hostname for nodes and then remove this flag.
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
	PerPublicKeyRateLimit = cli.DurationFlag{
		Name:     common.PrefixFlag(FlagPrefix, "per-public-key-rate-limit"),
		Usage:    "Rate limit interval for each public key",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envPrefix, "PER_PUBLIC_KEY_RATE_LIMIT"),
		Value:    24 * time.Hour,
	}
	EnableMetrics = cli.BoolFlag{
		Name:     common.PrefixFlag(FlagPrefix, "enable-metrics"),
		Usage:    "start metrics server",
		Required: true,
		EnvVar:   common.PrefixEnvVar(envPrefix, "ENABLE_METRICS"),
	}
	/* Optional Flags*/
	MetricsHTTPPort = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "metrics-http-port"),
		Usage:    "the http port which the metrics prometheus server is listening",
		Required: false,
		Value:    "9100",
		EnvVar:   common.PrefixEnvVar(envPrefix, "METRICS_HTTP_PORT"),
	}
	ChurnApprovalInterval = cli.DurationFlag{
		Name:     common.PrefixFlag(FlagPrefix, "churn-approval-interval"),
		Usage:    "If this interval is N mins, the churner will only approve a new churn request N mins after the previous approval",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envPrefix, "CHURN_APPROVAL_INTERVAL"),
		Value:    15 * time.Minute,
	}
)

var requiredFlags = []cli.Flag{
	HostnameFlag,
	GrpcPortFlag,
	BlsOperatorStateRetrieverFlag,
	EigenDAServiceManagerFlag,
	EnableMetrics,
}

var optionalFlags = []cli.Flag{
	PerPublicKeyRateLimit,
	MetricsHTTPPort,
	ChurnApprovalInterval,
}

// Flags contains the list of configuration options available to the binary.
var Flags []cli.Flag

func init() {
	Flags = append(requiredFlags, optionalFlags...)
	Flags = append(Flags, geth.EthClientFlags(envPrefix)...)
	Flags = append(Flags, common.LoggerCLIFlags(envPrefix, FlagPrefix)...)
	Flags = append(Flags, indexer.CLIFlags(envPrefix)...)
	Flags = append(Flags, thegraph.CLIFlags(envPrefix)...)
}
