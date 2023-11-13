package common

import "github.com/ethereum/go-ethereum/log"

type Logger interface {
	// New returns a new Logger that has this logger's context plus the given context
	New(ctx ...interface{}) Logger

	// SetHandler updates the logger to write records to the specified handler.
	SetHandler(h log.Handler)

	// Log a message at the trace level with context key/value pairs
	//
	// # Usage
	//
	//	log.Trace("msg")
	//	log.Trace("msg", "key1", val1)
	//	log.Trace("msg", "key1", val1, "key2", val2)
	Trace(msg string, ctx ...interface{})

	// Log a message at the debug level with context key/value pairs
	//
	// # Usage Examples
	//
	//	log.Debug("msg")
	//	log.Debug("msg", "key1", val1)
	//	log.Debug("msg", "key1", val1, "key2", val2)
	Debug(msg string, ctx ...interface{})

	// Log a message at the info level with context key/value pairs
	//
	// # Usage Examples
	//
	//	log.Info("msg")
	//	log.Info("msg", "key1", val1)
	//	log.Info("msg", "key1", val1, "key2", val2)
	Info(msg string, ctx ...interface{})

	// Log a message at the warn level with context key/value pairs
	//
	// # Usage Examples
	//
	//	log.Warn("msg")
	//	log.Warn("msg", "key1", val1)
	//	log.Warn("msg", "key1", val1, "key2", val2)
	Warn(msg string, ctx ...interface{})

	// Log a message at the error level with context key/value pairs
	//
	// # Usage Examples
	//
	//	log.Error("msg")
	//	log.Error("msg", "key1", val1)
	//	log.Error("msg", "key1", val1, "key2", val2)
	Error(msg string, ctx ...interface{})

	// Log a message at the crit level with context key/value pairs, and then exit.
	//
	// # Usage Examples
	//
	//	log.Crit("msg")
	//	log.Crit("msg", "key1", val1)
	//	log.Crit("msg", "key1", val1, "key2", val2)
	Crit(msg string, ctx ...interface{})

	// Fatal is an alias for Crit
	// Log a message at the crit level with context key/value pairs, and then exit.
	//
	// # Usage Examples
	//
	//	log.Fatal("msg")
	//	log.Fatal("msg", "key1", val1)
	//	log.Fatal("msg", "key1", val1, "key2", val2)
	Fatal(msg string, ctx ...interface{})

	// We add the below methods to be compliant with the eigensdk Logger interface
	Debugf(template string, args ...interface{})
	Infof(template string, args ...interface{})
	Warnf(template string, args ...interface{})
	Errorf(template string, args ...interface{})
	Critf(template string, args ...interface{})
	// eigensdk uses fatal instead of crit so we add it,
	// but should have same semantic as Critf
	Fatalf(template string, args ...interface{})
}
