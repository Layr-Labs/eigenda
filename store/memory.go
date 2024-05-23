package store

import (
	"context"
	"crypto/rand"
	"fmt"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/urfave/cli/v2"
)

const (
	MemStoreFlagName   = "memstore.enabled"
	ExpirationFlagName = "memstore.expiration"

	DefaultPruneInterval = 500 * time.Millisecond
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

// NewMemStore ... constructor
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
	e.Lock()
	defer e.Unlock()

	for commit, dur := range e.keyStarts {
		if time.Since(dur) >= e.cfg.BlobExpiration {
			delete(e.keyStarts, commit)
			delete(e.store, commit)
		}
	}

}

// Get fetches a value from the store.
func (e *MemStore) Get(ctx context.Context, commit []byte) ([]byte, error) {
	e.RLock()
	defer e.RUnlock()

	key := common.Bytes2Hex(commit)
	if _, exists := e.store[key]; !exists {
		return nil, fmt.Errorf("commitment key not found")
	}

	return e.store[key], nil
}

// Put inserts a value into the store.
func (e *MemStore) Put(ctx context.Context, value []byte) ([]byte, error) {
	e.Lock()
	defer e.Unlock()

	fingerprint := crypto.Keccak256Hash(value)
	// add some entropy to commit to emulate randomness seen in EigenDA
	// when generating operator BLS signature certificates
	entropy := make([]byte, 10)
	rand.Read(entropy)

	rawCommit := append(fingerprint.Bytes(), entropy...)
	commit := common.Bytes2Hex(rawCommit)

	if _, exists := e.store[commit]; exists {
		return nil, fmt.Errorf("commitment key already exists")
	}

	e.store[commit] = value
	// add expiration
	e.keyStarts[commit] = time.Now()

	return rawCommit, nil
}

func ReadConfig(ctx *cli.Context) MemStoreConfig {
	cfg := MemStoreConfig{
		/* Required Flags */
		Enabled:        ctx.Bool(MemStoreFlagName),
		BlobExpiration: ctx.Duration(ExpirationFlagName),
	}
	return cfg
}

func CLIFlags(envPrefix string) []cli.Flag {

	return []cli.Flag{
		&cli.BoolFlag{
			Name:    MemStoreFlagName,
			Usage:   "Whether to use mem-store for DA logic.",
			EnvVars: []string{"MEMSTORE_ENABLED"},
		},
		&cli.DurationFlag{
			Name:    ExpirationFlagName,
			Usage:   "Duration that a blob/commitment pair are allowed to live.",
			Value:   25 * time.Minute,
			EnvVars: []string{"MEMSTORE_EXPIRATION"},
		},
	}
}
