package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Layr-Labs/eigenda-proxy/flags"
	"github.com/Layr-Labs/eigenda-proxy/metrics"
	"github.com/Layr-Labs/eigenda-proxy/server"
	"github.com/urfave/cli/v2"

	"github.com/ethereum-optimism/optimism/op-service/ctxinterrupt"
	oplog "github.com/ethereum-optimism/optimism/op-service/log"
)

func StartProxySvr(cliCtx *cli.Context) error {
	log := oplog.NewLogger(oplog.AppOut(cliCtx), oplog.ReadCLIConfig(cliCtx)).New("role", "eigenda_proxy")
	oplog.SetGlobalLogHandler(log.Handler())
	log.Info("Starting EigenDA Proxy Server", "version", Version, "date", Date, "commit", Commit)

	cfg := server.ReadCLIConfig(cliCtx)
	if err := cfg.Check(); err != nil {
		return err
	}
	ctx, ctxCancel := context.WithCancel(cliCtx.Context)
	defer ctxCancel()

	configJSON, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	log.Info(fmt.Sprintf("Initializing EigenDA proxy server with config: %v", string(configJSON)))

	daRouter, err := server.LoadStoreRouter(ctx, cfg, log)
	if err != nil {
		return fmt.Errorf("failed to create store: %w", err)
	}
	m := metrics.NewMetrics("default")
	server := server.NewServer(cliCtx.String(flags.ListenAddrFlagName), cliCtx.Int(flags.PortFlagName), daRouter, log, m)

	if err := server.Start(); err != nil {
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
		log.Debug("starting metrics server", "addr", cfg.MetricsCfg.ListenAddr, "port", cfg.MetricsCfg.ListenPort)
		svr, err := m.StartServer(cfg.MetricsCfg.ListenAddr, cfg.MetricsCfg.ListenPort)
		if err != nil {
			return fmt.Errorf("failed to start metrics server: %w", err)
		}
		defer func() {
			if err := svr.Stop(context.Background()); err != nil {
				log.Error("failed to stop metrics server", "err", err)
			}
		}()
		log.Info("started metrics server", "addr", svr.Addr())
		m.RecordUp()
	}

	return ctxinterrupt.Wait(cliCtx.Context)
}
