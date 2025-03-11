package util

import (
	"fmt"
	"runtime/debug"
	"sync/atomic"

	"github.com/Layr-Labs/eigensdk-go/logging"
)

// TODO also close context

// DBPanic is a struct that permits the DB to "panic". There are many goroutines that function under the hood, and
// many of these threads could, in theory, encounter errors which are unrecoverable. In such situations, the desirable
// outcome is for the DB to report the error and then refuse to do additional work. If the DB is in a broken state,
// it is much better to refuse to do work than to continue to do work and potentially corrupt data.
//
// Even though this is called "panic", it is not the same as the panic that is built into Go. Once the DB "panics",
// all public methods will return a meaningful error, and the DB will refuse to do additional work.
type DBPanic struct {
	logger logging.Logger

	// If this is non-nil, the DB is in a "panic" state and will refuse to do additional work.
	error atomic.Pointer[error]
}

// NewDBPanic creates a new DBPanic struct.
func NewDBPanic(logger logging.Logger) *DBPanic {
	return &DBPanic{
		logger: logger,
	}
}

// IsOk returns true if the DB is in a good state, and false if the DB is in a "panic" state.
// The error returned is the error that caused the DB to panic, and does not indicate that
// the call to IsOk() failed. If the DB has panicked multiple times, the error returned will
// be the first error that caused the DB to panic.
func (p *DBPanic) IsOk() (bool, error) {
	err := p.error.Load()
	if err != nil {
		return false, *err
	}
	return true, nil
}

// Shutdown causes the DB to enter a "shutdown" state. Once the DB is in a "shutdown" state,
// it will refuse to do additional work.
func (p *DBPanic) Shutdown() {
	err := fmt.Errorf("DB is shut down")

	// don't overwrite the error if there is already an error stored
	p.error.CompareAndSwap(nil, &err)
}

// Panic time! Something just went very wrong. (╯°□°)╯︵ ┻━┻
//
// Panic causes the DB to enter a "panic" state. Once the DB is in a "panic" state, it will refuse to do
// additional work.
func (p *DBPanic) Panic(err error) {
	stackTrace := string(debug.Stack())

	p.logger.Errorf("LittDB encountered an unrecoverable error: %v\n%s", err, stackTrace)

	// only store the error if there isn't already an error stored
	p.error.CompareAndSwap(nil, &err)
}
