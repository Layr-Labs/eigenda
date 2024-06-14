package server

import (
	"context"

	"github.com/Layr-Labs/eigenda-proxy/verify"
	"github.com/Layr-Labs/eigenda/api/clients"
	"github.com/ethereum/go-ethereum/log"
)

func LoadStore(cfg CLIConfig, ctx context.Context, log log.Logger) (Store, error) {
	log.Info("Using eigenda backend")
	daCfg := cfg.EigenDAConfig
	vCfg := daCfg.VerificationCfg()

	verifier, err := verify.NewVerifier(vCfg, log)
	if err != nil {
		return nil, err
	}

	if vCfg.Verify {
		log.Info("Certificate verification with Ethereum enabled")
	} else {
		log.Warn("Verification disabled")
	}

	maxBlobLength, err := daCfg.GetMaxBlobLength()
	if err != nil {
		return nil, err
	}

	if cfg.EigenDAConfig.MemstoreEnabled {
		log.Info("Using memstore backend")
		return NewMemStore(ctx, verifier, log, maxBlobLength, cfg.EigenDAConfig.MemstoreBlobExpiration)
	}

	log.Info("Using EigenDA backend")
	client, err := clients.NewEigenDAClient(log, daCfg.ClientConfig)
	if err != nil {
		return nil, err
	}
	return NewEigenDAStore(
		ctx,
		client,
		verifier,
		maxBlobLength,
	)
}
