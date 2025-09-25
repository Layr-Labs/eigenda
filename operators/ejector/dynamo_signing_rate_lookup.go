package ejector

import (
	"time"

	"github.com/Layr-Labs/eigenda/api/grpc/validator"
)

var _ SigningRateLookup = (*dynamoSigningRateLookup)(nil)

// Uses batch information in dynamoDB to determine signing rates.
type dynamoSigningRateLookup struct{}

func (srl *dynamoSigningRateLookup) GetSigningRateInfo(
	timeSpan time.Duration,
) ([]*validator.ValidatorSigningRate, error) {
	//TODO implement me
	panic("implement me")
}
