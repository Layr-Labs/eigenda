package server

import (
	"context"
	"fmt"

	"github.com/Layr-Labs/eigenda-proxy/store"
	"github.com/Layr-Labs/eigenda-proxy/verify"
	"github.com/Layr-Labs/eigenda/api/clients"
	"github.com/ethereum/go-ethereum/log"
)

// populateTargets ... creates a list of storage backends based on the provided target strings
func populateTargets(targets []string, s3 *store.S3Store, redis *store.RedStore) []store.PrecomputedKeyStore {
	stores := make([]store.PrecomputedKeyStore, len(targets))

	for i, f := range targets {
		b := store.StringToBackendType(f)

		switch b {
		case store.Redis:
			stores[i] = redis

		case store.S3:
			stores[i] = s3

		case store.EigenDA, store.Memory:
			panic(fmt.Sprintf("Invalid target for fallback: %s", f))

		case store.Unknown:
			fallthrough

		default:
			panic(fmt.Sprintf("Unknown fallback target: %s", f))
		}
	}

	return stores
}

// LoadStoreRouter ... creates storage backend clients and instruments them into a storage routing abstraction
func LoadStoreRouter(ctx context.Context, cfg CLIConfig, log log.Logger) (store.IRouter, error) {
	// create S3 backend store (if enabled)
	var err error
	var s3 *store.S3Store
	var redis *store.RedStore

	if cfg.S3Config.Bucket != "" && cfg.S3Config.Endpoint != "" {
		log.Info("Using S3 backend")
		s3, err = store.NewS3(cfg.S3Config)
		if err != nil {
			return nil, err
		}
	}

	if cfg.RedisCfg.Endpoint != "" {
		log.Info("Using Redis backend")
		// create Redis backend store
		redis, err = store.NewRedisStore(&cfg.RedisCfg)
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
		eigenda, err = store.NewMemStore(ctx, verifier, log, store.MemStoreConfig{
			MaxBlobSizeBytes: maxBlobLength,
			BlobExpiration:   cfg.EigenDAConfig.MemstoreBlobExpiration,
			PutLatency:       cfg.EigenDAConfig.MemstorePutLatency,
			GetLatency:       cfg.EigenDAConfig.MemstoreGetLatency,
		})
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
				EthConfirmationDepth: uint64(cfg.EigenDAConfig.EthConfirmationDepth), // #nosec G115
				StatusQueryTimeout:   cfg.EigenDAConfig.ClientConfig.StatusQueryTimeout,
			},
		)
	}

	if err != nil {
		return nil, err
	}

	// determine read fallbacks
	fallbacks := populateTargets(cfg.EigenDAConfig.FallbackTargets, s3, redis)
	caches := populateTargets(cfg.EigenDAConfig.CacheTargets, s3, redis)

	log.Info("Creating storage router", "eigenda backend type", eigenda != nil, "s3 backend type", s3 != nil)
	return store.NewRouter(eigenda, s3, log, caches, fallbacks)
}
