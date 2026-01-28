package builder

import (
	"encoding/json"
	"fmt"
	"slices"
	"time"

	"github.com/Layr-Labs/eigenda/api/proxy/common"
	eigendaflags_v2 "github.com/Layr-Labs/eigenda/api/proxy/config/v2/eigendaflags"
	"github.com/Layr-Labs/eigenda/api/proxy/store"
	"github.com/Layr-Labs/eigenda/api/proxy/store/generated_key/memstore"
	"github.com/Layr-Labs/eigenda/api/proxy/store/generated_key/memstore/memconfig"
	"github.com/Layr-Labs/eigenda/api/proxy/store/secondary/s3"
	"github.com/urfave/cli/v2"
)

// Config ... Higher order config which bundles all configs for building
// the proxy store manager with necessary client context
type Config struct {
	StoreConfig store.Config

	// main storage configs
	ClientConfigV2 common.ClientConfigV2

	MemstoreConfig  *memconfig.SafeConfig
	MemstoreEnabled bool

	// secondary storage cfgs
	S3Config s3.Config

	// eth rpc retry count and delay
	RetryCount int
	RetryDelay time.Duration
}

// ReadConfig ... parses the Config from the provided flags or environment variables.
func ReadConfig(ctx *cli.Context) (Config, error) {
	storeConfig, err := store.ReadConfig(ctx)
	if err != nil {
		return Config{}, fmt.Errorf("read storage config: %w", err)
	}

	if slices.Contains(storeConfig.BackendsToEnable, common.V1EigenDABackend) {
		return Config{}, fmt.Errorf("V1 backend has been removed, please use V2")
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
		return Config{}, fmt.Errorf("V1 dispersal backend has been removed, please use V2")
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
		StoreConfig:     storeConfig,
		ClientConfigV2:  clientConfigV2,
		MemstoreConfig:  memstoreConfig,
		MemstoreEnabled: ctx.Bool(memstore.EnabledFlagName),
		S3Config:        s3.ReadConfig(ctx),
		RetryCount:      ctx.Int(eigendaflags_v2.EthRPCRetryCountFlagName),
		RetryDelay:      ctx.Duration(eigendaflags_v2.EthRPCRetryDelayIncrementFlagName),
	}

	return cfg, nil
}

// Check ... verifies that configuration values are adequately set
func (cfg *Config) Check() error {
	v1Enabled := slices.Contains(cfg.StoreConfig.BackendsToEnable, common.V1EigenDABackend)
	if v1Enabled {
		return fmt.Errorf("V1 backend has been removed, please use V2")
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

	return cfg.StoreConfig.Check()
}

func (cfg *Config) ToString() (string, error) {
	redacted := "******"

	// create a copy, otherwise the original values being redacted will be lost
	configCopy := *cfg

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
