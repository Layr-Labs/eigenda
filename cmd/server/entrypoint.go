package main

import (
	"context"
	"fmt"

	"github.com/Layr-Labs/eigenda-proxy/config"
	proxy_logging "github.com/Layr-Labs/eigenda-proxy/logging"
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
	configString, err := cfg.EigenDAConfig.ToString()
	if err != nil {
		return fmt.Errorf("convert config json to string: %w", err)
	}

	log.Infof("Initializing EigenDA proxy server with config (\"*****\" fields are hidden): %v", configString)

	metrics := proxy_metrics.NewMetrics("default")

	ctx, ctxCancel := context.WithCancel(cliCtx.Context)
	defer ctxCancel()

	proxyServer, err := server.BuildAndStartProxyServer(ctx, log, metrics, cfg)
	if err != nil {
		return fmt.Errorf("build proxy server: %w", err)
	}

	defer func() {
		if err := proxyServer.Stop(); err != nil {
			log.Error("failed to stop DA server", "err", err)
		}

		log.Info("successfully shutdown API server")
	}()

	if cfg.MetricsConfig.Enabled {
		log.Debug("starting metrics server", "addr", cfg.MetricsConfig.Host, "port", cfg.MetricsConfig.Port)
		svr, err := metrics.StartServer(cfg.MetricsConfig.Host, cfg.MetricsConfig.Port)
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
