package signingrate

import (
	"context"
	"time"

	"github.com/Layr-Labs/eigenda/api/grpc/validator"
)

var _ SigningRateStorage = (*noOpSigningRateStorage)(nil)

// A no-op implementation of the SigningRateStorage interface, for unit tests.
type noOpSigningRateStorage struct {
}

// Create a new no-op SigningRateStorage.
func NewNoOpSigningRateStorage() SigningRateStorage {
	return &noOpSigningRateStorage{}
}

func (n *noOpSigningRateStorage) StoreBuckets(ctx context.Context, buckets []*validator.SigningRateBucket) error {
	return nil
}

func (n *noOpSigningRateStorage) LoadBuckets(
	ctx context.Context,
	startTimestamp time.Time,
) ([]*validator.SigningRateBucket, error) {
	return make([]*validator.SigningRateBucket, 0), nil
}
