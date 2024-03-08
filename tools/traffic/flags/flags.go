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
	/* Required Flags */

	HostnameFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "disperser-hostname"),
		Usage:    "Hostname at which disperser service is available",
		Required: true,
		EnvVar:   common.PrefixEnvVar(envPrefix, "HOSTNAME"),
	}
	GrpcPortFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "disperser-port"),
		Usage:    "Port at which a disperser listens for grpc calls",
		Required: true,
		EnvVar:   common.PrefixEnvVar(envPrefix, "GRPC_PORT"),
	}
	TimeoutFlag = cli.DurationFlag{
		Name:     common.PrefixFlag(FlagPrefix, "timeout"),
		Usage:    "Amount of time to wait for GPRC",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envPrefix, "TIMEOUT"),
		Value:    10 * time.Second,
	}
	NumInstancesFlag = cli.UintFlag{
		Name:     common.PrefixFlag(FlagPrefix, "num-instances"),
		Usage:    "Number of generator instances to run in parallel",
		Required: true,
		EnvVar:   common.PrefixEnvVar(envPrefix, "NUM_INSTANCES"),
	}
	RequestIntervalFlag = cli.DurationFlag{
		Name:     common.PrefixFlag(FlagPrefix, "request-interval"),
		Usage:    "Duration between requests",
		Required: true,
		EnvVar:   common.PrefixEnvVar(envPrefix, "REQUEST_INTERVAL"),
		Value:    30 * time.Second,
	}
	DataSizeFlag = cli.Uint64Flag{
		Name:     common.PrefixFlag(FlagPrefix, "data-size"),
		Usage:    "Size of the data blob",
		Required: true,
		EnvVar:   common.PrefixEnvVar(envPrefix, "DATA_SIZE"),
	}
	AdversarialThresholdFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "adv-threshold"),
		Usage:    "Adversarial threshold between 0 and 100",
		Required: true,
		EnvVar:   common.PrefixEnvVar(envPrefix, "ADV_THRESHOLD"),
	}
	QuorumThresholdFlag = cli.StringFlag{
		Name:     common.PrefixFlag(FlagPrefix, "quorum-threshold"),
		Usage:    "Quorum threshold between 0 and 100",
		Required: true,
		EnvVar:   common.PrefixEnvVar(envPrefix, "QUORUM_THRESHOLD"),
	}
	RandomizeBlobsFlag = cli.BoolFlag{
		Name:     common.PrefixFlag(FlagPrefix, "randomize-blobs"),
		Usage:    "Whether to randomzie blob data",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envPrefix, "RANDOMIZE_BLOBS"),
	}
	InstanceLaunchIntervalFlag = cli.DurationFlag{
		Name:     common.PrefixFlag(FlagPrefix, "instance-launch-interva"),
		Usage:    "Duration between generator instance launches",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envPrefix, "INSTANCE_LAUNCH_INTERVAL"),
		Value:    1 * time.Second,
	}
	UseSecureGrpcFlag = cli.BoolFlag{
		Name:     common.PrefixFlag(FlagPrefix, "use-secure-grpc"),
		Usage:    "Whether to use secure grpc",
		Required: false,
		EnvVar:   common.PrefixEnvVar(envPrefix, "USE_SECURE_GRPC"),
	}
)

var requiredFlags = []cli.Flag{
	HostnameFlag,
	GrpcPortFlag,
	NumInstancesFlag,
	RequestIntervalFlag,
	DataSizeFlag,
	AdversarialThresholdFlag,
	QuorumThresholdFlag,
}

var optionalFlags = []cli.Flag{
	TimeoutFlag,
	RandomizeBlobsFlag,
	InstanceLaunchIntervalFlag,
	UseSecureGrpcFlag,
}

// Flags contains the list of configuration options available to the binary.
var Flags []cli.Flag

func init() {
	Flags = append(requiredFlags, optionalFlags...)
	Flags = append(Flags, common.LoggerCLIFlags(envPrefix, FlagPrefix)...)
}
