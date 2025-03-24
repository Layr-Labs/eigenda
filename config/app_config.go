package config

import (
	"fmt"

	"github.com/Layr-Labs/eigenda-proxy/common"
	"github.com/Layr-Labs/eigenda-proxy/config/v2/eigendaflags"
	"github.com/Layr-Labs/eigenda-proxy/metrics"
	"github.com/urfave/cli/v2"
)

// AppConfig ... Highest order config. Stores all relevant fields necessary for running both proxy & metrics servers.
type AppConfig struct {
	EigenDAConfig ProxyConfig
	SecretConfig  common.SecretConfigV2
	MetricsConfig metrics.Config
}

// Check checks config invariants, and returns an error if there is a problem with the config struct
func (c AppConfig) Check() error {
	err := c.EigenDAConfig.Check()
	if err != nil {
		return fmt.Errorf("check eigenDAConfig: %w", err)
	}

	if c.EigenDAConfig.ClientConfigV2.DisperseToV2 && !c.EigenDAConfig.MemstoreEnabled {
		err = c.SecretConfig.Check()
		if err != nil {
			return fmt.Errorf("check secret config: %w", err)
		}
	}

	return nil
}

func ReadCLIConfig(ctx *cli.Context) (AppConfig, error) {
	proxyConfig, err := ReadProxyConfig(ctx)
	if err != nil {
		return AppConfig{}, fmt.Errorf("read proxy config: %w", err)
	}

	secretConfig := eigendaflags.ReadSecretConfigV2(ctx)

	return AppConfig{
		EigenDAConfig: proxyConfig,
		SecretConfig:  secretConfig,
		MetricsConfig: metrics.ReadConfig(ctx),
	}, nil
}
