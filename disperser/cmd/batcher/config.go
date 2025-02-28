package main

import (
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/aws"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/core/thegraph"
	"github.com/Layr-Labs/eigenda/disperser/batcher"
	"github.com/Layr-Labs/eigenda/disperser/cmd/batcher/flags"
	"github.com/Layr-Labs/eigenda/disperser/common/blobstore"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/Layr-Labs/eigenda/indexer"
	"github.com/urfave/cli"
)

type Config struct {
	BatcherConfig    batcher.Config
	TimeoutConfig    batcher.TimeoutConfig
	BlobstoreConfig  blobstore.Config
	EthClientConfig  geth.EthClientConfig
	AwsClientConfig  aws.ClientConfig
	EncoderConfig    kzg.KzgConfig
	LoggerConfig     common.LoggerConfig
	MetricsConfig    batcher.MetricsConfig
	IndexerConfig    indexer.Config
	KMSKeyConfig     common.KMSKeyConfig
	ChainStateConfig thegraph.Config
	UseGraph         bool

	IndexerDataDir string

	BLSOperatorStateRetrieverAddr string
	EigenDAServiceManagerAddr     string

	EnableGnarkBundleEncoding bool
}

func NewConfig(ctx *cli.Context) (Config, error) {
	loggerConfig, err := common.ReadLoggerCLIConfig(ctx, flags.FlagPrefix)
	if err != nil {
		return Config{}, err
	}
	ethClientConfig := geth.ReadEthClientConfig(ctx)
	kmsConfig := common.ReadKMSKeyConfig(ctx, flags.FlagPrefix)
	if !kmsConfig.Disable {
		ethClientConfig = geth.ReadEthClientConfigRPCOnly(ctx)
	}
	config := Config{
		BlobstoreConfig: blobstore.Config{
			BucketName: ctx.GlobalString(flags.S3BucketNameFlag.Name),
			TableName:  ctx.GlobalString(flags.DynamoDBTableNameFlag.Name),
		},
		EthClientConfig: ethClientConfig,
		AwsClientConfig: aws.ReadClientConfig(ctx, flags.FlagPrefix),
		EncoderConfig:   kzg.ReadCLIConfig(ctx),
		LoggerConfig:    *loggerConfig,
		BatcherConfig: batcher.Config{
			PullInterval:             ctx.GlobalDuration(flags.PullIntervalFlag.Name),
			FinalizerInterval:        ctx.GlobalDuration(flags.FinalizerIntervalFlag.Name),
			FinalizerPoolSize:        ctx.GlobalInt(flags.FinalizerPoolSizeFlag.Name),
			EncoderSocket:            ctx.GlobalString(flags.EncoderSocket.Name),
			NumConnections:           ctx.GlobalInt(flags.NumConnectionsFlag.Name),
			EncodingRequestQueueSize: ctx.GlobalInt(flags.EncodingRequestQueueSizeFlag.Name),
			BatchSizeMBLimit:         ctx.GlobalUint(flags.BatchSizeLimitFlag.Name),
			SRSOrder:                 ctx.GlobalInt(flags.SRSOrderFlag.Name),
			MaxNumRetriesPerBlob:     ctx.GlobalUint(flags.MaxNumRetriesPerBlobFlag.Name),
			TargetNumChunks:          ctx.GlobalUint(flags.TargetNumChunksFlag.Name),
			MaxBlobsToFetchFromStore: ctx.GlobalInt(flags.MaxBlobsToFetchFromStoreFlag.Name),
			FinalizationBlockDelay:   ctx.GlobalUint(flags.FinalizationBlockDelayFlag.Name),
		},
		TimeoutConfig: batcher.TimeoutConfig{
			EncodingTimeout:     ctx.GlobalDuration(flags.EncodingTimeoutFlag.Name),
			AttestationTimeout:  ctx.GlobalDuration(flags.AttestationTimeoutFlag.Name),
			ChainReadTimeout:    ctx.GlobalDuration(flags.ChainReadTimeoutFlag.Name),
			ChainWriteTimeout:   ctx.GlobalDuration(flags.ChainWriteTimeoutFlag.Name),
			ChainStateTimeout:   ctx.GlobalDuration(flags.ChainStateTimeoutFlag.Name),
			TxnBroadcastTimeout: ctx.GlobalDuration(flags.TransactionBroadcastTimeoutFlag.Name),
		},
		MetricsConfig: batcher.MetricsConfig{
			HTTPPort:      ctx.GlobalString(flags.MetricsHTTPPort.Name),
			EnableMetrics: ctx.GlobalBool(flags.EnableMetrics.Name),
		},
		ChainStateConfig:              thegraph.ReadCLIConfig(ctx),
		UseGraph:                      ctx.Bool(flags.UseGraphFlag.Name),
		BLSOperatorStateRetrieverAddr: ctx.GlobalString(flags.BlsOperatorStateRetrieverFlag.Name),
		EigenDAServiceManagerAddr:     ctx.GlobalString(flags.EigenDAServiceManagerFlag.Name),
		IndexerDataDir:                ctx.GlobalString(flags.IndexerDataDirFlag.Name),
		IndexerConfig:                 indexer.ReadIndexerConfig(ctx),
		KMSKeyConfig:                  kmsConfig,
		EnableGnarkBundleEncoding:     ctx.Bool(flags.EnableGnarkBundleEncodingFlag.Name),
	}
	return config, nil
}
