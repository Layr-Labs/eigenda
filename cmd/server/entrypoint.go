package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Layr-Labs/eigenda-proxy/common"
	"github.com/Layr-Labs/eigenda-proxy/config"
	eigendaflags_v2 "github.com/Layr-Labs/eigenda-proxy/config/eigendaflags/v2"

	proxy_logging "github.com/Layr-Labs/eigenda-proxy/logging"
	"github.com/Layr-Labs/eigenda-proxy/store"
	"github.com/Layr-Labs/eigenda-proxy/store/generated_key/memstore/memconfig"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/gorilla/mux"

	proxy_metrics "github.com/Layr-Labs/eigenda-proxy/metrics"
	"github.com/Layr-Labs/eigenda-proxy/server"
	"github.com/urfave/cli/v2"

	"github.com/ethereum-optimism/optimism/op-service/ctxinterrupt"
)

func StartProxySvr(cliCtx *cli.Context) error {
	logCfg, err := proxy_logging.ReadLoggerCLIConfig(cliCtx)
	if err != nil {
		return err
	}

	log, err := proxy_logging.NewLogger(*logCfg)
	if err != nil {
		return err
	}

	log.Info("Starting EigenDA Proxy Server", "version", Version, "date", Date, "commit", Commit)

	cfg, err := config.ReadCLIConfig(cliCtx)
	if err != nil {
		return fmt.Errorf("read cli config: %w", err)
	}

	if err := cfg.Check(); err != nil {
		return err
	}
	err = prettyPrintConfig(cliCtx, log)
	if err != nil {
		return fmt.Errorf("failed to pretty print config: %w", err)
	}

	var secretConfig common.SecretConfigV2
	if cfg.EigenDAConfig.EdaClientConfigV2.Enabled {
		// secret config is kept entirely separate from the other config values, which may be printed
		secretConfig = eigendaflags_v2.ReadSecretConfigV2(cliCtx)
		if err := secretConfig.Check(); err != nil {
			return err
		}
	}

	metrics := proxy_metrics.NewMetrics("default")

	ctx, ctxCancel := context.WithCancel(cliCtx.Context)
	defer ctxCancel()

	memConfig := cfg.EigenDAConfig.MemstoreConfig
	if !cfg.EigenDAConfig.MemstoreEnabled {
		memConfig = nil
	}

	storageManager, err := store.NewStorageManagerBuilder(
		ctx,
		log,
		metrics,
		cfg.EigenDAConfig.StorageConfig,
		cfg.EigenDAConfig.EdaVerifierConfigV1,
		cfg.EigenDAConfig.EdaClientConfigV1,
		cfg.EigenDAConfig.EdaClientConfigV2,
		secretConfig,
		memConfig,
		cfg.EigenDAConfig.EdaClientConfigV2.PutRetries,
		cfg.EigenDAConfig.MaxBlobSizeBytes,
	).Build(ctx)
	if err != nil {
		return fmt.Errorf("failed to create store: %w", err)
	}

	server := server.NewServer(cfg.EigenDAConfig.ServerConfig, storageManager, log, metrics)
	r := mux.NewRouter()
	server.RegisterRoutes(r)
	if cfg.EigenDAConfig.MemstoreEnabled {
		memconfig.NewHandlerHTTP(log, cfg.EigenDAConfig.MemstoreConfig).RegisterMemstoreConfigHandlers(r)
	}

	if err := server.Start(r); err != nil {
		return fmt.Errorf("failed to start the DA server: %w", err)
	}

	log.Info("Started EigenDA proxy server")

	defer func() {
		if err := server.Stop(); err != nil {
			log.Error("failed to stop DA server", "err", err)
		}

		log.Info("successfully shutdown API server")
	}()

	if cfg.MetricsCfg.Enabled {
		log.Debug("starting metrics server", "addr", cfg.MetricsCfg.Host, "port", cfg.MetricsCfg.Port)
		svr, err := metrics.StartServer(cfg.MetricsCfg.Host, cfg.MetricsCfg.Port)
		if err != nil {
			return fmt.Errorf("failed to start metrics server: %w", err)
		}
		defer func() {
			if err := svr.Stop(context.Background()); err != nil {
				log.Error("failed to stop metrics server", "err", err)
			}
		}()
		log.Info("started metrics server", "addr", svr.Addr())
		metrics.RecordUp()
	}

	return ctxinterrupt.Wait(cliCtx.Context)
}

// TODO: we should probably just change EdaClientConfig struct definition in eigenda-client
func prettyPrintConfig(cliCtx *cli.Context, log logging.Logger) error {
	redacted := "******"

	// we read a new config which we modify to hide private info in order to log the rest
	cfg, err := config.ReadCLIConfig(cliCtx)
	if err != nil {
		return fmt.Errorf("read cli config: %w", err)
	}
	if cfg.EigenDAConfig.EdaClientConfigV1.SignerPrivateKeyHex != "" {
		// marshaling defined in client config
		cfg.EigenDAConfig.EdaClientConfigV1.SignerPrivateKeyHex = redacted
	}
	if cfg.EigenDAConfig.EdaClientConfigV1.EthRpcUrl != "" {
		// hiding as RPC providers typically use sensitive API keys within
		cfg.EigenDAConfig.EdaClientConfigV1.EthRpcUrl = redacted
	}
	if cfg.EigenDAConfig.StorageConfig.RedisConfig.Password != "" {
		cfg.EigenDAConfig.StorageConfig.RedisConfig.Password = redacted // masking Redis password
	}

	configJSON, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	log.Info(
		fmt.Sprintf(
			"Initializing EigenDA proxy server with config (\"*****\" fields are hidden): %v",
			string(configJSON)))

	return nil
}
