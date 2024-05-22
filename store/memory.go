package store

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/urfave/cli"
)

const (
	MemStoreName       = "memstore"
	MemStoreFlagName   = "enable"
	ExpirationFlagName = "expiration"

	DefaultPruneInterval = 1 * time.Second
)

type MemStoreConfig struct {
	Enabled        bool
	BlobExpiration time.Duration
}

// MemStore is a simple in-memory store for blobs which uses an expiration
// time to evict blobs to best emulate the ephemeral nature of blobs dispersed to
// EigenDA operators.
type MemStore struct {
	sync.RWMutex

	cfg       *MemStoreConfig
	keyStarts map[string]time.Time
	store     map[string][]byte
}

func NewMemStore(ctx context.Context, cfg *MemStoreConfig) (*MemStore, error) {
	store := &MemStore{
		cfg:       cfg,
		keyStarts: make(map[string]time.Time),
		store:     make(map[string][]byte),
	}

	if cfg.BlobExpiration != 0 {
		go store.EventLoop(ctx)
	}

	return store, nil
}

func (e *MemStore) EventLoop(ctx context.Context) {

	timer := time.NewTicker(DefaultPruneInterval)

	select {
	case <-ctx.Done():
		return

	case <-timer.C:
		e.pruneExpired()
	}

}

func (e *MemStore) pruneExpired() {
	e.RLock()
	defer e.RUnlock()

	for commit, dur := range e.keyStarts {
		if time.Since(dur) >= e.cfg.BlobExpiration {
			// prune expired blobs
			e.Lock()
			delete(e.keyStarts, commit)
			delete(e.store, commit)
			e.Unlock()
		}
	}

}

func (e *MemStore) Get(ctx context.Context, key []byte) ([]byte, error) {
	e.RLock()
	defer e.RUnlock()

	if _, exists := e.store[common.Bytes2Hex(key)]; !exists {
		return nil, fmt.Errorf("commitment key not found")
	}

	return e.store[string(key)], nil
}

func (e *MemStore) Put(ctx context.Context, value []byte) ([]byte, error) {
	e.Lock()
	defer e.Unlock()

	commit := crypto.Keccak256Hash(value)

	if _, exists := e.store[commit.String()]; !exists {
		return nil, fmt.Errorf("commitment key not found")
	}

	return commit.Bytes(), nil
}

func ReadConfig(ctx *cli.Context) MemStoreConfig {
	cfg := MemStoreConfig{
		/* Required Flags */
		Enabled:        ctx.Bool(MemStoreName),
		BlobExpiration: ctx.Duration(ExpirationFlagName),
	}
	return cfg
}

func CLIFlags(envPrefix string) []cli.Flag {

	return []cli.Flag{
		&cli.BoolFlag{
			Name:   MemStoreFlagName,
			Usage:  "Whether to use mem-store for DA logic.",
			EnvVar: "MEMSTORE_ENABLED",
		},
		&cli.DurationFlag{
			Name:   ExpirationFlagName,
			Usage:  "Duration that a blob/commitment pair are allowed to live.",
			Value:  25 * time.Minute,
			EnvVar: "MEMSTORE_EXPIRATION",
		},
	}
}
