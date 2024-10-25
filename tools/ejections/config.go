package ejections

import (
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/tools/ejections/flags"
	"github.com/urfave/cli"
)

type Config struct {
	LoggerConfig     common.LoggerConfig
	Days             int
	OperatorId       string
	SubgraphEndpoint string
	First            uint
	Skip             uint

	EthClientConfig               geth.EthClientConfig
	BLSOperatorStateRetrieverAddr string
	EigenDAServiceManagerAddr     string
}

func ReadConfig(ctx *cli.Context) *Config {
	return &Config{
		Days:                          ctx.Int(flags.DaysFlag.Name),
		OperatorId:                    ctx.String(flags.OperatorIdFlag.Name),
		SubgraphEndpoint:              ctx.String(flags.SubgraphEndpointFlag.Name),
		First:                         ctx.Uint(flags.FirstFlag.Name),
		Skip:                          ctx.Uint(flags.SkipFlag.Name),
		EthClientConfig:               geth.ReadEthClientConfig(ctx),
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

	return config, nil
}
