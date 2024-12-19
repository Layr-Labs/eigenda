package config

import (
	"errors"
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients"
	"github.com/Layr-Labs/eigenda/core/thegraph"
	"github.com/Layr-Labs/eigenda/retriever"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/urfave/cli"
)

// Config configures a traffic generator.
type Config struct {

	// Logging configuration.
	LoggingConfig common.LoggerConfig

	// Configuration for the disperser client.
	DisperserClientConfig *clients.Config

	// Configuration for the retriever client.
	RetrievalClientConfig *retriever.Config

	// Configuration for the graph.
	TheGraphConfig *thegraph.Config

	// Configuration for the EigenDA client.
	EigenDAClientConfig *clients.EigenDAClientConfig

	// Configures the traffic generator workers.
	WorkerConfig WorkerConfig

	// The port at which the metrics server listens for HTTP requests.
	MetricsHTTPPort string
	// The timeout for the node client.
	NodeClientTimeout time.Duration
	// The amount of time to sleep after launching each worker thread.
	InstanceLaunchInterval time.Duration
}

func NewConfig(ctx *cli.Context) (*Config, error) {
	loggerConfig, err := common.ReadLoggerCLIConfig(ctx, FlagPrefix)
	if err != nil {
		return nil, err
	}
	customQuorums := ctx.GlobalIntSlice(CustomQuorumNumbersFlag.Name)
	customQuorumsUint8 := make([]uint8, len(customQuorums))
	for i, q := range customQuorums {
		if q < 0 || q > 255 {
			return nil, errors.New("invalid custom quorum number")
		}
		customQuorumsUint8[i] = uint8(q)
	}

	retrieverConfig, err := retriever.NewConfig(ctx)
	if err != nil {
		return nil, err
	}

	config := &Config{
		DisperserClientConfig: &clients.Config{
			Hostname:          ctx.GlobalString(HostnameFlag.Name),
			Port:              ctx.GlobalString(GrpcPortFlag.Name),
			Timeout:           ctx.Duration(TimeoutFlag.Name),
			UseSecureGrpcFlag: ctx.GlobalBool(UseSecureGrpcFlag.Name),
		},

		RetrievalClientConfig: retrieverConfig,

		TheGraphConfig: &thegraph.Config{
			Endpoint:     ctx.String(TheGraphUrlFlag.Name),
			PullInterval: ctx.Duration(TheGraphPullIntervalFlag.Name),
			MaxRetries:   ctx.Int(TheGraphRetriesFlag.Name),
		},

		EigenDAClientConfig: &clients.EigenDAClientConfig{
			RPC:                 fmt.Sprintf("%s:%s", ctx.GlobalString(HostnameFlag.Name), ctx.GlobalString(GrpcPortFlag.Name)),
			SignerPrivateKeyHex: ctx.String(SignerPrivateKeyFlag.Name),
			DisableTLS:          ctx.GlobalBool(DisableTLSFlag.Name),
		},

		LoggingConfig: *loggerConfig,

		MetricsHTTPPort:   ctx.GlobalString(MetricsHTTPPortFlag.Name),
		NodeClientTimeout: ctx.Duration(NodeClientTimeoutFlag.Name),

		InstanceLaunchInterval: ctx.Duration(InstanceLaunchIntervalFlag.Name),

		WorkerConfig: WorkerConfig{
			NumWriteInstances:    ctx.GlobalUint(NumWriteInstancesFlag.Name),
			WriteRequestInterval: ctx.Duration(WriteRequestIntervalFlag.Name),
			DataSize:             ctx.GlobalUint64(DataSizeFlag.Name),
			RandomizeBlobs:       !ctx.GlobalBool(UniformBlobsFlag.Name),
			WriteTimeout:         ctx.Duration(WriteTimeoutFlag.Name),

			TrackerInterval:      ctx.Duration(VerifierIntervalFlag.Name),
			GetBlobStatusTimeout: ctx.Duration(GetBlobStatusTimeoutFlag.Name),

			NumReadInstances:             ctx.GlobalUint(NumReadInstancesFlag.Name),
			ReadRequestInterval:          ctx.Duration(ReadRequestIntervalFlag.Name),
			RequiredDownloads:            ctx.Float64(RequiredDownloadsFlag.Name),
			FetchBatchHeaderTimeout:      ctx.Duration(FetchBatchHeaderTimeoutFlag.Name),
			RetrieveBlobChunksTimeout:    ctx.Duration(RetrieveBlobChunksTimeoutFlag.Name),
			StatusTrackerChannelCapacity: ctx.Uint(VerificationChannelCapacityFlag.Name),

			EigenDAServiceManager: retrieverConfig.EigenDAServiceManagerAddr,
			SignerPrivateKey:      ctx.String(SignerPrivateKeyFlag.Name),
			CustomQuorums:         customQuorumsUint8,

			MetricsBlacklist:      ctx.StringSlice(MetricsBlacklistFlag.Name),
			MetricsFuzzyBlacklist: ctx.StringSlice(MetricsFuzzyBlacklistFlag.Name),
		},
	}

	err = config.EigenDAClientConfig.CheckAndSetDefaults()
	if err != nil {
		return nil, err
	}

	return config, nil
}
