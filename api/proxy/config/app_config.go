package config

import (
	"fmt"
	"slices"

	"github.com/Layr-Labs/eigenda/api/proxy/common"
	enabled_apis "github.com/Layr-Labs/eigenda/api/proxy/config/enablement"
	"github.com/Layr-Labs/eigenda/api/proxy/config/v2/eigendaflags"
	"github.com/Layr-Labs/eigenda/api/proxy/metrics"
	"github.com/Layr-Labs/eigenda/api/proxy/servers/arbitrum_altda"
	"github.com/Layr-Labs/eigenda/api/proxy/servers/rest"
	"github.com/Layr-Labs/eigenda/api/proxy/store/builder"
	"github.com/urfave/cli/v2"
)

// AppConfig is the highest order config. Stores all relevant fields necessary for running
// REST ALTDA, Arbitrum Custom DA, & metrics servers.
type AppConfig struct {
	StoreBuilderConfig builder.Config
	SecretConfig       common.SecretConfigV2

	EnabledServersConfig *enabled_apis.EnabledServersConfig

	ArbCustomDASvrCfg arbitrum_altda.Config
	RestSvrCfg        rest.Config
	MetricsSvrConfig  metrics.Config
}

// Check checks critical config invariants and returns an error
// if there is a problem with the config struct's expression
func (c AppConfig) Check() error {
	err := c.StoreBuilderConfig.Check()
	if err != nil {
		return fmt.Errorf("check eigenDAConfig: %w", err)
	}

	v2Enabled := slices.Contains(c.StoreBuilderConfig.StoreConfig.BackendsToEnable, common.V2EigenDABackend)
	if v2Enabled && !c.StoreBuilderConfig.MemstoreEnabled {
		err = c.SecretConfig.Check()
		if err != nil {
			return fmt.Errorf("check secret config: %w", err)
		}
	}

	err = c.EnabledServersConfig.Check()
	if err != nil {
		return fmt.Errorf("check enabled APIs: %w", err)
	}

	return nil
}

func ReadAppConfig(ctx *cli.Context) (AppConfig, error) {
	storeBuilderConfig, err := builder.ReadConfig(ctx)
	if err != nil {
		return AppConfig{}, fmt.Errorf("read proxy config: %w", err)
	}

	return AppConfig{
		StoreBuilderConfig:   storeBuilderConfig,
		SecretConfig:         eigendaflags.ReadSecretConfigV2(ctx),
		EnabledServersConfig: enabled_apis.ReadEnabledServersCfg(ctx),

		ArbCustomDASvrCfg: arbitrum_altda.ReadConfig(ctx),
		MetricsSvrConfig:  metrics.ReadConfig(ctx),
	}, nil
}
