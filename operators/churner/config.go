package churner

import (
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/operators/churner/flags"
	"github.com/urfave/cli"
)

type Config struct {
	EthClientConfig geth.EthClientConfig
	LoggerConfig    common.LoggerConfig
	GraphUrl        string
	MetricsConfig   MetricsConfig

	BLSOperatorStateRetrieverAddr string
	EigenDAServiceManagerAddr     string

	PerPublicKeyRateLimit time.Duration
}

func NewConfig(ctx *cli.Context) (*Config, error) {
	loggerConfig, err := common.ReadLoggerCLIConfig(ctx, flags.FlagPrefix)
	if err != nil {
		return nil, err
	}
	return &Config{
		EthClientConfig:               geth.ReadEthClientConfig(ctx),
		LoggerConfig:                  *loggerConfig,
		GraphUrl:                      ctx.GlobalString(flags.GraphUrlFlag.Name),
		BLSOperatorStateRetrieverAddr: ctx.GlobalString(flags.BlsOperatorStateRetrieverFlag.Name),
		EigenDAServiceManagerAddr:     ctx.GlobalString(flags.EigenDAServiceManagerFlag.Name),
		PerPublicKeyRateLimit:         ctx.GlobalDuration(flags.PerPublicKeyRateLimit.Name),
		MetricsConfig: MetricsConfig{
			HTTPPort:      ctx.GlobalString(flags.MetricsHTTPPort.Name),
			EnableMetrics: ctx.GlobalBool(flags.EnableMetrics.Name),
		},
	}, nil
}
