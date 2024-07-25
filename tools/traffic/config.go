package traffic

import (
	"errors"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/tools/traffic/flags"
	"github.com/urfave/cli"
)

// Config configures a traffic generator.
type Config struct {
	LoggingConfig common.LoggerConfig
	clients.Config
	SignerPrivateKey          string
	CustomQuorums             []uint8
	DisableTlS                bool
	MetricsHTTPPort           string
	EthClientHostname         string
	EthClientPort             string
	BlsOperatorStateRetriever string
	EigenDAServiceManager     string

	// The amount of time to sleep after launching each worker thread.
	InstanceLaunchInterval time.Duration

	// The number of worker threads that generate write traffic.
	NumWriteInstances uint
	// The period of the submission rate of new blobs for each write worker thread.
	WriteRequestInterval time.Duration
	// The Size of each blob dispersed, in bytes.
	DataSize uint64
	// If true, then each blob will contain unique random data. If false, the same random data
	// will be dispersed for each blob by a particular worker thread.
	RandomizeBlobs bool

	// TODO add to flags.go
	// The number of worker threads that generate read traffic.
	NumReadInstances uint
	// The period of the submission rate of read requests for each read worker thread.
	ReadRequestInterval time.Duration
	// For each blob, how many times should it be downloaded? If between 0.0 and 1.0, blob will be downloaded
	// 0 or 1 times with the specified probability (e.g. 0.2 means each blob has a 20% chance of being downloaded).
	// If greater than 1.0, then each blob will be downloaded the specified number of times.
	RequiredDownloads float64
	// The minimum amount of time that must pass after a blob is written prior to the first read attempt being made.
	ReadDelay time.Duration
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
		LoggingConfig: *loggerConfig,
		Config: *clients.NewConfig(
			ctx.GlobalString(flags.HostnameFlag.Name),
			ctx.GlobalString(flags.GrpcPortFlag.Name),
			ctx.Duration(flags.TimeoutFlag.Name),
			ctx.GlobalBool(flags.UseSecureGrpcFlag.Name),
		),
		SignerPrivateKey:          ctx.String(flags.SignerPrivateKeyFlag.Name),
		CustomQuorums:             customQuorumsUint8,
		DisableTlS:                ctx.GlobalBool(flags.DisableTLSFlag.Name),
		MetricsHTTPPort:           ctx.GlobalString(flags.MetricsHTTPPortFlag.Name),
		EthClientHostname:         ctx.GlobalString(flags.EthClientHostnameFlag.Name),
		EthClientPort:             ctx.GlobalString(flags.EthClientPortFlag.Name),
		BlsOperatorStateRetriever: ctx.String(flags.BLSOperatorStateRetrieverFlag.Name),
		EigenDAServiceManager:     ctx.String(flags.EigenDAServiceManagerFlag.Name),

		InstanceLaunchInterval: ctx.Duration(flags.InstanceLaunchIntervalFlag.Name),

		NumWriteInstances:    ctx.GlobalUint(flags.NumWriteInstancesFlag.Name),
		WriteRequestInterval: ctx.Duration(flags.WriteRequestIntervalFlag.Name),
		DataSize:             ctx.GlobalUint64(flags.DataSizeFlag.Name),
		RandomizeBlobs:       !ctx.GlobalBool(flags.UniformBlobsFlag.Name),

		NumReadInstances:    ctx.GlobalUint(flags.NumReadInstancesFlag.Name),
		ReadRequestInterval: ctx.Duration(flags.ReadRequestIntervalFlag.Name),
		RequiredDownloads:   ctx.Float64(flags.RequiredDownloadsFlag.Name),
	}, nil
}
