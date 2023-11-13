package retriever

import (
	"time"

	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/common/logging"
	"github.com/Layr-Labs/eigenda/core/encoding"
	"github.com/Layr-Labs/eigenda/indexer"
	"github.com/Layr-Labs/eigenda/retriever/flags"
	"github.com/urfave/cli"
)

type Config struct {
	EncoderConfig   encoding.EncoderConfig
	EthClientConfig geth.EthClientConfig
	LoggerConfig    logging.Config
	IndexerConfig   indexer.Config
	MetricsConfig   MetricsConfig

	IndexerDataDir                string
	Timeout                       time.Duration
	NumConnections                int
	BLSOperatorStateRetrieverAddr string
	EigenDAServiceManagerAddr     string
}

func NewConfig(ctx *cli.Context) *Config {
	return &Config{
		EncoderConfig:   encoding.ReadCLIConfig(ctx),
		EthClientConfig: geth.ReadEthClientConfig(ctx),
		LoggerConfig:    logging.ReadCLIConfig(ctx, flags.FlagPrefix),
		IndexerConfig:   indexer.ReadIndexerConfig(ctx),
		MetricsConfig: MetricsConfig{
			HTTPPort: ctx.GlobalString(flags.MetricsHTTPPortFlag.Name),
		},
		IndexerDataDir:                ctx.GlobalString(flags.IndexerDataDirFlag.Name),
		Timeout:                       ctx.Duration(flags.TimeoutFlag.Name),
		NumConnections:                ctx.Int(flags.NumConnectionsFlag.Name),
		BLSOperatorStateRetrieverAddr: ctx.GlobalString(flags.BlsOperatorStateRetrieverFlag.Name),
		EigenDAServiceManagerAddr:     ctx.GlobalString(flags.EigenDAServiceManagerFlag.Name),
	}
}
