package ephemeraldb

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/Layr-Labs/eigenda/api"
	"github.com/Layr-Labs/eigenda/api/proxy/common/proxyerrors"
	"github.com/Layr-Labs/eigenda/api/proxy/store/generated_key/memstore/memconfig"
	"github.com/Layr-Labs/eigensdk-go/logging"
)

const (
	DefaultPruneInterval = 500 * time.Millisecond
)

// DB ... An ephemeral && simple in-memory database used to emulate
// an EigenDA network for dispersal/retrieval operations.
type DB struct {
	// knobs used to express artificial conditions for testing
	config *memconfig.SafeConfig
	log    logging.Logger

	// mu guards the below fields
	mu        sync.RWMutex
	keyStarts map[string]time.Time // used for managing expiration
	store     map[string][]byte    // db
}

// New ... constructor
func New(ctx context.Context, cfg *memconfig.SafeConfig, log logging.Logger) *DB {
	db := &DB{
		config:    cfg,
		keyStarts: make(map[string]time.Time),
		store:     make(map[string][]byte),
		log:       log,
	}

	// if no expiration set then blobs will be persisted indefinitely
	if cfg.BlobExpiration() != 0 {
		db.log.Info("ephemeral db expiration enabled for payload entries.", "time", cfg.BlobExpiration)
		go db.pruningLoop(ctx)
	}

	return db
}

// InsertEntry ... inserts a value into the db provided a key
func (db *DB) InsertEntry(key []byte, value []byte) error {
	if db.config.PutReturnsFailoverError() {
		return api.NewErrorFailover(errors.New("ephemeral db in failover simulation mode"))
	}
	if uint64(len(value)) > db.config.MaxBlobSizeBytes() {
		return fmt.Errorf(
			"%w: blob length %d, max blob size %d",
			proxyerrors.ErrProxyOversizedBlob,
			len(value),
			db.config.MaxBlobSizeBytes())
	}

	time.Sleep(db.config.LatencyPUTRoute())
	db.mu.Lock()
	defer db.mu.Unlock()

	strKey := string(key)

	_, exists := db.store[strKey]
	if exists {
		return fmt.Errorf("payload key already exists in ephemeral db: %s", strKey)
	}

	db.store[strKey] = value
	// add expiration if applicable

	if db.config.BlobExpiration() > 0 {
		db.keyStarts[strKey] = time.Now()
	}

	return nil
}

// FetchEntry ... looks up a value from the db provided a key
func (db *DB) FetchEntry(key []byte) ([]byte, error) {
	time.Sleep(db.config.LatencyGETRoute())
	db.mu.RLock()
	defer db.mu.RUnlock()

	payload, exists := db.store[string(key)]

	if !exists {
		return nil, fmt.Errorf("payload not found for key: %s", hex.EncodeToString(key))
	}

	return payload, nil
}

// pruningLoop ... runs a background goroutine to prune expired blobs from the store on a regular interval.
func (db *DB) pruningLoop(ctx context.Context) {
	timer := time.NewTicker(DefaultPruneInterval)

	for {
		select {
		case <-ctx.Done():
			return

		case <-timer.C:
			db.pruneExpired()
		}
	}
}

// pruneExpired ... removes expired blobs from the store based on the expiration time.
func (db *DB) pruneExpired() {
	db.mu.Lock()
	defer db.mu.Unlock()

	for commit, dur := range db.keyStarts {
		if time.Since(dur) >= db.config.BlobExpiration() {
			delete(db.keyStarts, commit)
			delete(db.store, commit)

			db.log.Debug("blob pruned", "commit", commit)
		}
	}
}
