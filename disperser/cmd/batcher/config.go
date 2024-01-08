package main

import (
	"github.com/Layr-Labs/eigenda/common/aws"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/common/logging"
	"github.com/Layr-Labs/eigenda/core/encoding"
	"github.com/Layr-Labs/eigenda/disperser/batcher"
	"github.com/Layr-Labs/eigenda/disperser/cmd/batcher/flags"
	"github.com/Layr-Labs/eigenda/disperser/common/blobstore"
	"github.com/Layr-Labs/eigenda/indexer"
	"github.com/urfave/cli"
)

type Config struct {
	BatcherConfig   batcher.Config
	TimeoutConfig   batcher.TimeoutConfig
	BlobstoreConfig blobstore.Config
	EthClientConfig geth.EthClientConfig
	AwsClientConfig aws.ClientConfig
	EncoderConfig   encoding.EncoderConfig
	LoggerConfig    logging.Config
	MetricsConfig   batcher.MetricsConfig
	IndexerConfig   indexer.Config
	GraphUrl        string
	UseGraph        bool

	IndexerDataDir string

	BLSOperatorStateRetrieverAddr string
	EigenDAServiceManagerAddr     string
}

func NewConfig(ctx *cli.Context) Config {

	maxBlobsToFetchFromStore := ctx.GlobalInt(flags.MaxBlobsToFetchFromStoreFlag.Name)
	// Set Minimum Number if no value is set
	if maxBlobsToFetchFromStore == 0 {
		maxBlobsToFetchFromStore = 1
	}
	config := Config{
		BlobstoreConfig: blobstore.Config{
			BucketName: ctx.GlobalString(flags.S3BucketNameFlag.Name),
			TableName:  ctx.GlobalString(flags.DynamoDBTableNameFlag.Name),
		},
		EthClientConfig: geth.ReadEthClientConfig(ctx),
		AwsClientConfig: aws.ReadClientConfig(ctx, flags.FlagPrefix),
		EncoderConfig:   encoding.ReadCLIConfig(ctx),
		LoggerConfig:    logging.ReadCLIConfig(ctx, flags.FlagPrefix),
		BatcherConfig: batcher.Config{
			PullInterval:             ctx.GlobalDuration(flags.PullIntervalFlag.Name),
			FinalizerInterval:        ctx.GlobalDuration(flags.FinalizerIntervalFlag.Name),
			EncoderSocket:            ctx.GlobalString(flags.EncoderSocket.Name),
			NumConnections:           ctx.GlobalInt(flags.NumConnectionsFlag.Name),
			EncodingRequestQueueSize: ctx.GlobalInt(flags.EncodingRequestQueueSizeFlag.Name),
			BatchSizeMBLimit:         ctx.GlobalUint(flags.BatchSizeLimitFlag.Name),
			SRSOrder:                 ctx.GlobalInt(flags.SRSOrderFlag.Name),
			MaxNumRetriesPerBlob:     ctx.GlobalUint(flags.MaxNumRetriesPerBlobFlag.Name),
			TargetNumChunks:          ctx.GlobalUint(flags.TargetNumChunksFlag.Name),
			MaxBlobsToFetchFromStore: maxBlobsToFetchFromStore,
		},
		TimeoutConfig: batcher.TimeoutConfig{
			EncodingTimeout:    ctx.GlobalDuration(flags.EncodingTimeoutFlag.Name),
			AttestationTimeout: ctx.GlobalDuration(flags.AttestationTimeoutFlag.Name),
			ChainReadTimeout:   ctx.GlobalDuration(flags.ChainReadTimeoutFlag.Name),
			ChainWriteTimeout:  ctx.GlobalDuration(flags.ChainWriteTimeoutFlag.Name),
		},
		MetricsConfig: batcher.MetricsConfig{
			HTTPPort:      ctx.GlobalString(flags.MetricsHTTPPort.Name),
			EnableMetrics: ctx.GlobalBool(flags.EnableMetrics.Name),
		},
		UseGraph:                      ctx.Bool(flags.UseGraphFlag.Name),
		GraphUrl:                      ctx.GlobalString(flags.GraphUrlFlag.Name),
		BLSOperatorStateRetrieverAddr: ctx.GlobalString(flags.BlsOperatorStateRetrieverFlag.Name),
		EigenDAServiceManagerAddr:     ctx.GlobalString(flags.EigenDAServiceManagerFlag.Name),
		IndexerDataDir:                ctx.GlobalString(flags.IndexerDataDirFlag.Name),
		IndexerConfig:                 indexer.ReadIndexerConfig(ctx),
	}
	return config
}
