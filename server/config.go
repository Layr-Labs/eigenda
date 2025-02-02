package server

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/Layr-Labs/eigenda-proxy/flags/eigendaflags"
	"github.com/Layr-Labs/eigenda-proxy/metrics"
	"github.com/Layr-Labs/eigenda-proxy/store"
	"github.com/Layr-Labs/eigenda-proxy/store/generated_key/memstore"
	"github.com/Layr-Labs/eigenda-proxy/verify"
	"github.com/Layr-Labs/eigenda/api/clients"
)

type Config struct {
	EdaClientConfig clients.EigenDAClientConfig
	MemstoreConfig  memstore.Config
	StorageConfig   store.Config
	VerifierConfig  verify.Config
	PutRetries      uint

	MemstoreEnabled bool
}

// ReadConfig ... parses the Config from the provided flags or environment variables.
func ReadConfig(ctx *cli.Context) Config {
	edaClientConfig := eigendaflags.ReadConfig(ctx)
	return Config{
		EdaClientConfig: edaClientConfig,
		VerifierConfig:  verify.ReadConfig(ctx, edaClientConfig),
		PutRetries:      ctx.Uint(eigendaflags.PutRetriesFlagName),
		MemstoreEnabled: ctx.Bool(memstore.EnabledFlagName),
		MemstoreConfig:  memstore.ReadConfig(ctx),
		StorageConfig:   store.ReadConfig(ctx),
	}
}

// Check ... verifies that configuration values are adequately set
func (cfg *Config) Check() error {
	if !cfg.MemstoreEnabled {
		if cfg.EdaClientConfig.RPC == "" {
			return fmt.Errorf("using eigenda backend (memstore.enabled=false) but eigenda disperser rpc url is not set")
		}
	}

	// provide dummy values to eigenda client config. Since the client won't be called in this
	// mode it doesn't matter.
	if cfg.MemstoreEnabled {
		cfg.EdaClientConfig.SvcManagerAddr = "0x0000000000000000000000000000000000000000"
		cfg.EdaClientConfig.EthRpcUrl = "http://0.0.0.0:666"
	}

	if !cfg.MemstoreEnabled {
		if cfg.EdaClientConfig.SvcManagerAddr == "" {
			return fmt.Errorf("service manager address is required for communication with EigenDA")
		}
		if cfg.EdaClientConfig.EthRpcUrl == "" {
			return fmt.Errorf("eth prc url is required for communication with EigenDA")
		}
	}

	// cert verification is enabled
	// TODO: move this verification logic to verify/cli.go
	if cfg.VerifierConfig.VerifyCerts {
		if cfg.MemstoreEnabled {
			return fmt.Errorf("cannot enable cert verification when memstore is enabled. use --%s", verify.CertVerificationDisabledFlagName)
		}
		if cfg.VerifierConfig.RPCURL == "" {
			return fmt.Errorf("cert verification enabled but eth rpc is not set")
		}
		if cfg.VerifierConfig.SvcManagerAddr == "" {
			return fmt.Errorf("cert verification enabled but svc manager address is not set")
		}
	}

	return cfg.StorageConfig.Check()
}

type CLIConfig struct {
	EigenDAConfig Config
	MetricsCfg    metrics.CLIConfig
}

func ReadCLIConfig(ctx *cli.Context) CLIConfig {
	return CLIConfig{
		EigenDAConfig: ReadConfig(ctx),
		MetricsCfg:    metrics.ReadCLIConfig(ctx),
	}
}

func (c CLIConfig) Check() error {
	err := c.EigenDAConfig.Check()
	if err != nil {
		return err
	}
	return nil
}
