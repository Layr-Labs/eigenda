package main

import (
	"errors"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/aws"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/common/ratelimit"
	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/Layr-Labs/eigenda/disperser/apiserver"
	"github.com/Layr-Labs/eigenda/disperser/cmd/apiserver/flags"
	"github.com/Layr-Labs/eigenda/disperser/common/blobstore"
	"github.com/urfave/cli"
)

type Config struct {
	AwsClientConfig   aws.ClientConfig
	BlobstoreConfig   blobstore.Config
	ServerConfig      disperser.ServerConfig
	LoggerConfig      common.LoggerConfig
	MetricsConfig     disperser.MetricsConfig
	RatelimiterConfig ratelimit.Config
	RateConfig        apiserver.RateConfig
	EnableRatelimiter bool
	BucketTableName   string
	ShadowTableName   string
	ShadowSampleRate  float64
	BucketStoreSize   int
	EthClientConfig   geth.EthClientConfig
	MaxBlobSize       int

	BLSOperatorStateRetrieverAddr string
	EigenDAServiceManagerAddr     string
}

func NewConfig(ctx *cli.Context) (Config, error) {

	ratelimiterConfig, err := ratelimit.ReadCLIConfig(ctx, flags.FlagPrefix)
	if err != nil {
		return Config{}, err
	}

	rateConfig, err := apiserver.ReadCLIConfig(ctx)
	if err != nil {
		return Config{}, err
	}

	loggerConfig, err := common.ReadLoggerCLIConfig(ctx, flags.FlagPrefix)
	if err != nil {
		return Config{}, err
	}

	if ctx.GlobalFloat64(flags.ShadowSampleRateFlag.Name) < 0 || ctx.GlobalFloat64(flags.ShadowSampleRateFlag.Name) > 1 {
		return Config{}, errors.New("shadow-sample-rate must be between 0.0 and 1.0")
	}

	config := Config{
		AwsClientConfig: aws.ReadClientConfig(ctx, flags.FlagPrefix),
		ServerConfig: disperser.ServerConfig{
			GrpcPort:    ctx.GlobalString(flags.GrpcPortFlag.Name),
			GrpcTimeout: ctx.GlobalDuration(flags.GrpcTimeoutFlag.Name),
		},
		BlobstoreConfig: blobstore.Config{
			BucketName:       ctx.GlobalString(flags.S3BucketNameFlag.Name),
			TableName:        ctx.GlobalString(flags.DynamoDBTableNameFlag.Name),
			ShadowTableName:  ctx.GlobalString(flags.ShadowTableNameFlag.Name),
			ShadowSampleRate: ctx.GlobalFloat64(flags.ShadowSampleRateFlag.Name),
		},
		LoggerConfig: *loggerConfig,
		MetricsConfig: disperser.MetricsConfig{
			HTTPPort:      ctx.GlobalString(flags.MetricsHTTPPort.Name),
			EnableMetrics: ctx.GlobalBool(flags.EnableMetrics.Name),
		},
		RatelimiterConfig: ratelimiterConfig,
		RateConfig:        rateConfig,
		EnableRatelimiter: ctx.GlobalBool(flags.EnableRatelimiter.Name),
		BucketTableName:   ctx.GlobalString(flags.BucketTableName.Name),
		BucketStoreSize:   ctx.GlobalInt(flags.BucketStoreSize.Name),
		EthClientConfig:   geth.ReadEthClientConfigRPCOnly(ctx),
		MaxBlobSize:       ctx.GlobalInt(flags.MaxBlobSize.Name),

		BLSOperatorStateRetrieverAddr: ctx.GlobalString(flags.BlsOperatorStateRetrieverFlag.Name),
		EigenDAServiceManagerAddr:     ctx.GlobalString(flags.EigenDAServiceManagerFlag.Name),
	}
	return config, nil
}
