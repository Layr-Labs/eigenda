package server

import (
	"context"

	"github.com/Layr-Labs/eigenda-proxy/store"
	"github.com/Layr-Labs/eigenda-proxy/verify"
	"github.com/Layr-Labs/eigenda/api/clients"
	"github.com/ethereum/go-ethereum/log"
)

func LoadStore(cfg CLIConfig, ctx context.Context, log log.Logger) (store.Store, error) {
	log.Info("Using eigenda backend")
	daCfg := cfg.EigenDAConfig

	verifier, err := verify.NewVerifier(daCfg.KzgConfig())
	if err != nil {
		return nil, err
	}

	maxBlobLength, err := daCfg.GetMaxBlobLength()
	if err != nil {
		return nil, err
	}

	if cfg.MemStoreCfg.Enabled {
		log.Info("Using memstore backend")
		return store.NewMemStore(ctx, &cfg.MemStoreCfg, verifier, log, maxBlobLength)
	}

	client, err := clients.NewEigenDAClient(log, daCfg.ClientConfig)
	if err != nil {
		return nil, err
	}
	return store.NewEigenDAStore(
		ctx,
		client,
		verifier,
		maxBlobLength,
	)
}
