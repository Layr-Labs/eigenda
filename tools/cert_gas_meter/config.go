package cert_gas_meter

import (
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/thegraph"
	"github.com/Layr-Labs/eigenda/tools/cert_gas_meter/flags"
	"github.com/urfave/cli"
)

type Config struct {
	LoggerConfig       common.LoggerConfig
	BlockNumber        uint64
	Workers            int
	Timeout            time.Duration
	UseRetrievalClient bool
	QuorumIDs          []core.QuorumID
	TopN               uint
	OutputFormat       string
	OutputFile         string

	ChainStateConfig thegraph.Config
	EthClientConfig  geth.EthClientConfig

	EigenDADirectory           string
	OperatorStateRetrieverAddr string
	EigenDAServiceManagerAddr  string

	CertPath string
}

func ReadConfig(ctx *cli.Context) *Config {

	return &Config{
		ChainStateConfig: thegraph.ReadCLIConfig(ctx),
		EthClientConfig:  geth.ReadEthClientConfig(ctx),
		//EigenDADirectory:           ctx.GlobalString(flags.EigenDADirectoryFlag.Name),
		OperatorStateRetrieverAddr: ctx.GlobalString(flags.OperatorStateRetrieverFlag.Name),
		//EigenDAServiceManagerAddr:  ctx.GlobalString(flags.EigenDAServiceManagerFlag.Name),
		CertPath: ctx.GlobalString(flags.CertFileFlag.Name),
	}
}

func NewConfig(ctx *cli.Context) (*Config, error) {
	loggerConfig, err := common.ReadLoggerCLIConfig(ctx, flags.FlagPrefix)
	if err != nil {
		return nil, err
	}

	config := ReadConfig(ctx)
	config.LoggerConfig = *loggerConfig

	return config, nil
}
