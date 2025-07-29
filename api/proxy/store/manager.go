package store

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"

	_ "github.com/Layr-Labs/eigenda/api/clients/v2"
	"github.com/Layr-Labs/eigenda/api/proxy/common"
	"github.com/Layr-Labs/eigenda/api/proxy/common/types/certs"
	"github.com/Layr-Labs/eigenda/api/proxy/common/types/commitments"
	"github.com/Layr-Labs/eigenda/api/proxy/store/secondary"
	"github.com/Layr-Labs/eigenda/api/proxy/store/secondary/s3"
	"github.com/Layr-Labs/eigensdk-go/logging"
)

//go:generate mockgen -package mocks --destination ../test/mocks/manager.go . IManager

// IManager ... read/write interface
type IManager interface {
	// See [Manager.Put]
	Put(ctx context.Context, cm commitments.CommitmentMode, value []byte) ([]byte, error)
	// See [Manager.Get]
	Get(ctx context.Context, versionedCert certs.VersionedCert,
		cm commitments.CommitmentMode, opts common.GETOpts) ([]byte, error)
	// See [Manager.SetDispersalBackend]
	SetDispersalBackend(backend common.EigenDABackend)
	// See [Manager.GetDispersalBackend]
	GetDispersalBackend() common.EigenDABackend
	// See [Manager.PutOPKeccakPairInS3]
	PutOPKeccakPairInS3(ctx context.Context, key []byte, value []byte) error
	// See [Manager.GetOPKeccakValueFromS3]
	GetOPKeccakValueFromS3(ctx context.Context, key []byte) ([]byte, error)
}

// Manager ... storage backend routing layer
type Manager struct {
	log logging.Logger

	s3 *s3.Store // for op keccak256 commitment
	// For op generic commitments & standard commitments
	eigenda          common.EigenDAV1Store // v0 version byte
	eigendaV2        common.EigenDAV2Store // >= v1 version bytes
	dispersalBackend atomic.Value          // stores the EigenDABackend to write blobs to

	// secondary storage backends (caching and fallbacks)
	secondary secondary.ISecondary
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
	eigenda common.EigenDAV1Store,
	eigenDAV2 common.EigenDAV2Store,
	s3 *s3.Store,
	l logging.Logger,
	secondary secondary.ISecondary,
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

// Get fetches a value from a storage backend based on the (commitment mode, type).
// It also validates the value retrieved and returns an error if the value is invalid.
// If opts.ReturnEncodedPayload is true, it will return the encoded payload without decoding it.
func (m *Manager) Get(ctx context.Context,
	versionedCert certs.VersionedCert,
	cm commitments.CommitmentMode,
	opts common.GETOpts,
) ([]byte, error) {
	switch cm {
	case commitments.StandardCommitmentMode, commitments.OptimismGenericCommitmentMode:
		if versionedCert.Version == certs.V0VersionByte && m.eigenda == nil {
			return nil, errors.New("expected EigenDA V1 backend for DA commitment type with CertV0")
		}
		if versionedCert.Version == certs.V1VersionByte && m.eigendaV2 == nil {
			return nil, errors.New("expected EigenDA V2 backend for DA commitment type with CertV1")
		}

		// The eigendav2 verify function verifies the cert, but not the kzg commitment.
		// TODO: should we be verifying those when retrieving from secondary storages like S3...?
		verifyMethod, err := m.getVerifyMethod(versionedCert.Version)
		if err != nil {
			return nil, fmt.Errorf("get verify method: %w", err)
		}

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
				versionedCert.SerializedCert, false, verifyMethod, opts.L1InclusionBlockNum)
			if err == nil {
				return payload, nil
			}
			m.log.Warn("Failed to read payload from cache targets", "err", err)
		}

		// 2 - read maybeEncodedPayload from EigenDA
		maybeEncodedPayload, err := m.getFromCorrectEigenDABackend(ctx, versionedCert, opts)
		if err == nil {
			// Only backup to secondary storage if we're returning the decoded payload
			// since the secondary stores are currently hardcoded to store payloads only.
			// TODO: we could consider also storing encoded payloads under separate keys?
			if m.secondary.WriteOnCacheMissEnabled() && !opts.ReturnEncodedPayload {
				m.backupToSecondary(ctx, versionedCert.SerializedCert, maybeEncodedPayload)
			}

			return maybeEncodedPayload, nil
		}

		// 3 - read blob from fallbacks if enabled and data is non-retrievable from EigenDA
		// Only use fallbacks if we're not requesting encoded payload
		if m.secondary.FallbackEnabled() && !opts.ReturnEncodedPayload {
			maybeEncodedPayload, err = m.secondary.MultiSourceRead(ctx,
				versionedCert.SerializedCert, true, verifyMethod, opts.L1InclusionBlockNum)
			if err != nil {
				m.log.Error("Failed to read payload from fallback targets", "err", err)
				return nil, err
			}
		} else {
			return nil, err
		}
		return maybeEncodedPayload, err
	case commitments.OptimismKeccakCommitmentMode:
		// TODO: we should refactor the manager to not deal with keccak commitments at all.
		return nil, fmt.Errorf("INTERNAL BUG: call GetOPKeccakValueFromS3 instead")
	default:
		return nil, errors.New("could not determine which storage backend to route to based on unknown commitment mode")
	}
}

// Put ... inserts a value into a storage backend based on the commitment mode
func (m *Manager) Put(ctx context.Context, cm commitments.CommitmentMode, value []byte) ([]byte, error) {
	var commit []byte
	var err error

	// 1 - Put blob into primary storage backend
	switch cm {
	case commitments.OptimismGenericCommitmentMode, commitments.StandardCommitmentMode:
		commit, err = m.putToCorrectEigenDABackend(ctx, value)
		if err != nil {
			return nil, err
		}
	case commitments.OptimismKeccakCommitmentMode:
		// TODO: we should refactor the manager to not deal with keccak commitments at all.
		return nil, fmt.Errorf("INTERNAL BUG: call PutOPKeccakPairInS3 instead")
	default:
		return nil, fmt.Errorf("unknown commitment mode")
	}

	// 2 - Put blob into secondary storage backends
	if m.secondary.Enabled() {
		m.backupToSecondary(ctx, commit, value)
	}

	return commit, nil
}

func (m *Manager) backupToSecondary(ctx context.Context, commitment []byte, value []byte) {
	if m.secondary.AsyncWriteEntry() { // publish put notification to secondary's subscription on PutNotify topic
		m.log.Debug("Publishing data to async secondary stores")
		m.secondary.Topic() <- secondary.PutNotify{
			Commitment: commitment,
			Value:      value,
		}
		// secondary is available only for synchronous writes
	} else {
		m.log.Debug("Publishing data to single threaded secondary stores")
		err := m.secondary.HandleRedundantWrites(ctx, commitment, value)
		if err != nil {
			m.log.Error("Secondary insertions failed", "error", err.Error())
		}
	}
}

// getVerifyMethod returns the correct verify method based on commitment type
func (m *Manager) getVerifyMethod(certVersion certs.VersionByte) (
	func(context.Context, []byte, []byte, uint64) error,
	error,
) {
	v1VerifyWrapper := func(ctx context.Context, cert []byte, payload []byte, l1InclusionBlockNumber uint64) error {
		// we don't add the cert version because EigenDA V1 only supported [certs.V0VersionByte] Certs.
		return m.eigenda.Verify(ctx, cert, payload)
	}
	v2VerifyWrapper := func(ctx context.Context, cert []byte, payload []byte, l1InclusionBlockNumber uint64) error {
		return m.eigendaV2.VerifyCert(ctx, certs.NewVersionedCert(cert, certVersion), l1InclusionBlockNumber)
	}

	switch certVersion {
	case certs.V0VersionByte:
		return v1VerifyWrapper, nil
	case certs.V1VersionByte, certs.V2VersionByte:
		return v2VerifyWrapper, nil
	default:
		return nil, fmt.Errorf("commitment version unknown: %b", certVersion)
	}
}

// putToCorrectEigenDABackend ... disperses blob to EigenDA backend
func (m *Manager) putToCorrectEigenDABackend(ctx context.Context, value []byte) ([]byte, error) {
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
	opts common.GETOpts,
) ([]byte, error) {
	switch versionedCert.Version {
	case certs.V0VersionByte:
		m.log.Debug("Reading blob from EigenDAV1 backend")

		// V1 backend does not support returning encoded payloads
		if opts.ReturnEncodedPayload {
			return nil, fmt.Errorf("returning encoded payload is not supported for V0 certificates")
		}

		data, err := m.eigenda.Get(ctx, versionedCert.SerializedCert)
		if err == nil {
			err = m.eigenda.Verify(ctx, versionedCert.SerializedCert, data)
			if err != nil {
				return nil, err
			}
			return data, nil
		}

		return nil, err

	case certs.V1VersionByte, certs.V2VersionByte:
		// The cert must be verified before attempting to get the data, since the GET logic
		// assumes the cert is valid. Verify v2 doesn't require a payload
		// because the payload is checked inside the Get function below.
		err := m.eigendaV2.VerifyCert(ctx, versionedCert, opts.L1InclusionBlockNum)
		if err != nil {
			return nil, fmt.Errorf("verify EigenDACert: %w", err)
		}

		m.log.Debug("Reading blob from EigenDAV2 backend", "returnEncodedPayload", opts.ReturnEncodedPayload)
		maybeEncodedPayload, err := m.eigendaV2.Get(ctx, versionedCert, opts.ReturnEncodedPayload)
		if err != nil {
			return nil, fmt.Errorf("get payload from EigenDA V2 backend: %w", err)
		}

		return maybeEncodedPayload, nil
	default:
		return nil, fmt.Errorf("cert version unknown: %b", versionedCert.Version)
	}
}

// PutOPKeccakPairInS3 puts a key/value pair, where key=keccak(value), into S3.
// If key!=keccak(value), a Keccak256KeyValueMismatchError is returned.
// This is only used for OP keccak256 commitments.
func (m *Manager) PutOPKeccakPairInS3(ctx context.Context, key []byte, value []byte) error {
	if m.s3 == nil {
		return errors.New("S3 is disabled but is only supported for posting known commitment keys")
	}
	err := m.s3.Verify(ctx, key, value)
	if err != nil {
		return fmt.Errorf("s3 verify: %w", err)
	}
	err = m.s3.Put(ctx, key, value)
	if err != nil {
		return fmt.Errorf("s3 put: %w", err)
	}
	return nil
}

// GetOPKeccakValueFromS3 retrieves the value associated with the given key from S3.
// It verifies that the key=keccak(value) and returns an error if they don't match.
// Otherwise returns the value and nil.
func (m *Manager) GetOPKeccakValueFromS3(ctx context.Context, key []byte) ([]byte, error) {
	if m.s3 == nil {
		return nil, errors.New("expected S3 backend for OP keccak256 commitment type, but none configured")
	}

	// 1 - read blob from S3 backend
	m.log.Debug("Retrieving data from S3 backend")
	value, err := m.s3.Get(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("s3 get: %w", err)
	}

	// 2 - verify payload hash against commitment key digest
	err = m.s3.Verify(ctx, key, value)
	if err != nil {
		return nil, fmt.Errorf("s3 verify: %w", err)
	}
	return value, nil
}
