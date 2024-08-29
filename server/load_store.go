package server

import (
	"context"
	"fmt"

	"github.com/Layr-Labs/eigenda-proxy/store"
	"github.com/Layr-Labs/eigenda-proxy/verify"
	"github.com/Layr-Labs/eigenda/api/clients"
	"github.com/ethereum/go-ethereum/log"
)

// LoadStoreRouter ... creates storage backend clients and instruments them into a storage routing abstraction
func LoadStoreRouter(ctx context.Context, cfg CLIConfig, log log.Logger) (*store.Router, error) {
	// create S3 backend store (if enabled)
	var err error
	var s3 *store.S3Store
	if cfg.S3Config.Bucket != "" && cfg.S3Config.Endpoint != "" {
		log.Info("Using S3 backend")
		s3, err = store.NewS3(cfg.S3Config)
		if err != nil {
			return nil, err
		}
	}

	// create cert/data verification type
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

	// create EigenDA backend store
	var eigenda store.KeyGeneratedStore
	if cfg.EigenDAConfig.MemstoreEnabled {
		log.Info("Using mem-store backend for EigenDA")
		eigenda, err = store.NewMemStore(ctx, verifier, log, maxBlobLength, cfg.EigenDAConfig.MemstoreBlobExpiration)
	} else {
		var client *clients.EigenDAClient
		log.Info("Using EigenDA backend")
		client, err = clients.NewEigenDAClient(log, daCfg.ClientConfig)
		if err != nil {
			return nil, err
		}

		eigenda, err = store.NewEigenDAStore(
			client,
			verifier,
			log,
			&store.EigenDAStoreConfig{
				MaxBlobSizeBytes:     maxBlobLength,
				EthConfirmationDepth: uint64(cfg.EigenDAConfig.EthConfirmationDepth),
				StatusQueryTimeout:   cfg.EigenDAConfig.ClientConfig.StatusQueryTimeout,
			},
		)
	}

	if err != nil {
		return nil, err
	}

	// determine read fallbacks
	fallbacks := make([]store.PrecomputedKeyStore, len(cfg.EigenDAConfig.FallbackTargets))

	for i, f := range cfg.EigenDAConfig.FallbackTargets {
		b := store.StringToBackendType(f)

		switch b {
		case store.S3:
			fallbacks[i] = s3

		case store.EigenDA, store.Memory:
			return nil, fmt.Errorf("EigenDA cannot be used as a fallback target")

		case store.Redis:
			return nil, fmt.Errorf("redis is not supported yet")

		case store.Unknown:
			fallthrough

		default:
			panic(fmt.Sprintf("Unknown fallback target: %s", f))
		}
	}

	// determine caches for priority reads
	caches := make([]store.PrecomputedKeyStore, len(cfg.EigenDAConfig.CacheTargets))

	for i, f := range cfg.EigenDAConfig.CacheTargets {
		b := store.StringToBackendType(f)

		switch b {
		case store.S3:
			caches[i] = s3

		case store.EigenDA, store.Memory:
			return nil, fmt.Errorf("EigenDA cannot be used as a cache target")

		case store.Redis:
			return nil, fmt.Errorf("redis is not supported yet")

		case store.Unknown:
			fallthrough

		default:
			log.Warn("Unknown fallback target", "target", f)
		}
	}

	log.Info("Creating storage router", "eigenda backend type", eigenda != nil, "s3 backend type", s3 != nil)
	return store.NewRouter(eigenda, s3, log, caches, fallbacks)
}
