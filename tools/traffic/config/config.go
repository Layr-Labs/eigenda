package config

import (
	"errors"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients/v2"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/core/thegraph"
	"github.com/urfave/cli"
)

// Config configures a traffic generator.
type Config struct {
	// Logging configuration.
	LoggingConfig common.LoggerConfig

	// Configuration for the disperser client.
	DisperserClientConfig *clients.DisperserClientConfig

	// Signer private key
	SignerPrivateKey string

	// Configuration for the graph.
	TheGraphConfig *thegraph.Config

	// Configures the blob writers.
	BlobWriterConfig BlobWriterConfig

	// The port at which the metrics server listens for HTTP requests.
	MetricsHTTPPort string

	// The timeout for the node client.
	NodeClientTimeout time.Duration

	// The amount of time to sleep after launching each worker thread.
	InstanceLaunchInterval time.Duration
}

// BlobWriterConfig configures the blob writer.
type BlobWriterConfig struct {
	// The number of worker threads that generate write traffic.
	NumWriteInstances uint

	// The period of the submission rate of new blobs for each write worker thread.
	WriteRequestInterval time.Duration

	// The Size of each blob dispersed, in bytes.
	DataSize uint64

	// If true, then each blob will contain unique random data. If false, the same random data
	// will be dispersed for each blob by a particular worker thread.
	RandomizeBlobs bool

	// The amount of time to wait for a blob to be written.
	WriteTimeout time.Duration

	// Custom quorum numbers to use for the traffic generator.
	CustomQuorums []uint8
}

func NewConfig(ctx *cli.Context) (*Config, error) {
	loggerConfig, err := common.ReadLoggerCLIConfig(ctx, FlagPrefix)
	if err != nil {
		return nil, err
	}
	customQuorums := ctx.GlobalIntSlice(CustomQuorumNumbersFlag.Name)
	if len(customQuorums) == 0 {
		return nil, errors.New("no custom quorum numbers provided")
	}

	customQuorumsUint8 := make([]uint8, len(customQuorums))
	for i, q := range customQuorums {
		if q < 0 || q > 255 {
			return nil, errors.New("invalid custom quorum number")
		}
		customQuorumsUint8[i] = uint8(q)
	}

	config := &Config{
		DisperserClientConfig: &clients.DisperserClientConfig{
			Hostname:          ctx.GlobalString(HostnameFlag.Name),
			Port:              ctx.GlobalString(GrpcPortFlag.Name),
			UseSecureGrpcFlag: ctx.GlobalBool(UseSecureGrpcFlag.Name),
		},

		SignerPrivateKey: ctx.String(SignerPrivateKeyFlag.Name),
		LoggingConfig:    *loggerConfig,

		MetricsHTTPPort:   ctx.GlobalString(MetricsHTTPPortFlag.Name),
		NodeClientTimeout: ctx.Duration(NodeClientTimeoutFlag.Name),

		InstanceLaunchInterval: ctx.Duration(InstanceLaunchIntervalFlag.Name),

		BlobWriterConfig: BlobWriterConfig{
			NumWriteInstances:    ctx.GlobalUint(NumWriteInstancesFlag.Name),
			WriteRequestInterval: ctx.Duration(WriteRequestIntervalFlag.Name),
			DataSize:             ctx.GlobalUint64(DataSizeFlag.Name),
			RandomizeBlobs:       ctx.GlobalBool(RandomizeBlobsFlag.Name),
			WriteTimeout:         ctx.Duration(WriteTimeoutFlag.Name),
			CustomQuorums:        customQuorumsUint8,
		},
	}

	return config, nil
}
