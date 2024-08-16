package opscan

import (
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/tools/opscan/flags"
	"github.com/urfave/cli"
)

type Config struct {
	LoggerConfig     common.LoggerConfig
	MaxConnections   int
	OperatorId       string
	SubgraphEndpoint string
	Timeout          time.Duration
}

func ReadConfig(ctx *cli.Context) *Config {
	return &Config{
		Timeout:          ctx.Duration(flags.TimeoutFlag.Name),
		MaxConnections:   ctx.Int(flags.MaxConnectionsFlag.Name),
		OperatorId:       ctx.String(flags.OperatorIdFlag.Name),
		SubgraphEndpoint: ctx.String(flags.SubgraphEndpointFlag.Name),
	}
}

func NewConfig(ctx *cli.Context) (*Config, error) {
	loggerConfig, err := common.ReadLoggerCLIConfig(ctx, flags.FlagPrefix)
	if err != nil {
		return nil, err
	}

	config := ReadConfig(ctx)
	config.LoggerConfig = *loggerConfig

	return config, nil
}
