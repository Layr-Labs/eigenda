package semverscan

import (
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/core/eth"
	"github.com/Layr-Labs/eigenda/core/thegraph"
	"github.com/Layr-Labs/eigenda/tools/semverscan/flags"
	"github.com/urfave/cli"
)

type Config struct {
	LoggerConfig       common.LoggerConfig
	Workers            int
	OperatorId         string
	Timeout            time.Duration
	UseRetrievalClient bool

	ChainStateConfig thegraph.Config
	EthClientConfig  geth.EthClientConfig

	AddressDirectoryAddr          string
	BLSOperatorStateRetrieverAddr string
	EigenDAServiceManagerAddr     string
}

func ReadConfig(ctx *cli.Context) *Config {
	return &Config{
		Timeout:                       ctx.Duration(flags.TimeoutFlag.Name),
		Workers:                       ctx.Int(flags.WorkersFlag.Name),
		OperatorId:                    ctx.String(flags.OperatorIdFlag.Name),
		UseRetrievalClient:            ctx.Bool(flags.UseRetrievalClientFlag.Name),
		ChainStateConfig:              thegraph.ReadCLIConfig(ctx),
		EthClientConfig:               geth.ReadEthClientConfig(ctx),
		AddressDirectoryAddr:          ctx.GlobalString(flags.AddressDirectoryFlag.Name),
		BLSOperatorStateRetrieverAddr: ctx.GlobalString(flags.BlsOperatorStateRetrieverFlag.Name),
		EigenDAServiceManagerAddr:     ctx.GlobalString(flags.EigenDAServiceManagerFlag.Name),
	}
}

func NewConfig(ctx *cli.Context) (*Config, error) {
	loggerConfig, err := common.ReadLoggerCLIConfig(ctx, flags.FlagPrefix)
	if err != nil {
		return nil, err
	}

	config := ReadConfig(ctx)
	config.LoggerConfig = *loggerConfig

	// Validate that either address directory is provided OR both individual addresses are provided
	if err := eth.ValidateAddressConfig(config.AddressDirectoryAddr, config.BLSOperatorStateRetrieverAddr, config.EigenDAServiceManagerAddr); err != nil {
		return nil, err
	}

	return config, nil
}
