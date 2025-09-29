package ejector

import (
	"context"
	"time"

	"github.com/Layr-Labs/eigenda/api/grpc/validator"
	"github.com/Layr-Labs/eigenda/core/eth"
	"github.com/Layr-Labs/eigensdk-go/logging"
)

// TODO this should be the top level logic, should periodically loop over validators and decide which ones to eject

// Ejector is responsible for periodically evaluating validators and deciding which ones to eject.
type Ejector struct {
	ctx    context.Context
	logger logging.Logger

	// Responsible for executing ejections and managing the ejection lifecycle.
	ejectionManager *ThreadedEjectionManager

	// Used for looking up signing rates for V1.
	// TODO(cody.littley): remove this after V1 sunset
	signingRateLookupV1 SigningRateLookup

	// Used for looking up signing rates for V2.
	signingRateLookupV2 SigningRateLookup

	// The freqeuency with which to evaluate validators for ejection.
	period time.Duration

	// The time window over which to evaluate signing rates. Validators that have not signed any batches
	// in this time window are considered eligible for ejection.
	ejectionTimeWindow time.Duration
}

// NewEjector creates a new Ejector.
func NewEjector(
	ctx context.Context,
	logger logging.Logger,
	ejectionManager *ThreadedEjectionManager,
	signingRateLookupV1 SigningRateLookup,
	signingRateLookupV2 SigningRateLookup,
	period time.Duration,
	ejectionTimeWindow time.Duration,
) *Ejector {
	e := &Ejector{
		ctx:                 ctx,
		logger:              logger,
		ejectionManager:     ejectionManager,
		signingRateLookupV1: signingRateLookupV1,
		signingRateLookupV2: signingRateLookupV2,
	}

	go e.mainLoop()

	return e
}

// The main loop periodically evaluates validators for ejection.
func (e *Ejector) mainLoop() {
	ticker := time.NewTicker(e.period)
	defer ticker.Stop()

	for {
		select {
		case <-e.ctx.Done():
			e.logger.Info("ejector shutting down")
			return
		case <-ticker.C:
			e.evaluateValidators()
		}
	}
}

// evaluateValidators looks up signing rates and evaluates which validators should be ejected.
func (e *Ejector) evaluateValidators() {

	v1SigningRates, err := e.signingRateLookupV1.GetSigningRates(
		e.ejectionTimeWindow,
		nil, // all quorums
		ProtocolVersionV1,
		true, // omit perfect signers if possible
	)
	if err != nil {
		e.logger.Error("error looking up v1 signing rates", "error", err)
		return
	}

	v2SigningRates, err := e.signingRateLookupV2.GetSigningRates(
		e.ejectionTimeWindow,
		nil, // all quorums
		ProtocolVersionV2,
		true, // omit perfect signers if possible
	)
	if err != nil {
		e.logger.Error("error looking up v2 signing rates", "error", err)
		return
	}

	// Combine data from v1 and v2 lookups, since the validator is likely to cancel ejection if it is active in either.
	signingRates := combineSigningRateSlices(v1SigningRates, v2SigningRates)
	sortByUnsignedBytesDescending(signingRates)

	for _, signingRate := range signingRates {
		e.evaluateValidator(signingRate)
	}
}

// evaluateValidator evaluates a single validator's signing rate and decides whether to eject it.
func (e *Ejector) evaluateValidator(signingRate *validator.ValidatorSigningRate) {
	ejectable := IsEjectable(signingRate)
	if !ejectable {
		return
	}

	validatorAddress, err := eth.ValidatorIDToAddress()

	// TODO we need to get the mapping from validator ID to operator ID
	// TODO log signing rates

	e.logger.Info("Validator is eligible for ejection",
		"validatorID", signingRate.GetValidatorId())

	// e.ejectionManager.EjectValidator(core.OperatorID(signingRate.GetValidatorId()))
}
