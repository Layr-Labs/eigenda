package indexer

import (
	"fmt"
	"time"
)

type Config struct {
	// The frequency to pull data from The Graph.
	PullInterval time.Duration
}

// DefaultIndexerConfig returns the default indexer configuration.
func DefaultIndexerConfig() Config {
	return Config{
		PullInterval: 1 * time.Second,
	}
}

// Verify validates the indexer configuration.
func (c *Config) Verify() error {
	if c.PullInterval <= 0 {
		return fmt.Errorf("pull interval must be positive, got %v", c.PullInterval)
	}
	return nil
}
