package store

import (
	"context"
	"fmt"
	"os"
	"time"

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
// It uses atomic file operations (write to temp, then rename) to ensure consistency.
func (p *JSONPersister) Save(ctx context.Context) error {
	data, err := p.store.Snapshot()
	if err != nil {
		return fmt.Errorf("failed to create snapshot: %w", err)
	}

	tmpPath := p.path + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write temp file: %w", err)
	}

	if err := os.Rename(tmpPath, p.path); err != nil {
		return fmt.Errorf("failed to rename temp file: %w", err)
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
