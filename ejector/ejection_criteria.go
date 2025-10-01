package ejector

import (
	"github.com/Layr-Labs/eigenda/api/grpc/validator"
)

// Checks to see if a validator is eligible for ejection. A validator is considered to be eligible for ejection if
// 1. It has not signed for any batches in the time period (i.e. signed_batches == 0), and
// 2. There were batches it COULD have signed (i.e. unsigned_batches > 0)
//
// Note that the SLAs for validators requiest a much higher standard for ejection (i.e. 99.5% signing rate over a day).
// But due to constraints in the ejection protocol, there is little point in attempting to eject a validator that
// is currently online (even if it isn't doing a good job).
func IsEjectable(signingRate *validator.ValidatorSigningRate) bool {
	return signingRate.GetSignedBatches() == 0 && signingRate.GetUnsignedBatches() > 0
}
