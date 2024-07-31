package config

import (
	"errors"
	"fmt"
	"github.com/Layr-Labs/eigenda/api/clients"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/Layr-Labs/eigenda/retriever"
	"time"

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

	//LoggerConfig     common.LoggerConfig
	//IndexerConfig    indexer.Config
	//MetricsConfig    MetricsConfig
	//ChainStateConfig thegraph.Config
	//
	//IndexerDataDir                string
	//Timeout                       time.Duration
	//NumConnections                int
	//BLSOperatorStateRetrieverAddr string
	//EigenDAServiceManagerAddr     string
	//UseGraph                      bool

	// The private key to use for signing requests.
	SignerPrivateKey string
	// Custom quorum numbers to use for the traffic generator.
	CustomQuorums []uint8
	// Whether to disable TLS for an insecure connection.
	DisableTlS bool
	// The port at which the metrics server listens for HTTP requests.
	MetricsHTTPPort string

	// The address of the BLS operator state retriever smart contract, in hex.
	BlsOperatorStateRetriever string

	// The URL of the subgraph instance.
	TheGraphUrl string
	// The interval at which to pull data from the subgraph.
	TheGraphPullInterval time.Duration
	// The number of times to retry a subgraph request.
	TheGraphRetries uint

	// The number of connections to use for the retriever.
	RetrieverNumConnections uint
	// The timeout for the node client.
	NodeClientTimeout time.Duration

	// The amount of time to sleep after launching each worker thread.
	InstanceLaunchInterval time.Duration

	// Configures the traffic generator workers.
	WorkerConfig WorkerConfig
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
		DisperserClientConfig: &clients.Config{
			Hostname:          ctx.GlobalString(HostnameFlag.Name),
			Port:              ctx.GlobalString(GrpcPortFlag.Name),
			Timeout:           ctx.Duration(TimeoutFlag.Name),
			UseSecureGrpcFlag: ctx.GlobalBool(UseSecureGrpcFlag.Name),
		},

		// TODO refactor flags
		RetrievalClientConfig: &retriever.Config{
			EigenDAServiceManagerAddr: ctx.String(EigenDAServiceManagerFlag.Name),
			EncoderConfig: kzg.KzgConfig{
				G1Path:          ctx.String(EncoderG1PathFlag.Name),
				G2Path:          ctx.String(EncoderG2PathFlag.Name),
				CacheDir:        ctx.String(EncoderCacheDirFlag.Name),
				SRSOrder:        ctx.Uint64(EncoderSRSOrderFlag.Name),
				SRSNumberToLoad: ctx.Uint64(EncoderSRSNumberToLoadFlag.Name),
				NumWorker:       ctx.Uint64(EncoderNumWorkersFlag.Name),
			},
			EthClientConfig: geth.EthClientConfig{
				RPCURLs:    []string{fmt.Sprintf("%s:%s", ctx.GlobalString(EthClientHostnameFlag.Name), ctx.GlobalString(EthClientPortFlag.Name))},
				NumRetries: ctx.Int(EthClientRetriesFlag.Name),
			},
		},

		LoggingConfig:    *loggerConfig,
		SignerPrivateKey: ctx.String(SignerPrivateKeyFlag.Name),
		CustomQuorums:    customQuorumsUint8,
		DisableTlS:       ctx.GlobalBool(DisableTLSFlag.Name),
		MetricsHTTPPort:  ctx.GlobalString(MetricsHTTPPortFlag.Name),

		BlsOperatorStateRetriever: ctx.String(BLSOperatorStateRetrieverFlag.Name),

		TheGraphUrl:          ctx.String(TheGraphUrlFlag.Name),
		TheGraphPullInterval: ctx.Duration(TheGraphPullIntervalFlag.Name),
		TheGraphRetries:      ctx.Uint(TheGraphRetriesFlag.Name),

		RetrieverNumConnections: ctx.Uint(RetrieverNumConnectionsFlag.Name),
		NodeClientTimeout:       ctx.Duration(NodeClientTimeoutFlag.Name),

		InstanceLaunchInterval: ctx.Duration(InstanceLaunchIntervalFlag.Name),

		WorkerConfig: WorkerConfig{
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
