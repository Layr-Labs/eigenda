package litt

import "time"

type DBType int

const DiskDB DBType = 0
const MemDB DBType = 1

// Config is configuration for a litt.DB.
type Config struct {
	// The path where the database will store its files. If the path does not exist, it will be created.
	// If the path exists, the database will attempt to open the existing database at that path.
	Path string

	// The type of the DB. Choices are DiskDB and MemDB. Default is DiskDB.
	Type DBType

	// The default TTL for newly created tables (either ones with data on disk or new tables).
	// The default is 0 (no TTL).
	TTL time.Duration

	// The period between garbage collection runs. The default is 5 minutes.
	GCPeriod time.Duration

	// The time source used by the database. This can be substituted for an artificial time source
	// for testing purposes. The default is time.Now.
	TimeSource func() time.Time
}

// DefaultConfig returns a Config with default values.
func DefaultConfig() *Config {
	return &Config{
		TimeSource: time.Now,
		GCPeriod:   5 * time.Minute,
	}
}
