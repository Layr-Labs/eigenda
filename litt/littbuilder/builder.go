package littbuilder

import (
	"fmt"
	"github.com/Layr-Labs/eigenda/litt"
	"github.com/Layr-Labs/eigenda/litt/memtable"
	"time"
)

// NewDB builds a new litt.DB.
func NewDB(config *Config) (litt.DB, error) {
	var tb tableBuilder
	switch config.Type {
	case DiskDB:
		tb = func(timeSource func() time.Time, name string, ttl time.Duration) (litt.ManagedTable, error) {
			return nil, nil // TODO
		}
	case MemDB:
		tb = func(timeSource func() time.Time, name string, ttl time.Duration) (litt.ManagedTable, error) {
			return memtable.NewMemTable(timeSource, name, ttl), nil
		}
	default:
		return nil, fmt.Errorf("unsupported DB type: %v", config.Type)
	}

	database := newDB(config.TimeSource, config.TTL, config.GCPeriod, tb)
	return database, nil
}
