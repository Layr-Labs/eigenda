package config

import (
	"errors"
	"github.com/Layr-Labs/eigenda/tools/traffic/workers"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/urfave/cli"
)

// Config configures a traffic generator.
type Config struct {
	// Logging configuration.
	LoggingConfig common.LoggerConfig
	// The hostname of the disperser.
	DisperserHostname string
	// The port of the disperser.
	DisperserPort string
	// The timeout for the disperser.
	DisperserTimeout time.Duration
	// Whether to use a secure gRPC connection to the disperser.
	DisperserUseSecureGrpcFlag bool
	// The private key to use for signing requests.
	SignerPrivateKey string
	// Custom quorum numbers to use for the traffic generator.
	CustomQuorums []uint8
	// Whether to disable TLS for an insecure connection.
	DisableTlS bool
	// The port at which the metrics server listens for HTTP requests.
	MetricsHTTPPort string
	// The hostname of the Ethereum client.
	EthClientHostname string
	// The port of the Ethereum client.
	EthClientPort string
	// The address of the BLS operator state retriever smart contract, in hex.
	BlsOperatorStateRetriever string
	// The address of the EigenDA service manager smart contract, in hex.
	EigenDAServiceManager string
	// The number of times to retry an Ethereum client request.
	EthClientRetries uint
	// The URL of the subgraph instance.
	TheGraphUrl string
	// The interval at which to pull data from the subgraph.
	TheGraphPullInterval time.Duration
	// The number of times to retry a subgraph request.
	TheGraphRetries uint
	// The path to the encoder G1 binary.
	EncoderG1Path string
	// The path to the encoder G2 binary.
	EncoderG2Path string
	// The path to the encoder cache directory.
	EncoderCacheDir string
	// The SRS order to use for the encoder.
	EncoderSRSOrder uint64
	// The SRS number to load for the encoder.
	EncoderSRSNumberToLoad uint64
	// The number of worker threads to use for the encoder.
	EncoderNumWorkers uint64
	// The number of connections to use for the retriever.
	RetrieverNumConnections uint
	// The timeout for the node client.
	NodeClientTimeout time.Duration

	// The amount of time to sleep after launching each worker thread.
	InstanceLaunchInterval time.Duration

	// Configures the traffic generator workers.
	WorkerConfig workers.Config
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

	return &Config{
		LoggingConfig:              *loggerConfig,
		DisperserHostname:          ctx.GlobalString(HostnameFlag.Name),
		DisperserPort:              ctx.GlobalString(GrpcPortFlag.Name),
		DisperserTimeout:           ctx.Duration(TimeoutFlag.Name),
		DisperserUseSecureGrpcFlag: ctx.GlobalBool(UseSecureGrpcFlag.Name),
		SignerPrivateKey:           ctx.String(SignerPrivateKeyFlag.Name),
		CustomQuorums:              customQuorumsUint8,
		DisableTlS:                 ctx.GlobalBool(DisableTLSFlag.Name),
		MetricsHTTPPort:            ctx.GlobalString(MetricsHTTPPortFlag.Name),
		EthClientHostname:          ctx.GlobalString(EthClientHostnameFlag.Name),
		EthClientPort:              ctx.GlobalString(EthClientPortFlag.Name),
		BlsOperatorStateRetriever:  ctx.String(BLSOperatorStateRetrieverFlag.Name),
		EigenDAServiceManager:      ctx.String(EigenDAServiceManagerFlag.Name),
		EthClientRetries:           ctx.Uint(EthClientRetriesFlag.Name),
		TheGraphUrl:                ctx.String(TheGraphUrlFlag.Name),
		TheGraphPullInterval:       ctx.Duration(TheGraphPullIntervalFlag.Name),
		TheGraphRetries:            ctx.Uint(TheGraphRetriesFlag.Name),
		EncoderG1Path:              ctx.String(EncoderG1PathFlag.Name),
		EncoderG2Path:              ctx.String(EncoderG2PathFlag.Name),
		EncoderCacheDir:            ctx.String(EncoderCacheDirFlag.Name),
		EncoderSRSOrder:            ctx.Uint64(EncoderSRSOrderFlag.Name),
		EncoderSRSNumberToLoad:     ctx.Uint64(EncoderSRSNumberToLoadFlag.Name),
		EncoderNumWorkers:          ctx.Uint64(EncoderNumWorkersFlag.Name),
		RetrieverNumConnections:    ctx.Uint(RetrieverNumConnectionsFlag.Name),
		NodeClientTimeout:          ctx.Duration(NodeClientTimeoutFlag.Name),

		InstanceLaunchInterval: ctx.Duration(InstanceLaunchIntervalFlag.Name),

		WorkerConfig: workers.Config{
			NumWriteInstances:    ctx.GlobalUint(NumWriteInstancesFlag.Name),
			WriteRequestInterval: ctx.Duration(WriteRequestIntervalFlag.Name),
			DataSize:             ctx.GlobalUint64(DataSizeFlag.Name),
			RandomizeBlobs:       !ctx.GlobalBool(UniformBlobsFlag.Name),
			WriteTimeout:         ctx.Duration(WriteTimeoutFlag.Name),

			VerifierInterval:     ctx.Duration(VerifierIntervalFlag.Name),
			GetBlobStatusTimeout: ctx.Duration(GetBlobStatusTimeoutFlag.Name),

			NumReadInstances:            ctx.GlobalUint(NumReadInstancesFlag.Name),
			ReadRequestInterval:         ctx.Duration(ReadRequestIntervalFlag.Name),
			RequiredDownloads:           ctx.Float64(RequiredDownloadsFlag.Name),
			ReadOverflowTableSize:       ctx.Uint(ReadOverflowTableSizeFlag.Name),
			FetchBatchHeaderTimeout:     ctx.Duration(FetchBatchHeaderTimeoutFlag.Name),
			RetrieveBlobChunksTimeout:   ctx.Duration(RetrieveBlobChunksTimeoutFlag.Name),
			VerificationChannelCapacity: ctx.Uint(VerificationChannelCapacityFlag.Name),

			EigenDAServiceManager: ctx.String(EigenDAServiceManagerFlag.Name),
			SignerPrivateKey:      ctx.String(SignerPrivateKeyFlag.Name),
			CustomQuorums:         customQuorumsUint8,
		},
	}, nil
}
