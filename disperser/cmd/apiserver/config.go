package main

import (
	"fmt"

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

type DisperserVersion uint

const (
	V1 DisperserVersion = 1
	V2 DisperserVersion = 2
)

type Config struct {
	DisperserVersion  DisperserVersion
	AwsClientConfig   aws.ClientConfig
	BlobstoreConfig   blobstore.Config
	ServerConfig      disperser.ServerConfig
	LoggerConfig      common.LoggerConfig
	MetricsConfig     disperser.MetricsConfig
	RatelimiterConfig ratelimit.Config
	RateConfig        apiserver.RateConfig
	EnableRatelimiter bool
	BucketTableName   string
	BucketStoreSize   int
	EthClientConfig   geth.EthClientConfig
	MaxBlobSize       int

	BLSOperatorStateRetrieverAddr string
	EigenDAServiceManagerAddr     string
}

func NewConfig(ctx *cli.Context) (Config, error) {
	version := ctx.GlobalUint(flags.DisperserVersionFlag.Name)
	if version != uint(V1) && version != uint(V2) {
		return Config{}, fmt.Errorf("unknown disperser version %d", version)
	}

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

	config := Config{
		DisperserVersion: DisperserVersion(version),
		AwsClientConfig:  aws.ReadClientConfig(ctx, flags.FlagPrefix),
		ServerConfig: disperser.ServerConfig{
			GrpcPort:    ctx.GlobalString(flags.GrpcPortFlag.Name),
			GrpcTimeout: ctx.GlobalDuration(flags.GrpcTimeoutFlag.Name),
		},
		BlobstoreConfig: blobstore.Config{
			BucketName: ctx.GlobalString(flags.S3BucketNameFlag.Name),
			TableName:  ctx.GlobalString(flags.DynamoDBTableNameFlag.Name),
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
