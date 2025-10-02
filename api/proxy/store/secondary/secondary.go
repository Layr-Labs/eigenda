package secondary

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/Layr-Labs/eigenda/api/proxy/common"
	"github.com/Layr-Labs/eigenda/api/proxy/metrics"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/ethereum-optimism/optimism/op-service/retry"
	"github.com/ethereum/go-ethereum/crypto"
)

type MetricExpression = string

const (
	Miss    MetricExpression = "miss"
	Success MetricExpression = "success"
	Failed  MetricExpression = "failed"
)

type ISecondary interface {
	AsyncWriteEntry() bool
	Enabled() bool
	Topic() chan<- PutNotify
	CachingEnabled() bool
	FallbackEnabled() bool
	HandleRedundantWrites(ctx context.Context, commitment []byte, value []byte) error
	// verify fn signature has to match that of common/store.go's GeneratedKeyStore.Verify fn.
	MultiSourceRead(
		ctx context.Context, commitment []byte, fallback bool,
		verifyPayload func(context.Context, []byte, []byte) error,
	) ([]byte, error)
	WriteSubscriptionLoop(ctx context.Context)
	WriteOnCacheMissEnabled() bool
	ErrorOnInsertFailure() bool
}

// PutNotify ... notification received by primary manager to perform insertion across
// secondary storage backends
type PutNotify struct {
	Commitment []byte
	Value      []byte
}

// SecondaryManager ... routing abstraction for secondary storage backends
type SecondaryManager struct {
	log logging.Logger
	m   metrics.Metricer

	caches    []common.SecondaryStore
	fallbacks []common.SecondaryStore

	verifyLock           sync.RWMutex
	topic                chan PutNotify
	concurrentWrites     bool
	writeOnCacheMiss     bool
	errorOnInsertFailure bool
}

// NewSecondaryManager ... creates a new secondary storage manager
func NewSecondaryManager(
	log logging.Logger,
	m metrics.Metricer,
	caches []common.SecondaryStore,
	fallbacks []common.SecondaryStore,
	writeOnCacheMiss bool,
	errorOnInsertFailure bool,
) ISecondary {
	return &SecondaryManager{
		topic: make(
			chan PutNotify,
		), // channel is un-buffered which dispersing consumption across routines helps alleviate
		log:                  log,
		m:                    m,
		caches:               caches,
		fallbacks:            fallbacks,
		verifyLock:           sync.RWMutex{},
		writeOnCacheMiss:     writeOnCacheMiss,
		errorOnInsertFailure: errorOnInsertFailure,
	}
}

// Topic ...
func (sm *SecondaryManager) Topic() chan<- PutNotify {
	return sm.topic
}

func (sm *SecondaryManager) Enabled() bool {
	return sm.CachingEnabled() || sm.FallbackEnabled()
}

func (sm *SecondaryManager) CachingEnabled() bool {
	return len(sm.caches) > 0
}

func (sm *SecondaryManager) FallbackEnabled() bool {
	return len(sm.fallbacks) > 0
}

func (sm *SecondaryManager) WriteOnCacheMissEnabled() bool {
	return sm.CachingEnabled() && sm.writeOnCacheMiss
}

// ErrorOnInsertFailure returns whether secondary insertion failures should be returned as errors
// to the client, rather than being silently logged.
func (sm *SecondaryManager) ErrorOnInsertFailure() bool {
	return sm.errorOnInsertFailure
}

// HandleRedundantWrites ... writes to both sets of backends (i.e, fallback, cache)
// and returns an error based on the errorOnInsertFailure configuration:
//   - If errorOnInsertFailure is false (default): returns error only if ALL writes fail
//   - If errorOnInsertFailure is true: returns error if ANY write fails
func (sm *SecondaryManager) HandleRedundantWrites(ctx context.Context, commitment []byte, value []byte) error {
	sources := sm.caches
	sources = append(sources, sm.fallbacks...)

	key := crypto.Keccak256(commitment)
	successes := 0
	var errs []error

	for _, src := range sources {
		sm.log.Debug("Attempting to write to secondary storage", "backend", src.BackendType())
		cb := sm.m.RecordSecondaryRequest(src.BackendType().String(), http.MethodPut)

		// for added safety - we retry the insertion 5x using a default exponential backoff
		_, err := retry.Do[any](ctx, 5, retry.Exponential(),
			func() (any, error) {
				return 0, src.Put(
					ctx,
					key,
					value,
				) // this implementation assumes that all secondary clients are thread safe
			})
		if err != nil {
			sm.log.Warn("Failed to write to redundant target", "backend", src.BackendType(), "err", err)
			cb(Failed)
			errs = append(errs, fmt.Errorf("write to %s failed: %w", src.BackendType(), err))
		} else {
			successes++
			cb(Success)
		}
	}

	// If no writes succeeded at all, always return error
	if successes == 0 {
		return fmt.Errorf("failed to write blob to any redundant targets: %w", errors.Join(errs...))
	}

	// If errorOnInsertFailure is enabled and any writes failed (partial success), return error
	if sm.errorOnInsertFailure && len(errs) > 0 {
		return fmt.Errorf("failed to write to %d of %d secondary targets (error-on-secondary-insert-failure=true): %w",
			len(errs), len(sources), errors.Join(errs...))
	}

	return nil
}

// AsyncWriteEntry ... subscribes to put notifications posted to shared topic with primary manager
func (sm *SecondaryManager) AsyncWriteEntry() bool {
	return sm.concurrentWrites
}

// WriteSubscriptionLoop ... subscribes to put notifications posted to shared topic with primary manager
func (sm *SecondaryManager) WriteSubscriptionLoop(ctx context.Context) {
	sm.concurrentWrites = true

	for {
		select {
		case notif := <-sm.topic:
			err := sm.HandleRedundantWrites(context.Background(), notif.Commitment, notif.Value)
			if err != nil {
				sm.log.Error("Failed to write to redundant targets", "err", err)
			}

		case <-ctx.Done():
			sm.log.Debug("Terminating secondary event loop")
			return
		}
	}
}

// MultiSourceRead ... reads from a set of backends and returns the first successfully read blob
// NOTE: - this can also be parallelized when reading from multiple sources and discarding connections that fail
// - for complete optimization we can profile secondary storage backends to determine the fastest / most reliable and
// always route to it first
func (sm *SecondaryManager) MultiSourceRead(
	ctx context.Context,
	commitment []byte,
	fallback bool,
	verifyPayload func(context.Context, []byte, []byte) error,
) ([]byte, error) {
	var sources []common.SecondaryStore
	if fallback {
		sources = sm.fallbacks
	} else {
		sources = sm.caches
	}

	key := crypto.Keccak256(commitment)
	for _, src := range sources {
		cb := sm.m.RecordSecondaryRequest(src.BackendType().String(), http.MethodGet)
		data, err := src.Get(ctx, key)
		if err != nil {
			cb(Failed)
			sm.log.Warn("Failed to read from redundant target", "backend", src.BackendType(), "err", err)
			continue
		}

		if data == nil {
			cb(Miss)
			sm.log.Debug("No data found in redundant target", "backend", src.BackendType())
			continue
		}

		// verify cert:data using provided verification function
		sm.verifyLock.Lock()
		err = verifyPayload(ctx, commitment, data)
		if err != nil {
			cb(Failed)
			sm.log.Warn("Failed to verify blob", "err", err, "backend", src.BackendType())
			sm.verifyLock.Unlock()
			continue
		}
		sm.verifyLock.Unlock()
		cb(Success)
		return data, nil
	}
	return nil, errors.New("no data found in any redundant backend")
}
