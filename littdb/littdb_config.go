package littdb

import "time"

type LittDBConfig struct {
	// The path where the database will store its files. If the path does not exist, it will be created.
	// If the path exists, the database will attempt to open the existing database at that path.
	Path string

	// The time-to-live for values in the database. Values will be automatically deleted after this duration.
	// If the database is reloaded with a different TTL, the new TTL will apply to all values, even values
	// written with the old TTL. If zero, values will never be deleted. The default is zero.
	TTL time.Duration

	// The time source used by the database. This can be substituted for an artificial time source
	// for testing purposes. The default is time.Now.
	TimeSource func() time.Time
}

// DefaultLittDBConfig returns a LittDBConfig with default values.
func DefaultLittDBConfig() *LittDBConfig {
	return &LittDBConfig{
		TimeSource: time.Now,
	}
}
