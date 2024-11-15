package main

import (
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/aws"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/common/ratelimit"
	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/Layr-Labs/eigenda/disperser/apiserver"
	"github.com/Layr-Labs/eigenda/disperser/cmd/apiserver/flags"
	"github.com/Layr-Labs/eigenda/disperser/common/blobstore"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/urfave/cli"
)

type DisperserVersion uint

const (
	V1 DisperserVersion = 1
	V2 DisperserVersion = 2
)

type Config struct {
	DisperserVersion            DisperserVersion
	AwsClientConfig             aws.ClientConfig
	BlobstoreConfig             blobstore.Config
	ServerConfig                disperser.ServerConfig
	LoggerConfig                common.LoggerConfig
	MetricsConfig               disperser.MetricsConfig
	RatelimiterConfig           ratelimit.Config
	RateConfig                  apiserver.RateConfig
	EncodingConfig              kzg.KzgConfig
	EnableRatelimiter           bool
	EnablePaymentMeterer        bool
	UpdateInterval              int
	ChainReadTimeout            int
	ReservationsTableName       string
	OnDemandTableName           string
	GlobalRateTableName         string
	BucketTableName             string
	BucketStoreSize             int
	EthClientConfig             geth.EthClientConfig
	MaxBlobSize                 int
	MaxNumSymbolsPerBlob        uint
	OnchainStateRefreshInterval time.Duration

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
		RatelimiterConfig:           ratelimiterConfig,
		RateConfig:                  rateConfig,
		EncodingConfig:              kzg.ReadCLIConfig(ctx),
		EnableRatelimiter:           ctx.GlobalBool(flags.EnableRatelimiter.Name),
		EnablePaymentMeterer:        ctx.GlobalBool(flags.EnablePaymentMeterer.Name),
		ReservationsTableName:       ctx.GlobalString(flags.ReservationsTableName.Name),
		OnDemandTableName:           ctx.GlobalString(flags.OnDemandTableName.Name),
		GlobalRateTableName:         ctx.GlobalString(flags.GlobalRateTableName.Name),
		BucketTableName:             ctx.GlobalString(flags.BucketTableName.Name),
		BucketStoreSize:             ctx.GlobalInt(flags.BucketStoreSize.Name),
		UpdateInterval:              ctx.GlobalInt(flags.UpdateInterval.Name),
		ChainReadTimeout:            ctx.GlobalInt(flags.ChainReadTimeout.Name),
		EthClientConfig:             geth.ReadEthClientConfigRPCOnly(ctx),
		MaxBlobSize:                 ctx.GlobalInt(flags.MaxBlobSize.Name),
		MaxNumSymbolsPerBlob:        ctx.GlobalUint(flags.MaxNumSymbolsPerBlob.Name),
		OnchainStateRefreshInterval: ctx.GlobalDuration(flags.OnchainStateRefreshInterval.Name),

		BLSOperatorStateRetrieverAddr: ctx.GlobalString(flags.BlsOperatorStateRetrieverFlag.Name),
		EigenDAServiceManagerAddr:     ctx.GlobalString(flags.EigenDAServiceManagerFlag.Name),
	}
	return config, nil
}
