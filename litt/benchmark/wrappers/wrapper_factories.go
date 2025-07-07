package wrappers

import "github.com/Layr-Labs/eigenda/litt/benchmark/config"

// A WrapperFactory is a function that creates a DatabaseWrapper based on the provided configuration.
type WrapperFactory func(cfg *config.BenchmarkConfig) (DatabaseWrapper, error)

// WrapperFactories is a map of wrapper factories, where the key is the name of the database and the value is a function
// that creates a DatabaseWrapper for that database type.
var WrapperFactories = map[string]WrapperFactory{
	"littdb":  NewLittDBWrapper,
	"leveldb": NewLevelDBWrapper,
}
