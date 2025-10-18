package ejector

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/Layr-Labs/eigenda/api/grpc/validator"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/eth"
	"github.com/Layr-Labs/eigenda/core/eth/operatorstate"
	"github.com/Layr-Labs/eigensdk-go/logging"
)

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

	// The frequency with which to evaluate validators for ejection.
	period time.Duration

	// Defines the time window over which to evaluate signing metrics when deciding whether to eject a validator.
	ejectionCriteriaTimeWindow time.Duration

	// Used to convert validator IDs to validator addresses.
	validatorIDToAddressCache eth.ValidatorIDToAddressConverter

	// Used to look up the latest reference number.
	referenceBlockProvider eth.ReferenceBlockProvider

	// Used to fetch operator state.
	operatorStateCache operatorstate.OperatorStateCache
}

// NewEjector creates a new Ejector.
func NewEjector(
	ctx context.Context,
	logger logging.Logger,
	config *EjectorConfig,
	ejectionManager *ThreadedEjectionManager,
	signingRateLookupV1 SigningRateLookup,
	signingRateLookupV2 SigningRateLookup,
	validatorIDToAddressCache eth.ValidatorIDToAddressConverter,
	referenceBlockProvider eth.ReferenceBlockProvider,
	operatorStateCache operatorstate.OperatorStateCache,
) *Ejector {
	e := &Ejector{
		ctx:                        ctx,
		logger:                     logger,
		ejectionManager:            ejectionManager,
		signingRateLookupV1:        signingRateLookupV1,
		signingRateLookupV2:        signingRateLookupV2,
		period:                     config.EjectionPeriod,
		ejectionCriteriaTimeWindow: config.EjectionCriteriaTimeWindow,
		validatorIDToAddressCache:  validatorIDToAddressCache,
		referenceBlockProvider:     referenceBlockProvider,
		operatorStateCache:         operatorStateCache,
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
			err := e.evaluateValidators()
			if err != nil {
				e.logger.Error("error evaluating validators", "error", err)
			}
		}
	}
}

// evaluateValidators looks up signing rates and evaluates which validators should be ejected.
func (e *Ejector) evaluateValidators() error {

	v1SigningRates, err := e.signingRateLookupV1.GetSigningRates(
		e.ejectionCriteriaTimeWindow,
		nil, // all quorums
		ProtocolVersionV1,
		true, // omit perfect signers if possible (data API has inconsistent behavior across v1 and v2)
	)
	if err != nil {
		return fmt.Errorf("error looking up v1 signing rates: %w", err)
	}

	v2SigningRates, err := e.signingRateLookupV2.GetSigningRates(
		e.ejectionCriteriaTimeWindow,
		nil, // all quorums
		ProtocolVersionV2,
		true, // omit perfect signers if possible (data API has inconsistent behavior across v1 and v2)
	)
	if err != nil {
		return fmt.Errorf("error looking up v2 signing rates: %w", err)
	}

	// Combine data from v1 and v2 lookups, since the validator is likely to cancel ejection if it is active in either.
	signingRates, err := combineSigningRateSlices(v1SigningRates, v2SigningRates)
	if err != nil {
		return fmt.Errorf("error combining signing rates: %w", err)
	}
	sortByUnsignedBytesDescending(signingRates)

	for _, signingRate := range signingRates {
		err := e.evaluateValidator(signingRate)
		if err != nil {
			e.logger.Error("error evaluating validator", "validatorID", signingRate.GetValidatorId(), "error", err)
		}
	}

	return nil
}

// evaluateValidator evaluates a single validator's signing rate and decides whether to eject it.
func (e *Ejector) evaluateValidator(signingRate *validator.ValidatorSigningRate) error {
	isEjectable := signingRate.GetSignedBatches() == 0 && signingRate.GetUnsignedBatches() > 0
	if !isEjectable {
		return nil
	}

	if len(signingRate.GetValidatorId()) != 32 {
		return fmt.Errorf("invalid validator ID length: %d", len(signingRate.GetValidatorId()))
	}

	validatorID := core.OperatorID(signingRate.GetValidatorId()[:])
	validatorAddress, err := e.validatorIDToAddressCache.ValidatorIDToAddress(e.ctx, validatorID)
	if err != nil {
		return fmt.Errorf("error converting validator ID to address: %w", err)
	}

	e.logger.Info("Validator is eligible for ejection",
		"validatorID", signingRate.GetValidatorId(),
		"validatorAddress", validatorAddress.Hex(),
		"signedBatches", signingRate.GetSignedBatches(),
		"unsignedBatches", signingRate.GetUnsignedBatches(),
		"signedBytes", signingRate.GetSignedBytes(),
		"unsignedBytes", signingRate.GetUnsignedBytes(),
	)

	rbn, err := e.referenceBlockProvider.GetReferenceBlockNumber(e.ctx)
	if err != nil {
		return fmt.Errorf("error looking up latest reference block number: %w", err)
	}

	operatorState, err := e.operatorStateCache.GetOperatorState(
		e.ctx,
		rbn,
		nil, // all quorums
	)
	if err != nil {
		return fmt.Errorf("error looking up operator state: %w", err)
	}

	stakeFractions, err := getStakeFractionMap(validatorID, operatorState)
	if err != nil {
		return fmt.Errorf("error calculating stake fractions: %w", err)
	}

	// The ejection manager is responsible for deduplicating ejection requests, and deciding if
	// there are other factors that may prevent ejection (e.g. too many ejection attempts, etc.).
	err = e.ejectionManager.EjectValidator(validatorAddress, stakeFractions)
	if err != nil {
		return fmt.Errorf("error requesting ejection: %w", err)
	}

	return nil
}

// Get the stake fraction map for a given validator.
func getStakeFractionMap(
	validatorID core.OperatorID,
	operatorState *core.OperatorState,
) (map[core.QuorumID]float64, error) {
	stakeFractions := make(map[core.QuorumID]float64)

	for quorumID, operators := range operatorState.Operators {
		quorumStake := big.NewInt(0)
		for _, operatorInfo := range operators {
			quorumStake.Add(quorumStake, operatorInfo.Stake)
		}

		if quorumStake.Cmp(big.NewInt(0)) == 0 {
			// Ignore quorums with zero total stake to avoid division by zero
			continue
		}

		validatorInfo, ok := operators[validatorID]
		if !ok {
			// Validator is not part of this quorum.
			continue
		}

		stakeFraction := new(big.Rat).SetFrac(validatorInfo.Stake, quorumStake)
		floatStakeFraction, _ := stakeFraction.Float64()
		stakeFractions[quorumID] = floatStakeFraction
	}

	return stakeFractions, nil
}
