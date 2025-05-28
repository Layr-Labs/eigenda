package config

import (
	"fmt"
	"slices"

	"github.com/Layr-Labs/eigenda-proxy/common"
	"github.com/Layr-Labs/eigenda-proxy/config/v2/eigendaflags"
	"github.com/Layr-Labs/eigenda-proxy/metrics"
	"github.com/Layr-Labs/eigenda-proxy/server"
	"github.com/Layr-Labs/eigenda-proxy/store/builder"
	"github.com/urfave/cli/v2"
)

// AppConfig ... Highest order config. Stores all relevant fields necessary for running both proxy & metrics servers.
type AppConfig struct {
	StoreBuilderConfig  builder.Config
	SecretConfig        common.SecretConfigV2
	ServerConfig        server.Config
	MetricsServerConfig metrics.Config
}

// Check checks config invariants, and returns an error if there is a problem with the config struct
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

	return nil
}

func ReadAppConfig(ctx *cli.Context) (AppConfig, error) {
	storeBuilderConfig, err := builder.ReadConfig(ctx)
	if err != nil {
		return AppConfig{}, fmt.Errorf("read proxy config: %w", err)
	}

	return AppConfig{
		StoreBuilderConfig:  storeBuilderConfig,
		SecretConfig:        eigendaflags.ReadSecretConfigV2(ctx),
		ServerConfig:        server.ReadConfig(ctx),
		MetricsServerConfig: metrics.ReadConfig(ctx),
	}, nil
}
