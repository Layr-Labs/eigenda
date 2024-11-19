package main

import (
	"fmt"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/aws"
	"github.com/Layr-Labs/eigenda/disperser/cmd/encoder/flags"
	"github.com/Layr-Labs/eigenda/disperser/common/blobstore"
	"github.com/Layr-Labs/eigenda/disperser/encoder"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/Layr-Labs/eigenda/relay/chunkstore"
	"github.com/urfave/cli"
)

type EncoderVersion uint

const (
	V1 EncoderVersion = 1
	V2 EncoderVersion = 2
)

type Config struct {
	EncoderVersion   EncoderVersion
	AwsClientConfig  aws.ClientConfig
	BlobStoreConfig  blobstore.Config
	ChunkStoreConfig chunkstore.Config
	EncoderConfig    kzg.KzgConfig
	LoggerConfig     common.LoggerConfig
	ServerConfig     *encoder.ServerConfig
	MetricsConfig    encoder.MetrisConfig
}

func NewConfig(ctx *cli.Context) (Config, error) {
	version := ctx.GlobalUint(flags.EncoderVersionFlag.Name)
	if version != uint(V1) && version != uint(V2) {
		return Config{}, fmt.Errorf("unknown encoder version %d", version)
	}

	loggerConfig, err := common.ReadLoggerCLIConfig(ctx, flags.FlagPrefix)
	if err != nil {
		return Config{}, err
	}
	config := Config{
		EncoderVersion:  EncoderVersion(version),
		AwsClientConfig: aws.ReadClientConfig(ctx, flags.FlagPrefix),
		BlobStoreConfig: blobstore.Config{
			BucketName: ctx.GlobalString(flags.S3BucketNameFlag.Name),
		},
		ChunkStoreConfig: chunkstore.Config{
			BucketName: ctx.GlobalString(flags.S3BucketNameFlag.Name),
		},
		EncoderConfig: kzg.ReadCLIConfig(ctx),
		LoggerConfig:  *loggerConfig,
		ServerConfig: &encoder.ServerConfig{
			GrpcPort:                 ctx.GlobalString(flags.GrpcPortFlag.Name),
			MaxConcurrentRequests:    ctx.GlobalInt(flags.MaxConcurrentRequestsFlag.Name),
			RequestPoolSize:          ctx.GlobalInt(flags.RequestPoolSizeFlag.Name),
			EnableGnarkChunkEncoding: ctx.Bool(flags.EnableGnarkChunkEncodingFlag.Name),
			PreventReencoding:        ctx.Bool(flags.PreventReencodingFlag.Name),
		},
		MetricsConfig: encoder.MetrisConfig{
			HTTPPort:      ctx.GlobalString(flags.MetricsHTTPPort.Name),
			EnableMetrics: ctx.GlobalBool(flags.EnableMetrics.Name),
		},
	}
	return config, nil
}
