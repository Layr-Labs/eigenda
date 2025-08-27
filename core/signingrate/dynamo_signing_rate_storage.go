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

func (d *dynamoSigningRateStorage) StoreBucket(buckets []Bucket) error {
	//TODO implement me
	panic("implement me")
}

func (d *dynamoSigningRateStorage) LoadBuckets(startTimestamp time.Time) ([]Bucket, error) {
	//TODO implement me
	panic("implement me")
}
