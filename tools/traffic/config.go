package traffic

import (
	"errors"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/tools/traffic/flags"
	"github.com/urfave/cli"
)

type Config struct {
	clients.Config

	NumInstances           uint
	RequestInterval        time.Duration
	DataSize               uint64
	LoggingConfig          common.LoggerConfig
	RandomizeBlobs         bool
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
		Config: clients.Config{
			Hostname:          ctx.GlobalString(flags.HostnameFlag.Name),
			Port:              ctx.GlobalString(flags.GrpcPortFlag.Name),
			Timeout:           ctx.Duration(flags.TimeoutFlag.Name),
			UseSecureGrpcFlag: ctx.GlobalBool(flags.UseSecureGrpcFlag.Name),
		},
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
