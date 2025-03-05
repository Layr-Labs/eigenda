package config

import (
	"fmt"

	"github.com/Layr-Labs/eigenda-proxy/common"
	"github.com/Layr-Labs/eigenda-proxy/config/eigendaflags"
	eigendaflags_v2 "github.com/Layr-Labs/eigenda-proxy/config/eigendaflags/v2"
	"github.com/Layr-Labs/eigenda-proxy/server"
	"github.com/Layr-Labs/eigenda-proxy/store"
	"github.com/Layr-Labs/eigenda-proxy/store/generated_key/memstore"
	"github.com/Layr-Labs/eigenda-proxy/store/generated_key/memstore/memconfig"
	"github.com/Layr-Labs/eigenda-proxy/verify/v1"
	"github.com/urfave/cli/v2"

	"github.com/Layr-Labs/eigenda/api/clients"
)

// ProxyConfig ... Higher order config which bundles all configs for orchestrating
// the proxy server with necessary client context
type ProxyConfig struct {
	ServerConfig        server.Config
	EdaClientConfigV1   clients.EigenDAClientConfig
	EdaVerifierConfigV1 verify.Config

	EdaClientConfigV2 common.ClientConfigV2

	MemstoreConfig *memconfig.SafeConfig
	StorageConfig  store.Config

	// Enabling will turn on memstore for both V1 && V2
	//
	// Currently the memstore is not persisted to disk, and is lost when proxy is stopped. Testing migrations is
	// currently not supported because it requires restarting the proxy with new flags, which will cause the
	// memstore db to be wiped.
	//
	MemstoreEnabled bool

	EigenDAV2Enabled bool

	PutRetries       uint
	MaxBlobSizeBytes uint
}

// ReadProxyConfig ... parses the Config from the provided flags or environment variables.
func ReadProxyConfig(ctx *cli.Context) (ProxyConfig, error) {
	edaClientV1Config := eigendaflags.ReadConfig(ctx)
	edaClientV2Config, err := eigendaflags_v2.ReadClientConfigV2(ctx)
	if err != nil {
		return ProxyConfig{}, fmt.Errorf("read client config v2: %w", err)
	}

	cfg := ProxyConfig{
		ServerConfig: server.Config{
			DisperseV2: edaClientV2Config.Enabled,
			Host:       ctx.String(ListenAddrFlagName),
			Port:       ctx.Int(PortFlagName),
		},
		EdaClientConfigV1:   edaClientV1Config,
		EdaClientConfigV2:   edaClientV2Config,
		EdaVerifierConfigV1: verify.ReadConfig(ctx, edaClientV1Config),
		PutRetries:          ctx.Uint(eigendaflags.PutRetriesFlagName),
		MemstoreEnabled:     ctx.Bool(memstore.EnabledFlagName),
		MemstoreConfig:      memstore.ReadConfig(ctx),
		StorageConfig:       store.ReadConfig(ctx),
		EigenDAV2Enabled:    edaClientV2Config.Enabled,
		MaxBlobSizeBytes:    uint(verify.MaxBlobLengthBytes),
	}

	return cfg, nil
}

// Check ... verifies that configuration values are adequately set
func (cfg *ProxyConfig) Check() error {
	if cfg.MemstoreEnabled {
		// provide dummy values to eigenda client config. Since the client won't be called in this
		// mode it doesn't matter.
		cfg.EdaClientConfigV1.SvcManagerAddr = "0x0000000000000000000000000000000000000000"
		cfg.EdaClientConfigV1.EthRpcUrl = "http://0.0.0.0:666"
	} else {
		if cfg.EdaClientConfigV1.SvcManagerAddr == "" {
			return fmt.Errorf("service manager address is required for communication with EigenDA")
		}
		if cfg.EdaClientConfigV1.EthRpcUrl == "" {
			return fmt.Errorf("eth prc url is required for communication with EigenDA")
		}
		if cfg.EdaClientConfigV1.RPC == "" {
			return fmt.Errorf("using eigenda backend (memstore.enabled=false) but eigenda disperser rpc url is not set")
		}
	}

	// cert verification is enabled
	// TODO: move this verification logic to verify/cli.go
	if cfg.EdaVerifierConfigV1.VerifyCerts {
		if cfg.MemstoreEnabled {
			return fmt.Errorf(
				"cannot enable cert verification when memstore is enabled. use --%s",
				verify.CertVerificationDisabledFlagName)
		}
		if cfg.EdaVerifierConfigV1.RPCURL == "" {
			return fmt.Errorf("cert verification enabled but eth rpc is not set")
		}
		if cfg.EdaVerifierConfigV1.SvcManagerAddr == "" {
			return fmt.Errorf("cert verification enabled but svc manager address is not set")
		}
	}

	// V2 dispersal/retrieval enabled
	if cfg.EigenDAV2Enabled && !cfg.MemstoreEnabled {
		err := cfg.EdaClientConfigV2.Check()
		if err != nil {
			return err
		}
	}

	return cfg.StorageConfig.Check()
}
