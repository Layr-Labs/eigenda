package lib

import (
	"fmt"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/aws"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/core/thegraph"
	core "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/relay"
	"github.com/Layr-Labs/eigenda/relay/cmd/flags"
	"github.com/Layr-Labs/eigenda/relay/limiter"
	"github.com/urfave/cli"
)

// Config is the configuration for the relay Server.
type Config struct {

	// Log is the configuration for the logger. Default is common.DefaultLoggerConfig().
	Log common.LoggerConfig

	// Configuration for the AWS client. Default is aws.DefaultClientConfig().
	AWS aws.ClientConfig

	// BucketName is the name of the S3 bucket that stores blobs. Default is "relay".
	BucketName string

	// MetadataTableName is the name of the DynamoDB table that stores metadata. Default is "metadata".
	MetadataTableName string

	// RelayConfig is the configuration for the relay.
	RelayConfig relay.Config

	// Configuration for the graph indexer.
	EthClientConfig               geth.EthClientConfig
	EigenDADirectory              string
	BLSOperatorStateRetrieverAddr string
	EigenDAServiceManagerAddr     string
	ChainStateConfig              thegraph.Config
}

func NewConfig(ctx *cli.Context) (Config, error) {
	loggerConfig, err := common.ReadLoggerCLIConfig(ctx, flags.FlagPrefix)
	if err != nil {
		return Config{}, err
	}
	awsClientConfig := aws.ReadClientConfig(ctx, flags.FlagPrefix)
	relayKeys := ctx.IntSlice(flags.RelayKeysFlag.Name)
	if len(relayKeys) == 0 {
		return Config{}, fmt.Errorf("no relay keys specified")
	}

	config := Config{
		Log:               *loggerConfig,
		AWS:               awsClientConfig,
		BucketName:        ctx.String(flags.BucketNameFlag.Name),
		MetadataTableName: ctx.String(flags.MetadataTableNameFlag.Name),
		RelayConfig: relay.Config{
			RelayKeys:                  make([]core.RelayKey, len(relayKeys)),
			GRPCPort:                   ctx.Int(flags.GRPCPortFlag.Name),
			MaxGRPCMessageSize:         ctx.Int(flags.MaxGRPCMessageSizeFlag.Name),
			MetadataCacheSize:          ctx.Int(flags.MetadataCacheSizeFlag.Name),
			MetadataMaxConcurrency:     ctx.Int(flags.MetadataMaxConcurrencyFlag.Name),
			BlobCacheBytes:             ctx.Uint64(flags.BlobCacheBytes.Name),
			BlobMaxConcurrency:         ctx.Int(flags.BlobMaxConcurrencyFlag.Name),
			ChunkCacheBytes:            ctx.Uint64(flags.ChunkCacheBytesFlag.Name),
			ChunkMaxConcurrency:        ctx.Int(flags.ChunkMaxConcurrencyFlag.Name),
			MaxKeysPerGetChunksRequest: ctx.Int(flags.MaxKeysPerGetChunksRequestFlag.Name),
			RateLimits: limiter.Config{
				MaxGetBlobOpsPerSecond:          ctx.Float64(flags.MaxGetBlobOpsPerSecondFlag.Name),
				GetBlobOpsBurstiness:            ctx.Int(flags.GetBlobOpsBurstinessFlag.Name),
				MaxGetBlobBytesPerSecond:        ctx.Float64(flags.MaxGetBlobBytesPerSecondFlag.Name),
				GetBlobBytesBurstiness:          ctx.Int(flags.GetBlobBytesBurstinessFlag.Name),
				MaxConcurrentGetBlobOps:         ctx.Int(flags.MaxConcurrentGetBlobOpsFlag.Name),
				MaxGetChunkOpsPerSecond:         ctx.Float64(flags.MaxGetChunkOpsPerSecondFlag.Name),
				GetChunkOpsBurstiness:           ctx.Int(flags.GetChunkOpsBurstinessFlag.Name),
				MaxGetChunkBytesPerSecond:       ctx.Float64(flags.MaxGetChunkBytesPerSecondFlag.Name),
				GetChunkBytesBurstiness:         ctx.Int(flags.GetChunkBytesBurstinessFlag.Name),
				MaxConcurrentGetChunkOps:        ctx.Int(flags.MaxConcurrentGetChunkOpsFlag.Name),
				MaxGetChunkOpsPerSecondClient:   ctx.Float64(flags.MaxGetChunkOpsPerSecondClientFlag.Name),
				GetChunkOpsBurstinessClient:     ctx.Int(flags.GetChunkOpsBurstinessClientFlag.Name),
				MaxGetChunkBytesPerSecondClient: ctx.Float64(flags.MaxGetChunkBytesPerSecondClientFlag.Name),
				GetChunkBytesBurstinessClient:   ctx.Int(flags.GetChunkBytesBurstinessClientFlag.Name),
				MaxConcurrentGetChunkOpsClient:  ctx.Int(flags.MaxConcurrentGetChunkOpsClientFlag.Name),
			},
			AuthenticationKeyCacheSize:   ctx.Int(flags.AuthenticationKeyCacheSizeFlag.Name),
			AuthenticationDisabled:       ctx.Bool(flags.AuthenticationDisabledFlag.Name),
			GetChunksRequestMaxPastAge:   ctx.Duration(flags.GetChunksRequestMaxPastAgeFlag.Name),
			GetChunksRequestMaxFutureAge: ctx.Duration(flags.GetChunksRequestMaxFutureAgeFlag.Name),
			OnchainStateRefreshInterval:  ctx.Duration(flags.OnchainStateRefreshIntervalFlag.Name),
			Timeouts: relay.TimeoutConfig{
				GetChunksTimeout:               ctx.Duration(flags.GetChunksTimeoutFlag.Name),
				GetBlobTimeout:                 ctx.Duration(flags.GetBlobTimeoutFlag.Name),
				InternalGetMetadataTimeout:     ctx.Duration(flags.InternalGetMetadataTimeoutFlag.Name),
				InternalGetBlobTimeout:         ctx.Duration(flags.InternalGetBlobTimeoutFlag.Name),
				InternalGetProofsTimeout:       ctx.Duration(flags.InternalGetProofsTimeoutFlag.Name),
				InternalGetCoefficientsTimeout: ctx.Duration(flags.InternalGetCoefficientsTimeoutFlag.Name),
			},
			MetricsPort:           ctx.Int(flags.MetricsPortFlag.Name),
			EnableMetrics:         ctx.Bool(flags.EnableMetricsFlag.Name),
			EnablePprof:           ctx.Bool(flags.EnablePprofFlag.Name),
			PprofHttpPort:         ctx.Int(flags.PprofHttpPortFlag.Name),
			MaxConnectionAge:      ctx.Duration(flags.MaxConnectionAgeFlag.Name),
			MaxConnectionAgeGrace: ctx.Duration(flags.MaxConnectionAgeGraceFlag.Name),
			MaxIdleConnectionAge:  ctx.Duration(flags.MaxIdleConnectionAgeFlag.Name),
		},
		EthClientConfig:               geth.ReadEthClientConfigRPCOnly(ctx),
		EigenDADirectory:              ctx.String(flags.EigenDADirectoryFlag.Name),
		BLSOperatorStateRetrieverAddr: ctx.String(flags.BlsOperatorStateRetrieverAddrFlag.Name),
		EigenDAServiceManagerAddr:     ctx.String(flags.EigenDAServiceManagerAddrFlag.Name),
		ChainStateConfig:              thegraph.ReadCLIConfig(ctx),
	}
	for i, id := range relayKeys {
		config.RelayConfig.RelayKeys[i] = core.RelayKey(id)
	}
	return config, nil
}
