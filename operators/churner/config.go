package churner

import (
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/core/eth"
	"github.com/Layr-Labs/eigenda/core/thegraph"
	"github.com/Layr-Labs/eigenda/operators/churner/flags"
	"github.com/urfave/cli"
)

type Config struct {
	EthClientConfig  geth.EthClientConfig
	LoggerConfig     common.LoggerConfig
	MetricsConfig    MetricsConfig
	ChainStateConfig thegraph.Config

	AddressDirectoryAddr          string
	BLSOperatorStateRetrieverAddr string // Legacy field, use AddressDirectoryAddr instead
	EigenDAServiceManagerAddr     string // Legacy field, use AddressDirectoryAddr instead

	PerPublicKeyRateLimit time.Duration
	ChurnApprovalInterval time.Duration
}

func NewConfig(ctx *cli.Context) (*Config, error) {
	loggerConfig, err := common.ReadLoggerCLIConfig(ctx, flags.FlagPrefix)
	if err != nil {
		return nil, err
	}

	// Validate address configuration: either use address directory (preferred) or legacy individual addresses
	addressDirectoryAddr := ctx.GlobalString(flags.AddressDirectoryFlag.Name)
	blsOperatorStateRetrieverAddr := ctx.GlobalString(flags.BlsOperatorStateRetrieverFlag.Name)
	eigenDAServiceManagerAddr := ctx.GlobalString(flags.EigenDAServiceManagerFlag.Name)
	if err := eth.ValidateAddressConfig(addressDirectoryAddr, blsOperatorStateRetrieverAddr, eigenDAServiceManagerAddr); err != nil {
		return nil, err
	}

	return &Config{
		EthClientConfig:               geth.ReadEthClientConfig(ctx),
		LoggerConfig:                  *loggerConfig,
		ChainStateConfig:              thegraph.ReadCLIConfig(ctx),
		AddressDirectoryAddr:          addressDirectoryAddr,
		BLSOperatorStateRetrieverAddr: blsOperatorStateRetrieverAddr,
		EigenDAServiceManagerAddr:     eigenDAServiceManagerAddr,
		PerPublicKeyRateLimit:         ctx.GlobalDuration(flags.PerPublicKeyRateLimit.Name),
		ChurnApprovalInterval:         ctx.GlobalDuration(flags.ChurnApprovalInterval.Name),
		MetricsConfig: MetricsConfig{
			HTTPPort:      ctx.GlobalString(flags.MetricsHTTPPort.Name),
			EnableMetrics: ctx.GlobalBool(flags.EnableMetrics.Name),
		},
	}, nil
}
