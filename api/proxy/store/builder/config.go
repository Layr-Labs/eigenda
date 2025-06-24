package builder

import (
	"encoding/json"
	"fmt"
	"slices"

	"github.com/Layr-Labs/eigenda/api/proxy/common"
	"github.com/Layr-Labs/eigenda/api/proxy/config/eigendaflags"
	eigendaflags_v2 "github.com/Layr-Labs/eigenda/api/proxy/config/v2/eigendaflags"
	"github.com/Layr-Labs/eigenda/api/proxy/store"
	"github.com/Layr-Labs/eigenda/api/proxy/store/generated_key/eigenda/verify"
	"github.com/Layr-Labs/eigenda/api/proxy/store/generated_key/memstore"
	"github.com/Layr-Labs/eigenda/api/proxy/store/generated_key/memstore/memconfig"
	"github.com/Layr-Labs/eigenda/api/proxy/store/secondary/redis"
	"github.com/Layr-Labs/eigenda/api/proxy/store/secondary/s3"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/urfave/cli/v2"
)

// Config ... Higher order config which bundles all configs for building
// the proxy store manager with necessary client context
type Config struct {
	StoreConfig store.Config

	// main storage configs
	ClientConfigV1   common.ClientConfigV1
	VerifierConfigV1 verify.Config
	KzgConfig        kzg.KzgConfig
	ClientConfigV2   common.ClientConfigV2

	MemstoreConfig  *memconfig.SafeConfig
	MemstoreEnabled bool

	// secondary storage cfgs
	RedisConfig redis.Config
	S3Config    s3.Config
}

// ReadConfig ... parses the Config from the provided flags or environment variables.
func ReadConfig(ctx *cli.Context) (Config, error) {
	storeConfig, err := store.ReadConfig(ctx)
	if err != nil {
		return Config{}, fmt.Errorf("read storage config: %w", err)
	}

	var clientConfigV1 common.ClientConfigV1
	var verifierConfigV1 verify.Config
	if slices.Contains(storeConfig.BackendsToEnable, common.V1EigenDABackend) {
		clientConfigV1, err = eigendaflags.ReadClientConfigV1(ctx)
		if err != nil {
			return Config{}, fmt.Errorf("read client config v1: %w", err)
		}

		verifierConfigV1 = verify.ReadConfig(ctx, clientConfigV1)
	}

	var clientConfigV2 common.ClientConfigV2
	if slices.Contains(storeConfig.BackendsToEnable, common.V2EigenDABackend) {
		clientConfigV2, err = eigendaflags_v2.ReadClientConfigV2(ctx)
		if err != nil {
			return Config{}, fmt.Errorf("read client config v2: %w", err)
		}
	}

	var maxBlobSizeBytes uint64
	switch storeConfig.DispersalBackend {
	case common.V1EigenDABackend:
		maxBlobSizeBytes = clientConfigV1.MaxBlobSizeBytes
	case common.V2EigenDABackend:
		maxBlobSizeBytes = clientConfigV2.MaxBlobSizeBytes
	default:
		return Config{}, fmt.Errorf("unknown dispersal backend %s",
			common.EigenDABackendToString(storeConfig.DispersalBackend))
	}

	memstoreConfig, err := memstore.ReadConfig(ctx, maxBlobSizeBytes)
	if err != nil {
		return Config{}, fmt.Errorf("read memstore config: %w", err)
	}

	cfg := Config{
		StoreConfig:      storeConfig,
		ClientConfigV1:   clientConfigV1,
		VerifierConfigV1: verifierConfigV1,
		KzgConfig:        verify.ReadKzgConfig(ctx, maxBlobSizeBytes),
		ClientConfigV2:   clientConfigV2,
		MemstoreConfig:   memstoreConfig,
		MemstoreEnabled:  ctx.Bool(memstore.EnabledFlagName),
		RedisConfig:      redis.ReadConfig(ctx),
		S3Config:         s3.ReadConfig(ctx),
	}

	return cfg, nil
}

// Check ... verifies that configuration values are adequately set
func (cfg *Config) Check() error {
	v1Enabled := slices.Contains(cfg.StoreConfig.BackendsToEnable, common.V1EigenDABackend)
	if v1Enabled {
		err := cfg.checkV1Config()
		if err != nil {
			return fmt.Errorf("check v1 config: %w", err)
		}
	}

	v2Enabled := slices.Contains(cfg.StoreConfig.BackendsToEnable, common.V2EigenDABackend)
	if v2Enabled && !cfg.MemstoreEnabled {
		err := cfg.ClientConfigV2.Check()
		if err != nil {
			return fmt.Errorf("check v2 config: %w", err)
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

	return cfg.StoreConfig.Check()
}

func (cfg *Config) checkV1Config() error {
	if cfg.MemstoreEnabled {
		// provide dummy values to eigenda client config. Since the client won't be called in this
		// mode it doesn't matter.
		cfg.VerifierConfigV1.SvcManagerAddr = "0x0000000000000000000000000000000000000000"
		cfg.ClientConfigV1.EdaClientCfg.EthRpcUrl = "http://0.0.0.0:666"
	} else {
		if cfg.ClientConfigV1.EdaClientCfg.SvcManagerAddr == "" || cfg.VerifierConfigV1.SvcManagerAddr == "" {
			return fmt.Errorf("service manager address is required for communication with EigenDA")
		}
		if cfg.ClientConfigV1.EdaClientCfg.EthRpcUrl == "" {
			return fmt.Errorf("eth prc url is required for communication with EigenDA")
		}
		if cfg.ClientConfigV1.EdaClientCfg.RPC == "" {
			return fmt.Errorf("using eigenda backend (memstore.enabled=false) but eigenda disperser rpc url is not set")
		}
	}

	// cert verification is enabled
	// TODO: move this verification logic to verify/cli.go
	if cfg.VerifierConfigV1.VerifyCerts {
		if cfg.MemstoreEnabled {
			return fmt.Errorf(
				"cannot enable cert verification when memstore is enabled. use --%s",
				verify.CertVerificationDisabledFlagName)
		}
		if cfg.VerifierConfigV1.RPCURL == "" {
			return fmt.Errorf("cert verification enabled but eth rpc is not set")
		}
		if cfg.ClientConfigV1.EdaClientCfg.SvcManagerAddr == "" || cfg.VerifierConfigV1.SvcManagerAddr == "" {
			return fmt.Errorf("cert verification enabled but svc manager address is not set")
		}
	}

	return nil
}

func (cfg *Config) ToString() (string, error) {
	redacted := "******"

	// create a copy, otherwise the original values being redacted will be lost
	configCopy := *cfg

	if configCopy.ClientConfigV1.EdaClientCfg.SignerPrivateKeyHex != "" {
		configCopy.ClientConfigV1.EdaClientCfg.SignerPrivateKeyHex = redacted
	}
	if configCopy.ClientConfigV1.EdaClientCfg.EthRpcUrl != "" {
		// hiding as RPC providers typically use sensitive API keys within
		configCopy.ClientConfigV1.EdaClientCfg.EthRpcUrl = redacted
	}
	if configCopy.RedisConfig.Password != "" {
		configCopy.RedisConfig.Password = redacted
	}
	if configCopy.S3Config.AccessKeySecret != "" {
		configCopy.S3Config.AccessKeySecret = redacted
	}
	if configCopy.S3Config.AccessKeyID != "" {
		configCopy.S3Config.AccessKeyID = redacted
	}

	configJSON, err := json.MarshalIndent(configCopy, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal config: %w", err)
	}

	return string(configJSON), nil
}
