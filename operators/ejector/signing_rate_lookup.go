package ejector

import (
	"time"

	"github.com/Layr-Labs/eigenda/api/grpc/validator"
)

// A tool for looking up signing rates for validators.
type SigningRateLookup interface {
	// GetSigningRateInfo returns signing rate information for all validators over the given time span.
	GetSigningRateInfo(timeSpan time.Duration) ([]*validator.ValidatorSigningRate, error)
}
