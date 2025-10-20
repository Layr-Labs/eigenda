package store

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"

	_ "github.com/Layr-Labs/eigenda/api/clients/v2"
	_ "github.com/Layr-Labs/eigenda/api/clients/v2/payloadretrieval"
	"github.com/Layr-Labs/eigenda/api/proxy/common"
	"github.com/Layr-Labs/eigenda/api/proxy/common/types/certs"
	"github.com/Layr-Labs/eigenda/api/proxy/store/secondary"
	"github.com/Layr-Labs/eigensdk-go/logging"
)

//go:generate mockgen -package mocks --destination ../test/mocks/eigen_da_manager.go . IEigenDAManager

// IEigenDAManager handles EigenDA certificate operations
type IEigenDAManager interface {
	// See [EigenDAManager.Put]
	Put(ctx context.Context, value []byte) ([]byte, error)
	// See [EigenDAManager.Get]
	Get(ctx context.Context, versionedCert certs.VersionedCert, opts common.GETOpts) ([]byte, error)
	// See [EigenDAManager.SetDispersalBackend]
	SetDispersalBackend(backend common.EigenDABackend)
	// See [EigenDAManager.GetDispersalBackend]
	GetDispersalBackend() common.EigenDABackend
}

// EigenDAManager handles EigenDA certificate operations
type EigenDAManager struct {
	log logging.Logger

	eigenda          common.EigenDAV1Store // v0 version byte
	eigendaV2        common.EigenDAV2Store // >= v1 version bytes
	dispersalBackend atomic.Value          // stores the EigenDABackend to write blobs to

	// secondary storage backends (caching and fallbacks)
	secondary secondary.ISecondary
}

var _ IEigenDAManager = &EigenDAManager{}

// NewEigenDAManager creates a new EigenDAManager
func NewEigenDAManager(
	eigenda common.EigenDAV1Store,
	eigenDAV2 common.EigenDAV2Store,
	l logging.Logger,
	secondary secondary.ISecondary,
	dispersalBackend common.EigenDABackend,
) (*EigenDAManager, error) {
	// Enforce invariants
	if dispersalBackend == common.V2EigenDABackend && eigenDAV2 == nil {
		return nil, fmt.Errorf("EigenDA V2 dispersal enabled but no v2 store provided")
	}

	if dispersalBackend == common.V1EigenDABackend && eigenda == nil {
		return nil, fmt.Errorf("EigenDA dispersal enabled but no store provided")
	}

	manager := &EigenDAManager{
		log:       l,
		eigenda:   eigenda,
		eigendaV2: eigenDAV2,
		secondary: secondary,
	}
	manager.dispersalBackend.Store(dispersalBackend)
	return manager, nil
}

// GetDispersalBackend returns which EigenDA backend is currently being used for dispersal
func (m *EigenDAManager) GetDispersalBackend() common.EigenDABackend {
	val := m.dispersalBackend.Load()
	backend, ok := val.(common.EigenDABackend)
	if !ok {
		m.log.Error("Failed to convert dispersalBackend to EigenDABackend type", "value", val)
		return 0
	}
	return backend
}

// SetDispersalBackend sets which EigenDA backend to use for dispersal
func (m *EigenDAManager) SetDispersalBackend(backend common.EigenDABackend) {
	m.dispersalBackend.Store(backend)
}

// Get fetches a value from a storage backend based on the (commitment mode, type).
// It also validates the value retrieved and returns an error if the value is invalid.
// If opts.ReturnEncodedPayload is true, it will return the encoded payload without decoding it.
func (m *EigenDAManager) Get(ctx context.Context,
	versionedCert certs.VersionedCert,
	opts common.GETOpts,
) ([]byte, error) {
	switch versionedCert.Version {
	case certs.V0VersionByte:
		if m.eigenda == nil {
			return nil, errors.New("received CertV0 but EigenDA V1 client is not initialized")
		}
		return m.getEigenDAV1(ctx, versionedCert)
	case certs.V1VersionByte, certs.V2VersionByte:
		if m.eigendaV2 == nil {
			return nil, errors.New("received EigenDAV2 cert but EigenDA V2 client is not initialized")
		}
		return m.getEigenDAV2(ctx, versionedCert, opts)
	default:
		return nil, fmt.Errorf("cert version unknown: %b", versionedCert.Version)
	}
}

// getEigenDAV1 will attempt to retrieve a blob for the given versionedCert
// from cache, EigenDA V1, and fallback storage.
// TODO: we should also add the v1 RetrievalClient to retrieve from the validators directly
// in case the v1 disperser is down, the same way we do for v2.
func (m *EigenDAManager) getEigenDAV1(
	ctx context.Context,
	versionedCert certs.VersionedCert,
) ([]byte, error) {
	verifyFnForSecondary := func(ctx context.Context, cert []byte, payload []byte) error {
		// We don't add the cert version because EigenDA V1 only supports [certs.V0VersionByte] Certs.
		// We also don't use the l1InclusionBlockNumber because Recency check is only supported by EigenDA V2.
		// TODO: we should decouple the Verify function into a VerifyCert and VerifyBlob function,
		// and only verify the cert once here, before retrievals, and then only verify the blob commitment
		// after retrievals.
		return m.eigenda.Verify(ctx, cert, payload)
	}

	var readErrors []error
	// 1 - read payload from cache if enabled
	// Secondary storages (cache and fallback) store payloads instead of blobs.
	// TODO: would be nice to store blobs instead of payloads in secondary storages, such that we could standardize all
	// storages and make them all implement the [clients.PayloadRetriever] interface.
	// We could then get rid of the proxy notion of caches/fallbacks and only have storages.
	if m.secondary.CachingEnabled() {
		m.log.Debug("Retrieving payload from cached backends")
		payload, err := m.secondary.MultiSourceRead(ctx,
			versionedCert.SerializedCert, false, verifyFnForSecondary)
		if err == nil {
			return payload, nil
		}
		m.log.Warn("Failed to read payload from cache targets", "err", err)
		readErrors = append(readErrors, fmt.Errorf("read from cache targets: %w", err))
	}

	// 2 - read payload from EigenDA
	payload, err := m.eigenda.Get(ctx, versionedCert.SerializedCert)
	if err == nil {
		err = m.eigenda.Verify(ctx, versionedCert.SerializedCert, payload)
		if err != nil {
			return nil, fmt.Errorf("verify EigenDA V1 cert: %w", err)
		}
		if m.secondary.WriteOnCacheMissEnabled() {
			err = m.backupToSecondary(ctx, versionedCert.SerializedCert, payload)
			if err != nil {
				return nil, fmt.Errorf("backup to secondary on cache miss: %w", err)
			}
		}
		return payload, nil
	}
	readErrors = append(readErrors, fmt.Errorf("read from EigenDA backend: %w", err))

	// 3 - read blob from fallbacks if enabled and data is non-retrievable from EigenDA
	if m.secondary.FallbackEnabled() {
		payload, err = m.secondary.MultiSourceRead(ctx,
			versionedCert.SerializedCert, true, verifyFnForSecondary)
		if err == nil {
			return payload, nil
		}
		readErrors = append(readErrors, fmt.Errorf("read from fallback targets: %w", err))
	}

	return nil, fmt.Errorf("failed to read from all storage backends: %w", errors.Join(readErrors...))
}

// getEigenDAV2 will attempt to retrieve a blob for the given versionedCert
// from cache, EigenDA V2 relays, EigenDA V2 validators, and fallback storage.
func (m *EigenDAManager) getEigenDAV2(
	ctx context.Context,
	versionedCert certs.VersionedCert,
	opts common.GETOpts,
) ([]byte, error) {

	// The cert must be verified before attempting to get the data, since the GET logic
	// assumes the cert is valid. Verify v2 doesn't require a payload
	// because the payload is checked inside the Get function below.
	err := m.eigendaV2.VerifyCert(ctx, versionedCert, opts.L1InclusionBlockNum)
	if err != nil {
		return nil, fmt.Errorf("verify EigenDACert: %w", err)
	}

	verifyFnForSecondary := func(ctx context.Context, cert []byte, payload []byte) error {
		// This was previously using the VerifyCert function, which is pointless because it is now verified above,
		// and the cert only needs to be verified once.
		// TODO: implement a verify blob function, the same way it is implemented in [payloadretrieval.RelayPayloadRetriever]
		return nil
	}

	var readErrors []error
	// 1 - read payload from cache if enabled
	// Secondary storages (cache and fallback) store payloads instead of blobs.
	// For simplicity, we bypass secondary storages when requesting encoded payloads,
	// since those requests are only for secure integrations and run by provers/challengers.
	// TODO: would be nice to store blobs instead of payloads in secondary storages, such that we could standardize all
	// storages and make them all implement the [clients.PayloadRetriever] interface.
	// We could then get rid of the proxy notion of caches/fallbacks and only have storages.
	if m.secondary.CachingEnabled() && !opts.ReturnEncodedPayload {
		m.log.Debug("Retrieving payload from cached backends")
		payload, err := m.secondary.MultiSourceRead(ctx,
			versionedCert.SerializedCert, false, verifyFnForSecondary)
		if err == nil {
			return payload, nil
		}
		m.log.Warn("Failed to read payload from cache targets", "err", err)
		readErrors = append(readErrors, fmt.Errorf("read from cache targets: %w", err))
	}

	// 2 - read payloadOrEncodedPayload from EigenDA
	m.log.Debug("Reading blob from EigenDAV2 backend", "returnEncodedPayload", opts.ReturnEncodedPayload)
	payloadOrEncodedPayload, err := m.eigendaV2.Get(ctx, versionedCert, opts.ReturnEncodedPayload)
	if err == nil {
		// Only backup to secondary storage if we're returning the decoded payload
		// since the secondary stores are currently hardcoded to store payloads only.
		// TODO: we could consider also storing encoded payloads under separate keys?
		if m.secondary.WriteOnCacheMissEnabled() && !opts.ReturnEncodedPayload {
			err = m.backupToSecondary(ctx, versionedCert.SerializedCert, payloadOrEncodedPayload)
			if err != nil {
				return nil, fmt.Errorf("backup to secondary on cache miss: %w", err)
			}
		}
		return payloadOrEncodedPayload, nil
	}
	readErrors = append(readErrors, fmt.Errorf("read from EigenDA backend: %w", err))

	// 3 - read blob from fallbacks if enabled and data is non-retrievable from EigenDA
	// Only use fallbacks if we're not requesting encoded payload
	if m.secondary.FallbackEnabled() && !opts.ReturnEncodedPayload {
		payloadOrEncodedPayload, err = m.secondary.MultiSourceRead(ctx,
			versionedCert.SerializedCert, true, verifyFnForSecondary)
		if err == nil {
			return payloadOrEncodedPayload, nil
		}
		readErrors = append(readErrors, fmt.Errorf("read from fallback targets: %w", err))
	}

	return nil, fmt.Errorf("failed to read from all storage backends: %w", errors.Join(readErrors...))
}

// Put ... inserts a value into a storage backend based on the commitment mode
func (m *EigenDAManager) Put(ctx context.Context, value []byte) ([]byte, error) {
	var commit []byte
	var err error

	// 1 - Put blob into primary storage backend
	commit, err = m.putToCorrectEigenDABackend(ctx, value)
	if err != nil {
		return nil, err
	}

	// 2 - Put blob into secondary storage backends
	if m.secondary.Enabled() {
		err = m.backupToSecondary(ctx, commit, value)
		if err != nil {
			return nil, fmt.Errorf("backup to secondary storage: %w", err)
		}
	}

	return commit, nil
}

// putToCorrectEigenDABackend ... disperses blob to EigenDA backend
func (m *EigenDAManager) putToCorrectEigenDABackend(ctx context.Context, value []byte) ([]byte, error) {
	val := m.dispersalBackend.Load()
	backend, ok := val.(common.EigenDABackend)
	if !ok {
		return nil, fmt.Errorf("invalid dispersal backend type: %v", val)
	}

	if backend == common.V1EigenDABackend {
		if m.eigenda == nil {
			return nil, errors.New("EigenDA V1 dispersal requested but not configured")
		}
		return m.eigenda.Put(ctx, value) //nolint: wrapcheck
	}

	if backend == common.V2EigenDABackend {
		if m.eigendaV2 == nil {
			return nil, errors.New("EigenDA V2 dispersal requested but not configured")
		}
		return m.eigendaV2.Put(ctx, value) //nolint: wrapcheck
	}

	return nil, fmt.Errorf("unsupported dispersal backend: %v", backend)
}

// backupToSecondary writes data to secondary storage backends (caches and fallbacks).
// When errorOnInsertFailure is enabled and writes are synchronous, errors are returned
// to the caller to propagate as HTTP 500 responses. For async writes, errors are only logged.
func (m *EigenDAManager) backupToSecondary(ctx context.Context, commitment []byte, value []byte) error {
	if m.secondary.AsyncWriteEntry() { // publish put notification to secondary's subscription on PutNotify topic
		m.log.Debug("Publishing data to async secondary stores", "commitment", commitment)
		m.secondary.Topic() <- secondary.PutNotify{
			Commitment: commitment,
			Value:      value,
		}
		// Async writes cannot return errors to the client since they happen in background goroutines.
		// The configuration validation ensures errorOnInsertFailure is disabled when async mode is enabled.
		return nil
	}

	// Synchronous writes
	m.log.Debug("Publishing data to single threaded secondary stores")
	err := m.secondary.HandleRedundantWrites(ctx, commitment, value)
	if err != nil {
		m.log.Error("Secondary insertions failed", "error", err.Error())
		// Only propagate the error if errorOnInsertFailure is enabled.
		// This allows the caller to return HTTP 500 to the client.
		if m.secondary.ErrorOnInsertFailure() {
			return fmt.Errorf("a secondary storage write failed and error-on-secondary-insert-failure is enabled: %w", err)
		}
	}

	return nil
}
