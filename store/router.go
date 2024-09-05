//go:generate mockgen -package mocks --destination ../mocks/router.go . IRouter

package store

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/Layr-Labs/eigenda-proxy/commitments"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
)

type IRouter interface {
	Get(ctx context.Context, key []byte, cm commitments.CommitmentMode) ([]byte, error)
	Put(ctx context.Context, cm commitments.CommitmentMode, key, value []byte) ([]byte, error)

	GetEigenDAStore() KeyGeneratedStore
	GetS3Store() PrecomputedKeyStore
}

// Router ... storage backend routing layer
type Router struct {
	log     log.Logger
	eigenda KeyGeneratedStore
	s3      PrecomputedKeyStore

	caches    []PrecomputedKeyStore
	cacheLock sync.RWMutex

	fallbacks    []PrecomputedKeyStore
	fallbackLock sync.RWMutex
}

func NewRouter(eigenda KeyGeneratedStore, s3 PrecomputedKeyStore, l log.Logger,
	caches []PrecomputedKeyStore, fallbacks []PrecomputedKeyStore) (IRouter, error) {
	return &Router{
		log:          l,
		eigenda:      eigenda,
		s3:           s3,
		caches:       caches,
		cacheLock:    sync.RWMutex{},
		fallbacks:    fallbacks,
		fallbackLock: sync.RWMutex{},
	}, nil
}

// Get ... fetches a value from a storage backend based on the (commitment mode, type)
func (r *Router) Get(ctx context.Context, key []byte, cm commitments.CommitmentMode) ([]byte, error) {
	switch cm {
	case commitments.OptimismGeneric:

		if r.s3 == nil {
			return nil, errors.New("expected S3 backend for OP keccak256 commitment type, but none configured")
		}

		r.log.Debug("Retrieving data from S3 backend")
		value, err := r.s3.Get(ctx, key)
		if err != nil {
			return nil, err
		}

		err = r.s3.Verify(key, value)
		if err != nil {
			return nil, err
		}
		return value, nil

	case commitments.SimpleCommitmentMode, commitments.OptimismAltDA:
		if r.eigenda == nil {
			return nil, errors.New("expected EigenDA backend for DA commitment type, but none configured")
		}

		// 1 - read blob from cache if enabled
		if r.cacheEnabled() {
			r.log.Debug("Retrieving data from cached backends")
			data, err := r.multiSourceRead(ctx, key, false)
			if err == nil {
				return data, nil
			}
			r.log.Warn("Failed to read from cache targets", "err", err)
		}

		// 2 - read blob from EigenDA
		data, err := r.eigenda.Get(ctx, key)
		if err == nil {
			// verify
			err = r.eigenda.Verify(key, data)
			if err != nil {
				return nil, err
			}
			return data, nil
		}

		// 3 - read blob from fallbacks if enabled and data is non-retrievable from EigenDA
		if r.fallbackEnabled() {
			data, err = r.multiSourceRead(ctx, key, true)
			if err != nil {
				r.log.Error("Failed to read from fallback targets", "err", err)
				return nil, err
			}
		} else {
			return nil, err
		}

		return data, err

	default:
		return nil, errors.New("could not determine which storage backend to route to based on unknown commitment mode")
	}
}

// Put ... inserts a value into a storage backend based on the commitment mode
func (r *Router) Put(ctx context.Context, cm commitments.CommitmentMode, key, value []byte) ([]byte, error) {
	var commit []byte
	var err error

	switch cm {
	case commitments.OptimismGeneric: // caching and fallbacks are unsupported for this commitment mode
		return r.putWithKey(ctx, key, value)
	case commitments.OptimismAltDA, commitments.SimpleCommitmentMode:
		commit, err = r.putWithoutKey(ctx, value)
	default:
		return nil, fmt.Errorf("unknown commitment mode")
	}

	if err != nil {
		return nil, err
	}

	if r.cacheEnabled() || r.fallbackEnabled() {
		err = r.handleRedundantWrites(ctx, commit, value)
		if err != nil {
			log.Error("Failed to write to redundant backends", "err", err)
		}
	}

	return commit, nil
}

// handleRedundantWrites ... writes to both sets of backends (i.e, fallback, cache)
// and returns an error if NONE of them succeed
// NOTE: multi-target set writes are done at once to avoid re-invocation of the same write function at the same
// caller step for different target sets vs. reading which is done conditionally to segment between a cached read type
// vs a fallback read type
func (r *Router) handleRedundantWrites(ctx context.Context, commitment []byte, value []byte) error {
	r.cacheLock.RLock()
	r.fallbackLock.RLock()

	defer func() {
		r.cacheLock.RUnlock()
		r.fallbackLock.RUnlock()
	}()

	sources := r.caches
	sources = append(sources, r.fallbacks...)

	key := crypto.Keccak256(commitment)
	successes := 0

	for _, src := range sources {
		err := src.Put(ctx, key, value)
		if err != nil {
			r.log.Warn("Failed to write to redundant target", "backend", src.BackendType(), "err", err)
		} else {
			successes++
		}
	}

	if successes == 0 {
		return errors.New("failed to write blob to any redundant targets")
	}

	return nil
}

// multiSourceRead ... reads from a set of backends and returns the first successfully read blob
func (r *Router) multiSourceRead(ctx context.Context, commitment []byte, fallback bool) ([]byte, error) {
	var sources []PrecomputedKeyStore
	if fallback {
		r.fallbackLock.RLock()
		defer r.fallbackLock.RUnlock()

		sources = r.fallbacks
	} else {
		r.cacheLock.RLock()
		defer r.cacheLock.RUnlock()

		sources = r.caches
	}

	key := crypto.Keccak256(commitment)
	for _, src := range sources {
		data, err := src.Get(ctx, key)
		if err != nil {
			r.log.Warn("Failed to read from redundant target", "backend", src.BackendType(), "err", err)
			continue
		}
		// verify cert:data using EigenDA verification checks
		err = r.eigenda.Verify(commitment, data)
		if err != nil {
			log.Warn("Failed to verify blob", "err", err, "backend", src.BackendType())
			continue
		}

		return data, nil
	}
	return nil, errors.New("no data found in any redundant backend")
}

// putWithoutKey ... inserts a value into a storage backend that computes the key on-demand (i.e, EigenDA)
func (r *Router) putWithoutKey(ctx context.Context, value []byte) ([]byte, error) {
	if r.eigenda != nil {
		r.log.Debug("Storing data to EigenDA backend")
		return r.eigenda.Put(ctx, value)
	}

	return nil, errors.New("no DA storage backend found")
}

// putWithKey ... only supported for S3 storage backends using OP's alt-da keccak256 commitment type
func (r *Router) putWithKey(ctx context.Context, key []byte, value []byte) ([]byte, error) {
	if r.s3 == nil {
		return nil, errors.New("S3 is disabled but is only supported for posting known commitment keys")
	}

	err := r.s3.Verify(key, value)
	if err != nil {
		return nil, err
	}

	return key, r.s3.Put(ctx, key, value)
}

func (r *Router) fallbackEnabled() bool {
	return len(r.fallbacks) > 0
}

func (r *Router) cacheEnabled() bool {
	return len(r.caches) > 0
}

// GetEigenDAStore ...
func (r *Router) GetEigenDAStore() KeyGeneratedStore {
	return r.eigenda
}

// GetS3Store ...
func (r *Router) GetS3Store() PrecomputedKeyStore {
	return r.s3
}
