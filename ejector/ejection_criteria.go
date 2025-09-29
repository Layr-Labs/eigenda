package ejector

import (
	"github.com/Layr-Labs/eigenda/api/grpc/validator"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigensdk-go/logging"
)

// Checks to see if a validator is eligible for ejection. A threshold of 1.0 means that the validator must
// sign for 100% of the time to avoid ejection. A threshold of 0.9 means that the validator must sign for
// 90% of the time to avoid ejection, and so on. The threshold must be in the range (0, 1].
func IsEjectable(
	logger logging.Logger,
	signingRate *validator.ValidatorSigningRate,
	threshold float64,
) bool {

	eject := false

	validatorID := core.OperatorID(signingRate.GetValidatorId())

	totalBatches := signingRate.GetSignedBatches() + signingRate.GetUnsignedBatches()
	if totalBatches > 0 {
		signedFraction := float64(signingRate.GetSignedBatches()) / float64(totalBatches)
		if signedFraction < threshold {
			logger.Infof(
				"Validator %s is eligible for ejeciton: signed batch fraction %.4f is below threshold %.4f",
				validatorID.Hex(), signedFraction, threshold)
			eject = true
		}
	}

	totalBytes := signingRate.GetSignedBytes() + signingRate.GetUnsignedBytes()
	if totalBytes > 0 {
		signedFraction := float64(signingRate.GetSignedBytes()) / float64(totalBytes)
		if signedFraction < threshold {
			logger.Infof(
				"Validator %s is eligible for ejeciton: signed byte fraction %.4f is below threshold %.4f",
				validatorID.Hex(), signedFraction, threshold)
			eject = true
		}
	}

	return eject
}
