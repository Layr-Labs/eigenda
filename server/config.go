package server

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/Layr-Labs/eigenda-proxy/flags"
	"github.com/Layr-Labs/eigenda-proxy/flags/eigendaflags"
	"github.com/Layr-Labs/eigenda-proxy/store"
	"github.com/Layr-Labs/eigenda-proxy/store/generated_key/memstore"
	"github.com/Layr-Labs/eigenda-proxy/store/precomputed_key/redis"
	"github.com/Layr-Labs/eigenda-proxy/store/precomputed_key/s3"
	"github.com/Layr-Labs/eigenda-proxy/utils"
	"github.com/Layr-Labs/eigenda-proxy/verify"
	"github.com/Layr-Labs/eigenda/api/clients"

	opmetrics "github.com/ethereum-optimism/optimism/op-service/metrics"
)

type Config struct {
	EdaClientConfig clients.EigenDAClientConfig
	VerifierConfig  verify.Config

	MemstoreEnabled bool
	MemstoreConfig  memstore.Config

	// routing
	FallbackTargets []string
	CacheTargets    []string

	// secondary storage
	RedisConfig redis.Config
	S3Config    s3.Config
}

// ReadConfig ... parses the Config from the provided flags or environment variables.
func ReadConfig(ctx *cli.Context) Config {
	return Config{
		RedisConfig:     redis.ReadConfig(ctx),
		S3Config:        s3.ReadConfig(ctx),
		EdaClientConfig: eigendaflags.ReadConfig(ctx),
		VerifierConfig:  verify.ReadConfig(ctx),
		MemstoreEnabled: ctx.Bool(memstore.EnabledFlagName),
		MemstoreConfig:  memstore.ReadConfig(ctx),
		FallbackTargets: ctx.StringSlice(flags.FallbackTargetsFlagName),
		CacheTargets:    ctx.StringSlice(flags.CacheTargetsFlagName),
	}
}

// checkTargets ... verifies that a backend target slice is constructed correctly
func (cfg *Config) checkTargets(targets []string) error {
	if len(targets) == 0 {
		return nil
	}

	if utils.ContainsDuplicates(targets) {
		return fmt.Errorf("duplicate targets provided: %+v", targets)
	}

	for _, t := range targets {
		if store.StringToBackendType(t) == store.Unknown {
			return fmt.Errorf("unknown fallback target provided: %s", t)
		}
	}

	return nil
}

// Check ... verifies that configuration values are adequately set
func (cfg *Config) Check() error {
	if !cfg.MemstoreEnabled {
		if cfg.EdaClientConfig.RPC == "" {
			return fmt.Errorf("using eigenda backend (memstore.enabled=false) but eigenda disperser rpc url is not set")
		}
	}

	// cert verification is enabled
	// TODO: move this verification logic to verify/cli.go
	if cfg.VerifierConfig.VerifyCerts {
		if cfg.MemstoreEnabled {
			return fmt.Errorf("cannot enable cert verification when memstore is enabled")
		}
		if cfg.VerifierConfig.RPCURL == "" {
			return fmt.Errorf("cert verification enabled but eth rpc is not set")
		}
		if cfg.VerifierConfig.SvcManagerAddr == "" {
			return fmt.Errorf("cert verification enabled but svc manager address is not set")
		}
	}

	if cfg.S3Config.CredentialType == s3.CredentialTypeUnknown && cfg.S3Config.Endpoint != "" {
		return fmt.Errorf("s3 credential type must be set")
	}
	if cfg.S3Config.CredentialType == s3.CredentialTypeStatic {
		if cfg.S3Config.Endpoint != "" && (cfg.S3Config.AccessKeyID == "" || cfg.S3Config.AccessKeySecret == "") {
			return fmt.Errorf("s3 endpoint is set, but access key id or access key secret is not set")
		}
	}

	if cfg.RedisConfig.Endpoint == "" && cfg.RedisConfig.Password != "" {
		return fmt.Errorf("redis password is set, but endpoint is not")
	}

	err := cfg.checkTargets(cfg.FallbackTargets)
	if err != nil {
		return err
	}

	err = cfg.checkTargets(cfg.CacheTargets)
	if err != nil {
		return err
	}

	// verify that same target is not in both fallback and cache targets
	for _, t := range cfg.FallbackTargets {
		if utils.Contains(cfg.CacheTargets, t) {
			return fmt.Errorf("target %s is in both fallback and cache targets", t)
		}
	}

	return nil
}

type CLIConfig struct {
	EigenDAConfig Config
	MetricsCfg    opmetrics.CLIConfig
}

func ReadCLIConfig(ctx *cli.Context) CLIConfig {
	config := ReadConfig(ctx)
	return CLIConfig{
		EigenDAConfig: config,
		MetricsCfg:    opmetrics.ReadCLIConfig(ctx),
	}
}

func (c CLIConfig) Check() error {
	err := c.EigenDAConfig.Check()
	if err != nil {
		return err
	}
	return nil
}
