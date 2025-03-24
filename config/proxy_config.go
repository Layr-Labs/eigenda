package config

import (
	"encoding/json"
	"fmt"

	"github.com/Layr-Labs/eigenda-proxy/common"
	"github.com/Layr-Labs/eigenda-proxy/config/eigendaflags"
	eigendaflags_v2 "github.com/Layr-Labs/eigenda-proxy/config/v2/eigendaflags"
	"github.com/Layr-Labs/eigenda-proxy/store"
	"github.com/Layr-Labs/eigenda-proxy/store/generated_key/memstore"
	"github.com/Layr-Labs/eigenda-proxy/store/generated_key/memstore/memconfig"
	"github.com/Layr-Labs/eigenda-proxy/verify"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/urfave/cli/v2"
)

// ProxyConfig ... Higher order config which bundles all configs for orchestrating
// the proxy server with necessary client context
type ProxyConfig struct {
	ServerConfig     ServerConfig
	ClientConfigV1   common.ClientConfigV1
	VerifierConfigV1 verify.Config
	KzgConfig        kzg.KzgConfig
	ClientConfigV2   common.ClientConfigV2

	MemstoreConfig  *memconfig.SafeConfig
	MemstoreEnabled bool
	StorageConfig   store.Config
}

// ReadProxyConfig ... parses the Config from the provided flags or environment variables.
func ReadProxyConfig(ctx *cli.Context) (ProxyConfig, error) {
	clientConfigV1, err := eigendaflags.ReadClientConfigV1(ctx)
	if err != nil {
		return ProxyConfig{}, fmt.Errorf("read client config v1: %w", err)
	}

	verifierConfigV1 := verify.ReadConfig(ctx, clientConfigV1)

	clientConfigV2, err := eigendaflags_v2.ReadClientConfigV2(ctx)
	if err != nil {
		return ProxyConfig{}, fmt.Errorf("read client config v2: %w", err)
	}

	var maxBlobSizeBytes uint64
	if clientConfigV2.DisperseToV2 {
		maxBlobSizeBytes = clientConfigV2.MaxBlobSizeBytes
	} else {
		maxBlobSizeBytes = clientConfigV1.MaxBlobSizeBytes
	}

	kzgConfig := verify.ReadKzgConfig(ctx, maxBlobSizeBytes)

	memstoreConfig, err := memstore.ReadConfig(ctx, maxBlobSizeBytes)
	if err != nil {
		return ProxyConfig{}, fmt.Errorf("read memstore config: %w", err)
	}

	cfg := ProxyConfig{
		ServerConfig: ServerConfig{
			DisperseToV2: clientConfigV2.DisperseToV2,
			Host:         ctx.String(ListenAddrFlagName),
			Port:         ctx.Int(PortFlagName),
		},
		ClientConfigV1:   clientConfigV1,
		VerifierConfigV1: verifierConfigV1,
		KzgConfig:        kzgConfig,
		ClientConfigV2:   clientConfigV2,
		MemstoreConfig:   memstoreConfig,
		MemstoreEnabled:  ctx.Bool(memstore.EnabledFlagName),
		StorageConfig:    store.ReadConfig(ctx),
	}

	return cfg, nil
}

// Check ... verifies that configuration values are adequately set
func (cfg *ProxyConfig) Check() error {
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

	// V2 dispersal/retrieval enabled
	if cfg.ClientConfigV2.DisperseToV2 && !cfg.MemstoreEnabled {
		err := cfg.ClientConfigV2.Check()
		if err != nil {
			return err
		}
	}

	return cfg.StorageConfig.Check()
}

func (cfg *ProxyConfig) ToString() (string, error) {
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
	if configCopy.StorageConfig.RedisConfig.Password != "" {
		configCopy.StorageConfig.RedisConfig.Password = redacted
	}
	if configCopy.StorageConfig.S3Config.AccessKeySecret != "" {
		configCopy.StorageConfig.S3Config.AccessKeySecret = redacted
	}
	if configCopy.StorageConfig.S3Config.AccessKeyID != "" {
		configCopy.StorageConfig.S3Config.AccessKeyID = redacted
	}

	configJSON, err := json.MarshalIndent(configCopy, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal config: %w", err)
	}

	return string(configJSON), nil
}
