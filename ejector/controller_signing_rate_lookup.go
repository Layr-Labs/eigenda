package ejector

import (
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda/api/grpc/validator"
)

var _ SigningRateLookup = (*controllerSigningRateLookup)(nil)

// Looks up signing rates by asking the controller.
type controllerSigningRateLookup struct {
	// This is a placeholder. Will be implemented once the controller exposes an API for fetching signing rates.
}

func (srl *controllerSigningRateLookup) GetSigningRates(
	timeSpan time.Duration,
	version ProtocolVersion,
	omitPerfectSigners bool,
) ([]*validator.ValidatorSigningRate, error) {
	if version != ProtocolVersionV2 {
		return nil, fmt.Errorf("controller signing rate lookup only supports protocol version v2")
	}

	// TODO placeholder
	return nil, nil
}
