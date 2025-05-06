//go:generate mockgen -package mocks --destination ../mocks/manager.go . IManager

package store

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"

	"github.com/Layr-Labs/eigenda-proxy/common"
	"github.com/Layr-Labs/eigenda-proxy/common/types/certs"
	"github.com/Layr-Labs/eigenda-proxy/common/types/commitments"
	"github.com/Layr-Labs/eigensdk-go/logging"
)

// IManager ... read/write interface
type IManager interface {
	Get(ctx context.Context, versionedCert certs.VersionedCert, cm commitments.CommitmentMode) ([]byte, error)
	Put(ctx context.Context, cm commitments.CommitmentMode, key, value []byte) ([]byte, error)
	SetDispersalBackend(backend common.EigenDABackend)
	GetDispersalBackend() common.EigenDABackend
}

// Manager ... storage backend routing layer
type Manager struct {
	log logging.Logger

	s3 common.PrecomputedKeyStore // for op keccak256 commitment
	// For op generic commitments & standard commitments
	eigenda          common.EigenDAStore // v0 da commitment version
	eigendaV2        common.EigenDAStore // v1 da commitment version
	dispersalBackend atomic.Value        // stores the EigenDABackend to write blobs to

	// secondary storage backends (caching and fallbacks)
	secondary ISecondary
}

var _ IManager = &Manager{}

// GetDispersalBackend returns which EigenDA backend is currently being used for dispersal
func (m *Manager) GetDispersalBackend() common.EigenDABackend {
	val := m.dispersalBackend.Load()
	backend, ok := val.(common.EigenDABackend)
	if !ok {
		m.log.Error("Failed to convert dispersalBackend to EigenDABackend type", "value", val)
		return 0
	}
	return backend
}

// SetDispersalBackend sets which EigenDA backend to use for dispersal
func (m *Manager) SetDispersalBackend(backend common.EigenDABackend) {
	m.dispersalBackend.Store(backend)
}

// NewManager ... Init
func NewManager(
	eigenda common.EigenDAStore,
	eigenDAV2 common.EigenDAStore,
	s3 common.PrecomputedKeyStore,
	l logging.Logger,
	secondary ISecondary,
	dispersalBackend common.EigenDABackend,
) (*Manager, error) {
	// Enforce invariants
	if dispersalBackend == common.V2EigenDABackend && eigenDAV2 == nil {
		return nil, fmt.Errorf("EigenDA V2 dispersal enabled but no v2 store provided")
	}

	if dispersalBackend == common.V1EigenDABackend && eigenda == nil {
		return nil, fmt.Errorf("EigenDA dispersal enabled but no store provided")
	}

	manager := &Manager{
		log:       l,
		eigenda:   eigenda,
		eigendaV2: eigenDAV2,
		s3:        s3,
		secondary: secondary,
	}
	manager.dispersalBackend.Store(dispersalBackend)
	return manager, nil
}

// Get ... fetches a value from a storage backend based on the (commitment mode, type)
func (m *Manager) Get(ctx context.Context,
	versionedCert certs.VersionedCert,
	cm commitments.CommitmentMode,
) ([]byte, error) {
	switch cm {
	case commitments.OptimismKeccakCommitmentMode:
		if m.s3 == nil {
			return nil, errors.New("expected S3 backend for OP keccak256 commitment type, but none configured")
		}

		// 1 - read blob from S3 backend
		m.log.Debug("Retrieving data from S3 backend")
		// Using only the serialized cert without the version byte, since that is how we originally
		// implemented this feature before we had versioned certs, and need to remain backwards compatible.
		value, err := m.s3.Get(ctx, versionedCert.SerializedCert)
		if err != nil {
			return nil, err
		}

		// 2 - verify blob hash against commitment key digest
		err = m.s3.Verify(ctx, versionedCert.SerializedCert, value)
		if err != nil {
			return nil, err
		}
		return value, nil

	case commitments.StandardCommitmentMode, commitments.OptimismGenericCommitmentMode:
		if versionedCert.Version == certs.V0VersionByte && m.eigenda == nil {
			return nil, errors.New("expected EigenDA V1 backend for DA commitment type with CertV0")
		}
		if versionedCert.Version == certs.V1VersionByte && m.eigendaV2 == nil {
			return nil, errors.New("expected EigenDA V2 backend for DA commitment type with CertV1")
		}

		verifyMethod, err := m.getVerifyMethod(versionedCert.Version)
		if err != nil {
			return nil, fmt.Errorf("get verify method: %w", err)
		}

		// 1 - read blob from cache if enabled
		if m.secondary.CachingEnabled() {
			m.log.Debug("Retrieving data from cached backends")
			data, err := m.secondary.MultiSourceRead(ctx, versionedCert.SerializedCert, false, verifyMethod)
			if err == nil {
				return data, nil
			}

			m.log.Warn("Failed to read from cache targets", "err", err)
		}

		// 2 - read blob from EigenDA
		data, err := m.getFromCorrectEigenDABackend(ctx, versionedCert)
		if err == nil {
			return data, nil
		}

		m.log.Error(err.Error())

		// 3 - read blob from fallbacks if enabled and data is non-retrievable from EigenDA
		if m.secondary.FallbackEnabled() {
			data, err = m.secondary.MultiSourceRead(ctx, versionedCert.SerializedCert, true, verifyMethod)
			if err != nil {
				m.log.Error("Failed to read from fallback targets", "err", err)
				return nil, err
			}
		} else {
			return nil, err
		}
		return data, err

	default:
		return nil, errors.New("could not determine which storage backend to route to based on unknown commitment mode")
	}
}

// Put ... inserts a value into a storage backend based on the commitment mode
func (m *Manager) Put(ctx context.Context, cm commitments.CommitmentMode, key, value []byte) ([]byte, error) {
	var commit []byte
	var err error

	// 1 - Put blob into primary storage backend
	switch cm {
	case commitments.OptimismKeccakCommitmentMode: // caching and fallbacks are unsupported for this commitment mode
		return m.putKeccak256Mode(ctx, key, value)
	case commitments.OptimismGenericCommitmentMode, commitments.StandardCommitmentMode:
		commit, err = m.putEigenDAMode(ctx, value)
	default:
		return nil, fmt.Errorf("unknown commitment mode")
	}

	if err != nil {
		return nil, err
	}

	// 2 - Put blob into secondary storage backends
	if m.secondary.Enabled() &&
		m.secondary.AsyncWriteEntry() { // publish put notification to secondary's subscription on PutNotify topic
		m.log.Debug("Publishing data to async secondary stores")
		m.secondary.Topic() <- PutNotify{
			Commitment: commit,
			Value:      value,
		}
		// secondary is available only for synchronous writes
	} else if m.secondary.Enabled() && !m.secondary.AsyncWriteEntry() {
		m.log.Debug("Publishing data to single threaded secondary stores")
		err := m.secondary.HandleRedundantWrites(ctx, commit, value)
		if err != nil {
			m.log.Error("Secondary insertions failed", "error", err.Error())
		}
	}

	return commit, nil
}

// getVerifyMethod returns the correct verify method based on commitment type
func (m *Manager) getVerifyMethod(commitmentType certs.VersionByte) (
	func(context.Context, []byte, []byte) error,
	error,
) {
	switch commitmentType {
	case certs.V0VersionByte:
		return m.eigenda.Verify, nil
	case certs.V1VersionByte:
		return m.eigendaV2.Verify, nil
	default:
		return nil, fmt.Errorf("commitment version unknown: %b", commitmentType)
	}
}

// putEigenDAMode ... disperses blob to EigenDA backend
func (m *Manager) putEigenDAMode(ctx context.Context, value []byte) ([]byte, error) {
	val := m.dispersalBackend.Load()
	backend, ok := val.(common.EigenDABackend)
	if !ok {
		return nil, fmt.Errorf("invalid dispersal backend type: %v", val)
	}

	if backend == common.V1EigenDABackend {
		if m.eigenda == nil {
			return nil, errors.New("EigenDA V1 dispersal requested but not configured")
		}
		return m.eigenda.Put(ctx, value)
	}

	if backend == common.V2EigenDABackend {
		if m.eigendaV2 == nil {
			return nil, errors.New("EigenDA V2 dispersal requested but not configured")
		}
		return m.eigendaV2.Put(ctx, value)
	}

	return nil, fmt.Errorf("unsupported dispersal backend: %v", backend)
}

func (m *Manager) getFromCorrectEigenDABackend(
	ctx context.Context,
	versionedCert certs.VersionedCert,
) ([]byte, error) {
	switch versionedCert.Version {
	case certs.V0VersionByte:
		m.log.Debug("Reading blob from EigenDAV1 backend")
		data, err := m.eigenda.Get(ctx, versionedCert.SerializedCert)
		if err == nil {
			// verify v1 (payload, cert)
			err = m.eigenda.Verify(ctx, versionedCert.SerializedCert, data)
			if err != nil {
				return nil, err
			}
			return data, nil
		}

		return nil, err

	case certs.V1VersionByte:
		// The cert must be verified before attempting to get the data, since the GET logic
		// assumes the cert is valid. Verify v2 doesn't require a payload.
		err := m.eigendaV2.Verify(ctx, versionedCert.SerializedCert, nil)
		if err != nil {
			return nil, fmt.Errorf("verify EigenDACert: %w", err)
		}

		m.log.Debug("Reading blob from EigenDAV2 backend")
		data, err := m.eigendaV2.Get(ctx, versionedCert.SerializedCert)
		if err != nil {
			return nil, fmt.Errorf("get data from V2 backend: %w", err)
		}

		return data, nil
	default:
		return nil, fmt.Errorf("cert version unknown: %b", versionedCert.Version)
	}
}

// putKeccak256Mode ... put blob into S3 compatible backend
func (m *Manager) putKeccak256Mode(ctx context.Context, key []byte, value []byte) ([]byte, error) {
	if m.s3 == nil {
		return nil, errors.New("S3 is disabled but is only supported for posting known commitment keys")
	}

	err := m.s3.Verify(ctx, key, value)
	if err != nil {
		return nil, err
	}

	return key, m.s3.Put(ctx, key, value)
}
