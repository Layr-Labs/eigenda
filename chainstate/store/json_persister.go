package store

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Layr-Labs/eigenda/litt/util"
	"github.com/Layr-Labs/eigensdk-go/logging"
)

// JSONPersister handles periodic persistence of store state to a JSON file.
type JSONPersister struct {
	store  Store
	path   string
	logger logging.Logger
}

// NewJSONPersister creates a new JSON persister for the given store.
func NewJSONPersister(store Store, path string, logger logging.Logger) *JSONPersister {
	return &JSONPersister{
		store:  store,
		path:   path,
		logger: logger,
	}
}

// Save persists the current store state to the configured JSON file.
// It uses util.AtomicWrite (write to a swap file, fsync, rename, fsync the
// directory) so a crash mid-save can never leave a truncated or torn state
// file behind.
func (p *JSONPersister) Save(ctx context.Context) error {
	data, err := p.store.Snapshot()
	if err != nil {
		return fmt.Errorf("failed to create snapshot: %w", err)
	}

	if err := util.AtomicWrite(p.path, data, true); err != nil {
		return fmt.Errorf("failed to write state file: %w", err)
	}

	p.logger.Info("State persisted", "path", p.path, "size_bytes", len(data))
	return nil
}

// Load restores the store state from the configured JSON file.
// If the file doesn't exist, it returns without error (fresh start).
func (p *JSONPersister) Load(ctx context.Context) error {
	data, err := os.ReadFile(p.path)
	if err != nil {
		if os.IsNotExist(err) {
			p.logger.Info("No existing state file, starting fresh", "path", p.path)
			return nil
		}
		return fmt.Errorf("failed to read state file: %w", err)
	}

	if err := p.store.Restore(data); err != nil {
		return fmt.Errorf("failed to restore state: %w", err)
	}

	p.logger.Info("State restored", "path", p.path, "size_bytes", len(data))
	return nil
}

// StartPeriodicSave starts a background goroutine that periodically saves the store state.
// It also performs a final save when the context is cancelled.
func (p *JSONPersister) StartPeriodicSave(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := p.Save(ctx); err != nil {
				p.logger.Error("Failed to persist state", "error", err)
			}
		case <-ctx.Done():
			// Perform final save before shutdown
			p.logger.Info("Context cancelled, performing final state save")
			if err := p.Save(context.Background()); err != nil {
				p.logger.Error("Failed final state save", "error", err)
			}
			return
		}
	}
}
