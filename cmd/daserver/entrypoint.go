package main

import (
	"fmt"

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

	server := plasma.NewDAServer(cliCtx.String(ListenAddrFlagName), cliCtx.Int(PortFlagName), store, log)

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

	opio.BlockOnInterrupts()

	return nil
}
