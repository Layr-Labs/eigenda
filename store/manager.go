//go:generate mockgen -package mocks --destination ../mocks/manager.go . IManager

package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/Layr-Labs/eigenda-proxy/commitments"
	"github.com/Layr-Labs/eigenda-proxy/common"
	"github.com/Layr-Labs/eigensdk-go/logging"
)

// IManager ... read/write interface
type IManager interface {
	Get(ctx context.Context, key []byte, cm commitments.CommitmentMode) ([]byte, error)
	Put(ctx context.Context, cm commitments.CommitmentMode, key, value []byte) ([]byte, error)
}

// Manager ... storage backend routing layer
type Manager struct {
	log logging.Logger
	// primary storage backends
	eigenda common.GeneratedKeyStore   // ALT DA commitment type for OP mode && std commitment mode for standard /client
	s3      common.PrecomputedKeyStore // OP commitment mode && keccak256 commitment type

	// secondary storage backends (caching and fallbacks)
	secondary ISecondary
}

// NewManager ... Init
func NewManager(eigenda common.GeneratedKeyStore, s3 common.PrecomputedKeyStore, l logging.Logger,
	secondary ISecondary) (IManager, error) {
	return &Manager{
		log:       l,
		eigenda:   eigenda,
		s3:        s3,
		secondary: secondary,
	}, nil
}

// Get ... fetches a value from a storage backend based on the (commitment mode, type)
func (m *Manager) Get(ctx context.Context, key []byte, cm commitments.CommitmentMode) ([]byte, error) {
	switch cm {
	case commitments.OptimismKeccak:

		if m.s3 == nil {
			return nil, errors.New("expected S3 backend for OP keccak256 commitment type, but none configured")
		}

		// 1 - read blob from S3 backend
		m.log.Debug("Retrieving data from S3 backend")
		value, err := m.s3.Get(ctx, key)
		if err != nil {
			return nil, err
		}

		// 2 - verify blob hash against commitment key digest
		err = m.s3.Verify(ctx, key, value)
		if err != nil {
			return nil, err
		}
		return value, nil

	case commitments.Standard, commitments.OptimismGeneric:
		if m.eigenda == nil {
			return nil, errors.New("expected EigenDA backend for DA commitment type, but none configured")
		}

		// 1 - read blob from cache if enabled
		if m.secondary.CachingEnabled() {
			m.log.Debug("Retrieving data from cached backends")
			data, err := m.secondary.MultiSourceRead(ctx, key, false, m.eigenda.Verify)
			if err == nil {
				return data, nil
			}

			m.log.Warn("Failed to read from cache targets", "err", err)
		}

		// 2 - read blob from EigenDA
		data, err := m.eigenda.Get(ctx, key)
		if err == nil {
			// verify
			err = m.eigenda.Verify(ctx, key, data)
			if err != nil {
				return nil, err
			}
			return data, nil
		}

		// 3 - read blob from fallbacks if enabled and data is non-retrievable from EigenDA
		if m.secondary.FallbackEnabled() {
			data, err = m.secondary.MultiSourceRead(ctx, key, true, m.eigenda.Verify)
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
	case commitments.OptimismKeccak: // caching and fallbacks are unsupported for this commitment mode
		return m.putKeccak256Mode(ctx, key, value)
	case commitments.OptimismGeneric, commitments.Standard:
		commit, err = m.putEigenDAMode(ctx, value)
	default:
		return nil, fmt.Errorf("unknown commitment mode")
	}

	if err != nil {
		return nil, err
	}

	// 2 - Put blob into secondary storage backends
	if m.secondary.Enabled() && m.secondary.AsyncWriteEntry() { // publish put notification to secondary's subscription on PutNotify topic
		m.log.Debug("Publishing data to async secondary stores")
		m.secondary.Topic() <- PutNotify{
			Commitment: commit,
			Value:      value,
		}
	} else if m.secondary.Enabled() && !m.secondary.AsyncWriteEntry() { // secondary is available only for synchronous writes
		m.log.Debug("Publishing data to single threaded secondary stores")
		err := m.secondary.HandleRedundantWrites(ctx, commit, value)
		if err != nil {
			m.log.Error("Secondary insertions failed", "error", err.Error())
		}
	}

	return commit, nil
}

// putEigenDAMode ... disperses blob to EigenDA backend
func (m *Manager) putEigenDAMode(ctx context.Context, value []byte) ([]byte, error) {
	if m.eigenda != nil {
		m.log.Debug("Storing data to EigenDA backend")
		return m.eigenda.Put(ctx, value)
	}

	return nil, errors.New("no DA storage backend found")
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
