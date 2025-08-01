package memconfig

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients/v2/coretypes"
)

// Config contains properties that are used to configure the MemStore's behavior.
type Config struct {
	MaxBlobSizeBytes uint64
	BlobExpiration   time.Duration
	// artificial latency added for memstore backend to mimic eigenda's latency
	PutLatency time.Duration
	GetLatency time.Duration
	// when true, put requests will return an errorFailover error,
	// after sleeping PutLatency duration.
	// This can be used to simulate eigenda being down.
	PutReturnsFailoverError bool
	// if nil, store data from the POST method in the empheral db
	// if it is set to some derivation error, then the derivation error is stored as opposed
	// to the data from the POST method in the empheral db
	// TODO we use Put in the name to be consistent to the name "PutReturnsFailoverError",
	// but they should have been named as Post from HTTP verb
	OverwritePutWithDerivationError error
}

// MarshalJSON implements custom JSON marshaling for Config.
// This is needed because time.Duration is serialized to nanoseconds,
// which is hard to read.
// We only implement Marshal and not Unmarshal because this is only needed
// for the GET /memstore/config endpoint, which only reads the configuration.
// Patches are reads as ConfigUpdates instead to handle omitted fields.
func (c Config) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		MaxBlobSizeBytes                uint64
		BlobExpiration                  string
		PutLatency                      string
		GetLatency                      string
		PutReturnsFailoverError         bool
		OverwritePutWithDerivationError error
	}{
		MaxBlobSizeBytes:                c.MaxBlobSizeBytes,
		BlobExpiration:                  c.BlobExpiration.String(),
		PutLatency:                      c.PutLatency.String(),
		GetLatency:                      c.GetLatency.String(),
		PutReturnsFailoverError:         c.PutReturnsFailoverError,
		OverwritePutWithDerivationError: c.OverwritePutWithDerivationError,
	})
}

// SafeConfig handles thread-safe access to Config.
// It is uses by MemStore to read configuration values.
// and by the MemStore API to update configuration values.
type SafeConfig struct {
	mu     sync.RWMutex
	config Config
}

// Need this because we marshal the entire proxy config on startup
// to log it, and private fields are not marshalled.
func (sc *SafeConfig) MarshalJSON() ([]byte, error) {
	sc.mu.RLock()
	defer sc.mu.RUnlock()
	return json.Marshal(sc.config)
}

func NewSafeConfig(config Config) *SafeConfig {
	return &SafeConfig{
		config: config,
	}
}

func (sc *SafeConfig) LatencyPUTRoute() time.Duration {
	sc.mu.RLock()
	defer sc.mu.RUnlock()
	return sc.config.PutLatency
}
func (sc *SafeConfig) SetLatencyPUTRoute(latency time.Duration) {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	sc.config.PutLatency = latency
}

func (sc *SafeConfig) LatencyGETRoute() time.Duration {
	sc.mu.RLock()
	defer sc.mu.RUnlock()
	return sc.config.GetLatency
}
func (sc *SafeConfig) SetLatencyGETRoute(latency time.Duration) {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	sc.config.GetLatency = latency
}

func (sc *SafeConfig) PutReturnsFailoverError() bool {
	sc.mu.RLock()
	defer sc.mu.RUnlock()
	return sc.config.PutReturnsFailoverError
}
func (sc *SafeConfig) SetPUTReturnsFailoverError(returnsFailoverError bool) {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	sc.config.PutReturnsFailoverError = returnsFailoverError
}

func (sc *SafeConfig) BlobExpiration() time.Duration {
	sc.mu.RLock()
	defer sc.mu.RUnlock()
	return sc.config.BlobExpiration
}
func (sc *SafeConfig) SetBlobExpiration(expiration time.Duration) {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	sc.config.BlobExpiration = expiration
}

func (sc *SafeConfig) MaxBlobSizeBytes() uint64 {
	sc.mu.RLock()
	defer sc.mu.RUnlock()
	return sc.config.MaxBlobSizeBytes
}
func (sc *SafeConfig) SetMaxBlobSizeBytes(maxBlobSizeBytes uint64) {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	sc.config.MaxBlobSizeBytes = maxBlobSizeBytes
}

func (sc *SafeConfig) OverwritePutWithDerivationError() error {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	return sc.config.OverwritePutWithDerivationError
}

func (sc *SafeConfig) SetOverwritePutWithDerivationError(inputError error) error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	// both dynamic type and value are nil, i.e there is no error
	if inputError == nil {
		sc.config.OverwritePutWithDerivationError = nil
		return nil
	}

	// cast into an DerivationError
	var derivationError coretypes.DerivationError
	if !errors.As(inputError, &derivationError) {
		return fmt.Errorf("unable to cast error into an DerivationError: %w", inputError)
	}

	derivationError.Validate()

	sc.config.OverwritePutWithDerivationError = derivationError

	return nil
}

func (sc *SafeConfig) Config() Config {
	sc.mu.RLock()
	defer sc.mu.RUnlock()
	return sc.config
}

func (sc *SafeConfig) Update(config Config) {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	sc.config = config
}
