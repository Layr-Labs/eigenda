package ejector

import (
	"time"

	"github.com/Layr-Labs/eigenda/api/grpc/validator"
)

var _ SigningRateLookup = (*controllerSigningRateLookup)(nil)

// Looks up signing rates by asking the controller.
type controllerSigningRateLookup struct {
	// This is a placeholder. Will be implemented once the controller exposes an API for signing rates.
}

func (srl *controllerSigningRateLookup) GetSigningRateInfo(
	timeSpan time.Duration,
) ([]*validator.ValidatorSigningRate, error) {
	// TODO placeholder
	return nil, nil
}
