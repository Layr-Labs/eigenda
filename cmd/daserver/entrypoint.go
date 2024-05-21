package main

import (
	"context"
	"fmt"

	"github.com/Layr-Labs/op-plasma-eigenda/metrics"
	"github.com/urfave/cli/v2"

	plasma "github.com/Layr-Labs/op-plasma-eigenda"
	"github.com/Layr-Labs/op-plasma-eigenda/eigenda"
	plasma_store "github.com/Layr-Labs/op-plasma-eigenda/store"
	"github.com/Layr-Labs/op-plasma-eigenda/verify"
	oplog "github.com/ethereum-optimism/optimism/op-service/log"
	"github.com/ethereum-optimism/optimism/op-service/opio"
)

func StartDAServer(cliCtx *cli.Context) error {
	if err := CheckRequired(cliCtx); err != nil {
		return err
	}
	cfg := ReadCLIConfig(cliCtx)
	if err := cfg.Check(); err != nil {
		return err
	}
	m := metrics.NewMetrics("default")

	log := oplog.NewLogger(oplog.AppOut(cliCtx), oplog.ReadCLIConfig(cliCtx)).New("role", "eigenda_plasma_server")
	oplog.SetGlobalLogHandler(log.Handler())

	log.Info("Initializing EigenDA Plasma DA server...")

	daCfg := cfg.EigenDAConfig

	v, err := verify.NewVerifier(daCfg.KzgConfig())
	if err != nil {
		return err
	}

	store, err := plasma_store.NewEigenDAStore(
		cliCtx.Context,
		eigenda.NewEigenDAClient(
			log,
			daCfg,
		),
		v,
	)
	if err != nil {
		return fmt.Errorf("failed to create EigenDA store: %w", err)
	}
	server := plasma.NewDAServer(cliCtx.String(ListenAddrFlagName), cliCtx.Int(PortFlagName), store, log, m)

	if err := server.Start(); err != nil {
		return fmt.Errorf("failed to start the DA server")
	} else {
		log.Info("Started DA Server")
	}

	defer func() {
		if err := server.Stop(); err != nil {
			log.Error("failed to stop DA server", "err", err)
		}
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
