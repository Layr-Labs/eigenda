package main

import (
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/disperser/cmd/encoder/flags"
	"github.com/Layr-Labs/eigenda/disperser/encoder"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/urfave/cli"
)

type Config struct {
	EncoderConfig kzg.KzgConfig
	LoggerConfig  common.LoggerConfig
	ServerConfig  *encoder.ServerConfig
	MetricsConfig encoder.MetrisConfig
}

func NewConfig(ctx *cli.Context) (Config, error) {
	loggerConfig, err := common.ReadLoggerCLIConfig(ctx, flags.FlagPrefix)
	if err != nil {
		return Config{}, err
	}
	config := Config{
		EncoderConfig: kzg.ReadCLIConfig(ctx),
		LoggerConfig:  *loggerConfig,
		ServerConfig: &encoder.ServerConfig{
			GrpcPort:              ctx.GlobalString(flags.GrpcPortFlag.Name),
			MaxConcurrentRequests: ctx.GlobalInt(flags.MaxConcurrentRequestsFlag.Name),
			RequestPoolSize:       ctx.GlobalInt(flags.RequestPoolSizeFlag.Name),
		},
		MetricsConfig: encoder.MetrisConfig{
			HTTPPort:      ctx.GlobalString(flags.MetricsHTTPPort.Name),
			EnableMetrics: ctx.GlobalBool(flags.EnableMetrics.Name),
		},
	}
	return config, nil
}
