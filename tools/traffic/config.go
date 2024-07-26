package traffic

import (
	"errors"
	"github.com/Layr-Labs/eigenda/tools/traffic/workers"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/tools/traffic/flags"
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
	WorkerConfig workers.WorkerConfig
}

func NewConfig(ctx *cli.Context) (*Config, error) {
	loggerConfig, err := common.ReadLoggerCLIConfig(ctx, flags.FlagPrefix)
	if err != nil {
		return nil, err
	}
	customQuorums := ctx.GlobalIntSlice(flags.CustomQuorumNumbersFlag.Name)
	customQuorumsUint8 := make([]uint8, len(customQuorums))
	for i, q := range customQuorums {
		if q < 0 || q > 255 {
			return nil, errors.New("invalid custom quorum number")
		}
		customQuorumsUint8[i] = uint8(q)
	}

	return &Config{
		LoggingConfig:              *loggerConfig,
		DisperserHostname:          ctx.GlobalString(flags.HostnameFlag.Name),
		DisperserPort:              ctx.GlobalString(flags.GrpcPortFlag.Name),
		DisperserTimeout:           ctx.Duration(flags.TimeoutFlag.Name),
		DisperserUseSecureGrpcFlag: ctx.GlobalBool(flags.UseSecureGrpcFlag.Name),
		SignerPrivateKey:           ctx.String(flags.SignerPrivateKeyFlag.Name),
		CustomQuorums:              customQuorumsUint8,
		DisableTlS:                 ctx.GlobalBool(flags.DisableTLSFlag.Name),
		MetricsHTTPPort:            ctx.GlobalString(flags.MetricsHTTPPortFlag.Name),
		EthClientHostname:          ctx.GlobalString(flags.EthClientHostnameFlag.Name),
		EthClientPort:              ctx.GlobalString(flags.EthClientPortFlag.Name),
		BlsOperatorStateRetriever:  ctx.String(flags.BLSOperatorStateRetrieverFlag.Name),
		EigenDAServiceManager:      ctx.String(flags.EigenDAServiceManagerFlag.Name),
		EthClientRetries:           ctx.Uint(flags.EthClientRetriesFlag.Name),
		TheGraphUrl:                ctx.String(flags.TheGraphUrlFlag.Name),
		TheGraphPullInterval:       ctx.Duration(flags.TheGraphPullIntervalFlag.Name),
		TheGraphRetries:            ctx.Uint(flags.TheGraphRetriesFlag.Name),
		EncoderG1Path:              ctx.String(flags.EncoderG1PathFlag.Name),
		EncoderG2Path:              ctx.String(flags.EncoderG2PathFlag.Name),
		EncoderCacheDir:            ctx.String(flags.EncoderCacheDirFlag.Name),
		EncoderSRSOrder:            ctx.Uint64(flags.EncoderSRSOrderFlag.Name),
		EncoderSRSNumberToLoad:     ctx.Uint64(flags.EncoderSRSNumberToLoadFlag.Name),
		EncoderNumWorkers:          ctx.Uint64(flags.EncoderNumWorkersFlag.Name),
		RetrieverNumConnections:    ctx.Uint(flags.RetrieverNumConnectionsFlag.Name),
		NodeClientTimeout:          ctx.Duration(flags.NodeClientTimeoutFlag.Name),

		InstanceLaunchInterval: ctx.Duration(flags.InstanceLaunchIntervalFlag.Name),

		WorkerConfig: workers.WorkerConfig{
			NumWriteInstances:    ctx.GlobalUint(flags.NumWriteInstancesFlag.Name),
			WriteRequestInterval: ctx.Duration(flags.WriteRequestIntervalFlag.Name),
			DataSize:             ctx.GlobalUint64(flags.DataSizeFlag.Name),
			RandomizeBlobs:       !ctx.GlobalBool(flags.UniformBlobsFlag.Name),
			WriteTimeout:         ctx.Duration(flags.WriteTimeoutFlag.Name),

			VerifierInterval:     ctx.Duration(flags.VerifierIntervalFlag.Name),
			GetBlobStatusTimeout: ctx.Duration(flags.GetBlobStatusTimeoutFlag.Name),

			NumReadInstances:          ctx.GlobalUint(flags.NumReadInstancesFlag.Name),
			ReadRequestInterval:       ctx.Duration(flags.ReadRequestIntervalFlag.Name),
			RequiredDownloads:         ctx.Float64(flags.RequiredDownloadsFlag.Name),
			ReadOverflowTableSize:     uint32(ctx.Uint(flags.ReadOverflowTableSizeFlag.Name)),
			FetchBatchHeaderTimeout:   ctx.Duration(flags.FetchBatchHeaderTimeoutFlag.Name),
			RetrieveBlobChunksTimeout: ctx.Duration(flags.RetrieveBlobChunksTimeoutFlag.Name),

			EigenDAServiceManager: ctx.String(flags.EigenDAServiceManagerFlag.Name),
			SignerPrivateKey:      ctx.String(flags.SignerPrivateKeyFlag.Name),
			CustomQuorums:         customQuorumsUint8,
		},
	}, nil
}
