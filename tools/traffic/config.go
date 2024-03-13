package traffic

import (
	"time"

	"github.com/Layr-Labs/eigenda/clients"
	"github.com/Layr-Labs/eigenda/common/logging"
	"github.com/Layr-Labs/eigenda/tools/traffic/flags"
	"github.com/urfave/cli"
)

type Config struct {
	clients.Config

	NumInstances           uint
	RequestInterval        time.Duration
	DataSize               uint64
	ConfirmationThreshold  uint8
	AdversarialThreshold   uint8
	LoggingConfig          logging.Config
	RandomizeBlobs         bool
	InstanceLaunchInterval time.Duration
}

func NewConfig(ctx *cli.Context) *Config {
	return &Config{
		Config: *clients.NewConfig(
			ctx.GlobalString(flags.HostnameFlag.Name),
			ctx.GlobalString(flags.GrpcPortFlag.Name),
			ctx.Duration(flags.TimeoutFlag.Name),
			ctx.GlobalBool(flags.UseSecureGrpcFlag.Name),
		),
		NumInstances:           ctx.GlobalUint(flags.NumInstancesFlag.Name),
		RequestInterval:        ctx.Duration(flags.RequestIntervalFlag.Name),
		DataSize:               ctx.GlobalUint64(flags.DataSizeFlag.Name),
		ConfirmationThreshold:  uint8(ctx.GlobalUint(flags.QuorumThresholdFlag.Name)),
		AdversarialThreshold:   uint8(ctx.GlobalUint(flags.AdversarialThresholdFlag.Name)),
		LoggingConfig:          logging.ReadCLIConfig(ctx, flags.FlagPrefix),
		RandomizeBlobs:         ctx.GlobalBool(flags.RandomizeBlobsFlag.Name),
		InstanceLaunchInterval: ctx.Duration(flags.InstanceLaunchIntervalFlag.Name),
	}
}
