package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Layr-Labs/op-plasma-eigenda/metrics"
	"github.com/urfave/cli/v2"

	plasma "github.com/Layr-Labs/op-plasma-eigenda"
	"github.com/Layr-Labs/op-plasma-eigenda/eigenda"
	plasma_store "github.com/Layr-Labs/op-plasma-eigenda/store"
	"github.com/Layr-Labs/op-plasma-eigenda/verify"
	oplog "github.com/ethereum-optimism/optimism/op-service/log"
	"github.com/ethereum-optimism/optimism/op-service/opio"
)

type App struct {
	DAServer   *http.Server
	MetricsSvr *http.Server
}

func StartDAServer(cliCtx *cli.Context) error {
	println("CHECKING")
	if err := CheckRequired(cliCtx); err != nil {
		return err
	}
	println("Reading CLI CFG")
	cfg := ReadCLIConfig(cliCtx)
	if err := cfg.Check(); err != nil {
		return err
	}
	println("HERERERER")
	m := metrics.NewMetrics("default")

	log := oplog.NewLogger(oplog.AppOut(cliCtx), oplog.ReadCLIConfig(cliCtx)).New("role", "eigenda_plasma_server")
	oplog.SetGlobalLogHandler(log.Handler())

	log.Info("Initializing EigenDA Plasma DA server with config ...")

	var store plasma.PlasmaStore

	if cfg.FileStoreEnabled() {
		log.Info("Using file storage", "path", cfg.FileStoreDirPath)
		store = plasma_store.NewFileStore(cfg.FileStoreDirPath)
	} else if cfg.S3Enabled() {
		log.Info("Using S3 storage", "bucket", cfg.S3Bucket)
		s3, err := plasma_store.NewS3Store(cliCtx.Context, cfg.S3Bucket)
		if err != nil {
			return fmt.Errorf("failed to create S3 store: %w", err)
		}
		store = s3
	} else if cfg.EigenDAEnabled() {
		daCfg := cfg.EigenDAConfig

		v, err := verify.NewVerifier(daCfg.KzgConfig())
		if err != nil {
			return err
		}

		eigenda, err := plasma_store.NewEigenDAStore(
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
		store = eigenda
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
		metricsSrv, err := m.StartServer(cfg.MetricsCfg.ListenAddr, cfg.MetricsCfg.ListenPort)
		if err != nil {
			return fmt.Errorf("failed to start metrics server: %w", err)
		}
		defer func() {
			if err := metricsSrv.Stop(context.Background()); err != nil {
				log.Error("failed to stop metrics server", "err", err)
			}
		}()
		log.Info("started metrics server", "addr", metricsSrv.Addr())
		m.RecordUp()
	}

	opio.BlockOnInterrupts()

	return nil
}
