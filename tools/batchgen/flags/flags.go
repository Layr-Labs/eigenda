package flags

import (
	"github.com/Layr-Labs/eigenda/common"
	"github.com/urfave/cli"
)

const (
	FlagPrefix = "batchgen"
	envPrefix  = "BATCHGEN"
)

var (
	/* Required Flags */

	NumThreadsFlag = cli.UintFlag{
		Name:     "host-threads",
		Usage:    "Number of host threads in parallel",
		Required: false,
		Value:    1,
		EnvVar:   "HOST_THREADS",
	}
	ScalingFactorFlag = cli.UintFlag{
		Name:     "scaling-factor",
		Usage:    "Scaling factor applied to default batch size",
		Required: false,
		Value:    10,
		EnvVar:   "SCALING_FACTOR",
	}
)

var requiredFlags = []cli.Flag{
	NumThreadsFlag,
	ScalingFactorFlag,
}

var optionalFlags = []cli.Flag{}

// Flags contains the list of configuration options available to the binary.
var Flags []cli.Flag

func init() {
	Flags = append(requiredFlags, optionalFlags...)
	Flags = append(Flags, common.LoggerCLIFlags(envPrefix, FlagPrefix)...)
}
