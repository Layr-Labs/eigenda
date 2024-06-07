package flags

import (
	"github.com/Layr-Labs/eigenda/common"
	"github.com/urfave/cli"
)

const (
	FlagPrefix = ""
	envPrefix  = "BATCHGEN"
)

var (
	/* Required Flags */

	NumThreadsFlag = cli.UintFlag{
		Name:     common.PrefixFlag(FlagPrefix, "threads"),
		Usage:    "Number of host threads in parallel",
		Required: false,
		Value:    1,
		EnvVar:   common.PrefixEnvVar(envPrefix, "HOST_THREADS"),
	}
	ScalingFactorFlag = cli.UintFlag{
		Name:     common.PrefixFlag(FlagPrefix, "scaling-factor"),
		Usage:    "Scaling factor applied to default batch size",
		Required: false,
		Value:    100,
		EnvVar:   common.PrefixEnvVar(envPrefix, "SCALING_FACTOR"),
	}
	HostsFlag = cli.StringSliceFlag{
		Name:     common.PrefixFlag(FlagPrefix, "host"),
		Usage:    "host:port",
		Required: true,
	}
)

var requiredFlags = []cli.Flag{
	HostsFlag,
}

var optionalFlags = []cli.Flag{
	NumThreadsFlag,
	ScalingFactorFlag,
}

// Flags contains the list of configuration options available to the binary.
var Flags []cli.Flag

func init() {
	Flags = append(requiredFlags, optionalFlags...)
	Flags = append(Flags, common.LoggerCLIFlags(envPrefix, FlagPrefix)...)
}
