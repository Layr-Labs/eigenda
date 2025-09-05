package signingrate

import (
	"time"
)

var _ SigningRateStorage = (*dynamoSigningRateStorage)(nil)

// A DynamoDB implementation of the SigningRateStorage interface.
type dynamoSigningRateStorage struct {
}

// Instantiate a new DynamoDB-backed SigningRateStorage.
func NewDynamoSigningRateStorage() (SigningRateStorage, error) {
	// TODO
	return &dynamoSigningRateStorage{}, nil
}

func (d *dynamoSigningRateStorage) StoreBuckets(buckets []*SigningRateBucket) error {
	//TODO implement me
	panic("implement me")
}

func (d *dynamoSigningRateStorage) LoadBuckets(startTimestamp time.Time) ([]*SigningRateBucket, error) {
	//TODO implement me
	panic("implement me")
}
