package churner

import (
	"time"

	"github.com/Layr-Labs/eigenda/churner/flags"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/common/logging"
	"github.com/urfave/cli"
)

type Config struct {
	EthClientConfig geth.EthClientConfig
	LoggerConfig    logging.Config
	GraphUrl        string
	MetricsConfig   MetricsConfig

	BLSOperatorStateRetrieverAddr string
	EigenDAServiceManagerAddr     string

	PerPublicKeyRateLimit time.Duration
}

func NewConfig(ctx *cli.Context) *Config {
	return &Config{
		EthClientConfig:               geth.ReadEthClientConfig(ctx),
		LoggerConfig:                  logging.ReadCLIConfig(ctx, flags.FlagPrefix),
		GraphUrl:                      ctx.GlobalString(flags.GraphUrlFlag.Name),
		BLSOperatorStateRetrieverAddr: ctx.GlobalString(flags.BlsOperatorStateRetrieverFlag.Name),
		EigenDAServiceManagerAddr:     ctx.GlobalString(flags.EigenDAServiceManagerFlag.Name),
		PerPublicKeyRateLimit:         ctx.GlobalDuration(flags.PerPublicKeyRateLimit.Name),
		MetricsConfig: MetricsConfig{
			HTTPPort:      ctx.GlobalString(flags.MetricsHTTPPort.Name),
			EnableMetrics: ctx.GlobalBool(flags.EnableMetrics.Name),
		},
	}
}
