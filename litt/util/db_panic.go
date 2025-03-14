package util

import (
	"context"
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
	ctx    context.Context
	cancel context.CancelFunc

	logger logging.Logger

	// If this is non-nil, the DB is in a "panic" state and will refuse to do additional work.
	error atomic.Pointer[error]
}

// NewDBPanic creates a new DBPanic struct.
func NewDBPanic(ctx context.Context, logger logging.Logger) *DBPanic {
	ctx, cancel := context.WithCancel(ctx)

	return &DBPanic{
		ctx:    ctx,
		cancel: cancel,
		logger: logger,
	}
}

// Await waits for a value to be sent on a channel. If the channel sends a value, the value is returned. If the DB
// panics before the channel sends a value, an error is returned.
func Await[T any](panic *DBPanic, channel <-chan T) (T, error) {
	select {
	case value := <-channel:
		return value, nil
	case <-panic.ImmediateShutdownRequired():
		var zero T
		return zero, fmt.Errorf("DB context cancelled")
	}
}

// Send sends a value on a channel. If the value is sent, nil is returned. If the DB panics before the value is sent,
// an error is returned.
func Send[T any](panic *DBPanic, channel chan<- T, value T) error {
	select {
	case channel <- value:
		return nil
	case <-panic.ImmediateShutdownRequired():
		return fmt.Errorf("DB context cancelled")
	}
}

// SendAny is identical to Send, but works with channels of type `chan<- any`.
func SendAny[T any](panic *DBPanic, channel chan<- any, value T) error {
	select {
	case channel <- value:
		return nil
	case <-panic.ImmediateShutdownRequired():
		return fmt.Errorf("DB context cancelled")
	}
}

// ImmediateShutdownRequired returns a channel that is closed when the DB "panics". The channel might also be
// closed if the parent context is cancelled, and so this channel being closed can't be used to infer that the
// DB is in a panicked state. When this channel is closed, it is expected that all DB goroutines immediately shut down.
func (p *DBPanic) ImmediateShutdownRequired() <-chan struct{} {
	return p.ctx.Done()
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
// it will refuse to do additional work. Does not cancel the context.
func (p *DBPanic) Shutdown() {
	err := fmt.Errorf("DB is shut down")

	// don't overwrite the error if there is already an error stored
	p.error.CompareAndSwap(nil, &err)
}

// Panic time! Something just went very wrong. (╯°□°)╯︵ ┻━┻
//
// Panic causes the DB to enter a "panic" state. Once the DB is in a "panic" state, it will refuse to do
// additional work. As a result of this method, the context managed by the DB is cancelled.
func (p *DBPanic) Panic(err error) {
	stackTrace := string(debug.Stack())

	p.logger.Errorf("LittDB encountered an unrecoverable error: %v\n%s", err, stackTrace)

	// only store the error if there isn't already an error stored
	p.error.CompareAndSwap(nil, &err)

	// Always cancel the context, even if this is not the first error. It's possible that the first "error" was
	// actually a shutdown request, and we want to make sure that the context is always cancelled in the event
	// of an unexpected error.
	p.cancel()
}
