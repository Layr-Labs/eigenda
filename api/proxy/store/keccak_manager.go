package store

import (
	"context"
	"errors"
	"fmt"

	_ "github.com/Layr-Labs/eigenda/api/clients/v2"
	"github.com/Layr-Labs/eigenda/api/proxy/store/secondary/s3"
	"github.com/Layr-Labs/eigensdk-go/logging"
)

//go:generate mockgen -package mocks --destination ../test/mocks/keccak_manager.go . IKeccakManager

// IKeccakManager handles keccak/S3 operations
type IKeccakManager interface {
	// See [KeccakManager.PutOPKeccakPairInS3]
	PutOPKeccakPairInS3(ctx context.Context, key []byte, value []byte) error
	// See [KeccakManager.GetOPKeccakValueFromS3]
	GetOPKeccakValueFromS3(ctx context.Context, key []byte) ([]byte, error)
}

// KeccakManager handles keccak/S3 operations
type KeccakManager struct {
	log logging.Logger
	s3  *s3.Store // for op keccak256 commitment
}

var _ IKeccakManager = &KeccakManager{}

// NewKeccakManager creates a new KeccakManager
// s3 is optional, but if nil, the PutOPKeccakPairInS3 and GetOPKeccakValueFromS3 methods will return errors.
func NewKeccakManager(s3 *s3.Store, l logging.Logger) (*KeccakManager, error) {
	return &KeccakManager{
		log: l,
		s3:  s3,
	}, nil
}

// PutOPKeccakPairInS3 puts a key/value pair, where key=keccak(value), into S3.
// If key!=keccak(value), a Keccak256KeyValueMismatchError is returned.
// This is only used for OP keccak256 commitments.
func (m *KeccakManager) PutOPKeccakPairInS3(ctx context.Context, key []byte, value []byte) error {
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
func (m *KeccakManager) GetOPKeccakValueFromS3(ctx context.Context, key []byte) ([]byte, error) {
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
