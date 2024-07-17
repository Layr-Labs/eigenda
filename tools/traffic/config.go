package traffic

import (
	"errors"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/tools/traffic/flags"
	"github.com/urfave/cli"
)

// Config configures a traffic generator.
type Config struct {
	clients.Config

	// The number of worker threads that generate write traffic.
	NumInstances uint
	// The period of the submission rate of new blobs for each worker thread.
	RequestInterval time.Duration
	// The size of each blob dispersed, in bytes.
	DataSize uint64
	// Configures logging for the traffic generator.
	LoggingConfig common.LoggerConfig
	// If true, then each blob will contain unique random data. If false, the same random data
	// will be dispersed for each blob by a particular worker thread.
	RandomizeBlobs bool
	// The amount of time to sleep after launching each worker thread.
	InstanceLaunchInterval time.Duration

	SignerPrivateKey string
	CustomQuorums    []uint8
}

func NewConfig(ctx *cli.Context) (*Config, error) {
	loggerConfig, err := common.ReadLoggerCLIConfig(ctx, flags.FlagPrefix)
	if err != nil {
		return nil, err
	}
	customQuorums := ctx.GlobalIntSlice(flags.CustomQuorumNumbersFlag.Name)
	customQuorumsUint8 := make([]uint8, len(customQuorums))
	for i, q := range customQuorums {
		if q < 0 || q > 255 {
			return nil, errors.New("invalid custom quorum number")
		}
		customQuorumsUint8[i] = uint8(q)
	}
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
		LoggingConfig:          *loggerConfig,
		RandomizeBlobs:         ctx.GlobalBool(flags.RandomizeBlobsFlag.Name),
		InstanceLaunchInterval: ctx.Duration(flags.InstanceLaunchIntervalFlag.Name),
		SignerPrivateKey:       ctx.String(flags.SignerPrivateKeyFlag.Name),
		CustomQuorums:          customQuorumsUint8,
	}, nil
}
