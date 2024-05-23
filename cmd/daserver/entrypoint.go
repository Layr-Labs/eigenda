package main

import (
	"context"
	"fmt"

	"github.com/Layr-Labs/op-plasma-eigenda/eigenda"
	"github.com/Layr-Labs/op-plasma-eigenda/metrics"
	"github.com/Layr-Labs/op-plasma-eigenda/store"
	"github.com/Layr-Labs/op-plasma-eigenda/verify"
	"github.com/ethereum/go-ethereum/log"
	"github.com/urfave/cli/v2"

	plasma "github.com/Layr-Labs/op-plasma-eigenda"
	oplog "github.com/ethereum-optimism/optimism/op-service/log"
	"github.com/ethereum-optimism/optimism/op-service/opio"
)

func LoadStore(cfg CLIConfig, ctx context.Context, log log.Logger) (plasma.PlasmaStore, error) {
	if cfg.MemStoreCfg.Enabled {
		log.Info("Using memstore backend")
		return store.NewMemStore(ctx, &cfg.MemStoreCfg)
	}

	log.Info("Using eigenda backend")
	daCfg := cfg.EigenDAConfig

	v, err := verify.NewVerifier(daCfg.KzgConfig())
	if err != nil {
		return nil, err
	}

	return store.NewEigenDAStore(
		ctx,
		eigenda.NewEigenDAClient(
			log,
			daCfg,
		),
		v,
	)
}

func StartDAServer(cliCtx *cli.Context) error {
	if err := CheckRequired(cliCtx); err != nil {
		return err
	}
	cfg := ReadCLIConfig(cliCtx)
	if err := cfg.Check(); err != nil {
		return err
	}
	ctx, ctxCancel := context.WithCancel(cliCtx.Context)
	defer ctxCancel()

	m := metrics.NewMetrics("default")

	log := oplog.NewLogger(oplog.AppOut(cliCtx), oplog.ReadCLIConfig(cliCtx)).New("role", "eigenda_plasma_server")
	oplog.SetGlobalLogHandler(log.Handler())

	log.Info("Initializing EigenDA Plasma DA server...")

	da, err := LoadStore(cfg, ctx, log)
	if err != nil {
		return fmt.Errorf("failed to create store: %w", err)
	}
	server := plasma.NewDAServer(cliCtx.String(ListenAddrFlagName), cliCtx.Int(PortFlagName), da, log, m)

	if err := server.Start(); err != nil {
		return fmt.Errorf("failed to start the DA server")
	} else {
		log.Info("Started DA Server")
	}

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

	opio.BlockOnInterrupts()

	return nil
}
