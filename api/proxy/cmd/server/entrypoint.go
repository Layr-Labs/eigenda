package main

import (
	"context"
	"fmt"

	"github.com/Layr-Labs/eigenda/api/proxy/config"
	enabled_apis "github.com/Layr-Labs/eigenda/api/proxy/config/enablement"
	proxy_logging "github.com/Layr-Labs/eigenda/api/proxy/logging"
	proxy_metrics "github.com/Layr-Labs/eigenda/api/proxy/metrics"
	"github.com/Layr-Labs/eigenda/api/proxy/servers/arbitrum_altda"
	"github.com/Layr-Labs/eigenda/api/proxy/servers/rest"
	"github.com/Layr-Labs/eigenda/api/proxy/store/builder"
	"github.com/Layr-Labs/eigenda/api/proxy/store/generated_key/memstore/memconfig"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/urfave/cli/v2"

	"github.com/ethereum-optimism/optimism/op-service/ctxinterrupt"
)

// TODO: Explore better encapsulation patterns that binds common interfaces / usage patterns
// across the three servers (arb-altda, rest, metrics) that can be spun-up under the proxy service.
// Especially if there's ever a need for an additional stack specific ALT DA server type to be introduced.
func StartProxyService(cliCtx *cli.Context) error {
	logCfg, err := proxy_logging.ReadLoggerCLIConfig(cliCtx)
	if err != nil {
		return err
	}

	log, err := proxy_logging.NewLogger(*logCfg)
	if err != nil {
		return err
	}

	log.Info("Starting EigenDA Proxy Service", "version", Version, "date", Date, "commit", Commit)

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

	log.Infof("Initializing EigenDA proxy service with config (\"*****\" fields are hidden): %v", configString)

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

	restServer := rest.NewServer(cfg.RestSvrCfg, certMgr, keccakMgr, log, metrics)
	router := mux.NewRouter()
	restServer.RegisterRoutes(router)
	if cfg.StoreBuilderConfig.MemstoreEnabled {
		memconfig.NewHandlerHTTP(log, cfg.StoreBuilderConfig.MemstoreConfig).RegisterMemstoreConfigHandlers(router)
	}

	restEnabledCfg := cfg.EnabledServersConfig.RestAPIConfig
	if restEnabledCfg.Enabled() {
		if err := restServer.Start(router); err != nil {
			return fmt.Errorf("start proxy rest server: %w", err)
		}

		log.Info("Started EigenDA Proxy REST ALT DA server",
			enabled_apis.Admin, restEnabledCfg.Admin,
			enabled_apis.StandardCommitment, restEnabledCfg.StandardCommitment,
			enabled_apis.OpGenericCommitment, restEnabledCfg.OpGenericCommitment,
			enabled_apis.OpKeccakCommitment, restEnabledCfg.OpKeccakCommitment)

		defer func() {
			if err := restServer.Stop(); err != nil {
				log.Error("failed to stop REST ALT DA server", "err", err)
			} else {
				log.Info("Successfully shutdown REST ALT DA server")
			}

		}()
	}

	if cfg.EnabledServersConfig.ArbCustomDA {
		arbitrumRpcServer, err := arbitrum_altda.NewServer(ctx, &cfg.ArbCustomDASvrCfg)
		if err != nil {
			return fmt.Errorf("new arbitrum custom da json rpc server: %w", err)
		}

		if err := arbitrumRpcServer.Start(); err != nil {
			return fmt.Errorf("start arbitrum custom da json rpc server: %w", err)
		}

		defer func() {
			if err := arbitrumRpcServer.Stop(); err != nil {
				log.Error("failed to stop arbitrum custom da json rpc server", "err", err)
			} else {
				log.Info("Successfully shutdown Arbitrum Custom DA server")
			}
		}()

		log.Info("Started Arbitrum Custom DA JSON RPC server", "addr", arbitrumRpcServer.Addr())
	}

	if cfg.EnabledServersConfig.Metric {
		log.Info("Starting metrics server", "addr", cfg.MetricsSvrConfig.Host, "port", cfg.MetricsSvrConfig.Port)
		svr := proxy_metrics.NewServer(registry, cfg.MetricsSvrConfig)
		err := svr.Start()
		if err != nil {
			return fmt.Errorf("failed to start metrics server: %w", err)
		}
		defer func() {
			if err := svr.Stop(context.Background()); err != nil {
				log.Error("failed to stop metrics server", "err", err)
			} else {
				log.Info("Successfully shutdown Metrics server")
			}
		}()
		log.Info("started metrics server", "addr", svr.Addr())
		metrics.RecordUp()
	}

	return ctxinterrupt.Wait(cliCtx.Context)
}
