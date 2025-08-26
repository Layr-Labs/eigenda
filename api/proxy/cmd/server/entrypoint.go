package main

import (
	"context"
	"fmt"

	"github.com/Layr-Labs/eigenda/api/proxy/config"
	proxy_logging "github.com/Layr-Labs/eigenda/api/proxy/logging"
	proxy_metrics "github.com/Layr-Labs/eigenda/api/proxy/metrics"
	"github.com/Layr-Labs/eigenda/api/proxy/server"
	"github.com/Layr-Labs/eigenda/api/proxy/store/builder"
	"github.com/Layr-Labs/eigenda/api/proxy/store/generated_key/memstore/memconfig"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
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

	cfg, err := config.ReadAppConfig(cliCtx)
	if err != nil {
		return fmt.Errorf("read cli config: %w", err)
	}

	if err := cfg.Check(); err != nil {
		return err
	}
	configString, err := cfg.StoreBuilderConfig.ToString()
	if err != nil {
		return fmt.Errorf("convert config json to string: %w", err)
	}

	log.Infof("Initializing EigenDA proxy server with config (\"*****\" fields are hidden): %v", configString)

	registry := prometheus.NewRegistry()
	metrics := proxy_metrics.NewMetrics(registry)

	ctx, ctxCancel := context.WithCancel(cliCtx.Context)
	defer ctxCancel()

	certMgr, keccakMgr, err := builder.BuildManagers(
		ctx,
		log,
		metrics,
		cfg.StoreBuilderConfig,
		cfg.SecretConfig,
		registry,
	)
	if err != nil {
		return fmt.Errorf("build storage managers: %w", err)
	}

	proxyServer := server.NewServer(cfg.ServerConfig, certMgr, keccakMgr, log, metrics)
	router := mux.NewRouter()
	proxyServer.RegisterRoutes(router)
	if cfg.StoreBuilderConfig.MemstoreEnabled {
		memconfig.NewHandlerHTTP(log, cfg.StoreBuilderConfig.MemstoreConfig).RegisterMemstoreConfigHandlers(router)
	}

	if err := proxyServer.Start(router); err != nil {
		return fmt.Errorf("start proxy server: %w", err)
	}

	log.Info("Started EigenDA proxy server")

	defer func() {
		if err := proxyServer.Stop(); err != nil {
			log.Error("failed to stop DA server", "err", err)
		}

		log.Info("Successfully shutdown API server")
	}()

	if cfg.MetricsServerConfig.Enabled {
		log.Info("Starting metrics server", "addr", cfg.MetricsServerConfig.Host, "port", cfg.MetricsServerConfig.Port)
		svr := proxy_metrics.NewServer(registry, cfg.MetricsServerConfig)
		err := svr.Start()
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
